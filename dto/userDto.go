package dto

// update user request
type UpdateUserRequest struct {
	Email       string `form:"email" json:"email"`
	Password    string `form:"password" json:"password"`
	Name        string `form:"name" json:"name"`
	DateOfBirth string `form:"date_of_birth" json:"date_of_birth"`
	Gender      string `form:"gender" json:"gender"`
}
