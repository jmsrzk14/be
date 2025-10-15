package handlers

import (
	"bem_be/internal/models"
	"bem_be/internal/services"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/microcosm-cc/bluemonday"
)

// NewsHandler menangani request HTTP terkait berita
type NewsHandler struct {
	service *services.NewsService
}

// NewNewsHandler membuat handler berita baru
func NewNewsHandler(db *gorm.DB) *NewsHandler {
	return &NewsHandler{
		service: services.NewNewsService(db),
	}
}

// Helper untuk parsing optional uint dari form data
func parseOptionalUint(value string) *uint {
	if value == "" {
		return nil
	}
	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return nil
	}
	u := uint(val)
	return &u
}

// GetAllNews mengembalikan semua berita dengan pagination
func (h *NewsHandler) GetAllNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	newsList, total, err := h.service.GetAllNews(perPage, offset)
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
		"data": newsList,
	})
}

// GetNewsByID mengembalikan berita berdasarkan ID
func (h *NewsHandler) GetNewsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stats := c.Query("stats")
	var result interface{}

	if stats == "true" {
		result, err = h.service.GetNewsWithStats(uint(id))
	} else {
		result, err = h.service.GetNewsByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "News retrieved successfully",
		"data":    result,
	})
}

// buat sanitizer sekali untuk dipakai ulang
var htmlSanitizer = bluemonday.UGCPolicy()

// CreateNews membuat berita baru (dengan unggahan file opsional)
func (h *NewsHandler) CreateNews(c *gin.Context) {
	var news models.News

	news.Title = c.PostForm("title")
	// sanitize content agar aman menyimpan HTML rich-text
	rawContent := c.PostForm("content")
	news.Content = htmlSanitizer.Sanitize(rawContent)

	news.Category = c.PostForm("category")
	news.BEMID = parseOptionalUint(c.PostForm("bem_id"))
	news.AssociationID = parseOptionalUint(c.PostForm("association_id"))
	news.DepartmentID = parseOptionalUint(c.PostForm("department_id"))

	file, err := c.FormFile("image")
	if err == nil {
		uploadPath := "uploads/news"
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Tidak dapat membuat folder unggahan",
			})
			return
		}

		// buat nama file unik
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		filePath := filepath.Join(uploadPath, fileName)

		// simpan file ke folder
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal menyimpan file",
			})
			return
		}

		// simpan hanya nama file ke database
		news.ImageURL = fileName
	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Gagal memproses file: " + err.Error(),
		})
		return
	}

	if err := h.service.CreateNews(&news); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Berita berhasil dibuat",
		"data":    news,
	})
}

// UpdateNews memperbarui berita yang ada.
func (h *NewsHandler) UpdateNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	existingNews, err := h.service.GetNewsByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// hanya update content jika ada input baru; sanitize sebelum set
	if content := c.PostForm("content"); content != "" {
		existingNews.Content = htmlSanitizer.Sanitize(content)
	}
	// update other fields only if provided
	if title := c.PostForm("title"); title != "" {
		existingNews.Title = title
	}
	if category := c.PostForm("category"); category != "" {
		existingNews.Category = category
	}
	if v := c.PostForm("bem_id"); v != "" {
		existingNews.BEMID = parseOptionalUint(v)
	}
	if v := c.PostForm("association_id"); v != "" {
		existingNews.AssociationID = parseOptionalUint(v)
	}
	if v := c.PostForm("department_id"); v != "" {
		existingNews.DepartmentID = parseOptionalUint(v)
	}

	file, err := c.FormFile("image")
	if err == nil {
		if existingNews.ImageURL != "" {
			_ = os.Remove(existingNews.ImageURL)
		}

		uploadPath := "uploads/news"
		_ = os.MkdirAll(uploadPath, os.ModePerm)

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		filePath := filepath.Join(uploadPath, fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan file"})
			return
		}
		existingNews.ImageURL = filePath
	}

	if err := h.service.UpdateNews(existingNews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berita berhasil diperbarui",
		"data":    existingNews,
	})
}

// DeleteNews menghapus sebuah berita
func (h *NewsHandler) DeleteNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	if err := h.service.DeleteNews(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berita berhasil dihapus",
	})
}

// RestoreNews menangani permintaan untuk memulihkan berita yang telah di-soft-delete.
func (h *NewsHandler) RestoreNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	restoredNews, err := h.service.RestoreNews(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berita berhasil dipulihkan",
		"data":    restoredNews,
	})
}
