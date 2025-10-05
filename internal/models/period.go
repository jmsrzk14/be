package models

import (
	"time"

	"gorm.io/gorm"
)

// Organization represents a club in the system
type Period struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrganizationID int            `json:"organization_id" gorm:"not null"`
	Period         string         `gorm:"not null" json:"period"`
	Vision         string         `gorm:"not null" json:"vision"`
	Mission        string         `gorm:"not null" json:"mission"`
	Workplan       string         `json:"workplan" gorm:"type:text"`
	LeaderID       uint           `json:"leader_id"`
	Leader         *Student       `json:"leader" gorm:"foreignKey:ID;references:LeaderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CoLeaderID     uint           `json:"coleader_id"`
	CoLeader       *Student       `json:"coleader" gorm:"foreignKey:ID;references:CoLeaderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Secretary1ID   uint           `json:"secretary1_id"`
	Secretary1     *Student       `json:"secretary1" gorm:"foreignKey:ID;references:Secretary1;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Secretary2ID   uint           `json:"secretary_id"`
	Secretary2     *Student       `json:"secretary" gorm:"foreignKey:ID;references:Secretary2ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Treasurer1ID   uint           `json:"treasurer1_id"`
	Treasurer1     *Student       `json:"treasurer1" gorm:"foreignKey:ID;references:Treasurer1ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Treasurer2ID   uint           `json:"treasurer2_id"`
	Treasurer2     *Student       `json:"treasurer2" gorm:"foreignKey:ID;references:Treasurer2ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index;uniqueIndex:idx_courses_code_deleted_at" json:"deleted_at,omitempty"`
}
