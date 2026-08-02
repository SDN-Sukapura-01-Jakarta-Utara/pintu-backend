package repositories

import (
	"fmt"
	"pintu-backend/src/modules/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GetFormulirFilter represents filter parameters for GetAllWithFilter
type GetFormulirFilter struct {
	Judul           string
	StartDate       time.Time
	EndDate         time.Time
	CreatedByID     *uint
	CreatedByRole   *string
	AccessType      *string
	TargetUserTypes []string
	RombelID        *int
}

// GetFormulirParams represents parameters for GetAllWithFilter with filters
type GetFormulirParams struct {
	Filter GetFormulirFilter
	Limit  int
	Offset int
}

// FormulirRepository handles data operations for Formulir
type FormulirRepository interface {
	Create(data *models.Formulir) error
	CreatePertanyaan(pertanyaan *models.FormulirPertanyaan) error
	GetByID(id uint) (*models.Formulir, error)
	GetBySlug(slug string) (*models.Formulir, error)
	Update(data *models.Formulir) error
	DeletePertanyaanByFormulirID(formulirID uint) error
	GetPertanyaanByFormulirID(formulirID uint) ([]models.FormulirPertanyaan, error)
	CheckSlugExists(slug string) error
	BeginTransaction() *gorm.DB
	GetAllWithFilter(params GetFormulirParams) ([]models.Formulir, int64, error)
	GetFormulirByUser(role string, rombelID *int, startDate, endDate time.Time, judul string, limit, offset int) ([]models.Formulir, int64, error)
	GetUserByID(userID uint) (*models.User, error)
	GetKepegawaianByUsername(username string) (*models.Kepegawaian, error)
	GetKepegawaianByID(id uint) (*models.Kepegawaian, error)
	GetPesertaDidikByUsername(username string) (*models.PesertaDidik, error)
	GetPesertaDidikByID(id uint) (*models.PesertaDidik, error)
	GetAllResponsesByFormulirID(formulirID uint) ([]models.FormulirResponse, error)
	GetJawabanByResponseID(responseID uint) ([]models.FormulirResponseJawaban, error)
}

type FormulirRepositoryImpl struct {
	db *gorm.DB
}

// NewFormulirRepository creates a new Formulir repository
func NewFormulirRepository(db *gorm.DB) FormulirRepository {
	return &FormulirRepositoryImpl{db: db}
}

// Create creates a new Formulir record
func (r *FormulirRepositoryImpl) Create(data *models.Formulir) error {
	return r.db.Create(data).Error
}

// CreatePertanyaan creates a new FormulirPertanyaan record
func (r *FormulirRepositoryImpl) CreatePertanyaan(pertanyaan *models.FormulirPertanyaan) error {
	return r.db.Create(pertanyaan).Error
}

// GetByID retrieves Formulir by ID with pertanyaan
func (r *FormulirRepositoryImpl) GetByID(id uint) (*models.Formulir, error) {
	var data models.Formulir
	if err := r.db.Preload("Pertanyaan", func(db *gorm.DB) *gorm.DB {
		return db.Order("urutan ASC")
	}).First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// GetBySlug retrieves Formulir by slug with pertanyaan
func (r *FormulirRepositoryImpl) GetBySlug(slug string) (*models.Formulir, error) {
	var data models.Formulir
	if err := r.db.Preload("Pertanyaan", func(db *gorm.DB) *gorm.DB {
		return db.Order("urutan ASC")
	}).Where("slug = ?", slug).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// Update updates Formulir record
func (r *FormulirRepositoryImpl) Update(data *models.Formulir) error {
	return r.db.Save(data).Error
}

// DeletePertanyaanByFormulirID deletes all pertanyaan for a formulir
func (r *FormulirRepositoryImpl) DeletePertanyaanByFormulirID(formulirID uint) error {
	return r.db.Where("form_id = ?", formulirID).Delete(&models.FormulirPertanyaan{}).Error
}

// GetPertanyaanByFormulirID gets all pertanyaan for a formulir
func (r *FormulirRepositoryImpl) GetPertanyaanByFormulirID(formulirID uint) ([]models.FormulirPertanyaan, error) {
	var pertanyaan []models.FormulirPertanyaan
	if err := r.db.Where("form_id = ?", formulirID).Order("urutan ASC").Find(&pertanyaan).Error; err != nil {
		return nil, err
	}
	return pertanyaan, nil
}

// CheckSlugExists checks if a slug already exists
func (r *FormulirRepositoryImpl) CheckSlugExists(slug string) error {
	var count int64
	if err := r.db.Model(&models.Formulir{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // Slug exists
	}
	return gorm.ErrRecordNotFound // Slug doesn't exist
}

// BeginTransaction starts a new database transaction
func (r *FormulirRepositoryImpl) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

// GetAllWithFilter retrieves all formulir with filters and pagination
func (r *FormulirRepositoryImpl) GetAllWithFilter(params GetFormulirParams) ([]models.Formulir, int64, error) {
	var formulirs []models.Formulir
	var total int64

	query := r.db.Model(&models.Formulir{})

	// Filter by judul (partial match)
	if params.Filter.Judul != "" {
		query = query.Where("judul ILIKE ?", "%"+params.Filter.Judul+"%")
	}

	// Filter by created_by_id
	if params.Filter.CreatedByID != nil {
		query = query.Where("created_by_user_id = ?", *params.Filter.CreatedByID)
	}

	// Filter by created_by_role
	if params.Filter.CreatedByRole != nil {
		if *params.Filter.CreatedByRole == "" {
			// Empty string means NULL (from users table)
			query = query.Where("created_by_role IS NULL")
		} else {
			// Filter by specific role
			query = query.Where("created_by_role = ?", *params.Filter.CreatedByRole)
		}
	}

	// Filter by start_date (created_at >= start_date)
	if !params.Filter.StartDate.IsZero() {
		query = query.Where("created_at >= ?", params.Filter.StartDate)
	}

	// Filter by end_date (created_at <= end_date)
	if !params.Filter.EndDate.IsZero() {
		query = query.Where("created_at <= ?", params.Filter.EndDate)
	}

	// Filter by access_type
	if params.Filter.AccessType != nil {
		query = query.Where("access_type = ?", *params.Filter.AccessType)
	}

	// Filter by target_user_types (check if JSONB array contains any of the specified types)
	if len(params.Filter.TargetUserTypes) > 0 {
		// For PostgreSQL JSONB array: check if any element in target_user_types matches
		// Build OR conditions for each target user type
		// Using JSONB @> operator to check if array contains element
		orConditions := make([]string, len(params.Filter.TargetUserTypes))
		orArgs := make([]interface{}, len(params.Filter.TargetUserTypes))
		
		for i, userType := range params.Filter.TargetUserTypes {
			orConditions[i] = "target_user_types @> ?"
			// Marshal to JSON array format: ["pendidik"]
			orArgs[i] = fmt.Sprintf(`["%s"]`, userType)
		}
		
		// Combine with OR
		query = query.Where(strings.Join(orConditions, " OR "), orArgs...)
	}

	// Filter by rombel_id (check if rombel_ids is NULL or contains the specified rombel_id)
	if params.Filter.RombelID != nil {
		// If rombel_ids is NULL (all rombels) OR contains the specified rombel_id
		query = query.Where("(rombel_ids IS NULL OR rombel_ids @> ?)", fmt.Sprintf(`[%d]`, *params.Filter.RombelID))
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	if err := query.Order("created_at DESC").
		Offset(params.Offset).
		Limit(params.Limit).
		Preload("CreatedBy").
		Find(&formulirs).Error; err != nil {
		return nil, 0, err
	}

	return formulirs, total, nil
}

// GetUserByID retrieves user by ID
func (r *FormulirRepositoryImpl) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetKepegawaianByUsername retrieves kepegawaian by username
func (r *FormulirRepositoryImpl) GetKepegawaianByUsername(username string) (*models.Kepegawaian, error) {
	var kepegawaian models.Kepegawaian
	if err := r.db.Where("username = ?", username).First(&kepegawaian).Error; err != nil {
		return nil, err
	}
	return &kepegawaian, nil
}

// GetKepegawaianByID retrieves kepegawaian by ID
func (r *FormulirRepositoryImpl) GetKepegawaianByID(id uint) (*models.Kepegawaian, error) {
	var kepegawaian models.Kepegawaian
	if err := r.db.Where("id = ?", id).First(&kepegawaian).Error; err != nil {
		return nil, err
	}
	return &kepegawaian, nil
}

// GetPesertaDidikByUsername retrieves peserta didik by username
func (r *FormulirRepositoryImpl) GetPesertaDidikByUsername(username string) (*models.PesertaDidik, error) {
	var pesertaDidik models.PesertaDidik
	if err := r.db.Where("username = ?", username).First(&pesertaDidik).Error; err != nil {
		return nil, err
	}
	return &pesertaDidik, nil
}

// GetPesertaDidikByID retrieves peserta didik by ID
func (r *FormulirRepositoryImpl) GetPesertaDidikByID(id uint) (*models.PesertaDidik, error) {
	var pesertaDidik models.PesertaDidik
	if err := r.db.Where("id = ?", id).First(&pesertaDidik).Error; err != nil {
		return nil, err
	}
	return &pesertaDidik, nil
}

// GetAllResponsesByFormulirID retrieves all responses for a formulir
func (r *FormulirRepositoryImpl) GetAllResponsesByFormulirID(formulirID uint) ([]models.FormulirResponse, error) {
	var responses []models.FormulirResponse
	if err := r.db.Where("form_id = ?", formulirID).Find(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

// GetJawabanByResponseID retrieves all jawaban for a response
func (r *FormulirRepositoryImpl) GetJawabanByResponseID(responseID uint) ([]models.FormulirResponseJawaban, error) {
	var jawaban []models.FormulirResponseJawaban
	if err := r.db.Where("response_id = ?", responseID).Find(&jawaban).Error; err != nil {
		return nil, err
	}
	return jawaban, nil
}

// GetFormulirByUser retrieves formulir based on user role and filters
func (r *FormulirRepositoryImpl) GetFormulirByUser(role string, rombelID *int, startDate, endDate time.Time, judul string, limit, offset int) ([]models.Formulir, int64, error) {
	var formulirs []models.Formulir
	var total int64

	query := r.db.Model(&models.Formulir{})

	// Filter: is_active must be true
	query = query.Where("is_active = ?", true)

	// Filter by target_user_types: must contain the specified role
	query = query.Where("target_user_types @> ?::jsonb", fmt.Sprintf(`["%s"]`, role))

	// Filter by rombel_id for murid role
	if role == "murid" && rombelID != nil {
		// Check if rombel_ids is NULL (all rombels) OR contains the specified rombel_id
		// Using PostgreSQL JSONB contains operator @>
		query = query.Where("(rombel_ids IS NULL OR rombel_ids::jsonb @> ?::jsonb)", fmt.Sprintf(`[%d]`, *rombelID))
	}

	// Filter by date range (form must be active during the date range)
	now := time.Now()
	query = query.Where("(start_date IS NULL OR start_date <= ?)", now)
	query = query.Where("(end_date IS NULL OR end_date >= ?)", now)

	// Additional filter by start_date if provided (filter by created_at)
	if !startDate.IsZero() {
		query = query.Where("created_at >= ?", startDate)
	}

	// Additional filter by end_date if provided (filter by created_at)
	if !endDate.IsZero() {
		query = query.Where("created_at <= ?", endDate)
	}

	// Filter by judul (partial match)
	if judul != "" {
		query = query.Where("judul ILIKE ?", "%"+judul+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Preload("CreatedBy").
		Find(&formulirs).Error; err != nil {
		return nil, 0, err
	}

	return formulirs, total, nil
}
