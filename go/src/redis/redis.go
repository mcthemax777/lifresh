package redis

import (
	"encoding/json"
	"fmt"
	"lifresh/define"
	"log"
	"time"

	"github.com/go-redis/redis"
)

var RedisHandlerSG RedisHandlerImpl

var redisClient *redis.Client

type RedisInfo struct {
	user string
	pwd  string
	url  string
}

func init() {

	var localRedisInfo = RedisInfo{"root", "1234", "host.docker.internal:6379"}

	if define.OsType == define.OsTypeWindows || define.OsType == define.OsTypeMac {
		localRedisInfo = RedisInfo{"root", "1234", "localhost:6379"}
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

func GetRedisClient() *redis.Client { return redisClient }

func SubscribeStream(rdb *redis.Client, stream string) {
	group := "lifresh_workers"

	// 그룹 생성 (없으면)
	rdb.XGroupCreateMkStream(stream, group, "$")

	for {
		msgs, err := rdb.XReadGroup(&redis.XReadGroupArgs{
			Group:    group,
			Consumer: "worker-1",
			Streams:  []string{stream, ">"},
			Count:    10,
			Block:    0,
		}).Result()

		if err != nil {
			log.Printf("[StreamWorker] read err: %v", err)
			continue
		}

		for _, m := range msgs {
			for _, v := range m.Messages {
				fmt.Println("Received:", v.Values)
				// TODO: handle payload
			}
		}
	}
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
	AccountId  int
	ExpireTime time.Time
}

func (dh RedisHandlerImpl) SetSession(uid string, sid string, accountId int) error {

	expireDuration, _ := time.ParseDuration("1h")

	var sessionInfo SessionInfo
	sessionInfo.Sid = sid
	sessionInfo.AccountId = accountId
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
