package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// ArticleRepository handles data operations for Article
type ArticleRepository interface {
	Create(data *models.Article) error
	GetByID(id uint) (*models.Article, error)
	GetAll(limit int, offset int) ([]models.Article, int64, error)
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
