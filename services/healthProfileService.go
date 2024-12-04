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
	UpdateHealthProfile(profile *dto.HealthProfileDto) error
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
	if !healthProfile.IsDiabetic {
		riskAssessment, err := h.healthRepo.GetRiskAssessmentByUserID(user.ID)
		if err != nil {
			return nil, err
		}

		resp.DiabetesPrediction = &dto.DiabetesPrediction{
			RiskPercentage: riskAssessment.RiskScore,
			RiskLevel:      string(riskAssessment.RiskLevel),
			Note:           riskAssessment.Note,
		}
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

		// save risk assessment
		riskData := models.RiskAssessment{
			ProfileID: data.ID,
			RiskScore: Prediction.Percentage,
			Note:      Prediction.Note,
		}

		if Prediction.Percentage > 70 {
			riskData.RiskLevel = "High"
		} else if Prediction.Percentage > 50 {
			riskData.RiskLevel = "Medium"
		} else {
			riskData.RiskLevel = "Low"
		}

		if err := h.healthRepo.CreateRiskAssessment(&riskData); err != nil {
			return err
		}
	}

	return nil
}

// UpdateHealthProfile implements HealthProfileService.
func (h *healthProfileService) UpdateHealthProfile(req *dto.HealthProfileDto) error {
	// 1. Find user by ID
	userData, err := h.authRepo.GetUserById(req.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// 2. Get or create health profile
	healthProfile, err := h.healthRepo.GetHealthProfileByUserID(req.UserID)
	if err != nil {
		return fmt.Errorf("failed to get health profile: %w", err)
	}

	// 3. Update basic health profile data
	healthProfile.Height = req.Height
	healthProfile.Weight = req.Weight
	healthProfile.SmokingHistory = req.SmokingHistory
	healthProfile.HasHeartDisease = req.HasHeartDisease
	healthProfile.ActivityLevel = req.ActivityLevel
	healthProfile.BMI = req.Weight / ((req.Height / 100) * (req.Height / 100))

	// 4. Handle diabetes-related logic
	if req.IsDiabetic {
		// a. Update or create diabetes details
		err := h.updateOrCreateDiabetesDetails(healthProfile, req)
		if err != nil {
			return fmt.Errorf("failed to handle diabetes details: %w", err)
		}

		// b. Delete risk assessment
		err = h.deleteRiskAssessmentIfExists(req.UserID)
		if err != nil {
			return fmt.Errorf("failed to delete risk assessment: %w", err)
		}

		healthProfile.IsDiabetic = true
	} else {
		// a. Delete diabetes details if they exist
		err := h.deleteDiabetesDetailsIfExists(healthProfile.ID)
		if err != nil {
			return fmt.Errorf("failed to delete diabetes details: %w", err)
		}

		// b. Create or update risk assessment
		err = h.updateRiskAssessment(userData, healthProfile, req)
		if err != nil {
			return fmt.Errorf("failed to update risk assessment: %w", err)
		}

		healthProfile.IsDiabetic = false
	}

	// 5. Save updated health profile
	if err := h.healthRepo.UpdateHealthProfile(healthProfile); err != nil {
		return fmt.Errorf("failed to update health profile: %w", err)
	}

	return nil
}

// Helper function to update or create diabetes details
func (h *healthProfileService) updateOrCreateDiabetesDetails(profile *models.HealthProfile, req *dto.HealthProfileDto) error {
	if profile.IsDiabetic {
		// Update existing diabetes details
		diabetesData, err := h.healthRepo.GetDiabetesDetailsByProfileID(fmt.Sprintf("%d", profile.ID))
		if err != nil {
			return err
		}

		diabetesData.DiabeticType = req.DiabeticType
		diabetesData.InsulinLevel = req.InsulinLevel
		diabetesData.BloodPressure = req.BloodPressure

		return h.healthRepo.UpdateDiabetesDetails(diabetesData)
	}

	// Create new diabetes details
	diabetesData := models.DiabetesDetails{
		ProfileID:     profile.ID,
		DiabeticType:  req.DiabeticType,
		InsulinLevel:  req.InsulinLevel,
		BloodPressure: req.BloodPressure,
	}
	return h.healthRepo.CreateDiabetesDetails(&diabetesData)
}

// Helper function to delete diabetes details if they exist
func (h *healthProfileService) deleteDiabetesDetailsIfExists(profileID uint) error {
	diabetesData, err := h.healthRepo.GetDiabetesDetailsByProfileID(fmt.Sprintf("%d", profileID))
	if err == nil {
		return h.healthRepo.DeleteDiabetesDetails(diabetesData)
	}
	// Ignore error if diabetes details don't exist
	return nil
}

// Helper function to delete risk assessment if it exists
func (h *healthProfileService) deleteRiskAssessmentIfExists(userID string) error {
	risk, err := h.healthRepo.GetRiskAssessmentByUserID(userID)
	if err == nil {
		return h.healthRepo.DeleteAssessment(risk)
	}
	// Ignore error if risk assessment doesn't exist
	return nil
}

// Helper function to update risk assessment
func (h *healthProfileService) updateRiskAssessment(userData *models.User, healthProfile *models.HealthProfile, req *dto.HealthProfileDto) error {
	predictionData := dto.DiabetesPredictionRequest{
		SmokingHistory: string(req.SmokingHistory),
		BMI:            healthProfile.BMI,
		Age:            userData.Age,
		HeartDisease:   req.HasHeartDisease,
		Gender:         userData.Gender,
	}

	Prediction, err := h.recomendRepo.DiabetesPrediction(&predictionData)
	if err != nil {
		return err
	}

	// Get or create risk assessment
	risk, err := h.healthRepo.GetRiskAssessmentByUserID(req.UserID)
	if err != nil {
		// Assume error means no existing assessment, create a new one
		risk = &models.RiskAssessment{}
	}

	risk.ProfileID = healthProfile.ID
	risk.RiskScore = Prediction.Percentage
	risk.Note = Prediction.Note

	// Determine risk level
	switch {
	case Prediction.Percentage > 70:
		risk.RiskLevel = "High"
	case Prediction.Percentage > 50:
		risk.RiskLevel = "Medium"
	default:
		risk.RiskLevel = "Low"
	}

	// Save risk assessment
	return h.healthRepo.UpdateAssessment(risk)
}
