package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"gorm.io/gorm"

	"bem_be/internal/database"
	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

// ItemService is a service for item operations
type ItemService struct {
	repository *repositories.ItemRepository
	db         *gorm.DB
}

// NewItemService creates a new item service
func NewItemService(db *gorm.DB) *ItemService {
	return &ItemService{
		repository: repositories.NewItemRepository(),
	}
}

// CreateItem creates a new item
func (s *ItemService) CreateItemSarpras(item *models.Item, file *multipart.FileHeader) error {
	// bikin folder kalau belum ada
	if err := os.MkdirAll("uploads/items", os.ModePerm); err != nil {
		return err
	}

	// nama file unik
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filepath := "uploads/items/" + filename

	// simpan file
	if err := saveUploadedFile(file, filepath); err != nil {
		return err
	}

	// simpan path/filename ke struct
	item.Image = filename

	// simpan ke DB
	return s.repository.Create(item)
}

// UpdateItem updates an existing item
func (s *ItemService) UpdateItemSarpras(item *models.Item) error {
	// Check if item exists
	existingItem, err := s.repository.FindByID(item.ID)
	if err != nil {
		return err
	}
	if existingItem == nil {
		return errors.New("himpunan tidak ditemukan")
	}

	// Update item
	return s.repository.Update(item)
}

// GetItemByID gets a item by ID
func (s *ItemService) GetItemByIDSarpras(id uint) (*models.Item, error) {
	return s.repository.FindByID(id)
}

// GetAllItems gets all items
func (s *ItemService) GetAllItemsSarpras(limit, offset int, search string) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	query := database.DB.Model(&models.Item{}).
		Where("category = ?", "2")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (s *ItemService) GetAllItemsGuestSarpras() ([]models.Item, error) {
	return s.repository.GetAllItemsGuestSarpras()
}

// DeleteItem deletes a item
func (s *ItemService) DeleteItemSarpras(id uint) error {
	// Check if item exists
	item, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("gedung tidak ditemukan")
	}

	// Delete item (soft delete)
	return s.repository.DeleteByID(id)
}

// ItemWithStats represents a item with additional statistics
type ItemSarprasWithStats struct {
	Item      models.Item `json:"item"`
	RoomCount int64       `json:"room_count"`
}

// GetItemWithStats gets a item with its statistics
func (s *ItemService) GetItemWithStatsSarpras(id uint) (*ItemSarprasWithStats, error) {
	// Get item
	item, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}

	// Return item with stats
	return &ItemSarprasWithStats{
		Item: *item,
	}, nil
}

// CreateItem creates a new item
func (s *ItemService) CreateItemDepol(item *models.Item, file *multipart.FileHeader) error {
	// bikin folder kalau belum ada
	if err := os.MkdirAll("uploads/items", os.ModePerm); err != nil {
		return err
	}

	// nama file unik
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filepath := "uploads/items/" + filename

	// simpan file
	if err := saveUploadedFile(file, filepath); err != nil {
		return err
	}

	// simpan path/filename ke struct
	item.Image = filename

	// simpan ke DB
	return s.repository.Create(item)
}

// UpdateItem updates an existing item
func (s *ItemService) UpdateItemDepol(item *models.Item) error {
	// Check if item exists
	existingItem, err := s.repository.FindByID(item.ID)
	if err != nil {
		return err
	}
	if existingItem == nil {
		return errors.New("himpunan tidak ditemukan")
	}

	// Update item
	return s.repository.Update(item)
}

// GetItemByID gets a item by ID
func (s *ItemService) GetItemByIDDepol(id uint) (*models.Item, error) {
	return s.repository.FindByID(id)
}

// GetAllItems gets all items
func (s *ItemService) GetAllItemsDepol(limit, offset int, search string) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	query := database.DB.Model(&models.Item{}).
		Where("category = ?", "1")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (s *ItemService) GetAllItemsGuestDepol() ([]models.Item, error) {
	return s.repository.GetAllItemsGuestDepol()
}

// DeleteItem deletes a item
func (s *ItemService) DeleteItemDepol(id uint) error {
	// Check if item exists
	item, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("gedung tidak ditemukan")
	}

	// Delete item (soft delete)
	return s.repository.DeleteByID(id)
}

// ItemWithStats represents a item with additional statistics
type ItemDepolWithStats struct {
	Item      models.Item `json:"item"`
	RoomCount int64       `json:"room_count"`
}

// GetItemWithStats gets a item with its statistics
func (s *ItemService) GetItemWithStatsDepol(id uint) (*ItemDepolWithStats, error) {
	// Get item
	item, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}

	// Return item with stats
	return &ItemDepolWithStats{
		Item: *item,
	}, nil
}
