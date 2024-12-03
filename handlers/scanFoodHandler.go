package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/errors"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

type ScanFoodHandler struct {
	scanFoodService services.ScanFoodService
}

func NewScanFoodHandler(scanFoodService services.ScanFoodService) *ScanFoodHandler {
	return &ScanFoodHandler{
		scanFoodService: scanFoodService,
	}
}

// ScanFood is a handler to scan food
func (s *ScanFoodHandler) ScanFood(c *gin.Context) {
	// get file from request
	file, err := c.FormFile("image")

	if file != nil {
		if file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/png" {
			errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", "Profile picture must be in jpg or png format")
			return
		}
	}

	if err != nil {
		errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// call service to scan food
	scanFoodResponse, err := s.scanFoodService.ScanFood(file)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusInternalServerError, "Failed to scan food", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status":    true,
		"message":   "Food detected successfully",
		"food_list": scanFoodResponse.FoodList,
	})
}
