package handlers

import (
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

type RecomendationHandler struct {
	recomendationService services.RecomendationService
}

func NewRecomendationHandler(recomendationService services.RecomendationService) *RecomendationHandler {
	return &RecomendationHandler{
		recomendationService: recomendationService,
	}
}

func (r *RecomendationHandler) GetRecomendations() {
	//TODO: implement get recomendations
}
