package api

import (
	"db"
	"encoding/json"
	"models"
	"request"
	"response"
)

type addTaskHandler struct {
	SessionApiHandler
}

func NewAddTaskHandler() addTaskHandler {
	h := addTaskHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h addTaskHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddTaskReq
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
	task.DcfNo = req.DCFNo
	task.StartTime = req.StartTime
	task.FinishTime = req.FinishTime
	task.Score = req.Score
	task.Memo = req.Memo
	

	err = db.DBHandlerSG.InsertTask(&task)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "InsertTask")), err
	}

	
	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_TASK_RES)

	sendRes := res. (*response.AddTaskRes)
	sendRes.Task = task
	
	return ResponseToByteArray(sendRes), nil
}
