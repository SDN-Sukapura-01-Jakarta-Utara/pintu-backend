package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GetArticleFilter represents filter parameters for GetAllWithFilter
type GetArticleFilter struct {
	Judul            string
	StartDate        time.Time
	EndDate          time.Time
	Kategori         string
	Penulis          string
	StatusPublikasi  string
	Status           string
}

// GetArticleParams represents parameters for GetAllWithFilter with filters
type GetArticleParams struct {
	Filter GetArticleFilter
	Limit  int
	Offset int
}

// ArticleRepository handles data operations for Article
type ArticleRepository interface {
	Create(data *models.Article) error
	GetByID(id uint) (*models.Article, error)
	GetAll(limit int, offset int) ([]models.Article, int64, error)
	GetAllWithFilter(params GetArticleParams) ([]models.Article, int64, error)
	GetPublicLatest() ([]models.Article, error)
	Update(data *models.Article) error
	Delete(id uint) error
	DeleteByGambar(gambar string) error
}

type ArticleRepositoryImpl struct {
	db *gorm.DB
}

// NewArticleRepository creates a new Article repository
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &ArticleRepositoryImpl{db: db}
}

// Create creates a new Article record
func (r *ArticleRepositoryImpl) Create(data *models.Article) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Article by ID
func (r *ArticleRepositoryImpl) GetByID(id uint) (*models.Article, error) {
	var data models.Article
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Article records with pagination
func (r *ArticleRepositoryImpl) GetAll(limit int, offset int) ([]models.Article, int64, error) {
	var data []models.Article
	var total int64

	// Get total count
	if err := r.db.Model(&models.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data, ordered by created_at descending
	if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetAllWithFilter retrieves Article records with filters and pagination
func (r *ArticleRepositoryImpl) GetAllWithFilter(params GetArticleParams) ([]models.Article, int64, error) {
	var data []models.Article
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.Judul != "" {
		query = query.Where("LOWER(judul) LIKE ?", "%"+strings.ToLower(params.Filter.Judul)+"%")
	}
	if !params.Filter.StartDate.IsZero() && !params.Filter.EndDate.IsZero() {
		query = query.Where("tanggal >= ? AND tanggal <= ?", params.Filter.StartDate, params.Filter.EndDate)
	} else if !params.Filter.StartDate.IsZero() {
		query = query.Where("tanggal >= ?", params.Filter.StartDate)
	} else if !params.Filter.EndDate.IsZero() {
		query = query.Where("tanggal <= ?", params.Filter.EndDate)
	}
	if params.Filter.Kategori != "" {
		query = query.Where("LOWER(kategori) = ?", strings.ToLower(params.Filter.Kategori))
	}
	if params.Filter.Penulis != "" {
		query = query.Where("LOWER(penulis) LIKE ?", "%"+strings.ToLower(params.Filter.Penulis)+"%")
	}
	if params.Filter.StatusPublikasi != "" {
		query = query.Where("status_publikasi = ?", params.Filter.StatusPublikasi)
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Model(&models.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates Article record
func (r *ArticleRepositoryImpl) Update(data *models.Article) error {
	return r.db.Save(data).Error
}

// Delete deletes Article record by ID
func (r *ArticleRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Article{}, id).Error
}

// DeleteByGambar deletes Article record by gambar key
func (r *ArticleRepositoryImpl) DeleteByGambar(gambar string) error {
	return r.db.Where("gambar = ?", gambar).Delete(&models.Article{}).Error
}

// GetPublicLatest retrieves 10 latest published and active articles ordered by tanggal DESC
func (r *ArticleRepositoryImpl) GetPublicLatest() ([]models.Article, error) {
	var data []models.Article
	if err := r.db.Where("status = ? AND status_publikasi = ?", "active", "published").
		Order("tanggal DESC").
		Limit(10).
		Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
