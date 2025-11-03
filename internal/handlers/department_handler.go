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

// DepartmentHandler handles HTTP requests related to departments
type DepartmentHandler struct {
	service             *services.DepartmentService
	notificationService *services.NotificationService
}

// NewDepartmentHandler creates a new department handler
func NewDepartmentHandler(db *gorm.DB, notificationService *services.NotificationService) *DepartmentHandler {
	return &DepartmentHandler{
		service:             services.NewDepartmentService(db),
		notificationService: notificationService,
	}
}

func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("name")

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	departments, total, err := h.service.GetAllDepartments(perPage, offset, search)
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
		"Berhasil list mendapatkan data departments",
		metadata,
		departments,
	)

	c.JSON(http.StatusOK, response)
}

func (h *DepartmentHandler) GetAllOrganizations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("name")

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	departments, total, err := h.service.GetAllOrganizations(perPage, offset, search)
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
		"Berhasil list mendapatkan data departments",
		metadata,
		departments,
	)

	c.JSON(http.StatusOK, response)
}

func (h *DepartmentHandler) GetAllDepartmentsGuest(c *gin.Context) {
	// ambil semua data tanpa limit & offset
	departments, err := h.service.GetAllDepartmentsGuest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	// langsung response tanpa metadata
	response := utils.ResponseHandler(
		"success",
		"Berhasil mendapatkan data",
		departments,
	)

	c.JSON(http.StatusOK, response)
}

// GetDepartmentByID returns a department by ID
func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stats := c.Query("stats")
	var result interface{}

	if stats == "true" {
		result, err = h.service.GetDepartmentWithStats(uint(id))
	} else {
		result, err = h.service.GetDepartmentByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Department retrieved successfully",
		"data":    result,
	})
}

// CreateDepartment creates a new department
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var department models.Organization

	// ambil field manual (biar gak coba bind file ke string)
	department.Name = c.PostForm("name")
	department.ShortName = c.PostForm("short_name")

	department.CategoryID = 2

	// ambil file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logo file is required"})
		return
	}

	// kirim ke service
	if err := h.service.CreateDepartment(&department, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	title := "Departemen Baru: " + department.Name
	message := fmt.Sprintf("Departemen baru telah dibuat. Cek sekarang!")

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
		"message":      "Department created successfully",
		"data":         department,
		"notification": createdNotif,
	})
}

func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	type UpdateInput struct {
		Name      string `form:"name" binding:"omitempty"`
		ShortName string `form:"short_name" binding:"omitempty"`
	}
	var input UpdateInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	var department models.Organization
	department.ID = uint(id)

	if input.Name != "" {
		department.Name = input.Name
	}
	if input.ShortName != "" {
		department.ShortName = input.ShortName
	}

	// Handle image upload if provided
	file, err := c.FormFile("image")
	if err == nil {
		uploadPath := "uploads/departments"
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

		department.Image = fileName // This will be non-empty, so repo will update it
	}

	// Call service to update (repo will only update non-zero fields)
	if err := h.service.UpdateDepartment(&department); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Department updated successfully",
		"data":    department,
	})
}

// DeleteDepartment deletes a department
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteDepartment(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Department deleted successfully",
	})
}
