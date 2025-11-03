package models

import (
	"time"

	"gorm.io/gorm"
)

// BEM represents the main student executive board for a specific period.
type MPM struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Vision      string         `gorm:"not null" json:"vision"`
	Mission     string         `gorm:"not null" json:"mission"`
	LeaderID    uint           `json:"leader_id"`
	Leader      *Student       `json:"leader" gorm:"foreignKey:LeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CoLeaderID  uint           `json:"coleader_id"`
	CoLeader    *Student       `json:"coleader" gorm:"foreignKey:CoLeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SecretaryID uint           `json:"secretary_id"`
	Secretary   *Student       `json:"secretary" gorm:"foreignKey:SecretaryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Period      string         `json:"period" gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (MPM) TableName() string {
	return "mpm"
}
