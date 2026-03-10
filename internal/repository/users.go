package repository

import "anonbot/internal/database"

func CreateUser(id int64, username string) {

	_, _ = database.DB.Exec(
		`INSERT OR IGNORE INTO users(id, telegram_id, username) VALUES(?,?,?)`,
		id,
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

func GetUserByRef(code string) int64 {

	row := database.DB.QueryRow(`
	SELECT telegram_id
	FROM users
	WHERE ref_code = ?
	`, code)

	var id int64

	err := row.Scan(&id)

	if err != nil {
		return 0
	}

	return id
}