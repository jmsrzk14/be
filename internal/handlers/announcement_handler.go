package handlers

import (
	"net/http"
	"strconv"
	"math"
	"fmt"
	"path/filepath"
	"os"
	"time"

	"gorm.io/gorm"
	"github.com/gin-gonic/gin"

	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"
)

// AnnouncementHandler handles HTTP requests related to announcements
// di AnnouncementHandler
type AnnouncementHandler struct {
    service *services.AnnouncementService
    db      *gorm.DB
}

func NewAnnouncementHandler(db *gorm.DB) *AnnouncementHandler {
    return &AnnouncementHandler{
        service: services.NewAnnouncementService(db),
        db:      db,
    }
}


// GetAllAnnouncements returns all announcements
func (h *AnnouncementHandler) GetAllAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	announcements, total, err := h.service.GetAllAnnouncements(perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	metadata := utils.PaginationMetadata{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  int(total),
		TotalPages:  totalPages,
		Links: utils.PaginationLinks{
			First: fmt.Sprintf("/announcements?page=1&per_page=%d", perPage),
			Last:  fmt.Sprintf("/announcements?page=%d&per_page=%d", totalPages, perPage),
		},
	}

	response := utils.MetadataFormatResponse(
		"success",
		"Berhasil mendapatkan daftar pengumuman",
		metadata,
		announcements,
	)

	c.JSON(http.StatusOK, response)
}

// GetAnnouncementByID returns an announcement by ID
func (h *AnnouncementHandler) GetAnnouncementByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.service.GetAnnouncementWithStats(uint(id))
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

// CreateAnnouncement creates a new announcement (with optional file)
func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
    var announcement models.Announcement

    // Ambil external user id dari JWT context
    extUserID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Cari student berdasarkan external user id
    var student models.Student
    if err := h.db.Where("user_id = ?", extUserID).First(&student).Error; err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Student not found"})
        return
    }

    // Isi field announcement
    announcement.Title = c.PostForm("title")
    announcement.Content = c.PostForm("content")
    announcement.AuthorID = uint(student.UserID)   // pakai user_id dari students
    announcement.Position = student.Position       // 🔥 ambil position dari tabel students

    if student.OrganizationID != nil {
        orgID := uint(*student.OrganizationID)
        announcement.OrganizationID = &orgID
    }

    // Handle file upload
    file, err := c.FormFile("file")
    if err == nil {
        uploadPath := "uploads/announcements"
        os.MkdirAll(uploadPath, os.ModePerm)
        fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
        filePath := filepath.Join(uploadPath, fileName)

        if err := c.SaveUploadedFile(file, filePath); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
            return
        }
        announcement.FileURL = filePath
    }

    // Simpan ke DB
    if err := h.service.Createannouncement(&announcement); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "status":  "success",
        "message": "Announcement created successfully",
        "data":    announcement,
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
