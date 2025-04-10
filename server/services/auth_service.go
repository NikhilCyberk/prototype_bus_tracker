// services/auth_service.go
package services

import (
	"bustracking/models"
	"bustracking/utils"
	"errors"

	"gorm.io/gorm"
)

// AuthService handles authentication operations
type AuthService struct {
	DB *gorm.DB
}

// NewAuthService creates a new auth service
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

// RegisterUser registers a new user
func (s *AuthService) RegisterUser(username, password, role string) (*models.User, error) {
	// Check if username already exists
	var count int64
	s.DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return nil, errors.New("username already exists")
	}

	// Validate role
	if role != "driver" && role != "passenger" && role != "admin" {
		return nil, errors.New("invalid role")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	result := s.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// AuthenticateUser verifies credentials and returns user if valid
func (s *AuthService) AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	result := s.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, result.Error
	}

	// Verify password
	if !utils.VerifyPassword(user.Password, password) {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

// GetUserByID retrieves user by ID
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := s.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetAssignedBus gets the bus assigned to a driver
func (s *AuthService) GetAssignedBus(driverID uint) (*models.Bus, error) {
	var bus models.Bus
	result := s.DB.Where("driver_id = ?", driverID).First(&bus)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bus, nil
}
