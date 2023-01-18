package models

type SubCategory struct {
	SubCategoryNo  int    `gorm:"column:subCategoryNo;primary_key;auto_increment;not_null" json:"subCategoryNo"`
	CategoryType   int    `gorm:"column:categoryType" json:"categoryType"`
	MainCategoryNo int    `gorm:"column:mainCategoryNo" json:"mainCategoryNo"`
	Name           string `gorm:"column:name" json:"name"`
	PlannerNo      int    `gorm:"column:plannerNo" json:"_"`
	Priority       int    `gorm:"column:priority" json:"priority"`
}

func (SubCategory) TableName() string {
	return "SubCategory"
}
