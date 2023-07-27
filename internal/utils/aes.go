package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	log "github.com/sirupsen/logrus"
)

func aesEncrypt(aesKey []byte, plaintext []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

func aesDecrypt(aesKey []byte, ciphertext []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		log.Error("NewCipher ", err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Error("NewGCM ", err)
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Error("gcm.Open ", err)
		return nil, err
	}

	return plaintext, nil
}
