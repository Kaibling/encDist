package main

import ("log"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)


func encryptData(encryptionKey []byte, data []byte) []byte {

	//text := bytes.Repeat([]byte("i"), 96)

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
	return gcm.Seal(nonce, nonce, data, nil)
}



func decryptdata(decryptionKey []byte,ciphertext []byte) []byte {
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
    return plaintext
}
