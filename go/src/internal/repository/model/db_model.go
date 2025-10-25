package model

import (
	"database/sql/driver"
	"gorm.io/gorm"
	"lifresh/define"
	"lifresh/internal/core"
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
	CreatedAt   core.CustomTime
	UpdatedAt   core.CustomTime
}

// User is the domain profile. Root points to the root Folder.
type User struct {
	ID         define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	AccountID  define.SnowflakeID `gorm:"type:bigint;index;not null"`
	Nickname   string             `gorm:"type:varchar(60);not null"`
	Bio        string             `gorm:"type:varchar(255)"`
	ProfileURL string             `gorm:"type:varchar(255)"`
	CreatedAt  core.CustomTime
	UpdatedAt  core.CustomTime
}

// Folder represents a node that can contain child folders and plans.
type Folder struct {
	ID            define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	LocalID       string             `gorm:"type:varchar(255)"`
	ParentID      define.SnowflakeID `gorm:"type:bigint;index"`
	LocalParentID string             `gorm:"type:varchar(255)"`
	UserID        define.SnowflakeID `gorm:"type:bigint;index;not null"`
	Color         int                `gorm:"not null"`
	Name          string             `gorm:"type:varchar(120);not null"`
	Order         int                `gorm:"not null;default:0"`
	EntityVersion int64              `gorm:"not null;default:1"`
	GlobalVersion int64              `gorm:"not null;default:0"`
	CreatedAt     core.CustomTime
	UpdatedAt     core.CustomTime
}

// Plan captures configuration and child entities.
type Plan struct {
	ID            define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	LocalID       string             `gorm:"type:varchar(255)"`
	ParentID      define.SnowflakeID `gorm:"type:bigint;index"`
	LocalParentID string             `gorm:"type:varchar(255)"`
	UserID        define.SnowflakeID `gorm:"type:bigint;index;not null"`
	Color         int                `gorm:"not null"`
	Name          string             `gorm:"type:varchar(120);not null"`
	Description   string             `gorm:"type:varchar(255)"`
	StartDate     *core.CustomTime
	FinishDate    *core.CustomTime
	DateType      DateType `gorm:"not null"`

	MainRecordFieldID define.SnowflakeID `gorm:"type:bigint;index"`

	Order         int   `gorm:"not null;default:0"`
	EntityVersion int64 `gorm:"not null;default:1"`
	GlobalVersion int64 `gorm:"not null;default:0"`
	GroupVersion  int64 `gorm:"not null;default:0"`
	CreatedAt     core.CustomTime
	UpdatedAt     core.CustomTime
}

// Goal mirrors GoalEntity.
type Goal struct {
	ID                 define.SnowflakeID  `gorm:"type:bigint;primaryKey"`
	LocalID            string              `gorm:"type:varchar(255)"`
	PlanID             define.SnowflakeID  `gorm:"type:bigint;index"`
	LocalPlanID        string              `gorm:"type:varchar(255)"`
	CycleType          GoalRepeatCycleType `gorm:"not null"`
	Interval           int                 `gorm:"not null"`
	RecordFieldID      define.SnowflakeID  `gorm:"type:bigint"`
	LocalRecordFieldID string              `gorm:"type:varchar(255)"`
	Count              int                 `gorm:"not null"`
	StartDate          *core.CustomTime
	FinishDate         *core.CustomTime

	Order         int   `gorm:"not null;default:0"`
	EntityVersion int64 `gorm:"not null;default:1"`
	GlobalVersion int64 `gorm:"not null;default:0"`
	CreatedAt     core.CustomTime
	UpdatedAt     core.CustomTime
}

// RecordField corresponds to RecordFieldEntity; both regular and repeat fields live here.
type RecordField struct {
	ID            define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	LocalID       string             `gorm:"type:varchar(255)"`
	PlanID        define.SnowflakeID `gorm:"type:bigint;index"`
	LocalPlanID   string             `gorm:"type:varchar(255)"`
	Name          string             `gorm:"type:varchar(120);not null"`
	Type          string             `gorm:"type:varchar(60);not null"`
	IsRepeatField bool               `gorm:"not null;default:false"`
	Unit          string             `gorm:"type:varchar(30)"`

	Options []OptionItem `gorm:"foreignKey:RecordFieldID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Order         int   `gorm:"not null;default:0"`
	EntityVersion int64 `gorm:"not null;default:1"`
	GlobalVersion int64 `gorm:"not null;default:0"`
	CreatedAt     core.CustomTime
	UpdatedAt     core.CustomTime
}

// OptionItem is a tree under a RecordField.
type OptionItem struct {
	ID                 define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	LocalID            string             `gorm:"type:varchar(255)"`
	PlanID             define.SnowflakeID `gorm:"type:bigint;index"`
	LocalPlanID        string             `gorm:"type:varchar(255)"`
	ParentID           define.SnowflakeID `gorm:"type:bigint;index"`
	LocalParentID      string             `gorm:"type:varchar(255)"`
	RecordFieldID      define.SnowflakeID `gorm:"type:bigint;index;not null"`
	LocalRecordFieldID string             `gorm:"type:varchar(255)"`
	Name               string             `gorm:"type:varchar(120);not null"`

	Children []OptionItem `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Order         int   `gorm:"not null;default:0"`
	EntityVersion int64 `gorm:"not null;default:1"`
	GlobalVersion int64 `gorm:"not null;default:0"`
	CreatedAt     core.CustomTime
	UpdatedAt     core.CustomTime
}

