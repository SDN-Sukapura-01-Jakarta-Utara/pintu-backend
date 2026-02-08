package services

import (
	"errors"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"

	"golang.org/x/crypto/bcrypt"
)

// LoginService handles business logic for authentication
type LoginService interface {
	Login(req *dtos.LoginRequest) (*dtos.LoginResponse, error)
}

type LoginServiceImpl struct {
	repository repositories.LoginRepository
}

// NewLoginService creates a new Login service
func NewLoginService(repository repositories.LoginRepository) LoginService {
	return &LoginServiceImpl{repository: repository}
}

// Login authenticates user and returns JWT token
func (s *LoginServiceImpl) Login(req *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	// Get user by username
	user, err := s.repository.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("username atau password salah")
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("user tidak aktif")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Nama, user.RoleID, user.Status)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	// Parse accessible systems
	var accessibleSystems []string
	if systems, err := user.AccessibleSystems(); err == nil {
		accessibleSystems = systems
	}

	// Prepare response
	roleID := uint(0)
	if user.RoleID != nil {
		roleID = *user.RoleID
	}

	response := &dtos.LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		User: dtos.UserLoginResponse{
			ID:               user.ID,
			Nama:             user.Nama,
			Username:         user.Username,
			Status:           user.Status,
			RoleID:           roleID,
			AccessibleSystem: accessibleSystems,
			CreatedAt:        user.CreatedAt,
		},
	}

	// Set role name if exists
	if user.Role != nil {
		response.User.RoleName = user.Role.Name
	}

	return response, nil
}
