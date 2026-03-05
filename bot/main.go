package main

import (
	"anonbot/internal/bot"
	"anonbot/internal/database"
	"anonbot/internal/service"
	"log"
	"os"
	"strconv"
	"time"

	tb "gopkg.in/telebot.v4"
)

var OwnerID int64

func main() {

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN not set")
	}

	owner := os.Getenv("OWNER_ID")
	if owner != "" {
		id, _ := strconv.ParseInt(owner, 10, 64)
		OwnerID = id
	}

	bot.OwnerID = OwnerID

	database.Init()

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

	username := b.Me.Username

	b.Handle("/start", bot.StartHandler(username))
	b.Handle("/stats", bot.StatsHandler)

	b.Handle("/admin", bot.AdminHandler)
	b.Handle("/setad", bot.SetAdHandler)

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