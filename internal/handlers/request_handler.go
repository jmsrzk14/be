package handlers

import (
	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RequestHandler struct {
	service *services.RequestService
}

func NewRequestHandler(db *gorm.DB) *RequestHandler {
	return &RequestHandler{
		service: services.NewRequestService(db),
	}
}

func (h *RequestHandler) GetAllRequestsSarpras(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	category := 2

	requests, total, err := h.service.GetAllRequestsDepol(category, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	metadata := utils.PaginationMetadata{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  int(total),
		TotalPages:  totalPages,
		Links: utils.PaginationLinks{
			First: fmt.Sprintf("/requests?page=1&per_page=%d", perPage),
			Last:  fmt.Sprintf("/requests?page=%d&per_page=%d", totalPages, perPage),
		},
	}

	response := utils.MetadataFormatResponse(
		"success",
		"Berhasil mendapatkan daftar permintaan peminjaman kategori Depol",
		metadata,
		requests,
	)

	c.JSON(http.StatusOK, response)
}

func (h *RequestHandler) GetRequestByIDSarpras(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.service.GetRequestWithStatsSarpras(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Request retrieved successfully",
		"data":    result,
	})
}

func (h *RequestHandler) GetRequestsByUserIDSapras(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid user ID format", nil))
		return
	}

	requests, err := h.service.GetRequestsByRequesterIDSarpras(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseHandler("success", "Requests retrieved successfully", requests))
}

func (h *RequestHandler) CreateRequestSarpras(c *gin.Context) {
	userIDStr := c.PostForm("userID")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "userID is required", nil))
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid userID format", nil))
		return
	}

	studentPtr, err := h.service.GetStudentByUserIDSarpras(userID)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseHandler("error", "Student associated with this token not found", nil))
		return
	}
	student := *studentPtr

	uploadDir := filepath.Join("uploads", "requests")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Failed to create upload directory", nil))
		return
	}

	var request models.Request
	request.RequesterID = uint(student.UserID)
	request.RequestPlan = c.PostForm("startDate")
	request.ReturnPlan = c.PostForm("endDate")
	request.Activity = c.PostForm("tujuan")
	request.Location = c.PostForm("lokasi")
	request.Name = c.PostForm("nama")
	request.Category = 2

	itemsStr := c.PostForm("items")
	if itemsStr != "" {
		var itemIDs []int
		if err := json.Unmarshal([]byte(itemsStr), &itemIDs); err == nil {
			request.Item = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(itemIDs)), ","), "[]")
		} else {
			c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid items format", nil))
			return
		}
	}

	// === Upload file KTM ===
	ktmFile, err := c.FormFile("image_ktm")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "KTM image is required", nil))
		return
	}
	originalKtmFilename := strings.ReplaceAll(filepath.Base(ktmFile.Filename), " ", "_")
	ktmFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), originalKtmFilename)
	ktmPath := filepath.Join(uploadDir, ktmFilename)
	if err := c.SaveUploadedFile(ktmFile, ktmPath); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Failed to save KTM image", nil))
		return
	}
	request.ImageURLKTM = ktmFilename // ✅ hanya nama file

	// === Upload file Barang ===
	brgFile, err := c.FormFile("image_brg")
	if err == nil {
		originalBrgFilename := strings.ReplaceAll(filepath.Base(brgFile.Filename), " ", "_")
		brgFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), originalBrgFilename)
		brgPath := filepath.Join(uploadDir, brgFilename)
		if err := c.SaveUploadedFile(brgFile, brgPath); err != nil {
			c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Failed to save item image", nil))
			return
		}
		request.ImageURLBRG = brgFilename // ✅ hanya nama file
	}

	request.Status = "pending"

	if err := h.service.CreateRequestSarpras(&request); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Failed to create request: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.ResponseHandler("success", "Request created successfully", request))
}

func (h *RequestHandler) UpdateRequestSarpras(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var request *models.Request
	request, err = h.service.GetRequestByIDSarpras(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	if request.Status == "Approved" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update request after approval"})
		return
	}
	request.Name = c.PostForm("name")
	request.Item = c.PostForm("item")
	request.RequestPlan = c.PostForm("request_plan")
	request.ReturnPlan = c.PostForm("return_plan")
	request.UpdatedAt = time.Now()
	if err := h.service.UpdateRequestSarpras(request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Request updated successfully",
		"data":    request,
	})
}

func (h *RequestHandler) DeleteRequestSarpras(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var request *models.Request
	request, err = h.service.GetRequestByIDSarpras(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	if request.Status == "Approved" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete request after approval"})
		return
	}
	if err := h.service.DeleteRequestSarpras(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Request deleted successfully",
	})
}

