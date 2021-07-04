package models

type CF struct {
	CfNo      int    `gorm:"column:cfNo;primary_key;auto_increment;not_null" json:"cfNo"`
	PlannerNo int    `gorm:"column:plannerNo" json:"_"`
	Name      string `gorm:"column:name" json:"name"`
}

func (CF) TableName() string {
	return "CF"
}