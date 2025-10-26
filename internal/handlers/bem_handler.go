package handlers

import (
	"net/http"
	"strconv"
	"math"
	"fmt"
	"gorm.io/gorm"

	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"
	"github.com/gin-gonic/gin"
)

// bemHandler handles HTTP requests related to bems
type BemHandler struct {
	service *services.BemService
}

// NewbemHandler creates a new bem handler
func NewBemHandler(db *gorm.DB) *BemHandler {
	return &BemHandler{
		service: services.NewBemService(db),
	}
}

// GetAllbems returns all bems
func (h *BemHandler) GetAllBems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

    if page < 1 {
        page = 1
    }
    if perPage < 1 {
        perPage = 10
    }

    offset := (page - 1) * perPage

    // ambil data + total count
    students, total, err := h.service.GetAllBems(perPage, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
        return
    }

    totalPages := int(math.Ceil(float64(total) / float64(perPage)))

    // siapkan metadata
    metadata := utils.PaginationMetadata{
        CurrentPage: page,
        PerPage:     perPage,
        TotalItems:  int(total),
        TotalPages:  totalPages,
        Links: utils.PaginationLinks{
            First: fmt.Sprintf("/students?page=1&per_page=%d", perPage),
            Last:  fmt.Sprintf("/students?page=%d&per_page=%d", totalPages, perPage),
        },
    }

    // response dengan metadata
    response := utils.MetadataFormatResponse(
        "success",
        "Berhasil list mendapatkan data",
        metadata,
        students,
    )

    c.JSON(http.StatusOK, response)
}

// GetbemByID returns a bem by ID
func (h *BemHandler) GetBemByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stats := c.Query("stats")
	var result interface{}

	if stats == "true" {
		result, err = h.service.GetBemWithStats(uint(id))
	} else {
		result, err = h.service.GetBemByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bem not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Bem retrieved successfully",
		"data":    result,
	})
}

// Createbem creates a new bem
func (h *BemHandler) CreateBem(c *gin.Context) {
	var bem models.BEM

	if err := c.ShouldBindJSON(&bem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.CreateBem(&bem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Bem created successfully",
		"data":    bem,
	})
}

// Updatebem updates a bem
func (h *BemHandler) UpdateBem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var bem models.BEM
	if err := c.ShouldBindJSON(&bem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	bem.ID = uint(id)

	if err := h.service.UpdateBem(&bem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Bem updated successfully",
		"data":    bem,
	})
}

// Deletebem deletes a bem
func (h *BemHandler) DeleteBem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteBem(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Bem deleted successfully",
	})
} 

func (h *BemHandler) GetAllLeaders(c *gin.Context) {
	students, err := h.service.GetAllLeaders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	c.JSON(http.StatusOK, students)
}
