package models

import (
	"custom_time"
)


type TaskPlan struct {
	TaskPlanNo int       `gorm:"column:taskPlanNo;primary_key;auto_increment;not_null"`
	DcfNo      int       `gorm:"column:dcfNo"`
	StartTime  custom_time.CustomTime `gorm:"column:startTime"`
	FinishTime custom_time.CustomTime `gorm:"column:finishTime"`
	Priority   int    `gorm:"column:priority"`
	TodayNo    int    `gorm:"column:todayNo"`
}

func (TaskPlan) TableName() string {
	return "TaskPlan"
}