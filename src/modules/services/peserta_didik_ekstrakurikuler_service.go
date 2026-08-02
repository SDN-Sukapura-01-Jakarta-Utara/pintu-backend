package services

import (
	"errors"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// PesertaDidikEkstrakurikulerService handles business logic
type PesertaDidikEkstrakurikulerService interface {
	RegisterOrUpdateEkstrakurikuler(req *dtos.RegisterOrUpdateEkstrakurikulerRequest, userID uint) (*dtos.RegisterOrUpdateEkstrakurikulerResponse, error)
	GetEkstrakurikulerByPesertaDidik(req *dtos.GetEkstrakurikulerPesertaDidikRequest) (*dtos.GetEkstrakurikulerPesertaDidikResponse, error)
	GetAllEkstrakurikulerSiswa(req *dtos.GetAllEkstrakurikulerSiswaRequest) (*dtos.GetAllEkstrakurikulerSiswaResponse, error)
	RegisterAllEkstrakurikulerSiswa(req *dtos.RegisterAllEkstrakurikulerSiswaRequest, userID uint) (*dtos.RegisterAllEkstrakurikulerSiswaResponse, error)
	GetStatistikEkstrakurikuler(req *dtos.GetStatistikEkstrakurikulerRequest) (*dtos.GetStatistikEkstrakurikulerResponse, error)
}

type PesertaDidikEkstrakurikulerServiceImpl struct {
	repository          repositories.PesertaDidikEkstrakurikulerRepository
	rombelRepo          repositories.PesertaDidikRombelRepository
	ekstrakurikulerRepo repositories.EkstrakurikulerRepository
}

// NewPesertaDidikEkstrakurikulerService creates a new service
func NewPesertaDidikEkstrakurikulerService(
	repository repositories.PesertaDidikEkstrakurikulerRepository,
	rombelRepo repositories.PesertaDidikRombelRepository,
	ekstrakurikulerRepo repositories.EkstrakurikulerRepository,
) PesertaDidikEkstrakurikulerService {
	return &PesertaDidikEkstrakurikulerServiceImpl{
		repository:          repository,
		rombelRepo:          rombelRepo,
		ekstrakurikulerRepo: ekstrakurikulerRepo,
	}
}

func (s *PesertaDidikEkstrakurikulerServiceImpl) RegisterOrUpdateEkstrakurikuler(req *dtos.RegisterOrUpdateEkstrakurikulerRequest, userID uint) (*dtos.RegisterOrUpdateEkstrakurikulerResponse, error) {
	// Validate peserta_didik_rombel exists
	rombel, err := s.rombelRepo.GetByID(req.PesertaDidikRombelID)
	if err != nil {
		return nil, errors.New("peserta didik rombel tidak ditemukan")
	}

	// Validate peserta_didik_rombel is active
	if rombel.Status != "active" {
		return nil, errors.New("peserta didik rombel tidak aktif")
	}

	// Get existing registrations
	existingRegistrations, err := s.repository.GetAllByPesertaDidikRombel(req.PesertaDidikRombelID)
	if err != nil {
		return nil, errors.New("gagal mengambil data ekstrakurikuler existing")
	}

	// Create map of existing ekstrakurikuler IDs
	existingMap := make(map[uint]uint) // ekstrakurikuler_id -> registration_id
	for _, reg := range existingRegistrations {
		existingMap[reg.EkstrakurikulerID] = reg.ID
	}

	// Create map of new ekstrakurikuler IDs
	newMap := make(map[uint]bool)
	for _, ekskulID := range req.EkstrakurikulerIDs {
		newMap[ekskulID] = true
	}

	// Determine actions: add, remove, keep
	var toAdd []uint
	var toRemove []uint
	kept := 0

	// Find ekstrakurikuler to add
	for _, ekskulID := range req.EkstrakurikulerIDs {
		if _, exists := existingMap[ekskulID]; !exists {
			toAdd = append(toAdd, ekskulID)
		} else {
			kept++
		}
	}

	// Find ekstrakurikuler to remove
	for ekskulID, regID := range existingMap {
		if !newMap[ekskulID] {
			toRemove = append(toRemove, regID)
		}
	}

	// Validate all new ekstrakurikuler
	for _, ekskulID := range req.EkstrakurikulerIDs {
		ekskul, err := s.ekstrakurikulerRepo.GetByID(ekskulID)
		if err != nil {
			return nil, errors.New("ekstrakurikuler dengan id tidak ditemukan")
		}

		if ekskul.Status != "active" {
			return nil, errors.New("ekstrakurikuler " + ekskul.Name + " tidak aktif")
		}
	}

	// Remove old registrations
	for _, regID := range toRemove {
		if err := s.repository.Delete(regID); err != nil {
			return nil, errors.New("gagal menghapus pendaftaran lama")
		}
	}

	// Add new registrations
	var newRegistrations []models.PesertaDidikEkstrakurikuler
	for _, ekskulID := range toAdd {
		registration := models.PesertaDidikEkstrakurikuler{
			PesertaDidikRombelID: req.PesertaDidikRombelID,
			EkstrakurikulerID:    ekskulID,
			CreatedByID:          &userID,
		}
		newRegistrations = append(newRegistrations, registration)
	}

	if len(newRegistrations) > 0 {
		if err := s.repository.CreateBulk(newRegistrations); err != nil {
			return nil, errors.New("gagal menyimpan pendaftaran ekstrakurikuler: " + err.Error())
		}
	}

	// Get final registrations
	finalRegistrations, err := s.repository.GetAllByPesertaDidikRombel(req.PesertaDidikRombelID)
	if err != nil {
		return nil, errors.New("gagal mengambil data ekstrakurikuler final")
	}

	// Map to response
	var responses []dtos.PesertaDidikEkstrakurikulerResponse
	for _, reg := range finalRegistrations {
		ekskul, _ := s.ekstrakurikulerRepo.GetByID(reg.EkstrakurikulerID)

		responses = append(responses, dtos.PesertaDidikEkstrakurikulerResponse{
			ID:                   reg.ID,
			PesertaDidikRombelID: reg.PesertaDidikRombelID,
			EkstrakurikulerID:    reg.EkstrakurikulerID,
			Ekstrakurikuler: &dtos.EkstrakurikulerResponse{
				ID:       ekskul.ID,
				Name:     ekskul.Name,
				Kategori: ekskul.Kategori,
				Status:   ekskul.Status,
			},
			CreatedAt: reg.CreatedAt,
			UpdatedAt: reg.UpdatedAt,
		})
	}

	message := "Berhasil menyimpan pendaftaran ekstrakurikuler"
	if len(existingRegistrations) > 0 {
		message = "Berhasil memperbarui pendaftaran ekstrakurikuler"
	}

	response := &dtos.RegisterOrUpdateEkstrakurikulerResponse{
		Message:       message,
		Registrations: responses,
	}
	response.Summary.Added = len(toAdd)
	response.Summary.Removed = len(toRemove)
	response.Summary.Kept = kept

	return response, nil
}

func (s *PesertaDidikEkstrakurikulerServiceImpl) GetEkstrakurikulerByPesertaDidik(req *dtos.GetEkstrakurikulerPesertaDidikRequest) (*dtos.GetEkstrakurikulerPesertaDidikResponse, error) {
	// Validate peserta_didik_rombel exists
	_, err := s.rombelRepo.GetByID(req.PesertaDidikRombelID)
	if err != nil {
		return nil, errors.New("peserta didik rombel tidak ditemukan")
	}

	// Get all registrations
	registrations, err := s.repository.GetAllByPesertaDidikRombel(req.PesertaDidikRombelID)
	if err != nil {
		return nil, errors.New("gagal mengambil data ekstrakurikuler")
	}

	// Map to response
	var responses []dtos.PesertaDidikEkstrakurikulerResponse
	for _, reg := range registrations {
		ekskul, err := s.ekstrakurikulerRepo.GetByID(reg.EkstrakurikulerID)
		if err != nil {
			continue // Skip if ekstrakurikuler not found
		}

		responses = append(responses, dtos.PesertaDidikEkstrakurikulerResponse{
			ID:                   reg.ID,
			PesertaDidikRombelID: reg.PesertaDidikRombelID,
			EkstrakurikulerID:    reg.EkstrakurikulerID,
			Ekstrakurikuler: &dtos.EkstrakurikulerResponse{
				ID:       ekskul.ID,
				Name:     ekskul.Name,
				Kategori: ekskul.Kategori,
				Status:   ekskul.Status,
			},
			CreatedAt: reg.CreatedAt,
			UpdatedAt: reg.UpdatedAt,
		})
	}

	return &dtos.GetEkstrakurikulerPesertaDidikResponse{
		PesertaDidikRombelID: req.PesertaDidikRombelID,
		Ekstrakurikuler:      responses,
		TotalEkskul:          len(responses),
	}, nil
}

func (s *PesertaDidikEkstrakurikulerServiceImpl) GetAllEkstrakurikulerSiswa(req *dtos.GetAllEkstrakurikulerSiswaRequest) (*dtos.GetAllEkstrakurikulerSiswaResponse, error) {
	// Get all ekstrakurikuler registrations for rombel and tahun pelajaran
	registrations, err := s.repository.GetAllByRombelAndTahunPelajaran(req.RombelID, req.TahunPelajaranID)
	if err != nil {
		return nil, errors.New("gagal mengambil data ekstrakurikuler siswa")
	}

	// Group by peserta_didik_rombel_id
	siswaMap := make(map[uint]*dtos.SiswaEkstrakurikuler)

	for _, reg := range registrations {
		// Skip if peserta_didik_rombel or peserta_didik is nil
		if reg.PesertaDidikRombel == nil || reg.PesertaDidikRombel.PesertaDidik == nil {
			continue
		}

		pesertaDidikRombelID := reg.PesertaDidikRombelID
		
		// Initialize siswa if not exists
		if _, exists := siswaMap[pesertaDidikRombelID]; !exists {
			siswaMap[pesertaDidikRombelID] = &dtos.SiswaEkstrakurikuler{
				PesertaDidikRombelID: pesertaDidikRombelID,
				PesertaDidikID:       reg.PesertaDidikRombel.PesertaDidikID,
				NamaLengkap:          reg.PesertaDidikRombel.PesertaDidik.NamaLengkap,
				NISN:                 reg.PesertaDidikRombel.PesertaDidik.NISN,
				Ekstrakurikuler:      []dtos.PesertaDidikEkstrakurikulerResponse{},
			}
		}

		// Add ekstrakurikuler to siswa
		if reg.Ekstrakurikuler != nil {
			siswaMap[pesertaDidikRombelID].Ekstrakurikuler = append(
				siswaMap[pesertaDidikRombelID].Ekstrakurikuler,
				dtos.PesertaDidikEkstrakurikulerResponse{
					ID:                   reg.ID,
					PesertaDidikRombelID: reg.PesertaDidikRombelID,
					EkstrakurikulerID:    reg.EkstrakurikulerID,
					Ekstrakurikuler: &dtos.EkstrakurikulerResponse{
						ID:       reg.Ekstrakurikuler.ID,
						Name:     reg.Ekstrakurikuler.Name,
						Kategori: reg.Ekstrakurikuler.Kategori,
						Status:   reg.Ekstrakurikuler.Status,
					},
					CreatedAt: reg.CreatedAt,
					UpdatedAt: reg.UpdatedAt,
				},
			)
		}
	}

	// Convert map to slice and count ekstrakurikuler
	var siswaList []dtos.SiswaEkstrakurikuler
	for _, siswa := range siswaMap {
		siswa.TotalEkskul = len(siswa.Ekstrakurikuler)
		siswaList = append(siswaList, *siswa)
	}

	return &dtos.GetAllEkstrakurikulerSiswaResponse{
		RombelID:         req.RombelID,
		TahunPelajaranID: req.TahunPelajaranID,
		Siswa:            siswaList,
		TotalSiswa:       len(siswaList),
	}, nil
}


func (s *PesertaDidikEkstrakurikulerServiceImpl) RegisterAllEkstrakurikulerSiswa(req *dtos.RegisterAllEkstrakurikulerSiswaRequest, userID uint) (*dtos.RegisterAllEkstrakurikulerSiswaResponse, error) {
	response := &dtos.RegisterAllEkstrakurikulerSiswaResponse{
		Message: "Proses bulk register/update ekstrakurikuler selesai",
		Details: []struct {
			PesertaDidikRombelID uint   `json:"peserta_didik_rombel_id"`
			Status               string `json:"status"`
			Added                int    `json:"added"`
			Removed              int    `json:"removed"`
			Kept                 int    `json:"kept"`
			Error                string `json:"error,omitempty"`
		}{},
	}

	response.Summary.TotalSiswa = len(req.Siswa)

	// Process each student
	for _, siswaInput := range req.Siswa {
		detail := struct {
			PesertaDidikRombelID uint   `json:"peserta_didik_rombel_id"`
			Status               string `json:"status"`
			Added                int    `json:"added"`
			Removed              int    `json:"removed"`
			Kept                 int    `json:"kept"`
			Error                string `json:"error,omitempty"`
		}{
			PesertaDidikRombelID: siswaInput.PesertaDidikRombelID,
		}

		// Validate peserta_didik_rombel exists
		rombel, err := s.rombelRepo.GetByID(siswaInput.PesertaDidikRombelID)
		if err != nil {
			detail.Status = "failed"
			detail.Error = "peserta didik rombel tidak ditemukan"
			response.Details = append(response.Details, detail)
			response.Summary.FailedCount++
			continue
		}

		// Validate peserta_didik_rombel is active
		if rombel.Status != "active" {
			detail.Status = "failed"
			detail.Error = "peserta didik rombel tidak aktif"
			response.Details = append(response.Details, detail)
			response.Summary.FailedCount++
			continue
		}

		// Get existing registrations
		existingRegistrations, err := s.repository.GetAllByPesertaDidikRombel(siswaInput.PesertaDidikRombelID)
		if err != nil {
			detail.Status = "failed"
			detail.Error = "gagal mengambil data ekstrakurikuler existing"
			response.Details = append(response.Details, detail)
			response.Summary.FailedCount++
			continue
		}

		// Create map of existing ekstrakurikuler IDs
		existingMap := make(map[uint]uint) // ekstrakurikuler_id -> registration_id
		for _, reg := range existingRegistrations {
			existingMap[reg.EkstrakurikulerID] = reg.ID
		}

		// Create map of new ekstrakurikuler IDs
		newMap := make(map[uint]bool)
		for _, ekskulID := range siswaInput.EkstrakurikulerIDs {
			newMap[ekskulID] = true
		}

		// Determine actions: add, remove, keep
		var toAdd []uint
		var toRemove []uint
		kept := 0

		// Find ekstrakurikuler to add
		for _, ekskulID := range siswaInput.EkstrakurikulerIDs {
			if _, exists := existingMap[ekskulID]; !exists {
				toAdd = append(toAdd, ekskulID)
			} else {
				kept++
			}
		}

		// Find ekstrakurikuler to remove
		for ekskulID, regID := range existingMap {
			if !newMap[ekskulID] {
				toRemove = append(toRemove, regID)
			}
		}

		// Validate all new ekstrakurikuler
		validationFailed := false
		for _, ekskulID := range siswaInput.EkstrakurikulerIDs {
			ekskul, err := s.ekstrakurikulerRepo.GetByID(ekskulID)
			if err != nil {
				detail.Status = "failed"
				detail.Error = "ekstrakurikuler dengan id tidak ditemukan"
				validationFailed = true
				break
			}

			if ekskul.Status != "active" {
				detail.Status = "failed"
				detail.Error = "ekstrakurikuler " + ekskul.Name + " tidak aktif"
				validationFailed = true
				break
			}
		}

		if validationFailed {
			response.Details = append(response.Details, detail)
			response.Summary.FailedCount++
			continue
		}

		// Remove old registrations
		for _, regID := range toRemove {
			if err := s.repository.Delete(regID); err != nil {
				detail.Status = "failed"
				detail.Error = "gagal menghapus pendaftaran lama"
				response.Details = append(response.Details, detail)
				response.Summary.FailedCount++
				continue
			}
		}

		// Add new registrations
		var newRegistrations []models.PesertaDidikEkstrakurikuler
		for _, ekskulID := range toAdd {
			registration := models.PesertaDidikEkstrakurikuler{
				PesertaDidikRombelID: siswaInput.PesertaDidikRombelID,
				EkstrakurikulerID:    ekskulID,
				CreatedByID:          &userID,
			}
			newRegistrations = append(newRegistrations, registration)
		}

		if len(newRegistrations) > 0 {
			if err := s.repository.CreateBulk(newRegistrations); err != nil {
				detail.Status = "failed"
				detail.Error = "gagal menyimpan pendaftaran ekstrakurikuler"
				response.Details = append(response.Details, detail)
				response.Summary.FailedCount++
				continue
			}
		}

		// Success
		detail.Status = "success"
		detail.Added = len(toAdd)
		detail.Removed = len(toRemove)
		detail.Kept = kept
		response.Details = append(response.Details, detail)
		response.Summary.SuccessCount++
		response.Summary.TotalAdded += len(toAdd)
		response.Summary.TotalRemoved += len(toRemove)
		response.Summary.TotalKept += kept
	}

	return response, nil
}


func (s *PesertaDidikEkstrakurikulerServiceImpl) GetStatistikEkstrakurikuler(req *dtos.GetStatistikEkstrakurikulerRequest) (*dtos.GetStatistikEkstrakurikulerResponse, error) {
	response := &dtos.GetStatistikEkstrakurikulerResponse{
		TahunPelajaranID: req.TahunPelajaranID,
		RombelID:         req.RombelID,
	}

	// Get all ekstrakurikuler registrations
	registrations, err := s.repository.GetStatistikByTahunPelajaran(req.TahunPelajaranID, req.RombelID)
	if err != nil {
		return nil, errors.New("gagal mengambil data ekstrakurikuler")
	}

	// Get all students (to find those not joining any ekstrakurikuler)
	var allStudents []models.PesertaDidikRombel
	filterParams := repositories.GetPesertaDidikRombelParams{
		Filter: repositories.GetPesertaDidikRombelFilter{
			TahunPelajaranID: req.TahunPelajaranID,
			Status:           "active",
		},
		Limit:  10000, // Large number to get all
		Offset: 0,
	}
	
	if req.RombelID != nil {
		filterParams.Filter.RombelID = *req.RombelID
	}
	
	allStudents, _, err = s.rombelRepo.GetAllWithFilter(filterParams)
	if err != nil {
		return nil, errors.New("gagal mengambil data siswa")
	}

	// Set nama rombel if specific rombel requested
	if req.RombelID != nil && len(allStudents) > 0 {
		if allStudents[0].Rombel != nil {
			response.NamaRombel = allStudents[0].Rombel.Name
		}
	}

	// Create maps for processing
	studentsWithEkskul := make(map[uint]bool) // peserta_didik_rombel_id
	ekskulMap := make(map[uint]map[uint]int) // ekstrakurikuler_id -> rombel_id -> count
	rombelMap := make(map[uint]map[uint]int) // rombel_id -> ekstrakurikuler_id -> count
	ekskulNames := make(map[uint]string)
	ekskulKategori := make(map[uint]string)
	rombelNames := make(map[uint]string)
	ekskulSet := make(map[uint]bool)
	rombelSet := make(map[uint]bool)

	// Process registrations
	for _, reg := range registrations {
		if reg.PesertaDidikRombel == nil || reg.Ekstrakurikuler == nil {
			continue
		}

		studentsWithEkskul[reg.PesertaDidikRombelID] = true
		
		ekskulID := reg.EkstrakurikulerID
		rombelID := reg.PesertaDidikRombel.RombelID
		
		// Track ekstrakurikuler
		if _, exists := ekskulMap[ekskulID]; !exists {
			ekskulMap[ekskulID] = make(map[uint]int)
		}
		ekskulMap[ekskulID][rombelID]++
		ekskulNames[ekskulID] = reg.Ekstrakurikuler.Name
		ekskulKategori[ekskulID] = reg.Ekstrakurikuler.Kategori
		ekskulSet[ekskulID] = true
		
		// Track rombel
		if _, exists := rombelMap[rombelID]; !exists {
			rombelMap[rombelID] = make(map[uint]int)
		}
		rombelMap[rombelID][ekskulID]++
		
		if reg.PesertaDidikRombel.Rombel != nil {
			rombelNames[rombelID] = reg.PesertaDidikRombel.Rombel.Name
		}
		rombelSet[rombelID] = true
	}

	// Build statistik per ekstrakurikuler
	for ekskulID := range ekskulSet {
		stat := dtos.StatistikPerEkskul{
			EkstrakurikulerID:   ekskulID,
			NamaEkstrakurikuler: ekskulNames[ekskulID],
			Kategori:            ekskulKategori[ekskulID],
			TotalSiswa:          0,
			Rombel:              []struct {
				RombelID   uint   `json:"rombel_id"`
				NamaRombel string `json:"nama_rombel"`
				JumlahSiswa int   `json:"jumlah_siswa"`
			}{},
		}
		
		for rombelID, count := range ekskulMap[ekskulID] {
			stat.TotalSiswa += count
			stat.Rombel = append(stat.Rombel, struct {
				RombelID   uint   `json:"rombel_id"`
				NamaRombel string `json:"nama_rombel"`
				JumlahSiswa int   `json:"jumlah_siswa"`
			}{
				RombelID:   rombelID,
				NamaRombel: rombelNames[rombelID],
				JumlahSiswa: count,
			})
		}
		
		response.StatistikPerEkskul = append(response.StatistikPerEkskul, stat)
	}

	// Build statistik per rombel
	rombelTotalStudents := make(map[uint]int)
	for _, student := range allStudents {
		rombelTotalStudents[student.RombelID]++
	}
	
	for rombelID := range rombelSet {
		totalSiswa := rombelTotalStudents[rombelID]
		
		// Count students in this rombel who joined ekstrakurikuler
		siswaIkutEkskul := 0
		for _, student := range allStudents {
			if student.RombelID == rombelID && studentsWithEkskul[student.ID] {
				siswaIkutEkskul++
			}
		}
		
		persentase := 0.0
		if totalSiswa > 0 {
			persentase = (float64(siswaIkutEkskul) / float64(totalSiswa)) * 100
		}
		
		stat := dtos.StatistikPerRombel{
			RombelID:              rombelID,
			NamaRombel:            rombelNames[rombelID],
			TotalSiswa:            totalSiswa,
			SiswaIkutEkskul:       siswaIkutEkskul,
			SiswaTidakIkutEkskul:  totalSiswa - siswaIkutEkskul,
			PersentaseIkutEkskul:  persentase,
			Ekstrakurikuler:       []struct {
				EkstrakurikulerID   uint   `json:"ekstrakurikuler_id"`
				NamaEkstrakurikuler string `json:"nama_ekstrakurikuler"`
				JumlahSiswa         int    `json:"jumlah_siswa"`
			}{},
		}
		
		for ekskulID, count := range rombelMap[rombelID] {
			stat.Ekstrakurikuler = append(stat.Ekstrakurikuler, struct {
				EkstrakurikulerID   uint   `json:"ekstrakurikuler_id"`
				NamaEkstrakurikuler string `json:"nama_ekstrakurikuler"`
				JumlahSiswa         int    `json:"jumlah_siswa"`
			}{
				EkstrakurikulerID:   ekskulID,
				NamaEkstrakurikuler: ekskulNames[ekskulID],
				JumlahSiswa:         count,
			})
		}
		
		response.StatistikPerRombel = append(response.StatistikPerRombel, stat)
	}

	// Find students not joining any ekstrakurikuler
	for _, student := range allStudents {
		if !studentsWithEkskul[student.ID] {
			if student.PesertaDidik == nil || student.Rombel == nil {
				continue
			}
			
			response.SiswaTidakIkutEkskul = append(response.SiswaTidakIkutEkskul, dtos.SiswaTidakIkutEkskul{
				PesertaDidikRombelID: student.ID,
				PesertaDidikID:       student.PesertaDidikID,
				NamaLengkap:          student.PesertaDidik.NamaLengkap,
				NISN:                 student.PesertaDidik.NISN,
				RombelID:             student.RombelID,
				NamaRombel:           student.Rombel.Name,
			})
		}
	}

	// Calculate overall summary
	response.Summary.TotalSiswa = len(allStudents)
	response.Summary.TotalSiswaIkutEkskul = len(studentsWithEkskul)
	response.Summary.TotalSiswaTidakIkutEkskul = len(response.SiswaTidakIkutEkskul)
	response.Summary.TotalEkstrakurikuler = len(ekskulSet)
	response.Summary.TotalRombel = len(rombelSet)
	
	if response.Summary.TotalSiswa > 0 {
		response.Summary.PersentaseIkutEkskul = (float64(response.Summary.TotalSiswaIkutEkskul) / float64(response.Summary.TotalSiswa)) * 100
	}

	return response, nil
}
