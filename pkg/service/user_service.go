package services

import (
	"atommuse/backend/authentication-service/pkg/model"
	"atommuse/backend/authentication-service/pkg/repository"
	"atommuse/backend/authentication-service/pkg/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	UpdateUserByID(userID string, updateUser *model.RequestUpdateUser) error
	ChangePassword(userID, oldPassword, newPassword string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) UpdateUserByID(userID string, updateUser *model.RequestUpdateUser) error {
	return s.userRepo.UpdateUserByID(userID, updateUser)
}

func (s *userService) ChangePassword(userID, oldPassword, newPassword string) error {
	// Retrieve the user by ID
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err // Return error if user retrieval fails
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := utils.VerifyPassword(user.Password, oldPassword); err != nil {
		return errors.New("invalid old password")
	}

	// Hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Update the user's password in the database
	err = s.userRepo.UpdateUserPasswordByID(userID, string(hashedNewPassword))
	if err != nil {
		return err
	}

	return nil
}
