package repositories

import (
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"gorm.io/gorm"
)

type HealthProfileRepository interface {
	CreateHealthProfile(profile *models.HealthProfile) error
	CreateDiabetesDetails(details *models.DiabetesDetails) error
	CreateRiskAssessment(assessment *models.RiskAssessment) error
	GetRiskAssessmentByUserID(userID string) (*models.RiskAssessment, error)
	GetHealthProfileByUserID(userID string) (*models.HealthProfile, error)
	GetDiabetesDetailsByProfileID(profileID string) (*models.DiabetesDetails, error)

	UpdateHealthProfile(profile *models.HealthProfile) error
	UpdateAssessment(assessment *models.RiskAssessment) error
	UpdateDiabetesDetails(details *models.DiabetesDetails) error

	DeleteDiabetesDetails(details *models.DiabetesDetails) error
	DeleteAssessment(assessment *models.RiskAssessment) error
}
type healthProfileRepository struct {
	db *gorm.DB
}

func NewHealthProfileRepository(db *gorm.DB) HealthProfileRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &healthProfileRepository{
		db: db,
	}
}

// Create implements HealthProfileRepository.
func (h *healthProfileRepository) CreateHealthProfile(profile *models.HealthProfile) error {
	err := h.db.Create(&profile).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateDiabetesDetails implements HealthProfileRepository.
func (h *healthProfileRepository) CreateDiabetesDetails(details *models.DiabetesDetails) error {
	err := h.db.Create(&details).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateRiskAssessment implements HealthProfileRepository.
func (h *healthProfileRepository) CreateRiskAssessment(assessment *models.RiskAssessment) error {
	err := h.db.Create(&assessment).Error
	if err != nil {
		return err
	}
	return nil
}

// GetRiskAssessmentByUserID implements HealthProfileRepository.
func (h *healthProfileRepository) GetRiskAssessmentByUserID(userID string) (*models.RiskAssessment, error) {
	var assessment models.RiskAssessment
	var healthProfile models.HealthProfile

	err := h.db.Where("user_id = ?", userID).First(&healthProfile).Error
	if err != nil {
		return nil, err
	}

	err = h.db.Where("profile_id = ?", healthProfile.ID).First(&assessment).Error
	if err != nil {
		return nil, err
	}

	return &assessment, nil
}

// GetHealthProfileByUserID implements HealthProfileRepository.
func (h *healthProfileRepository) GetHealthProfileByUserID(userID string) (*models.HealthProfile, error) {
	var healthProfile models.HealthProfile
	err := h.db.Where("user_id = ?", userID).First(&healthProfile).Error
	if err != nil {
		return nil, err
	}
	return &healthProfile, nil
}

// GetDiabetesDetailsByProfileID implements HealthProfileRepository.
func (h *healthProfileRepository) GetDiabetesDetailsByProfileID(profileID string) (*models.DiabetesDetails, error) {
	var details models.DiabetesDetails
	err := h.db.Where("profile_id = ?", profileID).First(&details).Error
	if err != nil {
		return nil, err
	}
	return &details, nil
}

func (h *healthProfileRepository) UpdateHealthProfile(profile *models.HealthProfile) error {
	err := h.db.Save(&profile).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateAssessment implements HealthProfileRepository.
func (h *healthProfileRepository) UpdateAssessment(assessment *models.RiskAssessment) error {
	err := h.db.Save(&assessment).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateDiabetesDetails implements HealthProfileRepository.
func (h *healthProfileRepository) UpdateDiabetesDetails(details *models.DiabetesDetails) error {
	err := h.db.Save(&details).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteDIabetesDetails implements HealthProfileRepository.
func (h *healthProfileRepository) DeleteDiabetesDetails(details *models.DiabetesDetails) error {
	err := h.db.Delete(&details).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteAssessment implements HealthProfileRepository.
func (h *healthProfileRepository) DeleteAssessment(assessment *models.RiskAssessment) error {
	err := h.db.Delete(&assessment).Error
	if err != nil {
		return err
	}
	return nil
}
