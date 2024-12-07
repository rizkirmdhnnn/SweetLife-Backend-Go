package services

import (
	"context"
	"encoding/json"
	"errors"
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
	SearchFood(req *dto.FindFoodRequest) (*models.ScanFood, error)
	SaveFood(req *dto.SaveFoodRequest, userId string) error
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
			Name:         nutrition.Name,
			Unit:         total,
			Protein:      nutrition.Protein * float64(total),
			Sugar:        nutrition.Sugar * float64(total),
			Carbohydrate: nutrition.Carbohydrates * float64(total),
			Fat:          nutrition.Fat * float64(total),
		})
	}

	return response, nil
}

// Helper functions
func findFoodByName(foods []models.ScanFood, name string) (*models.ScanFood, error) {
	for _, food := range foods {
		if strings.EqualFold(*food.NameIndo, name) {
			return &food, nil
		}
	}
	return nil, fmt.Errorf("food %s not found", name)
}

// SearchFoodByName implements ScanFoodService.
func (s *scanFoodService) SearchFood(req *dto.FindFoodRequest) (*models.ScanFood, error) {
	// 1. find food by name from database where name = name and weight = weight
	food, err := s.scanRepo.SearchFoodFromDB(req.Name)
	if err != nil {
		// 2. If food not found, call ML service to scrape food data
		foodFromAPI, err := s.scanRepo.SearchFoodAPI(req.Name)
		if err != nil {
			return nil, err
		}

		// 4. Save food data to database
		foodData := models.Food{
			Name: foodFromAPI.Name,
		}
		if err := s.scanRepo.CreateFood(&foodData); err != nil {
			return nil, err
		}

		// 5. Save food nutrition data to database
		nutritions := models.FoodNutrition{
			FoodID:        foodData.ID,
			Calories:      foodFromAPI.Calories,
			Sugar:         foodFromAPI.Sugar,
			Fat:           foodFromAPI.Fat,
			Carbohydrates: foodFromAPI.Carbs,
			Proteins:      foodFromAPI.Protein,
			Weight:        foodFromAPI.Weight,
		}
		if err := s.scanRepo.CreateFoodNutrition(&nutritions); err != nil {
			return nil, err
		}

		nutritions = helper.CalculateNutrients(nutritions.Weight, &nutritions)
		food = &models.FoodWithNutritions{
			Food:      foodData,
			Nutrition: nutritions,
		}

	}

	// calculate nutrition data
	food.Nutrition = helper.CalculateNutrients(req.Weight, &food.Nutrition)

	// return food data
	data := &models.ScanFood{
		Name:          food.Food.Name,
		Calories:      food.Nutrition.Calories,
		Protein:       food.Nutrition.Proteins,
		Sugar:         food.Nutrition.Sugar,
		Carbohydrates: food.Nutrition.Carbohydrates,
		Fat:           food.Nutrition.Fat,
		Weight:        food.Nutrition.Weight,
	}

	return data, nil
}

// SaveFood implements ScanFoodService.
func (s *scanFoodService) SaveFood(req *dto.SaveFoodRequest, userId string) error {
	var histories []models.UserFoodHistory

	// 1. Proses makanan hasil scan
	if len(req.Scan) > 0 {
		// Ambil ID makanan hasil scan
		foodMap, err := s.scanRepo.GetFoodIDs(&req.Scan)
		if err != nil {
			return err
		}

		// Buat data user food history dari hasil scan
		for _, food := range req.Scan {
			histories = append(histories, models.UserFoodHistory{
				UserID: userId,
				FoodID: foodMap[food.Name],
				Unit:   food.Unit,
				Weight: nil, // Berat tidak digunakan di data hasil scan
			})
		}
	}

	// 2. Proses makanan tambahan
	if len(req.Additionall) > 0 {
		// Konversi data tambahan ke format ScanFood
		var additionalScanFoods []dto.ScanFood
		for _, food := range req.Additionall {
			additionalScanFoods = append(additionalScanFoods, dto.ScanFood{
				Name: food.Name,
			})
		}

		// Ambil ID makanan tambahan
		foodMap, err := s.scanRepo.GetFoodIDs(&additionalScanFoods)
		if err != nil {
			return err
		}

		// Buat data user food history dari makanan tambahan
		for _, food := range req.Additionall {
			histories = append(histories, models.UserFoodHistory{
				UserID: userId,
				FoodID: foodMap[food.Name],
				Unit:   1,
				Weight: &food.Weight,
			})
		}
	}

	// 3. Simpan semua data user food history ke database
	if len(histories) > 0 {
		if err := s.scanRepo.SaveUserFoodHistory(&histories); err != nil {
			if strings.Contains(err.Error(), "violates foreign key constraint") {
				return errors.New("food not found")
			}
			return err
		}
	}

	return nil
}
