package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type RemoveSubCategoryListHandler struct {
	SessionApiHandler
}

func NewRemoveSubCategoryListHandler() RemoveSubCategoryListHandler {
	h := RemoveSubCategoryListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h RemoveSubCategoryListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.RemoveSubCategoryListReq
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

	_, err = db.DBHandlerSG.DeleteSubCategoryList(planner.PlannerNo, req.SubCategoryNoList)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "InsertCF")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.REMOVE_SUB_CATEGORY_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
