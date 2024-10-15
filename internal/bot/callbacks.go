package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/somuthink/sirius_weather_bot/internal/db"
	"gorm.io/gorm"
)

func CallbackConfirm(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	callbackData := strings.Split(update.CallbackData(), ",")
	chatId := update.FromChat().ChatConfig().ChatID
	confirm := callbackData[0]
	messageId, err := strconv.Atoi(callbackData[2])
	if err != nil {
		return err
	}
	var text string
	switch string(confirm) {
	case "y":
		if err := db.InsertUsers(chatId, callbackData[1]); err != nil {
			return err
		}
		text = fmt.Sprintf("succefully set city to %s", callbackData[1])
		defer Choose(bot, update)
	case "n":
		text = "Type city name again"
		UserState[chatId] = "input"
	}
	msg := tgbotapi.NewEditMessageText(update.FromChat().ChatConfig().ChatID, messageId, text)
	_, err = bot.Send(msg)

	return err
}

func getButtonLabel(activeLabel, inactiveLabel string, isActive bool) string {
	if isActive {
		return activeLabel
	} else {
		return inactiveLabel
	}
}

func getButtonQuery(timeName string, isActive bool, messageId int) string {
	if isActive {
		return fmt.Sprintf("a,%s,%d", timeName, messageId)
	} else {
		return fmt.Sprintf("i,%s,%d", timeName, messageId)
	}
}

func ChooseKeyboardBuilder(chatId int64, messageId int) (tgbotapi.InlineKeyboardMarkup, error) {
	var chooseKeyboard tgbotapi.InlineKeyboardMarkup
	timetable, err := db.SelectUserTimeTable(chatId)

	if err == gorm.ErrRecordNotFound {
		fmt.Println("record not found")
	}

	minuteButton := tgbotapi.NewInlineKeyboardButtonData(
		getButtonLabel("⏰minute", "minute", timetable.Minute),
		getButtonQuery("minute", timetable.Minute, messageId),
	)

	morningButton := tgbotapi.NewInlineKeyboardButtonData(
		getButtonLabel("🌅morning", "morning", timetable.Morning),
		getButtonQuery("morning", timetable.Morning, messageId),
	)
	afternoonButton := tgbotapi.NewInlineKeyboardButtonData(
		getButtonLabel("🌇afternoon", "afternoon", timetable.Afternoon),
		getButtonQuery("afternoon", timetable.Afternoon, messageId),
	)

	eveningButton := tgbotapi.NewInlineKeyboardButtonData(
		getButtonLabel("🌃evening", "evening", timetable.Evening),
		getButtonQuery("evening", timetable.Evening, messageId),
	)
	chooseKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(minuteButton),
		tgbotapi.NewInlineKeyboardRow(morningButton, afternoonButton, eveningButton),
	)

	return chooseKeyboard, nil
}

func CallbackChoose(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	callbackData := strings.Split(update.CallbackData(), ",")
	choose := callbackData[0] == "a"
	if err := db.InsertUserTimeTable(update.FromChat().ID, callbackData[1], !choose); err != nil {
		return err
	}
	messageId, err := strconv.Atoi(callbackData[2])
	if err != nil {
		return err
	}

	msg := tgbotapi.NewEditMessageText(update.FromChat().ChatConfig().ChatID, messageId, "updated, choose weather sending timetable ")
	chooseKeyboard, err := ChooseKeyboardBuilder(update.FromChat().ChatConfig().ChatID, messageId)
	if err != nil {
		return err
	}
	msg.ReplyMarkup = &chooseKeyboard
	_, err = bot.Send(msg)

	return nil
}
