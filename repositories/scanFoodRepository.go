package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"gorm.io/gorm"
)

type ScanFoodRepository interface {
	ScanFood(image string) (*dto.ScanFoodClientResp, error)
	SearchFoodFromDB(name string) (*models.FoodWithNutritions, error)
	SearchFoodFromML(name string) (*dto.FindFoodClientResp, error)

	CreateFood(food *models.Food) error
	CreateFoodNutrition(foodNutrition *models.FoodNutrition) error
}

type scanFoodRepository struct {
	httpClient *http.Client
	db         *gorm.DB
}

// CreateFood implements ScanFoodRepository.

func NewScanFoodRepository(httpClient *http.Client, db *gorm.DB) ScanFoodRepository {
	return &scanFoodRepository{
		httpClient: httpClient,
		db:         db,
	}
}

func (s *scanFoodRepository) ScanFood(image string) (*dto.ScanFoodClientResp, error) {
	data := map[string]interface{}{
		"image": image,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Post("https://ml.sweetlife.my.id/scan-food", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var scanFoodResponse dto.ScanFoodClientResp
	err = json.NewDecoder(resp.Body).Decode(&scanFoodResponse)
	if err != nil {
		return nil, err
	}

	return &scanFoodResponse, nil
}

// SearchFoodFromML implements ScanFoodRepository.
func (s *scanFoodRepository) SearchFoodFromML(name string) (*dto.FindFoodClientResp, error) {
	data := map[string]interface{}{
		"name":   name,
		"weight": 100,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Post("https://ml.sweetlife.my.id/food_nutritions", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var findFoodResp dto.FindFoodClientResp
	err = json.NewDecoder(resp.Body).Decode(&findFoodResp)
	if err != nil {
		return nil, err
	}

	fmt.Print(findFoodResp)

	return &findFoodResp, nil
}

// SearchFood implements ScanFoodRepository.
func (s *scanFoodRepository) SearchFoodFromDB(name string) (*models.FoodWithNutritions, error) {
	var foodWithNutrition models.FoodWithNutritions

	food := s.db.Where("foods.name = ?", name).First(&foodWithNutrition.Food)
	if food.Error != nil {
		return nil, food.Error
	}

	nutritions := s.db.Where("food_id = ?", foodWithNutrition.Food.ID).First(&foodWithNutrition.Nutrition)
	if nutritions.Error != nil {
		return nil, nutritions.Error
	}

	return &foodWithNutrition, nil
}

func (s *scanFoodRepository) CreateFood(food *models.Food) error {
	result := s.db.Create(food)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// CreateFoodNutrition implements ScanFoodRepository.
func (s *scanFoodRepository) CreateFoodNutrition(foodNutrition *models.FoodNutrition) error {
	result := s.db.Create(foodNutrition)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
