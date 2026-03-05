package bot

import (
	"anonbot/internal/repository"
	"fmt"
	"strings"

	tb "gopkg.in/telebot.v4"
)

var OwnerID int64

func AdminHandler(c tb.Context) error {

	if c.Sender().ID != OwnerID {
		return nil
	}

	users := repository.CountUsers()

	msg := fmt.Sprintf(
		"🛠 <b>Admin panel</b>\n\n"+
			"👤 Users: <b>%d</b>\n\n"+
			"/setad - изменить рекламу",
		users,
	)

	return c.Send(msg)
}

func SetAdHandler(c tb.Context) error {

	if c.Sender().ID != OwnerID {
		return nil
	}

	text := c.Text()

	ad := strings.TrimPrefix(text, "/setad ")

	if ad == "" {
		return c.Send("Использование:\n/setad <html реклама>")
	}

	add = ad

	return c.Send("✅ Реклама обновлена")
}