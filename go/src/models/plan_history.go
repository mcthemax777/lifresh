package models

import "lifresh/custom_time"

type PlanHistory struct {
	PlanHistoryId int                    `gorm:"primary_key" json:"plan_history_id"`
	PlanId        int                    `json:"plan_id"`
	PlannerId     int                    `json:"planner_id"`
	Today         int                    `json:"today"`
	RecordList    string                 `json:"record_list"`
	UpdateDate    custom_time.CustomTime `json:"update_date"`
}

func (PlanHistory) TableName() string {
	return "plan_history"
}
