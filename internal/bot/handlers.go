package bot

import (
	"anonbot/internal/repository"
	"anonbot/internal/service"
	"bytes"
	"fmt"
	"html"
	"strconv"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
	tb "gopkg.in/telebot.v4"
)

var add string = "<a href='https://t.me/rolivikstarsbot?start=ref_1611018178'>Звезды по курсу 1⭐ - 1.42₽</a>"

const maxMessageLength = 1000

func StartHandler(botUsername string) func(tb.Context) error {

	return func(c tb.Context) error {

		user := c.Sender()

		repository.CreateUser(user.ID, user.Username)

		ref := repository.GetRefCode(user.ID)

		if ref == "" {

			ref = service.GenerateRef()

			repository.SetRefCode(user.ID, ref)
		}

		link := fmt.Sprintf(
			"https://t.me/%s?start=%s",
			botUsername,
			ref,
		)

		args := c.Args()

		if len(args) == 0 {

			msg := fmt.Sprintf(
				"📨 <b>Получай анонимные сообщения</b>\n\n"+
					"1️⃣ Поделись своей ссылкой\n"+
					"2️⃣ Люди будут писать тебе анонимно\n"+
					"3️⃣ Ты сможешь отвечать\n\n"+
					"🔗 <code>%s</code>\n\n%s",
				link,
				add,
			)

			btnShare := tb.InlineButton{
				Text: "📤 Поделиться ссылкой",
				URL:  link,
			}

			btnQR := tb.InlineButton{
				Text:   "📷 QR код",
				Unique: "qr",
			}

			btnStats := tb.InlineButton{
				Text:   "📊 Моя статистика",
				Unique: "stats",
			}

			btnHelp := tb.InlineButton{
				Text:   "ℹ️ Как это работает",
				Unique: "help",
			}

			markup := &tb.ReplyMarkup{
				InlineKeyboard: [][]tb.InlineButton{
					{btnShare},
					{btnQR},
					{btnStats, btnHelp},
				},
			}

			return c.Send(msg, markup)
		}

		targetID := repository.GetUserByRef(args[0])

		if targetID == 0 {
			return c.Send("❌ Неверная ссылка")
		}

		if targetID == user.ID {

			return c.Send(
				fmt.Sprintf(
					"Это твоя ссылка 🙂\n\n🔗 <code>%s</code>\n\n%s",
					link,
					add,
				),
			)
		}

		repository.SetSession(user.ID, targetID)

		return c.Send("✉️ Напиши сообщение\n\nОно будет отправлено анонимно")
	}
}

func TextHandler(c tb.Context) error {

	user := c.Sender()
	text := c.Text()

	if len(text) > maxMessageLength {
		return c.Send("❌ Сообщение слишком длинное (макс 1000 символов)")
	}

	replyID, replying := repository.GetReply(user.ID)

	if replying {

		safe := html.EscapeString(text)

		msg := fmt.Sprintf(
			"💬 <b>Ответ на анонимное сообщение</b>\n\n<blockquote><code>%s</code></blockquote>\n\n%s",
			safe,
			add,
		)

		service.Queue <- service.Job{
			UserID: replyID,
			Text:   msg,
		}

		repository.DeleteReply(user.ID)

		return c.Send("✅ Ответ отправлен")
	}

	targetID, ok := repository.GetSession(user.ID)

	if !ok {
		return c.Send("Открой ссылку пользователя чтобы отправить сообщение")
	}

	allow, wait := service.Allow(user.ID)

	if !allow {
		return c.Send(
			fmt.Sprintf(
				"⏳ Подождите %d секунд перед следующим сообщением",
				wait,
			),
		)
	}

	safe := html.EscapeString(text)

	messageID := repository.SaveMessage(user.ID, targetID, safe)

	msg := fmt.Sprintf(
		"📩 <b>Анонимное сообщение</b>\n\n<blockquote><code>%s</code></blockquote>\n\n%s",
		safe,
		add,
	)

	btn := tb.InlineButton{
		Text: "💬 Ответить",
		Data: fmt.Sprintf("reply:%d", messageID),
	}

	markup := &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{
			{btn},
		},
	}

	service.Queue <- service.Job{
		UserID: targetID,
		Text:   msg,
		Markup: markup,
	}

	repository.DeleteSession(user.ID)

	return c.Send("✅ Сообщение отправлено анонимно")
}

