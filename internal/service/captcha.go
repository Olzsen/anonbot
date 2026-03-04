package service

import (
	"fmt"
	"math/rand"
	"sync"
)

var captcha = struct {
	data map[int64]int
	mu   sync.Mutex
}{
	data: map[int64]int{},
}

func Create(userID int64) (string, int) {

	a := rand.Intn(5) + 1
	b := rand.Intn(5) + 1

	answer := a + b

	captcha.mu.Lock()
	captcha.data[userID] = answer
	captcha.mu.Unlock()

	return fmt.Sprintf("%d + %d = ?", a, b), answer
}

func Check(userID int64, input int) bool {

	captcha.mu.Lock()
	defer captcha.mu.Unlock()

	ans, ok := captcha.data[userID]

	if !ok {
		return true
	}

	if ans == input {
		delete(captcha.data, userID)
		return true
	}

	return false
}