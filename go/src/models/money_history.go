package models

import "lifresh/custom_time"

type MoneyHistory struct {
	MoneyHistoryId  int                    `gorm:"primary_key" json:"money_history_id"`
	MoneyId         int                    `json:"money_id"`
	MoneyCategoryId int                    `json:"money_category_id"`
	Cost            int                    `json:"cost"`
	WasteCost       int                    `json:"waste_cost"`
	Detail          string                 `json:"detail"`
	UseDate         custom_time.CustomTime `json:"use_date"`
	OpenFlag        int                    `json:"open_flag"`
	UpdateDate      custom_time.CustomTime `json:"update_date"`
}

func (MoneyHistory) TableName() string {
	return "money_history"
}
