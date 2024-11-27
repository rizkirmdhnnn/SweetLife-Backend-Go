package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/errors"
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

func (r *RecomendationHandler) GetRecomendations(c *gin.Context) {

	// get id from context
	userID := c.GetString("userID")

	// get recomendations
	foodRecomendations, err := r.recomendationService.GetFoodRecomendations(userID)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get recomendations", err.Error())
		return
	}

	// get exercise recomendations
	exerciseRecomendations, err := r.recomendationService.GetExerciseRecomendations(userID)
	if err != nil {
		errors.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get recomendations", err.Error())
		return
	}

	resp := dto.RecomendationDto{
		FoodRecomendation:       foodRecomendations,
		ExerciseRecommendations: exerciseRecomendations,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   resp,
	})

}
