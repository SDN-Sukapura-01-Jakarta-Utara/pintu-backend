package repositories

import (
	"pintu-backend/src/modules/models"
	"strings"

	"gorm.io/gorm"
)

// GetPesertaDidikFilter represents filter parameters for GetAllWithFilter
type GetPesertaDidikFilter struct {
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
	GetByIDs(ids []uint) ([]models.PesertaDidik, error)
	GetByIDWithDetails(id uint) (*models.PesertaDidik, error)
	GetByNIS(nis string) (*models.PesertaDidik, error)
	GetByNISN(nisn string) (*models.PesertaDidik, error)
	GetByNISAndTahunPelajaran(nis string, tahunPelajaranID *uint) (*models.PesertaDidik, error)
	GetByUsername(username string) (*models.PesertaDidik, error)
	GetByUsernameAndTahunPelajaran(username string, tahunPelajaranID *uint) (*models.PesertaDidik, error)
	GetAll(limit int, offset int) ([]models.PesertaDidik, int64, error)
	GetAllWithFilter(params GetPesertaDidikParams) ([]models.PesertaDidik, int64, error)
	GetAllActive() ([]models.PesertaDidik, error)
	Update(data *models.PesertaDidik) error
	UpdateInTransaction(tx interface{}, data *models.PesertaDidik) error
	Delete(id uint) error
	AssignRoles(pesertaDidikID uint, roleIDs []uint) error
	RemoveRoles(pesertaDidikID uint) error
	GetRombelByName(name string) (*models.Rombel, error)
	GetRombelByID(id uint) (*models.Rombel, error)
	GetTahunPelajaranByName(name string) (*models.TahunPelajaran, error)
	GetAllRombels() ([]models.Rombel, error)
	GetAllTahunPelajaran() ([]models.TahunPelajaran, error)
	GetAllTahunPelajaranAll() ([]models.TahunPelajaran, error)
	GetTotalSiswaByActiveTahunPelajaran() (int64, error)
	GetPesertaDidikByTahunPelajaran(tahunPelajaranID uint) ([]models.PesertaDidik, error)
	GetPesertaDidikByTahunPelajaranAndRombel(tahunPelajaranID uint, rombelID uint) ([]models.PesertaDidik, error)
	GetPesertaDidikByRombelID(rombelID uint) ([]models.PesertaDidik, error)
	GetPesertaDidikByRombelAndTahunPelajaran(rombelID uint, tahunPelajaranID uint, status string) ([]models.PesertaDidik, error)
	GetAllPesertaDidikActive() ([]models.PesertaDidik, error)
	UpdateBarcode(pesertaDidikID uint, barcode string) error
	UpdateWithTransaction(fn func(tx interface{}) error) error
	GetKepalaSekolah() (*models.Kepegawaian, error)
	GetVisiMisi() (*models.VisiMisi, error)
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

// GetByIDs retrieves multiple PesertaDidik by IDs
func (r *PesertaDidikRepositoryImpl) GetByIDs(ids []uint) ([]models.PesertaDidik, error) {
	var data []models.PesertaDidik
	if err := r.db.Where("id IN ?", ids).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetByIDWithDetails retrieves PesertaDidik by ID with roles preloaded
func (r *PesertaDidikRepositoryImpl) GetByIDWithDetails(id uint) (*models.PesertaDidik, error) {
	var data models.PesertaDidik
	if err := r.db.Preload("Roles.System").First(&data, id).Error; err != nil {
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

// GetByNISN retrieves PesertaDidik by NISN
func (r *PesertaDidikRepositoryImpl) GetByNISN(nisn string) (*models.PesertaDidik, error) {
	var data models.PesertaDidik
	if err := r.db.Where("nisn = ?", nisn).First(&data).Error; err != nil {
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

	// Get paginated data with preloaded relations, ordered by nama
	if err := r.db.Preload("Roles.System").
		Limit(limit).Offset(offset).Order("nama ASC").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetAllWithFilter retrieves PesertaDidik records with filters and pagination
// Filters: nama, nis, jenis_kelamin, nisn, tempat_lahir, nik, agama, status
// Order by: nama
func (r *PesertaDidikRepositoryImpl) GetAllWithFilter(params GetPesertaDidikParams) ([]models.PesertaDidik, int64, error) {
	var data []models.PesertaDidik
	var total int64

	query := r.db

	// Apply filters (hapus filter tahun_pelajaran_id dan rombel_id)
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

	// Get paginated data with preloaded relations, ordered by nama
	if err := query.Preload("Roles.System").
		Order("nama ASC").Limit(params.Limit).Offset(params.Offset).Find(&data).Error; err != nil {
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

// GetRombelByID retrieves Rombel by ID
func (r *PesertaDidikRepositoryImpl) GetRombelByID(id uint) (*models.Rombel, error) {
	var data models.Rombel
	if err := r.db.First(&data, id).Error; err != nil {
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

// GetTotalSiswaByActiveTahunPelajaran retrieves total count of peserta didik with active status only
func (r *PesertaDidikRepositoryImpl) GetTotalSiswaByActiveTahunPelajaran() (int64, error) {
	var total int64
	
	// Count peserta_didik where status = 'active' only
	err := r.db.Model(&models.PesertaDidik{}).
		Where("status = ?", "active").
		Count(&total).Error
	
	if err != nil {
		return 0, err
	}
	
	return total, nil
}

// GetPesertaDidikByTahunPelajaran retrieves all peserta didik by tahun pelajaran ID
func (r *PesertaDidikRepositoryImpl) GetPesertaDidikByTahunPelajaran(tahunPelajaranID uint) ([]models.PesertaDidik, error) {
	var data []models.PesertaDidik
	err := r.db.
		Distinct("peserta_didik.*").
		Joins("JOIN peserta_didik_rombel ON peserta_didik.id = peserta_didik_rombel.peserta_didik_id").
		Where("peserta_didik_rombel.tahun_pelajaran_id = ?", tahunPelajaranID).
		Where("peserta_didik_rombel.status = ?", "active").
		Where("peserta_didik.status = ?", "active").
		Order("peserta_didik.nis ASC").
		Find(&data).Error
	
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetPesertaDidikByTahunPelajaranAndRombel retrieves all peserta didik by tahun pelajaran ID and rombel ID
func (r *PesertaDidikRepositoryImpl) GetPesertaDidikByTahunPelajaranAndRombel(tahunPelajaranID uint, rombelID uint) ([]models.PesertaDidik, error) {
	var data []models.PesertaDidik
	err := r.db.
		Distinct("peserta_didik.*").
		Joins("JOIN peserta_didik_rombel ON peserta_didik.id = peserta_didik_rombel.peserta_didik_id").
		Where("peserta_didik_rombel.tahun_pelajaran_id = ? AND peserta_didik_rombel.rombel_id = ?", tahunPelajaranID, rombelID).
		Where("peserta_didik_rombel.status = ?", "active").
		Where("peserta_didik.status = ?", "active").
		Order("peserta_didik.nis ASC").
		Find(&data).Error
	
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetPesertaDidikByRombelID retrieves all peserta didik by rombel ID from peserta_didik_rombel table
func (r *PesertaDidikRepositoryImpl) GetPesertaDidikByRombelID(rombelID uint) ([]models.PesertaDidik, error) {
	var data []models.PesertaDidik
	err := r.db.
		Distinct("peserta_didik.*").
		Joins("JOIN peserta_didik_rombel ON peserta_didik.id = peserta_didik_rombel.peserta_didik_id").
		Where("peserta_didik_rombel.rombel_id = ? AND peserta_didik_rombel.status = ?", rombelID, "active").
		Where("peserta_didik.status = ?", "active").
		Order("peserta_didik.nis ASC").
		Find(&data).Error
	
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetPesertaDidikByRombelAndTahunPelajaran retrieves all peserta didik by rombel ID and tahun pelajaran ID from peserta_didik_rombel table
func (r *PesertaDidikRepositoryImpl) GetPesertaDidikByRombelAndTahunPelajaran(rombelID uint, tahunPelajaranID uint, status string) ([]models.PesertaDidik, error) {
	var data []models.PesertaDidik
	query := r.db.
		Distinct("peserta_didik.*").
		Joins("JOIN peserta_didik_rombel ON peserta_didik.id = peserta_didik_rombel.peserta_didik_id")
	
	// Apply filters
	if rombelID > 0 {
		query = query.Where("peserta_didik_rombel.rombel_id = ?", rombelID)
	}
	if tahunPelajaranID > 0 {
		query = query.Where("peserta_didik_rombel.tahun_pelajaran_id = ?", tahunPelajaranID)
	}
	if status != "" {
		query = query.Where("peserta_didik_rombel.status = ?", status)
	} else {
		// Default to active if no status specified
		query = query.Where("peserta_didik_rombel.status = ?", "active")
	}
	
	err := query.Order("peserta_didik.nis ASC").Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetAllPesertaDidikActive retrieves all active peserta didik
func (r *PesertaDidikRepositoryImpl) GetAllPesertaDidikActive() ([]models.PesertaDidik, error) {
	var data []models.PesertaDidik
	if err := r.db.Where("status = ?", "active").Order("nis ASC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateBarcode updates barcode and barcode_generated_at for a peserta didik
func (r *PesertaDidikRepositoryImpl) UpdateBarcode(pesertaDidikID uint, barcode string) error {
	return r.db.Model(&models.PesertaDidik{}).
		Where("id = ?", pesertaDidikID).
		Updates(map[string]interface{}{
			"barcode":             barcode,
			"barcode_generated_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

// UpdateWithTransaction executes a function within a database transaction
func (r *PesertaDidikRepositoryImpl) UpdateWithTransaction(fn func(tx interface{}) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// UpdateInTransaction updates a peserta didik record within a transaction
func (r *PesertaDidikRepositoryImpl) UpdateInTransaction(tx interface{}, data *models.PesertaDidik) error {
	txDB, ok := tx.(*gorm.DB)
	if !ok {
		return gorm.ErrInvalidTransaction
	}
	return txDB.Save(data).Error
}

// GetAllActive retrieves all peserta didik with status active
func (r *PesertaDidikRepositoryImpl) GetAllActive() ([]models.PesertaDidik, error) {
	var data []models.PesertaDidik
	if err := r.db.Where("status = ?", "active").Order("nama ASC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetKepalaSekolah retrieves Kepala Sekolah from kepegawaian table
func (r *PesertaDidikRepositoryImpl) GetKepalaSekolah() (*models.Kepegawaian, error) {
	var kepalaSekolah models.Kepegawaian
	if err := r.db.Where("jabatan = ?", "Kepala Sekolah").First(&kepalaSekolah).Error; err != nil {
		return nil, err
	}
	return &kepalaSekolah, nil
}

// GetVisiMisi retrieves the first visi misi record
func (r *PesertaDidikRepositoryImpl) GetVisiMisi() (*models.VisiMisi, error) {
	var visiMisi models.VisiMisi
	if err := r.db.First(&visiMisi).Error; err != nil {
		return nil, err
	}
	return &visiMisi, nil
}
