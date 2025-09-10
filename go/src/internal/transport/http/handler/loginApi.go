package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/auth"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
	"time"
)

type LoginHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewLoginHandler(authService *service.AuthService, userService *service.UserService) *LoginHandler {
	return &LoginHandler{authService: authService, userService: userService}
}

func (h LoginHandler) ApiCall(c *gin.Context) {
	_ = ApiHandlerFx[request.LoginReq, response.LoginRes](c, &request.LoginReq{}, func(req *request.LoginReq) (*response.LoginRes, *apperr.AppError) {

		account, err := h.authService.Auth(req.SocialToken, define.SocialType(req.SocialType))
		if err != nil {
			return nil, err
		}

		account, err = h.userService.Login(account)
		if err != nil {
			return nil, err
		}

		at, err := auth.SignAccess(account.ID, account.ProviderUID, 15*time.Minute)
		if err != nil {
			return nil, err
		}

		rt, _, err := auth.SignRefresh(account.ID, account.ProviderUID, 14*24*time.Hour)
		if err != nil {
			return nil, err
		}

		res := response.CreateSuccessResponse(response.LOGIN_RES)
		loginRes := res.(*response.LoginRes)
		loginRes.AccessToken = at
		loginRes.RefreshToken = rt
		loginRes.Account = account

		return loginRes, nil
	})
}
