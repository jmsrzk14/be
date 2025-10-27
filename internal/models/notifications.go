package models

import (
	"time"
)

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type UserNotification struct {
	ID             uint          `json:"id" gorm:"primaryKey"`
	Username       string        `json:"username" gorm:"index"`
	NotificationID uint          `json:"notification_id"`
	Notification   *Notification `json:"notification" gorm:"foreignKey:NotificationID;references:ID;constraint:OnDelete:CASCADE"`
	IsRead         bool          `json:"is_read" gorm:"default:false"`
}
