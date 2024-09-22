package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/models"
	"lifresh/request"
	"lifresh/response"
)

type AddSubCategoryListHandler struct {
	SessionApiHandler
}

func NewAddSubCategoryListHandler() AddSubCategoryListHandler {
	h := AddSubCategoryListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h AddSubCategoryListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddSubCategoryListReq
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

	var insertList []models.SubCategory
	var updateList []models.SubCategory

	for _, subCategory := range req.SubCategoryList {

		subCategory.PlannerNo = planner.PlannerId

		//신규 등록
		if subCategory.MainCategoryNo == 0 {

			insertList = append(insertList, subCategory)

		} else { //변경

			updateList = append(updateList, subCategory)

		}
	}

	if len(insertList) > 0 {
		err = db.DBHandlerSG.InsertSubCategory(&insertList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertSubCategory")), err
		}
	}

	if len(updateList) > 0 {
		err = db.DBHandlerSG.UpdateSubCategory(&updateList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdateSubCategory")), err
		}
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_SUB_CATEGORY_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
