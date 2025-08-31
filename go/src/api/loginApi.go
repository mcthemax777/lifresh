package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"lifresh/auth"
	"lifresh/db"
	"lifresh/define"
	"lifresh/request"
	"lifresh/response"
	"os"
	"strings"
	"time"
)

type LoginHandler struct {
}

func (h LoginHandler) ApiCall(c *gin.Context) {
	ApiCall(c, h.process)
}

func (h LoginHandler) process(reqBody []byte) ([]byte, error) {

	var req request.LoginReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	// 1) Google idToken이 오면 OIDC로 검증
	var uid string
	var email string
	var name string
	//// 로그인 핸들러 내부 (요지)
	switch {
	case req.SocialType == define.SocialTypeGuest:
		// 2) 일반/게스트 로그인
		if req.SocialToken == "" {
			uid = "guest-" + uuid.New().String()
		} else {
			uid = req.SocialToken
		}

		email = ""
		name = uid
	case req.SocialType == define.SocialTypeGoogle:
		ids := strings.Split(os.Getenv("GOOGLE_OAUTH_CLIENT_IDS"), ",")
		ctx, cancel := auth.ContextWithTimeout(5)
		defer cancel()
		claims, err := auth.VerifyGoogleIDTokenAllowList(ctx, req.SocialToken, ids)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "invalid_google_idtoken")), err
		}
		uid = "google-" + claims.Subject
		email = claims.Email
		name = claims.Name
		//picture := claims.Picture

	case req.SocialType == define.SocialTypeApple:
		// 여러 client_id 허용 (iOS 번들ID, 웹 Service ID 등)
		ids := strings.Split(os.Getenv("GOOGLE_OAUTH_CLIENT_IDS"), ",")
		ctx, cancel := auth.ContextWithTimeout(5)
		defer cancel()
		ac, err := auth.VerifyAppleIDTokenAllowList(ctx, req.SocialToken, ids, "")
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(201, "invalid_apple_idtoken")), err
		}
		uid = "apple-" + ac.Subject
		email = ac.Email
		name = "apple_none"
	}

	account, err := db.DBHandlerSG.GetAccountByUID(uid)

	if err != nil {
		//게스트 로그인인데 토큰 보냈으면 존재하는 유저라고 판단해야됨
		if req.SocialType == define.SocialTypeGuest && req.SocialToken != "" {
			return ResponseToByteArray(response.CreateFailResponse(202, "select error")), err
		}

		account, err = db.DBHandlerSG.InsertAccount(req.SocialType, uid, name, email)
		if err != nil {
			return ResponseToByteArray(response.CreateFailResponse(202, "insert error")), err
		}
	}

	at, err := auth.SignAccess(account.ID, uid, 15*time.Minute)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "access_token invalid")), err
	}
	rt, _, err := auth.SignRefresh(account.ID, uid, 14*24*time.Hour)
	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "refresh_token invalid")), err
	}

	//uid := strconv.Itoa(req.SocialType) + "-" + req.SocialToken
	////session id 생성
	//sid := uuid.New().String()
	//sid = strings.Replace(sid, "-", "", -1)
	//err = redis.RedisHandlerSG.SetSession(uid, sid, account.Id)
	//if err != nil {
	//	return ResponseToByteArray(response.CreateFailResponse(301, "redis error")), err
	//}

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.LOGIN_RES)

	loginRes := res.(*response.LoginRes)
	loginRes.AccessToken = at
	loginRes.RefreshToken = rt
	//loginRes.Sid = sid
	loginRes.Account = account

	return ResponseToByteArray(loginRes), nil
}
