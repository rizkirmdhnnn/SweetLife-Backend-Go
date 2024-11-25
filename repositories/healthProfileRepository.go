package repositories

import (
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"gorm.io/gorm"
)

type HealthProfileRepository interface {
	CreateHealthProfile(profile *models.HealthProfile) error
	CreateDiabetesDetails(details *models.DiabetesDetails) error
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
