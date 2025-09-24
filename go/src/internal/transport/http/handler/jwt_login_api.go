package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/auth"
	"lifresh/internal/core/apperr"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
	"time"
)

type JwtLoginHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewJwtLoginHandler(authService *service.AuthService, userService *service.UserService) *JwtLoginHandler {
	return &JwtLoginHandler{authService: authService, userService: userService}
}

func (h JwtLoginHandler) ApiCall(c *gin.Context) {
	_ = ApiHandlerFx[request.RefreshReq, response.LoginRes](c, &request.RefreshReq{}, func(req *request.RefreshReq) (*response.LoginRes, *apperr.AppError) {

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

		account, err := h.userService.GetAccountByID(cl.AccountID)
		if err != nil {
			return nil, err
		}

		res := response.CreateSuccessResponse(response.LOGIN_RES)
		loginRes := res.(*response.LoginRes)
		loginRes.AccessToken = at
		loginRes.RefreshToken = newRT
		loginRes.Account = account

		return loginRes, nil
	})
}
