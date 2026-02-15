package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// LoginRepository handles data operations for authentication
type LoginRepository interface {
	GetByUsername(username string) (*models.User, error)
}

type LoginRepositoryImpl struct {
	db *gorm.DB
}

// NewLoginRepository creates a new Login repository
func NewLoginRepository(db *gorm.DB) LoginRepository {
	return &LoginRepositoryImpl{db: db}
}

// GetByUsername retrieves user by username
func (r *LoginRepositoryImpl) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Roles").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