func PhotoHandler(c tb.Context) error {

	user := c.Sender()
	photo := c.Message().Photo

	replyID, replying := repository.GetReply(user.ID)

	if replying {

		service.Queue <- service.Job{
			UserID: replyID,
			Photo:  photo.FileID,
		}

		repository.DeleteReply(user.ID)

		return c.Send("✅ Ответ отправлен")
	}

	targetID, ok := repository.GetSession(user.ID)

	if !ok {
		return c.Send("Открой ссылку пользователя чтобы отправить сообщение")
	}

	messageID := repository.SaveMediaMessage(
		user.ID,
		targetID,
		"photo",
		photo.FileID,
	)

	btn := tb.InlineButton{
		Text: "💬 Ответить",
		Data: fmt.Sprintf("reply:%d", messageID),
	}

	markup := &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{
			{btn},
		},
	}

	service.Queue <- service.Job{
		UserID: targetID,
		Photo:  photo.FileID,
		Markup: markup,
	}

	repository.DeleteSession(user.ID)

	return c.Send("✅ Фото отправлено анонимно")
}

func VideoHandler(c tb.Context) error {

	user := c.Sender()
	video := c.Message().Video

	replyID, replying := repository.GetReply(user.ID)

	if replying {

		service.Queue <- service.Job{
			UserID: replyID,
			Video:  video.FileID,
		}

		repository.DeleteReply(user.ID)

		return c.Send("✅ Ответ отправлен")
	}

	targetID, ok := repository.GetSession(user.ID)

	if !ok {
		return c.Send("Открой ссылку пользователя чтобы отправить сообщение")
	}

	messageID := repository.SaveMediaMessage(
		user.ID,
		targetID,
		"video",
		video.FileID,
	)

	btn := tb.InlineButton{
		Text: "💬 Ответить",
		Data: fmt.Sprintf("reply:%d", messageID),
	}

	markup := &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{
			{btn},
		},
	}

	service.Queue <- service.Job{
		UserID: targetID,
		Video:  video.FileID,
		Markup: markup,
	}

	repository.DeleteSession(user.ID)

	return c.Send("✅ Видео отправлено анонимно")
}

func VoiceHandler(c tb.Context) error {

	user := c.Sender()
	voice := c.Message().Voice

	replyID, replying := repository.GetReply(user.ID)

	if replying {

		service.Queue <- service.Job{
			UserID: replyID,
			Voice:  voice.FileID,
		}

		repository.DeleteReply(user.ID)

		return c.Send("✅ Ответ отправлен")
	}

	targetID, ok := repository.GetSession(user.ID)

	if !ok {
		return c.Send("Открой ссылку пользователя чтобы отправить сообщение")
	}

	messageID := repository.SaveMediaMessage(
		user.ID,
		targetID,
		"voice",
		voice.FileID,
	)

	btn := tb.InlineButton{
		Text: "💬 Ответить",
		Data: fmt.Sprintf("reply:%d", messageID),
	}

	markup := &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{
			{btn},
		},
	}

	service.Queue <- service.Job{
		UserID: targetID,
		Voice:  voice.FileID,
		Markup: markup,
	}

	repository.DeleteSession(user.ID)

	return c.Send("✅ Голосовое отправлено анонимно")
}

func ReplyButton(c tb.Context) error {

	data := c.Callback().Data

	if !strings.HasPrefix(data, "reply:") {
		return nil
	}

	idStr := strings.TrimPrefix(data, "reply:")

	messageID, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		return nil
	}

	senderID, ok := repository.GetMessageSender(messageID)

	if !ok {
		c.Respond()
		return c.Send("❌ Сообщение не найдено")
	}

	repository.SetReply(c.Sender().ID, senderID)

	c.Respond()

	return c.Send("✏️ Напиши ответ на сообщение")
}

func StatsHandler(c tb.Context) error {

	user := c.Sender()

	received := repository.CountReceived(user.ID)
	sent := repository.CountSent(user.ID)
	today := repository.CountToday(user.ID)

	msg := fmt.Sprintf(
		"📊 <b>Твоя статистика</b>\n\n"+
			"📨 Получено сообщений: <b>%d</b>\n"+
			"📤 Отправлено сообщений: <b>%d</b>\n"+
			"📅 Сегодня получено: <b>%d</b>\n\n%s",
		received,
		sent,
		today,
		add,
	)

	return c.Send(msg)
}

func HelpHandler(c tb.Context) error {

	msg := "ℹ️ <b>Как работает бот</b>\n\n" +
		"1. Поделись своей ссылкой\n" +
		"2. Люди будут писать тебе\n" +
		"3. Ты получишь анонимные сообщения\n" +
		"4. Можно отвечать прямо из бота"

	return c.Send(fmt.Sprintf("%s\n\n%s", msg, add))
}

func QRHandler(botUsername string) func(tb.Context) error {

	return func(c tb.Context) error {

		user := c.Sender()

		ref := repository.GetRefCode(user.ID)

		link := fmt.Sprintf(
			"https://t.me/%s?start=%s",
			botUsername,
			ref,
		)

		qr, err := qrcode.Encode(link, qrcode.Medium, 256)

		if err != nil {
			return c.Send("Ошибка генерации QR")
		}

		file := &tb.Photo{
			File:    tb.FromReader(bytes.NewReader(qr)),
			Caption: "📷 QR код для получения анонимных сообщений",
		}

		return c.Send(file)
	}
}