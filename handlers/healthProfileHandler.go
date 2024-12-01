package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/errors"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

type HealthProfileHandler struct {
	healthProfileService services.HealthProfileService
}

func NewHealthProfileHandler(healthProfileService services.HealthProfileService) *HealthProfileHandler {
	return &HealthProfileHandler{
		healthProfileService: healthProfileService,
	}
}

// CreateHealthProfile is a handler to create health profile
func (h *HealthProfileHandler) CreateHealthProfile(c *gin.Context) {
	// get userID from context
	userID := c.GetString("userID")

	// get data from request
	var req dto.HealthProfileDto
	if err := c.ShouldBind(&req); err != nil {
		errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// set userID
	req.UserID = userID

	// call service to create health profile
	if err := h.healthProfileService.CreateHealthProfile(&req); err != nil {
		errors.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create health profile", err.Error())
		return
	}
	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "action success",
	})
}

// GetHealthProfile is a handler to get health profile
func (h *HealthProfileHandler) GetHealthProfile(c *gin.Context) {
	// get userID from context
	userID := c.GetString("userID")

	// call service to get health profile
	profile, err := h.healthProfileService.GetHealthProfile(userID)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get health profile", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "action success",
		"data":    profile,
	})
}
