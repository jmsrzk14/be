package models

import (
	"time"

	"gorm.io/gorm"
)

// Activity represents a general event that can be organized by any entity.
type Activity struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	DepartmentID  *uint          `json:"department_id,omitempty" gorm:"index"`
	AssociationID *uint          `json:"association_id,omitempty" gorm:"index"`
	BEMID         *uint          `json:"bem_id,omitempty" gorm:"index"`
	Title         string         `json:"title" gorm:"type:varchar(255);not null"`
	Description   string         `json:"description" gorm:"type:text"`
	StartDate     time.Time      `json:"start_date"`
	EndDate       time.Time      `json:"end_date"`
	Location      string         `json:"location" gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Activity) TableName() string {
	return "activities"
}

// Proposal represents a proposal document for an activity.
type Proposal struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	ActivityID uint           `json:"activity_id" gorm:"not null"`
	Activity   *Activity      `json:"activity,omitempty" gorm:"foreignKey:ActivityID"`
	FilePath   string         `json:"file_path" gorm:"type:varchar(255);not null"`
	Status     string         `json:"status" gorm:"type:varchar(20);default:'pending';comment:approved, pending, rejected"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Proposal) TableName() string {
	return "proposals"
}

// Report represents a final accountability report for an activity.
type Report struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	ActivityID uint           `json:"activity_id" gorm:"not null"`
	Activity   *Activity      `json:"activity,omitempty" gorm:"foreignKey:ActivityID"`
	FilePath   string         `json:"file_path" gorm:"type:varchar(255);not null"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Report) TableName() string {
	return "reports"
}

// Finance tracks the finances for an activity.
type Finance struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	ActivityID uint           `json:"activity_id" gorm:"not null"`
	Activity   *Activity      `json:"activity,omitempty" gorm:"foreignKey:ActivityID"`
	Income     float64        `json:"income" gorm:"type:decimal(15,2);default:0"`
	Expense    float64        `json:"expense" gorm:"type:decimal(15,2);default:0"`
	Balance    float64        `json:"balance" gorm:"type:decimal(15,2);default:0"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Finance) TableName() string {
	return "finances"
}

