package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"

	"gorm.io/gorm"
)

type AspirationRepository struct {
	db *gorm.DB
}

func NewAspirationRepository() *AspirationRepository {
	return &AspirationRepository{
		db: database.GetDB(),
	}
}

func (r *AspirationRepository) Create(aspiration *models.Aspiration) error {
	return r.db.Create(aspiration).Error
}

func (r *AspirationRepository) Update(aspiration *models.Aspiration) error {
	return r.db.Save(aspiration).Error
}

func (r *AspirationRepository) FindByID(id uint) (*models.Aspiration, error) {
	var aspiration models.Aspiration
	err := r.db.First(&aspiration, id).Error
	if err != nil {
		return nil, err
	}
	return &aspiration, nil
}

func (r *AspirationRepository) GetAllAspirations(limit, offset int) ([]models.Aspiration, int64, error) {
	var aspirations []models.Aspiration
	var total int64

	query := r.db.Model(&models.Aspiration{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&aspirations).Error; err != nil {
		return nil, 0, err
	}

	return aspirations, total, nil
}

func (r *AspirationRepository) DeleteByID(id uint) error {
	return r.db.Delete(&models.Aspiration{}, id).Error
}
