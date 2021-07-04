package response

import "models"

const (
	FAIL_RES = iota
	LOGIN_RES
	SIGN_UP_RES
	GET_TODAY_RES
	GET_CF_LIST_RES
	ADD_CF_RES
	ADD_TASK_PLAN_RES
	ADD_TASK_RES
)

type Response interface {
	init(code int, msg string)
}

// func init() {

// }

type BaseResponse struct {
	ResultCode int    `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
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
		
	case GET_TODAY_RES:
		var res GetTodayInfoRes
		res.init(successCode, successMsg)

		return &res
		
	case GET_CF_LIST_RES:
		var res GetCFListRes
		res.init(successCode, successMsg)

		return &res

	case ADD_CF_RES:
		var res AddCFRes
		res.init(successCode, successMsg)

		return &res

	case ADD_TASK_PLAN_RES:
		var res AddTaskPlanRes
		res.init(successCode, successMsg)

		return &res

	case ADD_TASK_RES:
		var res AddTaskRes
		res.init(successCode, successMsg)

		return &res

	default:
		return nil
	}
}

type FailRes struct {
	BaseResponse
}

func (res *FailRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type LoginRes struct {
	BaseResponse
	SessionId string `json:"sessionId"`
	Account   models.Account `json:"account"`
	Planner   models.Planner `json:"planner"`
	CFList	  []models.CF `json:"cfList"`
	MCFList	  []models.MCF `json:"mcfList"`
	DCFList	  []models.DCF `json:"dcfList"`
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

type GetCFListRes struct {
	BaseResponse
	CFList  []models.CF  `json:"cfList"`
	MCFList []models.MCF `json:"mcfList"`
	DCFList []models.DCF `json:"dcfList"`
}

func (res *GetCFListRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type GetTodayInfoRes struct {
	BaseResponse
	Today 		models.Today `json:"today"`
	TaskPlanList 	[]models.TaskPlan `json:"taskPlanList"`
	TaskList	  	[]models.Task `json:"taskList"`
}

func (res *GetTodayInfoRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type AddCFRes struct {
	BaseResponse
	CF	  models.CF `json:"cf"`
	MCF	  models.MCF `json:"mcf"`
	DCF	  models.DCF `json:"dcf"`
}

func (res *AddCFRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type AddTaskPlanRes struct {
	BaseResponse
	TaskPlan	  models.TaskPlan `json:"taskPlan"`
}

func (res *AddTaskPlanRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type AddTaskRes struct {
	BaseResponse
	Task	  models.Task `json:"task"`
}

func (res *AddTaskRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}