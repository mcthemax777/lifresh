package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
)

type PingHandler struct {
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h PingHandler) ApiCall(c *gin.Context) {
	_ = AuthHandlerFx[request.PingReq, response.PingRes](c, &request.PingReq{}, func(accountID define.SnowflakeID, req *request.PingReq) (*response.PingRes, *apperr.AppError) {

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.PING_RES)

		pingRes := res.(*response.PingRes)
		pingRes.AccountID = accountID

		return pingRes, nil
	})
}
