package request

import (
	"lifresh/models"
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

type GetUserAllDataReq struct {
	Sid string `json:"sid"`
}

type GetMainCategoryReq struct {
	Sid string `json:"sid"`
}

type GetSubCategoryReq struct {
	Sid string `json:"sid"`
}

type GetScheduleTaskReq struct {
	Sid string `json:"sid"`
}

type GetToDoTaskReq struct {
	Sid string `json:"sid"`
}

type GetMoneyTaskReq struct {
	Sid string `json:"sid"`
}

type AddMainCategoryListReq struct {
	Sid              string                `json:"sid"`
	MainCategoryList []models.MainCategory `json:"mainCategoryList"`
}

type AddSubCategoryListReq struct {
	Sid             string               `json:"sid"`
	SubCategoryList []models.SubCategory `json:"subCategoryList"`
}

type AddScheduleTaskListReq struct {
	Sid              string                `json:"sid"`
	ScheduleTaskList []models.ScheduleTask `json:"scheduleTaskList"`
}

type AddToDoTaskListReq struct {
	Sid          string            `json:"sid"`
	ToDoTaskList []models.ToDoTask `json:"todoTaskList"`
}

type AddMoneyTaskListReq struct {
	Sid           string             `json:"sid"`
	MoneyTaskList []models.MoneyTask `json:"moneyTaskList"`
}

type RemoveMainCategoryListReq struct {
	Sid                string `json:"sid"`
	MainCategoryNoList []int  `json:"mainCategoryNoList"`
}

type RemoveSubCategoryListReq struct {
	Sid               string `json:"sid"`
	SubCategoryNoList []int  `json:"subCategoryNoList"`
}

type RemoveScheduleTaskListReq struct {
	Sid                string `json:"sid"`
	ScheduleTaskNoList []int  `json:"subScheduleTaskNoList"`
}

type RemoveToDoTaskListReq struct {
	Sid            string `json:"sid"`
	ToDoTaskNoList []int  `json:"toDoTaskNoList"`
}

type RemoveMoneyTaskListReq struct {
	Sid             string `json:"sid"`
	MoneyTaskNoList []int  `json:"moneyTaskNoList"`
}
