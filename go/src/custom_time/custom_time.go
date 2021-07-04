package custom_time

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// CustomTime provides an example of how to declare a new time Type with a custom formatter.
// Note that time.Time methods are not available, if needed you can add and cast like the String method does
// Otherwise, only use in the json struct at marshal/unmarshal time.
type CustomTime time.Time

const ctLayout = "2006-01-02 15:04:05"

// UnmarshalJSON Parses the json string in the custom format
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(ctLayout, s)
	*ct = CustomTime(nt)
	return
}

// MarshalJSON writes a quoted string in the custom format
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

// String returns the time in the custom format
func (ct *CustomTime) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(ctLayout))
}

func (ct *CustomTime) Convert(t time.Time) (err error) {
	*ct = CustomTime(t)
	return
}

func Now() CustomTime {
	ct := CustomTime(time.Now())
	return ct
}

func NewCustomTime(hour, min, sec int) CustomTime {
    t := time.Date(0, time.January, 1, hour, min, sec, 0, time.UTC)
    return CustomTime(t)
}

func (t *CustomTime) Scan(value interface{}) error {
    switch v := value.(type) {
    case []byte:
        return t.UnmarshalText(string(v))
    case string:
        return t.UnmarshalText(v)
    case time.Time:
        *t = CustomTime(v)
    case nil:
        *t = CustomTime{}
    default:
        return fmt.Errorf("cannot sql.Scan() MyTime from: %#v", v)
    }
    return nil
}

func (t CustomTime) Value() (driver.Value, error) {
    return driver.Value(time.Time(t).Format(ctLayout)), nil
}

func (t *CustomTime) UnmarshalText(value string) error {
    dd, err := time.Parse(ctLayout, value)
    if err != nil {
        return err
    }
    *t = CustomTime(dd)
    return nil
}

func (CustomTime) GormDataType() string {
    return "TIME"
}