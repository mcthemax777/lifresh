package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"lifresh/auth"
	"lifresh/request"
	"lifresh/response"
	"time"
)

type RefreshHandler struct {
}

func (h RefreshHandler) ApiCall(c *gin.Context) {
	ApiCall(c, h.process)
}

func (h RefreshHandler) process(reqBody []byte) ([]byte, error) {
	var req request.RefreshReq
	err := json.Unmarshal(reqBody, &req)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	cl, err := auth.ParseRefresh(req.RefreshToken)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_refresh")), err
	}

	at, err := auth.SignAccess(cl.AccountID, cl.UID, 15*time.Minute)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "sign_access_error invalid")), err
	}

	// (권장) RT 회전
	newRT, _, err := auth.SignRefresh(cl.AccountID, cl.UID, 14*24*time.Hour)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "sign_refresh_error invalid")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.REFRESH_RES)

	refreshRes := res.(*response.RefreshRes)
	refreshRes.AccessToken = at
	refreshRes.RefreshToken = newRT

	return ResponseToByteArray(refreshRes), nil
}
