package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"

	"gorm.io/gorm"
)

// GetKepegawaianFilter represents filter parameters for GetAllWithFilter
type GetKepegawaianFilter struct {
	Nama     string
	Username string
	NIP      string
	NKKI     string
	Kategori string
	Jabatan  string
	RoleID   uint
	Status   string
}

// GetKepegawaianParams represents parameters for GetAllWithFilter with filters
type GetKepegawaianParams struct {
	Filter GetKepegawaianFilter
	Limit  int
	Offset int
}

// KepegawaianRepository handles data operations for Kepegawaian
type KepegawaianRepository interface {
	Create(data *models.Kepegawaian) error
	GetByID(id uint) (*models.Kepegawaian, error)
	GetByIDWithRoles(id uint) (*models.Kepegawaian, error)
	GetByNIP(nip string) (*models.Kepegawaian, error)
	GetByUsername(username string) (*models.Kepegawaian, error)
	GetAll(limit int, offset int) ([]models.Kepegawaian, int64, error)
	GetAllWithFilter(params GetKepegawaianParams) ([]models.Kepegawaian, int64, error)
	Update(data *models.Kepegawaian) error
	Delete(id uint) error
	AssignRoles(kepegawaianID uint, roleIDs []uint) error
	RemoveRoles(kepegawaianID uint) error
	GetTotalPendidik() (int64, error)
	GetTotalTendik() (int64, error)
}

type KepegawaianRepositoryImpl struct {
	db *gorm.DB
}

// NewKepegawaianRepository creates a new Kepegawaian repository
func NewKepegawaianRepository(db *gorm.DB) KepegawaianRepository {
	return &KepegawaianRepositoryImpl{db: db}
}

// Create creates a new Kepegawaian record
func (r *KepegawaianRepositoryImpl) Create(data *models.Kepegawaian) error {
	return r.db.Create(data).Error
}

// GetByID retrieves Kepegawaian by ID
func (r *KepegawaianRepositoryImpl) GetByID(id uint) (*models.Kepegawaian, error) {
	var data models.Kepegawaian
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByIDWithRoles retrieves Kepegawaian by ID with roles, bidang_studi, and rombel preloaded
func (r *KepegawaianRepositoryImpl) GetByIDWithRoles(id uint) (*models.Kepegawaian, error) {
	var data models.Kepegawaian
	if err := r.db.Preload("Roles.System").Preload("BidangStudi").Preload("RombelGuruKelas").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByNIP retrieves Kepegawaian by NIP
func (r *KepegawaianRepositoryImpl) GetByNIP(nip string) (*models.Kepegawaian, error) {
	var data models.Kepegawaian
	if err := r.db.Where("nip = ?", nip).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByUsername retrieves Kepegawaian by Username
func (r *KepegawaianRepositoryImpl) GetByUsername(username string) (*models.Kepegawaian, error) {
	var data models.Kepegawaian
	if err := r.db.Where("username = ?", username).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all Kepegawaian records with pagination
func (r *KepegawaianRepositoryImpl) GetAll(limit int, offset int) ([]models.Kepegawaian, int64, error) {
	var data []models.Kepegawaian
	var total int64

	// Get total count
	if err := r.db.Model(&models.Kepegawaian{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data, ordered by created_at descending
	if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetAllWithFilter retrieves Kepegawaian records with filters and pagination
func (r *KepegawaianRepositoryImpl) GetAllWithFilter(params GetKepegawaianParams) ([]models.Kepegawaian, int64, error) {
	var data []models.Kepegawaian
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.Nama != "" {
		query = query.Where("LOWER(nama) LIKE ?", "%"+strings.ToLower(params.Filter.Nama)+"%")
	}
	if params.Filter.Username != "" {
		query = query.Where("LOWER(username) LIKE ?", "%"+strings.ToLower(params.Filter.Username)+"%")
	}
	if params.Filter.NIP != "" {
		query = query.Where("LOWER(nip) LIKE ?", "%"+strings.ToLower(params.Filter.NIP)+"%")
	}
	if params.Filter.NKKI != "" {
		query = query.Where("LOWER(nkki) LIKE ?", "%"+strings.ToLower(params.Filter.NKKI)+"%")
	}
	if params.Filter.Kategori != "" {
		query = query.Where("LOWER(kategori) = ?", strings.ToLower(params.Filter.Kategori))
	}
	if params.Filter.Jabatan != "" {
		query = query.Where("LOWER(jabatan) LIKE ?", "%"+strings.ToLower(params.Filter.Jabatan)+"%")
	}
	if params.Filter.RoleID != 0 {
		query = query.Joins("INNER JOIN kepegawaian_roles ON kepegawaian.id = kepegawaian_roles.kepegawaian_id").
			Where("kepegawaian_roles.role_id = ?", params.Filter.RoleID)
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Model(&models.Kepegawaian{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data ordered by created_at DESC
	if err := query.Order("kepegawaian.created_at DESC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates Kepegawaian record
func (r *KepegawaianRepositoryImpl) Update(data *models.Kepegawaian) error {
	return r.db.Save(data).Error
}

// Delete deletes Kepegawaian record by ID
func (r *KepegawaianRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Kepegawaian{}, id).Error
}

// AssignRoles assigns multiple roles to a kepegawaian
func (r *KepegawaianRepositoryImpl) AssignRoles(kepegawaianID uint, roleIDs []uint) error {
	// Only process if there are roles to assign
	if len(roleIDs) == 0 {
		return nil
	}

	// Clear existing roles first
	if err := r.db.Table("kepegawaian_roles").Where("kepegawaian_id = ?", kepegawaianID).Delete(nil).Error; err != nil {
		return err
	}

	// Insert new roles (deduplicate first)
	roleMap := make(map[uint]bool)
	var uniqueRoleIDs []uint
	for _, roleID := range roleIDs {
		if !roleMap[roleID] {
			uniqueRoleIDs = append(uniqueRoleIDs, roleID)
			roleMap[roleID] = true
		}
	}

	for _, roleID := range uniqueRoleIDs {
		if err := r.db.Table("kepegawaian_roles").Create(map[string]interface{}{
			"kepegawaian_id": kepegawaianID,
			"role_id":        roleID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

// RemoveRoles removes all roles from a kepegawaian
func (r *KepegawaianRepositoryImpl) RemoveRoles(kepegawaianID uint) error {
	return r.db.Table("kepegawaian_roles").Where("kepegawaian_id = ?", kepegawaianID).Delete(nil).Error
}

// GetTotalPendidik retrieves total count of kepegawaian with kategori "Pendidik" and status "active"
func (r *KepegawaianRepositoryImpl) GetTotalPendidik() (int64, error) {
	var total int64
	
	err := r.db.Model(&models.Kepegawaian{}).
		Where("LOWER(kategori) = ?", "pendidik").
		Where("status = ?", "active").
		Count(&total).Error
	
	if err != nil {
		return 0, err
	}
	
	return total, nil
}

// GetTotalTendik retrieves total count of kepegawaian with kategori "Tenaga Kependidikan" and status "active"
func (r *KepegawaianRepositoryImpl) GetTotalTendik() (int64, error) {
	var total int64
	
	err := r.db.Model(&models.Kepegawaian{}).
		Where("LOWER(kategori) = ?", "tenaga kependidikan").
		Where("status = ?", "active").
		Count(&total).Error
	
	if err != nil {
		return 0, err
	}
	
	return total, nil
}
