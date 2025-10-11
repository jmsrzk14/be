package models

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Title          string         `json:"title" gorm:"size:255;not null"`
	Content        string         `json:"content" gorm:"type:text;not null"`
	FileURL        string         `json:"file_url,omitempty" gorm:"type:varchar(255);column:file_url"`
	OrganizationID *uint          `json:"organization_id,omitempty"` // nullable
	Organization   *Organization  `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	AuthorID       uint           `json:"author_id" gorm:"not null"`
	Author         *User          `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	StartDate      *time.Time     `json:"start_date,omitempty"`
	EndDate        *time.Time     `json:"end_date,omitempty"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}
