package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
)

type DeletePlanHandler struct {
	userService *service.UserService
}

func NewDeletePlanHandler(userService *service.UserService) *DeletePlanHandler {
	return &DeletePlanHandler{userService: userService}
}

func (h DeletePlanHandler) ApiCall(c *gin.Context) {
	_ = AuthHandlerFx[request.DeletePlanReq, response.BasicRes](c, &request.DeletePlanReq{}, func(accountID define.SnowflakeID, req *request.DeletePlanReq) (*response.BasicRes, *apperr.AppError) {
		_, err := h.userService.DeletePlan(req.Plan)

		if err != nil {
			return nil, err
		}

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.DELETE_PLAN_RES)

		return res.(*response.BasicRes), nil
	})
}
