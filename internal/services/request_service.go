package services

import (
	"bem_be/internal/models"
	"bem_be/internal/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type RequestService struct {
	repository  *repositories.RequestRepository
	studentRepo *repositories.StudentRepository
	itemRepo    *repositories.ItemRepository
	db          *gorm.DB
}

func NewRequestService(db *gorm.DB) *RequestService {
	return &RequestService{
		repository:  repositories.NewRequestRepository(),
		studentRepo: repositories.NewStudentRepository(),
		itemRepo:    repositories.NewItemRepository(),
		db:          db,
	}
}

func (s *RequestService) GetStudentByUserID(userID int) (*models.Student, error) {
	return s.studentRepo.FindByUserID(userID)
}

func (s *RequestService) CreateRequest(request *models.Request) error {
	return s.repository.Create(request)
}

func (s *RequestService) UpdateRequest(request *models.Request) error {
	existingRequest, err := s.repository.FindByID(request.ID)
	if err != nil {
		return err
	}
	if existingRequest == nil {
		return errors.New("request not found")
	}
	return s.repository.Update(request)
}

func (s *RequestService) GetRequestByID(id uint) (*models.Request, error) {
	return s.repository.FindByID(id)
}

func (s *RequestService) GetAllRequests(limit, offset int) ([]models.Request, int64, error) {
	return s.repository.GetAllRequests(limit, offset)
}

type RequestWithStats struct {
	Request   models.Request `json:"request"`
	ItemNames []string       `json:"item_names"`
}

func (s *RequestService) GetRequestWithStats(id uint) (*models.RequestWithStats, error) {
	// 1. Ambil data request
	request, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, errors.New("permintaan tidak ditemukan")
	}

	// 2. Ambil item IDs
	var itemIDs []uint

	// Coba decode JSON array dulu
	if err := json.Unmarshal([]byte(request.Item), &itemIDs); err != nil {
		// Kalau gagal, berarti formatnya kemungkinan "1,2"
		strIDs := strings.Split(request.Item, ",")
		for _, strID := range strIDs {
			strID = strings.TrimSpace(strID)
			if strID == "" {
				continue
			}
			idInt, convErr := strconv.Atoi(strID)
			if convErr == nil {
				itemIDs = append(itemIDs, uint(idInt))
			}
		}
	}

	if len(itemIDs) == 0 {
		return &models.RequestWithStats{
			Request:   *request,
			ItemNames: []string{},
		}, nil
	}

	// 3. Ambil data items berdasarkan itemIDs
	items, err := s.repository.FindItemsByIDs(itemIDs)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data item: %v", err)
	}

	// 4. Ambil hanya nama-nama item
	itemNames := make([]string, len(items))
	for i, item := range items {
		itemNames[i] = item.Name
	}

	// 5. Gabungkan hasil
	return &models.RequestWithStats{
		Request:   *request,
		ItemNames: itemNames,
	}, nil
}

func (s *RequestService) GetRequestsByRequesterID(requesterID uint) ([]models.Request, error) {
	requests, err := s.repository.FindAllByRequesterID(requesterID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data request: %v", err)
	}

	if len(requests) == 0 {
		return nil, errors.New("belum ada request untuk user ini")
	}

	return requests, nil
}

func (s *RequestService) DeleteRequest(id uint) error {
	request, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("request not found")
	}
	return s.repository.DeleteByID(id)
}

// Approve or Reject Request
func (s *RequestService) ProcessRequestStatus(requestID uint, status string, adminID int, reason string) (*models.Request, error) {
	// Cari request
	request, err := s.repository.FindByID(requestID)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, errors.New("Request tidak ditemukan")
	}

	// Update status dan admin ID
	request.Status = status
	request.ApproverID = uint(adminID)

	// Simpan alasan penolakan jika ada
	if status == "rejected" {
		request.Reason = reason
	}

	// Simpan perubahan
	if err := s.repository.Update(request); err != nil {
		return nil, err
	}

	return request, nil
}

func (s *RequestService) UpdateImageBarangAndStatus(id uint, fileName string, status string) error {
	return s.repository.UpdateImageBarangAndStatus(id, fileName, status)
}