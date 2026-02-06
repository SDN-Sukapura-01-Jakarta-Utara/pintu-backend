package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// UserRepository handles data operations for User
type UserRepository interface {
	Create(data *models.User) error
	GetByID(id uint) (*models.User, error)
	GetAll() ([]models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(data *models.User) error
	Delete(id uint) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new User repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create creates a new User record
func (r *UserRepositoryImpl) Create(data *models.User) error {
	return r.db.Create(data).Error
}

// GetByID retrieves User by ID
func (r *UserRepositoryImpl) GetByID(id uint) (*models.User, error) {
	var data models.User
	if err := r.db.Preload("Role").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all User records
func (r *UserRepositoryImpl) GetAll() ([]models.User, error) {
	var data []models.User
	if err := r.db.Preload("Role").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetByUsername retrieves user by username
func (r *UserRepositoryImpl) GetByUsername(username string) (*models.User, error) {
	var data models.User
	if err := r.db.Preload("Role").Where("username = ?", username).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates User record
func (r *UserRepositoryImpl) Update(data *models.User) error {
	return r.db.Save(data).Error
}

// Delete deletes User record by ID
func (r *UserRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
