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

	// Check if user has role with system_id = 1 (PINTU)
	var hasAccessToPINTU bool
	var roleID uint
	for _, role := range user.Roles {
		if role.SystemID != nil && *role.SystemID == 1 {
			hasAccessToPINTU = true
			if roleID == 0 {
				roleID = role.ID // Use first PINTU role for token
			}
		}
	}

	if !hasAccessToPINTU {
		return nil, errors.New("anda tidak memiliki akses ke sistem PINTU")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Nama, &roleID, user.Status)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	// Map roles to response
	roles := make([]dtos.RoleResponse, len(user.Roles))
	for i, role := range user.Roles {
		var system *dtos.SystemResponse
		if role.System != nil {
			system = &dtos.SystemResponse{
				ID:          role.System.ID,
				Nama:        role.System.Nama,
				Description: role.System.Description,
			}
		}

		roles[i] = dtos.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			SystemID:    role.SystemID,
			System:      system,
			Status:      role.Status,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
			CreatedByID: role.CreatedByID,
			UpdatedByID: role.UpdatedByID,
		}
	}

	response := &dtos.LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		User: dtos.UserLoginResponse{
			ID:        user.ID,
			Nama:      user.Nama,
			Username:  user.Username,
			Status:    user.Status,
			Roles:     roles,
			CreatedAt: user.CreatedAt,
		},
	}

	return response, nil
}
