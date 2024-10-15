package db

import (
	"github.com/somuthink/sirius_weather_bot/internal/models"
	// "gorm.io/gorm"
	// "gorm.io/gorm/clause"
)

func SelectUserTimeTable(chatId int64) (models.TimeTable, error) {
	timeTable := models.TimeTable{Tg_id: chatId}
	err := DB.FirstOrCreate(&timeTable).Error
	return timeTable, err
}

func InsertUserTimeTable(chatId int64, time string, val bool) error {
	var timeTable models.TimeTable
	err := DB.Model(&timeTable).Where("tg_id=?", chatId).Update(time, val).Error

	return err
}
