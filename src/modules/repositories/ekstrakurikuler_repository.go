package repositories

import (
	"fmt"
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// GetEkstrakurikulerFilter represents filter parameters for GetAll
type GetEkstrakurikulerFilter struct {
	Name     string
	KelasID  uint
	Kategori string
	Status   string
}

// GetEkstrakurikulerParams represents parameters for GetAll with filters
type GetEkstrakurikulerParams struct {
	Filter GetEkstrakurikulerFilter
	Limit  int
	Offset int
}

// EkstrakurikulerRepository handles data operations for Ekstrakurikuler
type EkstrakurikulerRepository interface {
	Create(data *models.Ekstrakurikuler) error
	GetByID(id uint) (*models.Ekstrakurikuler, error)
	GetAll(limit int, offset int) ([]models.Ekstrakurikuler, int64, error)
	GetAllWithFilter(params GetEkstrakurikulerParams) ([]models.Ekstrakurikuler, int64, error)
	GetByName(name string) (*models.Ekstrakurikuler, error)
	GetByKategori(kategori string) ([]models.Ekstrakurikuler, error)
	Update(data *models.Ekstrakurikuler) error
	Delete(id uint) error
}

type EkstrakurikulerRepositoryImpl struct {
	db *gorm.DB
}

// NewEkstrakurikulerRepository creates a new Ekstrakurikuler repository
func NewEkstrakurikulerRepository(db *gorm.DB) EkstrakurikulerRepository {
	return &EkstrakurikulerRepositoryImpl{db: db}
}

// Create creates a new Ekstrakurikuler record
func (r *EkstrakurikulerRepositoryImpl) Create(data *models.Ekstrakurikuler) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Ekstrakurikuler by ID
func (r *EkstrakurikulerRepositoryImpl) GetByID(id uint) (*models.Ekstrakurikuler, error) {
	var data models.Ekstrakurikuler
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Ekstrakurikuler records with pagination
func (r *EkstrakurikulerRepositoryImpl) GetAll(limit int, offset int) ([]models.Ekstrakurikuler, int64, error) {
	var data []models.Ekstrakurikuler
	var total int64

	// Get total count
	if err := r.db.Model(&models.Ekstrakurikuler{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByName retrieves Ekstrakurikuler by name
func (r *EkstrakurikulerRepositoryImpl) GetByName(name string) (*models.Ekstrakurikuler, error) {
	var data models.Ekstrakurikuler
	if err := r.db.Where("name = ?", name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByKategori retrieves all Ekstrakurikuler by kategori
func (r *EkstrakurikulerRepositoryImpl) GetByKategori(kategori string) ([]models.Ekstrakurikuler, error) {
	var data []models.Ekstrakurikuler
	if err := r.db.Where("kategori = ?", kategori).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates Ekstrakurikuler record
func (r *EkstrakurikulerRepositoryImpl) Update(data *models.Ekstrakurikuler) error {
	return r.db.Save(data).Error
}

// Delete deletes Ekstrakurikuler record by ID
func (r *EkstrakurikulerRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Ekstrakurikuler{}, id).Error
}

// GetAllWithFilter retrieves Ekstrakurikuler records with filters and pagination
func (r *EkstrakurikulerRepositoryImpl) GetAllWithFilter(params GetEkstrakurikulerParams) ([]models.Ekstrakurikuler, int64, error) {
	var data []models.Ekstrakurikuler
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+params.Filter.Name+"%")
	}
	if params.Filter.Kategori != "" {
		query = query.Where("kategori = ?", params.Filter.Kategori)
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}
	// For kelas_id, check if it exists in the JSONB array
	if params.Filter.KelasID != 0 {
		query = query.Where("kelas_ids::jsonb @> ?::jsonb", fmt.Sprintf(`[%d]`, params.Filter.KelasID))
	}

	// Get total count
	if err := query.Model(&models.Ekstrakurikuler{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
