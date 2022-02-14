package api

import (
	"errors"
	"strings"
)

// UserService contains the methods of the user service (methods are on below: functions wiht prefix (u *userService) ....)
type UserService interface {
	New(user NewUserRequest) error
}

// UserRepository is what lets our service do db operations without knowing anything about the implementation
// 規範repository有什麼功能(dependencies的地方)
type UserRepository interface {
	CreateUser(NewUserRequest) error
}

// 規範userService內含參數（同時參數內含的功能也適用上方規範內，意即這包內有的功能都適用）
type userService struct {
	storage UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		storage: userRepo,
	}
}

// below are methods
func (u *userService) New(user NewUserRequest) error {
	// do some basic validations
	if user.Email == "" {
		return errors.New("user service - email required")
	}

	if user.Name == "" {
		return errors.New("user service - name required")
	}

	if user.WeightGoal == "" {
		return errors.New("user service - weight goal required")
	}

	// do some basic normalisation
	user.Name = strings.ToLower(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	err := u.storage.CreateUser(user)

	if err != nil {
		return err
	}

	return nil
}
