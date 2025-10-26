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

// DepartmentService is a service for department operations
type DepartmentService struct {
	repository *repositories.DepartmentRepository
	db *gorm.DB
}

// NewDepartmentService creates a new department service
func NewDepartmentService(db *gorm.DB) *DepartmentService {
    return &DepartmentService{
        repository: repositories.NewDepartmentRepository(),
    }
}

// CreateDepartment creates a new department
func (s *DepartmentService) CreateDepartment(department *models.Organization, file *multipart.FileHeader) error {
	// bikin folder kalau belum ada
	if err := os.MkdirAll("uploads/departments", os.ModePerm); err != nil {
		return err
	}

	// nama file unik
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filepath := "uploads/departments/" + filename

	// simpan file
	if err := saveUploadedFile(file, filepath); err != nil {
		return err
	}

	// simpan path/filename ke struct
	department.Image = filename

	// simpan ke DB
	return s.repository.Create(department)
}

// UpdateDepartment updates an existing department
func (s *DepartmentService) UpdateDepartment(department *models.Organization) error {
	// Check if department exists
	existingDepartment, err := s.repository.FindByID(department.ID)
	if err != nil {
		return err
	}
	if existingDepartment == nil {
		return errors.New("himpunan tidak ditemukan")
	}

	// Update department
	return s.repository.Update(department)
}

// GetDepartmentByID gets a department by ID
func (s *DepartmentService) GetDepartmentByID(id uint) (*models.Organization, error) {
	return s.repository.FindByID(id)
}

// GetAllDepartments gets all departments
func (s *DepartmentService) GetAllDepartments(limit, offset int, search string) ([]models.Organization, int64, error) {
    return s.repository.GetAllDepartments(limit, offset, search)
}

func (s *DepartmentService) GetAllOrganizations(limit, offset int, search string) ([]models.Organization, int64, error) {
    return s.repository.GetAllOrganizations(limit, offset, search)
}

func (s *DepartmentService) GetAllDepartmentsGuest() ([]models.Organization, error) {
    return s.repository.GetAllDepartmentsGuest()
}

// DeleteDepartment deletes a department
func (s *DepartmentService) DeleteDepartment(id uint) error {
	// Check if department exists
	department, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	if department == nil {
		return errors.New("gedung tidak ditemukan")
	}

	// Delete department (soft delete)
	return s.repository.DeleteByID(id)
}

// DepartmentWithStats represents a department with additional statistics
type DepartmentWithStats struct {
	Department  models.Organization `json:"department"`
	RoomCount int64           `json:"room_count"`
}

// GetDepartmentWithStats gets a department with its statistics
func (s *DepartmentService) GetDepartmentWithStats(id uint) (*DepartmentWithStats, error) {
	// Get department
	department, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if department == nil {
		return nil, errors.New("gedung tidak ditemukan")
	}

	// Return department with stats
	return &DepartmentWithStats{
		Department:  *department,
	}, nil
}

// GetAllDepartmentsWithStats gets all departments with their statistics
// func (s *DepartmentService) GetAllDepartmentsWithStats() ([]DepartmentWithStats, error) {
// 	// Get all departments
// 	departments, err := s.repository.Get()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build response with stats
// 	result := make([]DepartmentWithStats, len(departments))
// 	for i, department := range departments {
		
// 		result[i] = DepartmentWithStats{
// 			Department:  department,
// 		}
// 	}

// 	return result, nil
// } 