package service

import (
	"log"
	"time"

	tb "gopkg.in/telebot.v4"
)

type Job struct {
	UserID int64

	Text string

	Photo string
	Video string
	Voice string

	Markup interface{}
}

var Queue chan Job

func StartWorker(bot *tb.Bot) {

	go func() {

		for job := range Queue {

			send(bot, job)

			time.Sleep(50 * time.Millisecond)
		}

	}()
}

func send(bot *tb.Bot, job Job) {

	user := &tb.User{
		ID: job.UserID,
	}

	var err error

	if job.Markup != nil {

		_, err = bot.Send(user, job.Text, job.Markup)

	} else {

		_, err = bot.Send(user, job.Text)

	}

	if err != nil {

		log.Println("send error:", err)

		go func() {
			time.Sleep(2 * time.Second)
			Queue <- job
		}()

	}
}