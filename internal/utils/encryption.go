package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
)

var AESKeyLength = 16

func generateCryptoKey(n int) ([]byte, error) {
	data := make([]byte, n)
	_, err := rand.Read(data)
	if err != nil {
		return nil, fmt.Errorf("genCryptoKey: %w", err)
	}
	return data, nil
}

func generateAESKey() ([]byte, error) {
	return generateCryptoKey(AESKeyLength)
}

func GetEncryptedMessage(rsaPubKey *rsa.PublicKey, data []byte) ([]byte, string, error) {
	aesKey, err := generateAESKey()
	if err != nil {
		return nil, "", fmt.Errorf("AES gen failed, %w", err)
	}

	aesMessage, err := aesEncrypt(aesKey, data)
	if err != nil {
		return nil, "", fmt.Errorf("aes encryption failed, %w", err)
	}

	encryptedKey, err := rsaEncrypt(rsaPubKey, aesKey)
	if err != nil {
		return nil, "", fmt.Errorf("rsa encryption failed, %w", err)
	}

	msg := make([]byte, base64.RawStdEncoding.EncodedLen(len(aesMessage)))
	base64.RawStdEncoding.Encode(msg, aesMessage)

	key := base64.RawStdEncoding.EncodeToString(encryptedKey)
	return msg, key, nil
}

func GetDecryptedMessage(rsaPrivateKey *rsa.PrivateKey, cypher []byte, aesEncryptedKey string) ([]byte, error) {
	encryptedMsg := make([]byte, base64.RawStdEncoding.DecodedLen(len(cypher)))
	_, err := base64.RawStdEncoding.Decode(encryptedMsg, cypher)
	if err != nil {
		return nil, err
	}

	encryptedKey, err := base64.RawStdEncoding.DecodeString(aesEncryptedKey)
	if err != nil {
		return nil, err
	}

	key, err := rsaDecrypt(rsaPrivateKey, encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("rsa decryption failed, %w", err)
	}

	msg, err := aesDecrypt(key, encryptedMsg)
	if err != nil {
		return nil, fmt.Errorf("aes decryption failed, %w", err)
	}

	return msg, nil
}
