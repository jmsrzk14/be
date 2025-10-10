package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:255;not null"`
	Category  int            `json:"category"`
	Image     string         `json:"image" gorm:"type:varchar(255)"`
	Amount    int            `json:"amount" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;uniqueIndex:idx_courses_code_deleted_at" json:"deleted_at,omitempty"`
}

func (Item) TableName() string {
	return "item"
}
