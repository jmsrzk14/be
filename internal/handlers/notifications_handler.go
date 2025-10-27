package handlers

import (
	"bem_be/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// Ambil semua notif untuk user (status read)
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
    username := c.Param("username")
    if username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Username diperlukan"})
        return
    }

    notifications, err := h.service.GetAllNotificationsForUser(username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":        "success",
        "notifications": notifications,
    })
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
    username := c.Param("username")
    notificationID := c.Param("notificationID")

    if username == "" || notificationID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Username dan notificationID diperlukan",
        })
        return
    }

    if err := h.service.MarkNotificationAsRead(username, notificationID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Notifikasi ditandai sebagai dibaca",
    })
}
