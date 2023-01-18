package models

import "lifresh/custom_time"

type MoneyTask struct {
	MoneyTaskNo    int                    `gorm:"column:moneyTaskNo;primary_key;auto_increment;not_null" json:"moneyTaskNo"`
	CategoryType   int                    `gorm:"column:categoryType" json:"categoryType"`
	MainCategoryNo int                    `gorm:"column:mainCategoryNo" json:"mainCategoryNo"`
	SubCategoryNo  int                    `gorm:"column:subCategoryNo" json:"subCategoryNo"`
	StartTime      custom_time.CustomTime `gorm:"column:startTime" json:"startTime"`
	EndTime        custom_time.CustomTime `gorm:"column:endTime" json:"endTime"`
	Money          int                    `gorm:"column:money" json:"money"`
	Detail         string                 `gorm:"column:detail" json:"detail"`
	MoneyManagerNo int                    `gorm:"column:moneyManagerNo" json:"moneyManagerNo"`
	PlannerNo      int                    `gorm:"column:plannerNo" json:"plannerNo"`
	TodayNo        int                    `gorm:"column:todayNo" json:"todayNo"`
	Priority       int                    `gorm:"column:priority" json:"priority"`
	OverMoney      int                    `gorm:"column:overMoney" json:"overMoney"`
}

func (MoneyTask) TableName() string {
	return "MoneyTask"
}
