package repository

import (
	"anonbot/internal/database"
)

func SetSession(userID int64, targetID int64) {

	query := `
	INSERT OR REPLACE INTO sessions(user_id, target_user)
	VALUES (?, ?)`

	database.DB.Exec(query, userID, targetID)
}

func GetSession(userID int64) (int64, bool) {

	query := `
	SELECT target_user
	FROM sessions
	WHERE user_id = ?`

	row := database.DB.QueryRow(query, userID)

	var target int64

	err := row.Scan(&target)

	if err != nil {
		return 0, false
	}

	return target, true
}

func DeleteSession(userID int64) {

	query := `
	DELETE FROM sessions
	WHERE user_id = ?`

	database.DB.Exec(query, userID)
}