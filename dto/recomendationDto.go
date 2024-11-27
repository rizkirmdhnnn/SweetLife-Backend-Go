package dto

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
	Diabetes          bool               `json:"diabetes"`
	FoodRecomendation [][]FoodClientResp `json:"food_recommendation"`
}

type FoodClientResp struct {
	Calories     float64 `json:"calories"`
	Carbohydrate float64 `json:"carbohydrate"`
	Fat          float64 `json:"fat"`
	Image        string  `json:"image"`
	Name         string  `json:"name"`
	Proteins     float64 `json:"proteins"`
}

type RecomendationDto struct {
	FoodRecomendation       []*FoodRecomendation    `json:"food_recommendation"`
	ExerciseRecommendations *ExerciseRecommendation `json:"exercise_recommendations"`
}

type ExerciseRecommendation struct {
	CaloriesBurned   float64         `json:"calories_burned"`
	ExerciseDuration float64         `json:"exercise_duration"`
	ExerciseList     []*ExerciseList `json:"exercise_list"`
}

type ExerciseList struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Image string `json:"image"`
}

type ExerciseRecommendationClientResp struct {
	CaloriesBurned     float64  `json:"calories_burned"`
	ExerciseCategories []string `json:"exercise_categories"`
	ExerciseDuration   float64  `json:"exercise_duration"`
}

type ExerciseRequest struct {
	Gender   string  `json:"gender"`
	Age      int     `json:"age"`
	Height   float64 `json:"height"`
	Diabetes bool    `json:"diabetes"`
	Bmi      float64 `json:"bmi"`
}

type FoodRecomendation struct {
	Name    string               `json:"name"`
	Details RecomendationDetails `json:"details"`
	Image   string               `json:"image"`
}

type RecomendationDetails struct {
	Calories     string `json:"calories"`
	Carbohydrate string `json:"carbohydrate"`
	Fat          string `json:"fat"`
	Proteins     string `json:"proteins"`
}
