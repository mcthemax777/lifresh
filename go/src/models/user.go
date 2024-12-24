package models

import "lifresh/custom_time"

type User struct {
	Id              int                    `gorm:"primary_key" json:"id"`
	AccountId       int                    `json:"account_id"`
	Nickname        string                 `json:"nickname"`
	ProfileImageUrl string                 `json:"profile_image_url"`
	UpdateDate      custom_time.CustomTime `json:"update_date"`
}

func (User) TableName() string {
	return "user"
}
