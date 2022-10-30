package models

import "lifresh/custom_time"

type ToDoTask struct {
	ToDoTaskNo    int                    `gorm:"column:toDoTaskNo;primary_key;auto_increment;not_null" json:"toDoTaskNo"`
	SubCategoryNo int                    `gorm:"column:subCategoryNo" json:"subCategoryNo"`
	StartTime     custom_time.CustomTime `gorm:"column:startTime" json:"startTime"`
	FinishDate    custom_time.CustomTime `gorm:"column:endTime" json:"endTime"`
	Score         int                    `gorm:"column:score" json:"score"`
	Detail        string                 `gorm:"column:detail" json:"detail"`
	PlannerNo     int                    `gorm:"column:plannerNo" json:"plannerNo"`
	TodayNo       int                    `gorm:"column:todayNo" json:"todayNo"`
}

func (ToDoTask) TableName() string {
	return "ToDoTask"
}
