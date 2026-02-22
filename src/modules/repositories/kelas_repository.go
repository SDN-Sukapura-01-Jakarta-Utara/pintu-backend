package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// GetKelasFilter represents filter parameters for GetAll
type GetKelasFilter struct {
	Name   string
	Status string
}

// GetKelasParams represents parameters for GetAll with filters
type GetKelasParams struct {
	Filter GetKelasFilter
	Limit  int
	Offset int
}

// KelasRepository handles data operations for Kelas
type KelasRepository interface {
	Create(data *models.Kelas) error
	GetByID(id uint) (*models.Kelas, error)
	GetAll(limit int, offset int) ([]models.Kelas, int64, error)
	GetAllWithFilter(params GetKelasParams) ([]models.Kelas, int64, error)
	GetByName(name string) (*models.Kelas, error)
	Update(data *models.Kelas) error
	Delete(id uint) error
}

type KelasRepositoryImpl struct {
	db *gorm.DB
}

// NewKelasRepository creates a new Kelas repository
func NewKelasRepository(db *gorm.DB) KelasRepository {
	return &KelasRepositoryImpl{db: db}
}

// Create creates a new Kelas record
func (r *KelasRepositoryImpl) Create(data *models.Kelas) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Kelas by ID
func (r *KelasRepositoryImpl) GetByID(id uint) (*models.Kelas, error) {
	var data models.Kelas
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Kelas records with pagination
func (r *KelasRepositoryImpl) GetAll(limit int, offset int) ([]models.Kelas, int64, error) {
	var data []models.Kelas
	var total int64

	// Get total count
	if err := r.db.Model(&models.Kelas{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByName retrieves Kelas by name
func (r *KelasRepositoryImpl) GetByName(name string) (*models.Kelas, error) {
	var data models.Kelas
	if err := r.db.Where("name = ?", name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates Kelas record
func (r *KelasRepositoryImpl) Update(data *models.Kelas) error {
	return r.db.Save(data).Error
}

// Delete deletes Kelas record by ID
func (r *KelasRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Kelas{}, id).Error
}

// GetAllWithFilter retrieves Kelas records with filters and pagination
func (r *KelasRepositoryImpl) GetAllWithFilter(params GetKelasParams) ([]models.Kelas, int64, error) {
	var data []models.Kelas
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+params.Filter.Name+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Model(&models.Kelas{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
