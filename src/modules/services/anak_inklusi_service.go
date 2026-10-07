package services

import (
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
)

type AnakInklusiService interface {
	Create(fileSuratDokter *multipart.FileHeader, req *dtos.AnakInklusiCreateRequest, userID uint) (*dtos.AnakInklusiResponse, error)
	GetAllWithFilter(params repositories.GetAnakInklusiParams) (*dtos.AnakInklusiListWithPaginationResponse, error)
	GetHistoryWithFilter(params repositories.GetAnakInklusiParams) (*dtos.AnakInklusiListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.AnakInklusiResponse, error)
	Update(id uint, fileSuratDokter *multipart.FileHeader, req *dtos.AnakInklusiUpdateRequest, userID uint) (*dtos.AnakInklusiResponse, error)
	Delete(id uint) error
}

type anakInklusiService struct {
	repo      repositories.AnakInklusiRepository
	r2Storage *utils.R2Storage
}

func NewAnakInklusiService(repo repositories.AnakInklusiRepository, r2Storage *utils.R2Storage) AnakInklusiService {
	return &anakInklusiService{
		repo:      repo,
		r2Storage: r2Storage,
	}
}

func (s *anakInklusiService) Create(fileSuratDokter *multipart.FileHeader, req *dtos.AnakInklusiCreateRequest, userID uint) (*dtos.AnakInklusiResponse, error) {
	// Check if peserta didik already exists in anak_inklusi
	existing, _ := s.repo.GetByPesertaDidikID(req.PesertaDidikID)
	if existing != nil {
		return nil, errors.New("peserta didik sudah terdaftar sebagai anak inklusi")
	}

	// Validate file if provided
	var filePath string
	if fileSuratDokter != nil {
		// Validate file type (pdf only)
		ext := strings.ToLower(filepath.Ext(fileSuratDokter.Filename))
		if ext != ".pdf" {
			return nil, errors.New("file surat dokter harus berformat PDF")
		}

		// Validate file size (max 5MB)
		if fileSuratDokter.Size > 5*1024*1024 {
			return nil, errors.New("ukuran file surat dokter maksimal 5MB")
		}

		// Upload to R2
		uploadedPath, err := s.r2Storage.UploadFile(fileSuratDokter, "inklusi/surat-dokter")
		if err != nil {
			return nil, fmt.Errorf("gagal mengunggah file: %w", err)
		}
		filePath = uploadedPath
	}

	// Parse dates
	var tanggalIdentifikasi, tanggalDiagnosa, tanggalKadaluarsaSurat *time.Time

	if req.TanggalIdentifikasi != "" {
		t, err := time.Parse("2006-01-02", req.TanggalIdentifikasi)
		if err != nil {
			// Clean up uploaded file if date parsing fails
			if filePath != "" {
				s.r2Storage.DeleteFile(filePath)
			}
			return nil, errors.New("format tanggal identifikasi tidak valid (gunakan YYYY-MM-DD)")
		}
		tanggalIdentifikasi = &t
	}

	if req.TanggalDiagnosa != "" {
		t, err := time.Parse("2006-01-02", req.TanggalDiagnosa)
		if err != nil {
			// Clean up uploaded file if date parsing fails
			if filePath != "" {
				s.r2Storage.DeleteFile(filePath)
			}
			return nil, errors.New("format tanggal diagnosa tidak valid (gunakan YYYY-MM-DD)")
		}
		tanggalDiagnosa = &t
	}

	if req.TanggalKadaluarsaSurat != "" {
		t, err := time.Parse("2006-01-02", req.TanggalKadaluarsaSurat)
		if err != nil {
			// Clean up uploaded file if date parsing fails
			if filePath != "" {
				s.r2Storage.DeleteFile(filePath)
			}
			return nil, errors.New("format tanggal kadaluarsa surat tidak valid (gunakan YYYY-MM-DD)")
		}
		tanggalKadaluarsaSurat = &t
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "identified"
	}

	// Create model
	anakInklusi := &models.AnakInklusi{
		PesertaDidikID:         req.PesertaDidikID,
		JenisHambatan:          req.JenisHambatan,
		TanggalIdentifikasi:    tanggalIdentifikasi,
		TanggalDiagnosa:        tanggalDiagnosa,
		TanggalKadaluarsaSurat: tanggalKadaluarsaSurat,
		FileSuratDokter:        filePath,
		Status:                 status,
		Catatan:                req.Catatan,
		CreatedByID:            &userID,
		UpdatedByID:            &userID,
	}

	// Save to database
	if err := s.repo.Create(anakInklusi); err != nil {
		// Clean up uploaded file if database save fails
		if filePath != "" {
			s.r2Storage.DeleteFile(filePath)
		}
		return nil, fmt.Errorf("gagal menyimpan data: %w", err)
	}

	// Get created data with relations
	created, err := s.repo.GetByID(anakInklusi.ID)
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	return s.toResponseDTO(created), nil
}

