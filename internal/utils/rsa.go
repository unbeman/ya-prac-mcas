package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

var ErrNoRSAKey = errors.New("no key provided")

func GetPublicKey(path string) (*rsa.PublicKey, error) {
	if len(path) == 0 {
		return nil, ErrNoRSAKey
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("GetPublicKey: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("GetPublicKey: failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("GetPublicKey: %w", err)
	}
	return pub.(*rsa.PublicKey), nil
}

func GetPrivateKey(path string) (*rsa.PrivateKey, error) {
	if len(path) == 0 {
		return nil, ErrNoRSAKey
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("GetPrivateKey: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("GetPrivateKey: failed to decode PEM block containing private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("GetPrivateKey: %w", err)
	}

	return priv, nil
}

func rsaEncrypt(key *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, key, data, nil)
}

func rsaDecrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, key, data, nil)
}
