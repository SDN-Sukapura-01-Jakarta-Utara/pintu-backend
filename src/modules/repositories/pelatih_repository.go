package repositories

import (
	"errors"
	
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// GetPelatihFilter represents filters for GetAllWithFilter
type GetPelatihFilter struct {
	Nama              string
	EkstrakurikulerID *uint
	Status            string
}

// GetPelatihParams represents params for GetAllWithFilter
type GetPelatihParams struct {
	Filter GetPelatihFilter
	Limit  int
	Offset int
}

// PelatihRepository handles data operations for Pelatih
type PelatihRepository interface {
	Create(pelatih *models.Pelatih) error
	CreatePelatihEkstrakurikuler(mapping *models.PelatihEkstrakurikuler) error
	CreatePelatihRole(mapping *models.PelatihRole) error
	GetByID(id uint) (*models.Pelatih, error)
	GetByUsername(username string) (*models.Pelatih, error)
	GetAllWithFilter(params GetPelatihParams) ([]models.Pelatih, int64, error)
	GetRoleByID(roleID uint) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	GetActivePelatihEkstrakurikuler(pelatihID uint) ([]uint, error)
	RestorePelatihEkstrakurikuler(pelatihID uint, ekskulID uint) error
	SoftDeletePelatihEkstrakurikuler(pelatihID uint, ekskulID uint) error
	Update(pelatih *models.Pelatih) error
	DeleteAllPelatihEkstrakurikuler(pelatihID uint) error
	DeleteAllPelatihRoles(pelatihID uint) error
	Delete(id uint) error
}

type PelatihRepositoryImpl struct {
	db *gorm.DB
}

// NewPelatihRepository creates a new repository
func NewPelatihRepository(db *gorm.DB) PelatihRepository {
	return &PelatihRepositoryImpl{db: db}
}

func (r *PelatihRepositoryImpl) Create(pelatih *models.Pelatih) error {
	return r.db.Create(pelatih).Error
}

func (r *PelatihRepositoryImpl) CreatePelatihEkstrakurikuler(mapping *models.PelatihEkstrakurikuler) error {
	return r.db.Create(mapping).Error
}

func (r *PelatihRepositoryImpl) CreatePelatihRole(mapping *models.PelatihRole) error {
	return r.db.Create(mapping).Error
}

func (r *PelatihRepositoryImpl) GetByID(id uint) (*models.Pelatih, error) {
	var pelatih models.Pelatih
	if err := r.db.
		Preload("Ekstrakurikuler", func(db *gorm.DB) *gorm.DB {
			// Use Distinct to avoid duplicates from many-to-many join
			return db.Distinct().
				Joins("JOIN pelatih_ekstrakurikuler ON pelatih_ekstrakurikuler.ekstrakurikuler_id = ekstrakurikuler.id").
				Where("pelatih_ekstrakurikuler.pelatih_id = ? AND pelatih_ekstrakurikuler.deleted_at IS NULL", id)
		}).
		Preload("Roles.System").
		Preload("Roles.Permissions").
		First(&pelatih, id).Error; err != nil {
		return nil, err
	}
	return &pelatih, nil
}

func (r *PelatihRepositoryImpl) GetByUsername(username string) (*models.Pelatih, error) {
	var pelatih models.Pelatih
	if err := r.db.Preload("Roles.System").Preload("Roles.Permissions").Where("username = ?", username).First(&pelatih).Error; err != nil {
		return nil, err
	}
	return &pelatih, nil
}

func (r *PelatihRepositoryImpl) GetAllWithFilter(params GetPelatihParams) ([]models.Pelatih, int64, error) {
	var pelatih []models.Pelatih
	var total int64

	query := r.db.Model(&models.Pelatih{})

	// Apply filters
	if params.Filter.Nama != "" {
		query = query.Where("pelatih.nama ILIKE ?", "%"+params.Filter.Nama+"%")
	}
	if params.Filter.Status != "" {
		query = query.Where("pelatih.status = ?", params.Filter.Status)
	}
	if params.Filter.EkstrakurikulerID != nil {
		// Use EXISTS subquery to avoid duplicate rows
		query = query.Where("EXISTS (SELECT 1 FROM pelatih_ekstrakurikuler WHERE pelatih_ekstrakurikuler.pelatih_id = pelatih.id AND pelatih_ekstrakurikuler.ekstrakurikuler_id = ? AND pelatih_ekstrakurikuler.deleted_at IS NULL)", *params.Filter.EkstrakurikulerID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get pelatih with roles preload
	if err := query.
		Preload("Roles.System").
		Preload("Roles.Permissions").
		Limit(params.Limit).
		Offset(params.Offset).
		Order("pelatih.created_at DESC").
		Find(&pelatih).Error; err != nil {
		return nil, 0, err
	}

	// Manually load ekstrakurikuler for each pelatih to ensure correct filtering
	for i := range pelatih {
		var ekstrakurikuler []models.Ekstrakurikuler
		if err := r.db.
			Distinct().
			Joins("JOIN pelatih_ekstrakurikuler ON pelatih_ekstrakurikuler.ekstrakurikuler_id = ekstrakurikuler.id").
			Where("pelatih_ekstrakurikuler.pelatih_id = ? AND pelatih_ekstrakurikuler.deleted_at IS NULL", pelatih[i].ID).
			Find(&ekstrakurikuler).Error; err != nil {
			return nil, 0, err
		}
		pelatih[i].Ekstrakurikuler = ekstrakurikuler
	}

	return pelatih, total, nil
}

func (r *PelatihRepositoryImpl) GetRoleByID(roleID uint) (*models.Role, error) {
	var role models.Role
	if err := r.db.Preload("System").Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *PelatihRepositoryImpl) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.Preload("System").Preload("Permissions").Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *PelatihRepositoryImpl) Update(pelatih *models.Pelatih) error {
	return r.db.Save(pelatih).Error
}

func (r *PelatihRepositoryImpl) GetActivePelatihEkstrakurikuler(pelatihID uint) ([]uint, error) {
	var mappings []models.PelatihEkstrakurikuler
	if err := r.db.Where("pelatih_id = ?", pelatihID).Find(&mappings).Error; err != nil {
		return nil, err
	}
	
	ekskulIDs := make([]uint, len(mappings))
	for i, mapping := range mappings {
		ekskulIDs[i] = mapping.EkstrakurikulerID
	}
	return ekskulIDs, nil
}

func (r *PelatihRepositoryImpl) RestorePelatihEkstrakurikuler(pelatihID uint, ekskulID uint) error {
	// Restore (undelete) if exists with deleted_at not null
	result := r.db.Model(&models.PelatihEkstrakurikuler{}).
		Where("pelatih_id = ? AND ekstrakurikuler_id = ? AND deleted_at IS NOT NULL", pelatihID, ekskulID).
		Update("deleted_at", nil)
	
	// If no rows affected, return error so caller knows to create new
	if result.RowsAffected == 0 {
		return errors.New("no deleted record found to restore")
	}
	
	return result.Error
}

func (r *PelatihRepositoryImpl) SoftDeletePelatihEkstrakurikuler(pelatihID uint, ekskulID uint) error {
	// Soft delete specific mapping
	return r.db.Where("pelatih_id = ? AND ekstrakurikuler_id = ?", pelatihID, ekskulID).
		Delete(&models.PelatihEkstrakurikuler{}).Error
}

func (r *PelatihRepositoryImpl) DeleteAllPelatihEkstrakurikuler(pelatihID uint) error {
	// Soft delete all ekstrakurikuler mappings for this pelatih
	return r.db.Where("pelatih_id = ?", pelatihID).Delete(&models.PelatihEkstrakurikuler{}).Error
}

func (r *PelatihRepositoryImpl) DeleteAllPelatihRoles(pelatihID uint) error {
	// Hard delete all role mappings for this pelatih (junction table)
	return r.db.Unscoped().Where("pelatih_id = ?", pelatihID).Delete(&models.PelatihRole{}).Error
}

func (r *PelatihRepositoryImpl) Delete(id uint) error {
	// Soft delete pelatih
	return r.db.Delete(&models.Pelatih{}, id).Error
}
