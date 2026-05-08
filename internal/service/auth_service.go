package service

import (
	"errors"

	"backend/internal/config"
	"backend/internal/dto"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/utils"
	"github.com/google/uuid"
)

type AuthService interface {
	Signup(req dto.SignupRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
}

type authService struct {
	repo      repository.UserRepository
	orderRepo repository.OrderRepository
	cfg       *config.Config
}

func NewAuthService(repo repository.UserRepository, orderRepo repository.OrderRepository, cfg *config.Config) AuthService {
	return &authService{repo: repo, orderRepo: orderRepo, cfg: cfg}
}

func (s *authService) Signup(req dto.SignupRequest) (*dto.AuthResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "customer",
	}

	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("email might already exist")
	}

	token, err := utils.GenerateToken(user.ID, user.Role, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	// Link any previous guest orders
	_ = s.orderRepo.LinkGuestOrdersByEmail(user.Email, user.ID)

	return &dto.AuthResponse{Token: token, Message: "User created successfully"}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Role, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	// Link any previous guest orders (in case they ordered as guest before logging in)
	_ = s.orderRepo.LinkGuestOrdersByEmail(user.Email, user.ID)

	return &dto.AuthResponse{Token: token}, nil
}

func (s *authService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
