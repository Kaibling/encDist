package tokenizer

import (
	"github.com/kaibling/encDist/libs"
	log "github.com/sirupsen/logrus"
	"github.com/rs/xid"
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
)

type Tokenizer struct {
	userBuffer map[string] *libs.FullUser
	configuration *libs.Configuration
}

//NewUserstore dff
func NewTokenizer(configuration *libs.Configuration ) *Tokenizer{

	returnUserStore := new(Tokenizer)
	returnUserStore.userBuffer = make(map[string] *libs.FullUser)
	returnUserStore.configuration = configuration
	log.Debug("Loaded Config")
	log.Debug(returnUserStore.configuration)
	libs.SQLiteInitDB(returnUserStore.configuration.DBpath)
	return returnUserStore
}

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

func (Tokenizer *Tokenizer) StartServer() {
	http.HandleFunc("/token", Tokenizer.tokenHandler)
	log.Info("server started on Port " + Tokenizer.configuration.BindingPort )
	http.ListenAndServe(":"+Tokenizer.configuration.BindingPort, nil)
}


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