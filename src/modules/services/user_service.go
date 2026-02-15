package services

import (
	"errors"

	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// UserService handles business logic for User
type UserService interface {
	Create(data *models.User) error
	GetByID(id uint) (*models.User, error)
	GetAll() ([]models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAllWithFilter(params repositories.GetUsersParams) ([]models.User, int64, error)
	Update(data *models.User) error
	Delete(id uint) error
}

type UserServiceImpl struct {
	repository     repositories.UserRepository
	roleRepository repositories.RoleRepository
}

// NewUserService creates a new User service
func NewUserService(repository repositories.UserRepository) UserService {
	return &UserServiceImpl{repository: repository}
}

// NewUserServiceWithRole creates a new User service with role repository
func NewUserServiceWithRole(repository repositories.UserRepository, roleRepo repositories.RoleRepository) UserService {
	return &UserServiceImpl{
		repository:     repository,
		roleRepository: roleRepo,
	}
}

// Create creates a new User
func (s *UserServiceImpl) Create(data *models.User) error {
	// Validate role exists if role_id is provided
	if data.RoleID != nil && *data.RoleID > 0 {
		if s.roleRepository != nil {
			role, err := s.roleRepository.GetByID(*data.RoleID)
			if err != nil || role == nil {
				return errors.New("role tidak ditemukan atau sudah dihapus")
			}
		}
	}

	return s.repository.Create(data)
}

// GetByID retrieves User by ID
func (s *UserServiceImpl) GetByID(id uint) (*models.User, error) {
	return s.repository.GetByID(id)
}

// GetAll retrieves all User
func (s *UserServiceImpl) GetAll() ([]models.User, error) {
	return s.repository.GetAll()
}

// GetByUsername retrieves user by username
func (s *UserServiceImpl) GetByUsername(username string) (*models.User, error) {
	return s.repository.GetByUsername(username)
}

// GetAllWithFilter retrieves users with filters and pagination
func (s *UserServiceImpl) GetAllWithFilter(params repositories.GetUsersParams) ([]models.User, int64, error) {
	return s.repository.GetAllWithFilter(params)
}

// Update updates User
func (s *UserServiceImpl) Update(data *models.User) error {
	// Validate role exists if role_id is provided
	if data.RoleID != nil && *data.RoleID > 0 {
		if s.roleRepository != nil {
			role, err := s.roleRepository.GetByID(*data.RoleID)
			if err != nil || role == nil {
				return errors.New("role tidak ditemukan atau sudah dihapus")
			}
		}
	}

	return s.repository.Update(data)
}

// Delete deletes User by ID
func (s *UserServiceImpl) Delete(id uint) error {
	// Validate user exists before delete
	user, err := s.repository.GetByID(id)
	if err != nil || user == nil {
		return errors.New("user tidak ditemukan atau sudah dihapus")
	}

	return s.repository.Delete(id)
}
