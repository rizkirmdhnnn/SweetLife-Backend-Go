package dto

import "github.com/rizkirmdhnnn/sweetlife-backend-go/models"

type MiniCoursePaginationResponse struct {
	Data       []models.MiniCourse `json:"data"`
	Pagination PaginationInfo      `json:"pagination"`
}

type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
} 