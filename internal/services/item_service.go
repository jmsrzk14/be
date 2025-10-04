package services

import (
	"gorm.io/gorm"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"bem_be/internal/models"
	"bem_be/internal/repositories"
)

// ItemService is a service for item operations
type ItemService struct {
	repository *repositories.ItemRepository
	db *gorm.DB
}

// NewItemService creates a new item service
func NewItemService(db *gorm.DB) *ItemService {
    return &ItemService{
        repository: repositories.NewItemRepository(),
    }
}

// CreateItem creates a new item
func (s *ItemService) CreateItem(item *models.Item, file *multipart.FileHeader) error {
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
func (s *ItemService) UpdateItem(item *models.Item) error {
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
func (s *ItemService) GetItemByID(id uint) (*models.Item, error) {
	return s.repository.FindByID(id)
}

// GetAllItems gets all items
func (s *ItemService) GetAllItems(limit, offset int, search string) ([]models.Item, int64, error) {
    return s.repository.GetAllItems(limit, offset, search)
}

func (s *ItemService) GetAllItemsGuest() ([]models.Item, error) {
    return s.repository.GetAllItemsGuest()
}

// DeleteItem deletes a item
func (s *ItemService) DeleteItem(id uint) error {
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
type ItemWithStats struct {
	Item  models.Item `json:"item"`
	RoomCount int64           `json:"room_count"`
}

// GetItemWithStats gets a item with its statistics
func (s *ItemService) GetItemWithStats(id uint) (*ItemWithStats, error) {
	// Get item
	item, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}

	// Return item with stats
	return &ItemWithStats{
		Item:  *item,
	}, nil
}