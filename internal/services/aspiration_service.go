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
	return s.repository.FindByID(id)
}

func (s *AspirationService) GetAllAspirations(limit, offset int) ([]models.Aspiration, int64, error) {
	return s.repository.GetAllAspirations(limit, offset)
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
	aspiration, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if aspiration == nil {
		return nil, errors.New("gambar tidak ditemukan")
	}
	return &AspirationWithStats{
		Aspiration: *aspiration,
	}, nil
}
