package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/errors"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

type MiniGroceryHandler struct {
	miniGroceryService services.MiniGroceryService
}

func NewMiniGroceryHandler(miniGroceryService services.MiniGroceryService) *MiniGroceryHandler {
	return &MiniGroceryHandler{
		miniGroceryService: miniGroceryService,
	}
}

// GetMiniCourse
func (h *MiniGroceryHandler) GetMiniGrocery(c *gin.Context) {
	// Get query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Parse page parameter
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusBadRequest, "invalid page parameter", "page must be a valid integer")
		return
	}

	// Parse limit parameter
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusBadRequest, "invalid limit parameter", "limit must be a valid integer")
		return
	}

	// Get mini course with pagination
	miniGrocery, err := h.miniGroceryService.GetMiniGroceryWithPagination(page, limit)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to get mini course", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   miniGrocery,
	})
}
