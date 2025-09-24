package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
)

type DeleteFolderHandler struct {
	userService *service.UserService
}

func NewDeleteFolderHandler(userService *service.UserService) *DeleteFolderHandler {
	return &DeleteFolderHandler{userService: userService}
}

func (h DeleteFolderHandler) ApiCall(c *gin.Context) {
	_ = AuthHandlerFx[request.DeleteFolderReq, response.BasicRes](c, &request.DeleteFolderReq{}, func(accountID define.SnowflakeID, req *request.DeleteFolderReq) (*response.BasicRes, *apperr.AppError) {

		_, err := h.userService.DeleteFolder(req.Folder)

		if err != nil {
			return nil, err
		}

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.DELETE_FOLDER_RES)

		return res.(*response.BasicRes), nil
	})
}
