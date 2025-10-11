package models

import (
	"time"

	"gorm.io/gorm"
)

// News represents an article or announcement.
type News struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	BEMID         *uint          `json:"bem_id,omitempty" gorm:"index"`
	AssociationID *uint          `json:"association_id,omitempty" gorm:"index"`
	DepartmentID  *uint          `json:"department_id,omitempty" gorm:"index"`
	Title         string         `json:"title" gorm:"type:varchar(255);not null"`
	Content       string         `json:"content" gorm:"type:text;not null"`
	Category      string         `json:"category" gorm:"type:varchar(100)"`
	ImageURL      string         `json:"image_url" gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (News) TableName() string {
	return "news"
}

// Aspiration represents feedback or suggestions from users.
// type Aspiration struct {
// 	ID        uint           `json:"id" gorm:"primaryKey"`
// 	UserID    uint           `json:"user_id" gorm:"not null"`
// 	User      *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
// 	Content   string         `json:"content" gorm:"type:text;not null"`
// 	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
// 	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
// 	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
// }

// func (Aspiration) TableName() string {
// 	return "aspirations"
// }
