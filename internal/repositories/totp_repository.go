package repositories

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	totpSetupURL  = "https://service-users.del.ac.id/api/v1/auth/totp/setup"
	totpVerifyURL = " "
)

// Struct untuk response setup
type TOTPSetupResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		QRCode string `json:"qrcode,omitempty"`
		Secret string `json:"secret,omitempty"`
	} `json:"data,omitempty"`
}

// Struct untuk response verifikasi
type TOTPVerifyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Repository
type TOTPRepository struct{}

func NewTOTPRepository() *TOTPRepository {
	return &TOTPRepository{}
}

// Fungsi utama: setup + auto verify jika sudah punya TOTP
func (r *TOTPRepository) GetOrVerifyTOTP(token string) (*TOTPSetupResponse, error) {
	// Kirim POST ke endpoint setup
	req, err := http.NewRequest("POST", totpSetupURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call TOTP setup: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var setupResp TOTPSetupResponse
	if err := json.Unmarshal(body, &setupResp); err != nil {
		return nil, fmt.Errorf("failed to parse TOTP setup response: %w", err)
	}

	// Jika user sudah punya TOTP → frontend nanti yang handle input kode verifikasi
	if setupResp.Status == "fail" && setupResp.Message == "Anda sudah memiliki TOTP." {
		fmt.Println("User sudah memiliki TOTP. Frontend harus meminta kode verifikasi dari user.")
		return &setupResp, nil
	}

	// Jika response bukan 200 OK, lempar error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TOTP setup failed with status code: %d, message: %s", resp.StatusCode, setupResp.Message)
	}

	return &setupResp, nil
}

// Fungsi untuk verifikasi TOTP
func (r *TOTPRepository) PostTOTPVerify(token, code string) (*TOTPVerifyResponse, error) {
	// Payload dikirim dari user (kode dari Google Authenticator)
	payload := map[string]string{
		"code": code,
	}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", totpVerifyURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call TOTP verify: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var verifyResp TOTPVerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return nil, fmt.Errorf("failed to parse TOTP verify response: %w", err)
	}

	// Jika status bukan 200 OK, kirim pesan error dari server kampus
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(verifyResp.Message)
	}

	return &verifyResp, nil
}
