package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/email"
	errors "github.com/rizkirmdhnnn/sweetlife-backend-go/errors"
	helper "github.com/rizkirmdhnnn/sweetlife-backend-go/helpers"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	VerifyAccount(id, token string) error
	Login(req *dto.LoginRequest, info *dto.DeviceInfo) (*dto.LoginResponse, error)
	ForgotPassword(req *dto.ForgotPasswordRequest) (*models.Password_reset_tokens, error)
	ResetPassword(req *dto.ResetPasswordRequest) error
	GetTokenByToken(token string) (*models.Password_reset_tokens, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.TokenResponse, error)
	Logout(req *dto.LogoutRequest) error
	ChangePassword(req *dto.ChangePasswordRequest) error
}

type authService struct {
	repository    repositories.AuthRepository
	email         email.EmailClient
	healthProfile repositories.HealthProfileRepository
}

func NewAuthService(r repositories.AuthRepository, healthProfile repositories.HealthProfileRepository, emailClient email.EmailClient) AuthService {
	if r == nil {
		panic("repository cannot be nil")
	}
	return &authService{
		repository:    r,
		healthProfile: healthProfile,
		email:         emailClient,
	}
}

// Login implements AuthService.
func (s *authService) Login(req *dto.LoginRequest, info *dto.DeviceInfo) (*dto.LoginResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest()
	}

	// Cari user berdasarkan email dan verifikasi password
	user, err := s.repository.FindUserByEmail(req.Email)
	if err != nil || !helper.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.ErrInvalidEmailOrPassword()
	}

	// Cek apakah user sudah terverifikasi
	if user.Verified_at == nil {
		return nil, errors.ErrUserNotVerified()
	}

	// Ambil semua refresh token berdasarkan user ID
	existingTokens, err := s.repository.GetListRefreshToken(user.ID)
	if err != nil {
		existingTokens = nil // Tidak ada token ditemukan
	}

	// check health profile
	healthProfile, _ := s.healthProfile.CheckHealthProfileExist(user.ID)

	// Cari token yang sesuai
	var tokenToUpdate *models.RefreshToken
	for _, token := range existingTokens {
		var deviceInfo dto.DeviceInfo
		if err := json.Unmarshal([]byte(token.DeviceInfo), &deviceInfo); err == nil && deviceInfo.Useragent == info.Useragent {
			// Jika ada token dengan device yang sama
			if !token.Revoked {
				// Jika token tidak revoked, gunakan untuk update
				tokenToUpdate = &token
				break
			} else {
				// Jika token revoked, abaikan dan lanjutkan
				tokenToUpdate = nil
			}
		}
	}

	if tokenToUpdate != nil {
		// Update token yang ditemukan
		tokenToUpdate.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
		tokenToUpdate.UpdatedAt = time.Now()
		if err := s.repository.UpdateRefreshToken(tokenToUpdate); err != nil {
			return nil, errors.ErrInternalServer()
		}

		// Generate access token
		accessToken, err := helper.GenerateTokenAccessToken(user)
		if err != nil {
			return nil, errors.ErrInternalServer()
		}

		// Berikan respons
		return createLoginResponse(user, accessToken, tokenToUpdate.RefreshToken, healthProfile), nil
	}

	// Jika tidak ada token valid atau revoked, buat token baru
	return s.createNewTokens(user, info, healthProfile)
}

// createNewTokens membuat token baru dan menyimpannya ke database
func (s *authService) createNewTokens(user *models.User, info *dto.DeviceInfo, healthProfile bool) (*dto.LoginResponse, error) {
	// Generate refresh token
	refreshToken, err := helper.GenerateTokenRefresh(user)
	if err != nil {
		return nil, errors.ErrInternalServer()
	}

	// Generate access token
	accessToken, err := helper.GenerateTokenAccessToken(user)
	if err != nil {
		return nil, errors.ErrInternalServer()
	}

	// Serialize device info
	deviceInfoJSON, err := json.Marshal(info)
	if err != nil {
		return nil, errors.ErrInternalServer()
	}

	// Simpan refresh token baru ke database
	newRefreshToken := &models.RefreshToken{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		DeviceInfo:   string(deviceInfoJSON),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.repository.CreateRefreshToken(newRefreshToken); err != nil {
		return nil, errors.ErrInternalServer()
	}

	// Berikan respons
	return createLoginResponse(user, accessToken, refreshToken, healthProfile), nil
}

// createLoginResponse membuat respons login
func createLoginResponse(user *models.User, accessToken, refreshToken string, healthProfile bool) *dto.LoginResponse {
	return &dto.LoginResponse{
		Token: dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Type:         "Bearer",
		},

		User: dto.UserResponse{
			ID:               user.ID,
			Name:             user.Name,
			Email:            user.Email,
			Gender:           user.Gender,
			HasHealthProfile: &healthProfile,
			DateOfBirth:      user.DateOfBirth.Format("2006-01-02"),
		},
	}
}

