package service

import (
	"sync"
	"time"
)

type limiter struct {
	last map[int64]time.Time
	mu   sync.Mutex
}

var RateLimiter = limiter{
	last: map[int64]time.Time{},
}

func Allow(userID int64) (bool, int) {

	RateLimiter.mu.Lock()
	defer RateLimiter.mu.Unlock()

	now := time.Now()

	last, ok := RateLimiter.last[userID]

	if !ok {
		RateLimiter.last[userID] = now
		return true, 0
	}

	diff := now.Sub(last)

	if diff < 10*time.Second {

		wait := int((10*time.Second - diff).Seconds())

		return false, wait
	}

	RateLimiter.last[userID] = now

	return true, 0
}