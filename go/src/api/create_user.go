package api

//
//import (
//	"encoding/json"
//	"lifresh/db"
//	"lifresh/request"
//	"lifresh/response"
//)
//
//type CreateUserHandler struct {
//	SessionApiHandler
//}
//
//func NewCreateUserHandler() CreateUserHandler {
//	h := CreateUserHandler{SessionApiHandler: NewSessionApiHandler()}
//	return h
//}
//
//func (h CreateUserHandler) process(reqBody []byte) ([]byte, error) {
//
//	var req request.CreateUserReq
//	err := json.Unmarshal(reqBody, &req)
//
//	if err != nil {
//		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
//	}
//
//	currentTime := CurrentTime()
//
//	accountId, err := h.checkSession(req.Uid, req.Sid, currentTime)
//
//	//세션 만료
//	if err != nil {
//		return ResponseToByteArray(response.CreateFailResponse(201, err.Error())), err
//	}
//
//	_, err = db.DBHandlerSG.InsertUserAndRootFolder(accountId, req.Nickname, req.RootFolderName)
//
//	if err != nil {
//		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
//	}
//
//	//전송할 데이터 만들기
//	res := response.CreateSuccessResponse(response.CREATE_USER_RES)
//
//	sendRes := res.(*response.CreateUserRes)
//
//	return ResponseToByteArray(sendRes), nil
//}
