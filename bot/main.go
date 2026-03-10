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

func main() {

	token := os.Getenv("BOT_TOKEN")
	owner := os.Getenv("OWNER_ID")

	if token == "" {
		log.Fatal("BOT_TOKEN not set")
	}

	ownerID, _ := strconv.ParseInt(owner, 10, 64)
	bot.OwnerID = ownerID

	database.Init()

	pref := tb.Settings{
		Token:     token,
		ParseMode: tb.ModeHTML,
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	}

	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	btnStats := tb.InlineButton{
		Text:   "📊 Моя статистика",
		Unique: "stats",
	}

	btnHelp := tb.InlineButton{
		Text:   "ℹ️ Как это работает",
		Unique: "help",
	}

	btnQR := tb.InlineButton{
		Text:   "📷 QR код",
		Unique: "qr",
	}

	b.Handle("/start", bot.StartHandler(b.Me.Username))
	b.Handle("/admin", bot.AdminHandler)
	b.Handle("/setad", bot.SetAdHandler)

	b.Handle(&btnStats, bot.StatsHandler)
	b.Handle(&btnHelp, bot.HelpHandler)
	b.Handle(&btnQR, bot.QRHandler(b.Me.Username))

	b.Handle(tb.OnText, bot.TextHandler)
	b.Handle(tb.OnPhoto, bot.PhotoHandler)
	b.Handle(tb.OnVideo, bot.VideoHandler)
	b.Handle(tb.OnVoice, bot.VoiceHandler)

	b.Handle(tb.OnCallback, bot.ReplyButton)

	service.StartCleanup()

	service.Queue = make(chan service.Job, 100)
	go service.StartSender(b)

	log.Println("Bot started")

	b.Start()
}