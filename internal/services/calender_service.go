package services

import (
	"errors"
	"time"
	"bem_be/internal/models"
	"bem_be/internal/repositories"

	"gorm.io/gorm"
)

type CalenderService struct {
	repository *repositories.CalenderRepository
	db         *gorm.DB
}

func NewCalenderService(db *gorm.DB) *CalenderService {
	repo := repositories.NewCalenderRepository(db)
	return &CalenderService{
		repository: repo,
		db:         db,
	}
}

func (s *CalenderService) CreateEvent(event *models.Calender) error {
	if event.Title == "" {
		return errors.New("judul event wajib diisi")
	}
	if event.EndTime.Before(event.StartTime) {
		return errors.New("waktu selesai harus setelah waktu mulai")
	}
	return (*s.repository).Create(event)
}

func (s *CalenderService) UpdateEvent(event *models.Calender) error {
	existingEvent, err := (*s.repository).GetByID(event.ID)
	if err != nil {
		return err
	}
	if existingEvent == nil {
		return errors.New("event tidak ditemukan")
	}

	// Validasi waktu
	if event.EndTime.Before(event.StartTime) {
		return errors.New("waktu selesai harus setelah waktu mulai")
	}

	return (*s.repository).Update(event)
}

func (s *CalenderService) GetEventByID(id uint) (*models.Calender, error) {
	return (*s.repository).GetByID(id)
}

func (s *CalenderService) GetAllEventsInRange(start, end string) ([]models.Calender, error) {
	return (*s.repository).GetEventsInRangeString(start, end)
}

func (s *CalenderService) GetEventsByMonthYear(month, year int) ([]models.Calender, error) {
	return (*s.repository).GetEventsByMonthYear(month, year)
}

func (s *CalenderService) DeleteEvent(id uint) error {
	event, err := (*s.repository).GetByID(id)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event tidak ditemukan")
	}
	return (*s.repository).Delete(id)
}

func (s *CalenderService) GetEventsCurrentMonth() ([]models.Calender, int, int, error) {
	now := time.Now()
	year, month := now.Year(), now.Month()

	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	events, err := s.repository.GetEventsCurrentMonth(start, end)
	if err != nil {
		return nil, 0, 0, err
	}

	return events, int(month), year, nil
}

func (s *CalenderService) GetAllEvents() ([]models.Calender, error) {
	events, err := s.repository.GetAllEvents()
	if err != nil {
		return nil, err
	}
	return events, nil
}
