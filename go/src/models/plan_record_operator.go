package models

import "lifresh/custom_time"

type PlanRecordOperator struct {
	Id                int                    `gorm:"primary_key" json:"id"`
	UserId            int                    `json:"user_id"`
	PlanRecordId      int                    `json:"plan_record_id"`
	Type              int                    `json:"type"`
	ChildPlanRecordId int                    `json:"child_plan_record_id"`
	Sort              int                    `json:"sort"`
	UpdateDate        custom_time.CustomTime `json:"update_date"`
}

func (PlanRecordOperator) TableName() string {
	return "plan_record_operator"
}
