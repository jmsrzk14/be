package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"
	"gorm.io/gorm"
)

// AssociationRepository is a repository for association operations
type AssociationRepository struct {
	db *gorm.DB
}

// NewAssociationRepository creates a new association repository
func NewAssociationRepository() *AssociationRepository {
	return &AssociationRepository{
		db: database.GetDB(),
	}
}

// Create creates a new association
func (r *AssociationRepository) Create(association *models.Organization) error {
	return r.db.Create(association).Error
}

// Update updates an existing association
func (r *AssociationRepository) Update(association *models.Organization) error {
	return r.db.Model(&models.Organization{}).
		Where("id = ?", association.ID).
		Omit("created_at, category").
		Updates(association).Error
}

// FindByID finds a association by ID
func (r *AssociationRepository) FindByID(id uint) (*models.Organization, error) {
	var association models.Organization
	err := r.db.First(&association, id).Error
	if err != nil {
		return nil, err
	}
	return &association, nil
}

func (r *AssociationRepository) FindByShortName(name string) (*models.Organization, error) {
	var association models.Organization
	err := r.db.Where("short_name = ?", name).First(&association).Error
	if err != nil {
		return nil, err
	}
	return &association, nil
}

// FindByName finds a association by code
func (r *AssociationRepository) FindByName(code string) (*models.Organization, error) {
	var association models.Organization
	err := r.db.Where("code = ?", code).First(&association).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &association, nil
}

// GetAllAssociations returns all associations from the database with optional search filter
func (r *AssociationRepository) GetAllAssociations(limit, offset int, search string) ([]models.Organization, int64, error) {
    var associations []models.Organization
    var total int64

    query := r.db.Model(&models.Organization{}).Where("category_id = ?", 3)

    if search != "" {
        likeSearch := "%" + search + "%"
        query = query.Where("LOWER(name) LIKE ?", likeSearch)
    }

    query.Count(&total)

    result := query.
        Order("name ASC").
        Limit(limit).
        Offset(offset).
        Find(&associations)

    return associations, total, result.Error
}


// DeleteByID deletes a association by ID
func (r *AssociationRepository) DeleteByID(id uint) error {
	// Use soft delete (don't use Unscoped())
	return r.db.Delete(&models.Organization{}, id).Error
}

// FindDeletedByName finds a soft-deleted association by code
func (r *AssociationRepository) FindDeletedByName(code string) (*models.Organization, error) {
	var association models.Organization
	err := r.db.Unscoped().Where("code = ? AND deleted_at IS NOT NULL", code).First(&association).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &association, nil
}

// RestoreByName restores a soft-deleted association by code
func (r *AssociationRepository) RestoreByName(code string) (*models.Organization, error) {
	// Find the deleted record
	deletedAssociation, err := r.FindDeletedByName(code)
	if err != nil {
		return nil, err
	}
	if deletedAssociation == nil {
		return nil, nil
	}
	
	// Restore the record
	if err := r.db.Unscoped().Model(&models.Organization{}).Where("id = ?", deletedAssociation.ID).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}
	
	// Return the restored record
	return r.FindByID(deletedAssociation.ID)
}

func (r *AssociationRepository) GetAllAssociationsGuest() ([]models.Organization, error) {
    var associations []models.Organization
    err := r.db.Where("category_id = ?", 3).Find(&associations).Error
    return associations, err
}

func (r *AssociationRepository) FindAdminByShortNameAndPeriod(shortName string, period string) (*models.Period, error) {
	var result models.Period

	err := r.db.
		Preload("Leader").
		Preload("CoLeader").
		Preload("Secretary1").
		Preload("Secretary2").
		Preload("Treasurer1").
		Preload("Treasurer2").
		Where("organizations.short_name = ? AND periods.period = ?", shortName, period).
		First(&result).Error

	if err != nil {
		return nil, err
	}
	return &result, nil
}