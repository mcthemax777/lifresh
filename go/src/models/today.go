package models

type Today struct {
	TodayNo    int    `gorm:"column:todayNo;primary_key;auto_increment;not_null"`
	PlannerNo  int    `gorm:"column:plannerNo"`
	TodayIndex int    `gorm:"column:todayIndex"`
	todayType  int    `gorm:"column:todayType"`
	Diary      string `gorm:"column:diary"`
}

func (Today) TableName() string {
	return "Today"
}
