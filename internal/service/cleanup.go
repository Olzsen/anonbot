package service

import (
	"anonbot/internal/database"
	"log"
	"time"
)

func StartCleanup() {

	go func() {

		for {

			time.Sleep(24 * time.Hour)

			query := `
			DELETE FROM messages
			WHERE created_at < datetime('now', '-7 days')
			`

			_, err := database.DB.Exec(query)

			if err != nil {
				log.Println("cleanup error:", err)
			} else {
				log.Println("old messages deleted")
			}

		}

	}()
}