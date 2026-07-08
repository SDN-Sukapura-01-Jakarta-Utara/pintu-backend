package services

import (
	"fmt"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"time"
)

// AbsensiScanService handles business logic for Absensi Scan
type AbsensiScanService interface {
	ScanAbsensi(req *dtos.AbsensiScanRequest) (*dtos.AbsensiScanResponse, error)
}

type AbsensiScanServiceImpl struct {
	repository repositories.AbsensiScanRepository
}

// NewAbsensiScanService creates a new Absensi Scan service
func NewAbsensiScanService(repository repositories.AbsensiScanRepository) AbsensiScanService {
	return &AbsensiScanServiceImpl{
		repository: repository,
	}
}

// ScanAbsensi processes attendance scanning
func (s *AbsensiScanServiceImpl) ScanAbsensi(req *dtos.AbsensiScanRequest) (*dtos.AbsensiScanResponse, error) {
	// 1. Get konfigurasi (direct from DB, no cache)
	config, err := s.repository.GetKonfigurasiAbsensi()
	if err != nil {
		return &dtos.AbsensiScanResponse{
			Success: false,
			Message: "Konfigurasi absensi belum diatur",
		}, nil
	}

	// 2. Validate barcode and get peserta didik
	pesertaDidik, err := s.repository.GetPesertaDidikByBarcode(req.Barcode)
	if err != nil {
		return &dtos.AbsensiScanResponse{
			Success: false,
			Message: "Barcode tidak ditemukan",
		}, nil
	}

	// 3. Validate status peserta didik (must be active)
	if pesertaDidik.Status != "active" {
		return &dtos.AbsensiScanResponse{
			Success: false,
			Message: "Siswa tidak aktif, tidak dapat melakukan absensi",
		}, nil
	}

	// 4. Get current time (with Asia/Jakarta timezone)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback: use FixedZone if LoadLocation fails (for containers without tzdata)
		loc = time.FixedZone("WIB", 7*60*60) // UTC+7
	}
	now := time.Now().In(loc)
	currentDate := now.Format("2006-01-02")
	currentTime := now.Format("15:04:05")

	// 5. Validate time range
	scanType, status, validationErr := s.validateScanTime(currentTime, config)
	if validationErr != nil {
		return &dtos.AbsensiScanResponse{
			Success: false,
			Message: validationErr.Error(),
		}, nil
	}

	// 6. Check existing absensi
	existingAbsensi, err := s.repository.GetAbsensiByPesertaDidikAndDate(pesertaDidik.ID, now)
	
	if err == nil {
		// Record exists, check if already scanned for this type
		if scanType == "datang" && existingAbsensi.JamDatang != nil {
			// Already scanned for datang
			statusValue := "unknown"
			if existingAbsensi.Status != nil {
				statusValue = *existingAbsensi.Status
			}
			
			return &dtos.AbsensiScanResponse{
				Success: true,
				Message: "Anda sudah melakukan absen datang hari ini",
				PesertaDidik: &dtos.PesertaDidikInfo{
					ID:   pesertaDidik.ID,
					Nama: pesertaDidik.Nama,
					NISN: pesertaDidik.NISN,
				},
				AbsensiInfo: &dtos.AbsensiInfo{
					Tanggal:   currentDate,
					JamDatang: existingAbsensi.JamDatang,
					JamPulang: existingAbsensi.JamPulang,
					Status:    statusValue,
					IsUpdate:  false,
				},
			}, nil
		}
		
		if scanType == "pulang" {
			// Check if already scanned datang first
			if existingAbsensi.JamDatang == nil {
				// Belum scan datang, tidak boleh scan pulang
				return &dtos.AbsensiScanResponse{
					Success: false,
					Message: "Anda belum melakukan absen datang, tidak dapat melakukan absen pulang",
				}, nil
			}
			
			// Already scanned for pulang
			if existingAbsensi.JamPulang != nil {
				statusValue := "unknown"
				if existingAbsensi.Status != nil {
					statusValue = *existingAbsensi.Status
				}
				
				return &dtos.AbsensiScanResponse{
					Success: true,
					Message: "Anda sudah melakukan absen pulang hari ini",
					PesertaDidik: &dtos.PesertaDidikInfo{
						ID:   pesertaDidik.ID,
						Nama: pesertaDidik.Nama,
						NISN: pesertaDidik.NISN,
					},
					AbsensiInfo: &dtos.AbsensiInfo{
						Tanggal:   currentDate,
						JamDatang: existingAbsensi.JamDatang,
						JamPulang: existingAbsensi.JamPulang,
						Status:    statusValue,
						IsUpdate:  false,
					},
				}, nil
			}
		}
	} else {
		// No existing record
		if scanType == "pulang" {
			// Trying to scan pulang without datang first
			return &dtos.AbsensiScanResponse{
				Success: false,
				Message: "Anda belum melakukan absen datang, tidak dapat melakukan absen pulang",
			}, nil
		}
	}
	
	// 7. Prepare absensi record
	var absensi *models.Absensi
	isUpdate := false

	if err != nil {
		// No existing record, create new
		absensi = &models.Absensi{
			PesertaDidikID: pesertaDidik.ID,
			Tanggal:        now,
		}
	} else {
		// Record exists but hasn't scanned for this type yet, update
		absensi = existingAbsensi
		isUpdate = true
	}

	// 8. Update fields based on scan type
	if scanType == "datang" {
		absensi.JamDatang = &currentTime
		absensi.Status = &status
	} else if scanType == "pulang" {
		absensi.JamPulang = &currentTime
		// Status tidak berubah saat pulang
	}

	// 9. Save to database using UPSERT
	if err := s.repository.UpsertAbsensi(absensi); err != nil {
		return &dtos.AbsensiScanResponse{
			Success: false,
			Message: "Gagal menyimpan data absensi",
		}, err
	}

	// 10. Build response
	return &dtos.AbsensiScanResponse{
		Success: true,
		Message: s.buildSuccessMessage(scanType, status, isUpdate),
		PesertaDidik: &dtos.PesertaDidikInfo{
			ID:   pesertaDidik.ID,
			Nama: pesertaDidik.Nama,
			NISN: pesertaDidik.NISN,
		},
		AbsensiInfo: &dtos.AbsensiInfo{
			Tanggal:   currentDate,
			JamDatang: absensi.JamDatang,
			JamPulang: absensi.JamPulang,
			Status:    *absensi.Status,
			IsUpdate:  isUpdate,
		},
	}, nil
}

