package models

type SubCategory struct {
	SubCategoryNo  int    `gorm:"column:subCategoryNo;primary_key;auto_increment;not_null" json:"subCategoryNo"`
	MainCategoryNo int    `gorm:"column:mainCategoryNo" json:"mainCategoryNo"`
	Name           string `gorm:"column:name" json:"name"`
	PlannerNo      int    `gorm:"column:plannerNo" json:"_"`
}

func (SubCategory) TableName() string {
	return "SubCategory"
}