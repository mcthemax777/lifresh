package models

import "lifresh/custom_time"

type ScheduleTask struct {
	ScheduleTaskNo int                    `gorm:"column:scheduleTaskNo;primary_key;auto_increment;not_null" json:"scheduleTaskNo"`
	SubCategoryNo  int                    `gorm:"column:subCategoryNo" json:"subCategoryNo"`
	StartTime      custom_time.CustomTime `gorm:"column:startTime" json:"startTime"`
	EndTime        custom_time.CustomTime `gorm:"column:endTime" json:"endTime"`
	Detail         string                 `gorm:"column:detail" json:"detail"`
	PlannerNo      int                    `gorm:"column:plannerNo" json:"plannerNo"`
	TodayNo        int                    `gorm:"column:todayNo" json:"todayNo"`
}

func (ScheduleTask) TableName() string {
	return "ScheduleTask"
}
