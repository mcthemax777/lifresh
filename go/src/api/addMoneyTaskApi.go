package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/models"
	"lifresh/request"
	"lifresh/response"
)

type AddMoneyTaskListHandler struct {
	SessionApiHandler
}

func NewAddMoneyTaskListHandler() AddMoneyTaskListHandler {
	h := AddMoneyTaskListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h AddMoneyTaskListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddMoneyTaskListReq
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

	var insertList []models.MoneyTask
	var updateList []models.MoneyTask

	for _, moneyTask := range req.MoneyTaskList {

		moneyTask.PlannerNo = planner.PlannerNo

		//신규 등록
		if moneyTask.MoneyTaskNo == 0 {

			insertList = append(insertList, moneyTask)

		} else { //변경

			updateList = append(updateList, moneyTask)
		}
	}

	if len(insertList) > 0 {
		err = db.DBHandlerSG.InsertMoneyTask(&insertList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertMoneyTask")), err
		}
	}

	if len(updateList) > 0 {
		err = db.DBHandlerSG.UpdateMoneyTask(&updateList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdateMoneyTask")), err
		}
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_MONEY_TASK_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
