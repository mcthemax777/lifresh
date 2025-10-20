package handler

import (
	"github.com/gin-gonic/gin"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/service"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
)

type CreatePlanHandler struct {
	userService *service.UserService
}

func NewCreatePlanHandler(userService *service.UserService) *CreatePlanHandler {
	return &CreatePlanHandler{userService: userService}
}

func (h CreatePlanHandler) ApiCall(c *gin.Context) {
	_ = AuthHandlerFx[request.CreatePlanReq, response.BasicRes](c, &request.CreatePlanReq{}, func(accountID define.SnowflakeID, req *request.CreatePlanReq) (*response.BasicRes, *apperr.AppError) {

		plan := req.Plan
		plan.ID = 0
		plan, err := h.userService.CreatePlan(plan)

		if err != nil {
			return nil, err
		}

		//전송할 데이터 만들기
		res := response.CreateSuccessResponse(response.CREATE_PLAN_RES)

		return res.(*response.BasicRes), nil
	})
}
