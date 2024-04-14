package services

import (
	"atommuse/backend/authentication-service/pkg/model"
	"atommuse/backend/authentication-service/pkg/repository"
	"atommuse/backend/authentication-service/pkg/utils"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	UpdateUserByID(userID string, updateUser *model.RequestUpdateUser) (string, error)
	ChangePassword(userID, oldPassword, newPassword string) error
	GetUserByID(userID string) (*model.User, error)
	BanUser(ctx context.Context, userID string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) UpdateUserByID(userID string, updateUser *model.RequestUpdateUser) (string, error) {
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

func (s *userService) GetUserByID(userID string) (*model.User, error) {
	return s.userRepo.GetUserByID(userID)
}

func (s *userService) BanUser(ctx context.Context, userID string) error {
	return s.userRepo.BanUser(ctx, userID)
}
