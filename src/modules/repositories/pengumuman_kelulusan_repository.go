package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// PengumumanKelulusanRepository handles data operations for PengumumanKelulusan
type PengumumanKelulusanRepository interface {
	Create(data *models.PengumumanKelulusan) error
	GetByID(id uint) (*models.PengumumanKelulusan, error)
	Update(data *models.PengumumanKelulusan) error
	GetFirst() (*models.PengumumanKelulusan, error)
}

type PengumumanKelulusanRepositoryImpl struct {
	db *gorm.DB
}

// NewPengumumanKelulusanRepository creates a new PengumumanKelulusan repository
func NewPengumumanKelulusanRepository(db *gorm.DB) PengumumanKelulusanRepository {
	return &PengumumanKelulusanRepositoryImpl{db: db}
}

// Create creates a new PengumumanKelulusan record
func (r *PengumumanKelulusanRepositoryImpl) Create(data *models.PengumumanKelulusan) error {
	return r.db.Create(data).Error
}

// GetByID retrieves PengumumanKelulusan by ID
func (r *PengumumanKelulusanRepositoryImpl) GetByID(id uint) (*models.PengumumanKelulusan, error) {
	var data models.PengumumanKelulusan
	if err := r.db.Preload("CreatedBy").Preload("UpdatedBy").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates a PengumumanKelulusan record
func (r *PengumumanKelulusanRepositoryImpl) Update(data *models.PengumumanKelulusan) error {
	return r.db.Save(data).Error
}

// GetFirst retrieves the PengumumanKelulusan record with ID 1
func (r *PengumumanKelulusanRepositoryImpl) GetFirst() (*models.PengumumanKelulusan, error) {
	var data models.PengumumanKelulusan
	// Always get record with ID = 1
	if err := r.db.Preload("CreatedBy").Preload("UpdatedBy").Where("id = ?", 1).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
