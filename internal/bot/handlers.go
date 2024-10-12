package bot

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/somuthink/sirius_weather_bot/internal/weather"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "write the name of the city using english letter in a new message")
	msg.ReplyToMessageID = update.Message.MessageID
	UserState[update.Message.From.ID] = "input"
	bot.Send(msg)
}

func Input(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	city, err := weather.CheckCityExists(update.Message.Text)

	if err == nil {
		text := fmt.Sprintf("you wanna choose %s, right?", city)
		msg.Text = text
		sent_msg, err := bot.Send(msg)
		if err != nil {
			return err
		}
		msg_edited := tgbotapi.NewEditMessageText(update.Message.Chat.ID, sent_msg.MessageID, text)

		confirmKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅", fmt.Sprintf("%s%d", "y", sent_msg.MessageID)),
				tgbotapi.NewInlineKeyboardButtonData("❌", fmt.Sprintf("%s%d", "n", sent_msg.MessageID)),
			))
		msg_edited.ReplyMarkup = &confirmKeyboard
		UserState[update.Message.From.ID] = "confirm"
		_, err = bot.Send(msg_edited)
		return err
	} else if err == weather.ErrNotExistingCity {
		msg.Text = "it seems like there isn`t any city with name like that. lets try again"
	} else {
		msg.Text = "there was an error in bot`s work"
		return err
	}

	_, err = bot.Send(msg)

	return err

}

func CallbackConfirm(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	callbackData := update.CallbackData()
	confirm := callbackData[0]
	messageID, err := strconv.Atoi(callbackData[1:])
	if err != nil {
		return err
	}
	var text string
	switch string(confirm) {
	case "y":
		UserState[update.CallbackQuery.From.ID] = "idle"
		text = "Succesufly set"
	case "n":
		UserState[update.CallbackQuery.From.ID] = "input"
		text = "Type city name again"
	}
	msg := tgbotapi.NewEditMessageText(update.FromChat().ChatConfig().ChatID, messageID, text)
	_, err = bot.Send(msg)

	return err
}
