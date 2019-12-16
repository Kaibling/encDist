package main

import (log "github.com/sirupsen/logrus"
)


func main() {

	data := []byte("The s geht wordsuper see thingy")
	dbpath := "./sqlitehier.db"
	UserStore1 := NewUserstore(dbpath)
	UserStore1.newUser("hans","123123")
	UserStore1.newUser("fritz","654321")

	daten1 := new(CryptoData)
	daten1.encryptData(data,UserStore1.getKeyPair("hans","123123").PublicKey,"hans")
	daten1.grantAccess("hans",*UserStore1.getKeyPair("hans","123123"),"fritz",UserStore1.getKeyPair("fritz","654321").PublicKey)
	enc := daten1.decryptData("fritz",*UserStore1.getKeyPair("fritz","654321"))
	log.Info(string(enc))
}