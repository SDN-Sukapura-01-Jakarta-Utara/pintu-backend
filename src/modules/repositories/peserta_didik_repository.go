package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"

	"gorm.io/gorm"
)

// GetPesertaDidikFilter represents filter parameters for GetAllWithFilter
type GetPesertaDidikFilter struct {
	TahunPelajaranID *uint
	RombelID         *uint
	Nama             string
	NIS              string
	JenisKelamin     string
	NISN             string
	TempatLahir      string
	NIK              string
	Agama            string
	Status           string
}

// GetPesertaDidikParams represents parameters for GetAllWithFilter with filters
type GetPesertaDidikParams struct {
	Filter GetPesertaDidikFilter
	Limit  int
	Offset int
}

// PesertaDidikRepository handles data operations for PesertaDidik
type PesertaDidikRepository interface {
	Create(data *models.PesertaDidik) error
	GetByID(id uint) (*models.PesertaDidik, error)
	GetByIDWithDetails(id uint) (*models.PesertaDidik, error)
	GetByNIS(nis string) (*models.PesertaDidik, error)
	GetByNISAndTahunPelajaran(nis string, tahunPelajaranID *uint) (*models.PesertaDidik, error)
	GetByUsername(username string) (*models.PesertaDidik, error)
	GetByUsernameAndTahunPelajaran(username string, tahunPelajaranID *uint) (*models.PesertaDidik, error)
	GetAll(limit int, offset int) ([]models.PesertaDidik, int64, error)
	GetAllWithFilter(params GetPesertaDidikParams) ([]models.PesertaDidik, int64, error)
	Update(data *models.PesertaDidik) error
	Delete(id uint) error
	AssignRoles(pesertaDidikID uint, roleIDs []uint) error
	RemoveRoles(pesertaDidikID uint) error
	GetRombelByName(name string) (*models.Rombel, error)
	GetTahunPelajaranByName(name string) (*models.TahunPelajaran, error)
	GetAllRombels() ([]models.Rombel, error)
	GetAllTahunPelajaran() ([]models.TahunPelajaran, error)
	GetAllTahunPelajaranAll() ([]models.TahunPelajaran, error)
	GetTotalSiswaByActiveTahunPelajaran() (int64, error)
}

type PesertaDidikRepositoryImpl struct {
	db *gorm.DB
}

// NewPesertaDidikRepository creates a new PesertaDidik repository
func NewPesertaDidikRepository(db *gorm.DB) PesertaDidikRepository {
	return &PesertaDidikRepositoryImpl{db: db}
}

// Create creates a new PesertaDidik record
func (r *PesertaDidikRepositoryImpl) Create(data *models.PesertaDidik) error {
	return r.db.Create(data).Error
}

