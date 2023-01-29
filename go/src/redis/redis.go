package redis

import (
	"context"
	"encoding/json"
	"lifresh/define"
	"time"

	"github.com/go-redis/redis"
)

var RedisHandlerSG RedisHandlerImpl

var redisClient *redis.Client

type RedisInfo struct {
	user     string
	pwd      string
	url      string
	engine   string
	database string
}

// type camelNamer struct
// {

// }

// func (cm *camelNamer) ColumnName(table, column string) string {
// 	return column
// }

// func (cm *camelNamer) TableName(table string) string {
// 	return table
// }

// func (cm *camelNamer) JoinTableName(table string) string {
// 	return table
// }

// func (cm *camelNamer) RelationshipFKName(ss schema.Relationship) string {
// 	return ""
// }

// func (cm *camelNamer) CheckerName(table, column string) string {
// 	return column
// }

// func (cm *camelNamer) IndexName(table, column string) string {
// 	return column
// }

var ctx = context.Background()
var SESSION_KEY = "session"

func init() {

	var localRedisInfo = RedisInfo{"root", "1234", "host.docker.internal:6379", "mysql", "Lifresh"}

	if define.OsType == define.OsTypeWindows {
		localRedisInfo = RedisInfo{"root", "1234", "localhost:6379", "mysql", "Lifresh"}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     localRedisInfo.url, // 접근 url 및 port
		Password: "",                 // password ""값은 없다는 뜻
		DB:       0,                  // 기본 DB 사용
	})

	_, err := client.Ping().Result()

	if err != nil {
		return
	}

	redisClient = client
}

// type RedisHandler interface {
// 	InsertAccount(userId string, password string) error
// 	Login(userId string, password string) error
// }

type RedisHandlerImpl struct {
	//dbConn *gorm.DB
}

type SessionInfo struct {
	Sid        string
	AccountNo  int
	ExpireTime time.Time
}

func (dh RedisHandlerImpl) SetSession(uid string, sid string, accountNo int) error {

	expireDuration, _ := time.ParseDuration("1h")

	var sessionInfo SessionInfo
	sessionInfo.Sid = sid
	sessionInfo.AccountNo = accountNo
	sessionInfo.ExpireTime = time.Now().Add(expireDuration)

	value, err := json.Marshal(sessionInfo)

	if err != nil {
		return err
	}

	err = redisClient.Set(uid, string(value), expireDuration).Err()

	if err != nil {
		return err
	}

	return nil
}

func (dh RedisHandlerImpl) GetSession(uid string) (SessionInfo, error) {

	var sessionInfo SessionInfo

	sessionInfoStr, err := redisClient.Get(uid).Result()

	if err != nil {
		return sessionInfo, err
	}

	err = json.Unmarshal([]byte(sessionInfoStr), &sessionInfo)

	return sessionInfo, err
}

func (dh RedisHandlerImpl) DeleteSession(sid string) {
	redisClient.Del(sid)
}
