package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type GetUserAllDataHandler struct {
	SessionApiHandler
}

func NewGetUserAllDataHandler() GetUserAllDataHandler {
	h := GetUserAllDataHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h GetUserAllDataHandler) process(reqBody []byte) ([]byte, error) {

	var req request.GetUserAllDataReq
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

	//category 가져오기
	mainCategoryList, err := db.DBHandlerSG.GetMainCategoryListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "GetMainCategoryListByPlannerNo")), err
	}
	subCategoryList, err := db.DBHandlerSG.GetSubCategoryListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "GetSubCategoryListByPlannerNo")), err
	}

	//task 가져오기
	scheduleTaskList, err := db.DBHandlerSG.GetScheduleTaskListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "planner db error")), err
	}
	toDoTaskList, err := db.DBHandlerSG.GetToDoTaskListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "planner db error")), err
	}
	moneyTaskList, err := db.DBHandlerSG.GetMoneyTaskListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "planner db error")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_USER_ALL_DATA_RES)

	sendRes := res.(*response.GetUserAllData)
	sendRes.Planner = planner
	sendRes.TodayList = todayList
	sendRes.MainCategoryList = mainCategoryList
	sendRes.SubCategoryList = subCategoryList
	sendRes.ScheduleTaskList = scheduleTaskList
	sendRes.ToDoTaskList = toDoTaskList
	sendRes.MoneyTaskList = moneyTaskList

	return ResponseToByteArray(sendRes), nil
}
