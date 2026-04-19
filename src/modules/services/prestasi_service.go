package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

type PrestasiService interface {
	Create(foto []*multipart.FileHeader, fotoThumbnails []string, req *dtos.PrestasiCreateRequest, userID uint) (*dtos.PrestasiResponse, error)
	GetByID(id uint) (*dtos.PrestasiResponse, error)
	GetAll(limit int, offset int) (*dtos.PrestasiListResponse, error)
	GetAllWithFilter(params repositories.GetPrestasiParams) (*dtos.PrestasiListWithPaginationResponse, error)
	Update(id uint, foto []*multipart.FileHeader, fotoThumbnails []string, req *dtos.PrestasiUpdateRequest, userID uint) (*dtos.PrestasiResponse, error)
	Delete(id uint) error
}

type PrestasiServiceImpl struct {
	repository repositories.PrestasiRepository
	r2Storage  *utils.R2Storage
}

// NewPrestasiService creates a new Prestasi service
func NewPrestasiService(repository repositories.PrestasiRepository, r2Storage *utils.R2Storage) PrestasiService {
	return &PrestasiServiceImpl{
		repository: repository,
		r2Storage:  r2Storage,
	}
}

// Create creates a new Prestasi with foto uploads to R2
func (s *PrestasiServiceImpl) Create(foto []*multipart.FileHeader, fotoThumbnails []string, req *dtos.PrestasiCreateRequest, userID uint) (*dtos.PrestasiResponse, error) {
	// Parse tanggal lomba
	tanggalLomba, err := time.Parse("2006-01-02", req.TanggalLomba)
	if err != nil {
		return nil, errors.New("invalid tanggal_lomba format (use YYYY-MM-DD)")
	}

	// Upload foto if provided
	var fotoItems []models.FotoItem
	if len(foto) > 0 {
		for i, file := range foto {
			if file == nil {
				continue
			}

			// Validate foto file
			if file.Size > 5*1024*1024 { // 5MB
				return nil, errors.New("each foto must not exceed 5MB")
			}

			allowedTypes := map[string]bool{
				"image/jpeg": true,
				"image/png":  true,
				"image/gif":  true,
				"image/webp": true,
			}
			contentType := file.Header.Get("Content-Type")
			if !allowedTypes[contentType] {
				return nil, errors.New("only image files are allowed for foto (jpeg, png, gif, webp)")
			}

			// Upload foto to R2
			fileKey, err := s.r2Storage.UploadFile(file, "prestasi")
			if err != nil {
				return nil, err
			}

			// Generate unique foto ID
			fotoID := fmt.Sprintf("foto_%d_%s", time.Now().UnixNano(), fileKey[len(fileKey)-8:])

			// Set thumbnail status based on fotoThumbnails array
			thumbnail := "inactive"
			if len(fotoThumbnails) > i && fotoThumbnails[i] == "active" {
				thumbnail = "active"
			} else if i == 0 && len(fotoThumbnails) == 0 {
				// Default: first image is active if no thumbnails specified
				thumbnail = "active"
			}

			fotoItems = append(fotoItems, models.FotoItem{
				ID:        fotoID,
				Filename:  file.Filename,
				URL:       fileKey,
				Size:      file.Size,
				Thumbnail: thumbnail,
			})
		}
	}

	// Convert fotoItems to JSON
	fotoJSON, _ := json.Marshal(fotoItems)

	// Create prestasi record
	data := &models.Prestasi{
		PesertaDidikID:    req.PesertaDidikID,
		Jenis:             req.Jenis,
		NamaGrup:          req.NamaGrup,
		NamaPrestasi:      req.NamaPrestasi,
		TingkatPrestasi:   req.TingkatPrestasi,
		Penyelenggara:     req.Penyelenggara,
		TanggalLomba:      tanggalLomba,
		Juara:             req.Juara,
		Keterangan:        req.Keterangan,
		Foto:              fotoJSON,
		EkstrakurikulerID: req.EkstrakurikulerID,
		TahunPelajaranID:  req.TahunPelajaranID,
		CreatedByID:       &userID,
	}

	if err := s.repository.Create(data); err != nil {
		// If database save fails, delete uploaded files
		for _, item := range fotoItems {
			_ = s.r2Storage.DeleteFile(item.URL)
		}
		return nil, err
	}

	// Create anggota tim if provided
	if len(req.AnggotaTim) > 0 {
		for _, anggota := range req.AnggotaTim {
			tahunPelajaranID := anggota.TahunPelajaranID
			if tahunPelajaranID == 0 {
				tahunPelajaranID = data.TahunPelajaranID
			}
			anggotaData := &models.AnggotaTimPrestasi{
				PrestasiID:       data.ID,
				PesertaDidikID:   anggota.PesertaDidikID,
				TahunPelajaranID: tahunPelajaranID,
				CreatedByID:      &userID,
			}
			if err := s.repository.CreateAnggotaTim(anggotaData); err != nil {
				return nil, fmt.Errorf("failed to create anggota tim: %w", err)
			}
		}
	}

	// Get the created data with relationships
	createdData, err := s.repository.GetByID(data.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(createdData), nil
}

// GetByID retrieves Prestasi by ID
func (s *PrestasiServiceImpl) GetByID(id uint) (*dtos.PrestasiResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all Prestasi
func (s *PrestasiServiceImpl) GetAll(limit int, offset int) (*dtos.PrestasiListResponse, error) {
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
	responses := make([]dtos.PrestasiResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.PrestasiListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// GetAllWithFilter retrieves Prestasi with filters and pagination
func (s *PrestasiServiceImpl) GetAllWithFilter(params repositories.GetPrestasiParams) (*dtos.PrestasiListWithPaginationResponse, error) {
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
	responses := make([]dtos.PrestasiResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.PrestasiListWithPaginationResponse{
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

// Update updates Prestasi
func (s *PrestasiServiceImpl) Update(id uint, foto []*multipart.FileHeader, fotoThumbnails []string, req *dtos.PrestasiUpdateRequest, userID uint) (*dtos.PrestasiResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("prestasi not found")
	}

	// Clear preloaded associations to prevent GORM from overriding FK values on Save
	existing.PesertaDidik = nil
	existing.Ekstrakurikuler = nil
	existing.TahunPelajaran = nil
	existing.AnggotaTimPrestasi = nil

	// Update basic fields if provided
	if req.PesertaDidikID != nil {
		existing.PesertaDidikID = req.PesertaDidikID
	}
	if req.Jenis != "" {
		existing.Jenis = req.Jenis
	}
	if req.NamaGrup != "" {
		existing.NamaGrup = req.NamaGrup
	}
	if req.NamaPrestasi != "" {
		existing.NamaPrestasi = req.NamaPrestasi
	}
	if req.TingkatPrestasi != "" {
		existing.TingkatPrestasi = req.TingkatPrestasi
	}
	if req.Penyelenggara != "" {
		existing.Penyelenggara = req.Penyelenggara
	}
	if req.TanggalLomba != "" {
		tanggalLomba, err := time.Parse("2006-01-02", req.TanggalLomba)
		if err != nil {
			return nil, errors.New("invalid tanggal_lomba format (use YYYY-MM-DD)")
		}
		existing.TanggalLomba = tanggalLomba
	}
	if req.Juara != "" {
		existing.Juara = req.Juara
	}
	if req.Keterangan != "" {
		existing.Keterangan = req.Keterangan
	}
	if req.EkstrakurikulerID != nil {
		existing.EkstrakurikulerID = req.EkstrakurikulerID
	}
	if req.TahunPelajaranID != 0 {
		existing.TahunPelajaranID = req.TahunPelajaranID
	}

	// Delete foto if specified
	if len(req.FotoToDelete) > 0 {
		var existingFotoItems []models.FotoItem
		if err := json.Unmarshal(existing.Foto, &existingFotoItems); err == nil {
			// Build map of foto IDs to delete
			deleteMap := make(map[string]bool)
			for _, fotoID := range req.FotoToDelete {
				deleteMap[fotoID] = true
			}

			// Filter out foto to delete and delete from R2
			var remainingFoto []models.FotoItem
			for _, fotoItem := range existingFotoItems {
				if deleteMap[fotoItem.ID] {
					// Delete from R2
					_ = s.r2Storage.DeleteFile(fotoItem.URL)
				} else {
					remainingFoto = append(remainingFoto, fotoItem)
				}
			}

			// Update foto array
			fotoJSON, _ := json.Marshal(remainingFoto)
			existing.Foto = fotoJSON
		}
	}

	// Add new foto if provided
	if len(foto) > 0 {
		var existingFotoItems []models.FotoItem
		// Get existing foto
		if err := json.Unmarshal(existing.Foto, &existingFotoItems); err != nil && existing.Foto != nil && len(existing.Foto) > 0 {
			// Only log, don't fail
		}

		// Upload and add new foto
		for i, file := range foto {
			if file == nil {
				continue
			}

			// Validate foto file
			if file.Size > 5*1024*1024 { // 5MB
				return nil, errors.New("each foto must not exceed 5MB")
			}

			allowedTypes := map[string]bool{
				"image/jpeg": true,
				"image/png":  true,
				"image/gif":  true,
				"image/webp": true,
			}
			contentType := file.Header.Get("Content-Type")
			if !allowedTypes[contentType] {
				return nil, errors.New("only image files are allowed for foto (jpeg, png, gif, webp)")
			}

			// Upload foto to R2
			fileKey, err := s.r2Storage.UploadFile(file, "prestasi")
			if err != nil {
				return nil, err
			}

			// Generate unique foto ID
			fotoID := fmt.Sprintf("foto_%d_%s", time.Now().UnixNano(), fileKey[len(fileKey)-8:])

			// Set thumbnail status based on fotoThumbnails array
			thumbnail := "inactive"
			if len(fotoThumbnails) > i && fotoThumbnails[i] == "active" {
				thumbnail = "active"
			} else if len(existingFotoItems) == 0 && i == 0 && len(fotoThumbnails) == 0 {
				// Default: first foto becomes active if no existing foto and no thumbnails specified
				thumbnail = "active"
			}

			existingFotoItems = append(existingFotoItems, models.FotoItem{
				ID:        fotoID,
				Filename:  file.Filename,
				URL:       fileKey,
				Size:      file.Size,
				Thumbnail: thumbnail,
			})
		}

		// Convert updated foto to JSON
		fotoJSON, _ := json.Marshal(existingFotoItems)
		existing.Foto = fotoJSON
	}

	// Update anggota tim if provided
	if len(req.AnggotaTim) > 0 {
		// Delete existing anggota tim
		_ = s.repository.DeleteAnggotaTimByPrestasiID(id)

		// Create new anggota tim
		for _, anggota := range req.AnggotaTim {
			tahunPelajaranID := anggota.TahunPelajaranID
			if tahunPelajaranID == 0 {
				tahunPelajaranID = existing.TahunPelajaranID
			}
			anggotaData := &models.AnggotaTimPrestasi{
				PrestasiID:       id,
				PesertaDidikID:   anggota.PesertaDidikID,
				TahunPelajaranID: tahunPelajaranID,
				CreatedByID:      &userID,
			}
			if err := s.repository.CreateAnggotaTim(anggotaData); err != nil {
				return nil, fmt.Errorf("failed to create anggota tim: %w", err)
			}
		}
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	// Get updated data with relationships
	updatedData, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(updatedData), nil
}

// Delete deletes Prestasi by ID
func (s *PrestasiServiceImpl) Delete(id uint) error {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return errors.New("prestasi not found")
	}

	// Delete all foto from R2
	var fotoItems []models.FotoItem
	if err := json.Unmarshal(existing.Foto, &fotoItems); err == nil {
		for _, fotoItem := range fotoItems {
			_ = s.r2Storage.DeleteFile(fotoItem.URL)
		}
	}

	// Delete anggota tim
	_ = s.repository.DeleteAnggotaTimByPrestasiID(id)

	// Delete from database
	return s.repository.Delete(id)
}

// mapToResponse maps model to DTO response
func (s *PrestasiServiceImpl) mapToResponse(data *models.Prestasi) *dtos.PrestasiResponse {
	// Map foto from JSON
	fotoItems := make([]dtos.FotoItemDTO, 0)
	var fotoModels []models.FotoItem
	if err := json.Unmarshal(data.Foto, &fotoModels); err == nil {
		for _, fotoItem := range fotoModels {
			fotoItems = append(fotoItems, dtos.FotoItemDTO{
				ID:        fotoItem.ID,
				Filename:  fotoItem.Filename,
				URL:       s.r2Storage.GetPublicURL(fotoItem.URL),
				Size:      fotoItem.Size,
				Thumbnail: fotoItem.Thumbnail,
			})
		}
	}

	// Map peserta didik
	var pesertaDidik *dtos.PesertaDidikResponse
	if data.PesertaDidik != nil {
		pesertaDidik = s.mapPesertaDidikToResponse(data.PesertaDidik)
	}

	// Map tahun pelajaran
	var tahunPelajaran *dtos.TahunPelajaranDetailResponse
	if data.TahunPelajaran != nil {
		tahunPelajaran = &dtos.TahunPelajaranDetailResponse{
			ID:             data.TahunPelajaran.ID,
			TahunPelajaran: data.TahunPelajaran.TahunPelajaran,
			Status:         data.TahunPelajaran.Status,
			CreatedAt:      data.TahunPelajaran.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      data.TahunPelajaran.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedByID:    data.TahunPelajaran.CreatedByID,
			UpdatedByID:    data.TahunPelajaran.UpdatedByID,
		}
	}

	// Map ekstrakurikuler
	var ekstrakurikuler *dtos.EkstrakurikulerDetailDTO
	if data.Ekstrakurikuler != nil {
		ekstrakurikuler = &dtos.EkstrakurikulerDetailDTO{
			ID:       data.Ekstrakurikuler.ID,
			Name:     data.Ekstrakurikuler.Name,
			Kategori: data.Ekstrakurikuler.Kategori,
			Status:   data.Ekstrakurikuler.Status,
		}
	}

	// Map anggota tim prestasi
	anggotaTimPrestasi := make([]dtos.AnggotaTimPrestasiDTO, 0)
	for _, anggota := range data.AnggotaTimPrestasi {
		var anggotaPesertaDidik *dtos.PesertaDidikResponse
		if anggota.PesertaDidik != nil {
			anggotaPesertaDidik = s.mapPesertaDidikToResponse(anggota.PesertaDidik)
		}

		var anggotaTahunPelajaran *dtos.TahunPelajaranDetailResponse
		if anggota.TahunPelajaran != nil {
			anggotaTahunPelajaran = &dtos.TahunPelajaranDetailResponse{
				ID:             anggota.TahunPelajaran.ID,
				TahunPelajaran: anggota.TahunPelajaran.TahunPelajaran,
				Status:         anggota.TahunPelajaran.Status,
				CreatedAt:      anggota.TahunPelajaran.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt:      anggota.TahunPelajaran.UpdatedAt.Format("2006-01-02 15:04:05"),
				CreatedByID:    anggota.TahunPelajaran.CreatedByID,
				UpdatedByID:    anggota.TahunPelajaran.UpdatedByID,
			}
		}

		anggotaTimPrestasi = append(anggotaTimPrestasi, dtos.AnggotaTimPrestasiDTO{
			ID:               anggota.ID,
			PrestasiID:       anggota.PrestasiID,
			PesertaDidikID:   anggota.PesertaDidikID,
			TahunPelajaranID: anggota.TahunPelajaranID,
			PesertaDidik:     anggotaPesertaDidik,
			TahunPelajaran:   anggotaTahunPelajaran,
			CreatedAt:        anggota.CreatedAt,
			UpdatedAt:        anggota.UpdatedAt,
		})
	}

	return &dtos.PrestasiResponse{
		ID:                 data.ID,
		PesertaDidikID:     data.PesertaDidikID,
		PesertaDidik:       pesertaDidik,
		Jenis:              data.Jenis,
		NamaGrup:           data.NamaGrup,
		NamaPrestasi:       data.NamaPrestasi,
		TingkatPrestasi:    data.TingkatPrestasi,
		Penyelenggara:      data.Penyelenggara,
		TanggalLomba:       data.TanggalLomba,
		Juara:              data.Juara,
		Keterangan:         data.Keterangan,
		Foto:               fotoItems,
		EkstrakurikulerID:  data.EkstrakurikulerID,
		Ekstrakurikuler:    ekstrakurikuler,
		TahunPelajaranID:   data.TahunPelajaranID,
		TahunPelajaran:     tahunPelajaran,
		AnggotaTimPrestasi: anggotaTimPrestasi,
		CreatedAt:          data.CreatedAt,
		UpdatedAt:          data.UpdatedAt,
		CreatedByID:        data.CreatedByID,
		UpdatedByID:        data.UpdatedByID,
	}
}

// mapPesertaDidikToResponse maps PesertaDidik model to response DTO
func (s *PrestasiServiceImpl) mapPesertaDidikToResponse(data *models.PesertaDidik) *dtos.PesertaDidikResponse {
	var tanggalLahir string
	if data.TanggalLahir != nil {
		tanggalLahir = data.TanggalLahir.Format("2006-01-02")
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
				CreatedAt: data.Rombel.Kelas.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: data.Rombel.Kelas.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		rombel = &dtos.RombelDetailResponse{
			ID:        data.Rombel.ID,
			Name:      data.Rombel.Name,
			Status:    data.Rombel.Status,
			KelasID:   data.Rombel.KelasID,
			Kelas:     kelas,
			CreatedAt: data.Rombel.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: data.Rombel.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// Map tahun pelajaran
	var tahunPelajaran *dtos.TahunPelajaranDetailResponse
	if data.TahunPelajaran != nil {
		tahunPelajaran = &dtos.TahunPelajaranDetailResponse{
			ID:             data.TahunPelajaran.ID,
			TahunPelajaran: data.TahunPelajaran.TahunPelajaran,
			Status:         data.TahunPelajaran.Status,
			CreatedAt:      data.TahunPelajaran.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      data.TahunPelajaran.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedByID:    data.TahunPelajaran.CreatedByID,
			UpdatedByID:    data.TahunPelajaran.UpdatedByID,
		}
	}

	return &dtos.PesertaDidikResponse{
		ID:               data.ID,
		Nama:             data.Nama,
		NIS:              data.NIS,
		JenisKelamin:     data.JenisKelamin,
		NISN:             data.NISN,
		TempatLahir:      data.TempatLahir,
		TanggalLahir:     tanggalLahir,
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
		CreatedAt:        data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        data.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedByID:      data.CreatedByID,
		UpdatedByID:      data.UpdatedByID,
	}
}