package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"
	"gorm.io/gorm"
)

// ClubRepository is a repository for club operations
type ClubRepository struct {
	db *gorm.DB
}

// NewClubRepository creates a new club repository
func NewClubRepository() *ClubRepository {
	return &ClubRepository{
		db: database.GetDB(),
	}
}

// Create creates a new club
func (r *ClubRepository) Create(club *models.Organization) error {
	return r.db.Create(club).Error
}

// Update updates an existing club
func (r *ClubRepository) Update(club *models.Organization) error {
	return r.db.Model(&models.Organization{}).
		Where("id = ?", club.ID).
		Omit("created_at, category").
		Updates(club).Error
}

// FindByID finds a club by ID
func (r *ClubRepository) FindByID(id uint) (*models.Organization, error) {
	var club models.Organization
	err := r.db.First(&club, id).Error
	if err != nil {
		return nil, err
	}
	return &club, nil
}

// FindByName finds a club by code
func (r *ClubRepository) FindByName(code string) (*models.Organization, error) {
	var club models.Organization
	err := r.db.Where("code = ?", code).First(&club).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &club, nil
}

// FindAll finds all clubs
func (r *ClubRepository) GetAllClubs(limit, offset int, search string) ([]models.Organization, int64, error) {
    var clubs []models.Organization
    var total int64

    query := r.db.Model(&models.Organization{}).Where("category_id = ?", 1)

    if search != "" {
        likeSearch := "%" + search + "%"
        query = query.Where("LOWER(name) LIKE ?", likeSearch)
    }

    query.Count(&total)

    result := query.
        Order("name ASC").
        Limit(limit).
        Offset(offset).
        Find(&clubs)

    return clubs, total, result.Error
}


// DeleteByID deletes a club by ID
func (r *ClubRepository) DeleteByID(id uint) error {
	// Use soft delete (don't use Unscoped())
	return r.db.Delete(&models.Organization{}, id).Error
}

// FindDeletedByName finds a soft-deleted club by code
func (r *ClubRepository) FindDeletedByName(code string) (*models.Organization, error) {
	var club models.Organization
	err := r.db.Unscoped().Where("code = ? AND deleted_at IS NOT NULL", code).First(&club).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &club, nil
}

// RestoreByName restores a soft-deleted club by code
func (r *ClubRepository) RestoreByName(code string) (*models.Organization, error) {
	// Find the deleted record
	deletedClub, err := r.FindDeletedByName(code)
	if err != nil {
		return nil, err
	}
	if deletedClub == nil {
		return nil, nil
	}
	
	// Restore the record
	if err := r.db.Unscoped().Model(&models.Organization{}).Where("id = ?", deletedClub.ID).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}
	
	// Return the restored record
	return r.FindByID(deletedClub.ID)
}

// // CheckNameExists checks if a code exists, including soft-deleted records
// func (r *ClubRepository) CheckNameExists(code string, excludeID uint) (bool, error) {
// 	var count int64
// 	query := r.db.Unscoped().Model(&models.Organization{}).Where("code = ?", code)
	
// 	// Exclude the current record if updating
// 	if excludeID > 0 {
// 		query = query.Where("id != ?", excludeID)
// 	}
	
// 	err := query.Count(&count).Error
// 	if err != nil {
// 		return false, err
// 	}
	
// 	return count > 0, nil
// } 

func (r *ClubRepository) GetAllClubsGuest() ([]models.Organization, error) {
    var clubs []models.Organization
    err := r.db.Where("category_id = ?", 1).Find(&clubs).Error
    return clubs, err
}
