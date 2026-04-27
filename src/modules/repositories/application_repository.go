package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// GetApplicationFilter represents filter parameters for GetAll
type GetApplicationFilter struct {
	Nama             string
	Status           string
	ShowInJumbotron  *bool
}

// GetApplicationParams represents parameters for GetAll with filters
type GetApplicationParams struct {
	Filter GetApplicationFilter
	Limit  int
	Offset int
}

// ApplicationRepository handles data operations for Application
type ApplicationRepository interface {
	Create(data *models.Application) error
	GetByID(id uint) (*models.Application, error)
	GetAll(limit int, offset int) ([]models.Application, int64, error)
	GetAllWithFilter(params GetApplicationParams) ([]models.Application, int64, error)
	GetByNama(nama string) (*models.Application, error)
	Update(data *models.Application) error
	Delete(id uint) error
	UnsetAllShowInJumbotron() error
	GetPublicList(showInJumbotron *bool) ([]models.Application, int64, error)
}

type ApplicationRepositoryImpl struct {
	db *gorm.DB
}

// NewApplicationRepository creates a new Application repository
func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &ApplicationRepositoryImpl{db: db}
}

// Create creates a new Application record
func (r *ApplicationRepositoryImpl) Create(data *models.Application) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Application by ID
func (r *ApplicationRepositoryImpl) GetByID(id uint) (*models.Application, error) {
	var data models.Application
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Application records with pagination
func (r *ApplicationRepositoryImpl) GetAll(limit int, offset int) ([]models.Application, int64, error) {
	var data []models.Application
	var total int64

	// Get total count
	if err := r.db.Model(&models.Application{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by nama ASC
	if err := r.db.Order("nama").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByNama retrieves Application by nama
func (r *ApplicationRepositoryImpl) GetByNama(nama string) (*models.Application, error) {
	var data models.Application
	if err := r.db.Where("nama = ?", nama).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates Application record
func (r *ApplicationRepositoryImpl) Update(data *models.Application) error {
	return r.db.Save(data).Error
}

// Delete deletes Application record by ID
func (r *ApplicationRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Application{}, id).Error
}

// GetAllWithFilter retrieves Application records with filters and pagination
func (r *ApplicationRepositoryImpl) GetAllWithFilter(params GetApplicationParams) ([]models.Application, int64, error) {
	var data []models.Application
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.Nama != "" {
		query = query.Where("nama ILIKE ?", "%"+params.Filter.Nama+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}
	if params.Filter.ShowInJumbotron != nil {
		query = query.Where("show_in_jumbotron = ?", *params.Filter.ShowInJumbotron)
	}

	// Get total count
	if err := query.Model(&models.Application{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by nama ASC
	if err := query.Order("nama").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// UnsetAllShowInJumbotron sets all applications show_in_jumbotron to false
func (r *ApplicationRepositoryImpl) UnsetAllShowInJumbotron() error {
	return r.db.Model(&models.Application{}).Where("show_in_jumbotron = ?", true).Update("show_in_jumbotron", false).Error
}

// GetPublicList retrieves all active applications for public display (sorted from oldest to newest)
func (r *ApplicationRepositoryImpl) GetPublicList(showInJumbotron *bool) ([]models.Application, int64, error) {
	var data []models.Application
	var total int64

	query := r.db.Where("status = ?", "active")

	// Apply filter if provided
	if showInJumbotron != nil {
		query = query.Where("show_in_jumbotron = ?", *showInJumbotron)
	}

	// Get total count
	if err := query.Model(&models.Application{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get all data ordered by created_at ASC (oldest to newest)
	if err := query.Order("created_at ASC").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
