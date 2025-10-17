package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
type AspirationResponse struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	Content      string `json:"content"`
	Priority     string `json:"priority_level"`
	UserID       uint   `json:"user_id"`
	StudentName  string `json:"student_name"`
	CreatedAt    string `json:"created_at"`
	
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

	var responseData []AspirationResponse
	for _, a := range aspirations {
		responseData = append(responseData, AspirationResponse{
			ID:          a.ID,
			Content:     a.Content,
			Title:       a.Title,
			Description: a.Description,
			Category:    a.Category,
			Priority:    a.PriorityLevel,
			StudentName: a.Student.FullName,
			CreatedAt: a.CreatedAt.Local().Format("2006-01-02 15:04:05"),
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

	// ✅ Ambil userID dari context (harus sama key-nya dengan yang diset di middleware)
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User belum login atau token tidak valid",
		})
		return
	}

	// ✅ Konversi userID ke int dengan aman
	var userID int
	switch v := userIDVal.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	case uint:
		userID = int(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Tipe userID tidak dikenali: %T", userIDVal),
		})
		return
	}

	// ✅ Deteksi content-type dan lakukan binding sesuai jenisnya
	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		if err := c.ShouldBindJSON(&aspiration); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid JSON body: " + err.Error(),
			})
			return
		}
	} else {
		// fallback untuk form-data
		aspiration.Title = c.PostForm("title")
		aspiration.Description = c.PostForm("description")
		aspiration.Category = c.PostForm("category")
		aspiration.Content = c.PostForm("content")
		aspiration.PriorityLevel = c.PostForm("priority_level")
	}

	// ✅ Trim spasi untuk memastikan data bersih
	aspiration.Title = strings.TrimSpace(aspiration.Title)
	aspiration.Description = strings.TrimSpace(aspiration.Description)
	aspiration.Category = strings.TrimSpace(aspiration.Category)
	aspiration.Content = strings.TrimSpace(aspiration.Content)
	aspiration.PriorityLevel = strings.TrimSpace(aspiration.PriorityLevel)

	// ✅ Set user_id
	aspiration.UserID = userID

	// ✅ Validasi input
	if aspiration.Title == "" ||
		aspiration.Description == "" ||
		aspiration.Category == "" ||
		aspiration.PriorityLevel == "" ||
		aspiration.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Semua field (title, description, category, content, priority_level) wajib diisi",
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
