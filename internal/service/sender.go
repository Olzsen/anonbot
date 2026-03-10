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

			photo := &tb.Photo{
				File: tb.File{FileID: job.Photo},
			}

			if job.Text != "" {
				photo.Caption = job.Text
			}

			_, err := bot.Send(user, photo, job.Markup)

			if err != nil {
				log.Println(err)
			}

			continue
		}

		if job.Video != "" {

			video := &tb.Video{
				File: tb.File{FileID: job.Video},
			}

			if job.Text != "" {
				video.Caption = job.Text
			}

			_, err := bot.Send(user, video, job.Markup)

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

		if job.Text != "" {

			_, err := bot.Send(user, job.Text, job.Markup)

			if err != nil {
				log.Println(err)
			}
		}

		time.Sleep(200 * time.Millisecond)
	}
}