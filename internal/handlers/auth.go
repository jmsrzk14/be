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
	// Get the user ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in token"})
		return
	}

	// Convert userID ke uint biar bisa query DB
	var userIDValue uint
	switch v := userID.(type) {
	case float64:
		userIDValue = uint(v)
	case int:
		userIDValue = uint(v)
	case uint:
		userIDValue = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID type"})
		return
	}

	// Query ke model Student
	var student models.Student
	if err := database.DB.Where("user_id = ?", userIDValue).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	var organization models.Organization
	if student.OrganizationID != nil && *student.OrganizationID != 0 { // pastikan organizationId ada
		if err := database.DB.Where("id = ?", student.OrganizationID).First(&organization).Error; err != nil {
			organization = models.Organization{}
		}
	}

	// Return data student + role
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"id":            student.ID,
		"name":          student.FullName,
		"email":         student.Email,
		"username":      student.UserName,
		"nim":           student.NIM,
		"study_program": student.StudyProgram,
		"image":         student.Image,
		"role":          role,
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
	// Ambil userID dari token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in token"})
		return
	}

	var userIDValue uint
	switch v := userID.(type) {
	case float64:
		userIDValue = uint(v)
	case int:
		userIDValue = uint(v)
	case uint:
		userIDValue = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID type"})
		return
	}

	// Ambil field text dari form-data
	linkedin := c.PostForm("linkedin")
	instagram := c.PostForm("instagram")
	whatsapp := c.PostForm("whatsapp")

	// Ambil file dari form-data
	file, err := c.FormFile("image")
	var imageName string

	if err == nil {
		// Nama file untuk disimpan di DB
		imageName = fmt.Sprintf("%d_%s", userIDValue, file.Filename)

		// Path lengkap untuk simpan di folder
		savePath := fmt.Sprintf("uploads/user/%s", imageName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
	}

	// Cari student
	var student models.Student
	if err := database.DB.Where("user_id = ?", userIDValue).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Update data
	if imageName != "" {
		student.Image = imageName // simpan hanya nama file
	}
	student.LinkedIn = linkedin
	student.Instagram = instagram
	student.WhatsApp = whatsapp

	if err := database.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Biar frontend gampang akses gambar, balikin URL lengkap
	imageURL := ""
	if student.Image != "" {
		imageURL = fmt.Sprintf("http://localhost:9090/uploads/user/%s", student.Image)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Profile updated successfully",
		"image":     student.Image, // nama file saja
		"image_url": imageURL,      // URL lengkap untuk akses
		"linkedin":  student.LinkedIn,
		"instagram": student.Instagram,
		"whatsapp":  student.WhatsApp,
	})
}
