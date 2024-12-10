package dto

import (
	"time"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
)

// update user request
type UpdateUserRequest struct {
	Email       string `form:"email" json:"email"`
	Password    string `form:"password" json:"password"`
	Name        string `form:"name" json:"name"`
	DateOfBirth string `form:"date_of_birth" json:"date_of_birth"`
	Gender      string `form:"gender" json:"gender"`
}

type FoodHistoryResponse struct {
	FoodHistory []FoodHistoryEntry `json:"food_history"`
}

type FoodHistoryEntry struct {
	Date          string              `json:"date"`
	TotalCalories float64             `json:"total_calories"`
	Entries       []FoodHistoryByDate `json:"entries"`
}

type FoodHistoryByDate struct {
	ID         int        `json:"id"`
	Date       *time.Time `json:"date,omitempty"`
	TotalUnits int        `json:"total_units"`
	FoodName   string     `json:"food_name"`
	Calories   float64    `json:"calories"`
	Time       string     `json:"time"`
}

type DailyCaloriesRequest struct {
	Gender        string               `json:"gender"`
	Weight        float64              `json:"weight"`
	Height        float64              `json:"height"`
	Age           int                  `json:"age"`
	ActivityLevel models.ActivityLevel `json:"activity_level"`
}

type Satisfication string

const (
	UNDER Satisfication = "UNDER"
	PASS  Satisfication = "PASS"
	OVER  Satisfication = "OVER"
)

type DailyProgessPerItem struct {
	Current       float64       `json:"current"`
	Percent       int           `json:"percent"`
	Satisfication Satisfication `json:"satisfication"`
	Target        float64       `json:"target"`
}

type DailyProgress struct {
	Calories DailyProgessPerItem `json:"calories"`
	Carbs    DailyProgessPerItem `json:"carbs"`
	Sugar    DailyProgessPerItem `json:"sugar"`
}

type DailyProgressResponse struct {
	Progress DailyProgress `json:"dailyProgress"`
	Status   struct {
		Message       string        `json:"message"`
		Satisfication Satisfication `json:"satisfication"`
	} `json:"status"`
	User UserRespStruct `json:"user"`
}

type UserRespStruct struct {
	Name         string               `json:"name"`
	DiabetesType *models.DiabeticType `json:"diabetesType,omitempty"`
	Diabetes     bool                 `json:"diabetes"`
}

type DailyNutrition struct {
	TotalCalories float64 `json:"total_calories"`
	TotalCarbs    float64 `json:"total_carbs"`
	TotalSugar    float64 `json:"total_sugar"`
}
