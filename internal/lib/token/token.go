package token

import (
	"crypto/rand"
	"encoding/hex"
)

func NewToken() (string, error) {
	token := make([]byte, 256)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
