package services

import (
	"context"
	"errors"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	helper "github.com/rizkirmdhnnn/sweetlife-backend-go/helpers"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
)

type UserService interface {
	// UpdateProfile
	UpdateProfile(id string, photoProfile *multipart.FileHeader, req *dto.UpdateUserRequest) error
	// Profile
	GetProfile(id string) (*dto.UserResponse, error)
}

type userService struct {
	userRepo    repositories.UserRepository
	authRepo    repositories.AuthRepository
	storageRepo repositories.StorageBucketRepository
}

func NewUserService(userRepo repositories.UserRepository, authRepo repositories.AuthRepository, storageRepo repositories.StorageBucketRepository) UserService {
	if userRepo == nil {
		panic("userRepo cannot be nil")
	}

	if authRepo == nil {
		panic("authRepo cannot be nil")
	}
	if storageRepo == nil {
		panic("storageRepo cannot be nil")
	}

	return &userService{
		userRepo:    userRepo,
		authRepo:    authRepo,
		storageRepo: storageRepo,
	}
}

// UpdateProfile implements UserService.
func (u *userService) UpdateProfile(id string, photoProfile *multipart.FileHeader, req *dto.UpdateUserRequest) error {
	// Get user by ID
	user, err := u.authRepo.GetUserById(id)
	if err != nil {
		return err
	}

	// Jika photoProfile disertakan, maka hapus file lama jika ada, lalu upload file baru
	if photoProfile != nil {
		// Hapus file lama jika user memiliki foto profile
		if user.ImageUrl != "" {
			url, err := url.Parse(user.ImageUrl)
			if err != nil {
				return errors.New("invalid file URL")
			}

			// Extract file path
			filePath := strings.TrimPrefix(url.Path, "/sweetlife-go/") // Ganti sesuai root bucket Anda
			if filePath == "" {
				return errors.New("invalid file path")
			}

			// Hapus file lama
			if err := u.storageRepo.DeleteFile(context.Background(), config.ENV.STORAGE_BUCKET, filePath); err != nil {
				return err
			}
		}

		// Upload foto profile baru
		fileName := helper.GenerateFileName(filepath.Ext(photoProfile.Filename))
		uploadPath := "website/photo-profile/"
		file, err := photoProfile.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		urlPhotoProfile, err := u.storageRepo.UploadFile(context.Background(), config.ENV.STORAGE_BUCKET, uploadPath+fileName, file)
		if err != nil {
			return err
		}

		// Set URL foto profile baru ke user
		user.ImageUrl = urlPhotoProfile
	}

	// Parse date
	date, _ := helper.ParsedDate(req.DateOfBirth)

	// Update data user
	user.Name = req.Name
	user.DateOfBirth = date
	user.Gender = req.Gender
	user.Updated_at = time.Now()
	user.Age, _ = helper.CalculateAge(date.Format("2006-01-02"))

	// Save user
	err = u.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

// GetProfile implements UserService.
func (u *userService) GetProfile(id string) (*dto.UserResponse, error) {
	// get user by id
	user, err := u.authRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	// create response
	res := dto.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		DateOfBirth:  user.DateOfBirth.Format("2006-01-02"),
		Gender:       user.Gender,
		PhotoProfile: &user.ImageUrl,
	}

	return &res, nil
}
