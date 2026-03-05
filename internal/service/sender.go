package service

import (
	"log"
	"time"

	tb "gopkg.in/telebot.v4"
)

func StartSender(bot *tb.Bot) {

	for job := range Queue {

		user := &tb.User{ID: job.UserID}

		_, err := bot.Send(user, job.Text, job.Markup)

		if err != nil {
			log.Println("send error:", err)
		}

		time.Sleep(300 * time.Millisecond)
	}
}