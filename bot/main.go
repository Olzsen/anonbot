package main

import (
	"anonbot/internal/bot"
	"anonbot/internal/database"
	"anonbot/internal/service"
	"log"
	"os"
	"time"

	tb "gopkg.in/telebot.v4"
)

func main() {

	database.Init()

	token := os.Getenv("BOT_TOKEN")

	if token == "" {
		log.Fatal("BOT_TOKEN not set")
	}

	pref := tb.Settings{
		Token: token,
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
		ParseMode: tb.ModeHTML,
	}

	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot initialized")

	service.StartWorker(b)

	service.StartCleanup()

	botUsername := b.Me.Username

	b.Handle("/start", bot.StartHandler(botUsername))
	b.Handle("/stats", bot.StatsHandler)

	b.Handle(
		tb.OnText,
		bot.AntiSpamMiddleware(
			bot.RateLimitMiddleware(
				bot.TextHandler,
			),
		),
	)

	b.Handle(tb.OnCallback, bot.ReplyButton)

	b.Handle(&tb.InlineButton{Unique: "help"}, bot.HelpHandler)
	b.Handle(&tb.InlineButton{Unique: "stats"}, bot.StatsHandler)

	log.Println("Bot started")

	b.Start()
}