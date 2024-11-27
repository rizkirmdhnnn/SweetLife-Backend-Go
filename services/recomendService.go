package services

import (
	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
)

type RecomendationService interface {
	GetFoodRecomendations(userid string) ([]*dto.FoodRecomendation, error)
	GetExerciseRecomendations(userid string) ([]*dto.ExerciseRecommendations, error)
}

type recomendationService struct {
	recomendationRepo repositories.RecomendationRepo
}

func NewRecomendationService(recomendationRepo repositories.RecomendationRepo) RecomendationService {
	if recomendationRepo == nil {
		panic("recomendationRepo cannot be nil")
	}
	return &recomendationService{
		recomendationRepo: recomendationRepo,
	}
}

// GetRecomendations implements RecomendationService.
func (r *recomendationService) GetFoodRecomendations(userid string) ([]*dto.FoodRecomendation, error) {
	panic("unimplemented")
}

// GetExerciseRecomendations implements RecomendationService.
func (r *recomendationService) GetExerciseRecomendations(userid string) ([]*dto.ExerciseRecommendations, error) {
	panic("unimplemented")
}
