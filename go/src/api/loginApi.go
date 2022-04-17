package api

import (
	"db"
	"encoding/json"
	"redis"
	"request"
	"response"
	"strings"

	"github.com/google/uuid"
)

type loginHandler struct {
}

func (h loginHandler) process(reqBody []byte) ([]byte, error) {

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
	
	//cf들 가져오기
	cfList, err := db.DBHandlerSG.GetCFListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "GetCFListByPlannerNo")), err
	}
	mcfList, err := db.DBHandlerSG.GetMCFListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "GetMCFListByPlannerNo")), err
	}
	dcfList, err := db.DBHandlerSG.GetDCFListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "GetDCFListByPlannerNo")), err
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
	
	loginRes := res. (*response.LoginRes)
	loginRes.SessionId = sid
	loginRes.Account = account
	loginRes.Planner = planner
	loginRes.CFList = cfList
	loginRes.MCFList = mcfList
	loginRes.DCFList = dcfList
	

	return ResponseToByteArray(loginRes), nil
}
