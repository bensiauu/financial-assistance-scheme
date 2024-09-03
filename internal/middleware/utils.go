package middleware

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecureKey(size int) (string, error) {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil
}
