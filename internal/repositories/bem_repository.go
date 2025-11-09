package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"

	"gorm.io/gorm"
)

// bemRepository is a repository for bem operations
type BemRepository struct {
	db *gorm.DB
}

// NewbemRepository creates a new bem repository
func NewBemRepository() *BemRepository {
	return &BemRepository{
		db: database.GetDB(),
	}
}

// Create creates a new bem
func (r *BemRepository) Create(bem *models.BEM) error {
	return r.db.Create(bem).Error
}

// Update updates an existing bem
func (r *BemRepository) Update(bem *models.BEM) error {
	return r.db.Save(bem).Error
}

// FindByID finds a bem by ID
func (r *BemRepository) FindByID(id uint) (*models.BEM, error) {
	var bem models.BEM
	err := r.db.First(&bem, id).Error
	if err != nil {
		return nil, err
	}
	return &bem, nil
}

// FindByName finds a bem by code
func (r *BemRepository) FindByName(code string) (*models.BEM, error) {
	var bem models.BEM
	err := r.db.Where("code = ?", code).First(&bem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &bem, nil
}

// FindAll finds all bems
func (r *BemRepository) GetAllBems(limit, offset int) ([]models.BEM, int64, error) {
	var bems []models.BEM
	var total int64

	query := r.db.Model(&models.BEM{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&bems).Error; err != nil {
		return nil, 0, err
	}

	return bems, total, nil
}

// DeleteByID deletes a bem by ID
func (r *BemRepository) DeleteByID(id uint) error {
	// Use soft delete (don't use Unscoped())
	return r.db.Delete(&models.BEM{}, id).Error
}

// FindDeletedByName finds a soft-deleted bem by code
func (r *BemRepository) FindDeletedByName(code string) (*models.BEM, error) {
	var bem models.BEM
	err := r.db.Unscoped().Where("code = ? AND deleted_at IS NOT NULL", code).First(&bem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &bem, nil
}

// RestoreByName restores a soft-deleted bem by code
func (r *BemRepository) RestoreByName(code string) (*models.BEM, error) {
	// Find the deleted record
	deletedBem, err := r.FindDeletedByName(code)
	if err != nil {
		return nil, err
	}
	if deletedBem == nil {
		return nil, nil
	}

	// Restore the record
	if err := r.db.Unscoped().Model(&models.BEM{}).Where("id = ?", deletedBem.ID).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}

	// Return the restored record
	return r.FindByID(deletedBem.ID)
}

func (r *BemRepository) GetAllLeaders() ([]models.Student, error) {
	var students []models.Student

	err := r.db.
		Joins("LEFT JOIN organizations ON organizations.id = students.organization_id").
		Preload("Organization").
		Where(`
			students.position IS NOT NULL
			AND students.position <> ''
			AND (
				students.position LIKE '%_bem%' OR
				students.position LIKE '%_mpm%' OR
				students.position = 'ketua_department'
			)
			AND students.position NOT LIKE '%_himpunan%'
			AND students.position NOT LIKE '%_ukm%'
			AND students.deleted_at IS NULL
		`).
		Order("students.position ASC").
		Find(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}

func (r *BemRepository) FindBemByPeriod(period string) (map[string]interface{}, error) {
	var bem models.BEM

	// Ambil data BEM berdasarkan period
	if err := r.db.Where("period = ?", period).First(&bem).Error; err != nil {
		return nil, err
	}

	// Ambil data mahasiswa berdasarkan ID jabatan
	var leader, coLeader, secretary1, secretary2, treasurer1, treasurer2 models.Student

	// Gunakan ID yang valid (tidak perlu pointer check karena uint default = 0)
	if bem.LeaderID != 0 {
		r.db.First(&leader, bem.LeaderID)
	}
	if bem.CoLeaderID != 0 {
		r.db.First(&coLeader, bem.CoLeaderID)
	}
	if bem.Secretary1ID != 0 {
		r.db.First(&secretary1, bem.Secretary1ID)
	}
	if bem.Secretary2ID != 0 {
		r.db.First(&secretary2, bem.Secretary2ID)
	}
	if bem.Treasurer1ID != 0 {
		r.db.First(&treasurer1, bem.Treasurer1ID)
	}
	if bem.Treasurer2ID != 0 {
		r.db.First(&treasurer2, bem.Treasurer2ID)
	}

	// Gabungkan hasilnya dalam map
	result := map[string]interface{}{
		"period":       bem.Period,
		"leader":       leader,
		"co_leader":    coLeader,
		"secretary_1":  secretary1,
		"secretary_2":  secretary2,
		"treasurer_1":  treasurer1,
		"treasurer_2":  treasurer2,
	}

	return result, nil
}
