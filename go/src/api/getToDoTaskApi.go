package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type GetToDoTaskHandler struct {
	SessionApiHandler
}

func NewGetToDoTaskHandler() GetToDoTaskHandler {
	h := GetToDoTaskHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h GetToDoTaskHandler) process(reqBody []byte) ([]byte, error) {

	var req request.GetToDoTaskReq
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

	todayList, err := db.DBHandlerSG.GetTodayListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "GetTodayListByPlannerNo")), err
	}

	toDoTaskList, err := db.DBHandlerSG.GetToDoTaskListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "planner db error")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_TO_DO_TASK_RES)

	sendRes := res.(*response.GetToDoTaskRes)
	sendRes.Planner = planner
	sendRes.TodayList = todayList
	sendRes.ToDoTaskList = toDoTaskList

	return ResponseToByteArray(sendRes), nil
}
