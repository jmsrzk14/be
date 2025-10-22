package services

import (
	"bem_be/internal/models"
	"bem_be/internal/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func (s *RequestService) GetStudentByUserIDSarpras(username string) (*models.Student, error) {
	return s.studentRepo.FindByUserID(username)
}

func (s *RequestService) CreateRequestSarpras(request *models.Request) error {
	return s.repository.CreateSarpras(request)
}

func (s *RequestService) UpdateRequestSarpras(request *models.Request) error {
	existingRequest, err := s.repository.FindByIDSarpras(request.ID)
	if err != nil {
		return err
	}
	if existingRequest == nil {
		return errors.New("request not found")
	}
	return s.repository.UpdateSarpras(request)
}

func (s *RequestService) GetRequestByIDSarpras(id uint) (*models.Request, error) {
	return s.repository.FindByIDSarpras(id)
}

func (s *RequestService) GetAllRequestsSarpras(category, limit, offset int, search string) ([]models.Request, int64, error) {
	return s.repository.GetAllRequestsSarpras(category, limit, offset, search)
}

type RequestWithStatsSarpras struct {
	Request   models.Request `json:"request"`
	ItemNames []string       `json:"item_names"`
}

func (s *RequestService) GetRequestWithStatsSarpras(id uint) (*models.RequestWithStats, error) {
	// 1. Ambil data request
	request, err := s.repository.FindByIDSarpras(id)
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
	items, err := s.repository.FindItemsByIDsSarpras(itemIDs)
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

func (s *RequestService) GetRequestsByRequesterIDSarpras(requesterID string) ([]models.Request, error) {
	requests, err := s.repository.FindAllByRequesterIDSarpras(requesterID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data request: %v", err)
	}

	if len(requests) == 0 {
		return nil, errors.New("belum ada request untuk user ini")
	}

	return requests, nil
}

func (s *RequestService) DeleteRequestSarpras(id uint) error {
	request, err := s.repository.FindByIDSarpras(id)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("request not found")
	}
	return s.repository.DeleteByIDSarpras(id)
}

// Approve or Reject Request
func (s *RequestService) ProcessRequestStatusSarpras(requestID uint, status string, adminID string, reason string) (*models.Request, error) {
	// Cari request
	request, err := s.repository.FindByIDSarpras(requestID)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, errors.New("Request tidak ditemukan")
	}

	// Update status dan admin ID
	request.Status = status
	request.ApproverID = adminID

	// Simpan alasan penolakan jika ada
	if status == "rejected" {
		request.Reason = reason
	}

	// Simpan perubahan
	if err := s.repository.UpdateSarpras(request); err != nil {
		return nil, err
	}

	return request, nil
}

func (s *RequestService) UpdateImageBarangAndStatusSarpras(id uint, filename, status string) error {
	// 1. Ambil data request
	request, err := s.repository.FindByIDSarpras(id)
	if err != nil {
		return err
	}
	if request == nil {
		return fmt.Errorf("request not found")
	}

	// 2. Decode item IDs dari request.Item
	var itemIDs []uint
	if err := json.Unmarshal([]byte(request.Item), &itemIDs); err != nil {
		// Jika gagal decode JSON, fallback ke format "1,2,3"
		parts := strings.Split(request.Item, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			if num, convErr := strconv.Atoi(p); convErr == nil {
				itemIDs = append(itemIDs, uint(num))
			}
		}
	}

	// 3. Ambil semua item berdasarkan ID yang ditemukan
	items, err := s.repository.FindItemsByIDsSarpras(itemIDs)
	if err != nil {
		return err
	}

	// 4. Kurangi stok setiap item
	for _, item := range items {
		if item.Amount > 0 {
			item.Amount -= 1
			if err := s.itemRepo.Update(&item); err != nil {
				return fmt.Errorf("gagal update stok untuk item %s: %v", item.Name, err)
			}
		}
	}

	// 5. Update request (gambar + status)
	request.ImageURLBRG = filename
	request.Status = status

	return s.repository.UpdateSarpras(request)
}

func (s *RequestService) UpdateStatusRequestSarpras(id uint, status string) error {
	return s.repository.UpdateStatusSarpras(id, status)
}

func (s *RequestService) ReturnedItemSarpras(id uint, returnedAt time.Time) error {
	return s.repository.UpdateStatusAndReturnTimeSarpras(id, "dikembalikan", returnedAt)
}

func (s *RequestService) UpdateItemStockOnTakenSarpras(id uint) error {
	// 1. Ambil data request
	request, err := s.repository.FindByIDSarpras(id)
	if err != nil {
		return err
	}
	if request == nil {
		return fmt.Errorf("request not found")
	}

	// 2. Decode item IDs dari kolom request.Item
	var itemIDs []uint
	if err := json.Unmarshal([]byte(request.Item), &itemIDs); err != nil {
		// Fallback: handle format "1,2,3"
		parts := strings.Split(request.Item, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			if num, convErr := strconv.Atoi(p); convErr == nil {
				itemIDs = append(itemIDs, uint(num))
			}
		}
	}

	// 3. Ambil semua item berdasarkan ID
	items, err := s.repository.FindItemsByIDsSarpras(itemIDs)
	if err != nil {
		return fmt.Errorf("failed to get items: %v", err)
	}

	// 4. Kurangi stok setiap item (1 per item)
	for _, item := range items {
		if item.Amount > 0 {
			item.Amount += 1
			if err := s.itemRepo.Update(&item); err != nil {
				return fmt.Errorf("failed to update stock for item %s: %v", item.Name, err)
			}
		}
	}

	// 5. Tidak ubah status, tidak ubah image
	return nil
}

