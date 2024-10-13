package db

import (
	"github.com/somuthink/sirius_weather_bot/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() error {

	var err error
	DB, err = gorm.Open(sqlite.Open("data/gorm.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&models.Users{}, &models.TimeTable{})
	if err != nil {
		return err
	}

	return nil
}
