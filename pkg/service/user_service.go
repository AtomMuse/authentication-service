package services

import (
	"atommuse/backend/authentication-service/pkg/model"
	"atommuse/backend/authentication-service/pkg/repository"
)

type UserService interface {
	UpdateUserByID(userID string, updateUser *model.RequestUpdateUser) error
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
