package service

import (
	"sync"
	"time"
)

type userSpam struct {
	count int
	time  time.Time
}

type ban struct {
	until time.Time
}

var spamMap = struct {
	users map[int64]*userSpam
	bans  map[int64]*ban
	mu    sync.Mutex
}{
	users: map[int64]*userSpam{},
	bans:  map[int64]*ban{},
}

func CheckSpam(userID int64) (bool, int64) {

	spamMap.mu.Lock()
	defer spamMap.mu.Unlock()

	now := time.Now()

	// проверяем бан
	if b, ok := spamMap.bans[userID]; ok {

		if now.Before(b.until) {

			return false, int64(b.until.Sub(now).Seconds())
		}

		delete(spamMap.bans, userID)
	}

	user, ok := spamMap.users[userID]

	if !ok {

		spamMap.users[userID] = &userSpam{
			count: 1,
			time:  now,
		}

		return true, 0
	}

	if now.Sub(user.time) > time.Minute {

		user.count = 1
		user.time = now

		return true, 0
	}

	user.count++

	if user.count > 60 {

		spamMap.bans[userID] = &ban{
			until: now.Add(24 * time.Hour),
		}

		delete(spamMap.users, userID)

		return false, int64(24 * time.Hour.Seconds())
	}

	return true, 0
}