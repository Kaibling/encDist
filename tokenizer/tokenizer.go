package tokenizer

import (
	"github.com/kaibling/encDist/libs"
	log "github.com/sirupsen/logrus"
	"github.com/rs/xid"
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//Tokenizer sd
type Tokenizer struct {
	userBuffer map[string] *libs.FullUser
	configuration *libs.Configuration
}

//NewTokenizer dff
func NewTokenizer(configuration *libs.Configuration ) *Tokenizer{

	returnTokenizer := new(Tokenizer)
	returnTokenizer.userBuffer = make(map[string] *libs.FullUser)
	returnTokenizer.configuration = configuration
	log.Debug("Loaded Config")
	log.Debug(returnTokenizer.configuration)
	libs.SQLiteInitTknDB(returnTokenizer.configuration.DBpath)
	return returnTokenizer
}

// NewUser dds
func (Tokenizer *Tokenizer) NewUser(name string, password string) {

	//check if user already exists
	checkuser := libs.SQLiteGetUser(Tokenizer.configuration.DBpath , name)
	if checkuser != nil {
		log.Printf("User does already exist: " + checkuser.Name)
		return
	}
	newUser := new(libs.User)
	newUser.Name = name
	newUser.Password = password
	libs.SQLiteaddUser(Tokenizer.configuration.DBpath , newUser,libs.GenerateRSAKeyPair())
}

//GetToken sds
func (Tokenizer *Tokenizer) GetToken(name string, password string) string {

	//check if user already exists
	log.Debugf("check user %s in tokenBuffer",name)
	for key,val := range Tokenizer.userBuffer {
		if val.Name == name {
			log.Debugf("User found in tokenBuffer")
			return key
		}
	}

	
	checkuser := libs.SQLiteGetFullUser(Tokenizer.configuration.DBpath , name,password)
	if checkuser == nil {
		log.Printf("User/Password does not match")
		return ""
	}
	guid := xid.New()
	log.Debugf("User found in DB. Generating GUId %s",guid.String())
	Tokenizer.userBuffer[guid.String()] = checkuser
	return guid.String()
}

//StartServer sds
func (Tokenizer *Tokenizer) StartServer() {
	http.HandleFunc("/token", Tokenizer.tokenHandler)
	http.HandleFunc("/encrypt", Tokenizer.encryptHandler)
	log.Info("server started on Port " + Tokenizer.configuration.BindingPort )
	http.ListenAndServe(":"+Tokenizer.configuration.BindingPort, nil)
}

// tokenHandler sd
func (Tokenizer *Tokenizer) tokenHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug(r)
	var responseUser libs.User
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	json.Unmarshal([]byte(buf.String()), &responseUser)
	log.Debug(responseUser)

	token := Tokenizer.GetToken(responseUser.Name,responseUser.Password)
	log.Debugf("User %s got Token %s",responseUser.Name,token)
	fmt.Fprintf(w,token)
}

func (Tokenizer *Tokenizer) encryptHandler(w http.ResponseWriter, r *http.Request) {

	//get data and token
	log.Debug(r)
	var responseData libs.PlainDataTransfer
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	json.Unmarshal([]byte(buf.String()), &responseData)

 	//encrypt data
	decryptUser := Tokenizer.userBuffer[responseData.Token]
	cryptoData := new(libs.CryptoData)
	cryptoData.EncryptData(responseData.Data,decryptUser.PrivateKey.PublicKey,decryptUser.Name)

	//send data to publisher
	CryptoDataTransfer := new(libs.CryptoDataTransfer)
	CryptoDataTransfer.Token = responseData.Token
	CryptoDataTransfer.CryptoData = *cryptoData
	
	jsonBytes,err := json.Marshal(CryptoDataTransfer)
	if err != nil {
		log.Warn(err)
	}
	publisherServer := "http://127.0.0.1:8071/publish"
	
	resp, err := http.Post(publisherServer, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatalln(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	log.Println(string(data))

	fmt.Fprintf(w,"OK")
}