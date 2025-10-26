package services

import (
	"gorm.io/gorm"
	"errors"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

// bemService is a service for bem operations
type BemService struct {
	repository *repositories.BemRepository
	db *gorm.DB
}

// NewbemService creates a new bem service
func NewBemService(db *gorm.DB) *BemService {
    return &BemService{
        repository: repositories.NewBemRepository(),
    }
}

// Createbem creates a new bem
func (s *BemService) CreateBem(bem *models.BEM) error {
	// Check if code exists (including soft-deleted)
	// exists, err := s.repository.CheckNameExists(bem.Name, 0)
	// if err != nil {
	// 	return err
	// }

	// if exists {
	// 	// Try to find a soft-deleted bem with this code
	// 	deletedbem, err := s.repository.FindDeletedByName(bem.Name)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if deletedbem != nil {
	// 		// Restore the soft-deleted bem with updated data
	// 		deletedbem.Name = bem.Name
			
	// 		// Restore the bem
	// 		restoredbem, err := s.repository.RestoreByName(bem.Name)
	// 		if err != nil {
	// 			return err
	// 		}
			
	// 		// Update with new data
	// 		restoredbem.Name = bem.Name
			
	// 		return s.repository.Update(restoredbem)
	// 	}
		
	// 	return errors.New("kode gedung sudah digunakan")
	// }

	// Create bem
	return s.repository.Create(bem)
}

// Updatebem updates an existing bem
func (s *BemService) UpdateBem(bem *models.BEM) error {
	// Check if bem exists
	existingBem, err := s.repository.FindByID(bem.ID)
	if err != nil {
		return err
	}
	if existingBem == nil {
		return errors.New("Bem tidak ditemukan")
	}

	// Update bem
	return s.repository.Update(bem)
}

// GetbemByID gets a bem by ID
func (s *BemService) GetBemByID(id uint) (*models.BEM, error) {
	return s.repository.FindByID(id)
}

// GetAllbems gets all bems
func (s *BemService) GetAllBems(limit, offset int) ([]models.BEM, int64, error) {
    return s.repository.GetAllBems(limit, offset)
}

// Deletebem deletes a bem
func (s *BemService) DeleteBem(id uint) error {
	// Check if bem exists
	bem, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if bem == nil {
		return errors.New("gedung tidak ditemukan")
	}

	// Delete bem (soft delete)
	return s.repository.DeleteByID(id)
}

// bemWithStats represents a bem with additional statistics
type BemWithStats struct {
	Bem  models.BEM `json:"bem"`
	RoomCount int64           `json:"room_count"`
}

// GetbemWithStats gets a bem with its statistics
func (s *BemService) GetBemWithStats(id uint) (*BemWithStats, error) {
	// Get bem
	bem, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if bem == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}

	// Return bem with stats
	return &BemWithStats{
		Bem:  *bem,
	}, nil
}

// GetAllbemsWithStats gets all bems with their statistics
// func (s *bemService) GetAllbemsWithStats() ([]bemWithStats, error) {
// 	// Get all bems
// 	bems, err := s.repository.Get()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build response with stats
// 	result := make([]bemWithStats, len(bems))
// 	for i, bem := range bems {
		
// 		result[i] = bemWithStats{
// 			bem:  bem,
// 		}
// 	}

// 	return result, nil
// } 

func (s *BemService) GetAllLeaders() ([]models.Student, error) {
	return s.repository.GetAllLeaders()
}

