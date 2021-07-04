package api

import (
	"db"
	"encoding/json"
	"request"
	"response"
)

type signUpHandler struct {
}

func (h signUpHandler) process(reqBody []byte) ([]byte, error) {
	var req request.SignUpReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	err = db.DBHandlerSG.InsertAccountAndPlanner(req.UserId, req.Password)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.SIGN_UP_RES)

	signUpRes := res. (*response.SignUpRes)

	return ResponseToByteArray(signUpRes), nil
}