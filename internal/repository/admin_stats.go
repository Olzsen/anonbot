package repository

import "anonbot/internal/database"

func CountMessages() int {

	var count int

	database.DB.QueryRow(
		`SELECT COUNT(*) FROM messages`,
	).Scan(&count)

	return count
}

func CountMessagesToday() int {

	var count int

	database.DB.QueryRow(
		`SELECT COUNT(*) FROM messages 
		WHERE date(created_at)=date('now')`,
	).Scan(&count)

	return count
}

func CountActiveToday() int {

	var count int

	database.DB.QueryRow(
		`SELECT COUNT(DISTINCT from_user) FROM messages 
		WHERE date(created_at)=date('now')`,
	).Scan(&count)

	return count
}

func CountUsers() int {

	var count int

	database.DB.QueryRow(
		`SELECT COUNT(*) FROM users`,
	).Scan(&count)

	return count
}