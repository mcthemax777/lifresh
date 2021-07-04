package models

type DCF struct {
	DcfNo     int    `gorm:"column:dcfNo;primary_key;auto_increment;not_null" json:"dcfNo"`
	DcfType   int    `gorm:"column:dcfType" json:"dcfType"`
	McfNo     int    `gorm:"column:mcfNo" json:"mcfNo"`
	Name      string `gorm:"column:name" json:"name"`
	Priority  int    `gorm:"column:priority" json:"priority"`
	PlannerNo int    `gorm:"column:plannerNo" json:"_"`
}

func (DCF) TableName() string {
	return "DCF"
}