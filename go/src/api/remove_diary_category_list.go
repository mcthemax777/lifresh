package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type RemoveDiaryCategoryListHandler struct {
	SessionApiHandler
}

func NewRemoveDiaryCategoryListHandler() RemoveDiaryCategoryListHandler {
	h := RemoveDiaryCategoryListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h RemoveDiaryCategoryListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.RemoveDiaryCategoryListReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	currentTime := CurrentTime()

	accountId, err := h.checkSession(req.Uid, req.Sid, currentTime)

	//세션 만료
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "session invalid")), err
	}

	diary, err := db.DBHandlerSG.GetDiaryByAccountId(accountId)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "plannerNo not")), err
	}

	err = db.DBHandlerSG.DeleteDiaryCategoryList(diary.DiaryId, req.DiaryCategoryIdList)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "InsertCF")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.REMOVE_DIARY_CATEGORY_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}