func (s *anakInklusiService) GetAllWithFilter(params repositories.GetAnakInklusiParams) (*dtos.AnakInklusiListWithPaginationResponse, error) {
	anakInklusiList, total, err := s.repo.GetAllWithFilter(params)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	var responses []dtos.AnakInklusiResponse
	for _, anakInklusi := range anakInklusiList {
		responses = append(responses, *s.toResponseDTO(&anakInklusi))
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	page := (params.Offset / params.Limit) + 1

	return &dtos.AnakInklusiListWithPaginationResponse{
		Data: responses,
		Pagination: dtos.PaginationResponse{
			Limit:      params.Limit,
			Offset:     params.Offset,
			Page:       page,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *anakInklusiService) GetHistoryWithFilter(params repositories.GetAnakInklusiParams) (*dtos.AnakInklusiListWithPaginationResponse, error) {
	anakInklusiList, total, err := s.repo.GetAllWithFilter(params)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	var responses []dtos.AnakInklusiResponse
	for _, anakInklusi := range anakInklusiList {
		responses = append(responses, *s.toResponseDTO(&anakInklusi))
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	page := (params.Offset / params.Limit) + 1

	return &dtos.AnakInklusiListWithPaginationResponse{
		Data: responses,
		Pagination: dtos.PaginationResponse{
			Limit:      params.Limit,
			Offset:     params.Offset,
			Page:       page,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *anakInklusiService) GetByID(id uint) (*dtos.AnakInklusiResponse, error) {
	anakInklusi, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("data anak inklusi tidak ditemukan")
	}

	return s.toResponseDTO(anakInklusi), nil
}

func (s *anakInklusiService) Update(id uint, fileSuratDokter *multipart.FileHeader, req *dtos.AnakInklusiUpdateRequest, userID uint) (*dtos.AnakInklusiResponse, error) {
	// Get existing data
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("data anak inklusi tidak ditemukan")
	}

	// Track old file for deletion
	oldFilePath := existing.FileSuratDokter

	// Handle file upload if new file provided
	var newFilePath string
	if fileSuratDokter != nil {
		// Validate file type (pdf only)
		ext := strings.ToLower(filepath.Ext(fileSuratDokter.Filename))
		if ext != ".pdf" {
			return nil, errors.New("file surat dokter harus berformat PDF")
		}

		// Validate file size (max 5MB)
		if fileSuratDokter.Size > 5*1024*1024 {
			return nil, errors.New("ukuran file surat dokter maksimal 5MB")
		}

		// Upload to R2
		uploadedPath, err := s.r2Storage.UploadFile(fileSuratDokter, "inklusi/surat-dokter")
		if err != nil {
			return nil, fmt.Errorf("gagal mengunggah file: %w", err)
		}
		newFilePath = uploadedPath

		// Delete old file if exists
		if oldFilePath != "" {
			s.r2Storage.DeleteFile(oldFilePath)
		}

		existing.FileSuratDokter = newFilePath
	}

	// Handle file deletion if requested
	if req.DeleteFile && newFilePath == "" {
		if oldFilePath != "" {
			s.r2Storage.DeleteFile(oldFilePath)
		}
		existing.FileSuratDokter = ""
	}

	// Update fields only if provided (not empty)
	if req.JenisHambatan != "" {
		existing.JenisHambatan = req.JenisHambatan
	}

	if req.TanggalIdentifikasi != "" {
		t, err := time.Parse("2006-01-02", req.TanggalIdentifikasi)
		if err != nil {
			// Clean up uploaded file if date parsing fails
			if newFilePath != "" {
				s.r2Storage.DeleteFile(newFilePath)
			}
			return nil, errors.New("format tanggal identifikasi tidak valid (gunakan YYYY-MM-DD)")
		}
		existing.TanggalIdentifikasi = &t
	}

	if req.TanggalDiagnosa != "" {
		t, err := time.Parse("2006-01-02", req.TanggalDiagnosa)
		if err != nil {
			// Clean up uploaded file if date parsing fails
			if newFilePath != "" {
				s.r2Storage.DeleteFile(newFilePath)
			}
			return nil, errors.New("format tanggal diagnosa tidak valid (gunakan YYYY-MM-DD)")
		}
		existing.TanggalDiagnosa = &t
	}

	if req.TanggalKadaluarsaSurat != "" {
		t, err := time.Parse("2006-01-02", req.TanggalKadaluarsaSurat)
		if err != nil {
			// Clean up uploaded file if date parsing fails
			if newFilePath != "" {
				s.r2Storage.DeleteFile(newFilePath)
			}
			return nil, errors.New("format tanggal kadaluarsa surat tidak valid (gunakan YYYY-MM-DD)")
		}
		existing.TanggalKadaluarsaSurat = &t
	}

	if req.Status != "" {
		existing.Status = req.Status
	}

	if req.Catatan != "" {
		existing.Catatan = req.Catatan
	}

	// Update metadata
	existing.UpdatedByID = &userID

	// Save to database
	if err := s.repo.Update(existing); err != nil {
		// Clean up uploaded file if database save fails
		if newFilePath != "" {
			s.r2Storage.DeleteFile(newFilePath)
		}
		return nil, fmt.Errorf("gagal menyimpan data: %w", err)
	}

	// Get updated data with relations
	updated, err := s.repo.GetByID(existing.ID)
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	return s.toResponseDTO(updated), nil
}

func (s *anakInklusiService) Delete(id uint) error {
	// Get existing data to check if exists and get file path
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("data anak inklusi tidak ditemukan")
	}

	// Delete file from R2 if exists
	if existing.FileSuratDokter != "" {
		if err := s.r2Storage.DeleteFile(existing.FileSuratDokter); err != nil {
			// Log error but continue with deletion
			fmt.Printf("Warning: Failed to delete file from R2: %v\n", err)
		}
	}

	// Delete from database (soft delete using GORM)
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus data: %w", err)
	}

	return nil
}

func (s *anakInklusiService) toResponseDTO(anakInklusi *models.AnakInklusi) *dtos.AnakInklusiResponse {
	resp := &dtos.AnakInklusiResponse{
		ID:                  anakInklusi.ID,
		PesertaDidikID:      anakInklusi.PesertaDidikID,
		JenisHambatan:       anakInklusi.JenisHambatan,
		FileSuratDokter:     anakInklusi.FileSuratDokter,
		FileSuratDokterURL:  s.r2Storage.GetPublicURL(anakInklusi.FileSuratDokter),
		Status:              anakInklusi.Status,
		Catatan:             anakInklusi.Catatan,
		CreatedAt:           anakInklusi.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:           anakInklusi.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if anakInklusi.TanggalIdentifikasi != nil {
		formatted := anakInklusi.TanggalIdentifikasi.Format("2006-01-02")
		resp.TanggalIdentifikasi = &formatted
	}

	if anakInklusi.TanggalDiagnosa != nil {
		formatted := anakInklusi.TanggalDiagnosa.Format("2006-01-02")
		resp.TanggalDiagnosa = &formatted
	}

	if anakInklusi.TanggalKadaluarsaSurat != nil {
		formatted := anakInklusi.TanggalKadaluarsaSurat.Format("2006-01-02")
		resp.TanggalKadaluarsaSurat = &formatted
	}

	if anakInklusi.PesertaDidik != nil {
		resp.PesertaDidik = &dtos.PesertaDidikSimple{
			ID:       anakInklusi.PesertaDidik.ID,
			NIS:      anakInklusi.PesertaDidik.NIS,
			NISN:     anakInklusi.PesertaDidik.NISN,
			Nama:     anakInklusi.PesertaDidik.Nama,
			Photo:    anakInklusi.PesertaDidik.Photo,
			PhotoURL: s.r2Storage.GetPublicURL(anakInklusi.PesertaDidik.Photo),
			NamaAyah: anakInklusi.PesertaDidik.NamaAyah,
			NamaIbu:  anakInklusi.PesertaDidik.NamaIbu,
			Status:   anakInklusi.PesertaDidik.Status,
		}
	}

	return resp
}
