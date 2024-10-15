package models

import "gorm.io/gorm"

type TimeTable struct {
	gorm.Model
	ID                                        uint  `gorm:"primaryKey;not null"`
	Tg_id                                     int64 `gorm:"unique"`
	Minute, Hour, Morning, Afternoon, Evening bool
}

type Users struct {
	ID    uint  `gorm:"primaryKey;not null"`
	Tg_id int64 `gorm:"unique"`
	City  string
}
