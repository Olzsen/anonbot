package bot

import (
	"anonbot/internal/service"
	"fmt"

	tb "gopkg.in/telebot.v4"
)

func RateLimitMiddleware(next tb.HandlerFunc) tb.HandlerFunc {

	return func(c tb.Context) error {

		user := c.Sender()

		if user == nil {
			return next(c)
		}

		allow, wait := service.Allow(user.ID)

		if !allow {

			msg := fmt.Sprintf(
				"⏳ Подождите %d секунд перед следующим сообщением",
				wait,
			)

			return c.Send(msg)
		}

		return next(c)
	}
}

func AntiSpamMiddleware(next tb.HandlerFunc) tb.HandlerFunc {

	return func(c tb.Context) error {

		user := c.Sender()

		ok, wait := service.CheckSpam(user.ID)

		if !ok {

			hours := wait / 3600

			return c.Send(
				fmt.Sprintf(
					"🚫 Вы временно заблокированы за спам\n\nПопробуйте снова через %d часов",
					hours,
				),
			)
		}

		return next(c)
	}
}