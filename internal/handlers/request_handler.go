package handlers

import (
	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"
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

func (h *RequestHandler) GetAllRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	request, total, err := h.service.GetAllRequests(perPage, offset)
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
		"Berhasil mendapatkan daftar permintaan peminjaman",
		metadata,
		request,
	)

	c.JSON(http.StatusOK, response)
}

func (h *RequestHandler) GetRequestByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := h.service.GetRequestWithStats(uint(id))
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

// func (h *RequestHandler) CreateRequest(c *gin.Context) {
// 	var request models.Request
// 	request.Name = c.PostForm("name")
// 	quantityStr := c.DefaultPostForm("quantity", "1")
// 	quantity, err := strconv.ParseUint(quantityStr, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
// 		return
// 	}
// 	request.Quantity = uint(quantity)
// 	request.RequestPlan = c.PostForm("request_plan")
// 	request.ReturnPlan = c.PostForm("return_plan")
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}
// 	request.RequesterID = int(userID.(uint))

// 	var student *models.Student
// 	student, err = h.service.GetStudentByUserID(request.RequesterID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
// 		return
// 	}
// 	request.OrganizationID = student.OrganizationID
// 	request.OrganizationName = student.Organization.Name
// 	request.Status = "Pending"
// 	request.CreatedAt = time.Now()
// 	request.UpdatedAt = time.Now()
// 	if err := h.service.CreateRequest(&request); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, gin.H{"message": "Request created successfully", "data": request})
// }

func (h *RequestHandler) CreateRequest(c *gin.Context) {
	// --- 1. Ambil User ID dari context ---
	userIDClaim, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler("error", "Unauthorized: userID not found in context", nil))
		return
	}

	// --- 2. Cari student berdasarkan User ID dari context (dengan type assertion yang aman) ---
	var userID int
	switch v := userIDClaim.(type) {
	case uint:
		userID = int(v)
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Invalid userID format in context", nil))
		return
	}

	var student models.Student
	studentPtr, err := h.service.GetStudentByUserID(userID)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseHandler("error", "Student associated with this token not found", nil))
		return
	}
	student = *studentPtr // Dereference pointer

	// --- 3. Buat direktori untuk uploads jika belum ada ---
	uploadDir := filepath.Join("uploads", "requests")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Failed to create upload directory", nil))
		return
	}

	// --- 4. Ambil data dari multipart/form-data ---
	var request models.Request
	request.Name = c.PostForm("name")
	request.Activity = c.PostForm("activity")
	request.Location = c.PostForm("location")
	request.RequestPlan = c.PostForm("request_plan")
	request.ReturnPlan = c.PostForm("return_plan")

	quantityStr := c.DefaultPostForm("quantity", "1")
	quantity, err := strconv.ParseUint(quantityStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid quantity format", nil))
		return
	}
	request.Quantity = uint(quantity)
	request.RequesterID = student.UserID

	// --- 5. Proses Upload File ---
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
	request.ImageURLKTM = ktmPath

	brgFile, err := c.FormFile("image_brg")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Item image is required", nil))
		return
	}
	originalBrgFilename := strings.ReplaceAll(filepath.Base(brgFile.Filename), " ", "_")
	brgFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), originalBrgFilename)
	brgPath := filepath.Join(uploadDir, brgFilename)
	if err := c.SaveUploadedFile(brgFile, brgPath); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Failed to save item image", nil))
		return
	}
	request.ImageURLBRG = brgPath

	// --- 6. Set nilai default dan panggil service ---
	request.Status = "pending"
	if err := h.service.CreateRequest(&request); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Failed to create request: "+err.Error(), nil))
		return
	}

	// --- 7. Kirim response sukses ---
	c.JSON(http.StatusCreated, utils.ResponseHandler("success", "Request created successfully", request))
}

func (h *RequestHandler) UpdateRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var request *models.Request
	request, err = h.service.GetRequestByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	if request.Status == "Approved" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update request after approval"})
		return
	}
	request.Name = c.PostForm("name")
	quantityStr := c.DefaultPostForm("quantity", "1")
	quantity, err := strconv.ParseUint(quantityStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}
	request.Quantity = uint(quantity)
	request.RequestPlan = c.PostForm("request_plan")
	request.ReturnPlan = c.PostForm("return_plan")
	request.UpdatedAt = time.Now()
	if err := h.service.UpdateRequest(request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Request updated successfully",
		"data":    request,
	})
}

func (h *RequestHandler) DeleteRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var request *models.Request
	request, err = h.service.GetRequestByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	if request.Status == "Approved" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete request after approval"})
		return
	}
	if err := h.service.DeleteRequest(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Request deleted successfully",
	})
}

func (h *RequestHandler) UpdateRequestStatus(c *gin.Context) {
	// 1. Ambil ID request dari parameter URL (misal: /api/admin/request/123/status)
	idStr := c.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid request ID format", nil))
		return
	}

	// 2. Ambil data 'status' dari body JSON yang dikirim admin
	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid input: 'status' field is required in JSON body", err.Error()))
		return
	}

	// 3. Validasi nilai status yang diizinkan
	if input.Status != "approved" && input.Status != "rejected" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler("error", "Invalid status value. Must be 'approved' or 'rejected'", nil))
		return
	}

	// 4. Ambil ID admin yang sedang login dari context (disediakan oleh middleware)
	userIDClaim, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler("error", "Unauthorized: Admin ID not found in context", nil))
		return
	}
	var adminID int
	switch v := userIDClaim.(type) {
	case uint:
		adminID = int(v)
	case float64:
		adminID = int(v)
	case int:
		adminID = v
	default:
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", "Invalid admin ID format in context", nil))
		return
	}

	// 5. Panggil service untuk melakukan pekerjaan berat
	updatedRequest, err := h.service.ProcessRequestStatus(uint(requestID), input.Status, adminID)
	if err != nil {
		// Tampilkan error yang jelas dari service (misal: "stok tidak cukup")
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	// 6. Kirim response sukses dengan data request yang sudah terupdate
	c.JSON(http.StatusOK, utils.ResponseHandler("success", "Request status updated successfully", updatedRequest))
}
