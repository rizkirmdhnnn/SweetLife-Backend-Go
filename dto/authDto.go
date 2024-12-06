package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type DeviceInfo struct {
	IP        string `json:"ip"`
	Useragent string `json:"useragent"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserResponse struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	Name             string `json:"name"`
	DateOfBirth      string `json:"date_of_birth"`
	Gender           string `json:"gender"`
	HasHealthProfile bool   `json:"has_health_profile"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"type"`
}

type LoginResponse struct {
	User  UserResponse  `json:"user"`
	Token TokenResponse `json:"token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LogoutRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type ChangePasswordRequest struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
