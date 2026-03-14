package services

import (
	"errors"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"golang.org/x/crypto/bcrypt"
)

type PesertaDidikService interface {
	Create(req *dtos.PesertaDidikCreateRequest, userID uint) (*dtos.PesertaDidikResponse, error)
	GetByID(id uint) (*dtos.PesertaDidikResponse, error)
	GetByNIS(nis string) (*dtos.PesertaDidikResponse, error)
	GetAll(limit int, offset int) (*dtos.PesertaDidikListResponse, error)
	GetAllWithFilter(params repositories.GetPesertaDidikParams) (*dtos.PesertaDidikListWithPaginationResponse, error)
	Update(id uint, req *dtos.PesertaDidikUpdateRequest, userID uint) (*dtos.PesertaDidikResponse, error)
	Delete(id uint) error
}

type PesertaDidikServiceImpl struct {
	repository repositories.PesertaDidikRepository
}

// NewPesertaDidikService creates a new PesertaDidik service
func NewPesertaDidikService(repository repositories.PesertaDidikRepository) PesertaDidikService {
	return &PesertaDidikServiceImpl{
		repository: repository,
	}
}

// Create creates a new PesertaDidik
func (s *PesertaDidikServiceImpl) Create(req *dtos.PesertaDidikCreateRequest, userID uint) (*dtos.PesertaDidikResponse, error) {
	// Check if NIS already exists in same tahun_pelajaran
	existing, _ := s.repository.GetByNISAndTahunPelajaran(req.NIS, req.TahunPelajaranID)
	if existing != nil {
		return nil, errors.New("NIS sudah ada di tahun pelajaran ini")
	}

	// Check if username already exists in same tahun_pelajaran (if provided)
	if req.Username != "" {
		existing, _ := s.repository.GetByUsernameAndTahunPelajaran(req.Username, req.TahunPelajaranID)
		if existing != nil {
			return nil, errors.New("username sudah ada di tahun pelajaran ini")
		}
	}

	// Hash password if provided
	var hashedPassword string
	if req.Password != "" {
		hp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("gagal hash password")
		}
		hashedPassword = string(hp)
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Parse tanggal_lahir (YYYY-MM-DD format)
	var tanggalLahir *time.Time
	if req.TanggalLahir != "" {
		t, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err != nil {
			return nil, errors.New("format tanggal_lahir tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalLahir = &t
	}

	// Create peserta didik record
	data := &models.PesertaDidik{
		Nama:               req.Nama,
		NIS:                req.NIS,
		JenisKelamin:       req.JenisKelamin,
		NISN:               req.NISN,
		TempatLahir:        req.TempatLahir,
		TanggalLahir:       tanggalLahir,
		NIK:                req.NIK,
		Agama:              req.Agama,
		Alamat:             req.Alamat,
		RT:                 req.RT,
		RW:                 req.RW,
		Kelurahan:          req.Kelurahan,
		Kecamatan:          req.Kecamatan,
		KodePos:            req.KodePos,
		NamaAyah:           req.NamaAyah,
		NamaIbu:            req.NamaIbu,
		RombelID:           req.RombelID,
		TahunPelajaranID:   req.TahunPelajaranID,
		Status:             status,
		Username:           req.Username,
		Password:           hashedPassword,
		CreatedByID:        &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	// Assign roles
	if err := s.repository.AssignRoles(data.ID, req.RoleIDs); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves PesertaDidik by ID with details
func (s *PesertaDidikServiceImpl) GetByID(id uint) (*dtos.PesertaDidikResponse, error) {
	data, err := s.repository.GetByIDWithDetails(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetByNIS retrieves PesertaDidik by NIS
func (s *PesertaDidikServiceImpl) GetByNIS(nis string) (*dtos.PesertaDidikResponse, error) {
	data, err := s.repository.GetByNIS(nis)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all PesertaDidik
func (s *PesertaDidikServiceImpl) GetAll(limit int, offset int) (*dtos.PesertaDidikListResponse, error) {
	// Set default limit and offset
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	data, total, err := s.repository.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	// Map to response
	responses := make([]dtos.PesertaDidikResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.PesertaDidikListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// GetAllWithFilter retrieves PesertaDidik with filters and pagination
func (s *PesertaDidikServiceImpl) GetAllWithFilter(params repositories.GetPesertaDidikParams) (*dtos.PesertaDidikListWithPaginationResponse, error) {
	// Validate and set default limit and offset
	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	data, total, err := s.repository.GetAllWithFilter(params)
	if err != nil {
		return nil, err
	}

	// Map to response
	responses := make([]dtos.PesertaDidikResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.PesertaDidikListWithPaginationResponse{
		Data: responses,
		Pagination: dtos.PaginationInfo{
			Limit:      params.Limit,
			Offset:     params.Offset,
			Page:       (params.Offset / params.Limit) + 1,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// Update updates PesertaDidik
func (s *PesertaDidikServiceImpl) Update(id uint, req *dtos.PesertaDidikUpdateRequest, userID uint) (*dtos.PesertaDidikResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("peserta didik tidak ditemukan")
	}

	// Update basic fields if provided
	if req.Nama != "" {
		existing.Nama = req.Nama
	}
	if req.JenisKelamin != "" {
		existing.JenisKelamin = req.JenisKelamin
	}
	if req.NIS != "" && req.NIS != existing.NIS {
		// Check if new NIS already exists for other peserta didik in same tahun_pelajaran
		// Only check if NIS actually changed
		tahunPelajaranID := existing.TahunPelajaranID
		if existing.TahunPelajaranID == nil && req.TahunPelajaranID == nil {
			// Skip check if both are null
		} else if existing.TahunPelajaranID != nil && req.TahunPelajaranID == nil {
			tahunPelajaranID = existing.TahunPelajaranID
		} else if existing.TahunPelajaranID == nil && req.TahunPelajaranID != nil {
			tahunPelajaranID = req.TahunPelajaranID
		}
		
		existingNIS, _ := s.repository.GetByNISAndTahunPelajaran(req.NIS, tahunPelajaranID)
		if existingNIS != nil {
			return nil, errors.New("NIS sudah ada di tahun pelajaran ini")
		}
		existing.NIS = req.NIS
	}
	if req.NISN != "" {
		existing.NISN = req.NISN
	}
	if req.TempatLahir != "" {
		existing.TempatLahir = req.TempatLahir
	}
	if req.TanggalLahir != "" {
		t, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err != nil {
			return nil, errors.New("format tanggal_lahir tidak valid, gunakan YYYY-MM-DD")
		}
		existing.TanggalLahir = &t
	}
	if req.NIK != "" {
		existing.NIK = req.NIK
	}
	if req.Agama != "" {
		existing.Agama = req.Agama
	}
	if req.Alamat != "" {
		existing.Alamat = req.Alamat
	}
	if req.RT != "" {
		existing.RT = req.RT
	}
	if req.RW != "" {
		existing.RW = req.RW
	}
	if req.Kelurahan != "" {
		existing.Kelurahan = req.Kelurahan
	}
	if req.Kecamatan != "" {
		existing.Kecamatan = req.Kecamatan
	}
	if req.KodePos != "" {
		existing.KodePos = req.KodePos
	}
	if req.NamaAyah != "" {
		existing.NamaAyah = req.NamaAyah
	}
	if req.NamaIbu != "" {
		existing.NamaIbu = req.NamaIbu
	}
	if req.RombelID != nil {
		existing.RombelID = req.RombelID
	}
	if req.TahunPelajaranID != nil {
		existing.TahunPelajaranID = req.TahunPelajaranID
	}
	if req.Username != "" && req.Username != existing.Username {
		// Check if new username already exists for other peserta didik in same tahun_pelajaran
		// Only check if username actually changed
		tahunPelajaranID := existing.TahunPelajaranID
		if existing.TahunPelajaranID == nil && req.TahunPelajaranID == nil {
			// Skip check if both are null
		} else if existing.TahunPelajaranID != nil && req.TahunPelajaranID == nil {
			tahunPelajaranID = existing.TahunPelajaranID
		} else if existing.TahunPelajaranID == nil && req.TahunPelajaranID != nil {
			tahunPelajaranID = req.TahunPelajaranID
		}
		
		existingUsername, _ := s.repository.GetByUsernameAndTahunPelajaran(req.Username, tahunPelajaranID)
		if existingUsername != nil {
			return nil, errors.New("username sudah ada di tahun pelajaran ini")
		}
		existing.Username = req.Username
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("gagal hash password")
		}
		existing.Password = string(hashedPassword)
	}
	if req.Status != "" {
		existing.Status = req.Status
	}

	existing.UpdatedByID = &userID

	// Update record
	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	// Update roles
	if len(req.RoleIDs) > 0 {
		if err := s.repository.AssignRoles(existing.ID, req.RoleIDs); err != nil {
			return nil, errors.New("gagal mengupdate roles peserta didik")
		}
	} else {
		// If no roles provided, remove all roles
		if err := s.repository.RemoveRoles(existing.ID); err != nil {
			return nil, errors.New("gagal menghapus roles peserta didik")
		}
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes PesertaDidik by ID
func (s *PesertaDidikServiceImpl) Delete(id uint) error {
	// Validate peserta didik exists before delete
	data, err := s.repository.GetByID(id)
	if err != nil || data == nil {
		return errors.New("peserta didik tidak ditemukan atau sudah dihapus")
	}

	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *PesertaDidikServiceImpl) mapToResponse(data *models.PesertaDidik) *dtos.PesertaDidikResponse {
	// Map roles
	roles := make([]dtos.RoleResponse, len(data.Roles))
	for i, role := range data.Roles {
		var system *dtos.SystemResponse
		if role.System != nil {
			system = &dtos.SystemResponse{
				ID:          role.System.ID,
				Nama:        role.System.Nama,
				Description: role.System.Description,
			}
		}

		roles[i] = dtos.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			SystemID:    role.SystemID,
			System:      system,
			Status:      role.Status,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
			CreatedByID: role.CreatedByID,
			UpdatedByID: role.UpdatedByID,
		}
	}

	// Map rombel
	var rombel *dtos.RombelDetailResponse
	if data.Rombel != nil {
		var kelas *dtos.KelasDetailResponse
		if data.Rombel.Kelas != nil {
			kelas = &dtos.KelasDetailResponse{
				ID:        data.Rombel.Kelas.ID,
				Name:      data.Rombel.Kelas.Name,
				Status:    data.Rombel.Kelas.Status,
				CreatedAt: data.Rombel.Kelas.CreatedAt.Format("2006-01-02T15:04:05Z"),
				UpdatedAt: data.Rombel.Kelas.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			}
		}
		rombel = &dtos.RombelDetailResponse{
			ID:        data.Rombel.ID,
			Name:      data.Rombel.Name,
			Status:    data.Rombel.Status,
			KelasID:   data.Rombel.KelasID,
			Kelas:     kelas,
			CreatedAt: data.Rombel.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: data.Rombel.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	// Map tahun pelajaran
	var tahunPelajaran *dtos.TahunPelajaranDetailResponse
	if data.TahunPelajaran != nil {
		tahunPelajaran = &dtos.TahunPelajaranDetailResponse{
			ID:             data.TahunPelajaran.ID,
			TahunPelajaran: data.TahunPelajaran.TahunPelajaran,
			Status:         data.TahunPelajaran.Status,
			CreatedAt:      data.TahunPelajaran.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      data.TahunPelajaran.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			CreatedByID:    data.TahunPelajaran.CreatedByID,
			UpdatedByID:    data.TahunPelajaran.UpdatedByID,
		}
	}

	tanggalLahirStr := ""
	if data.TanggalLahir != nil {
		tanggalLahirStr = data.TanggalLahir.Format("2006-01-02")
	}

	return &dtos.PesertaDidikResponse{
		ID:               data.ID,
		Nama:             data.Nama,
		NIS:              data.NIS,
		JenisKelamin:     data.JenisKelamin,
		NISN:             data.NISN,
		TempatLahir:      data.TempatLahir,
		TanggalLahir:     tanggalLahirStr,
		NIK:              data.NIK,
		Agama:            data.Agama,
		Alamat:           data.Alamat,
		RT:               data.RT,
		RW:               data.RW,
		Kelurahan:        data.Kelurahan,
		Kecamatan:        data.Kecamatan,
		KodePos:          data.KodePos,
		NamaAyah:         data.NamaAyah,
		NamaIbu:          data.NamaIbu,
		RombelID:         data.RombelID,
		Rombel:           rombel,
		TahunPelajaranID: data.TahunPelajaranID,
		TahunPelajaran:   tahunPelajaran,
		Status:           data.Status,
		Username:         data.Username,
		Roles:            roles,
		CreatedAt:        data.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        data.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		CreatedByID:      data.CreatedByID,
		UpdatedByID:      data.UpdatedByID,
	}
}
