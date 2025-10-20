package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
)

type UpdateUserHandler struct {
	userService *service.UserService
}

func NewUpdateUserHandler(userService *service.UserService) *UpdateUserHandler {
	return &UpdateUserHandler{userService: userService}
}

func (h UpdateUserHandler) ApiCall(c *gin.Context) {
	_ = AuthHandlerFx[request.UpdateUserReq, response.UpdateUserRes](c, &request.UpdateUserReq{}, func(accountID define.SnowflakeID, req *request.UpdateUserReq) (*response.UpdateUserRes, *apperr.AppError) {

		user := req.User
		user.AccountID = accountID
		newUser, err := h.userService.UpdateUser(&user)

		if err != nil {
			return nil, err
		}

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.UPDATE_USER_RES)
		updateUserRes := res.(*response.UpdateUserRes)
		updateUserRes.User = newUser

		return updateUserRes, nil
	})
}
