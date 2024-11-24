package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/errors"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/services"
)

// authHandler is a struct to handle auth request
type authHandler struct {
	authService services.AuthService
}

// NewAuthHandler is a constructor to create authHandler instance
func NewAuthHandler(authService services.AuthService) *authHandler {
	if authService == nil {
		panic("authService cannot be nil")
	}
	return &authHandler{
		authService: authService,
	}
}

// Register
func (h *authHandler) Register(c *gin.Context) {
	// get data from request
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// call service to register
	err := h.authService.Register(&req)
	if err != nil {
		if err.Error() == errors.ErrEmailAlreadyRegistered().Error() {
			errors.SendErrorResponse(c, http.StatusBadRequest, "failed to register", err.Error())
			return
		}
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to register", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "action success",
	})
}

// Login
func (h *authHandler) Login(c *gin.Context) {
	// get data from request
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// check if request body is empty
		if err.Error() == "EOF" {
			errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", "request body is empty")
			return
		}
		errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// get device info
	deviceInfo := dto.DeviceInfo{
		IP:        c.ClientIP(),
		Useragent: c.GetHeader("User-Agent"),
	}

	// call service to login
	res, err := h.authService.Login(&req, &deviceInfo)
	if err != nil {
		// check if user not verified
		if err.Error() == errors.ErrUserNotVerified().Error() {
			errors.SendErrorResponse(c, http.StatusBadRequest, "failed to login", err.Error())
			return
		}

		// check if email or password is invalid
		if err.Error() == errors.ErrInvalidEmailOrPassword().Error() {
			errors.SendErrorResponse(c, http.StatusBadRequest, "failed to login", err.Error())
			return
		}
		// check if failed to login
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to login", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "action success",
		"data":    res,
	})
}

// VerifyAccount
func (h *authHandler) VerifyAccount(c *gin.Context) {
	// get token and id from query
	token := c.Query("token")
	id := c.Param("id")

	// check if token or id is empty
	if token == "" || id == "" {
		errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", "token or id is empty")
		return
	}

	// call service to verify account
	err := h.authService.VerifyAccount(id, token)
	if err != nil {
		// check if token is invalid
		if err.Error() == errors.ErrInvalidToken().Error() {
			c.HTML(http.StatusOK, "error-verify-email.tmpl", nil)
			return
		}
		// check if failed to verify email
		c.HTML(http.StatusOK, "error-verify-email.tmpl", nil)
		return
	}
	// give success response
	c.HTML(http.StatusOK, "success-verify-email.tmpl", nil)
}

// ForgotPassword
func (h *authHandler) ForgotPassword(c *gin.Context) {
	// get data from request
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// call service to forgot password
	data, err := h.authService.ForgotPassword(&req)
	if err != nil {
		// check if email not found
		if err.Error() == errors.ErrEmailNotFound().Error() {
			errors.SendErrorResponse(c, http.StatusBadRequest, "failed to forgot password", err.Error())
			return
		}
		// check if failed to forgot password
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to forgot password", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "action success",
		"data": gin.H{
			"email":  data.Email,
			"expire": data.Expires_at.Format("2006-01-02 15:04:05"),
		},
	})

}

// Form Reset Password
func (h *authHandler) ShowResetPassword(c *gin.Context) {
	// get token from query
	token := c.Query("token")
	if token == "" {
		c.HTML(http.StatusOK, "error-reset-password.tmpl", gin.H{
			"message": "Link is invalid, please request a new link",
		})
		return
	}

	// check if token is valid
	_, err := h.authService.GetTokenByToken(token)
	if err != nil {
		c.HTML(http.StatusOK, "error-reset-password.tmpl", gin.H{
			"message": "link is invalid, please request a new link",
		})
		return
	}

	// give success response
	c.HTML(http.StatusOK, "reset-password.tmpl", gin.H{
		"token": token,
	})
}

// ResetPassword
func (h *authHandler) ResetPassword(c *gin.Context) {
	// get data from request
	token := c.PostForm("token")
	newPassword := c.PostForm("new-password")
	confirmPassword := c.PostForm("confirm-password")

	// check if token, new password, or confirm password is empty
	if token == "" || newPassword == "" || confirmPassword == "" {
		c.HTML(http.StatusOK, "error-reset-password.tmpl", gin.H{
			"message": "Invalid request data",
		})
	}
	// check if new password and confirm password is not same
	if newPassword != confirmPassword {
		c.HTML(http.StatusOK, "error-reset-password.tmpl", gin.H{
			"message": "Password and confirm password must be same",
		})
		return
	}

	// create request dto
	req := &dto.ResetPasswordRequest{
		Token:    token,
		Password: newPassword,
	}

	// call service to reset password
	err := h.authService.ResetPassword(req)
	if err != nil {
		// check if token is invalid
		if err.Error() == errors.ErrInvalidToken().Error() {
			c.HTML(http.StatusOK, "error-reset-password.tmpl", gin.H{
				"message": "Invalid token",
			})
			return
		}
		// check if failed to reset password
		c.HTML(http.StatusOK, "error-reset-password.tmpl", gin.H{
			"message": "Failed to reset password",
		})
		return
	}

	// give success response
	c.HTML(http.StatusOK, "success-reset-password.tmpl", nil)
}

// RefreshToken
func (h *authHandler) RefreshToken(c *gin.Context) {
	// get data from request
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// call service to refresh token
	res, err := h.authService.RefreshToken(&req)
	if err != nil {
		// check if token is invalid
		if err.Error() == errors.ErrInvalidToken().Error() {
			errors.SendErrorResponse(c, http.StatusBadRequest, "failed to refresh token", err.Error())
			return
		}
		// check if failed to refresh token
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to refresh token", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "action success",
		"data":    res,
	})
}

// Logout
func (h *authHandler) Logout(c *gin.Context) {
	// get token from header
	authHeader := c.GetHeader("Authorization")

	// check if authorization header is empty
	if authHeader == "" {
		errors.SendErrorResponse(c, http.StatusBadRequest, "failed to logout", "authorization header is required")
		return
	}

	// split authorization header
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		errors.SendErrorResponse(c, http.StatusBadRequest, "failed to logout", "authorization header is invalid")
		return

	}

	// create request dto
	req := dto.LogoutRequest{
		AccessToken: parts[1],
	}

	// call service to logout
	err := h.authService.Logout(&req)
	if err != nil {
		// check if token is invalid
		if err.Error() == errors.ErrInvalidToken().Error() {
			errors.SendErrorResponse(c, http.StatusBadRequest, "failed to logout", err.Error())
			return
		}
		// check if failed to logout
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to logout", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "action success",
	})
}

// ChangePassword
func (h *authHandler) ChangePassword(c *gin.Context) {
	// get data from request
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// check if request body is empty
		if err.Error() == "EOF" {
			errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", "request body is empty")
			return
		}
		// check if old password or new password is empty
		errors.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// get user id from context
	userID := c.GetString("userID")
	req.UserID = userID

	// call service to change password
	err := h.authService.ChangePassword(&req)
	if err != nil {
		// check if old password is invalid
		if err.Error() == errors.ErrInvalidPassword().Error() {
			errors.SendErrorResponse(c, http.StatusBadRequest, "failed to change password", err.Error())
			return
		}
		// check if failed to change password
		errors.SendErrorResponse(c, http.StatusInternalServerError, "failed to change password", err.Error())
		return
	}

	// give success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "action success",
	})
}
