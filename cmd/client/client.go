package main

import (
	"net/http"
	"encoding/json"
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/kaibling/encDist/libs"
		"io/ioutil"
    "net/url"
	"os"
	  "bufio"
  "fmt"
  "strings"
  "syscall"
  "os/exec"
  "errors"
	
)
type ClientConfig struct {
    TokenizerIP string 
    Username string
    userpassword string
    clientToken string
    SavedIdentifier []string
}

func main() {
	log.SetLevel(log.DebugLevel)
    clientConfig := ParseConfigFile("config.json")
	startConsole(clientConfig)

}
func startConsole(clientConfig *ClientConfig) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("  EncDist ")
	fmt.Printf("user: %s\n",clientConfig.Username)
	fmt.Println("------------")

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "q" || text == "quit" {
			return
		}
		commandArray := strings.Split(text," ")

		switch commandArray[0] {
			case "help":
				help()
			case "ls":
				getIdentifier(clientConfig)
			case "p": fallthrough
			case "password":
				menusetPassword(clientConfig)
			case "e": fallthrough
			case "encrypt":
				menuEncryptData(clientConfig,commandArray)
			case "d": fallthrough
			case "decrypt":
				menuDecryptData(clientConfig,commandArray)
			case "token": case "t":
				fmt.Println(getToken(clientConfig))
			case "createuser": 
				menuCreateUser(clientConfig,commandArray)
			default:
				fmt.Println("unknown command")
		}
	
	}

}

func menusetPassword(clientConfig *ClientConfig) {

	stty, _ := exec.LookPath("stty")
	sttyArgs := syscall.ProcAttr{
		"",
		[]string{},
		[]uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		nil,
	}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Password: ")
	// Disable echo.
	if stty != "" {
		syscall.ForkExec(stty, []string{"stty", "-echo"}, &sttyArgs)
	}

	// Get password.
	pass, _ := reader.ReadString('\n')
	pass = strings.Replace(pass, "\n", "", -1)
	fmt.Print("\n")

	// Enable echo.
	if stty != "" {
		syscall.ForkExec(stty, []string{"stty", "echo"}, &sttyArgs)
	}

	clientConfig.userpassword = pass 
	saveConfig("config.json",clientConfig)
	fmt.Printf("Password set\n")
}

func menuEncryptData(clientConfig *ClientConfig, commandArray []string) {
	if len(commandArray) < 2  {
		fmt.Println("command unknown")
		return
	}

	plainText := strings.Join(commandArray[1:], " ")
	identifier, err := encryptData(clientConfig,[]byte(plainText))
	if err != nil {
		fmt.Println(err)
	}
	
	clientConfig.SavedIdentifier = append(clientConfig.SavedIdentifier,identifier)
	saveConfig("config.json",clientConfig)
}

func menuCreateUser(clientConfig *ClientConfig, commandArray []string) {
	if len(commandArray) != 3  {
		fmt.Println("missing arguments")
		return
	}

	connectionString := "http://127.0.0.1:8070/user"
	postData := url.Values{}
	postData.Add("username", commandArray[1])
	postData.Add("password", commandArray[2])
	postData.Add("command", "CREATE")
	resp, err := http.Post(connectionString, "application/x-www-form-urlencoded; param=value", bytes.NewBufferString(postData.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if string(data) == "OK" {
		clientConfig.Username = commandArray[1]
		saveConfig("config.json",clientConfig)
	} else {
		fmt.Printf("User creation failed\n")
	}


}

func menuDecryptData(clientConfig *ClientConfig, commandArray []string) {

	if len(commandArray) < 2  {
		fmt.Println("command unknown")
		return
	}
	for _,val := range commandArray[1:] {
		plaindata,err := decryptData(clientConfig,val)
		if err != nil {
		log.Debug(err)
		}
		fmt.Println(plaindata)
	}
}



func getToken(clientConfig *ClientConfig) string {
	connectionString := clientConfig.TokenizerIP+"/token"
	if clientConfig.Username == "" {
		fmt.Printf("username not set")
		return ""
	}
	if clientConfig.userpassword == "" {
		fmt.Printf("password not set")
		return ""
	}

	sendUser := &libs.User{Name: clientConfig.Username,Password: clientConfig.userpassword}
	bytesRepresentation, err := json.Marshal(sendUser)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(connectionString, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if string(data) == "" {
		fmt.Println("no token recieved. User/Password invalid")
		return ""
	} else {
		log.Debugf("Token: %s",string(data))
		fmt.Println("Token recieved")
		clientConfig.clientToken = string(data)
    	return string(data)
	}
	
}

func encryptData(clientConfig *ClientConfig,data []byte) (string, error) {

	connectionString := clientConfig.TokenizerIP+"/encrypt"
	if clientConfig.clientToken == "" {
		return "", errors.New("no token provided")
	}
	responseData := &libs.PlainDataTransfer{Data: data,Token: clientConfig.clientToken}
	bytesRepresentation, err := json.Marshal(responseData)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(connectionString, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}
    returnData, err := ioutil.ReadAll(resp.Body)
    identifier := string(returnData)
    log.Debugf(" Identifier: %s",identifier)
    return identifier, nil
}

func decryptData(clientConfig *ClientConfig,identifier string) (string,error) {

	connectionString := clientConfig.TokenizerIP+"/decrypt"
	if clientConfig.clientToken == "" {
		fmt.Printf("no token provided\n")
		return "", errors.New("no token provided")
	}
	postData := url.Values{}
	postData.Add("token", clientConfig.clientToken)
	postData.Add("hash", identifier)
	resp, err := http.Post(connectionString, "application/x-www-form-urlencoded; param=value", bytes.NewBufferString(postData.Encode()))
	if err != nil {
		log.Warn("no proper response from server")
		log.Warn(err)
		return "", errors.New("no proper response from server")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if string(data) == "" {
		 return "", errors.New("no encryption possible")
	}
    return string(data), nil
}


func ParseConfigFile(filepath string) *ClientConfig {

	//default path found
	//create config file
	_, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		log.Println("file does not exists")

		fo, err := os.Create("config.json")
		libs.CheckErr(err)

		returnConfig := new(ClientConfig)
		returnConfig.TokenizerIP = "http://127.0.0.1:8070"
		returnConfig.Username = ""
		configString, err := json.Marshal(returnConfig)
		libs.CheckErr(err)

		_, err = fo.Write(configString)
		libs.CheckErr(err)

		defer fo.Close()
        log.Println("Configuration file created")
        return returnConfig
    } else {

        log.Println("opening configuration file: " + filepath)
        returnConfig := new(ClientConfig)
        file, err := os.Open(filepath)
        libs.CheckErr(err)
        decoder := json.NewDecoder(file)
        err = decoder.Decode(&returnConfig)
        libs.CheckErr(err)
        return returnConfig
    }
    
}


func saveConfig(configPath string, clientConfig *ClientConfig)  {

        fo, err := os.Create("config.json")
		configString, err := json.Marshal(clientConfig)
		libs.CheckErr(err)
		_, err = fo.Write(configString)
		libs.CheckErr(err)
		defer fo.Close()
		log.Debugf("Configuration file saved")
}

func getIdentifier(clientConfig *ClientConfig) {
    for _,ident := range clientConfig.SavedIdentifier {
        log.Println(ident)
    }
}

func help() {
	fmt.Println("does things")
	fmt.Println("ls				list identifier")
	fmt.Println("encrypt / e	encrpyt data")
	fmt.Println("get token / t  retrive token")
}