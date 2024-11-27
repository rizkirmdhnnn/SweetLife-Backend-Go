package services

import (
	"fmt"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
)

type RecomendationService interface {
	GetFoodRecomendations(userid string) ([]*dto.FoodRecomendation, error)
	GetExerciseRecomendations(userid string) ([]*dto.ExerciseRecommendations, error)
}

type recomendationService struct {
	recomendationRepo repositories.RecomendationRepo
	healthRepo        repositories.HealthProfileRepository
}

func NewRecomendationService(recomendationRepo repositories.RecomendationRepo, healthRepo repositories.HealthProfileRepository) RecomendationService {
	if recomendationRepo == nil {
		panic("recomendationRepo cannot be nil")
	}
	return &recomendationService{
		recomendationRepo: recomendationRepo,
		healthRepo:        healthRepo,
	}
}

// GetRecomendations implements RecomendationService.
func (r *recomendationService) GetFoodRecomendations(userid string) ([]*dto.FoodRecomendation, error) {

	//get risk percentage
	healthProfile, err := r.healthRepo.GetRiskAssessmentByUserID(userid)
	if err != nil {
		return nil, err
	}

	//get recomendation
	foodRecomendationClientResp, err := r.recomendationRepo.GetFoodRecomendations(float32(healthProfile.RiskScore))
	if err != nil {
		return nil, err
	}

	var foodRecomendations []*dto.FoodRecomendation
	for _, foodList := range foodRecomendationClientResp.FoodRecomendation {
		for _, food := range foodList {
			foodRec := dto.FoodRecomendation{
				Name: food.Name,
				Details: dto.RecomendationDetails{
					Carbohydrate: fmt.Sprintf("%.2f g", food.Carbohydrate),
					Calories:     fmt.Sprintf("%.2f kcal", food.Calories),
					Fat:          fmt.Sprintf("%.2f g", food.Fat),
					Proteins:     fmt.Sprintf("%.2f g", food.Proteins),
				},
				Image: dto.RecomendationImage{
					URL: food.Image,
				},
			}
			foodRecomendations = append(foodRecomendations, &foodRec)
		}
	}

	return foodRecomendations, nil
}

// GetExerciseRecomendations implements RecomendationService.
func (r *recomendationService) GetExerciseRecomendations(userid string) ([]*dto.ExerciseRecommendations, error) {
	panic("unimplemented")
}
