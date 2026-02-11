package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// VisiMisiRepository handles data operations for VisiMisi
type VisiMisiRepository interface {
	Create(data *models.VisiMisi) error
	GetByID(id uint) (*models.VisiMisi, error)
	GetAll(limit int, offset int) ([]models.VisiMisi, int64, error)
	Update(data *models.VisiMisi) error
	Delete(id uint) error
}

type VisiMisiRepositoryImpl struct {
	db *gorm.DB
}

// NewVisiMisiRepository creates a new VisiMisi repository
func NewVisiMisiRepository(db *gorm.DB) VisiMisiRepository {
	return &VisiMisiRepositoryImpl{db: db}
}

// Create creates a new VisiMisi record
func (r *VisiMisiRepositoryImpl) Create(data *models.VisiMisi) error {
	return r.db.Create(data).Error
}

// GetByID retrieves VisiMisi by ID
func (r *VisiMisiRepositoryImpl) GetByID(id uint) (*models.VisiMisi, error) {
	var data models.VisiMisi
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all VisiMisi records with pagination
func (r *VisiMisiRepositoryImpl) GetAll(limit int, offset int) ([]models.VisiMisi, int64, error) {
	var data []models.VisiMisi
	var total int64

	// Get total count
	if err := r.db.Model(&models.VisiMisi{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	if err := r.db.Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates VisiMisi record
func (r *VisiMisiRepositoryImpl) Update(data *models.VisiMisi) error {
	return r.db.Save(data).Error
}

// Delete deletes VisiMisi record by ID
func (r *VisiMisiRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.VisiMisi{}, id).Error
}
