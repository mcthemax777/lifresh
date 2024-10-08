package models

import "lifresh/custom_time"

type Planner struct {
	PlannerId  int                    `gorm:"primary_key" json:"planner_id"`
	AccountId  int                    `json:"account_id"`
	UpdateDate custom_time.CustomTime `json:"update_date"`
}

func (Planner) TableName() string {
	return "planner"
}
