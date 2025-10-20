package response

import (
	"lifresh/define"
	"lifresh/internal/domain"
)

const (
	FAIL_RES = iota
	LOGIN_RES
	REFRESH_RES
	PING_RES
	GET_USER_RES
	CREATE_USER_RES
	UPDATE_USER_RES
	GET_FOLDER_RES
	CREATE_FOLDER_RES
	UPDATE_FOLDER_RES
	DELETE_FOLDER_RES
	GET_PLAN_RES
	CREATE_PLAN_RES
	UPDATE_PLAN_RES
	DELETE_PLAN_RES

	SIGN_UP_RES
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

	case REFRESH_RES:
		var res RefreshRes
		res.init(successCode, successMsg)

		return &res

	case CREATE_USER_RES:
		var res CreateUserRes
		res.init(successCode, successMsg)

		return &res

	case UPDATE_USER_RES:
		var res UpdateUserRes
		res.init(successCode, successMsg)

		return &res

	case PING_RES:
		var res PingRes
		res.init(successCode, successMsg)

		return &res

	case GET_USER_RES:
		var res GetUserRes
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
	AccessToken  string          `json:"accessToken"`
	RefreshToken string          `json:"refreshToken"`
	IdToken      string          `json:"idToken"`
	Account      *domain.Account `json:"account"`
}

func (res *LoginRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type RefreshRes struct {
	BaseResponse
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (res *RefreshRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type PingRes struct {
	AccountID define.SnowflakeID `json:"accountId"`
	BaseResponse
}

func (res *PingRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type CreateUserRes struct {
	User *domain.User `json:"user"`
	BaseResponse
}

func (res *CreateUserRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type UpdateUserRes struct {
	User *domain.User `json:"user"`
	BaseResponse
}

func (res *UpdateUserRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}

type GetUserRes struct {
	User *domain.User `json:"user"`
	BaseResponse
}

func (res *GetUserRes) init(code int, msg string) {
	res.BaseResponse = BaseResponse{ResultCode: code, ResultMsg: msg}
}
