package models

import (
	"time"

	"gorm.io/gorm"
)

type Request struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"size:255;not null"`
	Category       int            `json:"category"`
	Item           string         `json:"item" gorm:"not null"`
	Activity       string         `json:"activity" gorm:"not null"`
	Location       string         `json:"location" gorm:"not null"`
	RequestPlan    string         `json:"request_plan" gorm:"not null"`
	ReturnPlan     string         `json:"return_plan" gorm:"not null"`
	RequesterID    string         `json:"requester_id" gorm:"not null"`
	ApproverID     string         `json:"approver_id" gorm:"default:null"`
	ImageURLKTM    string         `json:"image_url_ktm" gorm:"type:varchar(255)"`
	ImageURLBRG    string         `json:"image_url_barang" gorm:"type:varchar(255)"`
	Status         string         `json:"status" gorm:"default:'pending';not null"`
	Reason         string         `json:"reason"`
	ReturnAt       time.Time      `json:"return_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index;uniqueIndex:idx_courses_code_deleted_at" json:"deleted_at,omitempty"`
	OrganizationID uint           `json:"organization_id"`
	Organization   *Organization  `json:"organization" gorm:"foreignKey:ID;references:OrganizationID;constraint:OnUpdate:CASCADE"`
}

type RequestWithStats struct {
	Request   Request  `json:"request"`
	ItemNames []string `json:"item_names"`
}

func (Request) TableName() string {
	return "requests"
}
