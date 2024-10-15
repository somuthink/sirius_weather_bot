package main

import (
	"github.com/somuthink/sirius_weather_bot/internal/bot"
	"github.com/somuthink/sirius_weather_bot/internal/db"
)

func main() {
	db.Initialize()

	bot.HandleUpdates()

	// bot.Debug = true
}
