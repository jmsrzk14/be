package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"

	"gorm.io/gorm"
)

// mpmRepository is a repository for mpm operations
type MpmRepository struct {
	db *gorm.DB
}

// NewmpmRepository creates a new mpm repository
func NewMpmRepository() *MpmRepository {
	return &MpmRepository{
		db: database.GetDB(),
	}
}

// Create creates a new mpm
func (r *MpmRepository) Create(mpm *models.MPM) error {
	return r.db.Create(mpm).Error
}

// Update updates an existing mpm
func (r *MpmRepository) Update(mpm *models.MPM) error {
	return r.db.Save(mpm).Error
}

// FindByID finds a mpm by ID
func (r *MpmRepository) FindByID(id uint) (*models.MPM, error) {
	var mpm models.MPM
	err := r.db.First(&mpm, id).Error
	if err != nil {
		return nil, err
	}
	return &mpm, nil
}

// FindByName finds a mpm by code
func (r *MpmRepository) FindByName(code string) (*models.MPM, error) {
	var mpm models.MPM
	err := r.db.Where("code = ?", code).First(&mpm).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &mpm, nil
}

// FindAll finds all mpms
func (r *MpmRepository) GetAllMpms(limit, offset int) ([]models.MPM, int64, error) {
	var mpms []models.MPM
	var total int64

	query := r.db.Model(&models.MPM{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&mpms).Error; err != nil {
		return nil, 0, err
	}

	return mpms, total, nil
}

// DeleteByID deletes a mpm by ID
func (r *MpmRepository) DeleteByID(id uint) error {
	// Use soft delete (don't use Unscoped())
	return r.db.Delete(&models.MPM{}, id).Error
}

// FindDeletedByName finds a soft-deleted mpm by code
func (r *MpmRepository) FindDeletedByName(code string) (*models.MPM, error) {
	var mpm models.MPM
	err := r.db.Unscoped().Where("code = ? AND deleted_at IS NOT NULL", code).First(&mpm).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &mpm, nil
}

// RestoreByName restores a soft-deleted mpm by code
func (r *MpmRepository) RestoreByName(code string) (*models.MPM, error) {
	// Find the deleted record
	deletedMpm, err := r.FindDeletedByName(code)
	if err != nil {
		return nil, err
	}
	if deletedMpm == nil {
		return nil, nil
	}

	// Restore the record
	if err := r.db.Unscoped().Model(&models.MPM{}).Where("id = ?", deletedMpm.ID).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}

	// Return the restored record
	return r.FindByID(deletedMpm.ID)
}

func (r *MpmRepository) GetMPMByPeriod(period string) (*models.MPM, error) {
	var mpm models.MPM

	err := r.db.
		Preload("Leader").
		Preload("CoLeader").
		Preload("Secretary1").
		Preload("Secretary2").
		Preload("Treasurer1").
		Preload("Treasurer2").
		Where("period = ? AND deleted_at IS NULL", period).
		First(&mpm).Error

	if err != nil {
		return nil, err
	}

	return &mpm, nil
}

func (r *MpmRepository) FindMpmByPeriod(period string) (map[string]interface{}, error) {
	var mpm models.MPM

	if err := r.db.Where("period = ?", period).First(&mpm).Error; err != nil {
		return nil, err
	}

	// Ambil data mahasiswa berdasarkan ID jabatan
	var leader, coLeader, secretary models.Student

	// Gunakan ID yang valid (tidak perlu pointer check karena uint default = 0)
	if mpm.LeaderID != 0 {
		r.db.First(&leader, mpm.LeaderID)
	}
	if mpm.CoLeaderID != 0 {
		r.db.First(&coLeader, mpm.CoLeaderID)
	}
	if mpm.SecretaryID != 0 {
		r.db.First(&secretary, mpm.Secretary)
	}

	// Gabungkan hasilnya dalam map
	result := map[string]interface{}{
		"period":    mpm.Period,
		"leader":    leader,
		"co_leader": coLeader,
		"secretary": secretary,
	}

	return result, nil
}
