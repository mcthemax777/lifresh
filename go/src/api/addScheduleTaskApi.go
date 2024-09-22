package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/models"
	"lifresh/request"
	"lifresh/response"
)

type AddScheduleTaskListHandler struct {
	SessionApiHandler
}

func NewAddScheduleTaskListHandler() AddScheduleTaskListHandler {
	h := AddScheduleTaskListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h AddScheduleTaskListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddScheduleTaskListReq
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

	var insertList []models.ScheduleTask
	var updateList []models.ScheduleTask

	for _, scheduleTask := range req.ScheduleTaskList {

		scheduleTask.PlannerNo = planner.PlannerId

		//신규 등록
		if scheduleTask.ScheduleTaskNo == 0 {

			insertList = append(insertList, scheduleTask)

		} else { //변경

			updateList = append(updateList, scheduleTask)
		}
	}

	if len(insertList) > 0 {
		err = db.DBHandlerSG.InsertScheduleTask(&insertList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertScheduleTask")), err
		}
	}

	if len(updateList) > 0 {
		err = db.DBHandlerSG.UpdateScheduleTask(&updateList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdateScheduleTask")), err
		}
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_SCHEDULE_TASK_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
