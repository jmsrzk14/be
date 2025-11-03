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

// mpmHandler handles HTTP requests related to mpms
type MpmHandler struct {
	service *services.MpmService
}

// NewmpmHandler creates a new mpm handler
func NewMpmHandler(db *gorm.DB) *MpmHandler {
	return &MpmHandler{
		service: services.NewMpmService(db),
	}
}

// GetAllmpms returns all mpms
func (h *MpmHandler) GetAllMpms(c *gin.Context) {
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
    students, total, err := h.service.GetAllMpms(perPage, offset)
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

// GetmpmByID returns a mpm by ID
func (h *MpmHandler) GetMpmByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stats := c.Query("stats")
	var result interface{}

	if stats == "true" {
		result, err = h.service.GetMpmWithStats(uint(id))
	} else {
		result, err = h.service.GetMpmByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "mpm not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "mpm retrieved successfully",
		"data":    result,
	})
}

// Creatempm creates a new mpm
func (h *MpmHandler) CreateMpm(c *gin.Context) {
	var mpm models.MPM

	if err := c.ShouldBindJSON(&mpm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.CreateMpm(&mpm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "mpm created successfully",
		"data":    mpm,
	})
}

// Updatempm updates a mpm
func (h *MpmHandler) UpdateMpm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var mpm models.MPM
	if err := c.ShouldBindJSON(&mpm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	mpm.ID = uint(id)

	if err := h.service.UpdateMpm(&mpm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "mpm updated successfully",
		"data":    mpm,
	})
}

// Deletempm deletes a mpm
func (h *MpmHandler) DeleteMpm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteMpm(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mpm deleted successfully",
	})
} 

func (h *MpmHandler) GetMPMByPeriod(c *gin.Context) {
	period := c.Param("period")

	mpm, err := h.service.GetMPMByPeriod(period)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "mpm not found"})
		return
	}

	c.JSON(http.StatusOK, mpm)
}

func (h *MpmHandler) GetMpmPeriod(c *gin.Context) {
	period := c.Param("period")

	data, err := h.service.GetMpmPeriod(period)
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
