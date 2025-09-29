package services

import (
	"bem_be/internal/models"
	"bem_be/internal/repositories"
	"errors"

	"gorm.io/gorm"
)

// VisiMisiService adalah service untuk operasi berita.
type VisiMisiService struct {
	repository *repositories.VisiMisiRepository
	db         *gorm.DB
}

// NewVisiMisiService membuat service berita baru.
func NewVisiMisiService(db *gorm.DB) *VisiMisiService {
	return &VisiMisiService{
		repository: repositories.NewVisiMisiRepository(),
		db:         db,
	}
}

func (s *VisiMisiService) GetVisiMisiByPeriod(period string) (*models.BEM, error) {
    var visiMisi models.BEM

    err := s.db.Where(
		"period = ?",
		period,
	).First(&visiMisi).Error
    if err != nil {
        return nil, err
    }

    return &visiMisi, nil
}


func (s *VisiMisiService) GetVisiMisiById(id uint) (*models.BEM, error) {
    var visiMisi models.BEM

    err := s.db.Where(
		"leader_id = ? OR co_leader_id = ? OR secretary1_id = ? OR secretary2_id = ? OR treasurer1_id = ? OR treasurer2_id = ?",
		id, id, id, id, id, id,
	).First(&visiMisi).Error
    if err != nil {
        return nil, err
    }

    return &visiMisi, nil
}

// GetVisiMisiByID mendapatkan berita berdasarkan ID.
func (s *VisiMisiService) UpdateVisiMisiByID(id uint, visi string, misi string) (*models.BEM, error) {
	var visiMisi models.BEM

	// Cari record dulu
	err := s.db.Where(
		"leader_id = ? OR co_leader_id = ? OR secretary1_id = ? OR secretary2_id = ? OR treasurer1_id = ? OR treasurer2_id = ?",
		id, id, id, id, id, id,
	).First(&visiMisi).Error

	if err != nil {
		return nil, err
	}

	// Update field
	visiMisi.Vision = visi
	visiMisi.Mission = misi

	// Simpan perubahan
	if err := s.db.Save(&visiMisi).Error; err != nil {
		return nil, err
	}

	return &visiMisi, nil
}

// GetAllVisiMisi mendapatkan semua berita dengan pagination.
func (s *VisiMisiService) GetAllVisiMisi(limit, offset int) ([]models.Period, int64, error) {
	return s.repository.GetAllVisiMisi(limit, offset)
}

// DeleteVisiMisi menghapus sebuah berita.
func (s *VisiMisiService) DeleteVisiMisi(id uint) error {
	_, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("berita yang akan dihapus tidak ditemukan")
		}
		return err
	}
	return s.repository.DeleteByID(id)
}

// RestoreVisiMisi memulihkan berita dan mengembalikan data yang telah dipulihkan.
func (s *VisiMisiService) RestoreVisiMisi(id uint) (*models.Period, error) {
	restoredVisiMisi, err := s.repository.RestoreByID(id)
	if err != nil {
		return nil, err
	}
	if restoredVisiMisi == nil {
		return nil, errors.New("berita tidak ditemukan atau sudah aktif (tidak dihapus)")
	}
	return restoredVisiMisi, nil
}
