package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"bem_be/internal/auth"
	"bem_be/internal/models"
	"bem_be/internal/repositories"

	"github.com/gin-gonic/gin"
)

func CampusLogin(c *gin.Context) {
	var req models.CampusLoginRequest

	// Bind form data
	if err := c.ShouldBind(&req); err != nil {
		log.Printf("Error binding request data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	log.Printf("Campus login attempt for username: %s", req.Username)

	// Call campus login service
	campusResponse, err := auth.CampusLogin(req.SystemID, req.Platform, req.Info, req.Username, req.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Authentication failed"

		if errors.Is(err, auth.ErrCampusAuthFailed) {
			statusCode = http.StatusUnauthorized
			message = "Campus authentication failed"
		}

		log.Printf("Campus login failed: %v", err)
		c.JSON(statusCode, gin.H{"error": message})
		return
	}

	// Convert to standard login response
	loginResponse := auth.ConvertCampusResponseToLoginResponse(campusResponse)

	log.Printf("Login success for user: %s (role: %s)", loginResponse.User.Username, loginResponse.User.Role)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login berhasil",
		"data": gin.H{
			"token": loginResponse.Token,
		},
	})
}

func TOTPSetup(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	totpRepo := repositories.NewTOTPRepository()

	totpSetupResp, err := totpRepo.GetOrVerifyTOTP(token[7:])
	if err != nil {
		log.Printf("Failed to setup TOTP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("TOTP setup failed: %v", err),
		})
		return
	}

	log.Printf("TOTP setup response: %s - %s", totpSetupResp.Status, totpSetupResp.Message)

	c.JSON(http.StatusOK, gin.H{
		"status":  totpSetupResp.Status,
		"message": totpSetupResp.Message,
		"data": gin.H{
			"qrcode": totpSetupResp.Data.QRCode,
			"secret": totpSetupResp.Data.Secret,
		},
	})
}

func TOTPVerify(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid TOTP code"})
		return
	}

	totpRepo := repositories.NewTOTPRepository()
	verifyResp, err := totpRepo.PostTOTPVerify(token[7:], req.Code)
	if err != nil {
		log.Printf("TOTP verify failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  verifyResp.Status,
		"message": verifyResp.Message,
	})
}
