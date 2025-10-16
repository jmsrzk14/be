package handlers

import (
	"errors"
	"log"
	"net/http"

	"bem_be/internal/auth" // ✅ IMPORT AUTH PACKAGE
	"bem_be/internal/models"
	"bem_be/internal/repositories"

	"github.com/gin-gonic/gin"
)

// CampusLogin handles login requests for campus users (all roles)
func CampusLogin(c *gin.Context) {
	var req models.CampusLoginRequestV2

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	log.Printf("Login attempt - Device: %s, Platform: %s, User: %s", req.SystemID, req.Platform, req.Username)

	// Validasi
	if req.SystemID == "" || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "system_id, username, password required",
		})
		return
	}

	campusResponse, err := auth.CampusLoginV2(req.SystemID, req.Platform, req.Info, req.Username, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "Authentication failed"
		if errors.Is(err, auth.ErrCampusAuthFailed) {
			status = http.StatusUnauthorized
			msg = "Invalid credentials"
		}
		log.Printf("Login failed: %v", err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	loginResponse := auth.ConvertCampusResponseToLoginResponse(campusResponse)

	studentRepo := repositories.NewStudentRepository()
	student, err := studentRepo.FindByExternalUserUsername(loginResponse.User.Username)
	if err != nil {
		log.Printf("Student fetch error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile"})
		return
	}

	orgID := 0
	if student.OrganizationID != nil {
		orgID = *student.OrganizationID
	}

	response := models.OrderedLoginResponse{
		User:           loginResponse.User,
		Token:          loginResponse.Token,
		RefreshToken:   loginResponse.RefreshToken,
		Position:       student.Position,
		OrganizationID: orgID,
		SystemID:       req.SystemID,
		Platform:       req.Platform,
	}

	log.Printf("✅ Login SUCCESS - %s (Role: %s, Pos: %s)", 
		req.Username, loginResponse.User.Role, student.Position)

	c.JSON(http.StatusOK, response)
}