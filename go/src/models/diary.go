package models

import "lifresh/custom_time"

type Diary struct {
	DiaryId    int                    `gorm:"primary_key" json:"diary_id"`
	AccountId  int                    `json:"account_id"`
	UpdateDate custom_time.CustomTime `json:"update_date"`
}

func (Diary) TableName() string {
	return "diary"
}
