package services

import (
	"math"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
)

type MiniGroceryService interface {
	GetMiniGroceryWithPagination(page, limit int) (*dto.MiniGroceryPaginationResponse, error)
}

type miniGroceryService struct {
	miniGroceryRepo repositories.MiniGroceryRepository
}

func NewMiniGroceryService(miniGroceryRepo repositories.MiniGroceryRepository) MiniGroceryService {
	return &miniGroceryService{
		miniGroceryRepo: miniGroceryRepo,
	}
}

// GetMiniCourseWithPagination
func (s *miniGroceryService) GetMiniGroceryWithPagination(page, limit int) (*dto.MiniGroceryPaginationResponse, error) {
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
	miniGrocery, totalItems, err := s.miniGroceryRepo.GetMiniGroceryWithPagination(page, limit)
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	// Create response
	response := &dto.MiniGroceryPaginationResponse{
		Data: miniGrocery,
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
