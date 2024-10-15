package db

import (
	"github.com/somuthink/sirius_weather_bot/internal/models"
	// "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InsertUsers(chatId int64, city string) error {
	err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tg_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"city"}),
	}).Create(&models.Users{Tg_id: chatId, City: city}).Error

	return err
}

func SelectUserCity(chatId int64) (string, error) {
	var user models.Users
	err := DB.First(&user, "tg_id=?", chatId).Error

	return user.City, err
}
