package models

type Planner struct {
	PlannerNo   int    `gorm:"column:plannerNo;primary_key;auto_increment;not_null"`
	AccountNo   int    `gorm:"column:accountNo"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
}

func (Planner) TableName() string {
	return "Planner"
}