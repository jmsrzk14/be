package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"
	"gorm.io/gorm"
)

// ItemRepository is a repository for department operations
type ItemRepository struct {
	db *gorm.DB
}

// NewItemRepository creates a new department repository
func NewItemRepository() *ItemRepository {
	return &ItemRepository{
		db: database.GetDB(),
	}
}

// Create creates a new department
func (r *ItemRepository) Create(department *models.Item) error {
	return r.db.Create(department).Error
}

// Update updates an existing department
func (r *ItemRepository) Update(department *models.Item) error {
	return r.db.Save(department).Error
}

// FindByID finds a department by ID
func (r *ItemRepository) FindByID(id uint) (*models.Item, error) {
	var department models.Item
	err := r.db.First(&department, id).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

// FindByName finds a department by code
func (r *ItemRepository) FindByName(code string) (*models.Item, error) {
	var department models.Item
	err := r.db.Where("code = ?", code).First(&department).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &department, nil
}

// GetAllAssociations returns all associations from the database with optional search filter
func (r *ItemRepository) GetAllItemsSarpras(limit, offset int, search string) ([]models.Item, int64, error) {
    var departments []models.Item
    var total int64

    query := r.db.Model(&models.Item{})

    if search != "" {
        likeSearch := "%" + search + "%"
        query = query.Where("name LIKE ? AND category = 1", likeSearch)
    }

    query.Count(&total)

    result := query.
        Order("name ASC").
        Limit(limit).
        Offset(offset).
        Find(&departments)

    return departments, total, result.Error
}

func (r *ItemRepository) GetAllItemsDepol(limit, offset int, search string) ([]models.Item, int64, error) {
    var departments []models.Item
    var total int64

    query := r.db.Model(&models.Item{})

    if search != "" {
        likeSearch := "%" + search + "%"
        query = query.Where("name LIKE ? AND category = 2", likeSearch)
    }

    query.Count(&total)

    result := query.
        Order("name ASC").
        Limit(limit).
        Offset(offset).
        Find(&departments)

    return departments, total, result.Error
}

// DeleteByID deletes a department by ID
func (r *ItemRepository) DeleteByID(id uint) error {
	// Use soft delete (don't use Unscoped())
	return r.db.Delete(&models.Item{}, id).Error
}

// FindDeletedByName finds a soft-deleted department by code
func (r *ItemRepository) FindDeletedByName(code string) (*models.Item, error) {
	var department models.Item
	err := r.db.Unscoped().Where("code = ? AND deleted_at IS NOT NULL", code).First(&department).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &department, nil
}

// RestoreByName restores a soft-deleted department by code
func (r *ItemRepository) RestoreByName(code string) (*models.Item, error) {
	// Find the deleted record
	deletedItem, err := r.FindDeletedByName(code)
	if err != nil {
		return nil, err
	}
	if deletedItem == nil {
		return nil, nil
	}
	
	// Restore the record
	if err := r.db.Unscoped().Model(&models.Item{}).Where("id = ?", deletedItem.ID).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}
	
	// Return the restored record
	return r.FindByID(deletedItem.ID)
}

func (r *ItemRepository) GetAllItemsGuestSarpras() ([]models.Item, error) {
    var associations []models.Item
    err := r.db.Where("category_id = ?", 1).Find(&associations).Error
    return associations, err
}

func (r *ItemRepository) GetAllItemsGuestDepol() ([]models.Item, error) {
    var associations []models.Item
    err := r.db.Where("category_id = ?", 2).Find(&associations).Error
    return associations, err
}
