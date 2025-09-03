package service

import (
	"errors"

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

func (s *UserService) UpdateUser(userId uuid.UUID, user *model.UpdateUserRequest) (*model.UserResponse, error) {
	_, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// if user.Preference == nil || *user.Preference == "" {
	// 	return nil, errors.New("preference is required")
	// }
	// if *user.Preference != "CARDIO" && *user.Preference != "WEIGHT" {
	// 	return nil, errors.New("preference must be either 'CARDIO' or 'WEIGHT'")
	// }

	// if user.WeightUnit == nil || *user.WeightUnit == "" {
	// 	return nil, errors.New("weightUnit is required")
	// }
	// if *user.WeightUnit != "KG" && *user.WeightUnit != "LBS" {
	// 	return nil, errors.New("weightUnit must be either 'KG' or 'LBS'")
	// }

	// if user.HeightUnit == nil || *user.HeightUnit == "" {
	// 	return nil, errors.New("heightUnit is required")
	// }
	// if *user.HeightUnit != "CM" && *user.HeightUnit != "INCH" {
	// 	return nil, errors.New("heightUnit must be either 'CM' or 'INCH'")
	// }

	// if user.Weight == nil || *user.Weight < 10 || *user.Weight > 1000 {
	// 	return nil, errors.New("weight must be between 10 and 1000")
	// }

	// if user.Height == nil || *user.Height < 3 || *user.Height > 250 {
	// 	return nil, errors.New("height must be between 3 and 250")
	// }

	updated, err := s.userRepo.UpdateUser(userId, user)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
