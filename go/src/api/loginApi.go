package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"lifresh/db"
	"lifresh/redis"
	"lifresh/request"
	"lifresh/response"
	"strconv"
	"strings"
)

type LoginHandler struct {
}

func (h LoginHandler) process(reqBody []byte) ([]byte, error) {

	var req request.LoginReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	account, err := db.DBHandlerSG.GetAccountBySocialToken(req.SocialType, req.SocialToken)

	if err != nil {
		if account.SocialType == 0 {
			return ResponseToByteArray(response.CreateFailResponse(202, "invalid_parameter")), err
		}

		account, err = db.DBHandlerSG.InsertAccount(req.SocialType, req.SocialToken)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(202, "insert error")), err
		}
	}

	uid := strconv.Itoa(req.SocialType) + "-" + req.SocialToken
	//session id 생성
	sid := uuid.New().String()
	sid = strings.Replace(sid, "-", "", -1)

	err = redis.RedisHandlerSG.SetSession(uid, sid, account.AccountId)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(301, "redis error")), err
	}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.LOGIN_RES)

	loginRes := res.(*response.LoginRes)
	//임시로 userId 넣어줌(나중에 유니크한 아이디 생성해서 전달)
	loginRes.Uid = uid
	loginRes.Sid = sid
	loginRes.Account = account

	return ResponseToByteArray(loginRes), nil
}
