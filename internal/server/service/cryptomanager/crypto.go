package cryptomanager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// NewAESGCM creates a new AES-GCM cipher.
func NewAESGCM(key []byte) (cipher.AEAD, error) {
	aesgcm, err := generateAesgcm(key)
	if err != nil {
		return nil, fmt.Errorf("failed create new aes gcm: %w", err)
	}
	return aesgcm, nil
}

// Encrypt encrypts data using AES-GCM.
func Encrypt(aesgcm cipher.AEAD, data, nonce []byte) []byte {
	return aesgcm.Seal(nil, nonce, data, nil)
}

// Decrypt decrypts data using AES-GCM.
func Decrypt(aesgcm cipher.AEAD, data, nonce []byte) ([]byte, error) {
	result, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed decrypt data: %w", err)
	}

	return result, nil
}

func generateAesgcm(key []byte) (cipher.AEAD, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed create new aes block: %w", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("failed create new aes gcm: %w", err)
	}

	return aesgcm, nil
}

func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed generate random bytes: %w", err)
	}

	return b, nil
}
