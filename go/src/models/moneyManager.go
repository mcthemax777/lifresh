package models

type MoneyManager struct {
	MoneyManagerNo       int    `gorm:"column:moneyManagerNo;primary_key;auto_increment;not_null" json:"moneyManagerNo"`
	MoneyManagerType     int    `gorm:"column:moneyManagerType" json:"moneyManagerType"`
	LinkedMoneyManagerNo int    `gorm:"column:linkedMoneyManagerNo" json:"linkedMoneyManagerNo"`
	Money                int    `gorm:"column:money" json:"money"`
	Name                 string `gorm:"column:name" json:"name"`
	Detail               string `gorm:"column:detail" json:"detail"`
	CalcDate             int    `gorm:"column:calcDate" json:"calcDate"`
	PayDate              int    `gorm:"column:payDate" json:"payDate"`
	PlannerNo            int    `gorm:"column:plannerNo" json:"plannerNo"`
}

func (MoneyManager) TableName() string {
	return "MoneyManager"
}