// RepeatRule represents recurrence constraints.
type RepeatRule struct {
	ID             define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	LocalID        string             `gorm:"type:varchar(255)"`
	PlanID         define.SnowflakeID `gorm:"type:bigint;index"`
	LocalPlanID    string             `gorm:"type:varchar(255)"`
	Unit           RepeatUnit         `gorm:"not null"`
	Interval       int                `gorm:"not null"`
	PerRepeatCount int
	Months         IntArray    `gorm:"type:json"`
	Days           IntArray    `gorm:"type:json"`
	Weeks          IntArray    `gorm:"type:json"`
	Weekdays       IntArray    `gorm:"type:json"`
	Times          StringArray `gorm:"type:json"` // e.g. ["09:00","18:30"]
	StartDate      core.CustomTime
	EndDate        core.CustomTime

	Order         int   `gorm:"not null;default:0"`
	EntityVersion int64 `gorm:"not null;default:1"`
	GlobalVersion int64 `gorm:"not null;default:0"`
	CreatedAt     core.CustomTime
	UpdatedAt     core.CustomTime
}

// Record stores user-entered values (including start/finish date fields) as JSON.
// Expected shape (app contract):
//
//	{
//	  "fieldValues": { fieldId: any, startDateTime: string, finishDateTime: string },
//	  "repeatValues": [ { fieldId: any }, ... ]
//	}
type Record struct {
	ID          define.SnowflakeID `gorm:"type:bigint;primaryKey"`
	LocalID     string             `gorm:"type:varchar(255)"`
	PlanID      define.SnowflakeID `gorm:"type:bigint;index"`
	LocalPlanID string             `gorm:"type:varchar(255)"`
	Values      JSONMap            `gorm:"type:json;not null"`

	EntityVersion int64 `gorm:"not null;default:1"`
	GlobalVersion int64 `gorm:"not null;default:0"`
	CreatedAt     core.CustomTime
	UpdatedAt     core.CustomTime
}

// Statistics configuration for a plan.
type Statistics struct {
	ID            define.SnowflakeID  `gorm:"type:bigint;primaryKey"`
	LocalID       string              `gorm:"type:varchar(255)"`
	PlanID        define.SnowflakeID  `gorm:"type:bigint;index"`
	LocalPlanID   string              `gorm:"type:varchar(255)"`
	Name          string              `gorm:"type:varchar(120);not null"`
	ChartType     StatisticsChartType `gorm:"not null"`
	XAxisType     string              `gorm:"type:varchar(60);not null"`
	XAxisPath     StringArray         `gorm:"type:json"`
	YAxisField    define.SnowflakeID  `gorm:"type:bigint;not null"`
	XAxisIsRepeat bool                `gorm:"not null;default:false"`
	YAxisIsRepeat bool                `gorm:"not null;default:false"`

	Order         int   `gorm:"not null;default:0"`
	EntityVersion int64 `gorm:"not null;default:1"`
	GlobalVersion int64 `gorm:"not null;default:0"`
	CreatedAt     core.CustomTime
	UpdatedAt     core.CustomTime
}

// =====================================================
// Sync Meta Tables
// =====================================================

// 그룹 버전 (예: plan_records, folder_tree 등)
type DataVersion struct {
	GroupType string             `gorm:"primaryKey;type:varchar(50)"`
	GroupID   define.SnowflakeID `gorm:"primaryKey;type:bigint"`
	Version   int64              `gorm:"not null"`
	UpdatedAt core.CustomTime    `gorm:"autoUpdateTime"`
}

// 변경 로그 (증분 동기화)
type ChangeLog struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement"`
	GlobalVersion int64  `gorm:"uniqueIndex;not null"`
	EntityType    string
	EntityID      define.SnowflakeID
	GroupType     string
	GroupID       define.SnowflakeID
	Operation     string // create/update/delete
	ActorID       define.SnowflakeID
	Payload       JSONMap         `gorm:"type:json"`
	CreatedAt     core.CustomTime `gorm:"autoCreateTime"`
}

// Outbox for Redis Streams 발행
type Outbox struct {
	ID        uint64          `gorm:"primaryKey;autoIncrement"`
	EventID   string          `gorm:"uniqueIndex;type:varchar(36)"`
	Stream    string          `gorm:"type:varchar(100);not null"`
	Payload   JSONMap         `gorm:"type:json"`
	Processed bool            `gorm:"not null;default:false"`
	CreatedAt core.CustomTime `gorm:"autoCreateTime"`
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
		&DataVersion{},
		&ChangeLog{},
		&Outbox{},
	)
}
