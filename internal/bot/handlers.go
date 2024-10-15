package bot

import (
	"fmt"

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
	chatId := update.Message.Chat.ID

	msg := tgbotapi.NewMessage(chatId, "")
	city, err := weather.CheckCityExists(update.Message.Text)

	if err == nil {
		text := fmt.Sprintf("you wanna choose %s, right?", city)
		msg.Text = text
		sent_msg, err := bot.Send(msg)
		if err != nil {
			return err
		}
		msg_edited := tgbotapi.NewEditMessageText(chatId, sent_msg.MessageID, text)

		confirmKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅", fmt.Sprintf("y,%s,%d", city, sent_msg.MessageID)),
				tgbotapi.NewInlineKeyboardButtonData("❌", fmt.Sprintf("n,%s,%d", city, sent_msg.MessageID)),
			))
		msg_edited.ReplyMarkup = &confirmKeyboard
		UserState[chatId] = "confirm"
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

func Choose(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	chatId := update.FromChat().ChatConfig().ChatID

	msg := tgbotapi.NewMessage(chatId, "choose weather sending timetable")

	text := "choose weather sending timetable"
	msg.Text = text
	sent_msg, err := bot.Send(msg)
	if err != nil {
		return err
	}
	chooseKeyboard, err := ChooseKeyboardBuilder(chatId, sent_msg.MessageID)
	if err != nil {
		return err
	}
	msg_edited := tgbotapi.NewEditMessageText(chatId, sent_msg.MessageID, text)

	msg_edited.ReplyMarkup = &chooseKeyboard
	UserState[chatId] = "choose"
	_, err = bot.Send(msg_edited)
	return err
}
