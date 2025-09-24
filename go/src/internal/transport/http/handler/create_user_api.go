package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lifresh/define"
	"lifresh/internal/core"
	"lifresh/internal/core/apperr"
	"lifresh/internal/domain"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
)

type CreateUserHandler struct {
	userService *service.UserService
}

func NewCreateUserHandler(userService *service.UserService) *CreateUserHandler {
	return &CreateUserHandler{userService: userService}
}

func (h CreateUserHandler) ApiCall(c *gin.Context) {
	_ = AuthHandlerFx[request.CreateUserReq, response.CreateUserRes](c, &request.CreateUserReq{}, func(accountID define.SnowflakeID, req *request.CreateUserReq) (*response.CreateUserRes, *apperr.AppError) {

		nickname := req.User.Nickname
		bio := req.User.Bio
		profileURL := req.User.ProfileURL

		if err := ValidateNickname(nickname); err != nil {
			fmt.Println(err) // nil → 유효
		}

		user := &domain.User{
			ID:         0,
			AccountID:  accountID,
			Nickname:   nickname,
			Bio:        bio,
			ProfileURL: profileURL,
			CreatedAt:  core.Now(),
			UpdatedAt:  core.Now(),
		}

		user, err := h.userService.CreateUser(user)

		if err != nil {
			return nil, err
		}

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.CREATE_USER_RES)
		createUserRes := res.(*response.CreateUserRes)
		createUserRes.User = user

		return createUserRes, nil
	})
}
