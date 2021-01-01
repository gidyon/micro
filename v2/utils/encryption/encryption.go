package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// API is used for encryption and decryption
type API interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(cipher []byte) ([]byte, error)
}

type encryptionAPI struct {
	gcm cipher.AEAD
}

// NewAPI creates a symetric encryption/decryption API
func NewAPI(key []byte) (API, error) {
	gcm, err := createGCM(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcm: %v", err)
	}

	return &encryptionAPI{gcm}, nil
}

func (e *encryptionAPI) Encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return e.gcm.Seal(nonce, nonce, data, nil), nil
}

func (e *encryptionAPI) Decrypt(cipher []byte) ([]byte, error) {
	nonceSize := e.gcm.NonceSize()
	if len(cipher) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, cipher := cipher[:nonceSize], cipher[nonceSize:]
	return e.gcm.Open(nil, nonce, cipher, nil)
}

func createGCM(key []byte) (cipher.AEAD, error) {
	keyLen := len(key)
	switch {
	case keyLen > 32:
		key = key[:32]
	case keyLen > 24:
		key = key[:24]
	case keyLen > 16:
		key = key[:16]
	case keyLen < 16:
		return nil, errors.New("incorrect key: len must be 16, 24 or 32")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}

// ParseKey parses an encryption key truncating it to the nearest 16 or 24 or 32 bit. If key is less 16 returns error
func ParseKey(key []byte) ([]byte, error) {
	keyLen := len(key)
	switch {
	case keyLen < 16:
		return nil, errors.New("key length less that 16")
	case keyLen < 24:
		return key[:16], nil
	case keyLen < 32:
		return key[:24], nil
	default:
		return key[:32], nil
	}
}
