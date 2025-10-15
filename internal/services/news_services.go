package services

import (
	"bem_be/internal/models"
	"bem_be/internal/repositories"
	"errors"

	"gorm.io/gorm"
)

// NewsService adalah service untuk operasi berita.
type NewsService struct {
	repository *repositories.NewsRepository
}

// NewNewsService membuat service berita baru.
func NewNewsService(db *gorm.DB) *NewsService {
	return &NewsService{
		repository: repositories.NewNewsRepository(),
	}
}

// CreateNews membuat berita baru.
func (s *NewsService) CreateNews(news *models.News) error {
	if news.Title == "" || news.Content == "" {
		return errors.New("judul dan konten tidak boleh kosong")
	}
	return s.repository.Create(news)
}
type NewsWithStats struct {
	News  models.News `json:"news"`
	RoomCount int64           `json:"room_count"`
}
func (s *NewsService) GetNewsWithStats(id uint) (*NewsWithStats, error) {
	// Get news
	news, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if news == nil {
		return nil, errors.New("berita tidak ditemukan")
	}

	// Return news with stats
	return &NewsWithStats{
		News:  *news,
	}, nil
}

// UpdateNews memperbarui berita yang ada.
func (s *NewsService) UpdateNews(news *models.News) error {
	return s.repository.Update(news)
}

// GetNewsByID mendapatkan berita berdasarkan ID.
func (s *NewsService) GetNewsByID(id uint) (*models.News, error) {
	news, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("berita tidak ditemukan")
		}
		return nil, err
	}
	return news, nil
}

// GetAllNews mendapatkan semua berita dengan pagination.
func (s *NewsService) GetAllNews(limit, offset int) ([]models.News, int64, error) {
	return s.repository.GetAllNews(limit, offset)
}

// DeleteNews menghapus sebuah berita.
func (s *NewsService) DeleteNews(id uint) error {
	_, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("berita yang akan dihapus tidak ditemukan")
		}
		return err
	}
	return s.repository.DeleteByID(id)
}

// RestoreNews memulihkan berita dan mengembalikan data yang telah dipulihkan.
func (s *NewsService) RestoreNews(id uint) (*models.News, error) {
	restoredNews, err := s.repository.RestoreByID(id)
	if err != nil {
		return nil, err
	}
	if restoredNews == nil {
		return nil, errors.New("berita tidak ditemukan atau sudah aktif (tidak dihapus)")
	}
	return restoredNews, nil
}
