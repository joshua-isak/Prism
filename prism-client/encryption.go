package main


import (
	"fmt"
	"crypto/aes"
    "crypto/cipher"
	"crypto/rand"
	"io"
)

// NOTE: Most of the below is magic that I do not understand

// AES encrypt bytes of data with a 32 byte key
func encrypt(data []byte, key []byte) []byte {
	// generate a new aes cipher using our 32 byte long key
    c, err := aes.NewCipher(key)
    if err != nil {
        fmt.Println("aes.NewCipher failed:", err)
	}

	// Galois/Counter Mode
	gcm, err := cipher.NewGCM(c)
    if err != nil {
        fmt.Println("cipher.NewGCM failed:", err)
    }

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        fmt.Println(err)
	}

	code := gcm.Seal(nonce, nonce, data, nil)

	return code
}


// AES decrypt bytes of data with a 32 byte key
func decrypt(cipherData []byte, key []byte) []byte {
    var errMsg string

	c, err := aes.NewCipher(key)
    if err != nil {
        errMsg += "aes.NewCipher: " + err.Error()
    }

    gcm, err := cipher.NewGCM(c)
    if err != nil {
        errMsg += "cipher.NewGCM: " + err.Error()
    }

    nonceSize := gcm.NonceSize()
    if len(cipherData) < nonceSize {
        errMsg += "cipherData is less than nonceSize"
    }

    nonce, ciphertext := cipherData[:nonceSize], cipherData[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        errMsg += "gcm.Open: " + err.Error()
    }

    // If an error occured, return that error as a byte array instead of the decrypted message
    if len(errMsg) > 0 {
        return []byte(errMsg)
    }

	return []byte(plaintext)
}