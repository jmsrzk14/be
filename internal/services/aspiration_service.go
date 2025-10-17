package services

import (
	"errors"

	"gorm.io/gorm"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

type AspirationService struct {
	repository *repositories.AspirationRepository
	db         *gorm.DB
}

func NewAspirationService(db *gorm.DB) *AspirationService {
	return &AspirationService{
		repository: repositories.NewAspirationRepository(),
		db:         db, // ✅ tambahkan ini biar Preload jalan
	}
}

func (s *AspirationService) CreateAspiration(aspiration *models.Aspiration) error {
	return s.repository.Create(aspiration)
}

func (s *AspirationService) UpdateAspiration(aspiration *models.Aspiration) error {
	existingAspiration, err := s.repository.FindByID(aspiration.ID)
	if err != nil {
		return err
	}
	if existingAspiration == nil {
		return errors.New("Aspirasi tidak ditemukan")
	}
	return s.repository.Update(aspiration)
}

func (s *AspirationService) GetAspirationByID(id uint) (*models.Aspiration, error) {
	var aspiration models.Aspiration
	if err := s.db.Preload("Student").First(&aspiration, id).Error; err != nil {
		return nil, err
	}
	return &aspiration, nil
}

func (s *AspirationService) GetAllAspirations(limit, offset int) ([]models.Aspiration, int64, error) {
	var aspirations []models.Aspiration
	var total int64

	// ✅ Hitung total aspirasi
	if err := s.db.Model(&models.Aspiration{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ✅ Ambil aspirasi + preload student (join ke tabel students)
	if err := s.db.Preload("Student").Limit(limit).Offset(offset).Find(&aspirations).Error; err != nil {
		return nil, 0, err
	}

	// ✅ Kalau ada aspirasi yang belum punya student, hindari null pointer
	for i := range aspirations {
		if aspirations[i].Student.FullName == "" {
			aspirations[i].Student.FullName = "-"
		}
	}

	return aspirations, total, nil
}

func (s *AspirationService) DeleteAspiration(id uint) error {
	aspiration, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if aspiration == nil {
		return errors.New("aspirasi tidak ditemukan")
	}
	return s.repository.DeleteByID(id)
}

type AspirationWithStats struct {
	Aspiration models.Aspiration `json:"aspiration"`
	RoomCount  int64             `json:"room_count"`
}

func (s *AspirationService) GetAspirationWithStats(id uint) (*AspirationWithStats, error) {
	aspiration, err := s.GetAspirationByID(id)
	if err != nil {
		return nil, err
	}

	return &AspirationWithStats{
		Aspiration: *aspiration,
	}, nil
}
