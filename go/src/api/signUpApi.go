package api

import (
	"encoding/json"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
	"lifresh/internal/transport/http/handler"
)

type SignUpHandler struct {
}

func (h SignUpHandler) process(reqBody []byte) ([]byte, error) {
	var req request.SignUpReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return handler.ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	//err, _ = db.DBHandlerSG.InsertUserAndRootFolder(req.Nickname, req.RootFolderName)

	if err != nil {
		return handler.ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.SIGN_UP_RES)

	signUpRes := res.(*response.SignUpRes)

	return handler.ResponseToByteArray(signUpRes), nil
}
