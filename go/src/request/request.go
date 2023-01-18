package request

import (
	"lifresh/models"
)

type Request interface {
}

type BaseRequest struct {
	Uid string `json:"uid"`
	Sid string `json:"sid"`
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
	Uid string `json:"uid"`
	Sid string `json:"sid"`
}

type GetMainCategoryReq struct {
	Uid string `json:"uid"`
	Sid string `json:"sid"`
}

type GetSubCategoryReq struct {
	Uid string `json:"uid"`
	Sid string `json:"sid"`
}

type GetScheduleTaskReq struct {
	Uid string `json:"uid"`
	Sid string `json:"sid"`
}

type GetToDoTaskReq struct {
	Uid string `json:"uid"`
	Sid string `json:"sid"`
}

type GetMoneyTaskReq struct {
	Uid string `json:"uid"`
	Sid string `json:"sid"`
}

type AddMainCategoryListReq struct {
	Uid              string                `json:"uid"`
	Sid              string                `json:"sid"`
	MainCategoryList []models.MainCategory `json:"mainCategoryList"`
}

type AddSubCategoryListReq struct {
	Uid             string               `json:"uid"`
	Sid             string               `json:"sid"`
	SubCategoryList []models.SubCategory `json:"subCategoryList"`
}

type AddMoneyManagerListReq struct {
	Uid              string                `json:"uid"`
	Sid              string                `json:"sid"`
	MoneyManagerList []models.MoneyManager `json:"moneyManagerList"`
}

type AddScheduleTaskListReq struct {
	Uid              string                `json:"uid"`
	Sid              string                `json:"sid"`
	ScheduleTaskList []models.ScheduleTask `json:"scheduleTaskList"`
}

type AddToDoTaskListReq struct {
	Uid          string            `json:"uid"`
	Sid          string            `json:"sid"`
	ToDoTaskList []models.ToDoTask `json:"todoTaskList"`
}

type AddMoneyTaskListReq struct {
	Uid           string             `json:"uid"`
	Sid           string             `json:"sid"`
	MoneyTaskList []models.MoneyTask `json:"moneyTaskList"`
}

type RemoveMainCategoryListReq struct {
	Uid                string `json:"uid"`
	Sid                string `json:"sid"`
	MainCategoryNoList []int  `json:"mainCategoryNoList"`
}

type RemoveSubCategoryListReq struct {
	Uid               string `json:"uid"`
	Sid               string `json:"sid"`
	SubCategoryNoList []int  `json:"subCategoryNoList"`
}

type RemoveScheduleTaskListReq struct {
	Uid                string `json:"uid"`
	Sid                string `json:"sid"`
	ScheduleTaskNoList []int  `json:"subScheduleTaskNoList"`
}

type RemoveToDoTaskListReq struct {
	Uid            string `json:"uid"`
	Sid            string `json:"sid"`
	ToDoTaskNoList []int  `json:"toDoTaskNoList"`
}

type RemoveMoneyTaskListReq struct {
	Uid             string `json:"uid"`
	Sid             string `json:"sid"`
	MoneyTaskNoList []int  `json:"moneyTaskNoList"`
}
