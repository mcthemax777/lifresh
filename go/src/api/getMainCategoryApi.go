package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type GetMainCategoryHandler struct {
	SessionApiHandler
}

func NewGetMainCategoryHandler() GetMainCategoryHandler {
	h := GetMainCategoryHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h GetMainCategoryHandler) process(reqBody []byte) ([]byte, error) {

	var req request.GetMainCategoryReq
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

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_MAIN_CATEGORY_RES)

	sendRes := res.(*response.GetMainCategoryRes)
	sendRes.Planner = planner
	sendRes.TodayList = todayList
	sendRes.MainCategoryList = mainCategoryList

	return ResponseToByteArray(sendRes), nil
}
