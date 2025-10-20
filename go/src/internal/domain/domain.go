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
	ID            define.SnowflakeID `json:"id"`
	LocalID       string             `json:"localId"`
	ParentID      define.SnowflakeID `json:"parentId"`
	LocalParentID string             `json:"localParentId"`
	UserID        define.SnowflakeID `json:"userId"`
	Name          string             `json:"name"`
	Order         int                `json:"order"`
	Color         int                `json:"color"`
	Type          define.FileType    `json:"type"`
	CreatedAt     core.CustomTime    `json:"createdAt"`
	UpdatedAt     core.CustomTime    `json:"updatedAt"`
}

// Folder represents a node that can contain child folders and plans.
type Folder struct {
	SystemFile

	Children []ChildSystemFile `json:"children"`
}

// Plan captures configuration and child entities.
type Plan struct {
	SystemFile

	Description string           `json:"description"`
	StartDate   *core.CustomTime `json:"startDate"`
	FinishDate  *core.CustomTime `json:"finishDate"`
	DateType    define.DateType  `json:"dateType"`

	MainRecordFieldID      define.SnowflakeID `json:"mainRecordFieldId"`
	LocalMainRecordFieldID string             `json:"localMainRecordFieldId"`

	Goals        []*Goal        `json:"goals"`
	RecordFields []*RecordField `json:"recordFields"`
	RepeatRules  []*RepeatRule  `json:"repeatRules"`
	Records      []*Record      `json:"records"`
	Statistics   []*Statistics  `json:"statistics"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// Goal mirrors GoalEntity.
type Goal struct {
	ID                 define.SnowflakeID         `json:"id"`
	LocalID            string                     `json:"localId"`
	PlanID             define.SnowflakeID         `json:"planId"`
	LocalPlanID        string                     `json:"localPlanId"`
	CycleType          define.GoalRepeatCycleType `json:"cycleType"`
	Interval           int                        `json:"interval"`
	RecordFieldID      define.SnowflakeID         `json:"recordFieldId"`
	LocalRecordFieldID string                     `json:"localRecordFieldId"`
	Count              int                        `json:"count"`
	StartDate          *core.CustomTime           `json:"startDate"`
	FinishDate         *core.CustomTime           `json:"finishDate"`
	Order              int                        `json:"order"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// RecordField corresponds to RecordFieldEntity; both regular and repeat fields live here.
type RecordField struct {
	ID            define.SnowflakeID `json:"id"`
	LocalID       string             `json:"localId"`
	PlanID        define.SnowflakeID `json:"planId"`
	LocalPlanID   string             `json:"localPlanId"`
	Name          string             `json:"name"`
	Type          string             `json:"type"`
	IsRepeatField bool               `json:"isRepeatField"`
	Unit          string             `json:"unit"`
	Order         int                `json:"order"`

	Options []*OptionItem `json:"options"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// OptionItem is a tree under a RecordField.
type OptionItem struct {
	ID                 define.SnowflakeID `json:"id"`
	LocalID            string             `json:"localId"`
	PlanID             define.SnowflakeID `json:"planId"`
	LocalPlanID        string             `json:"localPlanId"`
	ParentID           define.SnowflakeID `json:"parentId"`
	LocalParentID      string             `json:"localParentId"`
	RecordFieldID      define.SnowflakeID `json:"recordFieldId"`
	LocalRecordFieldID string             `json:"localRecordFieldId"`
	Name               string             `json:"name"`
	Order              int                `json:"order"`

	Children []*OptionItem `json:"children"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// RepeatRule represents recurrence constraints.
type RepeatRule struct {
	ID             define.SnowflakeID `json:"id"`
	LocalID        string             `json:"localId"`
	PlanID         define.SnowflakeID `json:"planId"`
	LocalPlanID    string             `json:"localPlanId"`
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
	Order          int                `json:"order"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// Record stores user-entered values.
type Record struct {
	ID          define.SnowflakeID `json:"id"`
	LocalID     string             `json:"localId"`
	PlanID      define.SnowflakeID `json:"planId"`
	LocalPlanID string             `json:"localPlanId"`
	Values      map[string]any     `json:"values"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}

// Statistics configuration for a plan.
type Statistics struct {
	ID            define.SnowflakeID         `json:"id"`
	LocalID       string                     `json:"localId"`
	PlanID        define.SnowflakeID         `json:"planId"`
	LocalPlanID   string                     `json:"localPlanId"`
	Name          string                     `json:"name"`
	ChartType     define.StatisticsChartType `json:"chartType"`
	XAxisType     string                     `json:"xAxisType"`
	XAxisPath     []string                   `json:"xAxisPath"`
	YAxisField    define.SnowflakeID         `json:"yAxisField"`
	XAxisIsRepeat bool                       `json:"xAxisIsRepeat"`
	YAxisIsRepeat bool                       `json:"yAxisIsRepeat"`
	Order         int                        `json:"order"`

	CreatedAt core.CustomTime `json:"createdAt"`
	UpdatedAt core.CustomTime `json:"updatedAt"`
}
