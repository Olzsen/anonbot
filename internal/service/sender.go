package service

import (
	"log"
	"time"

	tb "gopkg.in/telebot.v4"
)

func StartSender(bot *tb.Bot) {

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

			_, err := bot.Send(user, &tb.Video{
				File:    tb.File{FileID: job.Video},
				Caption: "📩 <b>Анонимное сообщение</b>",
			}, job.Markup)

			if err != nil {
				log.Println(err)
			}

			continue
		}

		if job.Voice != "" {

			_, err := bot.Send(user, &tb.Voice{
				File: tb.File{FileID: job.Voice},
			}, job.Markup)

			if err != nil {
				log.Println(err)
			}

			continue
		}

		_, err := bot.Send(user, job.Text, job.Markup)

		if err != nil {
			log.Println(err)
		}

		time.Sleep(250 * time.Millisecond)
	}
}