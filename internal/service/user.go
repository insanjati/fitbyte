package service

import (
	"context"
	"fmt"

	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/repository"
	"github.com/insanjati/fitbyte/internal/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
	userUtils utils.PasswordHasher
	jwtService JwtService
}

func NewUserService(userRepo *repository.UserRepository, jwt JwtService) *UserService {
	return &UserService{userRepo: userRepo,
						userUtils: utils.NewPasswordHasher(),
						jwtService: jwt}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}


func (s *UserService) RegisterNewUser(ctx context.Context, payload model.User) (model.AuthResponse, error) {
	// Business Logiv
	// Check if email exists Import get email function

	// user, err := s.userRepo.GetUserByEmail(ctx, payload.Email)
	// if err != nil{
	// 	return model.User{}, ctx.Err()
	// }
	// if user.ID == "" {
	// 	return model.User{}, fmt.Errorf("Email already Exists")
	// } 
	// Hash Password

	hashedPassword, err := s.userUtils.EncryptPassword(payload.Password)
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("failed to encrypt password: %v", err)
	}
	payload.Password = hashedPassword
	// Payload exists here
	createdUser, err := s.userRepo.RegisterNewUser(ctx, payload)
	if err != nil {
		if ctx.Err() != nil{
			return model.AuthResponse{}, fmt.Errorf("Error" + ctx.Err().Error())
		}
		return model.AuthResponse{}, fmt.Errorf("failed to create user: %v", err)
	}

	

	token, err := s.jwtService.GenerateToken(&createdUser)
	if err != nil{
		return model.AuthResponse{}, fmt.Errorf("failed to create user: %v", err)
	}


	return model.AuthResponse{Email: createdUser.Email, Token:  token}, nil
}

func (s *UserService) Login(ctx context.Context, payload model.User) (model.AuthResponse, error){

	if payload.Email == ""{
		return model.AuthResponse{}, fmt.Errorf("Email is Required")
	}
	if payload.Password == ""{
		return model.AuthResponse{}, fmt.Errorf("Password is Required")
	}

	fmt.Println(payload.Email)
	user, err := s.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil{
		return model.AuthResponse{}, fmt.Errorf(err.Error())
	}

	if err:= s.userUtils.ComparePasswordHash(user.Password, payload.Password); err!=nil{
		return model.AuthResponse{}, fmt.Errorf("invalid credentials - password")
	}

	token, err:= s.jwtService.GenerateToken(&user)
	if err!=nil{
		return model.AuthResponse{}, fmt.Errorf("failed to ")
	}
	return model.AuthResponse{Email: payload.Email, Token: token}, nil
}