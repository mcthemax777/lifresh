package domain

import (
	"lifresh/define"
	"time"
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
	ID          define.SnowflakeID
	Name        string
	Email       string
	PhotoURL    string
	Social      define.SocialType
	ProviderUID string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	//User *User `gorm:"-" json:"user"`
}

// User is the domain profile. Root points to the root Folder.
type User struct {
	ID         define.SnowflakeID
	AccountID  define.SnowflakeID
	Nickname   string
	Bio        string
	ProfileURL string
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Root *Folder
}

// Folder represents a node that can contain child folders and plans.
type Folder struct {
	ID       define.SnowflakeID
	ParentID *define.SnowflakeID
	UserID   define.SnowflakeID
	Order    int
	Color    int
	Name     string

	ChildrenFolders []Folder
	Plans           []Plan

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Plan captures configuration and child entities.
type Plan struct {
	ID          define.SnowflakeID
	ParentID    *define.SnowflakeID
	UserID      define.SnowflakeID
	Order       int
	Color       int
	Name        string
	Description string
	StartDate   *time.Time
	FinishDate  *time.Time
	DateType    DateType

	MainRecordFieldID *define.SnowflakeID

	Goals        []Goal
	RecordFields []RecordField
	RepeatRules  []RepeatRule
	Records      []Record
	Statistics   []Statistics

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Goal mirrors GoalEntity.
type Goal struct {
	ID         define.SnowflakeID
	PlanID     define.SnowflakeID
	CycleType  GoalRepeatCycleType
	Interval   int
	RecordID   *define.SnowflakeID
	Count      int
	StartDate  *time.Time
	FinishDate *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// RecordField corresponds to RecordFieldEntity; both regular and repeat fields live here.
type RecordField struct {
	ID            define.SnowflakeID
	PlanID        define.SnowflakeID
	Name          string
	Type          string
	IsRepeatField bool
	Unit          *string

	Options []OptionItem

	CreatedAt time.Time
	UpdatedAt time.Time
}

// OptionItem is a tree under a RecordField.
type OptionItem struct {
	ID            define.SnowflakeID
	RecordFieldID define.SnowflakeID
	ParentID      *define.SnowflakeID
	Name          string

	Children []OptionItem

	CreatedAt time.Time
	UpdatedAt time.Time
}

// RepeatRule represents recurrence constraints.
type RepeatRule struct {
	ID             define.SnowflakeID
	PlanID         define.SnowflakeID
	Unit           RepeatUnit
	Interval       int
	PerRepeatCount *int
	Months         []int
	Days           []int
	Weeks          []int
	Weekdays       []int
	Times          []string
	StartDate      *time.Time
	EndDate        *time.Time

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
	ID     define.SnowflakeID
	PlanID define.SnowflakeID
	Values map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Statistics configuration for a plan.
type Statistics struct {
	ID            define.SnowflakeID
	PlanID        define.SnowflakeID
	Name          string
	ChartType     StatisticsChartType
	XAxisType     string
	XAxisPath     []string
	YAxisField    define.SnowflakeID
	XAxisIsRepeat bool
	YAxisIsRepeat bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
