package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"lifresh/define"
	"lifresh/internal/core/apperr"
	"lifresh/internal/transport/http/dto/response"
	"lifresh/lflog"
	"lifresh/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthHandlerFx[T any, R any](
	c *gin.Context,
	req *T,
	process func(accountID define.SnowflakeID, req *T) (*R, *apperr.AppError),
) error {
	return ApiHandlerFx(c, req, func(req1 *T) (*R, *apperr.AppError) {
		//accountID 가져오기
		accountID, _ := c.Get("account_id")

		//로직 실행
		return process(accountID.(define.SnowflakeID), req1)
	})
}

func ApiHandlerFx[T any, R any](
	c *gin.Context,
	req *T,
	process func(req *T) (*R, *apperr.AppError),
) error {
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		lflog.Logging(lflog.LogLevelInfo, err.Error())
		c.String(http.StatusOK, string(ResponseToByteArray(response.CreateFailResponse(301, "body error"))))
		return err
	}

	//받은 데이터 출력
	lflog.Logging(lflog.LogLevelInfo, string(body))

	//request dto 로 파싱
	err = json.Unmarshal(body, &req)

	if err != nil {
		c.String(http.StatusOK, string(ResponseToByteArray(response.CreateFailResponse(201, "invalid_json"))))
		return err
	}

	//로직 실행
	res, appError := process(req)

	if appError != nil {
		lflog.Logging(lflog.LogLevelInfo, appError.Error())
		c.String(http.StatusOK, string(ResponseToByteArray(response.CreateFailResponse(appError.Code, appError.Message))))
		return appError.Err
	}

	resBytes := ResponseToByteArray(res)

	resultLog := "{\"input\":" + string(body) + ", \"output\":" + string(resBytes) + "}"

	lflog.Logging(lflog.LogLevelInfo, resultLog)

	c.String(http.StatusOK, string(resBytes))

	return nil
}

func ApiCall(c *gin.Context, process func(b []byte) ([]byte, error)) {

	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		lflog.Logging(lflog.LogLevelInfo, err.Error())
		c.String(http.StatusOK, string(ResponseToByteArray(response.CreateFailResponse(301, "body error"))))
		return
	}

	//받은 데이터 출력
	lflog.Logging(lflog.LogLevelInfo, string(body))

	//로직 실행
	res, err := process(body)

	if err != nil {
		lflog.Logging(lflog.LogLevelInfo, err.Error())
	}

	resultLog := "{\"input\":" + string(body) + ", \"output\":" + string(res) + "}"

	lflog.Logging(lflog.LogLevelInfo, resultLog)

	c.String(http.StatusOK, string(res))
}

func ResponseToByteArray(res any) []byte {
	result, _ := json.Marshal(res)

	return result
}

type SessionApiHandler struct {
	CurrentTime time.Time
}

func NewSessionApiHandler() SessionApiHandler {
	sah := SessionApiHandler{}
	sah.CurrentTime = time.Now()

	return sah
}

func (sah *SessionApiHandler) checkSession(uid string, sid string, currentTime time.Time) (userNo int, err error) {

	sessionInfo, err := redis.RedisHandlerSG.GetSession(uid)

	if err != nil {
		return 0, err
	}

	//세션이 다르면 err
	if sessionInfo.Sid != sid {
		return 0, errors.New("session mismatch")
	}

	//세션은 남아있는데 만료시간이 넘었다면 nil이 아닌 다른 err로 보내줘야됨
	if sessionInfo.ExpireTime.Before(currentTime) {
		return 0, errors.New("session expired")
	}

	//세션 기간 재설정
	err = redis.RedisHandlerSG.SetSession(uid, sid, sessionInfo.AccountId)

	if err != nil {
		return 0, err
	}

	return sessionInfo.AccountId, nil
}

func CurrentTime() time.Time {
	return time.Now()
}
