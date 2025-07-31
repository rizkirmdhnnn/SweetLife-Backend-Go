package repositories

import (
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"gorm.io/gorm"
)

type MiniCourseRepository interface {
	GetMiniCourseWithPagination(page, limit int) ([]models.MiniCourse, int, error)
}

type miniCourseRepository struct {
	db *gorm.DB
}

func NewMiniCourseRepository(db *gorm.DB) MiniCourseRepository {
	return &miniCourseRepository{
		db: db,
	}
}

// GetMiniCourseWithPagination
func (r *miniCourseRepository) GetMiniCourseWithPagination(page, limit int) ([]models.MiniCourse, int, error) {
	var miniCourse []models.MiniCourse
	var total int64

	// Get total count
	err := r.db.Model(&models.MiniCourse{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated data
	err = r.db.Offset(offset).Limit(limit).Find(&miniCourse).Error
	if err != nil {
		return nil, 0, err
	}

	return miniCourse, int(total), nil
}
