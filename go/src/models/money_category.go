package models

import "lifresh/custom_time"

type MoneyCategory struct {
	MoneyCategoryId       int                    `gorm:"primary_key" json:"money_category_id"`
	MoneyId               int                    `json:"money_id"`
	ParentMoneyCategoryId int                    `json:"parent_money_category_id"`
	Name                  string                 `json:"name"`
	RGBColor              string                 `json:"rgb_color"`
	Priority              int                    `json:"priority"`
	OpenFlag              int                    `json:"open_flag"`
	UpdateDate            custom_time.CustomTime `json:"update_date"`
}

func (MoneyCategory) TableName() string {
	return "money_category"
}
