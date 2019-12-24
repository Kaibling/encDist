package libs


import (
	log "github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"crypto/rsa"
	"encoding/json"
	"os"
	"flag"

)

type Configuration struct {
	BindingIPAddress string
	BindingPort      string
	DBpath	         string
}

type PlainDataTransfer struct {
	Token string
	Data []byte
}

type CryptoDataTransfer struct {
	Token string
	CryptoData CryptoData
}

func ParseArguments() map[string]string {

	arguments := make(map[string]string)
	configString := flag.String("config", "config.json", "filepath to configuration file")
	flag.Parse()

	arguments["configFilePath"] = *configString
	log.Print("load command line arguments ")
	log.Print(arguments)

	return arguments
}

func ParseConfigurationFile(filepath string) *Configuration {

	if filepath == "config.json" {
		//default path found
		//create config file
		_, err := os.Stat(filepath)

		if os.IsNotExist(err) {
			log.Println("file does not exists")

			fo, err := os.Create("config.json")
			CheckErr(err)

			returnConfig := new(Configuration)
			returnConfig.BindingIPAddress = "127.0.0.1"
			returnConfig.BindingPort = "8070"
			returnConfig.DBpath = "./sqlitehier.db"

			configString, err := json.Marshal(returnConfig)
			CheckErr(err)

			_, err = fo.Write(configString)
			CheckErr(err)

			defer fo.Close()
			log.Println("Configuration file created")
			return returnConfig
		}

	}
	log.Println("opening configuration file: " + filepath)
	returnConfig := new(Configuration)
	file, err := os.Open(filepath)
	CheckErr(err)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&returnConfig)
	CheckErr(err)

	return returnConfig
}

//User sds
type FullUser struct {
	Name string
	Password string
	PrivateKey *rsa.PrivateKey
}

//User sds
type User struct {
	Name string
	id int32
	Password string
}


func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}
