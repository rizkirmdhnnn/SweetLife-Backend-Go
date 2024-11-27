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
