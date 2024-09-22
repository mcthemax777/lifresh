package models

import "lifresh/custom_time"

type DiaryHistory struct {
	DiaryHistoryId  int                    `gorm:"primary_key" json:"diary_history_id"`
	DiaryId         int                    `json:"diary_id"`
	DiaryCategoryId int                    `json:"diary_category_id"`
	Today           int                    `json:"today"`
	Content         string                 `json:"content"`
	ImageUrlList    string                 `json:"image_url_list"`
	Priority        int                    `json:"priority"`
	OpenFlag        int                    `json:"open_flag"`
	UpdateDate      custom_time.CustomTime `json:"update_date"`
}

func (DiaryHistory) TableName() string {
	return "diary_history"
}
