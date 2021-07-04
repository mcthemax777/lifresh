package api

import (
	"db"
	"encoding/json"
	"models"
	"request"
	"response"
)

type addTaskPlanHandler struct {
	SessionApiHandler
}

func NewAddTaskPlanHandler() addTaskPlanHandler {
	h := addTaskPlanHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h addTaskPlanHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddTaskPlanReq
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
	taskPlan.DcfNo = req.DCFNo
	taskPlan.StartTime = req.StartTime
	taskPlan.FinishTime = req.FinishTime
	taskPlan.Priority = req.Priority

	err = db.DBHandlerSG.InsertTaskPlan(&taskPlan)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "InsertTaskPlan")), err
	}

	
	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_TASK_PLAN_RES)

	sendRes := res. (*response.AddTaskPlanRes)
	sendRes.TaskPlan = taskPlan
	
	return ResponseToByteArray(sendRes), nil
}
