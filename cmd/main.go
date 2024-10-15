package main

import (
	"github.com/somuthink/sirius_weather_bot/internal/bot"
	"github.com/somuthink/sirius_weather_bot/internal/db"
	"github.com/somuthink/sirius_weather_bot/internal/sheduler"
)

func main() {
	db.Initialize()
	go sheduler.StartTickers()

	bot.HandleUpdates()

	// bot.Debug = true
}
