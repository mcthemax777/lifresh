package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
)

type CreateFolderHandler struct {
	userService *service.UserService
}

func NewCreateFolderHandler(userService *service.UserService) *CreateFolderHandler {
	return &CreateFolderHandler{userService: userService}
}

func (h CreateFolderHandler) ApiCall(c *gin.Context) {
	_ = AuthHandlerFx[request.CreateFolderReq, response.BasicRes](c, &request.CreateFolderReq{}, func(accountID define.SnowflakeID, req *request.CreateFolderReq) (*response.BasicRes, *apperr.AppError) {

		folder := req.Folder
		folder.ID = 0
		folder, err := h.userService.CreateFolder(folder)

		if err != nil {
			return nil, err
		}

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.CREATE_FOLDER_RES)

		return res.(*response.BasicRes), nil
	})
}
