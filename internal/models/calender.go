// models/event.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Calender struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Title          string         `gorm:"type:text;not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description,omitempty"`
	Location       string         `gorm:"type:text" json:"location,omitempty"`
	StartTime      time.Time      `gorm:"not null" json:"start"`
	EndTime        time.Time      `gorm:"not null" json:"end"`
	OrganizationID *int           `form:"organization_id" json:"organization_id"` // nullable
	Organization   *Organization  `json:"organization" gorm:"foreignKey:ID;references:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
