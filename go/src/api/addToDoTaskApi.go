package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/models"
	"lifresh/request"
	"lifresh/response"
)

type AddToDoTaskListHandler struct {
	SessionApiHandler
}

func NewAddToDoTaskListHandler() AddToDoTaskListHandler {
	h := AddToDoTaskListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h AddToDoTaskListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddToDoTaskListReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	currentTime := CurrentTime()

	accountNo, err := h.checkSession(req.Uid, req.Sid, currentTime)

	//세션 만료
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "session invalid")), err
	}

	planner, err := db.DBHandlerSG.GetPlannerByAccountNo(accountNo)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "plannerNo not")), err
	}

	var insertList []models.ToDoTask
	var updateList []models.ToDoTask

	for _, toDoTask := range req.ToDoTaskList {

		toDoTask.PlannerNo = planner.PlannerNo

		//신규 등록
		if toDoTask.ToDoTaskNo == 0 {

			insertList = append(insertList, toDoTask)

		} else { //변경

			updateList = append(updateList, toDoTask)

		}
	}

	if len(insertList) > 0 {
		err = db.DBHandlerSG.InsertToDoTask(&insertList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertToDoTask")), err
		}
	}

	if len(updateList) > 0 {
		err = db.DBHandlerSG.UpdateToDoTask(&updateList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdateToDoTask")), err
		}
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_TO_DO_TASK_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
