package repositories

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// LoginRepository handles data operations for authentication
type LoginRepository interface {
	GetByUsername(username string) (*models.User, error)
	GetKepegawaianByUsername(username string) (*models.Kepegawaian, error)
	GetPesertaDidikByUsername(username string) (*models.PesertaDidik, error)
	GetPelatihByUsername(username string) (*models.Pelatih, error)
	GetPesertaDidikRombelByStudentID(studentID uint) ([]models.PesertaDidikRombel, error)
}

type LoginRepositoryImpl struct {
	db *gorm.DB
}

// NewLoginRepository creates a new Login repository
func NewLoginRepository(db *gorm.DB) LoginRepository {
	return &LoginRepositoryImpl{db: db}
}

// GetByUsername retrieves user by username
func (r *LoginRepositoryImpl) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Roles.System").Preload("Roles.Permissions").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetKepegawaianByUsername retrieves kepegawaian by username
func (r *LoginRepositoryImpl) GetKepegawaianByUsername(username string) (*models.Kepegawaian, error) {
	var kepegawaian models.Kepegawaian
	if err := r.db.Preload("Roles.System").Preload("Roles.Permissions").Where("username = ?", username).First(&kepegawaian).Error; err != nil {
		return nil, err
	}
	return &kepegawaian, nil
}

// GetPesertaDidikByUsername retrieves peserta didik by username
func (r *LoginRepositoryImpl) GetPesertaDidikByUsername(username string) (*models.PesertaDidik, error) {
	var pesertaDidik models.PesertaDidik
	if err := r.db.Preload("Roles.System").Preload("Roles.Permissions").Where("username = ?", username).First(&pesertaDidik).Error; err != nil {
		return nil, err
	}
	return &pesertaDidik, nil
}

// GetPelatihByUsername retrieves pelatih by username
func (r *LoginRepositoryImpl) GetPelatihByUsername(username string) (*models.Pelatih, error) {
	var pelatih models.Pelatih
	if err := r.db.Preload("Roles.System").Preload("Roles.Permissions").Where("username = ?", username).First(&pelatih).Error; err != nil {
		return nil, err
	}
	return &pelatih, nil
}

// GetPesertaDidikRombelByStudentID retrieves all rombel for a student
func (r *LoginRepositoryImpl) GetPesertaDidikRombelByStudentID(studentID uint) ([]models.PesertaDidikRombel, error) {
	var rombelData []models.PesertaDidikRombel
	if err := r.db.Preload("Rombel.Kelas").Preload("TahunPelajaran").Where("peserta_didik_id = ?", studentID).Find(&rombelData).Error; err != nil {
		return nil, err
	}
	return rombelData, nil
}
