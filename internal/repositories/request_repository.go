package repositories

import (
	"bem_be/internal/database"
	"bem_be/internal/models"
	"errors"
	"time"

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

func (r *RequestRepository) CreateSarpras(request *models.Request) error {
	return r.db.Create(request).Error
}

func (r *RequestRepository) UpdateSarpras(request *models.Request) error {
	return r.db.Save(request).Error
}

func (r *RequestRepository) FindByIDSarpras(id uint) (*models.Request, error) {
	var request models.Request
	if err := r.db.First(&request, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

func (r *RequestRepository) FindAllByRequesterIDSarpras(requesterID uint) ([]models.Request, error) {
	var requests []models.Request
	if err := r.db.Where("requester_id = ? AND category = 2", requesterID).Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *RequestRepository) FindItemsByIDsSarpras(itemIDs []uint) ([]models.Item, error) {
	var items []models.Item
	if err := r.db.Where("id IN ?", itemIDs).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RequestRepository) FindItemsByUserIDsSarpras(itemIDs []uint) ([]models.Item, error) {
	var items []models.Item
	if err := r.db.Where("id IN ?", itemIDs).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RequestRepository) GetAllRequestsSarpras(limit, offset int) ([]models.Request, int64, error) {
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

func (r *RequestRepository) DeleteByIDSarpras(id uint) error {
	return r.db.Delete(&models.Request{}, id).Error
}

func (r *RequestRepository) UpdateImageBarangAndStatusSarpras(id uint, fileName string, status string) error {
	return r.db.Model(&models.Request{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"image_url_brg": fileName,
			"status":        status,
		}).Error
}

func (r *RequestRepository) UpdateStatusSarpras(id uint, status string) error {
	var req models.Request

	if err := r.db.First(&req, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("request not found")
		}
		return err
	}

	req.Status = status
	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

func (r *RequestRepository) UpdateStatusAndReturnTimeSarpras(id uint, status string, returnedAt time.Time) error {
	return r.db.Model(&models.Request{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":    status,
			"return_at": returnedAt,
		}).Error
}

func (r *RequestRepository) IncreaseItemAmountSarpras(itemID uint, amount int) error {
	return r.db.Model(&models.Item{}).
		Where("id = ?", itemID).
		Update("amount", gorm.Expr("amount + ?", amount)).Error
}

func (r *RequestRepository) CreateDepol(request *models.Request) error {
	return r.db.Create(request).Error
}

func (r *RequestRepository) UpdateDepol(request *models.Request) error {
	return r.db.Save(request).Error
}

func (r *RequestRepository) FindByIDDepol(id uint) (*models.Request, error) {
	var request models.Request
	if err := r.db.First(&request, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

func (r *RequestRepository) FindAllByRequesterIDDepol(requesterID uint) ([]models.Request, error) {
	var requests []models.Request
	if err := r.db.Where("requester_id = ? AND category = 1", requesterID).Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *RequestRepository) FindItemsByIDsDepol(itemIDs []uint) ([]models.Item, error) {
	var items []models.Item
	if err := r.db.Where("id IN ?", itemIDs).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RequestRepository) FindItemsByUserIDsDepol(itemIDs []uint) ([]models.Item, error) {
	var items []models.Item
	if err := r.db.Where("id IN ?", itemIDs).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RequestRepository) GetAllRequestsDepol(limit, offset int) ([]models.Request, int64, error) {
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

func (r *RequestRepository) DeleteByIDDepol(id uint) error {
	return r.db.Delete(&models.Request{}, id).Error
}

func (r *RequestRepository) UpdateImageBarangAndStatusDepol(id uint, fileName string, status string) error {
	return r.db.Model(&models.Request{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"image_url_brg": fileName,
			"status":        status,
		}).Error
}

func (r *RequestRepository) UpdateStatusDepol(id uint, status string) error {
	var req models.Request

	if err := r.db.First(&req, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("request not found")
		}
		return err
	}

	req.Status = status
	if err := r.db.Save(&req).Error; err != nil {
		return err
	}

	return nil
}

func (r *RequestRepository) UpdateStatusAndReturnTimeDepol(id uint, status string, returnedAt time.Time) error {
	return r.db.Model(&models.Request{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":    status,
			"return_at": returnedAt,
		}).Error
}

func (r *RequestRepository) IncreaseItemAmountDepol(itemID uint, amount int) error {
	return r.db.Model(&models.Item{}).
		Where("id = ?", itemID).
		Update("amount", gorm.Expr("amount + ?", amount)).Error
}
