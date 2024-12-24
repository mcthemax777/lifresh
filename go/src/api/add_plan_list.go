package api

import (
	"encoding/json"
	"lifresh/db"
	"lifresh/models"
	"lifresh/request"
	"lifresh/response"
)

type AddPlanListHandler struct {
	SessionApiHandler
}

func NewAddPlanListHandler() AddPlanListHandler {
	h := AddPlanListHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h AddPlanListHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddPlanListReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	currentTime := CurrentTime()

	accountId, err := h.checkSession(req.Uid, req.Sid, currentTime)

	//세션 만료
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, err.Error())), err
	}

	user, err := db.DBHandlerSG.GetUserByAccountId(accountId)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, err.Error())), err
	}

	userId := user.Id

	var insertPlanList []models.Plan
	var updatePlanList []models.Plan

	for _, plan := range req.PlanList {
		plan.UserId = userId

		//신규 등록
		if plan.Id == 0 {
			insertPlanList = append(insertPlanList, plan)
		} else {
			//변경
			updatePlanList = append(updatePlanList, plan)
		}
	}

	if len(insertPlanList) > 0 {
		err = db.DBHandlerSG.InsertPlanList(&insertPlanList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertPlanList")), err
		}
	}

	if len(updatePlanList) > 0 {
		err = db.DBHandlerSG.UpdatePlanList(&updatePlanList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdatePlanList")), err
		}
	}

	var insertPlanRecordList []models.PlanRecord
	var updatePlanRecordList []models.PlanRecord

	for _, planRecord := range req.PlanRecordList {

		planRecord.UserId = userId

		//신규 등록
		if planRecord.Id == 0 {
			insertPlanRecordList = append(insertPlanRecordList, planRecord)
		} else {
			//변경
			updatePlanRecordList = append(updatePlanRecordList, planRecord)
		}
	}

	if len(insertPlanRecordList) > 0 {
		err = db.DBHandlerSG.InsertPlanRecordList(&insertPlanRecordList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertPlanRecordList")), err
		}
	}

	if len(updatePlanRecordList) > 0 {
		err = db.DBHandlerSG.UpdatePlanRecordList(&updatePlanRecordList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdatePlanRecordList")), err
		}
	}

	var insertPlanRecordOperatorList []models.PlanRecordOperator
	var updatePlanRecordOperatorList []models.PlanRecordOperator

	for _, planRecordOperator := range req.PlanRecordOperatorList {

		planRecordOperator.UserId = userId

		//신규 등록
		if planRecordOperator.Id == 0 {
			insertPlanRecordOperatorList = append(insertPlanRecordOperatorList, planRecordOperator)
		} else {
			//변경
			updatePlanRecordOperatorList = append(updatePlanRecordOperatorList, planRecordOperator)
		}
	}

	if len(insertPlanRecordOperatorList) > 0 {
		err = db.DBHandlerSG.InsertPlanRecordOperatorList(&insertPlanRecordOperatorList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertPlanRecordOperatorList")), err
		}
	}

	if len(updatePlanRecordOperatorList) > 0 {
		err = db.DBHandlerSG.UpdatePlanRecordOperatorList(&updatePlanRecordOperatorList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdatePlanRecordOperatorList")), err
		}
	}

	var insertPlanGoalList []models.PlanGoal
	var updatePlanGoalList []models.PlanGoal

	for _, planGoal := range req.PlanGoalList {

		planGoal.UserId = userId

		//신규 등록
		if planGoal.Id == 0 {
			insertPlanGoalList = append(insertPlanGoalList, planGoal)
		} else {
			//변경
			updatePlanGoalList = append(updatePlanGoalList, planGoal)
		}
	}

	if len(insertPlanGoalList) > 0 {
		err = db.DBHandlerSG.InsertPlanGoalList(&insertPlanGoalList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertPlanGoalList")), err
		}
	}

	if len(updatePlanGoalList) > 0 {
		err = db.DBHandlerSG.UpdatePlanGoalList(&updatePlanGoalList)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "UpdatePlanGoalList")), err
		}
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_PLAN_LIST_RES)

	return ResponseToByteArray(res.(*response.BasicRes)), nil
}
