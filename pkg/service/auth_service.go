// services/auth_service.go
package services

import (
	"atommuse/backend/authentication-service/pkg/model"
	"atommuse/backend/authentication-service/pkg/repository"
	"atommuse/backend/authentication-service/pkg/utils"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
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
		return "", errors.New("Email or Password Mismatch")
	}

	if user == nil || !utils.ComparePassword(user.Password, request.Password) {
		return "", errors.New("Email or Password Mismatch")
	}

	userIDString := user.ID.Hex() // Convert ObjectID to string
	token, err := utils.GenerateToken(userIDString, user.UserName, user.Role, user.ProfileImage, user.FirstName, user.LastName)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Register(request model.RegisterRequest) error {
	// Check if the email already exists in the database
	existingUser, err := s.userRepo.GetUserByEmail(request.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return err // Return error if something goes wrong other than not finding the document
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

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
