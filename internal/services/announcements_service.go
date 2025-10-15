package services

import (
	"gorm.io/gorm"
	"errors"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

// announcementService is a service for announcement operations
type AnnouncementService struct {
	repository *repositories.AnnouncementRepository
	db *gorm.DB
}

// NewannouncementService creates a new announcement service
func NewAnnouncementService(db *gorm.DB) *AnnouncementService {
    return &AnnouncementService{
        repository: repositories.NewAnnouncementRepository(),
    }
}

// Createannouncement creates a new announcement
func (s *AnnouncementService) Createannouncement(announcement *models.Announcement) error {
	// Check if code exists (including soft-deleted)
	// exists, err := s.repository.CheckNameExists(announcement.Name, 0)
	// if err != nil {
	// 	return err
	// }

	// if exists {
	// 	// Try to find a soft-deleted announcement with this code
	// 	deletedannouncement, err := s.repository.FindDeletedByName(announcement.Name)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if deletedannouncement != nil {
	// 		// Restore the soft-deleted announcement with updated data
	// 		deletedannouncement.Name = announcement.Name
			
	// 		// Restore the announcement
	// 		restoredannouncement, err := s.repository.RestoreByName(announcement.Name)
	// 		if err != nil {
	// 			return err
	// 		}
			
	// 		// Update with new data
	// 		restoredannouncement.Name = announcement.Name
			
	// 		return s.repository.Update(restoredannouncement)
	// 	}
		
	// 	return errors.New("kode gedung sudah digunakan")
	// }

	// Create announcement
	return s.repository.Create(announcement)
}

// Updateannouncement updates an existing announcement
func (s *AnnouncementService) Updateannouncement(announcement *models.Announcement) error {
	// Check if announcement exists
	existingAnnouncement, err := s.repository.FindByID(announcement.ID)
	if err != nil {
		return err
	}
	if existingAnnouncement == nil {
		return errors.New("himpunan tidak ditemukan")
	}

	// Update announcement
	return s.repository.Update(announcement)
}

// GetannouncementByID gets a announcement by ID
func (s *AnnouncementService) GetAnnouncementByID(id uint) (*models.Announcement, error) {
	announcement, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("berita tidak ditemukan")
		}
		return nil, err
	}
	return announcement, nil
}
// GetAllannouncements gets all announcements
func (s *AnnouncementService) GetAllAnnouncements(limit, offset int) ([]models.Announcement, int64, error) {
    return s.repository.GetAllAnnouncements(limit, offset)
}


// Deleteannouncement deletes a announcement
func (s *AnnouncementService) DeleteAnnouncement(id uint) error {
	// Check if announcement exists
	announcement, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if announcement == nil {
		return errors.New("gedung tidak ditemukan")
	}

	// Delete announcement (soft delete)
	return s.repository.DeleteByID(id)
}

// announcementWithStats represents a announcement with additional statistics
type AnnouncementWithStats struct {
	Announcement  models.Announcement `json:"announcement"`
	RoomCount int64           `json:"room_count"`
}

// GetannouncementWithStats gets a announcement with its statistics
func (s *AnnouncementService) GetAnnouncementWithStats(id uint) (*AnnouncementWithStats, error) {
	// Get announcement
	announcement, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if announcement == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}

	// Return announcement with stats
	return &AnnouncementWithStats{
		Announcement:  *announcement,
	}, nil
}

// GetAllannouncementsWithStats gets all announcements with their statistics
// func (s *announcementService) GetAllannouncementsWithStats() ([]announcementWithStats, error) {
// 	// Get all announcements
// 	announcements, err := s.repository.Get()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build response with stats
// 	result := make([]announcementWithStats, len(announcements))
// 	for i, announcement := range announcements {
		
// 		result[i] = announcementWithStats{
// 			announcement:  announcement,
// 		}
// 	}

// 	return result, nil
// } 