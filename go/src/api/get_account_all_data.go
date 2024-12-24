package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type GetAccountAllDataHandler struct {
	SessionApiHandler
}

func NewGetAccountAllDataHandler() GetAccountAllDataHandler {
	h := GetAccountAllDataHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h GetAccountAllDataHandler) process(reqBody []byte) ([]byte, error) {

	var req request.GetAccountAllDataReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	currentTime := CurrentTime()

	accountId, err := h.checkSession(req.Uid, req.Sid, currentTime)

	//세션 만료
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, err.Error())), err
	}

	user, err := db.DBHandlerSG.GetUserByAccountId(accountId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, err.Error())), err
	}

	planList, err := db.DBHandlerSG.GetPlanListByUserId(user.Id)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, err.Error())), err
	}
	planRecordList, err := db.DBHandlerSG.GetPlanRecordListByUserId(user.Id)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, err.Error())), err
	}
	planRecordOperatorList, err := db.DBHandlerSG.GetPlanRecordOperatorListByUserId(user.Id)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, err.Error())), err
	}
	planGoalList, err := db.DBHandlerSG.GetPlanGoalListByUserId(user.Id)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, err.Error())), err
	}
	
	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_ACCOUNT_ALL_DATA_RES)

	sendRes := res.(*response.GetAccountAllData)
	sendRes.User = user
	sendRes.PlanList = planList
	sendRes.PlanRecordList = planRecordList
	sendRes.PlanRecordOperatorList = planRecordOperatorList
	sendRes.PlanGoalList = planGoalList

	return ResponseToByteArray(sendRes), nil
}
