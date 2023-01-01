package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type RemoveToDoTaskListHandler struct {
	SessionApiHandler
}

func NewRemoveToDoTaskListHandler() RemoveToDoTaskListHandler {
	h := RemoveToDoTaskListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h RemoveToDoTaskListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.RemoveToDoTaskListReq
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

	_, err = db.DBHandlerSG.DeleteToDoTaskList(planner.PlannerNo, req.ToDoTaskNoList)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "InsertCF")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.REMOVE_TO_DO_TASK_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
