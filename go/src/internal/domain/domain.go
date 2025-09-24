package domain

import (
	"lifresh/define"
	"lifresh/internal/core"
)

// =====================================================
// Core models (IDs are BIGINT Snowflake)
// =====================================================

type ChildSystemFile interface{}

// Account: authentication account; separate from domain User profile.
type Account struct {
	ID          define.SnowflakeID `json:"id"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	PhotoURL    string             `json:"photoUrl"`
	SocialType  define.SocialType  `json:"socialType"`
	ProviderUID string             `json:"providerUid"`
	CreatedAt   core.CustomTime    `json:"createdAt"`
	UpdatedAt   core.CustomTime    `json:"updatedAt"`

	User *User `json:"user"`
}

// User is the domain profile. Root points to the root Folder.
type User struct {
	ID         define.SnowflakeID `json:"id"`
	AccountID  define.SnowflakeID `json:"accountId"`
	Nickname   string             `json:"nickname"`
	Bio        string             `json:"bio"`
	ProfileURL string             `json:"profileUrl"`
	CreatedAt  core.CustomTime    `json:"createdAt"`
	UpdatedAt  core.CustomTime    `json:"updatedAt"`

	Root *Folder `json:"root"`
}

type SystemFile struct {
	ID        define.SnowflakeID `json:"id"`
	Name      string             `json:"name"`
	ParentID  define.SnowflakeID `json:"parentId"`
	UserID    define.SnowflakeID `json:"userId"`
	Order     int                `json:"order"`
	Color     int                `json:"color"`
	Type      define.FileType    `json:"type"`
	CreatedAt core.CustomTime    `json:"createdAt"`
	UpdatedAt core.CustomTime    `json:"updatedAt"`
}

// Folder represents a node that can contain child folders and plans.
type Folder struct {
	SystemFile

	Children []ChildSystemFile `json:"children"`
}

// Plan captures configuration and child entities.
type Plan struct {
	SystemFile

	Description string          `json:"description"`
	StartDate   core.CustomTime `json:"startDate"`
	FinishDate  core.CustomTime `json:"finishDate"`
	DateType    define.DateType `json:"dateType"`

	MainRecordFieldID define.SnowflakeID `json:"mainRecordFieldId"`

	Goals        []Goal        `json:"goals"`
	RecordFields []RecordField `json:"recordFields"`
	RepeatRules  []RepeatRule  `json:"repeatRules"`
	Records      []Record      `json:"records"`
	Statistics   []Statistics  `json:"statistics"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// Goal mirrors GoalEntity.
type Goal struct {
	ID         define.SnowflakeID         `json:"id"`
	PlanID     define.SnowflakeID         `json:"planId"`
	CycleType  define.GoalRepeatCycleType `json:"cycleType"`
	Interval   int                        `json:"interval"`
	RecordID   define.SnowflakeID         `json:"recordId"`
	Count      int                        `json:"count"`
	StartDate  core.CustomTime            `json:"startDate"`
	FinishDate core.CustomTime            `json:"finishDate"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// RecordField corresponds to RecordFieldEntity; both regular and repeat fields live here.
type RecordField struct {
	ID            define.SnowflakeID `json:"id"`
	PlanID        define.SnowflakeID `json:"planId"`
	Name          string             `json:"name"`
	Type          string             `json:"type"`
	IsRepeatField bool               `json:"isRepeatField"`
	Unit          string             `json:"unit"`

	Options []OptionItem `json:"options"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// OptionItem is a tree under a RecordField.
type OptionItem struct {
	ID            define.SnowflakeID `json:"id"`
	RecordFieldID define.SnowflakeID `json:"recordFieldId"`
	ParentID      define.SnowflakeID `json:"parentId"`
	Name          string             `json:"name"`

	Children []OptionItem `json:"children"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// RepeatRule represents recurrence constraints.
type RepeatRule struct {
	ID             define.SnowflakeID `json:"id"`
	PlanID         define.SnowflakeID `json:"planId"`
	Unit           define.RepeatUnit  `json:"unit"`
	Interval       int                `json:"interval"`
	PerRepeatCount int                `json:"perRepeatCount"`
	Months         []int              `json:"months"`
	Days           []int              `json:"days"`
	Weeks          []int              `json:"weeks"`
	Weekdays       []int              `json:"weekdays"`
	Times          []string           `json:"times"`
	StartDate      core.CustomTime    `json:"startDate"`
	EndDate        core.CustomTime    `json:"endDate"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// Record stores user-entered values.
type Record struct {
	ID     define.SnowflakeID `json:"id"`
	PlanID define.SnowflakeID `json:"planId"`
	Values map[string]any     `json:"values"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// Statistics configuration for a plan.
type Statistics struct {
	ID            define.SnowflakeID         `json:"id"`
	PlanID        define.SnowflakeID         `json:"planId"`
	Name          string                     `json:"name"`
	ChartType     define.StatisticsChartType `json:"chartType"`
	XAxisType     string                     `json:"xAxisType"`
	XAxisPath     []string                   `json:"xAxisPath"`
	YAxisField    define.SnowflakeID         `json:"yAxisField"`
	XAxisIsRepeat bool                       `json:"xAxisIsRepeat"`
	YAxisIsRepeat bool                       `json:"yAxisIsRepeat"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}
