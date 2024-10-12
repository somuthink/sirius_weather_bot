package models

import "gorm.io/gorm"

type TimeTable struct {
	gorm.Model
	ID                                        int
	Minute, Hour, Morning, Afternoon, Evening int64
}

type Users struct {
	Tg_id int64
	City  string
}
