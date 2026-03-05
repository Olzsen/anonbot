package service

import (
	"crypto/rand"
	"math/big"
)

const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenerateRef() string {

	length := 8

	result := make([]byte, length)

	for i := range result {

		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[n.Int64()]

	}

	return string(result)
}