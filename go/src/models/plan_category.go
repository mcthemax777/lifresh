package models

import "lifresh/custom_time"

type PlanCategory struct {
	PlanCategoryId       int                    `gorm:"primary_key" json:"plan_category_id"`
	PlannerId            int                    `json:"planner_id"`
	ParentPlanCategoryId int                    `json:"parent_plan_category_id"`
	Name                 string                 `json:"name"`
	RGBColor             string                 `json:"rgb_color"`
	Priority             int                    `json:"priority"`
	OpenFlag             int                    `json:"open_flag"`
	UpdateDate           custom_time.CustomTime `json:"update_date"`
}

func (PlanCategory) TableName() string {
	return "plan_category"
}