func (s *RequestService) GetStudentByUserIDDepol(username string) (*models.Student, error) {
	return s.studentRepo.FindByUserID(username)
}

func (s *RequestService) CreateRequestDepol(request *models.Request) error {
	return s.repository.CreateDepol(request)
}

func (s *RequestService) UpdateRequestDepol(request *models.Request) error {
	existingRequest, err := s.repository.FindByIDDepol(request.ID)
	if err != nil {
		return err
	}
	if existingRequest == nil {
		return errors.New("request not found")
	}
	return s.repository.UpdateDepol(request)
}

func (s *RequestService) GetRequestByIDDepol(id uint) (*models.Request, error) {
	return s.repository.FindByIDDepol(id)
}

func (s *RequestService) GetAllRequestsDepol(category, limit, offset int, search string) ([]models.Request, int64, error) {
	return s.repository.GetAllRequestsDepol(category, limit, offset, search)
}

type RequestWithStatsDepol struct {
	Request   models.Request `json:"request"`
	ItemNames []string       `json:"item_names"`
}

func (s *RequestService) GetRequestWithStatsDepol(id uint) (*models.RequestWithStats, error) {
	// 1. Ambil data request
	request, err := s.repository.FindByIDDepol(id)
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
	items, err := s.repository.FindItemsByIDsDepol(itemIDs)
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

func (s *RequestService) GetRequestsByRequesterIDDepol(username string) ([]models.Request, error) {
	requests, err := s.repository.FindAllByRequesterIDDepol(username)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data request: %v", err)
	}

	if len(requests) == 0 {
		return nil, errors.New("belum ada request untuk user ini")
	}

	return requests, nil
}

func (s *RequestService) DeleteRequestDepol(id uint) error {
	request, err := s.repository.FindByIDDepol(id)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("request not found")
	}
	return s.repository.DeleteByIDDepol(id)
}

// Approve or Reject Request
func (s *RequestService) ProcessRequestStatusDepol(requestID uint, status string, adminID string, reason string) (*models.Request, error) {
	// Cari request
	request, err := s.repository.FindByIDDepol(requestID)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, errors.New("Request tidak ditemukan")
	}

	// Update status dan admin ID
	request.Status = status
	request.ApproverID = adminID

	// Simpan alasan penolakan jika ada
	if status == "rejected" {
		request.Reason = reason
	}

	// Simpan perubahan
	if err := s.repository.UpdateDepol(request); err != nil {
		return nil, err
	}

	return request, nil
}

func (s *RequestService) UpdateImageBarangAndStatusDepol(id uint, filename, status string) error {
	// 1. Ambil data request
	request, err := s.repository.FindByIDDepol(id)
	if err != nil {
		return err
	}
	if request == nil {
		return fmt.Errorf("request not found")
	}

	// 2. Decode item IDs dari request.Item
	var itemIDs []uint
	if err := json.Unmarshal([]byte(request.Item), &itemIDs); err != nil {
		// Jika gagal decode JSON, fallback ke format "1,2,3"
		parts := strings.Split(request.Item, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			if num, convErr := strconv.Atoi(p); convErr == nil {
				itemIDs = append(itemIDs, uint(num))
			}
		}
	}

	// 3. Ambil semua item berdasarkan ID yang ditemukan
	items, err := s.repository.FindItemsByIDsDepol(itemIDs)
	if err != nil {
		return err
	}

	// 4. Kurangi stok setiap item
	for _, item := range items {
		if item.Amount > 0 {
			item.Amount -= 1
			if err := s.itemRepo.Update(&item); err != nil {
				return fmt.Errorf("gagal update stok untuk item %s: %v", item.Name, err)
			}
		}
	}

	// 5. Update request (gambar + status)
	request.ImageURLBRG = filename
	request.Status = status

	return s.repository.UpdateDepol(request)
}

func (s *RequestService) UpdateStatusRequestDepol(id uint, status string) error {
	return s.repository.UpdateStatusDepol(id, status)
}

func (s *RequestService) ReturnedItemDepol(id uint, returnedAt time.Time) error {
	return s.repository.UpdateStatusAndReturnTimeDepol(id, "dikembalikan", returnedAt)
}

func (s *RequestService) UpdateItemStockOnTakenDepol(id uint) error {
	// 1. Ambil data request
	request, err := s.repository.FindByIDDepol(id)
	if err != nil {
		return err
	}
	if request == nil {
		return fmt.Errorf("request not found")
	}

	// 2. Decode item IDs dari kolom request.Item
	var itemIDs []uint
	if err := json.Unmarshal([]byte(request.Item), &itemIDs); err != nil {
		// Fallback: handle format "1,2,3"
		parts := strings.Split(request.Item, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			if num, convErr := strconv.Atoi(p); convErr == nil {
				itemIDs = append(itemIDs, uint(num))
			}
		}
	}

	// 3. Ambil semua item berdasarkan ID
	items, err := s.repository.FindItemsByIDsDepol(itemIDs)
	if err != nil {
		return fmt.Errorf("failed to get items: %v", err)
	}

	// 4. Kurangi stok setiap item (1 per item)
	for _, item := range items {
		if item.Amount > 0 {
			item.Amount += 1
			if err := s.itemRepo.Update(&item); err != nil {
				return fmt.Errorf("failed to update stock for item %s: %v", item.Name, err)
			}
		}
	}

	// 5. Tidak ubah status, tidak ubah image
	return nil
}
