package repositories

import (
	"errors"
	"time"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"gorm.io/gorm"
)

// thanks copilot for generating documentation

// AuthRepository defines the interface for authentication-related operations.
// It includes methods for user management, password reset tokens, and refresh tokens.
type AuthRepository interface {
	// CreateUser creates a new user in the repository.
	// req: The user details to be created.
	// Returns an error if the operation fails.
	CreateUser(req *models.User) error

	// FindUserByEmail retrieves a user by their email address.
	// email: The email address of the user to be retrieved.
	// Returns the user details and an error if the operation fails.
	FindUserByEmail(email string) (*models.User, error)

	// GetUserById retrieves a user by their ID.
	// id: The ID of the user to be retrieved.
	// Returns the user details and an error if the operation fails.
	GetUserById(id string) (*models.User, error)

	// DeleteUserById deletes a user by their ID.
	// id: The ID of the user to be deleted.
	// Returns an error if the operation fails.
	DeleteUserById(id string) error

	// VerifyUser verifies a user's email address.
	// email: The email address of the user to be verified.
	// Returns an error if the operation fails.
	VerifyUser(email string) error

	// CreateToken creates a new password reset token.
	// req: The password reset token details to be created.
	// Returns an error if the operation fails.
	CreateToken(req *models.Password_reset_tokens) error

	// GetResetTokenByToken retrieves a password reset token by its token value.
	// token: The token value of the password reset token to be retrieved.
	// Returns the password reset token details and an error if the operation fails.
	GetResetTokenByToken(token string) (*models.Password_reset_tokens, error)

	// GetResetTokenByEmail retrieves a password reset token by the associated email address.
	// email: The email address associated with the password reset token to be retrieved.
	// Returns the password reset token details and an error if the operation fails.
	GetResetTokenByEmail(email string) (*models.Password_reset_tokens, error)

	// DeleteResetTokenByToken deletes a password reset token by its token value.
	// token: The token value of the password reset token to be deleted.
	// Returns an error if the operation fails.
	DeleteResetTokenByToken(token string) error

	// DeleteResetTokenByEmail deletes a password reset token by the associated email address.
	// email: The email address associated with the password reset token to be deleted.
	// Returns an error if the operation fails.
	DeleteResetTokenByEmail(email string) error

	// UpdateUser updates the details of an existing user.
	// user: The updated user details.
	// Returns an error if the operation fails.
	UpdateUser(user *models.User) error

	// CreateRefreshToken creates a new refresh token.
	// token: The refresh token details to be created.
	// Returns an error if the operation fails.
	CreateRefreshToken(token *models.RefreshToken) error

	// GetRefreshTokenByToken retrieves a refresh token by its token value.
	// token: The token value of the refresh token to be retrieved.
	// Returns the refresh token details and an error if the operation fails.
	GetRefreshTokenByToken(token string) (*models.RefreshToken, error)

	// GetRefreshTokenByUserId retrieves a refresh token by the associated user ID.
	// id: The ID of the user associated with the refresh token to be retrieved.
	// Returns the refresh token details and an error if the operation fails.
	GetRefreshTokenByUserId(id string) (*models.RefreshToken, error)

	// GetListRefreshToken retrieves a list of refresh tokens by the associated user ID.
	// id: The ID of the user associated with the refresh tokens to be retrieved.
	// Returns a list of refresh tokens and an error if the operation fails.
	GetListRefreshToken(id string) ([]models.RefreshToken, error)

	// UpdateRefreshToken updates the details of an existing refresh token.
	// token: The updated refresh token details.
	// Returns an error if the operation fails.
	UpdateRefreshToken(token *models.RefreshToken) error

	// DeleteRefreshTokenByToken deletes a refresh token by its token value.
	// token: The token value of the refresh token to be deleted.
	// Returns an error if the operation fails.
	DeleteRefreshTokenByToken(token string) error
}

type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new instance of AuthRepository.
func NewAuthRepository(db *gorm.DB) AuthRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &authRepository{
		db: db,
	}
}

