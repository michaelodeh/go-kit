package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func Hash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func Compare(value string, hash string) bool {
	return Hash(value) == hash
}

func GenerateRandomHash() string {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		panic("failed to generate secure random bytes")
	}

	return hex.EncodeToString(b)
}
