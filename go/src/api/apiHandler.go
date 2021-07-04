package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lflog"
	"net/http"
	"redis"
	"response"
	"time"

	"github.com/gin-gonic/gin"
)

var handlerMap map[string]apiHandler

func init() {
	handlerMap = make(map[string]apiHandler)
	fmt.Println("api init all")
	handlerMap["login"] = loginHandler{}
	handlerMap["signUp"] = signUpHandler{}
	handlerMap["getToday"] = NewGetTodayHandler()
	handlerMap["getCFList"] = NewGetCFListHandler()
	handlerMap["addCF"] = NewAddCFHandler()
	handlerMap["addTaskPlan"] = NewAddTaskPlanHandler()
	handlerMap["addTask"] = NewAddTaskHandler()
}

func ApiCall(c *gin.Context) {
	handler := handlerMap[c.Param("name")]

	if handler == nil {
		return 
	}

	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		lflog.Logging(1, err.Error())
		c.String(http.StatusOK, string(ResponseToByteArray(response.CreateFailResponse(301, "body error"))))
		return
	}

	//받은 데이터 출력
	lflog.Logging(1, string(body))

	//기본 세팅(현재 시간, 유저 정보 등등...)
	

	//로직 실행
	res, err := handler.process(body)

	if err != nil {
		lflog.Logging(1, err.Error())
	}

	lflog.Logging(1, string(res))

	c.String(http.StatusOK, string(res))
}

type apiHandler interface {
	process(reqBody []byte) ([]byte, error)
}

func ResponseToByteArray(res response.Response) []byte {
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

func (sah *SessionApiHandler) checkSession(sid string, currentTime time.Time) (userNo int, err error) {

	sessionInfo, err := redis.RedisHandlerSG.GetSession(sid)

	if err != nil {
		return 0, nil
	}

	if sessionInfo.ExpireTime.Before(currentTime) {
		return 0, nil
	}

	return sessionInfo.AccountNo, nil
}

func CurrentTime() time.Time {
	return time.Now()
}