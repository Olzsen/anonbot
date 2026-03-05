package repository

import "anonbot/internal/database"

func CreateUser(id int64, username string) {

	_, _ = database.DB.Exec(
		`INSERT OR IGNORE INTO users (id, username) VALUES (?, ?)`,
		id,
		username,
	)
}

func CountUsers() int {

	var count int

	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM users`,
	).Scan(&count)

	if err != nil {
		return 0
	}

	return count
}