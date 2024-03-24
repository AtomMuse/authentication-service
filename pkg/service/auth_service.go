// services/auth_service.go
package services

import (
	"atommuse/backend/authentication-service/pkg/model"
	"atommuse/backend/authentication-service/pkg/repository"
	"atommuse/backend/authentication-service/pkg/utils"
	"errors"
)

type AuthService interface {
	Login(request model.LoginRequest) (string, error)
	Register(request model.RegisterRequest) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Login(request model.LoginRequest) (string, error) {
	user, err := s.userRepo.GetUserByEmail(request.Email)
	if err != nil {
		return "", err
	}

	if user == nil || !utils.ComparePassword(user.Password, request.Password) {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.UserName, user.Role, user.ProfileImage, user.FirstName, user.LastName)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Register(request model.RegisterRequest) error {
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}

	user := model.User{
		Email:        request.Email,
		UserName:     request.UserName,
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Password:     hashedPassword,
		ProfileImage: request.ProfileImage,
	}

	err = s.userRepo.CreateUser(&user)
	if err != nil {
		return err
	}

	return nil
}
