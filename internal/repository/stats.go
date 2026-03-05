package repository

import "anonbot/internal/database"

func CountReceived(userID int64) int {

	var count int

	database.DB.QueryRow(
		`SELECT COUNT(*) FROM messages WHERE to_user=?`,
		userID,
	).Scan(&count)

	return count
}

func CountSent(userID int64) int {

	var count int

	database.DB.QueryRow(
		`SELECT COUNT(*) FROM messages WHERE from_user=?`,
		userID,
	).Scan(&count)

	return count
}

func CountToday(userID int64) int {

	var count int

	database.DB.QueryRow(
		`SELECT COUNT(*) FROM messages 
		WHERE to_user=? 
		AND date(created_at)=date('now')`,
		userID,
	).Scan(&count)

	return count
}