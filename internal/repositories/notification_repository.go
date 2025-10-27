package repositories

import (
	"bem_be/internal/models"
	"strconv"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateNotification(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) GetUserNotifications(username string) ([]models.UserNotification, error) {
	var userNotifs []models.UserNotification
	err := r.db.Preload("Notification").
		Where("username = ?", username).
		Order("id DESC").
		Find(&userNotifs).Error
	return userNotifs, err
}

func (r *NotificationRepository) MarkAsRead(username string, notificationID uint) error {
	return r.db.Model(&models.UserNotification{}).
		Where("username = ? AND notification_id = ?", username, notificationID).
		Update("is_read", true).Error
}

func (r *NotificationRepository) ExistsUserNotification(username string, notificationID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserNotification{}).
		Where("username = ? AND notification_id = ?", username, notificationID).
		Count(&count).Error
	return count > 0, err
}

func (r *NotificationRepository) CreateUserNotification(userNotification *models.UserNotification) error {
	return r.db.Create(userNotification).Error
}

// Ambil semua notifikasi
func (r *NotificationRepository) GetAllNotifications() ([]models.Notification, error) {
	var notifications []models.Notification
	if err := r.db.Order("created_at desc").Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

// Ambil semua UserNotification untuk user tertentu dalam bentuk map[notification_id]bool
func (r *NotificationRepository) GetUserNotificationsMap(username string) (map[uint]bool, error) {
	var userNotifs []models.UserNotification
	result := make(map[uint]bool)

	if err := r.db.Where("username = ?", username).Find(&userNotifs).Error; err != nil {
		return nil, err
	}

	for _, u := range userNotifs {
		result[u.NotificationID] = true
	}

	return result, nil
}

func (r *NotificationRepository) CreateUserNotificationIfNotExists(username, notificationID string) error {
    // cek apakah sudah ada
    var existing models.UserNotification
    err := r.db.Where("username = ? AND notification_id = ?", username, notificationID).First(&existing).Error
    if err == nil {
        // sudah ada, tidak perlu buat baru
        return nil
    }
    if err != gorm.ErrRecordNotFound {
        return err
    }

    // buat record baru karena belum ada
    userNotif := models.UserNotification{
        Username:       username,
        NotificationID: parseUint(notificationID),
        IsRead:         true,
    }

    return r.db.Create(&userNotif).Error
}

// helper untuk konversi string ke uint
func parseUint(s string) uint {
    id, _ := strconv.ParseUint(s, 10, 64)
    return uint(id)
}