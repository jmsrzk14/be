package models

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	EntityType  string         `json:"entity_type"`  // organization, announcement, news
	EntityID    uint           `json:"entity_id"`
	Type        string         `json:"type"`         // create, update, delete
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

