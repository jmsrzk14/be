package handlers

import (
	"bem_be/internal/services"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// visimisiHandler menangani request HTTP terkait berita
type visimisiHandler struct {
	service *services.VisiMisiService
}

// NewvisimisiHandler membuat handler berita baru
func NewVisiMisiHandler(db *gorm.DB) *visimisiHandler {
	return &visimisiHandler{
		service: services.NewVisiMisiService(db),
	}
}

// GetAllVisiMisi mengembalikan semua berita dengan pagination
func (h *visimisiHandler) GetAllVisiMisi(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	visimisiList, total, err := h.service.GetAllVisiMisi(perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var totalPages int
	if perPage > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(perPage)))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mendapatkan daftar berita",
		"metadata": gin.H{
			"current_page": page,
			"per_page":     perPage,
			"total_items":  total,
			"total_pages":  totalPages,
		},
		"data": visimisiList,
	})
}

// GetVisiMisiByID mengembalikan berita berdasarkan ID
// func (h *visimisiHandler) GetVisiMisiByID(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseUint(idStr, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
// 		return
// 	}

// 	visimisi, err := h.service.GetVisiMisiByID(uint(id))
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  "success",
// 		"message": "Berita berhasil didapatkan",
// 		"data":    visimisi,
// 	})
// }

func (h *visimisiHandler) GetVisiMisiById(c *gin.Context) {
	userId := c.Param("id")

    id, err := strconv.ParseUint(userId, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
        return
    }

    visiMisi, err := h.service.GetVisiMisiById(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "status":  "error",
            "message": "Visi misi untuk periode ini tidak ditemukan",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Data visi misi ditemukan",
        "data":    visiMisi,
    })
}

func (h *visimisiHandler) GetVisiMisiByPeriod(c *gin.Context) {
	period := c.Param("period")

    visiMisi, err := h.service.GetVisiMisiByPeriod(period)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "status":  "error",
            "message": "Visi misi untuk periode ini tidak ditemukan",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Data visi misi ditemukan",
        "data":    visiMisi,
    })
}

// UpdateVisiMisi memperbarui berita yang ada.
func (h *visimisiHandler) UpdateVisiMisi(c *gin.Context) {
    userId := c.Param("id")

    id, err := strconv.ParseUint(userId, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
        return
    }

    var req struct {
        Visi string `json:"visi"`
        Misi string `json:"misi"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format JSON tidak valid"})
        return
    }

    updatedVisiMisi, err := h.service.UpdateVisiMisiByID(uint(id), req.Visi, req.Misi)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Visi misi berhasil diperbarui",
        "data":    updatedVisiMisi,
    })
}

// DeleteVisiMisi menghapus sebuah berita
func (h *visimisiHandler) DeleteVisiMisi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	if err := h.service.DeleteVisiMisi(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berita berhasil dihapus",
	})
}

// RestoreVisiMisi menangani permintaan untuk memulihkan berita yang telah di-soft-delete.
func (h *visimisiHandler) RestoreVisiMisi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	restoredVisiMisi, err := h.service.RestoreVisiMisi(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berita berhasil dipulihkan",
		"data":    restoredVisiMisi,
	})
}
