package repositories

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"gorm.io/gorm"
)

type ScanFoodRepository interface {
	ScanFood(image string) (*dto.ScanFoodClientResp, error)
	SearchFoodFromDB(name string) (*models.FoodWithNutritions, error)
	SearchFoodAPI(foodName string) (*dto.FoodNutritionResponse, error)

	GetFoodIDs(foodNames *[]dto.ScanFood) (map[string]uint, error)

	CreateFood(food *models.Food) error
	CreateFoodNutrition(foodNutrition *models.FoodNutrition) error

	SaveUserFoodHistory(food *[]models.UserFoodHistory) error
}

type scanFoodRepository struct {
	httpClient *http.Client
	db         *gorm.DB
	apiKey     string
}

// CreateFood implements ScanFoodRepository.
func NewScanFoodRepository(httpClient *http.Client, db *gorm.DB, apiKey string) ScanFoodRepository {
	return &scanFoodRepository{
		httpClient: httpClient,
		db:         db,
		apiKey:     apiKey,
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

// SearchFoodAPI implements ScanFoodRepository.
func (s *scanFoodRepository) SearchFoodAPI(foodName string) (*dto.FoodNutritionResponse, error) {
	url := fmt.Sprintf("https://api.nal.usda.gov/fdc/v1/foods/search?api_key=%s&query=%s", s.apiKey, strings.ReplaceAll(foodName, " ", "%20"))
	resp, err := s.httpClient.Get(url)
	if err != nil {
		log.Fatal("Failed to fetch data:", err)
	}

	var food dto.FoodNutritionResponseClient
	err = json.NewDecoder(resp.Body).Decode(&food)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	// if no food data found.
	if len(food.Foods) == 0 {
		return nil, errors.New("food not found")
	}

	// Extract the first food data.
	foodData := food.Foods[0]

	// Nutrient extraction using a map
	targetNutrients := map[string]*float64{
		"Energy":                      new(float64),
		"Protein":                     new(float64),
		"Total lipid (fat)":           new(float64),
		"Carbohydrate, by difference": new(float64),
		"Total Sugars":                new(float64),
	}

	// Extract the nutrient value.
	for _, nutrient := range foodData.FoodNutrients {
		if target, exists := targetNutrients[nutrient.NutrientName]; exists {
			*target = nutrient.Value
		}
	}

	// Create the response data.
	data := dto.FoodNutritionResponse{
		Name:     foodName,
		Calories: *targetNutrients["Energy"],
		Protein:  *targetNutrients["Protein"],
		Fat:      *targetNutrients["Total lipid (fat)"],
		Carbs:    *targetNutrients["Carbohydrate, by difference"],
		Sugar:    *targetNutrients["Total Sugars"],
		Weight:   100,
	}

	return &data, nil
}

// CreateFood implements ScanFoodRepository.
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

// GetFoodIDs implements ScanFoodRepository.
func (s *scanFoodRepository) GetFoodIDs(foodNames *[]dto.ScanFood) (map[string]uint, error) {
	var foods []models.Food
	var names []string

	// Ekstrak nama makanan dari slice dto.ScanFood
	for _, food := range *foodNames {
		names = append(names, food.Name)
	}

	// Query untuk mengambil ID berdasarkan nama makanan
	if err := s.db.Where("name IN ?", names).Find(&foods).Error; err != nil {
		return nil, err
	}

	// Buat map nama makanan ke ID
	foodMap := make(map[string]uint)
	for _, food := range foods {
		foodMap[food.Name] = food.ID
	}

	return foodMap, nil
}

// SaveUserFoodHistory implements ScanFoodRepository.
func (s *scanFoodRepository) SaveUserFoodHistory(food *[]models.UserFoodHistory) error {
	result := s.db.Create(food)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
