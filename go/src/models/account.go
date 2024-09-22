package models

import "lifresh/custom_time"

type Account struct {
	AccountId   int                    `gorm:"primary_key" json:"account_id"`
	SocialType  int                    `json:"social_type"`
	SocialToken string                 `json:"social_token"`
	CreateDate  custom_time.CustomTime `json:"create_date"`
	UpdateDate  custom_time.CustomTime `json:"update_date"`
}

func (Account) TableName() string {
	return "account"
}
