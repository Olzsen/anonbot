package bot

import (
	"anonbot/internal/repository"
	"anonbot/internal/service"
	"fmt"
	"html"
	"strconv"
	"strings"

	tb "gopkg.in/telebot.v4"
)

var add string = "<a href='https://t.me/rolivikstarsbot?start=ref_1611018178'>Звезды по курсу 1⭐ - 1.42₽</a>"

const maxMessageLength = 1000

func StartHandler(botUsername string) func(tb.Context) error {

	return func(c tb.Context) error {

		user := c.Sender()

		ref := service.GenerateRef()

		repository.CreateUser(user.ID, user.Username, ref)

		userRef := repository.GetRefCode(user.ID)

		link := fmt.Sprintf(
			"https://t.me/%s?start=%s",
			botUsername,
			userRef,
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

		return c.Send(
			fmt.Sprintf("✉️ <b>Напиши сообщение</b>\n\nОно будет отправлено анонимно.\n\n%s", add),
		)
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

func ReplyButton(c tb.Context) error {

	data := c.Callback().Data

	if !strings.HasPrefix(data, "reply:") {
		return nil
	}

	idStr := strings.TrimPrefix(data, "reply:")
	messageID, _ := strconv.ParseInt(idStr, 10, 64)

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