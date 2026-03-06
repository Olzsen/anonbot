package service

import (
	"anonbot/internal/database"
	"log"
	"time"
)

func StartCleanup() {

	go func() {

		for {

			_, err := database.DB.Exec(
				`DELETE FROM messages
				 WHERE created_at < datetime('now','-30 day')`,
			)

			if err != nil {
				log.Println("cleanup error:", err)
			}

			time.Sleep(12 * time.Hour)
		}

	}()
}