// CreateUser implements AuthRepository.
func (r *authRepository) DeleteUserById(id string) error {
	err := r.db.Where("id = ?", id).Delete(&models.User{}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetUserById implements AuthRepository.
func (r *authRepository) GetUserById(id string) (*models.User, error) {
	var user models.User
	err := r.db.Raw("SELECT id, email, password, name, verified_at, gender FROM users WHERE id = ?", id).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail implements AuthRepository.
func (r *authRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Raw("SELECT id, email, password, name FROM users WHERE email = ?", email).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// VerifyUser implements AuthRepository.
func (r *authRepository) VerifyUser(email string) error {
	user, err := r.FindUserByEmail(email)
	if err != nil {
		return err
	}

	now := time.Now()
	user.Verified_at = &now
	err = r.db.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}

// CreateUser implements AuthRepository.
func (r *authRepository) CreateUser(user *models.User) error {
	err := r.db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *authRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// CreateToken implements AuthRepository.
func (r *authRepository) CreateToken(req *models.Password_reset_tokens) error {
	err := r.db.Create(&req).Error
	if err != nil {
		return err
	}
	return nil
}

// GetTokenByToken implements AuthRepository.
func (r *authRepository) GetResetTokenByToken(token string) (*models.Password_reset_tokens, error) {
	var tokenData models.Password_reset_tokens
	err := r.db.Where("token = ?", token).First(&tokenData).Error
	if err != nil {
		return nil, err
	}
	return &tokenData, nil
}

// GetTokenByEmail implements AuthRepository.
func (r *authRepository) GetResetTokenByEmail(email string) (*models.Password_reset_tokens, error) {
	var tokenData models.Password_reset_tokens
	err := r.db.Where("email = ?", email).First(&tokenData).Error
	if err != nil {
		return nil, err
	}
	return &tokenData, nil
}

// UpdateUser implements AuthRepository.
func (r *authRepository) UpdateUser(user *models.User) error {
	err := r.db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteTokenByToken implements AuthRepository.
func (r *authRepository) DeleteResetTokenByToken(token string) error {
	err := r.db.Where("token = ?", token).Delete(&models.Password_reset_tokens{}).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteTokenByEmail implements AuthRepository.
func (r *authRepository) DeleteResetTokenByEmail(email string) error {
	err := r.db.Where("email = ?", email).Delete(&models.Password_reset_tokens{}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetRefreshTokenByToken implements AuthRepository.
func (r *authRepository) GetRefreshTokenByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Raw("SELECT id, user_id, refresh_token, expires_at, created_at, revoked, device_info, updated_at FROM refresh_tokens WHERE refresh_token = ?", token).Scan(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	if refreshToken.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &refreshToken, nil
}

// SaveRefreshToken implements AuthRepository.
func (r *authRepository) CreateRefreshToken(token *models.RefreshToken) error {
	err := r.db.Create(&token).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateRefreshToken implements AuthRepository.
func (r *authRepository) UpdateRefreshToken(token *models.RefreshToken) error {
	err := r.db.Save(&token).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteRefreshTokenByToken implements AuthRepository.
func (r *authRepository) DeleteRefreshTokenByToken(token string) error {
	err := r.db.Where("refresh_token = ?", token).Delete(&models.RefreshToken{}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetRefreshTokenById implements AuthRepository.
func (r *authRepository) GetRefreshTokenByUserId(id string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Raw("SELECT id, user_id, refresh_token, expires_at, created_at, revoked, device_info, updated_at FROM refresh_tokens WHERE user_id = ?", id).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *authRepository) GetListRefreshToken(id string) ([]models.RefreshToken, error) {
	var refreshToken []models.RefreshToken
	err := r.db.Raw("SELECT id, user_id, refresh_token, expires_at, created_at, revoked, device_info, updated_at FROM refresh_tokens WHERE user_id = ?", id).Scan(&refreshToken).Error
	// err := r.db.Where("user_id = ?", id).Find(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}
