package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

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
		"Berhasil mendapatkan daftar pengumuman",
		metadata,
		aspirations,
	)

	c.JSON(http.StatusOK, response)
}

func (h *AspirationHandler) CreateAspiration(c *gin.Context) {
	var aspiration models.Aspiration
	aspiration.Title = c.PostForm("title")
	aspiration.Description = c.PostForm("description")
	aspiration.Category = c.PostForm("category")
	aspiration.PriorityLevel = c.PostForm("priority_level")
	if aspiration.Title == "" || aspiration.Description == "" || aspiration.Category == "" || aspiration.PriorityLevel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (title, description, category, priority) are required"})
		return
	}
	if err := h.service.CreateAspiration(&aspiration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Aspirasi berhasil dikirim secara anonim",
		"data":    aspiration,
	})
}
