package dto

import (
	"time"
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
