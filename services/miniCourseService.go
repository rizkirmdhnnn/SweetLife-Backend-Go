package services

import (
	"math"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
)

type MiniCourseService interface {
	GetMiniCourseWithPagination(page, limit int) (*dto.MiniCoursePaginationResponse, error)
}

type miniCourseService struct {
	miniCourseRepo repositories.MiniCourseRepository
}

func NewMiniCourseService(miniCourseRepo repositories.MiniCourseRepository) MiniCourseService {
	return &miniCourseService{
		miniCourseRepo: miniCourseRepo,
	}
}

// GetMiniCourseWithPagination
func (s *miniCourseService) GetMiniCourseWithPagination(page, limit int) (*dto.MiniCoursePaginationResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Get paginated data
	miniCourse, totalItems, err := s.miniCourseRepo.GetMiniCourseWithPagination(page, limit)
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	// Create response
	response := &dto.MiniCoursePaginationResponse{
		Data: miniCourse,
		Pagination: dto.PaginationInfo{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
	}

	return response, nil
}
