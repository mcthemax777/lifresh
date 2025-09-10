package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/auth"
	"lifresh/internal/core/apperr"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
	"time"
)

type RefreshHandler struct {
}

func NewRefreshHandler() *RefreshHandler {
	return &RefreshHandler{}
}

func (h RefreshHandler) ApiCall(c *gin.Context) {
	_ = ApiHandlerFx[request.RefreshReq, response.RefreshRes](c, &request.RefreshReq{}, func(req *request.RefreshReq) (*response.RefreshRes, *apperr.AppError) {

		cl, err := auth.ParseRefresh(req.RefreshToken)
		if err != nil {
			return nil, err
		}

		at, err := auth.SignAccess(cl.AccountID, cl.UID, 15*time.Minute)
		if err != nil {
			return nil, err
		}

		// (권장) RT 회전
		newRT, _, err := auth.SignRefresh(cl.AccountID, cl.UID, 14*24*time.Hour)
		if err != nil {
			return nil, err
		}

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.REFRESH_RES)

		refreshRes := res.(*response.RefreshRes)
		refreshRes.AccessToken = at
		refreshRes.RefreshToken = newRT

		return refreshRes, nil
	})
}
