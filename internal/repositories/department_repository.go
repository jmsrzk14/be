package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"
	"gorm.io/gorm"
)

// DepartmentRepository is a repository for department operations
type DepartmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository creates a new department repository
func NewDepartmentRepository() *DepartmentRepository {
	return &DepartmentRepository{
		db: database.GetDB(),
	}
}

// Create creates a new department
func (r *DepartmentRepository) Create(department *models.Organization) error {
	return r.db.Create(department).Error
}

// Update updates an existing department
func (r *DepartmentRepository) Update(department *models.Organization) error {
	return r.db.Model(&models.Organization{}).
		Where("id = ?", department.ID).
		Omit("created_at, category").
		Updates(department).Error
}
// FindByID finds a department by ID
func (r *DepartmentRepository) FindByID(id uint) (*models.Organization, error) {
	var department models.Organization
	err := r.db.First(&department, id).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

// FindByName finds a department by code
func (r *DepartmentRepository) FindByName(code string) (*models.Organization, error) {
	var department models.Organization
	err := r.db.Where("code = ?", code).First(&department).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &department, nil
}

// GetAllDepartments returns all departments from the database with optional search filter
func (r *DepartmentRepository) GetAllDepartments(limit, offset int, search string) ([]models.Organization, int64, error) {
    var departments []models.Organization
    var total int64

    query := r.db.Model(&models.Organization{}).Where("category_id = ?", 2)

    if search != "" {
        likeSearch := "%" + search + "%"
        query = query.Where("LOWER(name) LIKE ?", likeSearch)
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
func (r *DepartmentRepository) DeleteByID(id uint) error {
	// Use soft delete (don't use Unscoped())
	return r.db.Delete(&models.Organization{}, id).Error
}

// FindDeletedByName finds a soft-deleted department by code
func (r *DepartmentRepository) FindDeletedByName(code string) (*models.Organization, error) {
	var department models.Organization
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
func (r *DepartmentRepository) RestoreByName(code string) (*models.Organization, error) {
	// Find the deleted record
	deletedDepartment, err := r.FindDeletedByName(code)
	if err != nil {
		return nil, err
	}
	if deletedDepartment == nil {
		return nil, nil
	}
	
	// Restore the record
	if err := r.db.Unscoped().Model(&models.Organization{}).Where("id = ?", deletedDepartment.ID).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}
	
	// Return the restored record
	return r.FindByID(deletedDepartment.ID)
}

// // CheckNameExists checks if a code exists, including soft-deleted records
// func (r *DepartmentRepository) CheckNameExists(code string, excludeID uint) (bool, error) {
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

func (r *DepartmentRepository) GetAllDepartmentsGuest() ([]models.Organization, error) {
    var departments []models.Organization
    err := r.db.Where("category_id = ?", 2).Find(&departments).Error
    return departments, err
}
