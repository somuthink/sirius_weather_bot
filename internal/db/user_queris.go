package db

import (
	"github.com/somuthink/sirius_weather_bot/internal/models"
	// "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InserUser(userId int64, city string) error {

	err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tg_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"city"}),
	}).Create(&models.Users{Tg_id: userId, City: city}).Error

	return err

}
