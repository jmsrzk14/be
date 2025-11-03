package handlers

import (
	"net/http"

	"bem_be/internal/models"
	"bem_be/internal/repositories"

	"github.com/gin-gonic/gin"
)

// ===== Seeder =====
func CreateStatusSeeder(repo *repositories.StatusRepository) {
	count, err := repo.CountByID(1)
	if err != nil {
		panic(err)
	}

	// jika belum ada data, maka buat data awal dengan status 0
	if count == 0 {
		status := models.Status_Aspirations{Status: 0}
		if err := repo.Create(&status); err != nil {
			panic(err)
		}
	}
}

// ===== Get Status =====
func GetStatusAspirations(c *gin.Context, repo *repositories.StatusRepository) {
	var status models.Status_Aspirations
	if err := repo.DB.First(&status, 1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "status not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": status})
}

// ===== Update (toggle) Status =====
func UpdateStatusAspirations(c *gin.Context, repo *repositories.StatusRepository) {
	var status models.Status_Aspirations
	if err := repo.DB.First(&status, 1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "status not found"})
		return
	}

	// toggle status: 0 -> 1 atau 1 -> 0
	if status.Status == 0 {
		status.Status = 1
	} else {
		status.Status = 0
	}

	if err := repo.DB.Save(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "status updated successfully",
		"data":    status,
	})
}
