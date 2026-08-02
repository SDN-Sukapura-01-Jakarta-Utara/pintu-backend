package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// FormulirResponseRepository handles data operations for FormulirResponse
type FormulirResponseRepository interface {
	Create(response *models.FormulirResponse) error
	CreateJawaban(jawaban *models.FormulirResponseJawaban) error
	CountResponsesByFormulirID(formulirID uint) (int64, error)
	GetByFormulirIDAndUserID(formulirID uint, userID uint) (*models.FormulirResponse, error)
	GetByFormulirIDAndIPAddress(formulirID uint, ipAddress string) (*models.FormulirResponse, error)
	GetAllResponsesByFormulirID(formulirID uint) ([]models.FormulirResponse, error)
	GetResponseByFormulirIDAndUserID(formulirID uint, userID uint) (*models.FormulirResponse, error)
	GetKepegawaianByUsername(username string) (*models.Kepegawaian, error)
	GetPesertaDidikByUsername(username string) (*models.PesertaDidik, error)
	GetUserByID(userID uint) (*models.User, error)
	GetResponseByID(responseID uint) (*models.FormulirResponse, error)
	GetJawabanByResponseID(responseID uint) ([]models.FormulirResponseJawaban, error)
	UpdateResponse(response *models.FormulirResponse) error
	UpdateJawaban(jawaban *models.FormulirResponseJawaban) error
	DeleteJawaban(jawabanID uint) error
	BeginTransaction() *gorm.DB
}

type FormulirResponseRepositoryImpl struct {
	db *gorm.DB
}

// NewFormulirResponseRepository creates a new repository
func NewFormulirResponseRepository(db *gorm.DB) FormulirResponseRepository {
	return &FormulirResponseRepositoryImpl{db: db}
}

// Create creates a new response record
func (r *FormulirResponseRepositoryImpl) Create(response *models.FormulirResponse) error {
	return r.db.Create(response).Error
}

// CreateJawaban creates a new answer record
func (r *FormulirResponseRepositoryImpl) CreateJawaban(jawaban *models.FormulirResponseJawaban) error {
	return r.db.Create(jawaban).Error
}

// CountResponsesByFormulirID counts total responses for a form
func (r *FormulirResponseRepositoryImpl) CountResponsesByFormulirID(formulirID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.FormulirResponse{}).Where("form_id = ?", formulirID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetByFormulirIDAndUserID gets response by form ID and user ID (check if user already submitted)
func (r *FormulirResponseRepositoryImpl) GetByFormulirIDAndUserID(formulirID uint, userID uint) (*models.FormulirResponse, error) {
	var response models.FormulirResponse
	if err := r.db.Where("form_id = ? AND submitted_by_user_id = ?", formulirID, userID).First(&response).Error; err != nil {
		return nil, err
	}
	return &response, nil
}

// GetByFormulirIDAndIPAddress gets response by form ID and IP address (check if IP already submitted)
func (r *FormulirResponseRepositoryImpl) GetByFormulirIDAndIPAddress(formulirID uint, ipAddress string) (*models.FormulirResponse, error) {
	var response models.FormulirResponse
	if err := r.db.Where("form_id = ? AND ip_address = ?", formulirID, ipAddress).First(&response).Error; err != nil {
		return nil, err
	}
	return &response, nil
}

// GetAllResponsesByFormulirID gets all responses for a form with related data
func (r *FormulirResponseRepositoryImpl) GetAllResponsesByFormulirID(formulirID uint) ([]models.FormulirResponse, error) {
	var responses []models.FormulirResponse
	if err := r.db.Where("form_id = ?", formulirID).
		Preload("Jawaban.Pertanyaan").
		Preload("SubmittedBy").
		Order("submitted_at DESC").
		Find(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

// GetResponseByFormulirIDAndUserID gets a specific user's response for a form
func (r *FormulirResponseRepositoryImpl) GetResponseByFormulirIDAndUserID(formulirID uint, userID uint) (*models.FormulirResponse, error) {
	var response models.FormulirResponse
	if err := r.db.Where("form_id = ? AND submitted_by_user_id = ?", formulirID, userID).
		Preload("Jawaban.Pertanyaan").
		Preload("SubmittedBy").
		First(&response).Error; err != nil {
		return nil, err
	}
	return &response, nil
}

// GetKepegawaianByUsername retrieves kepegawaian by username
func (r *FormulirResponseRepositoryImpl) GetKepegawaianByUsername(username string) (*models.Kepegawaian, error) {
	var kepegawaian models.Kepegawaian
	if err := r.db.Where("username = ?", username).First(&kepegawaian).Error; err != nil {
		return nil, err
	}
	return &kepegawaian, nil
}

// GetPesertaDidikByUsername retrieves peserta didik by username
func (r *FormulirResponseRepositoryImpl) GetPesertaDidikByUsername(username string) (*models.PesertaDidik, error) {
	var pesertaDidik models.PesertaDidik
	if err := r.db.Where("username = ?", username).First(&pesertaDidik).Error; err != nil {
		return nil, err
	}
	return &pesertaDidik, nil
}

// GetUserByID retrieves user by ID
func (r *FormulirResponseRepositoryImpl) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetResponseByID retrieves response by ID with jawaban
func (r *FormulirResponseRepositoryImpl) GetResponseByID(responseID uint) (*models.FormulirResponse, error) {
	var response models.FormulirResponse
	if err := r.db.Preload("Jawaban").Preload("SubmittedBy").Where("id = ?", responseID).First(&response).Error; err != nil {
		return nil, err
	}
	return &response, nil
}

// GetJawabanByResponseID retrieves all jawaban for a response
func (r *FormulirResponseRepositoryImpl) GetJawabanByResponseID(responseID uint) ([]models.FormulirResponseJawaban, error) {
	var jawaban []models.FormulirResponseJawaban
	if err := r.db.Where("response_id = ?", responseID).Find(&jawaban).Error; err != nil {
		return nil, err
	}
	return jawaban, nil
}

// UpdateResponse updates a response record
func (r *FormulirResponseRepositoryImpl) UpdateResponse(response *models.FormulirResponse) error {
	return r.db.Save(response).Error
}

// UpdateJawaban updates a jawaban record
func (r *FormulirResponseRepositoryImpl) UpdateJawaban(jawaban *models.FormulirResponseJawaban) error {
	return r.db.Save(jawaban).Error
}

// DeleteJawaban deletes a jawaban record
func (r *FormulirResponseRepositoryImpl) DeleteJawaban(jawabanID uint) error {
	return r.db.Delete(&models.FormulirResponseJawaban{}, jawabanID).Error
}

// BeginTransaction starts a new database transaction
func (r *FormulirResponseRepositoryImpl) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}
