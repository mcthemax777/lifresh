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

const (
	RouterAuth = "/v1/auth"
	RouterApi  = "/v1/api"
)
const (
	ApiLogin   = "/login"
	ApiRefresh = "/refresh"

	ApiPing    = "/ping"
	ApiGetItem = "/item"
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
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
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
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
		*id = SnowflakeID(i)
		return nil
	case string:
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
