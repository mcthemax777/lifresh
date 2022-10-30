package models

type MainCategory struct {
	MainCategoryNo int    `gorm:"column:mainCategoryNo;primary_key;auto_increment;not_null" json:"mainCategoryNo"`
	PlannerNo      int    `gorm:"column:plannerNo" json:"plannerNo"`
	CategoryType   int    `gorm:"column:categoryType" json:"categoryType"`
	Name           string `gorm:"column:name" json:"name"`
}

func (MainCategory) TableName() string {
	return "MainCategory"
}
