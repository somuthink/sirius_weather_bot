package db

import (
	"fmt"

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
	err := DB.Model(models.TimeTable{}).Where("tg_id=?", chatId).Update(time, val).Error

	return err
}

func SelectAllTimeUsers(time string) ([]int64, error) {
	var res []int64
	return res, DB.Model(models.TimeTable{}).Where(fmt.Sprintf("%s = ?", time), true).Pluck("tg_id", &res).Error
}
