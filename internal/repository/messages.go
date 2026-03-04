package repository

import (
	"anonbot/internal/database"
	"time"
)

func CanSendMessage(userID int64) (bool, int) {

	query := `
	SELECT created_at
	FROM messages
	WHERE from_user = ?
	ORDER BY created_at DESC
	LIMIT 1
	`

	row := database.DB.QueryRow(query, userID)

	var createdAt string

	err := row.Scan(&createdAt)

	if err != nil {
		return true, 0
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return true, 0
	}

	diff := time.Since(t)

	if diff < 10*time.Second {

		remain := 10 - int(diff.Seconds())

		return false, remain
	}

	return true, 0
}

func SaveMessage(from int64, to int64, text string) int64 {

	query := `
	INSERT INTO messages(from_user, to_user, text)
	VALUES (?, ?, ?)
	`

	res, _ := database.DB.Exec(query, from, to, text)

	id, _ := res.LastInsertId()

	return id
}

func GetMessageSender(messageID int64) (int64, bool) {

	query := `
	SELECT from_user
	FROM messages
	WHERE id = ?
	`

	row := database.DB.QueryRow(query, messageID)

	var sender int64

	err := row.Scan(&sender)

	if err != nil {
		return 0, false
	}

	return sender, true
}