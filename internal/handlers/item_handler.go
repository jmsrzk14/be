package handlers

import (
	"math"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"bem_be/internal/models"
	"bem_be/internal/services"
	"bem_be/internal/utils"

	"github.com/gin-gonic/gin"
)

// ItemHandler handles HTTP requests related to items
type ItemHandler struct {
	service *services.ItemService
}

// NewItemHandler creates a new item handler
func NewItemHandler(db *gorm.DB) *ItemHandler {
	return &ItemHandler{
		service: services.NewItemService(db),
	}
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	var item models.Item

	item.Name = c.PostForm("name")
	amountStr := c.PostForm("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid amount"})
		return
	}
	item.Amount = amount

	// ambil file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logo file is required"})
		return
	}

	// kirim ke service
	if err := h.service.CreateItem(&item, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Item created successfully",
		"data":    item,
	})
}

func (h *ItemHandler) GetAllItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("name") // pencarian pakai param ?name=

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	items, total, err := h.service.GetAllItems(perPage, offset, search)
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
	}

	response := utils.MetadataFormatResponse(
		"success",
		"Berhasil list mendapatkan data associations",
		metadata,
		items,
	)

	c.JSON(http.StatusOK, response)
}

func (h *ItemHandler) GetAllItemsGuest(c *gin.Context) {
	// ambil semua data tanpa limit & offset
	items, err := h.service.GetAllItemsGuest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseHandler("error", err.Error(), nil))
		return
	}

	// langsung response tanpa metadata
	response := utils.ResponseHandler(
		"success",
		"Berhasil mendapatkan data",
		items,
	)

	c.JSON(http.StatusOK, response)
}

// GetItemByID returns a item by ID
func (h *ItemHandler) GetItemByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stats := c.Query("stats")
	var result interface{}

	if stats == "true" {
		result, err = h.service.GetItemWithStats(uint(id))
	} else {
		result, err = h.service.GetItemByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Item retrieved successfully",
		"data":    result,
	})
}

func (h *ItemHandler) UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	item.ID = uint(id)

	if err := h.service.UpdateItem(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Item updated successfully",
		"data":    item,
	})
}

// DeleteItem deletes a item
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteItem(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Item deleted successfully",
	})
}
