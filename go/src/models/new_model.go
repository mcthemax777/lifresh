package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"lifresh/define"
	"time"

	"gorm.io/gorm"
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
	ID          define.SnowflakeID `gorm:"type:bigint;primaryKey" json:"id"`
	Name        *string            `gorm:"type:varchar(100)" json:"name"`
	Email       *string            `gorm:"type:varchar(191);uniqueIndex" json:"email"`
	PhotoURL    *string            `gorm:"type:varchar(255)" json:"photoUrl"`
	Social      SocialType         `gorm:"not null" json:"socialType"`
	ProviderUID *string            `gorm:"type:varchar(191);uniqueIndex" json:"providerUid"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`

	//User *User `gorm:"-" json:"user"`
}

// User is the domain profile. Root points to the root Folder.
type User struct {
	ID         define.SnowflakeID `gorm:"type:bigint;primaryKey" json:"id"`
	AccountID  define.SnowflakeID `gorm:"type:bigint;index;not null" json:"accountId"`
	Nickname   string             `gorm:"type:varchar(60);not null" json:"nickname"`
	Bio        string             `gorm:"type:varchar(255)" json:"bio"`
	ProfileURL *string            `gorm:"type:varchar(255)" json:"profileUrl"`
	CreatedAt  time.Time          `json:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt"`

	Root *Folder `gorm:"-" json:"root"`
}

// Folder represents a node that can contain child folders and plans.
type Folder struct {
	ID       define.SnowflakeID  `gorm:"type:bigint;primaryKey" json:"id"`
	ParentID *define.SnowflakeID `gorm:"type:bigint;index" json:"parentId"`
	UserID   define.SnowflakeID  `gorm:"type:bigint;index;not null" json:"userId"`
	Order    int                 `gorm:"not null;default:0" json:"order"`
	Color    int                 `gorm:"not null" json:"color"`
	Name     string              `gorm:"type:varchar(120);not null" json:"name"`

	ChildrenFolders []Folder `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"childrenFolders"`
	Plans           []Plan   `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"plans"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Plan captures configuration and child entities.
type Plan struct {
	ID          define.SnowflakeID  `gorm:"type:bigint;primaryKey" json:"id"`
	ParentID    *define.SnowflakeID `gorm:"type:bigint;index" json:"parentId"` // Folder.ID
	UserID      define.SnowflakeID  `gorm:"type:bigint;index;not null" json:"userId"`
	Order       int                 `gorm:"not null;default:0" json:"order"`
	Color       int                 `gorm:"not null" json:"color"`
	Name        string              `gorm:"type:varchar(120);not null" json:"name"`
	Description string              `gorm:"type:varchar(255)" json:"description"`
	StartDate   *time.Time          `json:"startDate"`
	FinishDate  *time.Time          `json:"finishDate"`
	DateType    DateType            `gorm:"not null" json:"dateType"`

	MainRecordFieldID *define.SnowflakeID `gorm:"type:bigint;index" json:"mainRecordFieldId"`

	Goals        []Goal        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"goals"`
	RecordFields []RecordField `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"recordFields"`
	RepeatRules  []RepeatRule  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"repeatRules"`
	Records      []Record      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"records"`
	Statistics   []Statistics  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"statistics"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Goal mirrors GoalEntity.
type Goal struct {
	ID         define.SnowflakeID  `gorm:"type:bigint;primaryKey" json:"id"`
	PlanID     define.SnowflakeID  `gorm:"type:bigint;index;not null" json:"planId"`
	CycleType  GoalRepeatCycleType `gorm:"not null" json:"cycleType"`
	Interval   int                 `gorm:"not null" json:"interval"`
	RecordID   *define.SnowflakeID `gorm:"type:bigint" json:"recordId"`
	Count      int                 `gorm:"not null" json:"count"`
	StartDate  *time.Time          `json:"startDate"`
	FinishDate *time.Time          `json:"finishDate"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RecordField corresponds to RecordFieldEntity; both regular and repeat fields live here.
type RecordField struct {
	ID            define.SnowflakeID `gorm:"type:bigint;primaryKey" json:"id"`
	PlanID        define.SnowflakeID `gorm:"type:bigint;index;not null" json:"planId"`
	Name          string             `gorm:"type:varchar(120);not null" json:"name"`
	Type          string             `gorm:"type:varchar(60);not null" json:"type"`
	IsRepeatField bool               `gorm:"not null;default:false" json:"isRepeatField"`
	Unit          *string            `gorm:"type:varchar(30)" json:"unit"`

	Options []OptionItem `gorm:"foreignKey:RecordFieldID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"options"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// OptionItem is a tree under a RecordField.
type OptionItem struct {
	ID            define.SnowflakeID  `gorm:"type:bigint;primaryKey" json:"id"`
	RecordFieldID define.SnowflakeID  `gorm:"type:bigint;index;not null" json:"recordFieldId"`
	ParentID      *define.SnowflakeID `gorm:"type:bigint;index" json:"parentId"`
	Name          string              `gorm:"type:varchar(120);not null" json:"name"`

	Children []OptionItem `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"children"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RepeatRule represents recurrence constraints.
type RepeatRule struct {
	ID             define.SnowflakeID `gorm:"type:bigint;primaryKey" json:"id"`
	PlanID         define.SnowflakeID `gorm:"type:bigint;index;not null" json:"planId"`
	Unit           RepeatUnit         `gorm:"not null" json:"unit"`
	Interval       int                `gorm:"not null" json:"interval"`
	PerRepeatCount *int               `json:"perRepeatCount"`
	Months         IntArray           `gorm:"type:json" json:"months"`
	Days           IntArray           `gorm:"type:json" json:"days"`
	Weeks          IntArray           `gorm:"type:json" json:"weeks"`
	Weekdays       IntArray           `gorm:"type:json" json:"weekdays"`
	Times          StringArray        `gorm:"type:json" json:"times"` // e.g. ["09:00","18:30"]
	StartDate      *time.Time         `json:"startDate"`
	EndDate        *time.Time         `json:"endDate"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Record stores user-entered values (including start/finish date fields) as JSON.
// Expected shape (app contract):
//
//	{
//	  "fieldValues": { fieldId: any, startDateTime: string, finishDateTime: string },
//	  "repeatValues": [ { fieldId: any }, ... ]
//	}
type Record struct {
	ID     define.SnowflakeID `gorm:"type:bigint;primaryKey" json:"id"`
	PlanID define.SnowflakeID `gorm:"type:bigint;index;not null" json:"planId"`
	Values JSONMap            `gorm:"type:json;not null" json:"values"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Statistics configuration for a plan.
type Statistics struct {
	ID            define.SnowflakeID  `gorm:"type:bigint;primaryKey" json:"id"`
	PlanID        define.SnowflakeID  `gorm:"type:bigint;index;not null" json:"planId"`
	Name          string              `gorm:"type:varchar(120);not null" json:"name"`
	ChartType     StatisticsChartType `gorm:"not null" json:"chartType"`
	XAxisType     string              `gorm:"type:varchar(60);not null" json:"xAxisType"`
	XAxisPath     StringArray         `gorm:"type:json" json:"xAxisPath"`
	YAxisField    define.SnowflakeID  `gorm:"type:bigint;not null" json:"yAxisField"`
	XAxisIsRepeat bool                `gorm:"not null;default:false" json:"xAxisIsRepeat"`
	YAxisIsRepeat bool                `gorm:"not null;default:false" json:"yAxisIsRepeat"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
