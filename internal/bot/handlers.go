package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/somuthink/sirius_weather_bot/internal/db"
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

	chatId := update.Message.Chat.ID

	msg := tgbotapi.NewMessage(chatId, "choose weather sending timetable")

	chooseKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄Every Minute", fmt.Sprintf("minute")),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌅Morning", fmt.Sprintf("morning")),
			tgbotapi.NewInlineKeyboardButtonData("🌇Afternoon", fmt.Sprintf("afternoon")),
			tgbotapi.NewInlineKeyboardButtonData("🌃Evening", fmt.Sprintf("evening")),
		),
	)

	msg.ReplyMarkup = chooseKeyboard

	_, err := bot.Send(msg)

	return err
}

func CallbackConfirm(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	callbackData := strings.Split(update.CallbackData(), ",")
	userId := update.CallbackQuery.From.ID
	confirm := callbackData[0]
	messageID, err := strconv.Atoi(callbackData[2])
	if err != nil {
		return err
	}
	var text string
	switch string(confirm) {
	case "y":
		if err := db.InsertUsers(userId, callbackData[1]); err != nil {
			return err
		}
		text = fmt.Sprintf("succefully set city to %s", callbackData[1])
		UserState[userId] = "choose"
	case "n":
		text = "Type city name again"
		UserState[userId] = "input"
	}
	msg := tgbotapi.NewEditMessageText(update.FromChat().ChatConfig().ChatID, messageID, text)
	_, err = bot.Send(msg)

	return err
}