func (s *authService) VerifyAccount(id, token string) error {
	user, err := s.repository.GetUserById(id)
	if err != nil {
		return errors.ErrInvalidRequest()
	}

	if user.Verified_at != nil {
		return errors.ErrInvalidRequest()
	}

	hash := helper.GenerateHash(user.Email, config.ENV.APP_KEY)

	if hash != token {
		return errors.ErrInvalidRequest()
	}

	// update data
	err = s.repository.VerifyUser(user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) Register(req *dto.RegisterRequest) error {
	if req == nil {
		return errors.ErrInvalidRequest()
	}

	existingUser, err := s.repository.FindUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return errors.ErrEmailAlreadyRegistered()
	}

	passwordHash, err := helper.HashPassword(req.Password)
	if err != nil {
		return err
	}
	hash := helper.GenerateHash(req.Email, config.ENV.APP_KEY)

	user := &models.User{
		Email:      req.Email,
		Password:   passwordHash,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	err = s.repository.CreateUser(user)
	if err != nil {
		return err
	}

	// find user by email
	user, err = s.repository.FindUserByEmail(req.Email)
	if err != nil {
		return err
	}

	// struct for email template
	type EmailData struct {
		VerificationLink string
	}

	// load html template
	tmpl, err := template.ParseFiles("templates/email/email-verification.tmpl")
	if err != nil {
		return err
	}

	// create email data
	emailData := &EmailData{
		VerificationLink: fmt.Sprintf("%s/api/v1/auth/verify/%s?token=%s", config.ENV.APP_HOST, user.ID, hash),
	}

	// create email body
	var emailBody strings.Builder
	err = tmpl.Execute(&emailBody, emailData)
	if err != nil {
		return err
	}

	go func() {
		maxRetries := 3
		var err error
		for i := 0; i < maxRetries; i++ {
			err = s.email.SendEmail(user.Email, "SweetLife - Email Verification", emailBody.String())
			if err == nil {
				fmt.Println("success send email verification to: ", user.Email)
				break
			}
			fmt.Println("Failed to send email to: ", user.Email, " Retrying...")
			time.Sleep(2 * time.Second) // Tunggu sebelum mencoba ulang
		}
		if err != nil {
			fmt.Println("Failed to send email to: ", user.Email)
			fmt.Println("Error: ", err)
			fmt.Println("Deleting user: ", user.ID)
			err := s.repository.DeleteUserById(user.ID)
			if err != nil {
				fmt.Println("Failed to delete user: ", user.ID)
			}
		}

	}()
	return nil
}

func (s *authService) ForgotPassword(req *dto.ForgotPasswordRequest) (*models.Password_reset_tokens, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest()
	}

	// find user by email
	user, err := s.repository.FindUserByEmail(req.Email)
	if err != nil {
		return nil, errors.ErrEmailNotFound()
	}

	// generate token
	token := helper.GenerateHash(user.Email, config.ENV.APP_KEY)
	tokenScema := &models.Password_reset_tokens{
		Email:      user.Email,
		Token:      token,
		Expires_at: time.Now().Add(time.Hour * 1),
	}

	//find token by email
	_, err = s.repository.GetResetTokenByEmail(user.Email)
	if err == nil {
		// delete token
		err = s.repository.DeleteResetTokenByEmail(user.Email)
		if err != nil {
			return nil, err
		}
	}

	// save token to database
	err = s.repository.CreateToken(tokenScema)
	if err != nil {
		return nil, err
	}

	// struct for email template
	type EmailData struct {
		ResetPasswordLink string
		Name              string
	}

	// load html template
	tmpl, err := template.ParseFiles("templates/email/reset-password.tmpl")
	if err != nil {
		return nil, err
	}

	// create email data
	emailData := &EmailData{
		ResetPasswordLink: fmt.Sprintf("%s/api/v1/auth/reset-password?token=%s", config.ENV.APP_HOST, token),
		Name:              user.Name,
	}

	// create email body
	var emailBody strings.Builder
	err = tmpl.Execute(&emailBody, emailData)
	if err != nil {
		return nil, err
	}

	// send email
	go func() {
		maxRetries := 3
		var err error
		for i := 0; i < maxRetries; i++ {
			s.email.SendEmail(user.Email, "SweetLife - Reset Password", emailBody.String())
			if err == nil {
				fmt.Println("success send email reset password to: ", user.Email)
				break
			}
			fmt.Println("Failed to send email to: ", user.Email, " Retrying...")
			time.Sleep(2 * time.Second)
		}
		if err != nil {
			fmt.Println("Failed to send email to: ", user.Email)
			fmt.Println("Error: ", err)
		}

	}()
	return tokenScema, nil
}

