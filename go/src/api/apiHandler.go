package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"lifresh/define"
	"lifresh/lflog"
	"lifresh/redis"
	"lifresh/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var HandlerMap map[string]BaseHandler

//var logError error
//var logger *fluent.Fluent

func init() {
	HandlerMap = make(map[string]BaseHandler)
	HandlerMap[define.ApiLogin] = LoginHandler{}
	HandlerMap[define.ApiRefresh] = RefreshHandler{}
	//HandlerMap[define.ApiRoot] = GetFolderHandler{}
	//handlerMap["signUp"] = SignUpHandler{}
	//handlerMap["createUser"] = NewCreateUserHandler()
	//handlerMap["getAccountAllData"] = NewGetAccountAllDataHandler()
	//handlerMap["addPlanList"] = NewAddPlanListHandler()
	//handlerMap["addPlanRecordList"] = NewGetAccountAllDataHandler()
	//handlerMap["addPlanRecordOperatorList"] = NewGetAccountAllDataHandler()
	//handlerMap["addPlanGoalList"] = NewGetAccountAllDataHandler()
	//handlerMap["removePlanList"] = NewGetAccountAllDataHandler()
	//handlerMap["removePlanRecordList"] = NewGetAccountAllDataHandler()
	//handlerMap["removePlanRecordOperatorList"] = NewGetAccountAllDataHandler()
	//handlerMap["removePlanGoalList"] = NewGetAccountAllDataHandler()
	//handlerMap["addDiaryCategoryList"] = NewAddDiaryCategoryListHandler()
	//handlerMap["addDiaryHistoryList"] = NewAddDiaryHistoryListHandler()
	//handlerMap["removeDiaryCategoryList"] = NewRemoveDiaryCategoryListHandler()
	//handlerMap["removeDiaryHistoryList"] = NewRemoveDiaryHistoryListHandler()
	//handlerMap["getUserAllData"] = NewGetUserAllDataHandler()
	//handlerMap["getMainCategoryList"] = NewGetMainCategoryHandler()
	//handlerMap["getSubCategoryList"] = NewGetSubCategoryHandler()
	//handlerMap["getScheduleTaskList"] = NewGetScheduleTaskHandler()
	//handlerMap["getToDoTaskList"] = NewGetToDoTaskHandler()
	//handlerMap["getMoneyTaskList"] = NewGetMoneyTaskHandler()
	//handlerMap["addMainCategoryList"] = NewAddMainCategoryListHandler()
	//handlerMap["addSubCategoryList"] = NewAddSubCategoryListHandler()
	//handlerMap["addMoneyManagerList"] = NewAddMoneyManagerListHandler()
	//handlerMap["addScheduleTaskList"] = NewAddScheduleTaskListHandler()
	//handlerMap["addToDoTaskList"] = NewAddToDoTaskListHandler()
	//handlerMap["addMoneyTaskList"] = NewAddMoneyTaskListHandler()
	//handlerMap["removeMainCategoryList"] = NewRemoveMainCategoryListHandler()
	//handlerMap["removeSubCategoryList"] = NewRemoveSubCategoryListHandler()
	//handlerMap["removeScheduleTaskList"] = NewRemoveScheduleTaskListHandler()
	//handlerMap["removeToDoTaskList"] = NewRemoveToDoTaskListHandler()
	//handlerMap["removeMoneyTaskList"] = NewRemoveMoneyTaskListHandler()
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

type BaseHandler interface {
	ApiCall(c *gin.Context)
}

type ApiHandler struct {
	authHandler AuthHandler
}

func (ah *ApiHandler) apiCall(c *gin.Context, process func(accountId int, b []byte) ([]byte, error)) {
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		lflog.Logging(lflog.LogLevelInfo, err.Error())
		c.String(http.StatusOK, string(ResponseToByteArray(response.CreateFailResponse(301, "body error"))))
		return
	}

	//받은 데이터 출력
	lflog.Logging(lflog.LogLevelInfo, string(body))

	//기본 세팅(현재 시간, 유저 정보 등등...)
	accountId := c.GetInt("account_id")

	//로직 실행
	res, err := process(accountId, body)

	if err != nil {
		lflog.Logging(lflog.LogLevelInfo, err.Error())
	}

	resultLog := "{\"input\":" + string(body) + ", \"output\":" + string(res) + "}"

	lflog.Logging(lflog.LogLevelInfo, resultLog)

	c.String(http.StatusOK, string(res))
}

type AuthHandler struct {
}

func (ah *AuthHandler) apiCall(c *gin.Context, process func(b []byte) ([]byte, error)) {

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

//
//type apiHandler interface {
//	process(reqBody []byte) ([]byte, error)
//}

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
