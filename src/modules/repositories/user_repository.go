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
	AssignRoles(userID uint, roleIDs []uint) error
	RemoveRoles(userID uint) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

// GetUsersFilter represents filters for getting users
type GetUsersFilter struct {
	Nama     string
	Username string
	RoleIDs  []uint
	SystemID uint
	Status   string
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
	if err := r.db.Preload("Roles.System").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all User records
func (r *UserRepositoryImpl) GetAll() ([]models.User, error) {
	var data []models.User
	if err := r.db.Preload("Roles.System").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetByUsername retrieves user by username
func (r *UserRepositoryImpl) GetByUsername(username string) (*models.User, error) {
	var data models.User
	if err := r.db.Preload("Roles.System").Where("username = ?", username).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAllWithFilter retrieves users with filters and pagination
func (r *UserRepositoryImpl) GetAllWithFilter(params GetUsersParams) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Preload("Roles.System")

	// Apply filters
	if params.Filter.Nama != "" {
		query = query.Where("LOWER(users.nama) LIKE ?", "%"+strings.ToLower(params.Filter.Nama)+"%")
	}
	if params.Filter.Username != "" {
		query = query.Where("LOWER(users.username) LIKE ?", "%"+strings.ToLower(params.Filter.Username)+"%")
	}
	if len(params.Filter.RoleIDs) > 0 || params.Filter.SystemID > 0 {
		// Filter by roles and/or system
		query = query.Joins("INNER JOIN user_roles ON users.id = user_roles.user_id").
			Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
			Select("DISTINCT users.*")
		
		if len(params.Filter.RoleIDs) > 0 {
			query = query.Where("user_roles.role_id IN ?", params.Filter.RoleIDs)
		}
		if params.Filter.SystemID > 0 {
			query = query.Where("roles.system_id = ?", params.Filter.SystemID)
		}
	}
	if params.Filter.Status != "" {
		query = query.Where("users.status = ?", params.Filter.Status)
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
		"nama":          data.Nama,
		"username":      data.Username,
		"status":        data.Status,
		"updated_by_id": data.UpdatedByID,
		"updated_at":    data.UpdatedAt,
	})
	return result.Error
}

// Delete deletes User record by ID
func (r *UserRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// AssignRoles assigns multiple roles to a user
func (r *UserRepositoryImpl) AssignRoles(userID uint, roleIDs []uint) error {
	// Clear existing roles first
	if err := r.db.Table("user_roles").Where("user_id = ?", userID).Delete(nil).Error; err != nil {
		return err
	}

	// Insert new roles
	for _, roleID := range roleIDs {
		if err := r.db.Table("user_roles").Create(map[string]interface{}{
			"user_id": userID,
			"role_id": roleID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

// RemoveRoles removes all roles from a user
func (r *UserRepositoryImpl) RemoveRoles(userID uint) error {
	return r.db.Table("user_roles").Where("user_id = ?", userID).Delete(nil).Error
}
