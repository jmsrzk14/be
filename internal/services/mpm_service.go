package services

import (
	"gorm.io/gorm"
	"errors"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

// mpmService is a service for mpm operations
type MpmService struct {
	repository *repositories.MpmRepository
	db *gorm.DB
}

// NewmpmService creates a new mpm service
func NewMpmService(db *gorm.DB) *MpmService {
    return &MpmService{
        repository: repositories.NewMpmRepository(),
    }
}

// Creatempm creates a new mpm
func (s *MpmService) CreateMpm(mpm *models.MPM) error {
	// Check if code exists (including soft-deleted)
	// exists, err := s.repository.CheckNameExists(mpm.Name, 0)
	// if err != nil {
	// 	return err
	// }

	// if exists {
	// 	// Try to find a soft-deleted mpm with this code
	// 	deletedmpm, err := s.repository.FindDeletedByName(mpm.Name)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if deletedmpm != nil {
	// 		// Restore the soft-deleted mpm with updated data
	// 		deletedmpm.Name = mpm.Name
			
	// 		// Restore the mpm
	// 		restoredmpm, err := s.repository.RestoreByName(mpm.Name)
	// 		if err != nil {
	// 			return err
	// 		}
			
	// 		// Update with new data
	// 		restoredmpm.Name = mpm.Name
			
	// 		return s.repository.Update(restoredmpm)
	// 	}
		
	// 	return errors.New("kode gedung sudah digunakan")
	// }

	// Create mpm
	return s.repository.Create(mpm)
}

// Updatempm updates an existing mpm
func (s *MpmService) UpdateMpm(mpm *models.MPM) error {
	// Check if mpm exists
	existingMpm, err := s.repository.FindByID(mpm.ID)
	if err != nil {
		return err
	}
	if existingMpm == nil {
		return errors.New("mpm tidak ditemukan")
	}

	// Update mpm
	return s.repository.Update(mpm)
}

// GetmpmByID gets a mpm by ID
func (s *MpmService) GetMpmByID(id uint) (*models.MPM, error) {
	return s.repository.FindByID(id)
}

// GetAllmpms gets all mpms
func (s *MpmService) GetAllMpms(limit, offset int) ([]models.MPM, int64, error) {
    return s.repository.GetAllMpms(limit, offset)
}

// Deletempm deletes a mpm
func (s *MpmService) DeleteMpm(id uint) error {
	// Check if mpm exists
	mpm, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if mpm == nil {
		return errors.New("gedung tidak ditemukan")
	}

	// Delete mpm (soft delete)
	return s.repository.DeleteByID(id)
}

// mpmWithStats represents a mpm with additional statistics
type MpmWithStats struct {
	Mpm  models.MPM `json:"mpm"`
	RoomCount int64           `json:"room_count"`
}

// GetmpmWithStats gets a mpm with its statistics
func (s *MpmService) GetMpmWithStats(id uint) (*MpmWithStats, error) {
	// Get mpm
	mpm, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if mpm == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}

	// Return mpm with stats
	return &MpmWithStats{
		Mpm:  *mpm,
	}, nil
}

// GetAllmpmsWithStats gets all mpms with their statistics
// func (s *mpmService) GetAllmpmsWithStats() ([]mpmWithStats, error) {
// 	// Get all mpms
// 	mpms, err := s.repository.Get()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build response with stats
// 	result := make([]mpmWithStats, len(mpms))
// 	for i, mpm := range mpms {
		
// 		result[i] = mpmWithStats{
// 			mpm:  mpm,
// 		}
// 	}

// 	return result, nil
// } 

func (s *MpmService) GetMPMByPeriod(period string) (*models.MPM, error) {
	return s.repository.GetMPMByPeriod(period)
}