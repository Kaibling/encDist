package main

import ("log"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)


func encryptData(encryptionKey []byte, data []byte) []byte {

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
	return gcm.Seal(nonce, nonce, data, nil)
}

func padKey( key []byte) []byte {

	for i := len(key); i <  32; i++ {
			key = append(key,key[i%len(key)])
	}
	return key
} 

func decryptdata(decryptionKey []byte,ciphertext []byte) []byte {
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
    return plaintext
}
