package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
)

type GetUserHandler struct {
}

func NewGetUserHandler() *GetUserHandler {
	return &GetUserHandler{}
}

func (h GetUserHandler) ApiCall(c *gin.Context) {
	_ = AuthHandlerFx[request.GetUserReq, response.PingRes](c, &request.GetUserReq{}, func(accountID define.SnowflakeID, req *request.GetUserReq) (*response.PingRes, *apperr.AppError) {

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.PING_RES)

		pingRes := res.(*response.PingRes)
		pingRes.AccountID = accountID

		return pingRes, nil
	})
}
