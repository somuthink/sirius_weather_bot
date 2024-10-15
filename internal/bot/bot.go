package bot

import (
	// "fmt"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/somuthink/sirius_weather_bot/internal/sheduler"
)

var UserState map[int64]string

func HandleUpdates() {
	UserState = make(map[int64]string)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	go sheduler.StartTickers(bot)

	u := tgbotapi.NewUpdate(0)
	// u.AllowedUpdates = []string{"callback_query", "message"}
	u.Timeout = 60

	// bot.Debug = true

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.CallbackData() != "" {
			state, _ := UserState[update.CallbackQuery.From.ID]

			switch state {
			case "confirm":
				err = CallbackConfirm(bot, update)
			case "choose":
				err = CallbackChoose(bot, update)
			}

			if err != nil {
				log.Print(err)
			}

			continue
		}

		state, _ := UserState[update.Message.From.ID]
		switch state {
		case "input":
			if err := Input(bot, update); err != nil {
				log.Fatal(err)
			}
			continue
		}
		if !update.Message.IsCommand() && state == "idle" {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "sry, i work only with commands"))
		}

		switch update.Message.Command() {
		case "start":
			Start(bot, update)
		case "choose":
			Choose(bot, update)
		}

		fmt.Println()

	}
}
