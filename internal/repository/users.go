package repository

import "anonbot/internal/database"

func CreateUser(id int64, username string, ref string) {

	_, _ = database.DB.Exec(
		`INSERT OR IGNORE INTO users(id, username, ref_code)
		 VALUES(?,?,?)`,
		id,
		username,
		ref,
	)
}

func GetRefCode(userID int64) string {

	var code string

	database.DB.QueryRow(
		`SELECT ref_code FROM users WHERE id=?`,
		userID,
	).Scan(&code)

	return code
}

func GetUserByRef(ref string) int64 {

	var id int64

	database.DB.QueryRow(
		`SELECT id FROM users WHERE ref_code=?`,
		ref,
	).Scan(&id)

	return id
}

func CountUsers() int {

	var count int

	database.DB.QueryRow(
		`SELECT COUNT(*) FROM users`,
	).Scan(&count)

	return count
}