// GetByID retrieves PesertaDidik by ID
func (r *PesertaDidikRepositoryImpl) GetByID(id uint) (*models.PesertaDidik, error) {
	var data models.PesertaDidik
	if err := r.db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByIDWithDetails retrieves PesertaDidik by ID with roles, rombel, and tahun_pelajaran preloaded
func (r *PesertaDidikRepositoryImpl) GetByIDWithDetails(id uint) (*models.PesertaDidik, error) {
	var data models.PesertaDidik
	if err := r.db.Preload("Roles.System").Preload("Rombel.Kelas").Preload("TahunPelajaran").First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByNIS retrieves PesertaDidik by NIS
func (r *PesertaDidikRepositoryImpl) GetByNIS(nis string) (*models.PesertaDidik, error) {
	var data models.PesertaDidik
	if err := r.db.Where("nis = ?", nis).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByNISAndTahunPelajaran retrieves PesertaDidik by NIS and TahunPelajaranID
func (r *PesertaDidikRepositoryImpl) GetByNISAndTahunPelajaran(nis string, tahunPelajaranID *uint) (*models.PesertaDidik, error) {
	var data models.PesertaDidik
	query := r.db.Where("nis = ?", nis)
	
	if tahunPelajaranID != nil && *tahunPelajaranID != 0 {
		query = query.Where("tahun_pelajaran_id = ?", *tahunPelajaranID)
	}
	
	if err := query.First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByUsername retrieves PesertaDidik by Username
func (r *PesertaDidikRepositoryImpl) GetByUsername(username string) (*models.PesertaDidik, error) {
	var data models.PesertaDidik
	if err := r.db.Where("username = ?", username).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByUsernameAndTahunPelajaran retrieves PesertaDidik by Username and TahunPelajaranID
func (r *PesertaDidikRepositoryImpl) GetByUsernameAndTahunPelajaran(username string, tahunPelajaranID *uint) (*models.PesertaDidik, error) {
	var data models.PesertaDidik
	query := r.db.Where("username = ?", username)
	
	if tahunPelajaranID != nil && *tahunPelajaranID != 0 {
		query = query.Where("tahun_pelajaran_id = ?", *tahunPelajaranID)
	}
	
	if err := query.First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAll retrieves all PesertaDidik records with pagination and preload relations
func (r *PesertaDidikRepositoryImpl) GetAll(limit int, offset int) ([]models.PesertaDidik, int64, error) {
	var data []models.PesertaDidik
	var total int64

	// Get total count
	if err := r.db.Model(&models.PesertaDidik{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data with preloaded relations, ordered by rombel_id and nama
	if err := r.db.Preload("Roles.System").Preload("Rombel.Kelas").Preload("TahunPelajaran").
		Limit(limit).Offset(offset).Order("rombel_id ASC, nama ASC").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetAllWithFilter retrieves PesertaDidik records with filters and pagination
// Filters: tahun_pelajaran_id, rombel_id, nama, nis, jenis_kelamin, nisn, tempat_lahir, nik, agama, status
// Order by: rombel_id, nama
func (r *PesertaDidikRepositoryImpl) GetAllWithFilter(params GetPesertaDidikParams) ([]models.PesertaDidik, int64, error) {
	var data []models.PesertaDidik
	var total int64

	query := r.db

	// Apply filters
	if params.Filter.TahunPelajaranID != nil && *params.Filter.TahunPelajaranID != 0 {
		query = query.Where("tahun_pelajaran_id = ?", *params.Filter.TahunPelajaranID)
	}
	if params.Filter.RombelID != nil && *params.Filter.RombelID != 0 {
		query = query.Where("rombel_id = ?", *params.Filter.RombelID)
	}
	if params.Filter.Nama != "" {
		query = query.Where("LOWER(nama) LIKE ?", "%"+strings.ToLower(params.Filter.Nama)+"%")
	}
	if params.Filter.NIS != "" {
		query = query.Where("LOWER(nis) LIKE ?", "%"+strings.ToLower(params.Filter.NIS)+"%")
	}
	if params.Filter.JenisKelamin != "" {
		query = query.Where("jenis_kelamin = ?", params.Filter.JenisKelamin)
	}
	if params.Filter.NISN != "" {
		query = query.Where("LOWER(nisn) LIKE ?", "%"+strings.ToLower(params.Filter.NISN)+"%")
	}
	if params.Filter.TempatLahir != "" {
		query = query.Where("LOWER(tempat_lahir) LIKE ?", "%"+strings.ToLower(params.Filter.TempatLahir)+"%")
	}
	if params.Filter.NIK != "" {
		query = query.Where("LOWER(nik) LIKE ?", "%"+strings.ToLower(params.Filter.NIK)+"%")
	}
	if params.Filter.Agama != "" {
		query = query.Where("LOWER(agama) = ?", strings.ToLower(params.Filter.Agama))
	}
	if params.Filter.Status != "" {
		query = query.Where("status = ?", params.Filter.Status)
	}

	// Get total count
	if err := query.Model(&models.PesertaDidik{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data with preloaded relations, ordered by rombel_id ASC, nama ASC
	if err := query.Preload("Roles.System").Preload("Rombel.Kelas").Preload("TahunPelajaran").
		Order("rombel_id ASC, nama ASC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// Update updates PesertaDidik record
func (r *PesertaDidikRepositoryImpl) Update(data *models.PesertaDidik) error {
	return r.db.Save(data).Error
}

// Delete deletes PesertaDidik record by ID
func (r *PesertaDidikRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.PesertaDidik{}, id).Error
}

// AssignRoles assigns multiple roles to a peserta didik
func (r *PesertaDidikRepositoryImpl) AssignRoles(pesertaDidikID uint, roleIDs []uint) error {
	// Only process if there are roles to assign
	if len(roleIDs) == 0 {
		return nil
	}

	// Clear existing roles first
	if err := r.db.Table("peserta_didik_roles").Where("peserta_didik_id = ?", pesertaDidikID).Delete(nil).Error; err != nil {
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
		if err := r.db.Table("peserta_didik_roles").Create(map[string]interface{}{
			"peserta_didik_id": pesertaDidikID,
			"role_id":          roleID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

// RemoveRoles removes all roles from a peserta didik
func (r *PesertaDidikRepositoryImpl) RemoveRoles(pesertaDidikID uint) error {
	return r.db.Table("peserta_didik_roles").Where("peserta_didik_id = ?", pesertaDidikID).Delete(nil).Error
}

// GetRombelByName retrieves Rombel by name
func (r *PesertaDidikRepositoryImpl) GetRombelByName(name string) (*models.Rombel, error) {
	var data models.Rombel
	if err := r.db.Where("LOWER(name) = ?", strings.ToLower(name)).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetTahunPelajaranByName retrieves TahunPelajaran by tahun_pelajaran value
func (r *PesertaDidikRepositoryImpl) GetTahunPelajaranByName(name string) (*models.TahunPelajaran, error) {
	var data models.TahunPelajaran
	if err := r.db.Where("LOWER(tahun_pelajaran) = ?", strings.ToLower(name)).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetAllRombels retrieves all active Rombel records ordered by name
func (r *PesertaDidikRepositoryImpl) GetAllRombels() ([]models.Rombel, error) {
	var data []models.Rombel
	if err := r.db.Where("status = ?", "active").Order("name ASC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetAllTahunPelajaran retrieves all active TahunPelajaran records
func (r *PesertaDidikRepositoryImpl) GetAllTahunPelajaran() ([]models.TahunPelajaran, error) {
	var data []models.TahunPelajaran
	if err := r.db.Where("status = ?", "active").Order("tahun_pelajaran ASC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetAllTahunPelajaranAll retrieves all TahunPelajaran records regardless of status
func (r *PesertaDidikRepositoryImpl) GetAllTahunPelajaranAll() ([]models.TahunPelajaran, error) {
	var data []models.TahunPelajaran
	if err := r.db.Order("tahun_pelajaran ASC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetTotalSiswaByActiveTahunPelajaran retrieves total count of peserta didik with active tahun pelajaran and active status
func (r *PesertaDidikRepositoryImpl) GetTotalSiswaByActiveTahunPelajaran() (int64, error) {
	var total int64
	
	// Join with tahun_pelajaran table and count peserta_didik where:
	// - tahun_pelajaran.status = 'active'
	// - peserta_didik.status = 'active'
	err := r.db.Model(&models.PesertaDidik{}).
		Joins("JOIN tahun_pelajaran ON peserta_didik.tahun_pelajaran_id = tahun_pelajaran.id").
		Where("tahun_pelajaran.status = ?", "active").
		Where("peserta_didik.status = ?", "active").
		Count(&total).Error
	
	if err != nil {
		return 0, err
	}
	
	return total, nil
}
