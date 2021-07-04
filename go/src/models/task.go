package models

import "custom_time"

type Task struct {
	TaskNo int       `gorm:"column:taskNo;primary_key;auto_increment;not_null"`
	DcfNo      int       `gorm:"column:dcfNo"`
	StartTime  custom_time.CustomTime `gorm:"column:startTime"`
	FinishTime custom_time.CustomTime `gorm:"column:finishTime"`
	Score   int    `gorm:"column:score"`
	Memo   string    `gorm:"column:memo"`
	TodayNo    int    `gorm:"column:todayNo"`
}

func (Task) TableName() string {
	return "Task"
}