// ResetPassword implements AuthService.
func (s *authService) ResetPassword(req *dto.ResetPasswordRequest) error {
	if req == nil {
		return errors.ErrInvalidRequest()
	}

	// find token by token
	token, err := s.repository.GetResetTokenByToken(req.Token)
	if err != nil {
		return errors.ErrInvalidRequest()
	}

	// check token expired
	if token.Expires_at.Before(time.Now()) {
		return errors.ErrInvalidRequest()
	}

	// find user by email
	user, err := s.repository.FindUserByEmail(token.Email)
	if err != nil {
		return errors.ErrInvalidRequest()
	}

	// hash password
	passwordHash, err := helper.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// update user password
	user.Password = passwordHash
	err = s.repository.UpdateUser(user)
	if err != nil {
		return err
	}

	// delete token
	err = s.repository.DeleteResetTokenByToken(token.Token)
	if err != nil {
		return err
	}

	return nil
}

// GetTokenByToken implements AuthService.
func (s *authService) GetTokenByToken(token string) (*models.Password_reset_tokens, error) {
	return s.repository.GetResetTokenByToken(token)
}

// RefreshToken implements AuthService.
func (s *authService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest()
	}

	// find token by token
	token, err := s.repository.GetRefreshTokenByToken(req.RefreshToken)
	if err != nil {
		return nil, errors.ErrInvalidToken()
	}

	// check revoked
	if token.Revoked {
		return nil, errors.ErrInvalidToken()
	}

	// check token expired
	if token.ExpiresAt.Before(time.Now()) {
		// revoke token
		token.Revoked = true
		err = s.repository.UpdateRefreshToken(token)
		if err != nil {
			return nil, errors.ErrInternalServer()
		}

		return nil, errors.ErrInvalidToken()
	}

	// find user by email
	userid := string(token.UserID)
	user, err := s.repository.GetUserById(userid)
	if err != nil {
		return nil, errors.ErrInvalidToken()
	}

	// generate token
	newToken, err := helper.GenerateTokenAccessToken(user)
	if err != nil {
		return nil, errors.ErrInternalServer()
	}

	// generate refresh token
	refreshToken, err := helper.GenerateTokenRefresh(user)
	if err != nil {
		return nil, errors.ErrInternalServer()
	}

	// update data
	token.RefreshToken = refreshToken
	token.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
	token.UpdatedAt = time.Now()

	// save refresh token
	err = s.repository.UpdateRefreshToken(token)
	if err != nil {
		return nil, errors.ErrInternalServer()
	}

	// return response
	return &dto.TokenResponse{
		AccessToken:  newToken,
		RefreshToken: refreshToken,
		Type:         "Bearer",
	}, nil
}

// Logout implements AuthService.
func (s *authService) Logout(req *dto.LogoutRequest) error {
	if req == nil {
		return errors.ErrInvalidRequest()
	}

	// get id from token
	claims, err := helper.GetClaimsFromToken(req.AccessToken)
	if err != nil {
		fmt.Print("Error: ", err)
		return errors.ErrInvalidToken()
	}

	// find token by user id
	token, err := s.repository.GetRefreshTokenByUserId(claims.ID)
	if err != nil {
		return errors.ErrInvalidToken()
	}

	// delete token
	err = s.repository.DeleteRefreshTokenByToken(token.RefreshToken)
	if err != nil {
		return errors.ErrInternalServer()
	}

	return nil
}

// ChangePassword implements AuthService.
func (s *authService) ChangePassword(req *dto.ChangePasswordRequest) error {
	if req == nil {
		return errors.ErrInvalidRequest()
	}

	// get user by id
	user, err := s.repository.GetUserById(req.UserID)
	if err != nil {
		return errors.ErrInvalidRequest()
	}

	// check old password
	if !helper.CheckPasswordHash(req.OldPassword, user.Password) {
		return errors.ErrInvalidPassword()
	}

	// hash new password
	passwordHash, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// update user password
	user.Password = passwordHash
	err = s.repository.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
