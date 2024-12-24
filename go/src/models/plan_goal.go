package models

import "lifresh/custom_time"

type PlanGoal struct {
	Id                int                    `gorm:"primary_key" json:"id"`
	UserId            int                    `json:"user_id"`
	PlanId            int                    `json:"plan_id"`
	Name              string                 `json:"name"`
	CycleType         int                    `json:"cycle_type"`
	CycleNum          int                    `json:"cycle_num"`
	CycleOperatorType int                    `json:"cycle_operator_type"`
	Value             int                    `json:"value"`
	Unit              string                 `json:"unit"`
	Sort              int                    `json:"sort"`
	UpdateDate        custom_time.CustomTime `json:"update_date"`
}

func (PlanGoal) TableName() string {
	return "plan_goal"
}
