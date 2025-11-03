package repositories

import (
	"gorm.io/gorm"
	"bem_be/internal/models"
	"bem_be/internal/database"
)

type StatusRepository struct {
	DB *gorm.DB
}

// NewStatusRepository creates a new StatusRepository
func NewStatusRepository() *StatusRepository {
	return &StatusRepository{
		DB: database.DB,
	}
}

func (r *StatusRepository) CreateStatus(user *models.Status_Aspirations) error {
	return r.DB.Create(user).Error
}

func (r *StatusRepository) CountByID(id uint) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Status_Aspirations{}).Where("id = ?", id).Count(&count).Error
	return count, err
}

func (r *StatusRepository) Create(status *models.Status_Aspirations) error {
	return r.DB.Create(status).Error
}
