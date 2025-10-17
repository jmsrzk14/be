package auth

import (
	"encoding/json"
	"errors"
	"log"
	"io"
	"strings"
	"net/http"
	"net/url"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

const (
	CampusAuthURL = "https://service-users.del.ac.id/api/v1/auth/login"
)

var (
	// ErrCampusAuthFailed is returned when campus authentication fails
	ErrCampusAuthFailed = errors.New("campus authentication failed")
)

// CampusLogin handles authentication with the campus API for all roles
// This includes students, lecturers, and employees
func CampusLogin(system_id, platform, info, username, password string) (*models.CampusLoginResponse, error) {
	log.Printf("Attempting campus login for username: %s", username)

	form := url.Values{}
	form.Set("system_id", system_id)
	form.Set("platform", platform)
	form.Set("info", info)
	form.Set("username", username)
	form.Set("password", password)

	request, err := http.NewRequest("POST", CampusAuthURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	log.Printf("Received response status: %d", response.StatusCode)

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		log.Printf("Campus API returned status %d: %s", response.StatusCode, string(bodyBytes))
		return nil, ErrCampusAuthFailed
	}

	var loginResponse models.CampusLoginResponse
	if err := json.NewDecoder(response.Body).Decode(&loginResponse); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, err
	}

	if err := SaveCampusUserToDatabase(&loginResponse, password); err != nil {
		log.Printf("Error saving user to database: %v", err)
		return nil, err
	}

	return &loginResponse, nil
}

// SaveCampusUserToDatabase creates or updates a user record for a campus user
func SaveCampusUserToDatabase(campusResponse *models.CampusLoginResponse, password string) error {
	// Initialize user repository if needed
	if UserRepository == nil {
		log.Printf("Initializing UserRepository")
		UserRepository = repositories.NewUserRepository()
	}

	log.Printf("User created successfully")
	return nil
}

// ConvertCampusResponseToLoginResponse converts campus login response to standard login response
func ConvertCampusResponseToLoginResponse(campusResponse *models.CampusLoginResponse) *models.LoginResponse {
	// Initialize user repository if needed
	if UserRepository == nil {
		log.Printf("Initializing UserRepository")
		UserRepository = repositories.NewUserRepository()
	}

	// Return login response - ensure token and refreshToken are correctly set
	// This is critical for frontend compatibility
	return &models.LoginResponse{
		Token:        campusResponse.Data.Token,        // JWT token for authorization
	}
}