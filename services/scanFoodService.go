package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	helper "github.com/rizkirmdhnnn/sweetlife-backend-go/helpers"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
)

type ScanFoodService interface {
	ScanFood(file *multipart.FileHeader) (*dto.ScanFoodResponse, error)
}

type scanFoodService struct {
	scanRepo    repositories.ScanFoodRepository
	storegeRepo repositories.StorageBucketRepository
}

func NewScanFoodService(scanRepo repositories.ScanFoodRepository, storageRepo repositories.StorageBucketRepository) ScanFoodService {
	return &scanFoodService{
		scanRepo:    scanRepo,
		storegeRepo: storageRepo,
	}
}

func (s *scanFoodService) ScanFood(file *multipart.FileHeader) (*dto.ScanFoodResponse, error) {
	// Generate unique file name
	fileName := helper.GenerateFileName(filepath.Ext(file.Filename))
	uploadPath := "website/scan-food/"
	image, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer image.Close()

	// Upload file to storage
	url, err := s.storegeRepo.UploadFile(context.Background(), config.ENV.STORAGE_BUCKET, uploadPath+fileName, image)
	if err != nil {
		return nil, err
	}

	fmt.Println("URL: ", url)

	// Call ML API to scan food
	scanFoodResponse, err := s.scanRepo.ScanFood(url)
	if err != nil {
		return nil, err
	}

	// Open file with nutrition data
	fileNutrition, err := os.Open("data/nutritions.json")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer fileNutrition.Close()

	// Parse the JSON data
	var foodListClientResp []models.ScanFood
	if err := json.NewDecoder(fileNutrition).Decode(&foodListClientResp); err != nil {
		log.Fatalf("failed to decode JSON: %v", err)
	}

	// Group food items by name and calculate total amount
	foodTotals := make(map[string]int)
	for _, food := range scanFoodResponse.Objects {
		foodTotals[food.Name] += food.Unit
	}

	// Create response
	response := &dto.ScanFoodResponse{
		IsDetected: len(foodTotals) > 0,
		FoodList:   []dto.FoodList{},
	}

	// Match grouped foods with nutrition data
	for name, total := range foodTotals {
		// Find nutrition by name
		nutrition, err := findFoodByName(foodListClientResp, name)
		if err != nil {
			fmt.Printf("Food not found: %s\n", name)
			continue // Skip if food not found
		}

		// Multiply nutrition values by total
		response.FoodList = append(response.FoodList, dto.FoodList{
			Name:         name,
			Amount:       total,
			Weight:       nutrition.Berat * float64(total),
			Calories:     nutrition.Kalori * float64(total),
			Protein:      nutrition.Protein * float64(total),
			Sugar:        nutrition.Gula * float64(total),
			Carbohydrate: nutrition.Karbohidrat * float64(total),
			Fat:          nutrition.Lemak * float64(total),
		})
	}

	return response, nil
}

// Helper functions
func findFoodByName(foods []models.ScanFood, name string) (*models.ScanFood, error) {
	for _, food := range foods {
		if strings.EqualFold(food.Name, name) {
			return &food, nil
		}
	}
	return nil, fmt.Errorf("food %s not found", name)
}
