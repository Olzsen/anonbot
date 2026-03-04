package service

import (
	tb "gopkg.in/telebot.v4"
)

type Job struct {
	UserID int64
	Text   string
}

var Queue = make(chan Job, 1000)

func StartWorker(bot *tb.Bot) {

	go func() {

		for job := range Queue {

			bot.Send(&tb.User{ID: job.UserID}, job.Text)

		}

	}()
}