package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/models"
	"lifresh/request"
	"lifresh/response"
)

type AddDiaryHistoryListHandler struct {
	SessionApiHandler
}

func NewAddDiaryHistoryListHandler() AddDiaryHistoryListHandler {
	h := AddDiaryHistoryListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h AddDiaryHistoryListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddDiaryHistoryListReq
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

	var insertList []models.DiaryHistory
	var updateList []models.DiaryHistory

	for _, diaryHistory := range req.DiaryHistoryList {

		diaryHistory.DiaryId = diary.DiaryId

		//신규 등록
		if diaryHistory.DiaryCategoryId == 0 {

			insertList = append(insertList, diaryHistory)

		} else { //변경

			updateList = append(updateList, diaryHistory)

		}
	}

	if len(insertList) > 0 {
		err = db.DBHandlerSG.InsertDiaryHistoryList(&insertList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertDiaryCategoryList")), err
		}
	}

	if len(updateList) > 0 {
		err = db.DBHandlerSG.UpdateDiaryHistoryList(&updateList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdateDiaryCategoryList")), err
		}
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_DIARY_HISTORY_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
