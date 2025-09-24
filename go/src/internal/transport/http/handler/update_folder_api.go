package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
)

type UpdateFolderHandler struct {
	userService *service.UserService
}

func NewUpdateFolderHandler(userService *service.UserService) *UpdateFolderHandler {
	return &UpdateFolderHandler{userService: userService}
}

func (h UpdateFolderHandler) ApiCall(c *gin.Context) {
	_ = AuthHandlerFx[request.UpdateFolderReq, response.BasicRes](c, &request.UpdateFolderReq{}, func(accountID define.SnowflakeID, req *request.UpdateFolderReq) (*response.BasicRes, *apperr.AppError) {

		folder := req.Folder
		folder, err := h.userService.UpdateFolder(folder)

		if err != nil {
			return nil, err
		}

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.UPDATE_FOLDER_RES)

		return res.(*response.BasicRes), nil
	})
}
