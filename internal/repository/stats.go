package repository

import "anonbot/internal/database"

func CountMessages(userID int64) int {

	query := `
	SELECT COUNT(*)
	FROM messages
	WHERE to_user = ?
	`

	row := database.DB.QueryRow(query, userID)

	var count int

	row.Scan(&count)

	return count
}