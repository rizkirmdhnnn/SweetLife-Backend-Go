package repositories

import (
	"github.com/rizkirmdhnnn/sweetlife-backend-go/models"
	"gorm.io/gorm"
)

// UserRepository is a contract of user repository
type UserRepository interface {
	Update(user *models.User) error
}

// userRepository is a struct to store db connection
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository is a constructor to create user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &userRepository{
		db: db,
	}
}

// Implement Update method of userRepository
func (r *userRepository) Update(user *models.User) error {
	err := r.db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}
