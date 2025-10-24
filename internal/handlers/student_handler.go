package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"bem_be/internal/database"
	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"

	"github.com/gin-gonic/gin"
)

// StudentHandler handles HTTP requests related to students
type StudentHandler struct {
	service *services.StudentService
}

// NewStudentHandler creates a new student handler
func NewStudentHandler(db *gorm.DB, campusAuth *services.CampusAuthService) *StudentHandler {
	return &StudentHandler{
		service: services.NewStudentService(db, campusAuth),
	}
}

func (h *StudentHandler) GetAllStudents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// ambil query params
	search := c.Query("name") // pakai param "name" untuk search
	studyProgram := c.Query("study_program")
	yearEnrolledStr := c.Query("year_enrolled")
	yearEnrolled, _ := strconv.Atoi(yearEnrolledStr) // default 0 kalau kosong / invalid

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// lempar semua filter ke service
	students, total, err := h.service.GetAllStudents(perPage, offset, search, studyProgram, yearEnrolled)
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
		"Berhasil list mendapatkan data",
		metadata,
		students,
	)

	c.JSON(http.StatusOK, response)
}

// GetStudentByID returns a student by ID
func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	student, err := h.service.GetStudentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Ambil data organisasi berdasarkan OrganizationID
	var organizationName string
	if student.OrganizationID != nil && *student.OrganizationID != 0 {
		var organization models.Organization
		if err := database.DB.First(&organization, *student.OrganizationID).Error; err == nil {
			organizationName = organization.Name
		}
	}

	// Gabungkan ke response JSON
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Student retrieved successfully",
		"data": gin.H{
			"id":              student.ID,
			"name":            student.FullName,
			"email":           student.Email,
			"organization_id": student.OrganizationID,
			"organization":    organizationName, // tampilkan nama organisasi
			// tambahkan field lain dari student kalau perlu
		},
	})
}

// GetStudentByUserID returns a student by their user ID from the campus system
func (h *StudentHandler) GetStudentByUserID(c *gin.Context) {
	username := c.Param("username")

	student, err := h.service.GetStudentByUserID(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Student retrieved successfully",
		"data":    student,
	})
}

// SyncStudents syncs students from the campus API
func (h *StudentHandler) SyncStudents(c *gin.Context) {
	// Sync students using the service
	count, err := h.service.SyncStudents()
	if err != nil {
		errMsg := err.Error()
		statusCode := http.StatusInternalServerError
		responseMsg := "Failed to sync students"

		// Check for specific errors to provide better messages
		if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded") {
			statusCode = http.StatusGatewayTimeout
			responseMsg = "Connection to campus API timed out"
		} else if strings.Contains(errMsg, "connection refused") {
			statusCode = http.StatusServiceUnavailable
			responseMsg = "Campus API service unavailable"
		}

		c.JSON(statusCode, gin.H{
			"status":  "error",
			"message": responseMsg,
			"error":   errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Students synced successfully",
		"data": gin.H{
			"count": count,
		},
	})
}

func (h *StudentHandler) AssignStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var body struct {
		OrganizationID        int    `json:"organization_id"`
		OrganizationShortName string `json:"organization_shortname"`
		PositionTitle         string `json:"position_title"`
		Role                  string `json:"role"`
		Category              string `json:"category"`
		Period                string `json:"period"` // period ikut ditampung
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var result interface{}

	switch strings.ToLower(body.Category) {
	case "bem":
		bem, err := h.service.AssignToBem(uint(id), body.Role, body.PositionTitle, body.Period)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result = bem

	case "mpm":
		mpm, err := h.service.AssignToMpm(uint(id), body.Role, body.PositionTitle, body.Period)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result = mpm

	default:
		period, err := h.service.AssignToPeriod(uint(id), body.OrganizationID, body.Role, body.Period)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result = period
	}
	c.JSON(http.StatusOK, result)
}
