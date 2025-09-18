package model

import (
	"database/sql/driver"
	"gorm.io/gorm"
	"lifresh/define"
	"time"
)

import (
	"encoding/json"
	"errors"
)

// Optional: plug your own generator at runtime (service layer)
// var DefaultIDGen interface{ NextID() SnowflakeID }

// =====================================================
// Helper JSON types for MySQL JSON columns
// =====================================================

type JSONMap map[string]any

type IntArray []int

type StringArray []string

func (m JSONMap) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
func (m *JSONMap) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, m)
	case string:
		return json.Unmarshal([]byte(v), m)
	case nil:
		*m = JSONMap{}
		return nil
	}
	return errors.New("unsupported type for JSONMap")
}

func (a IntArray) Value() (driver.Value, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
func (a *IntArray) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	case nil:
		*a = IntArray{}
		return nil
	}
	return errors.New("unsupported type for IntArray")
}

func (a StringArray) Value() (driver.Value, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
func (a *StringArray) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	case nil:
		*a = StringArray{}
		return nil
	}
	return errors.New("unsupported type for StringArray")
}

// =====================================================
// Enums (stored as SMALLINT/INT in DB)
// =====================================================

type SocialType int32

const (
	SocialGuest SocialType = iota
	SocialGoogle
	SocialApple
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

// =====================================================
// Core models (IDs are BIGINT Snowflake)
// =====================================================

// Account: authentication account; separate from domain User profile.
// NOTE: ProviderUID keeps the external provider's user id (string).
type Account struct {
	ID          define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	Name        string             `gorm:"type:varchar(100)"`
	Email       string             `gorm:"type:varchar(191);uniqueIndex"`
	PhotoURL    string             `gorm:"type:varchar(255)"`
	SocialType  int8               `gorm:"not null" `
	ProviderUID string             `gorm:"type:varchar(191);uniqueIndex"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// User is the domain profile. Root points to the root Folder.
type User struct {
	ID         define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	AccountID  define.SnowflakeID `gorm:"type:bigint;index;not null"`
	Nickname   string             `gorm:"type:varchar(60);not null"`
	Bio        string             `gorm:"type:varchar(255)"`
	ProfileURL string             `gorm:"type:varchar(255)"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Root *Folder `gorm:"-" json:"root"`
}

// Folder represents a node that can contain child folders and plans.
type Folder struct {
	ID        define.SnowflakeID  `gorm:"type:bigint;primaryKey"`
	ParentID  *define.SnowflakeID `gorm:"type:bigint;index"`
	UserID    define.SnowflakeID  `gorm:"type:bigint;index;not null"`
	Order     int                 `gorm:"not null;default:0"`
	Color     int                 `gorm:"not null"`
	Name      string              `gorm:"type:varchar(120);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Plan captures configuration and child entities.
type Plan struct {
	ID          define.SnowflakeID  `gorm:"type:bigint;primaryKey"`
	ParentID    *define.SnowflakeID `gorm:"type:bigint;index"`
	UserID      define.SnowflakeID  `gorm:"type:bigint;index;not null"`
	Order       int                 `gorm:"not null;default:0"`
	Color       int                 `gorm:"not null"`
	Name        string              `gorm:"type:varchar(120);not null"`
	Description string              `gorm:"type:varchar(255)"`
	StartDate   time.Time
	FinishDate  time.Time
	DateType    DateType `gorm:"not null"`

	MainRecordFieldID *define.SnowflakeID `gorm:"type:bigint;index"`

	Goals        []Goal        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RecordFields []RecordField `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RepeatRules  []RepeatRule  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Records      []Record      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Statistics   []Statistics  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Goal mirrors GoalEntity.
type Goal struct {
	ID         define.SnowflakeID  `gorm:"type:bigint;primaryKey"`
	PlanID     define.SnowflakeID  `gorm:"type:bigint;index;not null"`
	CycleType  GoalRepeatCycleType `gorm:"not null"`
	Interval   int                 `gorm:"not null"`
	RecordID   *define.SnowflakeID `gorm:"type:bigint"`
	Count      int                 `gorm:"not null"`
	StartDate  time.Time
	FinishDate time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// RecordField corresponds to RecordFieldEntity; both regular and repeat fields live here.
type RecordField struct {
	ID            define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	PlanID        define.SnowflakeID `gorm:"type:bigint;index;not null"`
	Name          string             `gorm:"type:varchar(120);not null"`
	Type          string             `gorm:"type:varchar(60);not null"`
	IsRepeatField bool               `gorm:"not null;default:false"`
	Unit          string             `gorm:"type:varchar(30)"`

	Options []OptionItem `gorm:"foreignKey:RecordFieldID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// OptionItem is a tree under a RecordField.
type OptionItem struct {
	ID            define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	RecordFieldID define.SnowflakeID `gorm:"type:bigint;index;not null"`
	ParentID      define.SnowflakeID `gorm:"type:bigint;index"`
	Name          string             `gorm:"type:varchar(120);not null"`

	Children []OptionItem `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// RepeatRule represents recurrence constraints.
type RepeatRule struct {
	ID             define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	PlanID         define.SnowflakeID `gorm:"type:bigint;index;not null"`
	Unit           RepeatUnit         `gorm:"not null"`
	Interval       int                `gorm:"not null"`
	PerRepeatCount int
	Months         IntArray    `gorm:"type:json"`
	Days           IntArray    `gorm:"type:json"`
	Weeks          IntArray    `gorm:"type:json"`
	Weekdays       IntArray    `gorm:"type:json"`
	Times          StringArray `gorm:"type:json"` // e.g. ["09:00","18:30"]
	StartDate      time.Time
	EndDate        time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Record stores user-entered values (including start/finish date fields) as JSON.
// Expected shape (app contract):
//
//	{
//	  "fieldValues": { fieldId: any, startDateTime: string, finishDateTime: string },
//	  "repeatValues": [ { fieldId: any }, ... ]
//	}
type Record struct {
	ID     define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	PlanID define.SnowflakeID `gorm:"type:bigint;index;not null"`
	Values JSONMap            `gorm:"type:json;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Statistics configuration for a plan.
type Statistics struct {
	ID            define.SnowflakeID  `gorm:"type:bigint;primaryKey"`
	PlanID        define.SnowflakeID  `gorm:"type:bigint;index;not null"`
	Name          string              `gorm:"type:varchar(120);not null"`
	ChartType     StatisticsChartType `gorm:"not null"`
	XAxisType     string              `gorm:"type:varchar(60);not null"`
	XAxisPath     StringArray         `gorm:"type:json"`
	YAxisField    define.SnowflakeID  `gorm:"type:bigint;not null"`
	XAxisIsRepeat bool                `gorm:"not null;default:false"`
	YAxisIsRepeat bool                `gorm:"not null;default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// =====================================================
// AutoMigrate helper
// =====================================================
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Account{},
		&User{},
		&Folder{},
		&Plan{},
		&Goal{},
		&RecordField{},
		&OptionItem{},
		&RepeatRule{},
		&Record{},
		&Statistics{},
	)
}
