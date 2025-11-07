package handlers

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"

	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"

	"github.com/gin-gonic/gin"
)

// AssociationHandler handles HTTP requests related to associations
type AssociationHandler struct {
	service             *services.AssociationService
	notificationService *services.NotificationService
}

// NewAssociationHandler creates a new association handler
func NewAssociationHandler(db *gorm.DB, notificationService *services.NotificationService) *AssociationHandler {
	return &AssociationHandler{
		service:             services.NewAssociationService(db),
		notificationService: notificationService,
	}
}

func (h *AssociationHandler) GetAllAssociations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("name") // pencarian pakai param ?name=

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	associations, total, err := h.service.GetAllAssociations(perPage, offset, search)
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
	}

	response := utils.MetadataFormatResponse(
		"success",
		"Berhasil list mendapatkan data associations",
		metadata,
		associations,
	)

	c.JSON(http.StatusOK, response)
}

func (h *AssociationHandler) GetAllAssociationsGuest(c *gin.Context) {
	// ambil semua data tanpa limit & offset
	associations, err := h.service.GetAllAssociationsGuest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	// langsung response tanpa metadata
	response := utils.ResponseHandler(
		"success",
		"Berhasil mendapatkan data",
		associations,
	)

	c.JSON(http.StatusOK, response)
}

// GetAssociationByID returns a association by ID
func (h *AssociationHandler) GetAssociationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stats := c.Query("stats")
	var result interface{}

	if stats == "true" {
		result, err = h.service.GetAssociationWithStats(uint(id))
	} else {
		result, err = h.service.GetAssociationByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Association not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Association retrieved successfully",
		"data":    result,
	})
}

func (h *AssociationHandler) GetAssociationByShortName(c *gin.Context) {
	shortName := c.Param("shortName")

	// Optional query param, misal nanti dipakai untuk stats
	stats := c.Query("stats")
	_ = stats // sementara tidak digunakan, bisa dihapus jika tidak dipakai

	result, err := h.service.GetAssociationByShortName(shortName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Association not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Association retrieved successfully",
		"data":    result,
	})
}

// CreateAssociation creates a new association
func (h *AssociationHandler) CreateAssociation(c *gin.Context) {
	var association models.Organization

	// ambil field manual (biar gak coba bind file ke string)
	association.Name = c.PostForm("name")
	association.ShortName = c.PostForm("short_name")

	association.CategoryID = 3

	// ambil file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logo file is required"})
		return
	}

	// kirim ke service
	if err := h.service.CreateAssociation(&association, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	title := "Himpunan Baru: " + association.Name
	message := fmt.Sprintf("Himpunan baru telah dibuat. Cek sekarang!")

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

	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"message":      "Association created successfully",
		"data":         association,
		"notification": createdNotif,
	})
}

// UpdateAssociation updates a association
func (h *AssociationHandler) UpdateAssociation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	type UpdateInput struct {
		Name      string `form:"name" binding:"omitempty" json:"name"`
		ShortName string `form:"short_name" binding:"omitempty" json:"short_name"`
	}
	var input UpdateInput
	if err := c.ShouldBind(&input); err != nil {
		_ = c.ShouldBind(&input)
	}

	var association models.Organization
	association.ID = uint(id)

	if input.Name != "" {
		association.Name = input.Name
	}
	if input.ShortName != "" {
		association.ShortName = input.ShortName
	}

	// Handle image upload if provided
	file, err := c.FormFile("image")
	if err == nil {
		uploadPath := "uploads/associations"
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		// Generate unique filename
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		filePath := filepath.Join(uploadPath, fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		association.Image = fileName // This will be non-empty, so repo will update it
	}

	// Call service to update (repo will only update non-zero fields)
	if err := h.service.UpdateAssociation(&association); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Association updated successfully",
		"data":    association,
	})
}

// DeleteAssociation deletes a association
func (h *AssociationHandler) DeleteAssociation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteAssociation(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Association deleted successfully",
	})
}

func (h *AssociationHandler) GetAdminAssociations(c *gin.Context) {
	shortName := c.Param("shortName")
	period := c.Param("period")

	data, err := h.service.GetAdminAssociation(shortName, period)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
