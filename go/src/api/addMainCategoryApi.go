package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/models"
	"lifresh/request"
	"lifresh/response"
)

type AddMainCategoryListHandler struct {
	SessionApiHandler
}

func NewAddMainCategoryListHandler() AddMainCategoryListHandler {
	h := AddMainCategoryListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h AddMainCategoryListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddMainCategoryListReq
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

	var insertList []models.MainCategory
	var updateList []models.MainCategory

	for _, mainCategory := range req.MainCategoryList {

		mainCategory.PlannerNo = planner.PlannerId

		//신규 등록
		if mainCategory.MainCategoryNo == 0 {
			insertList = append(insertList, mainCategory)
		} else { //변경
			updateList = append(updateList, mainCategory)
		}
	}

	if len(insertList) > 0 {
		err = db.DBHandlerSG.InsertMainCategory(&insertList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertMainCategoryList")), err
		}
	}

	if len(updateList) > 0 {
		err = db.DBHandlerSG.UpdateMainCategory(&updateList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdateMainCategoryList")), err
		}
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_MAIN_CATEGORY_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
