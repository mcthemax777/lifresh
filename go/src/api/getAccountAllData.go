package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type GetAccountAllDataHandler struct {
	SessionApiHandler
}

func NewGetAccountAllDataHandler() GetAccountAllDataHandler {
	h := GetAccountAllDataHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h GetAccountAllDataHandler) process(reqBody []byte) ([]byte, error) {

	var req request.GetAccountAllDataReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	currentTime := CurrentTime()

	accountId, err := h.checkSession(req.Uid, req.Sid, currentTime)

	//세션 만료
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, err.Error())), err
	}

	planner, err := db.DBHandlerSG.GetPlannerByAccountId(accountId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, err.Error())), err
	}
	planCategoryList, err := db.DBHandlerSG.GetPlanCategoryListByPlannerId(planner.PlannerId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, err.Error())), err
	}
	planList, err := db.DBHandlerSG.GetPlanListByPlannerId(planner.PlannerId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, err.Error())), err
	}
	planHistoryList, err := db.DBHandlerSG.GetPlanHistoryListByPlannerId(planner.PlannerId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, err.Error())), err
	}

	money, err := db.DBHandlerSG.GetMoneyByAccountId(accountId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, err.Error())), err
	}
	moneyCategoryList, err := db.DBHandlerSG.GetMoneyCategoryListByMoneyId(money.MoneyId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, err.Error())), err
	}
	moneyHistoryList, err := db.DBHandlerSG.GetMoneyHistoryListByMoneyId(money.MoneyId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, err.Error())), err
	}

	diary, err := db.DBHandlerSG.GetDiaryByAccountId(accountId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, err.Error())), err
	}
	diaryCategoryList, err := db.DBHandlerSG.GetDiaryCategoryListByMoneyId(diary.DiaryId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, err.Error())), err
	}
	diaryHistoryList, err := db.DBHandlerSG.GetDiaryHistoryListByMoneyId(diary.DiaryId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, err.Error())), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_ACCOUNT_ALL_DATA_RES)

	sendRes := res.(*response.GetAccountAllData)
	sendRes.Planner = planner
	sendRes.PlanCategoryList = planCategoryList
	sendRes.PlanList = planList
	sendRes.PlanHistoryList = planHistoryList
	sendRes.Money = money
	sendRes.MoneyCategoryList = moneyCategoryList
	sendRes.MoneyHistoryList = moneyHistoryList
	sendRes.Diary = diary
	sendRes.DiaryCategoryList = diaryCategoryList
	sendRes.DiaryHistoryList = diaryHistoryList

	return ResponseToByteArray(sendRes), nil
}
