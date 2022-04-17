package api

import (
	"db"
	"encoding/json"
	"models"
	"request"
	"response"
)

type removeTaskPlanHandler struct {
	SessionApiHandler
}

func NewRemoveTaskPlanHandler() removeTaskPlanHandler {
	h := removeTaskPlanHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h removeTaskPlanHandler) process(reqBody []byte) ([]byte, error) {

	var req request.RemoveTaskPlanReq
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


	today, err := db.DBHandlerSG.GetTodayByTodayNo(req.TodayNo, planner.PlannerNo)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "today not")), err
	}

	var taskPlan models.TaskPlan
	taskPlan.TodayNo = today.TodayNo
	taskPlan.TaskPlanNo = req.TaskPlanNo

	err = db.DBHandlerSG.RemoveTaskPlan(&taskPlan)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "RemoveTaskPlan")), err
	}

	
	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.REMOVE_TASK_PLAN_RES)

	sendRes := res. (*response.RemoveTaskPlanRes)
	
	return ResponseToByteArray(sendRes), nil
}
