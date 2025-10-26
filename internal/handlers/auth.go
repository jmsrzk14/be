package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"bem_be/internal/auth"
	"bem_be/internal/database"
	"bem_be/internal/models"

	"github.com/gin-gonic/gin"
)

// Login handles the login request
func Login(c *gin.Context) {
	var req models.LoginRequest

	// Validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Attempt to login
	response, err := auth.Login(req.Username, req.Password)
	if err != nil {
		var statusCode int
		var message string

		// Handle different error types
		switch {
		case errors.Is(err, auth.ErrUserNotFound), errors.Is(err, auth.ErrInvalidCredentials):
			statusCode = http.StatusUnauthorized
			message = "Invalid username or password"
		default:
			statusCode = http.StatusInternalServerError
			message = "An error occurred during login"
		}

		c.JSON(statusCode, gin.H{"error": message})
		return
	}

	// Use custom response struct to ensure the correct field order
	orderedResponse := models.OrderedLoginResponse{
		User:         response.User,
		Token:        response.Token,
		RefreshToken: response.RefreshToken,
	}

	// Set content type
	c.Header("Content-Type", "application/json")

	// Manually marshal to JSON to ensure field order
	jsonBytes, err := json.Marshal(orderedResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating response"})
		return
	}

	// Write the response
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(jsonBytes)
}

// RefreshToken handles token refresh requests
func RefreshToken(c *gin.Context) {
	var req models.RefreshRequest

	// Validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Attempt to refresh the token
	response, err := auth.RefreshToken(req.RefreshToken)
	if err != nil {
		var statusCode int
		var message string

		// Handle different error types
		switch {
		case errors.Is(err, auth.ErrInvalidToken):
			statusCode = http.StatusUnauthorized
			message = "Invalid or expired refresh token"
		case errors.Is(err, auth.ErrUserNotFound):
			statusCode = http.StatusUnauthorized
			message = "User not found"
		default:
			statusCode = http.StatusInternalServerError
			message = "An error occurred during token refresh"
		}

		c.JSON(statusCode, gin.H{"error": message})
		return
	}

	// Use custom response struct to ensure the correct field order
	orderedResponse := models.OrderedLoginResponse{
		User:         response.User,
		Token:        response.Token,
		RefreshToken: response.RefreshToken,
	}

	// Set content type
	c.Header("Content-Type", "application/json")

	// Manually marshal to JSON to ensure field order
	jsonBytes, err := json.Marshal(orderedResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating response"})
		return
	}

	// Write the response
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(jsonBytes)
}

// GetCurrentUser returns the currently logged-in user
func GetCurrentUser(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	if request.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// 🔹 Query student berdasarkan username
	var student models.Student
	if err := database.DB.Where("user_name = ?", request.Username).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// 🔹 Ambil data organisasi jika ada
	var organization models.Organization
	if student.OrganizationID != nil && *student.OrganizationID != 0 {
		if err := database.DB.Where("id = ?", student.OrganizationID).First(&organization).Error; err != nil {
			organization = models.Organization{}
		}
	}

	// 🔹 Kirim response JSON
	c.JSON(http.StatusOK, gin.H{
		"id":            student.ID,
		"name":          student.FullName,
		"email":         student.Email,
		"username":      student.UserName,
		"nim":           student.NIM,
		"study_program": student.StudyProgram,
		"image":         student.Image,
		"linkedin":      student.LinkedIn,
		"instagram":     student.Instagram,
		"whatsapp":      student.WhatsApp,
		"faculty":       student.Faculty,
		"year_enrolled": student.YearEnrolled,
		"status":        student.Status,
		"position":      student.Position,
		"organization": gin.H{
			"id":   organization.ID,
			"name": organization.Name,
		},
	})
}

func EditProfile(c *gin.Context) {
	// Ambil username dari JSON body
	var request struct {
		Username  string `form:"username"`
		LinkedIn  string `form:"linkedin"`
		Instagram string `form:"instagram"`
		WhatsApp  string `form:"whatsapp"`
	}

	// Validasi input JSON
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if request.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// Cari student berdasarkan username
	var student models.Student
	if err := database.DB.Where("user_name = ?", request.Username).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Ambil file dari form-data (opsional)
	file, err := c.FormFile("image")
	var imageName string

	if err == nil {
		imageName = fmt.Sprintf("%d_%s", student.ID, file.Filename)
		savePath := fmt.Sprintf("uploads/user/%s", imageName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
	}

	// Update data student
	if imageName != "" {
		student.Image = imageName
	}
	student.LinkedIn = request.LinkedIn
	student.Instagram = request.Instagram
	student.WhatsApp = request.WhatsApp

	if err := database.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Buat URL lengkap untuk gambar (biar frontend bisa akses)
	imageURL := ""
	if student.Image != "" {
		imageURL = fmt.Sprintf("http://localhost:9090/uploads/user/%s", student.Image)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Profile updated successfully",
		"image":     student.Image,
		"image_url": imageURL,
		"linkedin":  student.LinkedIn,
		"instagram": student.Instagram,
		"whatsapp":  student.WhatsApp,
	})
}
