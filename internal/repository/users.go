package repository

import "anonbot/internal/database"

func CreateUser(id int64, username string) {

	query := `
	INSERT OR IGNORE INTO users(telegram_id, username)
	VALUES (?, ?)`

	database.DB.Exec(query, id, username)
}