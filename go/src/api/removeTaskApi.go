package api

import (
	"db"
	"encoding/json"
	"models"
	"request"
	"response"
)

type removeTaskHandler struct {
	SessionApiHandler
}

func NewRemoveTaskHandler() removeTaskHandler {
	h := removeTaskHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h removeTaskHandler) process(reqBody []byte) ([]byte, error) {

	var req request.RemoveTaskReq
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

	var task models.Task
	task.TodayNo = today.TodayNo
	task.TaskNo = req.TaskNo

	err = db.DBHandlerSG.RemoveTask(&task)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "RemoveTask")), err
	}

	
	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.REMOVE_TASK_RES)

	sendRes := res. (*response.RemoveTaskRes)
	
	return ResponseToByteArray(sendRes), nil
}
