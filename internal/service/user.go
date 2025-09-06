package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/insanjati/fitbyte/internal/cache"
	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/repository"
	"github.com/insanjati/fitbyte/internal/utils"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo   *repository.UserRepository
	cache      *cache.Redis
	userUtils  utils.PasswordHasher
	jwtService JwtService
}

func NewUserService(userRepo *repository.UserRepository, cache *cache.Redis, jwt JwtService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		cache:      cache,
		userUtils:  utils.NewPasswordHasher(),
		jwtService: jwt,
	}
}

func (s *UserService) FindUserById(userId uuid.UUID) (*model.UserResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:id:%s", userId.String())

	// Try cache first using GetAs
	var cachedUser model.UserResponse
	if err := s.cache.GetAs(ctx, cacheKey, &cachedUser); err == nil {
		return &cachedUser, nil
	}

	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, user, 10*time.Minute)

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

	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		log.Printf("WARN: failed to invalidate cache for user %s: %v", userId, err)
	}

	if err := s.cache.SetExp(context.Background(), cacheKey, updated, 10*time.Minute); err != nil {
		log.Printf("WARN: failed to cache updated user %s: %v", userId, err)
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
			return model.AuthResponse{}, fmt.Errorf("context error: %v", ctx.Err())
		}
		return model.AuthResponse{}, fmt.Errorf("failed to create user: %v", err)
	}

	token, err := s.jwtService.GenerateToken(&createdUser)
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("failed to generate token: %v", err)
	}

	return model.AuthResponse{Email: createdUser.Email, Token: token}, nil
}

func (s *UserService) Login(ctx context.Context, payload model.User) (model.AuthResponse, error) {
	if payload.Email == "" {
		return model.AuthResponse{}, fmt.Errorf("email is required")
	}
	if payload.Password == "" {
		return model.AuthResponse{}, fmt.Errorf("password is required")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("user not found")
	}

	if err := s.userUtils.ComparePasswordHash(user.Password, payload.Password); err != nil {
		return model.AuthResponse{}, fmt.Errorf("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(&user)
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("failed to generate token: %v", err)
	}

	return model.AuthResponse{Email: payload.Email, Token: token}, nil
}
