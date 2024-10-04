package models

import "lifresh/custom_time"

type Profile struct {
	ProfileId       int                    `gorm:"primary_key" json:"profile_id"`
	AccountId       int                    `json:"account_id"`
	Nickname        string                 `json:"nickname"`
	ProfileImageUrl string                 `json:"profile_image_url"`
	UpdateDate      custom_time.CustomTime `json:"update_date"`
}

func (Profile) TableName() string {
	return "profile"
}
