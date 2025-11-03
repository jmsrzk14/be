// handlers/event_handler.go
package handlers

import (
	"net/http"
	"strconv"
	"time"
	"fmt"

	"bem_be/internal/models"
	"bem_be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventHandler struct {
	service *services.CalenderService
	db      *gorm.DB
	notificationService *services.NotificationService
}

func NewEventHandler(db *gorm.DB, notificationService *services.NotificationService) *EventHandler {
	return &EventHandler{
		service: services.NewCalenderService(db),
		db:      db,
		notificationService: notificationService,
	}
}

// GET /events/month?month=10&year=2025
func (h *EventHandler) GetEventsByMonth(c *gin.Context) {
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	if monthStr == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month and year are required"})
		return
	}

	month, err1 := strconv.Atoi(monthStr)
	year, err2 := strconv.Atoi(yearStr)
	if err1 != nil || err2 != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month or year"})
		return
	}

	// Tentukan rentang waktu awal dan akhir bulan
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond) // hari terakhir bulan itu

	var events []models.Calender
	// ambil event yang overlap rentang bulan
	if err := h.db.Where("end_time >= ? AND start_time <= ?", start, end).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) GetEventsCurrentMonth(c *gin.Context) {
	events, month, year, err := h.service.GetEventsCurrentMonth()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"month":  month,
		"year":   year,
		"events": events,
	})
}

// POST /events
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var payload models.Calender
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validation dasar
	if payload.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	if payload.EndTime.Before(payload.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end must be after start"})
		return
	}

	// Simpan (timestamps otomatis oleh GORM)
	if err := h.db.Create(&payload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	title := "Event Baru dari " + payload.Organization.ShortName
	message := fmt.Sprintf("%s menambahkan membuat kegiatan %s. Cek sekarang!", payload.Organization.ShortName, payload.Title)

	// Buat instance Notification
	notification := &models.Notification{
		Title:   title,
		Message: message,
	}

	// Simpan ke database menggunakan service
	createdNotif, err := h.notificationService.CreateNotification(notification.Title, notification.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat notifikasi berita",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":       "success",
		"message":      "Event berhasil dibuat",
		"data":         payload,
		"notification": createdNotif,
	})
}

// PUT /events/:id
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var existing models.Calender
	if err := h.db.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	var payload models.Calender
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// update fields
	existing.Title = payload.Title
	existing.Description = payload.Description
	existing.Location = payload.Location
	existing.StartTime = payload.StartTime
	existing.EndTime = payload.EndTime

	if existing.EndTime.Before(existing.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end must be after start"})
		return
	}

	if err := h.db.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existing)
}

// DELETE /events/:id
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := h.db.Delete(&models.Calender{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

