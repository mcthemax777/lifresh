package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/define"
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

	var updateMoneyTaskNoList []int
	updateMoneyManagerList := make(map[int]int)

	var updateMoneyManagerNoList []int

	for _, moneyTask := range req.MoneyTaskList {

		moneyTask.PlannerNo = planner.PlannerNo

		//신규 등록
		if moneyTask.MoneyTaskNo == 0 {

			insertList = append(insertList, moneyTask)

		} else { //변경

			updateList = append(updateList, moneyTask)
			updateMoneyTaskNoList = append(updateMoneyTaskNoList, moneyTask.MoneyTaskNo)
		}

		money := moneyTask.Money

		//지출일 경우 마이너스 추가
		if moneyTask.CategoryType == define.CategoryTypeMoneyMinus {
			money *= -1
		}

		//업데이트할 자산 저장
		_, exists := updateMoneyManagerList[moneyTask.MoneyManagerNo]
		if exists {
			updateMoneyManagerList[moneyTask.MoneyManagerNo] += money
		} else {
			updateMoneyManagerList[moneyTask.MoneyManagerNo] = money
		}

		updateMoneyManagerNoList = append(updateMoneyManagerNoList, moneyTask.MoneyManagerNo)
	}

	//업데이트일 경우 자산금액 계산 다시하기
	updateMoneyTaskList, err := db.DBHandlerSG.GetMoneyTaskListByPlannerNoAndMoneyTaskNoList(planner.PlannerNo, updateMoneyTaskNoList)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "moneyManager not")), err
	}

	for _, moneyTask := range updateMoneyTaskList {

		money := moneyTask.Money

		//수입일 경우 마이너스 추가
		if moneyTask.CategoryType == define.CategoryTypeMoneyPlus {
			money *= -1
		}

		_, exists := updateMoneyManagerList[moneyTask.MoneyManagerNo]
		if exists {
			updateMoneyManagerList[moneyTask.MoneyManagerNo] += money
		} else {
			updateMoneyManagerList[moneyTask.MoneyManagerNo] = money
		}

		updateMoneyManagerNoList = append(updateMoneyManagerNoList, moneyTask.MoneyManagerNo)
	}

	//사용한 자산 확인
	moneyManagerList, err := db.DBHandlerSG.GetMoneyManagerListByPlannerNoAndMoneyManagerNoList(planner.PlannerNo, updateMoneyManagerNoList)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "moneyManager not")), err
	}

	for index, moneyManager := range moneyManagerList {

		_, exists := updateMoneyManagerList[moneyManager.MoneyManagerNo]
		if exists {
			moneyManagerList[index].Money += updateMoneyManagerList[moneyManager.MoneyManagerNo]
		} else {
			return ResponseToByteArray(response.CreateFailResponse(201, "moneyManager not exist - "+string(rune(moneyManager.MoneyManagerNo)))), nil
		}
	}

	err = db.DBHandlerSG.UpdateMoneyManager(&moneyManagerList)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "UpdateMoneyManager")), err
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
