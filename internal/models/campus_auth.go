package models

import (
	"encoding/json"
	"log"
)

type CampusPosition struct {
	StrukturJabatanID int    `json:"struktur_jabatan_id"`
	Jabatan           string `json:"jabatan"`
}

func (u *CampusUser) GetJabatanPositions() []CampusPosition {
	var positions []CampusPosition
	err := json.Unmarshal(u.Jabatan, &positions)
	if err != nil {
		log.Printf("Jabatan is not an array of positions: %v", err)
		return []CampusPosition{}
	}
	return positions
}

func (u *CampusUser) GetJabatanString() string {
	var jabatanString string
	err := json.Unmarshal(u.Jabatan, &jabatanString)
	if err != nil {
		log.Printf("Jabatan is not a string: %v", err)
		return ""
	}
	return jabatanString
}

type CampusLoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type CampusLoginRequestV2 struct {
	SystemID string `json:"system_id" binding:"required"`
	Platform string `json:"platform" binding:"required"`
	Info     string `json:"info" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CampusLoginResponse struct {
	Result        bool        `json:"result"`
	Token         string      `json:"token"`
	RefreshToken  string      `json:"refresh_token"`
	Error         string      `json:"error,omitempty"`
	User          CampusUser  `json:"user"`  
}

type CampusUser struct {
	UserID   int             `json:"user_id"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Role     string          `json:"role"`
	Status   int             `json:"status"`
	Jabatan  json.RawMessage `json:"jabatan"` 
}

type OrderedLoginResponse struct {
	User           User    `json:"user"`
	Token          string  `json:"token"`
	RefreshToken   string  `json:"refresh_token"`
	Position       string  `json:"position"`
	OrganizationID int     `json:"organization_id"`
	SystemID       string  `json:"system_id"`
	Platform       string  `json:"platform"`
}