// 1. Ambil ID request dari parameter URL
func (h *RequestHandler) UpdateRequestSarprasStatus(c *gin.Context) {
	idStr := c.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid request ID format", nil))
		return
	}

	// 2. Ambil data dari body JSON
	var input struct {
		Status string `json:"status" binding:"required"`
		UserID uint   `json:"user_id" binding:"required"`
		Reason string `json:"reason"` // optional, hanya dipakai saat rejected
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid input: 'status' and 'user_id' are required", err.Error()))
		return
	}

	// 3. Validasi nilai status
	if input.Status != "approved" && input.Status != "rejected" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid status value. Must be 'approved' or 'rejected'", nil))
		return
	}

	// 4. Jika status "rejected", pastikan alasan diberikan
	if input.Status == "rejected" && strings.TrimSpace(input.Reason) == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Reason is required when rejecting a request", nil))
		return
	}

	// 5. Kirim ke service untuk diproses
	updatedRequest, err := h.service.ProcessRequestStatusSarpras(uint(requestID), input.Status, int(input.UserID), input.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	// 6. Response sukses
	c.JSON(http.StatusOK, utils.ResponseHandler("success", "Request status updated successfully", updatedRequest))
}

func (h *RequestHandler) UploadImageBarangSarpras(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	// Ambil file dari form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Buat nama file unik
	filename := fmt.Sprintf("%d_%s", id, file.Filename)
	saveDir := "uploads/request/barang"
	savePath := fmt.Sprintf("%s/%s", saveDir, filename)

	// Pastikan folder ada
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// Simpan file ke disk
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update nama file dan ubah status jadi 'diambil'
	if err := h.service.UpdateImageBarangAndStatusSarpras(uint(id), filename, "diambil"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "File uploaded and status updated to 'diambil'",
		"data": gin.H{
			"file_name": filename,
			"status":    "diambil",
		},
	})
}

func (h *RequestHandler) ReturnBarangSarpras(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	// Waktu pengembalian sekarang
	returnTime := time.Now()

	// 1️⃣ Tambahkan kembali stok barang
	if err := h.service.UpdateItemStockOnTakenSarpras(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore item stock"})
		return
	}

	// 2️⃣ Update status & kolom ReturnedAt
	if err := h.service.ReturnedItemSarpras(uint(id), returnTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Barang berhasil dikembalikan dan stok diperbarui",
	})
}

func (h *RequestHandler) EndRequestBarangSarpras(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	// Ubah status request jadi "selesai"
	if err := h.service.UpdateStatusRequestSarpras(uint(id), "selesai"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Request status updated to 'selesai'",
	})
}

func (h *RequestHandler) GetAllRequestsDepol(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	category := 1

	requests, total, err := h.service.GetAllRequestsDepol(category, perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	metadata := utils.PaginationMetadata{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  int(total),
		TotalPages:  totalPages,
		Links: utils.PaginationLinks{
			First: fmt.Sprintf("/requests?page=1&per_page=%d", perPage),
			Last:  fmt.Sprintf("/requests?page=%d&per_page=%d", totalPages, perPage),
		},
	}

	response := utils.MetadataFormatResponse(
		"success",
		"Berhasil mendapatkan daftar permintaan peminjaman kategori Depol",
		metadata,
		requests,
	)

	c.JSON(http.StatusOK, response)
}

func (h *RequestHandler) GetRequestByIDDepol(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.service.GetRequestWithStatsDepol(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Request retrieved successfully",
		"data":    result,
	})
}

func (h *RequestHandler) GetRequestsByUserIDDepol(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid user ID format", nil))
		return
	}

	requests, err := h.service.GetRequestsByRequesterIDDepol(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseHandler("success", "Requests retrieved successfully", requests))
}

