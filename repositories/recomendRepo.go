package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
)

type RecomendationRepo interface {
	GetFoodRecomendations(diabetPercentage float32) (*dto.FoodRecomendationClientResp, error)
	//TODO: Get Exercice recomendation
	// GetExerciceRecomendations(diabetPercentage float32) ([]*dto.ExerciseRecommendationClientResp, error)
	DiabetesPrediction(data *dto.DiabetesPredictionRequest) (*dto.DiabetesPredictionClientResp, error)
}

type recomendationRepo struct {
	httpClient *http.Client
}

func NewRecomendationRepo(httpClient *http.Client) RecomendationRepo {
	if httpClient == nil {
		panic("httpClient cannot be nil")
	}
	return &recomendationRepo{
		httpClient: httpClient,
	}
}

func (r *recomendationRepo) GetFoodRecomendations(diabetPercentage float32) (*dto.FoodRecomendationClientResp, error) {
	data := map[string]interface{}{
		"diabetes_percentage": diabetPercentage,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, err
	}

	fmt.Print("Request Body: ")

	resp, err := r.httpClient.Post("https://ml.sweetlife.my.id//food_recommendation", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	fmt.Println("Response Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Decode JSON ke dalam struktur Go
	var recomendation dto.FoodRecomendationClientResp
	err = json.Unmarshal(body, &recomendation)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	fmt.Println("Recomendation:", recomendation)

	return &recomendation, nil
}

func (r *recomendationRepo) DiabetesPrediction(data *dto.DiabetesPredictionRequest) (*dto.DiabetesPredictionClientResp, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, err
	}

	resp, err := r.httpClient.Post("https://ml.sweetlife.my.id/diabetes_predict", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var prediction dto.DiabetesPredictionClientResp
	err = json.Unmarshal(body, &prediction)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	return &prediction, nil
}
