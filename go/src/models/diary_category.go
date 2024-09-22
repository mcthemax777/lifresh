package models

import "lifresh/custom_time"

type DiaryCategory struct {
	DiaryCategoryId int                    `gorm:"primary_key" json:"diary_category_id"`
	DiaryId         int                    `json:"diary_id"`
	Name            string                 `json:"name"`
	RGBColor        string                 `json:"rgb_color"`
	Priority        int                    `json:"priority"`
	OpenFlag        int                    `json:"open_flag"`
	UpdateDate      custom_time.CustomTime `json:"update_date"`
}

func (DiaryCategory) TableName() string {
	return "diary_category"
}
