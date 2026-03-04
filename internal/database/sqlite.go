package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {

	var err error

	DB, err = sql.Open("sqlite3", "./data/bot.db")

	if err != nil {
		log.Fatal(err)
	}

	// оптимизация соединений
	DB.SetMaxOpenConns(1)

	optimize()

	createTables()
}

func optimize() {

	pragmas := []string{

		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA temp_store = MEMORY;",
		"PRAGMA mmap_size = 30000000000;",
	}

	for _, pragma := range pragmas {

		_, err := DB.Exec(pragma)

		if err != nil {
			log.Println("pragma error:", err)
		}
	}
}

func createTables() {

	createUsers := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		telegram_id INTEGER UNIQUE,
		username TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createMessages := `
	CREATE TABLE IF NOT EXISTS messages(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_user INTEGER,
		to_user INTEGER,
		text TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createSessions := `
	CREATE TABLE IF NOT EXISTS sessions(
		user_id INTEGER PRIMARY KEY,
		target_user INTEGER
	);`

	createReplies := `
	CREATE TABLE IF NOT EXISTS replies(
		user_id INTEGER PRIMARY KEY,
		target_user INTEGER
	);`

	DB.Exec(createUsers)
	DB.Exec(createMessages)
	DB.Exec(createSessions)
	DB.Exec(createReplies)
}