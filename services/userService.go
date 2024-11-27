package services

import (
	"time"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/dto"
	helper "github.com/rizkirmdhnnn/sweetlife-backend-go/helpers"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
)

type UserService interface {
	// UpdateProfile
	UpdateProfile(id, urlPhotoProfile string, req *dto.UpdateUserRequest) error
	// Profile
	GetProfile(id string) (*dto.UserResponse, error)
}

type userService struct {
	userRepo repositories.UserRepository
	authRepo repositories.AuthRepository
}

func NewUserService(userRepo repositories.UserRepository, authRepo repositories.AuthRepository) UserService {
	if userRepo == nil {
		panic("userRepo cannot be nil")
	}

	if authRepo == nil {
		panic("authRepo cannot be nil")
	}

	return &userService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

// UpdateUser implements UserService.
func (u *userService) UpdateProfile(id, urlPhotoProfile string, req *dto.UpdateUserRequest) error {
	// get user by id
	user, err := u.authRepo.GetUserById(id)
	if err != nil {
		return err
	}

	// parse date
	date, _ := helper.ParsedDate(req.DateOfBirth)

	// update user
	user.Name = req.Name
	user.DateOfBirth = date
	user.Gender = req.Gender
	user.Updated_at = time.Now()
	user.ImageUrl = urlPhotoProfile
	user.Age, _ = helper.CalculateAge(date.Format("2006-01-02"))

	// save user
	err = u.userRepo.Update(user)
	if err != nil {
		return err
	}
	return nil
}

// GetProfile implements UserService.
func (u *userService) GetProfile(id string) (*dto.UserResponse, error) {
	// get user by id
	user, err := u.authRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	// create response
	res := dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		DateOfBirth: user.DateOfBirth.Format("2006-01-02"),
		Gender:      user.Gender,
	}

	return &res, nil
}
