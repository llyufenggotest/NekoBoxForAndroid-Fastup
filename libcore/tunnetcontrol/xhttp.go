package tunnetcontrol

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const xhttpAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateXHTTPPath() (string, error) {
	lengthOffset, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		return "", err
	}
	length := 16 + int(lengthOffset.Int64())
	result := make([]byte, len("/api/v1/sync/")+length)
	copy(result, "/api/v1/sync/")
	for index := len("/api/v1/sync/"); index < len(result); index++ {
		character, err := rand.Int(rand.Reader, big.NewInt(int64(len(xhttpAlphabet))))
		if err != nil {
			return "", err
		}
		result[index] = xhttpAlphabet[character.Int64()]
	}
	if len(result) < 29 || len(result) > 38 {
		return "", errors.New("invalid TunNet XHTTP path length")
	}
	return string(result), nil
}
