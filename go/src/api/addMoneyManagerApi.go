package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/models"
	"lifresh/request"
	"lifresh/response"
)

type AddMoneyManagerListHandler struct {
	SessionApiHandler
}

func NewAddMoneyManagerListHandler() AddMoneyManagerListHandler {
	h := AddMoneyManagerListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h AddMoneyManagerListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddMoneyManagerListReq
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

	var insertList []models.MoneyManager
	var updateList []models.MoneyManager

	for _, moneyManager := range req.MoneyManagerList {

		moneyManager.PlannerNo = planner.PlannerNo

		//신규 등록
		if moneyManager.MoneyManagerNo == 0 {
			insertList = append(insertList, moneyManager)
		} else { //변경
			updateList = append(updateList, moneyManager)
		}
	}

	if len(insertList) > 0 {
		err = db.DBHandlerSG.InsertMoneyManager(&insertList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertMoneyManager")), err
		}
	}

	if len(updateList) > 0 {
		err = db.DBHandlerSG.UpdateMoneyManager(&updateList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdateMoneyManager")), err
		}
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_MONEY_MANAGER_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
