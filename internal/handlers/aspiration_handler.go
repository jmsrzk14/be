package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"

	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"
)

type AspirationHandler struct {
	service *services.AspirationService
	db      *gorm.DB
}

func NewAspirationHandler(db *gorm.DB) *AspirationHandler {
	return &AspirationHandler{
		service: services.NewAspirationService(db),
		db:      db,
	}
}

type Aspiration struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	UserName      string `json:"user_name" gorm:"type:varchar(20);not null"`
	Title         string `json:"title" gorm:"not null"`
	Description   string `json:"description" gorm:"type:text;not null"`
	Category      string `json:"category" gorm:"type:text;not null"`
	Content       string `json:"content" gorm:"type:text;not null"`
	PriorityLevel string `json:"priority_level" gorm:"type:text;not null"`
	CreatedAt     string `json:"created_at"`
	StudentName   string `json:"student_name"`
}

func (h *AspirationHandler) GetAllAspirations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	aspirations, total, err := h.service.GetAllAspirations(perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	var responseData []Aspiration
	for _, a := range aspirations {
		responseData = append(responseData, Aspiration{
			ID:          a.ID,
			Content:     a.Content,
			Title:       a.Title,
			Description: a.Description,
			Category:    a.Category,
			PriorityLevel:    a.PriorityLevel,
			StudentName: a.Student.FullName,
			CreatedAt:   a.CreatedAt.Local().Format("2006-01-02 15:04:05"),
			// ambil nama mahasiswa
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	metadata := utils.PaginationMetadata{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  int(total),
		TotalPages:  totalPages,
		Links: utils.PaginationLinks{
			First: fmt.Sprintf("/aspirations?page=1&per_page=%d", perPage),
			Last:  fmt.Sprintf("/aspirations?page=%d&per_page=%d", totalPages, perPage),
		},
	}

	response := utils.MetadataFormatResponse(
		"success",
		"Berhasil mendapatkan daftar aspirasi",
		metadata,
		responseData,
	)

	c.JSON(http.StatusOK, response)
}
func (h *AspirationHandler) CreateAspiration(c *gin.Context) {
	var aspiration models.Aspiration

	// ✅ Baca body JSON dari frontend
	if err := c.ShouldBindJSON(&aspiration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format JSON tidak valid: " + err.Error(),
		})
		return
	}

	// ✅ Validasi input wajib
	if aspiration.Title == "" ||
		aspiration.Description == "" ||
		aspiration.Category == "" ||
		aspiration.PriorityLevel == "" ||
		aspiration.Content == "" ||
		aspiration.UserName == "" { // pastikan username dikirim dari frontend
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Semua field wajib diisi (title, description, category, content, priority_level, username)",
		})
		return
	}

	// ✅ Simpan ke DB
	if err := h.service.CreateAspiration(&aspiration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// ✅ Response sukses
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Aspirasi berhasil dikirim",
		"data":    aspiration,
	})
}
