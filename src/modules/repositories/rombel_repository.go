package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// GetRombelFilter represents filter parameters for GetAll
type GetRombelFilter struct {
	Name    string
	Status  string
	KelasID uint
}

// GetRombelParams represents parameters for GetAll with filters
type GetRombelParams struct {
	Filter GetRombelFilter
	Limit  int
	Offset int
}

// RombelRepository handles data operations for Rombel
type RombelRepository interface {
	Create(data *models.Rombel) error
	GetByID(id uint) (*models.Rombel, error)
	GetAll(limit int, offset int) ([]models.Rombel, int64, error)
	GetAllWithFilter(params GetRombelParams) ([]models.Rombel, int64, error)
	GetByName(name string) (*models.Rombel, error)
	GetByKelasID(kelasID uint) ([]models.Rombel, error)
	Update(data *models.Rombel) error
	Delete(id uint) error
}

type RombelRepositoryImpl struct {
	db *gorm.DB
}

// NewRombelRepository creates a new Rombel repository
func NewRombelRepository(db *gorm.DB) RombelRepository {
	return &RombelRepositoryImpl{db: db}
}

// Create creates a new Rombel record
func (r *RombelRepositoryImpl) Create(data *models.Rombel) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Rombel by ID
func (r *RombelRepositoryImpl) GetByID(id uint) (*models.Rombel, error) {
	var data models.Rombel
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Rombel records with pagination
func (r *RombelRepositoryImpl) GetAll(limit int, offset int) ([]models.Rombel, int64, error) {
	var data []models.Rombel
	var total int64

	// Get total count
	if err := r.db.Model(&models.Rombel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC with Kelas relationship
	if err := r.db.Preload("Kelas").Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByName retrieves Rombel by name
func (r *RombelRepositoryImpl) GetByName(name string) (*models.Rombel, error) {
	var data models.Rombel
	if err := r.db.Where("name = ?", name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByKelasID retrieves all Rombel by Kelas ID
func (r *RombelRepositoryImpl) GetByKelasID(kelasID uint) ([]models.Rombel, error) {
	var data []models.Rombel
	if err := r.db.Where("kelas_id = ?", kelasID).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates Rombel record
func (r *RombelRepositoryImpl) Update(data *models.Rombel) error {
	return r.db.Save(data).Error
}

// Delete deletes Rombel record by ID
func (r *RombelRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Rombel{}, id).Error
}

// GetAllWithFilter retrieves Rombel records with filters and pagination
func (r *RombelRepositoryImpl) GetAllWithFilter(params GetRombelParams) ([]models.Rombel, int64, error) {
	var data []models.Rombel
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+params.Filter.Name+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}
	if params.Filter.KelasID != 0 {
		query = query.Where("kelas_id = ?", params.Filter.KelasID)
	}

	// Get total count
	if err := query.Model(&models.Rombel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC with Kelas relationship
	if err := query.Preload("Kelas").Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
