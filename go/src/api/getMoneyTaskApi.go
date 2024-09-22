package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type GetMoneyTaskHandler struct {
	SessionApiHandler
}

func NewGetMoneyTaskHandler() GetMoneyTaskHandler {
	h := GetMoneyTaskHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h GetMoneyTaskHandler) process(reqBody []byte) ([]byte, error) {

	var req request.GetMoneyTaskReq
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

	todayList, err := db.DBHandlerSG.GetTodayListByPlannerNo(planner.PlannerId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "GetTodayListByPlannerNo")), err
	}

	moneyTaskList, err := db.DBHandlerSG.GetMoneyTaskListByPlannerNo(planner.PlannerId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "planner db error")), err
	}

	//category 가져오기
	mainCategoryList, err := db.DBHandlerSG.GetMainCategoryListByPlannerNoAndCategoryTypeList(planner.PlannerId, []int{0, 1})
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "GetMainCategoryListByPlannerNoAndCategoryTypeList")), err
	}

	var mainCategoryNoList []int

	for _, mainCategory := range mainCategoryList {

		mainCategoryNoList = append(mainCategoryNoList, mainCategory.MainCategoryNo)

	}

	subCategoryList, err := db.DBHandlerSG.GetSubCategoryListByPlannerNoAndMainCategoryNoList(planner.PlannerId, mainCategoryNoList)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "GetSubCategoryListByPlannerNoAndMainCategoryNoList")), err
	}

	//자산 가져오기
	moneyManagerList, err := db.DBHandlerSG.GetMoneyManagerListByPlannerNo(planner.PlannerId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "GetMoneyManagerListByPlannerNo")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_MONEY_TASK_RES)

	sendRes := res.(*response.GetMoneyTaskRes)
	sendRes.Planner = planner
	sendRes.TodayList = todayList
	sendRes.MainCategoryList = mainCategoryList
	sendRes.SubCategoryList = subCategoryList
	sendRes.MoneyTaskList = moneyTaskList
	sendRes.MoneyManagerList = moneyManagerList

	return ResponseToByteArray(sendRes), nil
}
