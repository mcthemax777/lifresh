package api

import (
	"db"
	"encoding/json"
	"request"
	"response"
)

type RemoveCFHandler struct {
	SessionApiHandler
}

func NewRemoveCFHandler() RemoveCFHandler {
	h := RemoveCFHandler{SessionApiHandler: NewSessionApiHandler()}
	return h
}

func (h RemoveCFHandler) process(reqBody []byte) ([]byte, error) {

	var req request.AddCFListReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	currentTime := CurrentTime()

	accountNo, err := h.checkSession(req.Sid, currentTime)
	
	//세션 만료
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "session invalid")), err
	}


	planner, err := db.DBHandlerSG.GetPlannerByAccountNo(accountNo)
		
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "plannerNo not")), err
	}


	
	if(req.CF.CfNo == 0) {
		
		cfNo, err := db.DBHandlerSG.InsertCF(req.CF.Name, planner.PlannerNo)
		
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertCF")), err
		}

		req.CF.CfNo = cfNo
	}

	if(req.MCF.McfNo == 0) {
		
		mcfNo, err := db.DBHandlerSG.InsertMCF(req.MCF.Name, req.CF.CfNo, planner.PlannerNo)
		
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertMCF")), err
		}

		req.MCF.McfNo = mcfNo
	}

	if(req.DCF.DcfNo == 0) {
		
		dcfNo, err := db.DBHandlerSG.InsertDCF(req.DCF.Name, req.MCF.McfNo, req.DCF.DcfType, req.DCF.Priority, planner.PlannerNo)
		
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "InsertDCF")), err
		}

		req.DCF.DcfNo = dcfNo
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.ADD_CF_RES)

	loginRes := res. (*response.AddCFRes)
	loginRes.CF = req.CF
	loginRes.MCF = req.MCF
	loginRes.DCF = req.DCF
	
	
	return ResponseToByteArray(loginRes), nil
}
