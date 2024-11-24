package errors

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// struct for error response
type ErrorResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// function to create new error response
func NewErrorResponse(message string, errDetail string) ErrorResponse {
	return ErrorResponse{
		Status:  false,
		Message: message,
		Error:   errDetail,
	}
}

// function to send error response
func SendErrorResponse(c *gin.Context, status int, message, errDetail string) {
	c.JSON(status, NewErrorResponse(message, errDetail))
}

func ErrEmailAlreadyRegistered() error {
	return errors.New("email already registered")
}

func ErrInvalidRequest() error {
	return errors.New("invalid request data")
}

func ErrInvalidToken() error {
	return errors.New("invalid token")
}

func ErrTokenExpired() error {
	return errors.New("token has expired")
}

func ErrInvalidEmailOrPassword() error {
	return errors.New("invalid email or password")
}

func ErrInternalServer() error {
	return errors.New("internal server error")
}

func ErrEmailNotFound() error {
	return errors.New("email not found")
}

func ErrAuthHeaderNotFound() error {
	return errors.New("authorization header not found")
}

func ErrUserNotVerified() error {
	return errors.New("user not verified")
}

func ErrRefreshTokenRevoked() error {
	return errors.New("refresh token has been revoked")
}

func ErrRefreshTokenExpired() error {
	return errors.New("refresh token has expired")
}

func ErrInvalidPassword() error {
	return errors.New("invalid password")
}
