package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/models"
	"lifresh/request"
	"lifresh/response"
)

type AddDiaryCategoryListHandler struct {
	SessionApiHandler
}

func NewAddDiaryCategoryListHandler() AddDiaryCategoryListHandler {
	h := AddDiaryCategoryListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h AddDiaryCategoryListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddDiaryCategoryListReq
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

	diary, err := db.DBHandlerSG.GetDiaryByAccountId(accountId)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "diary not")), err
	}

	var insertList []models.DiaryCategory
	var updateList []models.DiaryCategory

	for _, diaryCategory := range req.DiaryCategoryList {

		diaryCategory.DiaryId = diary.DiaryId

		//신규 등록
		if diaryCategory.DiaryCategoryId == 0 {

			insertList = append(insertList, diaryCategory)

		} else { //변경

			updateList = append(updateList, diaryCategory)

		}
	}

	if len(insertList) > 0 {
		err = db.DBHandlerSG.InsertDiaryCategoryList(&insertList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertDiaryCategoryList")), err
		}
	}

	if len(updateList) > 0 {
		err = db.DBHandlerSG.UpdateDiaryCategoryList(&updateList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdateDiaryCategoryList")), err
		}
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_DIARY_CATEGORY_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
