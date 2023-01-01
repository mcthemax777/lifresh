package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type GetSubCategoryHandler struct {
	SessionApiHandler
}

func NewGetSubCategoryHandler() GetSubCategoryHandler {
	h := GetSubCategoryHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h GetSubCategoryHandler) process(reqBody []byte) ([]byte, error) {

	var req request.GetSubCategoryReq
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

	todayList, err := db.DBHandlerSG.GetTodayListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "GetTodayListByPlannerNo")), err
	}

	subCategoryList, err := db.DBHandlerSG.GetSubCategoryListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "GetSubCategoryListByPlannerNo")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_SUB_CATEGORY_RES)

	sendRes := res.(*response.GetSubCategoryRes)
	sendRes.Planner = planner
	sendRes.TodayList = todayList
	sendRes.SubCategoryList = subCategoryList

	return ResponseToByteArray(sendRes), nil
}
