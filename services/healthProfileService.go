package services

import (
	"errors"
	"fmt"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
)

type HealthProfileService interface {
	CreateHealthProfile(profile *dto.HealthProfileDto) error
	GetHealthProfile(userID string) (*dto.HealthProfileResponse, error)
}

type healthProfileService struct {
	healthRepo   repositories.HealthProfileRepository
	authRepo     repositories.AuthRepository
	recomendRepo repositories.RecomendationRepo
}

func NewHealthProfileService(healthRepo repositories.HealthProfileRepository, authRepo repositories.AuthRepository, recomend repositories.RecomendationRepo) HealthProfileService {
	return &healthProfileService{
		healthRepo:   healthRepo,
		authRepo:     authRepo,
		recomendRepo: recomend,
	}
}

// GetHealthProfile implements HealthProfileService.
func (h *healthProfileService) GetHealthProfile(userID string) (*dto.HealthProfileResponse, error) {
	// find user by id
	user, err := h.authRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	// get health profile
	healthProfile, err := h.healthRepo.GetHealthProfileByUserID(userID)
	if err != nil {
		return nil, err
	}

	resp := dto.HealthProfileResponse{
		UserID:          healthProfile.UserID,
		Height:          healthProfile.Height,
		Weight:          healthProfile.Weight,
		IsDiabetic:      healthProfile.IsDiabetic,
		SmokingHistory:  healthProfile.SmokingHistory,
		HasHeartDisease: healthProfile.HasHeartDisease,
		ActivityLevel:   healthProfile.ActivityLevel,
	}

	// if user is diabetic, get diabetes details
	if healthProfile.IsDiabetic {
		diabetesDetails, err := h.healthRepo.GetDiabetesDetailsByProfileID(fmt.Sprintf("%d", healthProfile.ID))
		if err != nil {
			return nil, err
		}
		resp.DiabetesDetails = &dto.DiabetesDetails{
			DiabeticType:  diabetesDetails.DiabeticType,
			InsulinLevel:  diabetesDetails.InsulinLevel,
			BloodPressure: diabetesDetails.BloodPressure,
		}
	}

	//get assessment
	riskAssessment, err := h.healthRepo.GetRiskAssessmentByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	resp.DiabetesPrediction = dto.DiabetesPrediction{
		RiskPercentage: riskAssessment.RiskScore,
		RiskLevel:      "aman|sedang|tinggi",
		Note:           riskAssessment.Note,
	}

	return &resp, nil
}

// CreateHealthProfile implements HealthProfileService.
func (h *healthProfileService) CreateHealthProfile(profile *dto.HealthProfileDto) error {
	if profile.Height <= 0 || profile.Weight <= 0 {
		return errors.New("height and weight must be greater than zero")
	}

	//find user by id
	userData, err := h.authRepo.GetUserById(profile.UserID)
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

	if !profile.IsDiabetic {
		// predict diabetes risk
		predictionData := dto.DiabetesPredictionRequest{
			SmokingHistory: string(profile.SmokingHistory),
			BMI:            data.BMI,
			Age:            userData.Age,
			HeartDisease:   profile.HasHeartDisease,
			Gender:         userData.Gender,
		}

		Prediction, err := h.recomendRepo.DiabetesPrediction(&predictionData)
		if err != nil {
			return err
		}

		// fmt.Print(Prediction)

		// save risk assessment
		riskData := models.RiskAssessment{
			ProfileID: data.ID,
			RiskScore: Prediction.Percentage,
			Note:      Prediction.Note,
		}

		if err := h.healthRepo.CreateRiskAssessment(&riskData); err != nil {
			return err
		}
	}

	return nil
}
