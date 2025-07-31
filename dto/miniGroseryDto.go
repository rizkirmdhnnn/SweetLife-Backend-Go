package dto

import "github.com/rizkirmdhnnn/sweetlife-backend-go/models"

type MiniGroceryPaginationResponse struct {
	Data       []models.MiniGrocery `json:"data"`
	Pagination PaginationInfo       `json:"pagination"`
}
