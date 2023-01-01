package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/request"
	"lifresh/response"
)

type GetScheduleTaskHandler struct {
	SessionApiHandler
}

func NewGetScheduleTaskHandler() GetScheduleTaskHandler {
	h := GetScheduleTaskHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h GetScheduleTaskHandler) process(reqBody []byte) ([]byte, error) {

	var req request.GetScheduleTaskReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	currentTime := CurrentTime()

	accountNo, err := h.checkSession(req.Uid, req.Sid, currentTime)

	//세션 만료
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "session invalid")), err
	}

	planner, err := db.DBHandlerSG.GetPlannerByAccountNo(accountNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "plannerNo not")), err
	}

	todayList, err := db.DBHandlerSG.GetTodayListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "GetTodayListByPlannerNo")), err
	}

	scheduleTaskList, err := db.DBHandlerSG.GetScheduleTaskListByPlannerNo(planner.PlannerNo)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "planner db error")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_SCHEDULE_TASK_RES)

	sendRes := res.(*response.GetScheduleTaskRes)
	sendRes.Planner = planner
	sendRes.TodayList = todayList
	sendRes.ScheduleTaskList = scheduleTaskList

	return ResponseToByteArray(sendRes), nil
}
