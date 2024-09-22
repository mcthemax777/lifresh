package models

import "lifresh/custom_time"

type Plan struct {
	PlanId                int                    `gorm:"primary_key" json:"plan_id"`
	PlannerId             int                    `json:"planner_id"`
	PlanCategoryId        int                    `json:"plan_category_id"`
	Name                  string                 `json:"name"`
	RGBColor              string                 `json:"rgb_color"`
	Priority              int                    `json:"priority"`
	RepeatType            int                    `json:"repeat_type"`
	RepeatValueList       string                 `json:"repeat_value_list"`
	RecordTypeList        string                 `json:"record_type_list"`
	RecordCombineTypeList string                 `json:"record_combine_type_list"`
	DisplayRecordType     int                    `json:"display_record_type"`
	StartDate             custom_time.CustomTime `json:"start_date"`
	FinishDate            custom_time.CustomTime `json:"finish_date"`
	OpenFlag              int                    `json:"open_flag"`
	UpdateDate            custom_time.CustomTime `json:"update_date"`
}

func (Plan) TableName() string {
	return "plan"
}
