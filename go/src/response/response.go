package response

import "lifresh/models"

const (
	FAIL_RES = iota
	LOGIN_RES
	SIGN_UP_RES
	GET_ACCOUNT_ALL_DATA_RES
	ADD_DIARY_CATEGORY_RES
	ADD_DIARY_HISTORY_RES
	REMOVE_DIARY_CATEGORY_RES
	REMOVE_DIARY_HISTORY_RES
	GET_USER_ALL_DATA_RES
	GET_TODAY_RES
	GET_MAIN_CATEGORY_RES
	GET_SUB_CATEGORY_RES
	GET_SCHEDULE_TASK_RES
	GET_TO_DO_TASK_RES
	GET_MONEY_TASK_RES
	ADD_MAIN_CATEGORY_RES
	ADD_SUB_CATEGORY_RES
	ADD_MONEY_MANAGER_RES
	ADD_SCHEDULE_TASK_RES
	ADD_TO_DO_TASK_RES
	ADD_MONEY_TASK_RES
	REMOVE_MAIN_CATEGORY_RES
	REMOVE_SUB_CATEGORY_RES
	REMOVE_SCHEDULE_TASK_RES
	REMOVE_TO_DO_TASK_RES
	REMOVE_MONEY_TASK_RES
)

type Response interface {
	init(code int, msg string)
}

// func init() {

// }

type BaseResponse struct {
	ResultCode int    `json:"result_code"`
	ResultMsg  string `json:"result_msg"`
}

func CreateFailResponse(code int, msg string) Response {
	var res FailRes
	res.init(code, msg)

	return &res
}

func CreateSuccessResponse(resType int) Response {

	successCode := 100
	successMsg := "success"
	switch resType {
	case LOGIN_RES:
		var res LoginRes
		res.init(successCode, successMsg)

		return &res

	case SIGN_UP_RES:
		var res SignUpRes
		res.init(successCode, successMsg)

		return &res
	case GET_ACCOUNT_ALL_DATA_RES:
		var res GetAccountAllData
		res.init(successCode, successMsg)

		return &res

	case GET_USER_ALL_DATA_RES:
		var res GetUserAllData
		res.init(successCode, successMsg)

		return &res

	case GET_MAIN_CATEGORY_RES:
		var res GetMainCategoryRes
		res.init(successCode, successMsg)

		return &res

	case GET_SUB_CATEGORY_RES:
		var res GetSubCategoryRes
		res.init(successCode, successMsg)

		return &res

	case GET_SCHEDULE_TASK_RES:
		var res GetScheduleTaskRes
		res.init(successCode, successMsg)

		return &res

	case GET_TO_DO_TASK_RES:
		var res GetToDoTaskRes
		res.init(successCode, successMsg)

		return &res

	case GET_MONEY_TASK_RES:
		var res GetMoneyTaskRes
		res.init(successCode, successMsg)

		return &res

	default:
		var res BasicRes
		res.init(successCode, successMsg)

		return &res
	}
}

type FailRes struct {
	BaseResponse
}

func (res *FailRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type BasicRes struct {
	BaseResponse
}

func (res *BasicRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type LoginRes struct {
	BaseResponse
	Uid     string         `json:"uid"`
	Sid     string         `json:"sid"`
	Account models.Account `json:"account"`
}

func (res *LoginRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type SignUpRes struct {
	BaseResponse
}

func (res *SignUpRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type GetAccountAllData struct {
	BaseResponse
	Profile          models.Profile        `json:"profile"`
	Planner          models.Planner        `json:"planner"`
	PlanCategoryList []models.PlanCategory `json:"plan_category_list"`
	PlanList         []models.Plan         `json:"plan_list"`
	PlanHistoryList  []models.PlanHistory  `json:"plan_history_list"`

	Money             models.Money           `json:"money"`
	MoneyCategoryList []models.MoneyCategory `json:"money_category_list"`
	MoneyHistoryList  []models.MoneyHistory  `json:"money_history_list"`

	Diary             models.Diary           `json:"diary"`
	DiaryCategoryList []models.DiaryCategory `json:"diary_category_list"`
	DiaryHistoryList  []models.DiaryHistory  `json:"diary_history_list"`
}

func (res *GetAccountAllData) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type GetUserAllData struct {
	BaseResponse
	Planner          models.Planner        `json:"planner"`
	TodayList        []models.Today        `json:"todayList"`
	MainCategoryList []models.MainCategory `json:"mainCategoryList"`
	SubCategoryList  []models.SubCategory  `json:"subCategoryList"`
	ScheduleTaskList []models.ScheduleTask `json:"scheduleTaskList"`
	ToDoTaskList     []models.ToDoTask     `json:"toDoTaskList"`
	MoneyTaskList    []models.MoneyTask    `json:"moneyTaskList"`
}

func (res *GetUserAllData) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type GetMainCategoryRes struct {
	BaseResponse
	Planner          models.Planner        `json:"planner"`
	TodayList        []models.Today        `json:"todayList"`
	MainCategoryList []models.MainCategory `json:"mainCategoryList"`
}

func (res *GetMainCategoryRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type GetSubCategoryRes struct {
	BaseResponse
	Planner         models.Planner       `json:"planner"`
	TodayList       []models.Today       `json:"todayList"`
	SubCategoryList []models.SubCategory `json:"subCategoryList"`
}

func (res *GetSubCategoryRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type GetScheduleTaskRes struct {
	BaseResponse
	Planner          models.Planner        `json:"planner"`
	TodayList        []models.Today        `json:"todayList"`
	ScheduleTaskList []models.ScheduleTask `json:"scheduleTaskList"`
}

func (res *GetScheduleTaskRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type GetToDoTaskRes struct {
	BaseResponse
	Planner      models.Planner    `json:"planner"`
	TodayList    []models.Today    `json:"todayList"`
	ToDoTaskList []models.ToDoTask `json:"toDoTaskList"`
}

func (res *GetToDoTaskRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type GetMoneyTaskRes struct {
	BaseResponse
	Planner          models.Planner        `json:"planner"`
	TodayList        []models.Today        `json:"todayList"`
	MainCategoryList []models.MainCategory `json:"mainCategoryList"`
	SubCategoryList  []models.SubCategory  `json:"subCategoryList"`
	MoneyTaskList    []models.MoneyTask    `json:"moneyTaskList"`
	MoneyManagerList []models.MoneyManager `json:"moneyManagerList"`
}

func (res *GetMoneyTaskRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}
