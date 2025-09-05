package service

import (
	"encoding/json"
	"errors"
	"time"

	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/repository"
	"github.com/insanjati/fitbyte/internal/utils"
)

type UserService struct {
	userRepo   *repository.UserRepository
	cacheRepo  repository.CacheRepository
	userUtils  utils.PasswordHasher
	jwtService JwtService
}

func NewUserService(userRepo *repository.UserRepository, cache repository.CacheRepository, jwt JwtService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		cacheRepo:  cache,
		userUtils:  utils.NewPasswordHasher(),
		jwtService: jwt,
	}
}

func (s *UserService) FindUserById(userId uuid.UUID) (*model.UserResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:id:%s", userId.String())

	if cached, err := s.cacheRepo.Get(ctx, cacheKey); err == nil && cached != "" {
		var user model.UserResponse
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	_ = s.cacheRepo.Set(ctx, cacheKey, user, 10*time.Minute)

	return user, nil
}

func (s *UserService) UpdateUser(userId uuid.UUID, user *model.UpdateUserRequest) (*model.UserResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:id:%s", userId.String())

	prevUser, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Name == nil || *user.Name == "" {
		user.Name = prevUser.Name
	}
	if user.ImageUri == nil || *user.ImageUri == "" {
		user.ImageUri = prevUser.ImageUri
	}

	updated, err := s.userRepo.UpdateUser(userId, user)
	if err != nil {
		return nil, err
	}

	if err := s.cacheRepo.Delete(ctx, cacheKey); err != nil {
		fmt.Println("WARN: gagal hapus cache key", cacheKey, err)
	}

	if err := s.cacheRepo.Set(ctx, cacheKey, updated, 10*time.Minute); err != nil {
		fmt.Println("WARN: gagal set cache key", cacheKey, err)
	}

	return updated, nil
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
		if ctx.Err() != nil {
			return model.AuthResponse{}, fmt.Errorf("Error" + ctx.Err().Error())
		}
		return model.AuthResponse{}, fmt.Errorf("failed to create user: %v", err)
	}

	token, err := s.jwtService.GenerateToken(&createdUser)
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("failed to create user: %v", err)
	}

	return model.AuthResponse{Email: createdUser.Email, Token: token}, nil
}

func (s *UserService) Login(ctx context.Context, payload model.User) (model.AuthResponse, error) {

	if payload.Email == "" {
		return model.AuthResponse{}, fmt.Errorf("Email is Required")
	}
	if payload.Password == "" {
		return model.AuthResponse{}, fmt.Errorf("Password is Required")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf(err.Error())
	}

	if err := s.userUtils.ComparePasswordHash(user.Password, payload.Password); err != nil {
		return model.AuthResponse{}, fmt.Errorf("invalid credentials - password")
	}

	token, err := s.jwtService.GenerateToken(&user)
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("failed to ")
	}
	return model.AuthResponse{Email: payload.Email, Token: token}, nil
}
