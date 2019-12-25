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
	
)
type ClientConfig struct {
    TokenizerIP string 
    Username string
    userpassword string
    clientToken string
    SavedIdentifier []string
}

func main() {
    clientConfig := ParseConfigFile("config.json")
    clientConfig.TokenizerIP = "http://127.0.0.1:8070"
    clientConfig.Username = "hans"
    clientConfig.userpassword = "hanspwd"
	saveConfig("config.json",clientConfig)
	
	startConsole(clientConfig)


/*
	clientConfig.clientToken = getToken(clientConfig)
    identifier := encryptData(clientConfig,[]byte("hansi"))
    clientConfig.SavedIdentifier = append(clientConfig.SavedIdentifier,identifier)
    saveConfig("config.json",clientConfig)
    plainText := decryptData(clientConfig,identifier)
	log.Print(plainText)
	*/

}
func startConsole(clientConfig *ClientConfig) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("  EncDist")
	fmt.Println("------------")

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "q" || text == "quit" {
			return
		}
		commandArray := strings.Split(text," ")
		if commandArray[0] == "d" || commandArray[0] =="decrypt" {
			if len(commandArray) < 2  {
				fmt.Println("command unknown")
				continue
			}
			for _,val := range commandArray[1:] {
				fmt.Println(decryptData(clientConfig,val))
			}
		} else {
		
			switch text {
				case "help":
					help()
				case "ls":
					getIdentifier(clientConfig)
				case "encrypt": case "e":
					menuEncryptData(clientConfig)
				case "get token": case "t":
					fmt.Println(getToken(clientConfig))
				default:
					fmt.Println("unknown command")
			}
		}
	
	}

}

func menuEncryptData(clientConfig *ClientConfig) {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("encrpyt> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "q" || text == "quit" {
			return
		}
		identifier:= encryptData(clientConfig,[]byte(text))
		clientConfig.SavedIdentifier = append(clientConfig.SavedIdentifier,identifier)
    	saveConfig("config.json",clientConfig)
	}
}


func getToken(clientConfig *ClientConfig) string {
	connectionString := clientConfig.TokenizerIP+"/token"
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
	log.Debugf("Token: %s",string(data))
	clientConfig.clientToken = string(data)
    return string(data)
}

func encryptData(clientConfig *ClientConfig,data []byte) string {

	connectionString := clientConfig.TokenizerIP+"/encrypt"
	if clientConfig.clientToken == "" {
		log.Errorf("no token provided")
		return ""
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
    return identifier
}

func decryptData(clientConfig *ClientConfig,identifier string) string {

	connectionString := clientConfig.TokenizerIP+"/decrypt"
	if clientConfig.clientToken == "" {
		log.Errorf("no token provided")
		return ""
	}
	postData := url.Values{}
	postData.Add("token", clientConfig.clientToken)
	postData.Add("hash", identifier)
	resp, err := http.Post(connectionString, "application/x-www-form-urlencoded; param=value", bytes.NewBufferString(postData.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
    return string(data)
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