package libs

import (
	log "github.com/sirupsen/logrus"
	"crypto/rand"
	"crypto/rsa"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"crypto/sha1"
	"encoding/hex"
	"errors"
)


//CryptoData sd 
type CryptoData struct {
	Data []byte
	Keys map[string][]byte
}

func (CryptoData *CryptoData) EncryptData(data []byte, publickey rsa.PublicKey,username string) {

	//todo: generate password
	aeskey := make([]byte, 32)
	rand.Read(aeskey)
	encryptedData, err := AESencryptData(aeskey,data)
	if err != nil {
		log.Debug(err)
		return
	}
	CryptoData.Data = encryptedData
	CryptoData.Keys = make(map[string][]byte)
	encryptedKey := RSAEncrypt(aeskey,&publickey)
	CryptoData.Keys[username] = encryptedKey

}

func (CryptoData *CryptoData) showKeys() {

	for key,val := range CryptoData.Keys {
		log.Printf("%s -> %v",key ,val )
	}

}

func (CryptoData *CryptoData) DecryptData(username string, privKey rsa.PrivateKey) ([]byte,error) {

	aeskey,err := RSADecrypt(CryptoData.Keys[username],&privKey)
	if err != nil {
		return []byte(""),err
	}
	plain,err := AESdecryptdata(aeskey,CryptoData.Data)
	if err != nil {
		return []byte(""),err
	}
	return plain,nil
}

func (CryptoData *CryptoData) grantAccess(username string, privKey rsa.PrivateKey, username2 string, pubKey rsa.PublicKey) {

	aeskey,err := RSADecrypt(CryptoData.Keys[username],&privKey)
	if err != nil {
		return
	}
	encryptedKey := RSAEncrypt(aeskey,&pubKey)
	CryptoData.Keys[username2] = encryptedKey
}



//AESencryptData asds
func AESencryptData(encryptionKey []byte, data []byte) ([]byte, error) {
	if len(encryptionKey) == 0 {
		//log.Debugf("AES no Encryption Key provided \n")
		return []byte(""),errors.New("AES no Encryption Key provided")
	}

	if len(encryptionKey) != 32 {
		encryptionKey = padKey(encryptionKey)
	}
    c, err := aes.NewCipher(encryptionKey)
    if err != nil {
        log.Println(err)
    }

    gcm, err := cipher.NewGCM(c)
    if err != nil {
        log.Println(err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        log.Println(err)
    }
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func padKey( key []byte) []byte {

	for i := len(key); i <  32; i++ {
			key = append(key,key[i%len(key)])
	}
	return key
} 

//AESdecryptdata ds
func AESdecryptdata(decryptionKey []byte,ciphertext []byte) ([]byte, error) {

	if len(decryptionKey) != 32 {
		decryptionKey = padKey(decryptionKey)
	}
    c, err := aes.NewCipher(decryptionKey)
    if err != nil {
        log.Println(err)
    }

    gcm, err := cipher.NewGCM(c)
    if err != nil {
        log.Println(err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        log.Println(err)
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        log.Println(err)
    }
    return plaintext,nil
}


//RSAEncrypt a
func RSAEncrypt(data []byte, publickey *rsa.PublicKey  ) []byte {

	label := []byte("")
	sha1Hash := sha1.New()
	log.Println("Encrypting Data ....")
	encryptedmsg, err := rsa.EncryptOAEP(sha1Hash, rand.Reader, publickey, data, label)
	if err != nil {
		log.Println(err)
	}
	log.Println("Encryption finished")
	return encryptedmsg
}


//GenerateRSAKeyPair sd
func GenerateRSAKeyPair() *rsa.PrivateKey {
	// generate private key
	log.Println("generating new Key...")
	returnPrivatekey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Printf("Generation error: %s",err)
	}

	returnPrivatekey.Precompute()
	err = returnPrivatekey.Validate()
	if err != nil {
		log.Println(err)
	}
	log.Println("Key generation finished")
	return returnPrivatekey
}

//RSADecrypt sda
func  RSADecrypt(cipherText []byte, privateKey *rsa.PrivateKey) ([]byte,error) {
	label := []byte("")
	sha1Hash := sha1.New()
	log.Println("RSA Decrypting Data ....")
	decryptedmsg, err := rsa.DecryptOAEP(sha1Hash, rand.Reader, privateKey, cipherText, label)
	if err != nil {
		return []byte(""),err
	}
	log.Println("RSA Decryption finished")

	return decryptedmsg, nil
}


func SHA1HashString(data []byte) string {
	h := sha1.New()

	h.Write(data)
	 bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
