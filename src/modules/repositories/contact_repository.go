package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// ContactRepository handles data operations for Contact
type ContactRepository interface {
	Create(data *models.Contact) error
	GetByID(id uint) (*models.Contact, error)
	GetAll() ([]models.Contact, error)
	Update(data *models.Contact) error
	Delete(id uint) error
}

type ContactRepositoryImpl struct {
	db *gorm.DB
}

// NewContactRepository creates a new Contact repository
func NewContactRepository(db *gorm.DB) ContactRepository {
	return &ContactRepositoryImpl{db: db}
}

// Create creates a new Contact record
func (r *ContactRepositoryImpl) Create(data *models.Contact) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Contact by ID
func (r *ContactRepositoryImpl) GetByID(id uint) (*models.Contact, error) {
	var data models.Contact
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Contact records
func (r *ContactRepositoryImpl) GetAll() ([]models.Contact, error) {
	var data []models.Contact
	if err := r.db.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates Contact record
func (r *ContactRepositoryImpl) Update(data *models.Contact) error {
	return r.db.Save(data).Error
}

// Delete deletes Contact record by ID
func (r *ContactRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Contact{}, id).Error
}
