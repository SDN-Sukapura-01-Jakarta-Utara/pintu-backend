package services

import (
	"errors"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
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
	// Try to get user from users table first
	user, userErr := s.repository.GetByUsername(req.Username)
	
	// If user not found in users table, try kepegawaian table
	if userErr != nil {
		kepegawaian, kepErr := s.repository.GetKepegawaianByUsername(req.Username)
		if kepErr != nil {
			return nil, errors.New("username atau password salah")
		}
		
		// Verify password for kepegawaian
		if err := bcrypt.CompareHashAndPassword([]byte(kepegawaian.Password), []byte(req.Password)); err != nil {
			return nil, errors.New("username atau password salah")
		}

		// Check if kepegawaian is active
		if kepegawaian.Status != "active" {
			return nil, errors.New("user tidak aktif")
		}

		// Filter roles that have system_id = 1 (PINTU)
		var pintuRoles []models.Role
		for _, role := range kepegawaian.Roles {
			if role.SystemID != nil && *role.SystemID == 1 {
				pintuRoles = append(pintuRoles, role)
			}
		}

		// Check if kepegawaian has at least one PINTU role
		if len(pintuRoles) == 0 {
			return nil, errors.New("anda tidak memiliki akses ke sistem PINTU")
		}

		// Use first PINTU role for token
		roleID := pintuRoles[0].ID

		// Generate JWT token for kepegawaian
		token, err := utils.GenerateToken(kepegawaian.ID, kepegawaian.Username, kepegawaian.Nama, &roleID, kepegawaian.Status)
		if err != nil {
			return nil, errors.New("gagal membuat token")
		}

		// Map PINTU roles and extract permissions
		roles := make([]dtos.RoleResponse, len(pintuRoles))
		permissionMap := make(map[string]bool)
		
		for i, role := range pintuRoles {
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

			// Collect unique permissions from PINTU roles
			for _, permission := range role.Permissions {
				permissionMap[permission.Name] = true
			}
		}

		// Convert permission map to slice
		permissions := make([]string, 0, len(permissionMap))
		for permName := range permissionMap {
			permissions = append(permissions, permName)
		}

		response := &dtos.LoginResponse{
			Token:       token,
			ExpiresAt:   time.Now().Add(24 * time.Hour),
			Permissions: permissions,
			User: dtos.UserLoginResponse{
				ID:        kepegawaian.ID,
				Nama:      kepegawaian.Nama,
				Username:  kepegawaian.Username,
				Status:    kepegawaian.Status,
				Roles:     roles,
				CreatedAt: kepegawaian.CreatedAt,
			},
		}

		return response, nil
	}

	// User found in users table, proceed with normal user login
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("username atau password salah")
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("user tidak aktif")
	}

	// Filter roles that have system_id = 1 (PINTU)
	var pintuRoles []models.Role
	for _, role := range user.Roles {
		if role.SystemID != nil && *role.SystemID == 1 {
			pintuRoles = append(pintuRoles, role)
		}
	}

	// Check if user has at least one PINTU role
	if len(pintuRoles) == 0 {
		return nil, errors.New("anda tidak memiliki akses ke sistem PINTU")
	}

	// Use first PINTU role for token
	roleID := pintuRoles[0].ID

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Nama, &roleID, user.Status)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	// Map PINTU roles and extract permissions
	roles := make([]dtos.RoleResponse, len(pintuRoles))
	permissionMap := make(map[string]bool) // To track unique permissions
	
	for i, role := range pintuRoles {
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

		// Collect unique permissions from PINTU roles
		for _, permission := range role.Permissions {
			permissionMap[permission.Name] = true
		}
	}

	// Convert permission map to slice
	permissions := make([]string, 0, len(permissionMap))
	for permName := range permissionMap {
		permissions = append(permissions, permName)
	}

	response := &dtos.LoginResponse{
		Token:       token,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		Permissions: permissions,
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
