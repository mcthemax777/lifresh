package models

import "custom_time"

type Today struct {
	TodayNo int       `gorm:"column:todayNo;primary_key;auto_increment;not_null"`
	PlannerNo int       `gorm:"column:plannerNo"`
	TodayDate custom_time.CustomTime `gorm:"column:todayDate"`
	Diary   string    `gorm:"column:diary"`
}

func (Today) TableName() string {
	return "Today"
}