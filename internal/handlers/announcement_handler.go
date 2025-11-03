package handlers

import "strings"
import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"bem_be/internal/models"
	"bem_be/internal/services"
)

// AnnouncementHandler handles HTTP requests related to announcements
// di AnnouncementHandler
type AnnouncementHandler struct {
	service *services.AnnouncementService
	db      *gorm.DB
	notificationService *services.NotificationService
}

func NewAnnouncementHandler(db *gorm.DB, notificationService *services.NotificationService) *AnnouncementHandler {
	return &AnnouncementHandler{
		service: services.NewAnnouncementService(db),
		notificationService: notificationService,
	}
}

// GetAllAnnouncements returns all announcements
// GetAllAnnouncements returns all announcements with pagination and optional filters
func (h *AnnouncementHandler) GetAllAnnouncement(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	title := c.DefaultQuery("title", "")
	content := c.DefaultQuery("content", "")
	category := c.DefaultQuery("category", "")

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	announcementList, total, err := h.service.GetAllAnnouncements(page, perPage, title, content, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mendapatkan daftar berita",
		"metadata": gin.H{
			"current_page": page,
			"per_page":     perPage,
			"total_items":  total,
			"total_pages":  totalPages,
		},
		"filters": gin.H{
			"title":    title,
			"content":  content,
			"category": category,
		},
		"data": announcementList,
	})
}


// GetAnnouncementByID returns an announcement by ID
func (h *AnnouncementHandler) GetAnnouncementByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stats := c.Query("stats")
	var result interface{}

	if stats == "true" {
		result, err = h.service.GetAnnouncementWithStats(uint(id))
	} else {
		result, err = h.service.GetAnnouncementByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Announcement retrieved successfully",
		"data":    result,
	})
}
func formatPosition(pos string) string {
	// Ganti underscore dengan spasi, lalu capital setiap kata
	pos = strings.ReplaceAll(pos, "_", " ")
	return strings.Title(pos) // "ketua_bem" -> "Ketua Bem"
}

// CreateAnnouncement creates a new announcement (with optional file)
func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	var announcement models.Announcement

	extUserID, exists := c.Get("userID")

	var authorID uint = 0
	var orgID *uint = nil

	if exists {
		var student models.Student
		if err := h.db.Where("user_id = ?", extUserID).First(&student).Error; err == nil {
			authorID = uint(student.UserID)
			if student.OrganizationID != nil {
				orgIDVal := uint(*student.OrganizationID)
				orgID = &orgIDVal
			}
		}
	}

	// Isi field announcement
	announcement.Title = c.PostForm("title")
	announcement.Content = c.PostForm("content")
	announcement.AuthorID = authorID
	announcement.OrganizationID = orgID

	layout := "2006-01-02"
	if startDateStr := c.PostForm("start_date"); startDateStr != "" {
		startDate, err := time.Parse(layout, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use YYYY-MM-DD"})
			return
		}
		announcement.StartDate = &startDate
	}

	if endDateStr := c.PostForm("end_date"); endDateStr != "" {
		endDate, err := time.Parse(layout, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use YYYY-MM-DD"})
			return
		}
		announcement.EndDate = &endDate
	}

	file, err := c.FormFile("file")
	if err == nil {
		uploadPath := "uploads"
		os.MkdirAll(uploadPath, os.ModePerm)

		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		filePath := filepath.Join(uploadPath, fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		announcement.FileURL = fileName
	}

	if err := h.service.Createannouncement(&announcement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	title := "Pengumuman Baru: " + announcement.Title
	message := fmt.Sprintf("Pengumuman baru telah dibuat. Cek sekarang!")

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
		"status":  "success",
		"message": "Announcement created successfully",
		"data":    announcement,
		"notification": createdNotif,
	})
}

// UpdateAnnouncement updates an existing announcement (with optional file update)
func (h *AnnouncementHandler) UpdateAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var announcement models.Announcement
	announcement.ID = uint(id)
	announcement.Title = c.PostForm("title")
	announcement.Content = c.PostForm("content")

	file, err := c.FormFile("file")
	if err == nil {
		uploadPath := "uploads/announcements"
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create upload folder"})
			return
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		filePath := filepath.Join(uploadPath, fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		announcement.FileURL = filePath
	}

	if err := h.service.Updateannouncement(&announcement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Announcement updated successfully",
		"data":    announcement,
	})
}

// DeleteAnnouncement deletes an announcement
func (h *AnnouncementHandler) DeleteAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteAnnouncement(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Announcement deleted successfully",
	})
}
