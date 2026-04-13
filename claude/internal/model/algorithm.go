package model

import (
	"time"

	"gorm.io/gorm"
)

// Algorithm represents the algorithms table (decision-tree nodes).
type Algorithm struct {
	ID                     uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Text                   string         `gorm:"type:text;not null" json:"text"`
	ActionType             NullString     `gorm:"type:enum('bool','input','date');null" json:"actionType"`
	Positive               NullString     `gorm:"type:varchar(255);null" json:"positive"`
	Negative               NullString     `gorm:"type:varchar(255);null" json:"negative"`
	OnPositive             *uint          `gorm:"null" json:"onPositive"`
	OnNegative             *uint          `gorm:"null" json:"onNegative"`
	NextMove               *uint          `gorm:"null" json:"nextMove"`
	TempPlan               NullString     `gorm:"type:varchar(255);null" json:"tempPlan"`
	TempPlanDirP           NullString     `gorm:"type:varchar(255);null" json:"tempPlanDirP"`
	TempPlanDirN           NullString     `gorm:"type:varchar(255);null" json:"tempPlanDirN"`
	ConditionalFactor      NullString     `gorm:"type:varchar(255);null" json:"conditionalFactor"`
	ConditionalOperator    NullString     `gorm:"type:varchar(10);null" json:"conditionalOperator"`
	ConditionalValue       NullString     `gorm:"type:varchar(255);null" json:"conditionalValue"`
	StateValue             NullString     `gorm:"type:varchar(255);null" json:"stateValue"`
	Label                  NullString     `gorm:"type:varchar(255);null" json:"label"`
	ProgestogenPossible    NullString     `gorm:"type:enum('true','false');null" json:"progestogenPossible"`
	ProgestogenPossibleDir NullString     `gorm:"type:enum('positive','negative');null" json:"progestogenPossibleDir"`
	Delay                  int            `gorm:"type:int;not null;default:0" json:"delay"`
	Series                 *int           `gorm:"null" json:"series"`
	Active                 string         `gorm:"type:enum('Y','N');default:'N'" json:"active"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
}

// TableName overrides the table name.
func (Algorithm) TableName() string { return "algorithms" }
