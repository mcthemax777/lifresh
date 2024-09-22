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
	SocialType  int    `json:"social_type"`
	SocialToken string `json:"social_token"`
}

type SignUpReq struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

type GetAccountAllDataReq struct {
	Uid string `json:"uid"`
	Sid string `json:"sid"`
}

type AddDiaryCategoryListReq struct {
	Uid               string                 `json:"uid"`
	Sid               string                 `json:"sid"`
	DiaryCategoryList []models.DiaryCategory `json:"diary_category_list"`
}
type AddDiaryHistoryListReq struct {
	Uid              string                `json:"uid"`
	Sid              string                `json:"sid"`
	DiaryHistoryList []models.DiaryHistory `json:"diary_history_list"`
}

type RemoveDiaryCategoryListReq struct {
	Uid                 string `json:"uid"`
	Sid                 string `json:"sid"`
	DiaryCategoryIdList []int  `json:"diary_category_id_list"`
}

type RemoveDiaryHistoryListReq struct {
	Uid                string `json:"uid"`
	Sid                string `json:"sid"`
	DiaryHistoryIdList []int  `json:"diary_history_id_list"`
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
