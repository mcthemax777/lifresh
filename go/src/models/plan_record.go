package models

import "lifresh/custom_time"

type PlanRecord struct {
	Id           int                    `gorm:"primary_key" json:"id"`
	UserId       int                    `json:"user_id"`
	PlanId       int                    `json:"plan_id"`
	Name         string                 `json:"name"`
	DefaultValue int                    `json:"default_value"`
	Unit         string                 `json:"unit"`
	Sort         int                    `json:"sort"`
	UpdateDate   custom_time.CustomTime `json:"update_date"`
}

func (PlanRecord) TableName() string {
	return "plan_record"
}
