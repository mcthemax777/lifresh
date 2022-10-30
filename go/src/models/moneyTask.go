package models

import "lifresh/custom_time"

type MoneyTask struct {
	MoneyTaskNo   int                    `gorm:"column:moneyTaskNo;primary_key;auto_increment;not_null" json:"moneyTaskNo"`
	SubCategoryNo int                    `gorm:"column:subCategoryNo" json:"subCategoryNo"`
	StartTime     custom_time.CustomTime `gorm:"column:startTime" json:"startTime"`
	EndTime       custom_time.CustomTime `gorm:"column:endTime" json:"endTime"`
	Money         int                    `gorm:"column:money" json:"money"`
	Detail        string                 `gorm:"column:detail" json:"detail"`
	PlannerNo     int                    `gorm:"column:plannerNo" json:"plannerNo"`
	TodayNo       int                    `gorm:"column:todayNo" json:"todayNo"`
}

func (MoneyTask) TableName() string {
	return "MoneyTask"
}
