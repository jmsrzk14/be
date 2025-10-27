package services

import (
	"bem_be/internal/models"
	"bem_be/internal/repositories"
	"time"
)

type NotificationService struct {
	repo *repositories.NotificationRepository
}

func NewNotificationService(repo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// Ambil semua notif (umum)
func (s *NotificationService) GetAllNotifications() ([]models.Notification, error) {
	return s.repo.GetAllNotifications()
}

// Ambil notif user, termasuk status is_read
func (s *NotificationService) GetUserNotifications(username string) ([]models.UserNotification, error) {
	return s.repo.GetUserNotifications(username)
}

func (s *NotificationService) CreateNotification(title, message string) (*models.Notification, error) {
	notification := &models.Notification{
		Title:   title,
		Message: message,
	}
	err := s.repo.CreateNotification(notification)
	return notification, err
}

func (s *NotificationService) MarkNotificationAsRead(username, notificationID string) error {
    return s.repo.CreateUserNotificationIfNotExists(username, notificationID)
}

type NotificationWithRead struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	IsRead    bool      `json:"is_read"`
}

func (s *NotificationService) GetAllNotificationsForUser(username string) ([]NotificationWithRead, error) {
	notifications, err := s.repo.GetAllNotifications()
	if err != nil {
		return nil, err
	}

	// ambil UserNotifications untuk user
	readMap, err := s.repo.GetUserNotificationsMap(username)
	if err != nil {
		return nil, err
	}

	result := make([]NotificationWithRead, len(notifications))
	for i, n := range notifications {
		result[i] = NotificationWithRead{
			ID:        n.ID,
			Title:     n.Title,
			Message:   n.Message,
			CreatedAt: n.CreatedAt,
			IsRead:    readMap[n.ID], // true jika user sudah baca
		}
	}

	return result, nil
}
