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
	dbFilePath string
}

//NewUserstore dff
func NewTokenizer(dbFilePath string ) *Tokenizer{

	returnUserStore := new(Tokenizer)
	returnUserStore.userBuffer = make(map[string] *libs.FullUser)
	returnUserStore.dbFilePath = dbFilePath
	libs.SQLiteInitDB(dbFilePath)
	return returnUserStore
}

func (Tokenizer *Tokenizer) NewUser(name string, password string) {

	//check if user already exists
	checkuser := libs.SQLiteGetUser(Tokenizer.dbFilePath , name)
	if checkuser != nil {
		log.Printf("User does already exist: " + checkuser.Name)
		return
	}
	newUser := new(libs.User)
	newUser.Name = name
	newUser.Password = password
	libs.SQLiteaddUser(Tokenizer.dbFilePath , newUser,libs.GenerateRSAKeyPair())
}

func (Tokenizer *Tokenizer) GetToken(name string, password string) string {

	//check if user already exists
	
	checkuser := libs.SQLiteGetFullUser(Tokenizer.dbFilePath , name,password)
	if checkuser == nil {
		log.Printf("User/Password does not match")
		return ""
	}
	guid := xid.New()
	return guid.String()
}

func (Tokenizer *Tokenizer) StartServer() {
	http.HandleFunc("/token", Tokenizer.tokenHandler)
	log.Info("server started on Port 8070" )
	http.ListenAndServe(":8070", nil)
}


func (Tokenizer *Tokenizer) tokenHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Body)

	var responseUser libs.User
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	json.Unmarshal([]byte(buf.String()), &responseUser)
	log.Debug(responseUser)
	token := Tokenizer.GetToken(responseUser.Name,responseUser.Password)
	log.Debugf("User %s got Token %s",responseUser.Name,token)
	fmt.Fprintf(w,token)
}