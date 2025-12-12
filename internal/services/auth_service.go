package services

import (
	"errors"
	"strings"

	"github.com/hoanghnt/TaskManagementAPI/internal/models"
	"github.com/hoanghnt/TaskManagementAPI/internal/repository"
	"github.com/hoanghnt/TaskManagementAPI/internal/utils"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtExpiry int
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, jwtExpiry int) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtExpiry: jwtExpiry,
	}
}

// Register registers a new user
func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Validate input
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.FullName = strings.TrimSpace(req.FullName)

	// Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Save user to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, s.jwtExpiry)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Login authenticates a user and returns JWT token
func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Validate input
	req.Username = strings.TrimSpace(req.Username)

	// Find user by username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, s.jwtExpiry)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// GetUserByID retrieves user by ID
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
