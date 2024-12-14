package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/errors"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

// UserHandler is a struct to handle user request
type UserHandler struct {
	userService    services.UserService
	storageService services.StorageBucketService
}

// NewUserHandler is a constructor to create UserHandler instance
func NewUserHandler(userService services.UserService, storageService services.StorageBucketService) *UserHandler {
	if userService == nil {
		panic("userService cannot be nil")
	}
	return &UserHandler{
		userService:    userService,
		storageService: storageService,
	}
}

// UpdateProfile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// get userID from context
	userID := c.GetString("userID")

	/// get data from request
	var req dto.UpdateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	photoProfile, _ := c.FormFile("profile_picture")
	// check if format not jpg or png
	if photoProfile != nil {
		if photoProfile.Header.Get("Content-Type") != "image/jpeg" && photoProfile.Header.Get("Content-Type") != "image/png" {
			errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", "Profile picture must be in jpg or png format")
			return
		}
	}

	// call service to update user
	err := h.userService.UpdateProfile(userID, photoProfile, &req)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to update user", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "action success",
	})
}

// UpdateProfile
func (h *UserHandler) UpdatePhotoProfile(c *gin.Context) {
	// get userID from context
	userID := c.GetString("userID")

	photoProfile, _ := c.FormFile("profile_picture")
	// check if format not jpg or png
	if photoProfile != nil {
		if photoProfile.Header.Get("Content-Type") != "image/jpeg" && photoProfile.Header.Get("Content-Type") != "image/png" {
			errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", "Profile picture must be in jpg or png format")
			return
		}
	}

	// call service to update user
	err := h.userService.UpdatePhotoProfile(userID, photoProfile)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to update user", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "action success",
	})
}

// GetProfile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// get userID from context
	userID := c.GetString("userID")

	// call service to get user profile
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to get user profile", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   user,
	})
}

// GetHistory
func (h *UserHandler) GetHistory(c *gin.Context) {
	// get userID from context
	userID := c.GetString("userID")

	foodHistory, err := h.userService.GetFoodHistoryWithPagination(userID)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to get user history", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   foodHistory,
	})
}

// GetDashboard
func (h *UserHandler) GetDashboard(c *gin.Context) {
	// get userID from context
	userID := c.GetString("userID")

	dashboard, err := h.userService.GetDashboard(userID)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to get user dashboard", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   dashboard,
	})
}
