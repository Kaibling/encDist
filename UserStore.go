
package main

import 	(
	"crypto/rsa"
	log "github.com/sirupsen/logrus"
)


//User sds
type User struct {
	name string
	id int32
	password string
}

//UserStore sds
type UserStore struct {
	users map[string] *User
	dbFilePath string
}

//NewUserstore dff
func NewUserstore(dbFilePath string ) *UserStore{

	returnUserStore := new(UserStore)
	returnUserStore.users = make(map[string] *User)
	returnUserStore.dbFilePath = dbFilePath
	SQLiteInitDB(dbFilePath)
	return returnUserStore
}

func (UserStore *UserStore) newUser(name string, password string) {

	//check if user already exists
	checkuser := SQLiteGetUser(UserStore.dbFilePath , name)
	if checkuser != nil {
		log.Printf("User does already exist: " + checkuser.name)
		return
	}
	newUser := new(User)
	newUser.name = name
	newUser.password = password

	SQLiteaddUser(UserStore.dbFilePath , newUser,GenerateRSAKeyPair())

}

func (UserStore *UserStore) getKeyPair( username string , password string ) *rsa.PrivateKey {
	returnKeyPair,err :=  SQLiteGetKeyPair(UserStore.dbFilePath , username,password)
	if err != nil {
		log.Println(err)
	}
	return returnKeyPair
}
