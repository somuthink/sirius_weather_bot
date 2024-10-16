package pkg

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/somuthink/sirius_weather_bot/internal/db"
	"github.com/somuthink/sirius_weather_bot/internal/weather"
)

func CurrentWeather(bot *tgbotapi.BotAPI, chatId int64) error {
	city, err := db.SelectUserCity(chatId)
	if err != nil {
		return err
	}

	weather_resp, err := weather.WeatherRequest(city)
	if err != nil {
		return err
	}

	text := fmt.Sprintf("~~ %s%s ~~\n \ncurrent temperature in *%s* is `%.2f°C` \n\nto change frequency of messages type /choose", weather_resp.Current.Condition.Text, weather_resp.GetConditionEmoji(), city, weather_resp.Current.Temp_c)

	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "Markdown"

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
