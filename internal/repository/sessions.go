package repository

import "anonbot/internal/database"

func SetSession(userID int64, targetID int64) {

	database.DB.Exec(`
	INSERT OR REPLACE INTO sessions(user_id,target_id,expires_at)
	VALUES(?,?,datetime('now','+5 minutes'))
	`,
		userID,
		targetID)
}

func GetSession(userID int64) (int64, bool) {

	row := database.DB.QueryRow(`
	SELECT target_id
	FROM sessions
	WHERE user_id=?
	AND expires_at > datetime('now')
	`,
		userID)

	var id int64

	err := row.Scan(&id)

	if err != nil {
		return 0, false
	}

	return id, true
}

func DeleteSession(userID int64) {

	database.DB.Exec(`
	DELETE FROM sessions
	WHERE user_id=?
	`, userID)
}