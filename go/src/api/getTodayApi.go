package api

import (
	"db"
	"encoding/json"
	"request"
	"response"
)

type getTodayHandler struct {
	SessionApiHandler
}

func NewGetTodayHandler() getTodayHandler {
	h := getTodayHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h getTodayHandler) process(reqBody []byte) ([]byte, error) {

	var req request.GetTodayInfoReq
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

	//get today
	today, err := db.DBHandlerSG.GetTodayByPlannerNoAndTime(planner.PlannerNo, currentTime)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "today not")), err
	}

	//get taskplan list
	taskPlanList, err := db.DBHandlerSG.GetTaskPlanListByToadyNoAndPlannerNo(today.TodayNo, planner.PlannerNo)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "taskPlanList")), err
	}

	//get task list
	taskList, err := db.DBHandlerSG.GetTaskListByToadyNoAndPlannerNo(today.TodayNo, planner.PlannerNo)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "taskList")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_TODAY_RES)

	sendRes := res. (*response.GetTodayInfoRes)
	sendRes.Today = today
	sendRes.TaskPlanList = taskPlanList
	sendRes.TaskList = taskList
	
	
	return ResponseToByteArray(sendRes), nil
}
