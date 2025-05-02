package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateId() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
