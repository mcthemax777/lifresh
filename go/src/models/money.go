package models

import "lifresh/custom_time"

type Money struct {
	MoneyId    int                    `gorm:"primary_key" json:"money_id"`
	AccountId  int                    `json:"account_id"`
	UpdateDate custom_time.CustomTime `json:"update_date"`
}

func (Money) TableName() string {
	return "money"
}
