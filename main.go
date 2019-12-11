package main

import "log"
import "encoding/hex"
func main() {

	data := []byte("My Ssfsdfklsdjosljf lsdfs dfs sdf sdfuff")
	//key := []byte("passphrasewhichneedstobe32bytes!")
	key := []byte("Ksein")



	encrData := encryptData(key,data)
	log.Println(hex.EncodeToString(encrData)  )
	plainText := decryptdata(key,encrData)
	log.Println(string(plainText))


}