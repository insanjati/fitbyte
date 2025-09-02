package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/repository"
	"github.com/insanjati/fitbyte/internal/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}


func (s *UserService) RegisterNewUser(c *gin.Context, payload model.User) (model.User, error) {
	// Business Logiv
	// Check if email exists Import get email function
	user, err := s.userRepo.GetUserByEmail(c, payload.Email)
	if user.ID != 0 {
		return model.User{}, fmt.Errorf("Email already Exists")
	} 
	if err != nil{
		return model.User{}, c.Err()
	}
	// Hash Password
	payload.Password, _ = utils.NewPasswordHasher().EncryptPassword(payload.Password)

	// Get Email only
	createdUser, err := s.userRepo.RegisterNewUser(c, payload)
	if err != nil {
		if c.Err() != nil{
			return model.User{}, c.Err()
		}
		return model.User{}, fmt.Errorf("failed to create user: %v", err)
	}

	return createdUser, nil
}

// Function is email valid()

