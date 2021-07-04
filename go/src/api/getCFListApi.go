package api

import (
	"db"
	"encoding/json"
	"request"
	"response"
)

type getCFList struct {
	SessionApiHandler
}

func NewGetCFListHandler() getCFList {
	h := getCFList{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h getCFList) process(reqBody []byte) ([]byte, error) {

	var req request.GetCFListReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	currentTime := CurrentTime()

	accountNo, err := h.checkSession(req.Sid, currentTime)
	
	//세션 만료
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "session invalid")), err
	}


	planner, err := db.DBHandlerSG.GetPlannerByAccountNo(accountNo)
		
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "plannerNo not")), err
	}

	//cf들 가져오기
	cfList, err := db.DBHandlerSG.GetCFListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "planner db error")), err
	}
	mcfList, err := db.DBHandlerSG.GetMCFListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "planner db error")), err
	}
	dcfList, err := db.DBHandlerSG.GetDCFListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "planner db error")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_CF_LIST_RES)

	sendRes := res. (*response.GetCFListRes)
	sendRes.CFList = cfList
	sendRes.MCFList = mcfList
	sendRes.DCFList = dcfList
	
	
	return ResponseToByteArray(sendRes), nil
}
