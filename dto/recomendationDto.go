package dto

type RecomendationDto struct {
	FoodRecomendation       FoodRecomendation       `json:"food_recomendation"`
	ExerciseRecommendations ExerciseRecommendations `json:"exercise_recommendations"`
}

type FoodRecomendation struct {
	Name    string               `json:"name"`
	Details RecomendationDetails `json:"details"`
	Image   RecomendationImage   `json:"image"`
}

type RecomendationDetails struct {
	Calories     float64 `json:"calories"`
	Carbohydrate float64 `json:"carbohydrate"`
	Fat          float64 `json:"fat"`
	Proteins     float64 `json:"proteins"`
}

type RecomendationImage struct {
	URL string `json:"url"`
}

type ExerciseRecommendations struct {
	Name string `json:"name"`
}

type DiabetesPredictionRequest struct {
	Age            int     `json:"age"`
	HeartDisease   bool    `json:"heart_disease"`
	SmokingHistory string  `json:"smoking_history"`
	BMI            float64 `json:"bmi"`
	Gender         string  `json:"gender"`
}

type DiabetesPredictionClientResp struct {
	Percentage float64 `json:"percentage"`
	Note       string  `json:"note"`
}

type FoodRecomendationClientResp struct {
	Calories     float64 `json:"calories"`
	Carbohydrate float64 `json:"carbohydrate"`
	Fat          float64 `json:"fat"`
	Image        string  `json:"image"`
	Name         string  `json:"name"`
	Proteins     float64 `json:"proteins"`
}
