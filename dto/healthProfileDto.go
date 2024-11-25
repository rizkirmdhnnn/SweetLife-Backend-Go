package dto

import "github.com/rizkirmdhnnn/sweetlife-backend-go/models"

// HealthProfileDto is a data transfer object for health profile
type HealthProfileDto struct {
	UserID          string               `json:"user_id"`
	Height          float64              `json:"height"`
	Weight          float64              `json:"weight"`
	IsDiabetic      bool                 `json:"is_diabetic"`
	IsSmoker        bool                 `json:"is_smoker"`
	HasHeartDisease bool                 `json:"has_heart_disease"`
	ActivityLevel   models.ActivityLevel `json:"activity_level"`

	// Diabetes details
	DiabeticType  models.DiabeticType `json:"diabetic_type"`
	InsulinLevel  float64             `json:"insulin_level"`
	BloodPressure uint                `json:"blood_pressure"`
}
