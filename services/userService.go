package services

import (
	"context"
	"errors"
	"fmt"
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
	GetFoodHistoryWithPagination(userID string) (*dto.FoodHistoryResponse, error)

	// GetDashboard
	GetDashboard(userID string) (*dto.DailyProgressResponse, error)
}

type userService struct {
	userRepo    repositories.UserRepository
	authRepo    repositories.AuthRepository
	storageRepo repositories.StorageBucketRepository
	healthRepo  repositories.HealthProfileRepository
}

func NewUserService(userRepo repositories.UserRepository, authRepo repositories.AuthRepository, storageRepo repositories.StorageBucketRepository, healthRepo repositories.HealthProfileRepository) UserService {
	if userRepo == nil {
		panic("userRepo cannot be nil")
	}

	if authRepo == nil {
		panic("authRepo cannot be nil")
	}
	if storageRepo == nil {
		panic("storageRepo cannot be nil")
	}
	if healthRepo == nil {
		panic("healthRepo cannot be nil")
	}

	return &userService{
		userRepo:    userRepo,
		authRepo:    authRepo,
		storageRepo: storageRepo,
		healthRepo:  healthRepo,
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

	// jika dateOfBirth itu tahun sekarang
	if date.Year() == time.Now().Year() {
		return errors.New("invalid date of birth")
	}

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

// TODO: ini untuk makanan yang dimasukin manual total kalorinya belum bener
// harusnya total kalorinya berdasarkan weightnya, klo yang sekarang masih berdasarkan total unit yang default weightnya 100
func (s *userService) GetFoodHistoryWithPagination(userID string) (*dto.FoodHistoryResponse, error) {
	// Get food history
	foodHistory, err := s.userRepo.GetFoodHistory(userID)
	if err != nil {
		return nil, err
	}

	// create map to store food history by date
	foodHistoryMap := make(map[string]*dto.FoodHistoryEntry)

	// Iterate over food history
	for _, entry := range foodHistory {
		// Format date
		formattedDate := entry.Date.Format("2006-01-02")

		// If date not exists in map, create new entry
		if _, exists := foodHistoryMap[formattedDate]; !exists {
			foodHistoryMap[formattedDate] = &dto.FoodHistoryEntry{
				Date:          formattedDate,
				TotalCalories: 0,
				Entries:       []dto.FoodHistoryByDate{},
			}
		}

		// Add calories to total calories
		foodHistoryMap[formattedDate].TotalCalories += entry.Calories

		// Add entry to entries
		foodHistoryMap[formattedDate].Entries = append(foodHistoryMap[formattedDate].Entries, dto.FoodHistoryByDate{
			ID:         entry.ID,
			FoodName:   entry.FoodName,
			Calories:   entry.Calories,
			Time:       entry.Time,
			TotalUnits: entry.TotalUnits,
		})
	}

	// Convert map to slice
	var finalFoodHistory []dto.FoodHistoryEntry
	for _, entry := range foodHistoryMap {
		finalFoodHistory = append(finalFoodHistory, *entry)
	}

	// Create response
	response := &dto.FoodHistoryResponse{
		FoodHistory: finalFoodHistory,
	}

	return response, nil
}

// GetDashboard implements UserService.
func (u *userService) GetDashboard(userID string) (*dto.DailyProgressResponse, error) {
	// Get user
	user, err := u.authRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	// Get user health profile
	profile, err := u.healthRepo.GetHealthProfileByUserID(userID)
	if err != nil {
		return nil, err
	}

	var userResp dto.UserRespStruct
	userResp.Name = user.Name
	userResp.Diabetes = profile.IsDiabetic

	if profile.IsDiabetic {
		diabetesDetails, err := u.healthRepo.GetDiabetesDetailsByProfileID(fmt.Sprintf("%d", profile.ID))
		if err != nil {
			return nil, err
		}
		userResp.DiabetesType = &diabetesDetails.DiabeticType
	}

	// Get user daily calories
	dailyCalories, err := helper.CalculateDailyCalories(dto.DailyCaloriesRequest{
		Height:        profile.Height,
		Weight:        profile.Weight,
		Gender:        user.Gender,
		Age:           user.Age,
		ActivityLevel: profile.ActivityLevel,
	})
	if err != nil {
		return nil, err
	}

	// Get user daily sugar
	dailySugar := helper.CalculateDailySugar(dailyCalories, profile.IsDiabetic)

	// Get user daily nutrition
	dailyCarb := helper.CalculateDialyCarbs(dailyCalories)

	// Get user daily nutrition
	dailyNutrition, err := u.userRepo.GetDailyNutrition(userID)
	if err != nil {
		return nil, err
	}

	// Dalam fungsi GetDashboard
	caloriesSatisfication := helper.DetermineSatisfication(
		float64(dailyNutrition.TotalCalories),
		float64(dailyCalories),
	)

	carbsSatisfication := helper.DetermineSatisfication(
		float64(dailyNutrition.TotalCarbs),
		float64(dailyCarb),
	)

	sugarSatisfication := helper.DetermineSatisfication(
		float64(dailyNutrition.TotalSugar),
		float64(dailySugar),
	)

	// Get user daily progress
	dailyProgress := &dto.DailyProgressResponse{
		Progress: dto.DailyProgress{
			Calories: dto.DailyProgessPerItem{
				Current:       dailyNutrition.TotalCalories,
				Percent:       int(float64(dailyNutrition.TotalCalories) / float64(dailyCalories) * 100),
				Satisfication: caloriesSatisfication,
				Target:        dailyCalories,
			},
			Carbs: dto.DailyProgessPerItem{
				Current:       dailyNutrition.TotalCarbs,
				Percent:       int(float64(dailyNutrition.TotalCarbs) / float64(dailyCarb) * 100),
				Satisfication: carbsSatisfication,
				Target:        dailyCarb,
			},
			Sugar: dto.DailyProgessPerItem{
				Current:       dailyNutrition.TotalSugar,
				Percent:       int(float64(dailyNutrition.TotalSugar) / float64(dailySugar) * 100),
				Satisfication: sugarSatisfication,
				Target:        dailySugar,
			},
		},
		Status: struct {
			Message       string            `json:"message"`
			Satisfication dto.Satisfication `json:"satisfication"`
		}{
			Message: helper.DetermineOverallMessage(caloriesSatisfication, carbsSatisfication, sugarSatisfication),
			Satisfication: helper.DetermineOverallSatisfication(
				dto.DailyProgessPerItem{
					Satisfication: caloriesSatisfication,
				},
				dto.DailyProgessPerItem{
					Satisfication: carbsSatisfication,
				},
				dto.DailyProgessPerItem{
					Satisfication: sugarSatisfication,
				},
			),
		},
		User: userResp,
	}

	return dailyProgress, nil
}
