package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"lifresh/db"
	"lifresh/redis"
	"lifresh/request"
	"lifresh/response"
	"strings"
)

type LoginHandler struct {
}

func (h LoginHandler) process(reqBody []byte) ([]byte, error) {

	var req request.LoginReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	account, err := db.DBHandlerSG.GetAccountByUserId(req.UserId, req.Password)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(202, "invalid_parameter")), err
	}

	//데이터 가져오기
	planner, err := db.DBHandlerSG.GetPlannerByAccountNo(account.AccountNo)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "planner db error")), err
	}

	//session id 생성
	sid := uuid.New().String()
	sid = strings.Replace(sid, "-", "", -1)

	err = redis.RedisHandlerSG.SetSession(sid, account.AccountNo)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "redis error")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.LOGIN_RES)

	loginRes := res.(*response.LoginRes)
	loginRes.Sid = sid
	loginRes.Account = account
	loginRes.Planner = planner

	return ResponseToByteArray(loginRes), nil
}
