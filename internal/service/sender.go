package service

import (
	"log"
	"time"

	tb "gopkg.in/telebot.v4"
)

func StartSender(bot *tb.Bot) {

	for i := 0; i < 5; i++ {

		go worker(bot)
	}
}

func worker(bot *tb.Bot) {

	for job := range Queue {

		user := &tb.User{ID: job.UserID}

		if job.Photo != "" {

			_, err := bot.Send(user, &tb.Photo{
				File:    tb.File{FileID: job.Photo},
				Caption: "📩 <b>Анонимное сообщение</b>",
			}, job.Markup)

			if err != nil {
				log.Println(err)
			}

			continue
		}

		if job.Video != "" {
			bot.Send(user, &tb.Video{
				File: tb.File{FileID: job.Video},
			}, job.Markup)
			continue
		}

		if job.Voice != "" {
			bot.Send(user, &tb.Voice{
				File: tb.File{FileID: job.Voice},
			}, job.Markup)
			continue
		}

		bot.Send(user, job.Text, job.Markup)

		time.Sleep(200 * time.Millisecond)
	}
}