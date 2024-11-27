package services

import (
	"errors"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
)

type HealthProfileService interface {
	CreateHealthProfile(profile *dto.HealthProfileDto) error
}

type healthProfileService struct {
	healthRepo repositories.HealthProfileRepository
	authRepo   repositories.AuthRepository
}

func NewHealthProfileService(healthRepo repositories.HealthProfileRepository, authRepo repositories.AuthRepository) HealthProfileService {
	return &healthProfileService{
		healthRepo: healthRepo,
		authRepo:   authRepo,
	}
}

// CreateHealthProfile implements HealthProfileService.
func (h *healthProfileService) CreateHealthProfile(profile *dto.HealthProfileDto) error {
	if profile.Height <= 0 || profile.Weight <= 0 {
		return errors.New("height and weight must be greater than zero")
	}

	//find user by id
	_, err := h.authRepo.GetUserById(profile.UserID)
	if err != nil {
		return err
	}

	// create health profile data
	data := models.HealthProfile{
		UserID:          profile.UserID,
		Height:          profile.Height,
		Weight:          profile.Weight,
		IsDiabetic:      profile.IsDiabetic,
		SmokingHistory:  profile.SmokingHistory,
		HasHeartDisease: profile.HasHeartDisease,
		ActivityLevel:   profile.ActivityLevel,
	}
	data.BMI = profile.Weight / ((profile.Height / 100) * (profile.Height / 100)) // Calculate BMI

	// call repository to save the data
	if err := h.healthRepo.CreateHealthProfile(&data); err != nil {
		return err
	}

	// create diabetes details if user is diabetic
	if profile.IsDiabetic {
		diabetesData := models.DiabetesDetails{
			ProfileID:     data.ID,
			DiabeticType:  profile.DiabeticType,
			InsulinLevel:  profile.InsulinLevel,
			BloodPressure: profile.BloodPressure,
		}
		if err := h.healthRepo.CreateDiabetesDetails(&diabetesData); err != nil {
			return err
		}
	}

	return nil
}
