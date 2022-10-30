package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type RemoveMoneyTaskListHandler struct {
	SessionApiHandler
}

func NewRemoveMoneyTaskListHandler() RemoveMoneyTaskListHandler {
	h := RemoveMoneyTaskListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h RemoveMoneyTaskListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.RemoveMoneyTaskListReq
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

	_, err = db.DBHandlerSG.DeleteMoneyTaskList(planner.PlannerNo, req.MoneyTaskNoList)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "InsertCF")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.REMOVE_MONEY_TASK_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
