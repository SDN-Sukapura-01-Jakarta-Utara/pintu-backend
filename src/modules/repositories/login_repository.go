package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// LoginRepository handles data operations for authentication
type LoginRepository interface {
	GetByUsername(username string) (*models.User, error)
	GetKepegawaianByUsername(username string) (*models.Kepegawaian, error)
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
	if err := r.db.Preload("Roles.System").Preload("Roles.Permissions").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetKepegawaianByUsername retrieves kepegawaian by username
func (r *LoginRepositoryImpl) GetKepegawaianByUsername(username string) (*models.Kepegawaian, error) {
	var kepegawaian models.Kepegawaian
	if err := r.db.Preload("Roles.System").Preload("Roles.Permissions").Where("username = ?", username).First(&kepegawaian).Error; err != nil {
		return nil, err
	}
	return &kepegawaian, nil
}
