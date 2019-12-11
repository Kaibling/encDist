package main

import "log"
import "encoding/hex"
func main() {

	data := []byte("My Suff")
    key := []byte("passphrasewhichneedstobe32bytes!")


	encrData := encryptData(key,data)
	log.Println(hex.EncodeToString(encrData)  )
	plainText := decryptdata(key,encrData)
	log.Println(string(plainText))

}