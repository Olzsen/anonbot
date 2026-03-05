package repository

import "anonbot/internal/database"

func CreateUser(id int64, username string) {

	_, _ = database.DB.Exec(
		`INSERT OR IGNORE INTO users(id, username) VALUES(?,?)`,
		id,
		username,
	)
}

func SetRefCode(userID int64, ref string) {

	_, _ = database.DB.Exec(
		`UPDATE users SET ref_code=? WHERE id=?`,
		ref,
		userID,
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