// validateScanTime validates if current time is within allowed range
func (s *AbsensiScanServiceImpl) validateScanTime(currentTime string, config *models.KonfigurasiAbsensi) (scanType string, status string, err error) {
	// Check if within "datang" time range
	if currentTime >= config.JamDatangMulai && currentTime <= config.JamDatangSelesai {
		scanType = "datang"
		
		// Check if tepat waktu or terlambat
		if currentTime <= config.JamMaxDatang {
			status = "tepat_waktu"
		} else {
			status = "terlambat"
		}
		return scanType, status, nil
	}

	// Check if within "pulang" time range
	if currentTime >= config.JamPulangMulai && currentTime <= config.JamPulangSelesai {
		scanType = "pulang"
		status = "" // Status tidak berubah saat pulang
		return scanType, status, nil
	}

	// Outside allowed time range - provide helpful error message
	return "", "", fmt.Errorf("Scan absensi hanya dapat dilakukan pada jam yang ditentukan. Jam sekarang: %s. Rentang datang: %s-%s. Rentang pulang: %s-%s", 
		currentTime, config.JamDatangMulai, config.JamDatangSelesai, config.JamPulangMulai, config.JamPulangSelesai)
}

// buildSuccessMessage creates success message based on scan type and status
func (s *AbsensiScanServiceImpl) buildSuccessMessage(scanType, status string, isUpdate bool) string {
	action := "Absensi"
	if isUpdate {
		action = "Update absensi"
	}

	if scanType == "datang" {
		if status == "tepat_waktu" {
			return fmt.Sprintf("%s berhasil - Datang tepat waktu", action)
		}
		return fmt.Sprintf("%s berhasil - Datang terlambat", action)
	}
	
	return fmt.Sprintf("%s berhasil - Pulang tercatat", action)
}
