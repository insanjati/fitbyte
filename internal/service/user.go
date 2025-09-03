package service

import (
	"github.com/google/uuid"
	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) FindUserById(userId uuid.UUID) (*model.UserResponse, error) {
	users, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return users, nil
}
