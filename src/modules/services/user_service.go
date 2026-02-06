package services

import (
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// UserService handles business logic for User
type UserService interface {
	Create(data *models.User) error
	GetByID(id uint) (*models.User, error)
	GetAll() ([]models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(data *models.User) error
	Delete(id uint) error
}

type UserServiceImpl struct {
	repository repositories.UserRepository
}

// NewUserService creates a new User service
func NewUserService(repository repositories.UserRepository) UserService {
	return &UserServiceImpl{repository: repository}
}

// Create creates a new User
func (s *UserServiceImpl) Create(data *models.User) error {
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

// Update updates User
func (s *UserServiceImpl) Update(data *models.User) error {
	return s.repository.Update(data)
}

// Delete deletes User by ID
func (s *UserServiceImpl) Delete(id uint) error {
	return s.repository.Delete(id)
}
