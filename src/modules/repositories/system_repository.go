package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"

	"gorm.io/gorm"
)

// SystemRepository handles data operations for System
type SystemRepository interface {
	Create(data *models.System) error
	GetByID(id uint) (*models.System, error)
	GetAll() ([]models.System, error)
	GetAllWithFilter(params GetSystemsParams) ([]models.System, int64, error)
	Update(data *models.System) error
	Delete(id uint) error
}

type SystemRepositoryImpl struct {
	db *gorm.DB
}

// GetSystemsFilter represents filters for getting systems
type GetSystemsFilter struct {
	Nama   string
	Status string
}

// GetSystemsParams represents parameters for getting systems with filters
type GetSystemsParams struct {
	Filter GetSystemsFilter
	Limit  int
	Offset int
}

// NewSystemRepository creates a new System repository
func NewSystemRepository(db *gorm.DB) SystemRepository {
	return &SystemRepositoryImpl{db: db}
}

// Create creates a new System record
func (r *SystemRepositoryImpl) Create(data *models.System) error {
	return r.db.Create(data).Error
}

// GetByID retrieves System by ID
func (r *SystemRepositoryImpl) GetByID(id uint) (*models.System, error) {
	var data models.System
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all System records
func (r *SystemRepositoryImpl) GetAll() ([]models.System, error) {
	var data []models.System
	if err := r.db.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetAllWithFilter retrieves systems with filters and pagination
func (r *SystemRepositoryImpl) GetAllWithFilter(params GetSystemsParams) ([]models.System, int64, error) {
	var systems []models.System
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.Nama != "" {
		query = query.Where("LOWER(nama) LIKE ?", "%"+strings.ToLower(params.Filter.Nama)+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Count total
	if err := query.Model(&models.System{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch data with pagination
	if err := query.Limit(params.Limit).Offset(params.Offset).Find(&systems).Error; err != nil {
		return nil, 0, err
	}

	return systems, total, nil
}

// Update updates System record
func (r *SystemRepositoryImpl) Update(data *models.System) error {
	result := r.db.Model(&models.System{}).Where("id = ?", data.ID).Updates(map[string]interface{}{
		"nama":           data.Nama,
		"description":    data.Description,
		"status":         data.Status,
		"updated_by_id":  data.UpdatedByID,
		"updated_at":     data.UpdatedAt,
	})
	return result.Error
}

// Delete deletes System record by ID
func (r *SystemRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.System{}, id).Error
}
