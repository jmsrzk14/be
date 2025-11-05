package repositories

import (
	"errors"
	"time"
	"bem_be/internal/models"

	"gorm.io/gorm"
)

type CalenderRepository struct {
	db *gorm.DB
}

func NewCalenderRepository(db *gorm.DB) *CalenderRepository {
	return &CalenderRepository{db: db}
}

// ✅ Create Event
func (r *CalenderRepository) Create(event *models.Calender) error {
	if err := r.db.Create(event).Error; err != nil {
		return err
	}
	return nil
}

// ✅ Update Event
func (r *CalenderRepository) Update(event *models.Calender) error {
	if err := r.db.Save(event).Error; err != nil {
		return err
	}
	return nil
}

// ✅ Delete Event by ID
func (r *CalenderRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Calender{}, id).Error; err != nil {
		return err
	}
	return nil
}

// ✅ Find Event by ID
func (r *CalenderRepository) GetByID(id uint) (*models.Calender, error) {
	var event models.Calender
	if err := r.db.First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

// ✅ Get Events in Range (pakai tipe time.Time)
func (r *CalenderRepository) GetEventsInRange(start, end time.Time) ([]models.Calender, error) {
	var events []models.Calender
	if err := r.db.Where("end_time >= ? AND start_time <= ?", start, end).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// ✅ Get Events in Range (pakai string input — helper untuk service)
func (r *CalenderRepository) GetEventsInRangeString(startStr, endStr string) ([]models.Calender, error) {
	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return nil, errors.New("format start time tidak valid")
		}
	}
	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return nil, errors.New("format end time tidak valid")
		}
	}

	return r.GetEventsInRange(start, end)
}

// ✅ Get Events by Month & Year
func (r *CalenderRepository) GetEventsByMonthYear(month, year int) ([]models.Calender, error) {
	var events []models.Calender
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0) // Awal bulan berikutnya

	if err := r.db.Where("end_time >= ? AND start_time < ?", start, end).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *CalenderRepository) GetEventsCurrentMonth(start, end time.Time) ([]models.Calender, error) {
	var events []models.Calender
	if err := r.db.Preload("Organization").
		Where("end_time >= ? AND start_time <= ?", start, end).
		Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *CalenderRepository) GetAllEvents() ([]models.Calender, error) {
	var events []models.Calender
	if err := r.db.Preload("Organization").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
