package service

import tb "gopkg.in/telebot.v4"

type Job struct {
	UserID int64
	Text   string

	Photo string
	Video string
	Voice string

	Markup *tb.ReplyMarkup
}

var Queue = make(chan Job, 100)