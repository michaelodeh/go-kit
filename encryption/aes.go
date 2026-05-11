package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type Encrypter struct {
	key []byte
}

func NewEncrypter(key string) *Encrypter {
	hash := sha256.Sum256([]byte(key))
	return &Encrypter{key: hash[:]}
}

// Encrypt takes a plaintext string and returns an encrypted byte slice
func (a *Encrypter) Encrypt(plaintext string) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Returns nonce + ciphertext
	return gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

func (a *Encrypter) EncryptToString(plaintext string) (string, error) {
	encrypted, err := a.Encrypt(plaintext)
	if err != nil {
		return "", err
	}
	// Convert binary to safe base64 string
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (a *Encrypter) DecryptFromString(encryptedText string) (string, error) {
	// Decode from base64 back to binary
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}
	return a.Decrypt(ciphertext)
}

// Decrypt takes the ciphertext byte slice and returns the original string
func (a *Encrypter) Decrypt(ciphertext []byte) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, encryptedMessage := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}
