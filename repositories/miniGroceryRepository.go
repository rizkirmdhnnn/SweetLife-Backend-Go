package repositories

import (
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"gorm.io/gorm"
)

type MiniGroceryRepository interface {
	GetMiniGroceryWithPagination(page, limit int) ([]models.MiniGrocery, int, error)
}

type miniGroceryRepository struct {
	db *gorm.DB
}

func NewMiniGroceryRepository(db *gorm.DB) MiniGroceryRepository {
	return &miniGroceryRepository{
		db: db,
	}
}

// GetMiniCourseWithPagination
func (r *miniGroceryRepository) GetMiniGroceryWithPagination(page, limit int) ([]models.MiniGrocery, int, error) {
	var miniGrocery []models.MiniGrocery
	var total int64

	// Get total count
	err := r.db.Model(&models.MiniGrocery{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated data
	err = r.db.Offset(offset).Limit(limit).Find(&miniGrocery).Error
	if err != nil {
		return nil, 0, err
	}

	return miniGrocery, int(total), nil
}
