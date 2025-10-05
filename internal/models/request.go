package models

import (
	"time"

	"gorm.io/gorm"
)

type Request struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"size:255;not null"`
	Item           string         `json:"item" gorm:"not null"`
	Activity       string         `json:"activity" gorm:"not null"`
	Location       string         `json:"location" gorm:"not null"`
	RequestPlan    string         `json:"request_plan" gorm:"not null"`
	ReturnPlan     string         `json:"return_plan" gorm:"not null"`
	RequesterID    uint           `json:"requester_id" gorm:"not null"`
	ApproverID     uint           `json:"approver_id" gorm:"default:null"`
	Requester      *User          `json:"requester,omitempty" gorm:"foreignKey:RequesterID"`
	Approver       *User          `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`
	ImageURLKTM    string         `json:"image_url_ktm" gorm:"type:varchar(255)"`
	ImageURLBRG    string         `json:"image_url_barang" gorm:"type:varchar(255)"`
	Status         string         `json:"status" gorm:"type:enum('pending', 'approved', 'rejected');default:'pending';not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index;uniqueIndex:idx_courses_code_deleted_at" json:"deleted_at,omitempty"`
	OrganizationID uint           `json:"organization_id"`
	Organization   *Organization  `json:"organization" gorm:"foreignKey:ID;references:OrganizationID;constraint:OnUpdate:CASCADE"`
}

func (Request) TableName() string {
	return "requests"
}
