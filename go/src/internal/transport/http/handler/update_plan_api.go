package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
)

type UpdatePlanHandler struct {
	userService *service.UserService
}

func NewUpdatePlanHandler(userService *service.UserService) *UpdatePlanHandler {
	return &UpdatePlanHandler{userService: userService}
}

func (h UpdatePlanHandler) ApiCall(c *gin.Context) {
	_ = AuthHandlerFx[request.UpdatePlanReq, response.BasicRes](c, &request.UpdatePlanReq{}, func(accountID define.SnowflakeID, req *request.UpdatePlanReq) (*response.BasicRes, *apperr.AppError) {

		plan := req.Plan
		plan, err := h.userService.UpdatePlan(plan)

		if err != nil {
			return nil, err
		}

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.UPDATE_PLAN_RES)

		return res.(*response.BasicRes), nil
	})
}
