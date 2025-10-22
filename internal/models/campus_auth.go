package models

import (
	"encoding/json"
	"log"
)

// CampusPosition represents a position/role in the campus system
type CampusPosition struct {
	StrukturJabatanID int    `json:"struktur_jabatan_id"`
	Jabatan           string `json:"jabatan"`
}

// CampusUser represents a user in the campus system
type CampusUser struct {
	UserID   int             `json:"user_id"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Role     string          `json:"role"`
	Status   int             `json:"status"`
	Jabatan  json.RawMessage `json:"jabatan"` // Using RawMessage to handle different types
}

// GetJabatanPositions returns the jabatan as []CampusPosition if it's an array,
// or an empty array if it's a string or any other type
func (u *CampusUser) GetJabatanPositions() []CampusPosition {
	var positions []CampusPosition
	err := json.Unmarshal(u.Jabatan, &positions)
	if err != nil {
		// If not an array, it might be a string or something else - just return empty array
		log.Printf("Jabatan is not an array of positions: %v", err)
		return []CampusPosition{}
	}
	return positions
}

// GetJabatanString returns the jabatan as string if it's a string,
// or empty string if it's an array or any other type
func (u *CampusUser) GetJabatanString() string {
	var jabatanString string
	err := json.Unmarshal(u.Jabatan, &jabatanString)
	if err != nil {
		// If not a string, it might be an array or something else - just return empty string
		log.Printf("Jabatan is not a string: %v", err)
		return ""
	}
	return jabatanString
}

// CampusLoginRequest represents the campus login request
type CampusLoginRequest struct {
	SystemID string `form:"system_id" json:"system_id" binding:"required"`
	Platform string `form:"platform" json:"platform" binding:"required"`
	Info     string `form:"info" json:"info" binding:"required"`
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// CampusLoginResponse represents the response from the campus login API
type CampusLoginResponse struct {
	Status  string `json:"status"`
	Token   string `json:"token"`
	Message string `json:"message"`
	Data    struct {
		Token    string `json:"token"`
		Username string `json:"username"`
		UserID   string `json:"user_id"`
		Role     string `json:"role"`
		Email    string `json:"email"`
		Status   int    `json:"status"`
	} `json:"data"`
}

type TOTPData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	QRCode  string `json:"qrCode,omitempty"`
	Secret  string `json:"secret,omitempty"`
}
