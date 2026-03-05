package repository

var sessions = map[int64]int64{}

func SetSession(userID int64, targetID int64) {
	sessions[userID] = targetID
}

func GetSession(userID int64) (int64, bool) {
	target, ok := sessions[userID]
	return target, ok
}

func DeleteSession(userID int64) {
	delete(sessions, userID)
}