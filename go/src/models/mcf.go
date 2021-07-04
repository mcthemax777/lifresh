package models

type MCF struct {
	McfNo     int    `gorm:"column:mcfNo;primary_key;auto_increment;not_null" json:"mcfNo"`
	CfNo      int    `gorm:"column:cfNo" json:"cfNo"`
	Name      string `gorm:"column:name" json:"name"`
	PlannerNo int    `gorm:"column:plannerNo" json:"_"`
}

func (MCF) TableName() string {
	return "MCF"
}