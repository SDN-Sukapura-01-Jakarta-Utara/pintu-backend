package services

import (
	"encoding/json"
	"errors"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
)

// StrukturOrganisasiService handles business logic for StrukturOrganisasi
type StrukturOrganisasiService interface {
	Create(req *dtos.StrukturOrganisasiCreateRequest, userID uint) (*dtos.StrukturOrganisasiResponse, error)
	GetByID(id uint) (*dtos.StrukturOrganisasiResponse, error)
	GetAll(limit int, offset int) (*dtos.StrukturOrganisasiListResponse, error)
	GetAllWithFilter(params repositories.GetStrukturOrganisasiParams) (*dtos.StrukturOrganisasiListWithPaginationResponse, error)
	Update(req *dtos.StrukturOrganisasiUpdateRequest, userID uint) (*dtos.StrukturOrganisasiResponse, error)
	Delete(id uint) error
	GetPublic() ([]dtos.StrukturOrganisasiGroupedResponse, error)
}

type StrukturOrganisasiServiceImpl struct {
	repository repositories.StrukturOrganisasiRepository
	pegawaiRepo repositories.KepegawaianRepository
}

// NewStrukturOrganisasiService creates a new StrukturOrganisasi service
func NewStrukturOrganisasiService(repository repositories.StrukturOrganisasiRepository, pegawaiRepo repositories.KepegawaianRepository) StrukturOrganisasiService {
	return &StrukturOrganisasiServiceImpl{
		repository: repository,
		pegawaiRepo: pegawaiRepo,
	}
}

