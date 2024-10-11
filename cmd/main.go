package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/somuthink/sirius_weather_bot/internal"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)

	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	state := "idle"

	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID
		if !update.Message.IsCommand() && state == "idle" {
			msg.Text = "sry, i work only with commands"
		}

		if state == "city input" {
			city, err := internal.CheckCityExists(update.Message.Text)
			if err == nil {
				msg.Text = fmt.Sprintf("you wanna choose %s, right?", city)
				state = "idle"

			} else if err == internal.ErrNotExistingCity {
				msg.Text = "it seems like there isn`t any city with name like that"
			} else {
				msg.Text = "there was an error in bot`s work"
				log.Println(err)
			}

		}

		if update.Message.Command() == "start" {
			state = "city input"
			msg.Text = "write the name of the city using english letter in a new message"
		}
		bot.Send(msg)
	}
}
