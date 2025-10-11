package models

import (
	"time"

	"gorm.io/gorm"
)

type Aspiration struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Title         string         `json:"title" gorm:"not null"`
	Description   string         `json:"description" gorm:"type:text;not null"`
	Category      string         `json:"category" gorm:"type:text;not null"`
	PriorityLevel string         `json:"priority_level" gorm:"type:text;not null"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Aspiration) TableName() string {
	return "aspirations"
}
