package service

import (
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenerateRef() string {

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 8)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}