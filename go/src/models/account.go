package models

import (
	"lifresh/custom_time"
)

type Account struct {
	AccountNo  int                    `gorm:"column:accountNo;primary_key;auto_increment;not_null" json:"accountNo"`
	UserId     string                 `gorm:"column:userId" json:"userId"`
	Password   string                 `gorm:"column:password" json:"password"`
	CreateTime custom_time.CustomTime `gorm:"column:createTime" json:"createTime"`
}

func (Account) TableName() string {
	return "Account"
}
