package services

import (
	"bem_be/internal/models"
	"bem_be/internal/repositories"
	"errors"
	// "strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	Request models.Request `json:"request"`
}

func (s *RequestService) GetRequestWithStats(id uint) (*RequestWithStats, error) {
	request, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, errors.New("permintaan tidak ditemukan")
	}
	return &RequestWithStats{
		Request: *request,
	}, nil
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
func (s *RequestService) ProcessRequestStatus(requestID uint, newStatus string, adminID int) (*models.Request, error) {
	var finalRequest *models.Request

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Ambil & Kunci request yang akan diupdate
		var request models.Request
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&request, requestID).Error; err != nil {
			return errors.New("request not found")
		}

		// 2. Validasi status
		if request.Status != "pending" {
			return errors.New("request has already been processed")
		}

		// 3. Logika jika disetujui
		if newStatus == "approved" {
			// Panggil ItemRepository yang ada.
			// Asumsi: request.Name SAMA DENGAN item.Code
			// Kita harus query item di dalam transaksi untuk mengunci barisnya.
			// Jadi, kita tidak bisa langsung pakai s.itemRepo.FindByName(request.Name)
			// Kita akan query langsung ke tabelnya.
			var item models.Item
			if err := tx.Table("item").Where("code = ?", request.Name).Clauses(clause.Locking{Strength: "UPDATE"}).First(&item).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return errors.New("master item with code '" + request.Name + "' not found")
				}
				return err
			}
		}

		// 4. Update status request
		request.Status = newStatus
		// request.ApproverID = adminID
		request.UpdatedAt = time.Now()

		if err := tx.Save(&request).Error; err != nil {
			return err
		}

		finalRequest = &request
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repository.FindByID(finalRequest.ID)
}
