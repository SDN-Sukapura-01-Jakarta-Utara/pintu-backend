package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"

	"gorm.io/gorm"
)

// UserRepository handles data operations for User
type UserRepository interface {
	Create(data *models.User) error
	GetByID(id uint) (*models.User, error)
	GetAll() ([]models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAllWithFilter(params GetUsersParams) ([]models.User, int64, error)
	Update(data *models.User) error
	Delete(id uint) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

// GetUsersFilter represents filters for getting users
type GetUsersFilter struct {
	Nama             string
	Username         string
	RoleID           uint
	Status           string
	AccessibleSystem string
}

// GetUsersParams represents parameters for getting users with filters
type GetUsersParams struct {
	Filter GetUsersFilter
	Limit  int
	Offset int
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

// GetAllWithFilter retrieves users with filters and pagination
func (r *UserRepositoryImpl) GetAllWithFilter(params GetUsersParams) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Preload("Role")

	// Apply filters
	if params.Filter.Nama != "" {
		query = query.Where("LOWER(nama) LIKE ?", "%"+strings.ToLower(params.Filter.Nama)+"%")
	}
	if params.Filter.Username != "" {
		query = query.Where("LOWER(username) LIKE ?", "%"+strings.ToLower(params.Filter.Username)+"%")
	}
	if params.Filter.RoleID > 0 {
		query = query.Where("role_id = ?", params.Filter.RoleID)
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}
	if params.Filter.AccessibleSystem != "" {
		query = query.Where("accessible_system::text LIKE ?", "%"+params.Filter.AccessibleSystem+"%")
	}

	// Count total
	if err := query.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch data with pagination
	if err := query.Limit(params.Limit).Offset(params.Offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update updates User record
func (r *UserRepositoryImpl) Update(data *models.User) error {
	// Use Update with map to explicitly set fields
	result := r.db.Model(&models.User{}).Where("id = ?", data.ID).Updates(map[string]interface{}{
		"nama":               data.Nama,
		"username":           data.Username,
		"role_id":            data.RoleID,
		"accessible_system":  data.AccessibleSystem,
		"status":             data.Status,
		"updated_by_id":      data.UpdatedByID,
		"updated_at":         data.UpdatedAt,
	})
	return result.Error
}

// Delete deletes User record by ID
func (r *UserRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
