package handlers

import (
	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
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
	var request models.Request
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
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	request.RequesterID = int(userID.(uint))
	var student *models.Student
	student, err = h.service.GetStudentByUserID(request.RequesterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	log.Printf("Student OrganizationID: %d, OrganizationName: %s", student.OrganizationID, student.Organization.Name)

	if student.OrganizationID == nil || *student.OrganizationID == 0 {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Student is not assigned to any organization"})
    return
}

	var orgID int
if student.OrganizationID != nil {
    orgID = *student.OrganizationID
}
request.OrganizationID = orgID

if student.Organization != nil {
    request.OrganizationName = student.Organization.Name
}

request.Status = "Pending"
request.CreatedAt = time.Now()
request.UpdatedAt = time.Now()

if err := h.service.CreateRequest(&request); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}


	c.JSON(http.StatusCreated, gin.H{
		"message": "Request created successfully",
		"data":    request,
	})
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
