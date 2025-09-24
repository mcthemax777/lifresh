package define

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
)

const (
	OsTypeLinux   = 1
	OsTypeWindows = 2
	OsTypeMac     = 3
)

var OsType int

func init() {
	os := runtime.GOOS
	switch os {
	case "windows":
		OsType = OsTypeWindows
	case "linux":
		OsType = OsTypeLinux
	case "macos":
	case "darwin":
		OsType = OsTypeMac
	default:
		panic(0)
	}

	fmt.Printf("os type - %d\n", OsType)
}

type SocialType int8

const (
	SocialTypeGuest SocialType = iota
	SocialTypeGoogle
	SocialTypeApple
)

type GoalRepeatCycleType int32

const (
	CycleDay GoalRepeatCycleType = iota
	CycleWeek
	CycleMonth
	CycleYear
	CycleWhole
)

type DateType int32

const (
	DateTypeDateTime DateType = iota
	DateTypePeriod
)

type RepeatUnit int32

const (
	RepeatDay RepeatUnit = iota
	RepeatWeek
	RepeatMonth
	RepeatYear
	RepeatCountOnly
)

type StatisticsChartType int32

const (
	ChartPie StatisticsChartType = iota
	ChartBar
)

type FileType int8

const (
	FileTypeFolder FileType = iota
	FileTypePlan
)

const (
	RouterAuth = "/v1/auth"
	RouterApi  = "/v1/api"
)
const (
	ApiLogin    = "/login"
	ApiJwtLogin = "/jwt-login"
	ApiRefresh  = "/refresh"

	ApiPing         = "/ping"
	ApiGetUser      = "/user"
	ApiCreateUser   = "/create-user"
	ApiCreateFolder = "/create-folder"
	ApiUpdateFolder = "/update-folder"
	ApiDeleteFolder = "/delete-folder"
)

// =====================================================
// SnowflakeID: stored as BIGINT in DB, serialized as string in JSON
// =====================================================
type SnowflakeID int64

// Marshal as JSON string (to avoid JS number precision issues)
func (id SnowflakeID) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%d", id))
}

// Accept both string and number in requests
func (id *SnowflakeID) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		*id = SnowflakeID(0)
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			//일단 공백오면 이쪽으로 오니 0으로 세팅
			*id = SnowflakeID(0)
			return nil
		}
		*id = SnowflakeID(v)
		return nil
	}
	var i int64
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	*id = SnowflakeID(i)
	return nil
}

// DB driver hooks
func (id SnowflakeID) Value() (driver.Value, error) { return int64(id), nil }
func (id *SnowflakeID) Scan(value any) error {
	switch v := value.(type) {
	case int64:
		*id = SnowflakeID(v)
		return nil
	case []byte:
		if len(v) == 0 {
			*id = SnowflakeID(0)
			return nil
		}
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
		*id = SnowflakeID(i)
		return nil
	case string:
		if len(v) == 0 {
			*id = SnowflakeID(0)
			return nil
		}
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		*id = SnowflakeID(i)
		return nil
	case nil:
		*id = 0
		return nil
	}
	return fmt.Errorf("unsupported Scan type for SnowflakeID: %T", value)
}

// 공통적으로 ID가 비어있을 경우 새로 생성
func IfZero[T comparable](val T, fallback T) T {
	var zero T
	if val == zero {
		return fallback
	}
	return val
}
