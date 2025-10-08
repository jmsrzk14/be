package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"
	"errors"

	"gorm.io/gorm"
)

type RequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository() *RequestRepository {
	return &RequestRepository{
		db: database.GetDB(),
	}
}

func (r *RequestRepository) Create(request *models.Request) error {
	return r.db.Create(request).Error
}

func (r *RequestRepository) Update(request *models.Request) error {
	return r.db.Save(request).Error
}

func (r *RequestRepository) FindByID(id uint) (*models.Request, error) {
	var request models.Request
	if err := r.db.First(&request, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

func (r *RequestRepository) FindAllByRequesterID(requesterID uint) ([]models.Request, error) {
	var requests []models.Request
	if err := r.db.Where("requester_id = ?", requesterID).Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *RequestRepository) FindItemsByIDs(itemIDs []uint) ([]models.Item, error) {
	var items []models.Item
	if err := r.db.Where("id IN ?", itemIDs).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RequestRepository) FindItemsByUserIDs(itemIDs []uint) ([]models.Item, error) {
	var items []models.Item
	if err := r.db.Where("id IN ?", itemIDs).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RequestRepository) GetAllRequests(limit, offset int) ([]models.Request, int64, error) {
	var requests []models.Request
	var total int64

	query := r.db.Model(&models.Request{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

func (r *RequestRepository) DeleteByID(id uint) error {
	return r.db.Delete(&models.Request{}, id).Error
}

func (r *RequestRepository) UpdateImageBarangAndStatus(id uint, fileName string, status string) error {
	return r.db.Model(&models.Request{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"image_url_brg": fileName,
			"status":        status,
		}).Error
}