// Create creates a new StrukturOrganisasi
func (s *StrukturOrganisasiServiceImpl) Create(req *dtos.StrukturOrganisasiCreateRequest, userID uint) (*dtos.StrukturOrganisasiResponse, error) {
	// Validate that either pegawai_id or nama_non_pegawai is provided
	if req.PegawaiID == nil && req.NamaNonPegawai == "" {
		return nil, errors.New("either pegawai_id or nama_non_pegawai must be provided")
	}

	// If pegawai_id is provided, validate it exists
	if req.PegawaiID != nil {
		_, err := s.pegawaiRepo.GetByID(*req.PegawaiID)
		if err != nil {
			return nil, errors.New("pegawai not found")
		}
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	data := &models.StrukturOrganisasi{
		PegawaiID:         req.PegawaiID,
		NamaNonPegawai:    req.NamaNonPegawai,
		JabatanNonPegawai: req.JabatanNonPegawai,
		Urutan:            req.Urutan,
		Relasi:            req.Relasi,
		Status:            status,
		CreatedByID:       &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves StrukturOrganisasi by ID
func (s *StrukturOrganisasiServiceImpl) GetByID(id uint) (*dtos.StrukturOrganisasiResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("struktur organisasi not found")
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all StrukturOrganisasi
func (s *StrukturOrganisasiServiceImpl) GetAll(limit int, offset int) (*dtos.StrukturOrganisasiListResponse, error) {
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
	responses := make([]dtos.StrukturOrganisasiResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.StrukturOrganisasiListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// GetAllWithFilter retrieves StrukturOrganisasi with filters and pagination
func (s *StrukturOrganisasiServiceImpl) GetAllWithFilter(params repositories.GetStrukturOrganisasiParams) (*dtos.StrukturOrganisasiListWithPaginationResponse, error) {
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
	responses := make([]dtos.StrukturOrganisasiResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.StrukturOrganisasiListWithPaginationResponse{
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

// Update updates StrukturOrganisasi
func (s *StrukturOrganisasiServiceImpl) Update(req *dtos.StrukturOrganisasiUpdateRequest, userID uint) (*dtos.StrukturOrganisasiResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("struktur organisasi not found")
	}

	// Validate if pegawai_id is being updated (explicitly set in request)
	if req.PegawaiIDSet {
		if req.PegawaiID != nil {
			_, err := s.pegawaiRepo.GetByID(*req.PegawaiID)
			if err != nil {
				return nil, errors.New("pegawai not found")
			}
		}
		existing.PegawaiID = req.PegawaiID
	}

	// Update nama_non_pegawai if provided
	if req.NamaNonPegawai != nil {
		existing.NamaNonPegawai = *req.NamaNonPegawai
	}

	// Update jabatan_non_pegawai if provided
	if req.JabatanNonPegawai != nil {
		existing.JabatanNonPegawai = *req.JabatanNonPegawai
	}

	// Update urutan if provided
	if req.Urutan != nil {
		existing.Urutan = *req.Urutan
	}

	// Update relasi if provided
	if req.Relasi != nil {
		existing.Relasi = *req.Relasi
	}

	// Update status if provided
	if req.Status != nil {
		existing.Status = *req.Status
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes StrukturOrganisasi by ID
func (s *StrukturOrganisasiServiceImpl) Delete(id uint) error {
	_, err := s.repository.GetByID(id)
	if err != nil {
		return errors.New("struktur organisasi not found")
	}
	return s.repository.Delete(id)
}

// GetPublic retrieves all active StrukturOrganisasi for public display, grouped by urutan
func (s *StrukturOrganisasiServiceImpl) GetPublic() ([]dtos.StrukturOrganisasiGroupedResponse, error) {
	data, err := s.repository.GetAllPublic()
	if err != nil {
		return nil, err
	}

	// Group by urutan for urutan 1 and 2
	groupedMap := make(map[int][]dtos.StrukturOrganisasiPublicResponse)
	// Special grouping for urutan 3 (guru kelas) by kelas
	guruKelasMap := make(map[string][]dtos.PegawaiPublicDetailResponse)
	guruKelasRelasiMap := make(map[string]string) // Store relasi for each kelas
	// Special grouping for urutan 4 (guru mapel) by bidang studi
	guruMapelMap := make(map[string][]dtos.PegawaiPublicDetailResponse)
	guruMapelRelasiMap := make(map[string]string) // Store relasi for each bidang studi
	// Grouping for urutan 5+ by jabatan
	byJabatanMap := make(map[int]map[string][]dtos.StrukturOrganisasiPublicResponse)

	for _, item := range data {
		response := dtos.StrukturOrganisasiPublicResponse{
			NamaNonPegawai:    item.NamaNonPegawai,
			JabatanNonPegawai: item.JabatanNonPegawai,
			Urutan:            item.Urutan,
			Relasi:            item.Relasi,
		}

		// Add detailed Pegawai data if available
		if item.Pegawai != nil {
			pegawaiDetail := dtos.PegawaiPublicDetailResponse{
				NamaLengkap: item.Pegawai.Nama,
				NIP:         item.Pegawai.NIP,
				NKKI:        item.Pegawai.NKKI,
				Jabatan:     item.Pegawai.Jabatan,
				Kategori:    item.Pegawai.Kategori,
			}

			// Urutan 3: Guru Kelas (group by kelas)
			if item.Urutan == 3 && item.Pegawai.Kategori == "Pendidik" {
				kelasMengajar := make([]dtos.RombelWithKelasResponse, 0)
				
				// Add rombel_guru_kelas_id (guru kelas)
				if item.Pegawai.RombelGuruKelas != nil {
					namaKelas := ""
					if item.Pegawai.RombelGuruKelas.Kelas != nil {
						namaKelas = item.Pegawai.RombelGuruKelas.Kelas.Name
					}
					kelasMengajar = append(kelasMengajar, dtos.RombelWithKelasResponse{
						ID:        item.Pegawai.RombelGuruKelas.ID,
						Rombel:    item.Pegawai.RombelGuruKelas.Name,
						NamaKelas: namaKelas,
						Status:    item.Pegawai.RombelGuruKelas.Status,
					})
					
					// Group by nama kelas
					pegawaiDetail.KelasMengajar = kelasMengajar
					guruKelasMap[namaKelas] = append(guruKelasMap[namaKelas], pegawaiDetail)
					// Store relasi (use first item's relasi if not set)
					if guruKelasRelasiMap[namaKelas] == "" {
						guruKelasRelasiMap[namaKelas] = item.Relasi
					}
				}
			} else if item.Urutan == 4 && item.Pegawai.Kategori == "Pendidik" {
				// Urutan 4: Guru Bidang Studi (group by bidang studi)
				if item.Pegawai.BidangStudi != nil {
					pegawaiDetail.BidangStudi = item.Pegawai.BidangStudi.Name
					// Group by bidang studi
					guruMapelMap[item.Pegawai.BidangStudi.Name] = append(guruMapelMap[item.Pegawai.BidangStudi.Name], pegawaiDetail)
					// Store relasi (use first item's relasi if not set)
					if guruMapelRelasiMap[item.Pegawai.BidangStudi.Name] == "" {
						guruMapelRelasiMap[item.Pegawai.BidangStudi.Name] = item.Relasi
					}
				}
			} else if item.Urutan >= 5 {
				// Urutan 5+: Group by jabatan
				if item.Pegawai.BidangStudi != nil {
					pegawaiDetail.BidangStudi = item.Pegawai.BidangStudi.Name
				}

				if item.Pegawai.Kategori == "Pendidik" {
					kelasMengajar := make([]dtos.RombelWithKelasResponse, 0)
					
					// Add rombel_guru_kelas_id (guru kelas)
					if item.Pegawai.RombelGuruKelas != nil {
						namaKelas := ""
						if item.Pegawai.RombelGuruKelas.Kelas != nil {
							namaKelas = item.Pegawai.RombelGuruKelas.Kelas.Name
						}
						kelasMengajar = append(kelasMengajar, dtos.RombelWithKelasResponse{
							ID:        item.Pegawai.RombelGuruKelas.ID,
							Rombel:    item.Pegawai.RombelGuruKelas.Name,
							NamaKelas: namaKelas,
							Status:    item.Pegawai.RombelGuruKelas.Status,
						})
					}
					
					// Add rombel_bidang_studi (guru bidang studi) - parse JSON array
					var rombelBidangStudiIDs []uint
					if err := json.Unmarshal(item.Pegawai.RombelBidangStudi, &rombelBidangStudiIDs); err == nil {
						// Fetch rombel details for each ID
						for _, rombelID := range rombelBidangStudiIDs {
							rombel, err := s.pegawaiRepo.GetRombelByID(rombelID)
							if err == nil && rombel != nil {
								namaKelas := ""
								if rombel.Kelas != nil {
									namaKelas = rombel.Kelas.Name
								}
								kelasMengajar = append(kelasMengajar, dtos.RombelWithKelasResponse{
									ID:        rombel.ID,
									Rombel:    rombel.Name,
									NamaKelas: namaKelas,
									Status:    rombel.Status,
								})
							}
						}
					}
					
					pegawaiDetail.KelasMengajar = kelasMengajar
				}

				response.Pegawai = &pegawaiDetail
				
				// Group by jabatan for urutan 5+
				if byJabatanMap[item.Urutan] == nil {
					byJabatanMap[item.Urutan] = make(map[string][]dtos.StrukturOrganisasiPublicResponse)
				}
				byJabatanMap[item.Urutan][item.Pegawai.Jabatan] = append(byJabatanMap[item.Urutan][item.Pegawai.Jabatan], response)
			} else {
				// Urutan 1 and 2: tampilkan semua (bidang studi dan kelas mengajar jika ada)
				if item.Pegawai.BidangStudi != nil {
					pegawaiDetail.BidangStudi = item.Pegawai.BidangStudi.Name
				}

				if item.Pegawai.Kategori == "Pendidik" {
					kelasMengajar := make([]dtos.RombelWithKelasResponse, 0)
					
					// Add rombel_guru_kelas_id (guru kelas)
					if item.Pegawai.RombelGuruKelas != nil {
						namaKelas := ""
						if item.Pegawai.RombelGuruKelas.Kelas != nil {
							namaKelas = item.Pegawai.RombelGuruKelas.Kelas.Name
						}
						kelasMengajar = append(kelasMengajar, dtos.RombelWithKelasResponse{
							ID:        item.Pegawai.RombelGuruKelas.ID,
							Rombel:    item.Pegawai.RombelGuruKelas.Name,
							NamaKelas: namaKelas,
							Status:    item.Pegawai.RombelGuruKelas.Status,
						})
					}
					
					// Add rombel_bidang_studi (guru bidang studi) - parse JSON array
					var rombelBidangStudiIDs []uint
					if err := json.Unmarshal(item.Pegawai.RombelBidangStudi, &rombelBidangStudiIDs); err == nil {
						// Fetch rombel details for each ID
						for _, rombelID := range rombelBidangStudiIDs {
							rombel, err := s.pegawaiRepo.GetRombelByID(rombelID)
							if err == nil && rombel != nil {
								namaKelas := ""
								if rombel.Kelas != nil {
									namaKelas = rombel.Kelas.Name
								}
								kelasMengajar = append(kelasMengajar, dtos.RombelWithKelasResponse{
									ID:        rombel.ID,
									Rombel:    rombel.Name,
									NamaKelas: namaKelas,
									Status:    rombel.Status,
								})
							}
						}
					}
					
					pegawaiDetail.KelasMengajar = kelasMengajar
				}

				response.Pegawai = &pegawaiDetail
				// Add to grouped map for urutan 1 and 2
				groupedMap[item.Urutan] = append(groupedMap[item.Urutan], response)
			}
		} else {
			// Non-pegawai data
			if item.Urutan >= 5 {
				// Group by jabatan_non_pegawai for urutan 5+
				if byJabatanMap[item.Urutan] == nil {
					byJabatanMap[item.Urutan] = make(map[string][]dtos.StrukturOrganisasiPublicResponse)
				}
				byJabatanMap[item.Urutan][item.JabatanNonPegawai] = append(byJabatanMap[item.Urutan][item.JabatanNonPegawai], response)
			} else {
				groupedMap[item.Urutan] = append(groupedMap[item.Urutan], response)
			}
		}
	}

	// Convert map to sorted array
	var groupedResponses []dtos.StrukturOrganisasiGroupedResponse
	
	// Add urutan 1 and 2
	for urutan, items := range groupedMap {
		groupedResponses = append(groupedResponses, dtos.StrukturOrganisasiGroupedResponse{
			Urutan: urutan,
			Data:   items,
		})
	}

	// Add urutan 3 with guru kelas grouping
	if len(guruKelasMap) > 0 {
		var guruKelasGroups []dtos.GuruKelasGroupResponse
		for namaKelas, guruList := range guruKelasMap {
			guruKelasGroups = append(guruKelasGroups, dtos.GuruKelasGroupResponse{
				NamaKelas: namaKelas,
				Relasi:    guruKelasRelasiMap[namaKelas],
				Guru:      guruList,
			})
		}
		groupedResponses = append(groupedResponses, dtos.StrukturOrganisasiGroupedResponse{
			Urutan:    3,
			GuruKelas: guruKelasGroups,
		})
	}

	// Add urutan 4 with guru mapel grouping
	if len(guruMapelMap) > 0 {
		var guruMapelGroups []dtos.GuruMapelGroupResponse
		for bidangStudi, guruList := range guruMapelMap {
			guruMapelGroups = append(guruMapelGroups, dtos.GuruMapelGroupResponse{
				BidangStudi: bidangStudi,
				Relasi:      guruMapelRelasiMap[bidangStudi],
				Guru:        guruList,
			})
		}
		groupedResponses = append(groupedResponses, dtos.StrukturOrganisasiGroupedResponse{
			Urutan:    4,
			GuruMapel: guruMapelGroups,
		})
	}

	// Add urutan 5+ with jabatan grouping
	for urutan, jabatanMap := range byJabatanMap {
		var jabatanGroups []dtos.JabatanGroupResponse
		for jabatan, items := range jabatanMap {
			jabatanGroups = append(jabatanGroups, dtos.JabatanGroupResponse{
				Jabatan: jabatan,
				Data:    items,
			})
		}
		groupedResponses = append(groupedResponses, dtos.StrukturOrganisasiGroupedResponse{
			Urutan:    urutan,
			ByJabatan: jabatanGroups,
		})
	}

	// Sort by urutan
	for i := 0; i < len(groupedResponses); i++ {
		for j := i + 1; j < len(groupedResponses); j++ {
			if groupedResponses[i].Urutan > groupedResponses[j].Urutan {
				groupedResponses[i], groupedResponses[j] = groupedResponses[j], groupedResponses[i]
			}
		}
	}

	return groupedResponses, nil
}

// mapToResponse maps model to DTO response
func (s *StrukturOrganisasiServiceImpl) mapToResponse(data *models.StrukturOrganisasi) *dtos.StrukturOrganisasiResponse {
	response := &dtos.StrukturOrganisasiResponse{
		ID:                data.ID,
		PegawaiID:         data.PegawaiID,
		NamaNonPegawai:    data.NamaNonPegawai,
		JabatanNonPegawai: data.JabatanNonPegawai,
		Urutan:            data.Urutan,
		Relasi:            data.Relasi,
		Status:            data.Status,
		CreatedAt:         data.CreatedAt,
		UpdatedAt:         data.UpdatedAt,
		CreatedByID:       data.CreatedByID,
		UpdatedByID:       data.UpdatedByID,
	}

	// Add Pegawai data if available
	if data.Pegawai != nil {
		response.Pegawai = &dtos.PegawaiSimpleResponse{
			ID:      data.Pegawai.ID,
			Nama:    data.Pegawai.Nama,
			Jabatan: data.Pegawai.Jabatan,
			Status:  data.Pegawai.Status,
		}
	}

	return response
}
