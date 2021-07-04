package request

import (
	"custom_time"
	"models"
)


type Request interface {

}

type LoginReq struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

type SignUpReq struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

type GetTodayInfoReq struct {
	Sid string `json:"sid"` 
	Time custom_time.CustomTime `json:"time"`
}

type GetMonthInfoReq struct {
	Sid string `json:"sid"` 
	Time custom_time.CustomTime `json:"time"`
}

type GetCFListReq struct {
	Sid string `json:"sid"` 
}

type AddCFListReq struct {
	Sid string `json:"sid"` 
	CF  models.CF  `json:"cf"`
	MCF models.MCF `json:"mcf"`
	DCF models.DCF `json:"dcf"`
}

type AddTaskPlanReq struct {
	Sid string `json:"sid"` 
	DCFNo     int `json:"dcfNo"`
	StartTime custom_time.CustomTime `json:"startTime"`
	FinishTime custom_time.CustomTime `json:"finishTime"`
	Priority int    `json:"priority"`
	TodayNo int    `json:"todayNo"`
}

type AddTaskReq struct {
	Sid string `json:"sid"` 
	DCFNo     int `json:"dcfNo"`
	StartTime custom_time.CustomTime `json:"startTime"`
	FinishTime custom_time.CustomTime `json:"finishTime"`
	Score int    `json:"score"`
	Memo string    `json:"memo"`
	TodayNo int    `json:"todayNo"`
}


