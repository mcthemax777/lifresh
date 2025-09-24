package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/domain"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
)

type GetUserHandler struct {
	userService *service.UserService
}

func NewGetUserHandler(userService *service.UserService) *GetUserHandler {
	return &GetUserHandler{userService: userService}
}

func (h GetUserHandler) ApiCall(c *gin.Context) {
	_ = AuthHandlerFx[request.GetUserReq, response.GetUserRes](c, &request.GetUserReq{}, func(accountID define.SnowflakeID, req *request.GetUserReq) (*response.GetUserRes, *apperr.AppError) {
		var userID define.SnowflakeID
		_ = userID.Scan(req.UserID)

		var user *domain.User
		var err *apperr.AppError
		if userID == 0 {
			accountID, _ := c.Get("account_id")
			user, err = h.userService.GetUserByAccountID(accountID.(define.SnowflakeID))
		} else {
			user, err = h.userService.GetUser(userID)
		}
		if err != nil {
			return nil, err
		}

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.GET_USER_RES)

		getUserRes := res.(*response.GetUserRes)
		getUserRes.User = user

		return getUserRes, nil
	})
}
