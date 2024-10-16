package sheduler

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/somuthink/sirius_weather_bot/internal/db"

	"github.com/somuthink/sirius_weather_bot/internal/pkg"
)

func sendWeather(bot *tgbotapi.BotAPI, time string) error {
	chatIds, err := db.SelectAllTimeUsers(time)
	if err != nil {
		return err
	}
	for _, chatId := range chatIds {
		if err = pkg.CurrentWeather(bot, chatId); err != nil {
			return err
		}
	}
	return nil
}

func StartTickers(bot *tgbotapi.BotAPI) {
	minuteTicker := time.NewTicker(60 * time.Second)
	defer minuteTicker.Stop()

	morningTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 8, 0, 0, 0, time.Local)
	if morningTime.Before(time.Now()) {
		morningTime = morningTime.Add(24 * time.Hour)
	}
	morningTicker := time.NewTicker(morningTime.Sub(time.Now()))
	defer morningTicker.Stop()

	afternoonTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 0, 0, 0, time.Local)
	if afternoonTime.Before(time.Now()) {
		afternoonTime = afternoonTime.Add(24 * time.Hour)
	}
	afternoonTicker := time.NewTicker(afternoonTime.Sub(time.Now()))
	defer afternoonTicker.Stop()

	eveningTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 18, 0, 0, 0, time.Local)
	if eveningTime.Before(time.Now()) {
		eveningTime = eveningTime.Add(24 * time.Hour)
	}
	eveningTicker := time.NewTicker(eveningTime.Sub(time.Now()))
	defer eveningTicker.Stop()
	var err error
	for {
		select {
		case <-minuteTicker.C:
			err = sendWeather(bot, "minute")
		case <-morningTicker.C:
			err = sendWeather(bot, "morning")
			morningTicker.Reset(24 * time.Hour)
		case <-afternoonTicker.C:
			err = sendWeather(bot, "afternoon")
			afternoonTicker.Reset(24 * time.Hour)
		case <-eveningTicker.C:
			err = sendWeather(bot, "evening")
			eveningTicker.Reset(24 * time.Hour)
		}

		if err != nil {
			log.Print(err)
		}
	}
}
