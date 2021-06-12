package models

import (
	"time"
)

type Account struct {
	AccountNo  int
	UserId     string
	Password   string
	CreateTime time.Time
}
