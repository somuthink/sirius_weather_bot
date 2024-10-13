package models

import "gorm.io/gorm"

type TimeTable struct {
	gorm.Model
	ID                                        uint `gorm:"primaryKey;not null"`
	Minute, Hour, Morning, Afternoon, Evening int64
}

type Users struct {
	ID    uint  `gorm:"primaryKey;not null"`
	Tg_id int64 `gorm:"UNIQUE"`
	City  string
}