func (h *RequestHandler) CreateRequestDepol(c *gin.Context) {
	userIDStr := c.PostForm("userID")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "userID is required", nil))
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid userID format", nil))
		return
	}

	studentPtr, err := h.service.GetStudentByUserIDDepol(userID)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseHandler("error", "Student associated with this token not found", nil))
		return
	}
	student := *studentPtr

	uploadDir := filepath.Join("uploads", "requests")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Failed to create upload directory", nil))
		return
	}

	var request models.Request
	request.RequesterID = uint(student.UserID)
	request.RequestPlan = c.PostForm("startDate")
	request.ReturnPlan = c.PostForm("endDate")
	request.Activity = c.PostForm("tujuan")
	request.Location = c.PostForm("lokasi")
	request.Name = c.PostForm("nama")
	request.Category = 1

	itemsStr := c.PostForm("items")
	if itemsStr != "" {
		var itemIDs []int
		if err := json.Unmarshal([]byte(itemsStr), &itemIDs); err == nil {
			request.Item = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(itemIDs)), ","), "[]")
		} else {
			c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid items format", nil))
			return
		}
	}

	// === Upload file KTM ===
	ktmFile, err := c.FormFile("image_ktm")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "KTM image is required", nil))
		return
	}
	originalKtmFilename := strings.ReplaceAll(filepath.Base(ktmFile.Filename), " ", "_")
	ktmFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), originalKtmFilename)
	ktmPath := filepath.Join(uploadDir, ktmFilename)
	if err := c.SaveUploadedFile(ktmFile, ktmPath); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Failed to save KTM image", nil))
		return
	}
	request.ImageURLKTM = ktmFilename // ✅ hanya nama file

	// === Upload file Barang ===
	brgFile, err := c.FormFile("image_brg")
	if err == nil {
		originalBrgFilename := strings.ReplaceAll(filepath.Base(brgFile.Filename), " ", "_")
		brgFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), originalBrgFilename)
		brgPath := filepath.Join(uploadDir, brgFilename)
		if err := c.SaveUploadedFile(brgFile, brgPath); err != nil {
			c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Failed to save item image", nil))
			return
		}
		request.ImageURLBRG = brgFilename // ✅ hanya nama file
	}

	request.Status = "pending"

	if err := h.service.CreateRequestDepol(&request); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Failed to create request: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.ResponseHandler("success", "Request created successfully", request))
}

func (h *RequestHandler) UpdateRequestDepol(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var request *models.Request
	request, err = h.service.GetRequestByIDDepol(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	if request.Status == "Approved" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update request after approval"})
		return
	}
	request.Name = c.PostForm("name")
	request.Item = c.PostForm("item")
	request.RequestPlan = c.PostForm("request_plan")
	request.ReturnPlan = c.PostForm("return_plan")
	request.UpdatedAt = time.Now()
	if err := h.service.UpdateRequestDepol(request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Request updated successfully",
		"data":    request,
	})
}

func (h *RequestHandler) DeleteRequestDepol(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var request *models.Request
	request, err = h.service.GetRequestByIDDepol(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	if request.Status == "Approved" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete request after approval"})
		return
	}
	if err := h.service.DeleteRequestDepol(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Request deleted successfully",
	})
}

// 1. Ambil ID request dari parameter URL
func (h *RequestHandler) UpdateRequestDepolStatus(c *gin.Context) {
	idStr := c.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid request ID format", nil))
		return
	}

	// 2. Ambil data dari body JSON
	var input struct {
		Status string `json:"status" binding:"required"`
		UserID uint   `json:"user_id" binding:"required"`
		Reason string `json:"reason"` // optional, hanya dipakai saat rejected
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid input: 'status' and 'user_id' are required", err.Error()))
		return
	}

	// 3. Validasi nilai status
	if input.Status != "approved" && input.Status != "rejected" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid status value. Must be 'approved' or 'rejected'", nil))
		return
	}

	// 4. Jika status "rejected", pastikan alasan diberikan
	if input.Status == "rejected" && strings.TrimSpace(input.Reason) == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Reason is required when rejecting a request", nil))
		return
	}

	// 5. Kirim ke service untuk diproses
	updatedRequest, err := h.service.ProcessRequestStatusDepol(uint(requestID), input.Status, int(input.UserID), input.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	// 6. Response sukses
	c.JSON(http.StatusOK, utils.ResponseHandler("success", "Request status updated successfully", updatedRequest))
}

func (h *RequestHandler) UploadImageBarangDepol(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	// Ambil file dari form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Buat nama file unik
	filename := fmt.Sprintf("%d_%s", id, file.Filename)
	saveDir := "uploads/request/barang"
	savePath := fmt.Sprintf("%s/%s", saveDir, filename)

	// Pastikan folder ada
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// Simpan file ke disk
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update nama file dan ubah status jadi 'diambil'
	if err := h.service.UpdateImageBarangAndStatusDepol(uint(id), filename, "diambil"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "File uploaded and status updated to 'diambil'",
		"data": gin.H{
			"file_name": filename,
			"status":    "diambil",
		},
	})
}

func (h *RequestHandler) ReturnBarangDepol(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	// Waktu pengembalian sekarang
	returnTime := time.Now()

	// 1️⃣ Tambahkan kembali stok barang
	if err := h.service.UpdateItemStockOnTakenDepol(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore item stock"})
		return
	}

	// 2️⃣ Update status & kolom ReturnedAt
	if err := h.service.ReturnedItemDepol(uint(id), returnTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Barang berhasil dikembalikan dan stok diperbarui",
	})
}

func (h *RequestHandler) EndRequestBarangDepol(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	// Ubah status request jadi "selesai"
	if err := h.service.UpdateStatusRequestDepol(uint(id), "selesai"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Request status updated to 'selesai'",
	})
}
