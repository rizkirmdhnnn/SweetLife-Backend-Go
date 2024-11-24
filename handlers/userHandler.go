package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/errors"
	helper "github.com/rizkirmdhnnn/sweetlife-backend-go/helpers"
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

	// Process profile picture (optional)
	file, err := c.FormFile("profile_picture")
	var profilePictureURL string

	// check if profile picture is uploaded
	if err == nil {
		// Generate a unique file name
		fileName := helper.GenerateFileName(filepath.Ext(file.Filename))
		uploadPath := "website/photo-profile/"

		// save to cloud storage
		profilePictureURL, err = h.storageService.UploadFile(file, uploadPath, fileName)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload profile picture", "details": err.Error()})
			return
		}
	}

	// call service to update user
	err = h.userService.UpdateProfile(userID, profilePictureURL, &req)
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
