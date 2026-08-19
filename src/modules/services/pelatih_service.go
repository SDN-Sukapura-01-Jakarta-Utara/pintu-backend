package services

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"strings"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"

	"golang.org/x/crypto/bcrypt"
)

// PelatihService handles business logic for Pelatih
type PelatihService interface {
	Create(req *dtos.PelatihCreateRequest, userID uint) (*dtos.PelatihResponse, error)
	GetAllWithFilter(params repositories.GetPelatihParams) (*dtos.PelatihListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.PelatihResponse, error)
	Update(req *dtos.PelatihUpdateRequest, fotoProfil *multipart.FileHeader, sertifikatFiles []*multipart.FileHeader, userID uint) (*dtos.PelatihResponse, error)
	Delete(id uint) error
}

type PelatihServiceImpl struct {
	repository          repositories.PelatihRepository
	ekstrakurikulerRepo repositories.EkstrakurikulerRepository
	r2Storage           *utils.R2Storage
}

// NewPelatihService creates a new service
func NewPelatihService(
	repository repositories.PelatihRepository,
	ekstrakurikulerRepo repositories.EkstrakurikulerRepository,
) PelatihService {
	return &PelatihServiceImpl{
		repository:          repository,
		ekstrakurikulerRepo: ekstrakurikulerRepo,
		r2Storage:           utils.NewSieksaR2Storage(), // Use SIEKSA R2 bucket
	}
}

func (s *PelatihServiceImpl) Create(req *dtos.PelatihCreateRequest, userID uint) (*dtos.PelatihResponse, error) {
	// Validate all ekstrakurikuler exist and are active
	if len(req.EkstrakurikulerIDs) == 0 {
		return nil, errors.New("minimal harus ada 1 ekstrakurikuler")
	}

	for _, ekskulID := range req.EkstrakurikulerIDs {
		ekskul, err := s.ekstrakurikulerRepo.GetByID(ekskulID)
		if err != nil {
			return nil, errors.New("ekstrakurikuler tidak ditemukan")
		}

		if ekskul.Status != "active" {
			return nil, errors.New("ekstrakurikuler tidak aktif")
		}
	}

	// Hash password if provided
	var hashedPassword string
	if req.Password != nil && *req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("gagal mengenkripsi password")
		}
		hashedPassword = string(hashed)
	}

	// Check username uniqueness if provided
	if req.Username != nil && *req.Username != "" {
		existing, _ := s.repository.GetByUsername(*req.Username)
		if existing != nil {
			return nil, errors.New("username sudah digunakan")
		}
	}

	// Create pelatih
	pelatih := &models.Pelatih{
		Nama:        req.Nama,
		Username:    req.Username,
		Password:    hashedPassword,
		Telepon:     req.Telepon,
		Alamat:      req.Alamat,
		Keahlian:    req.Keahlian,
		Status:      req.Status,
		CreatedByID: &userID,
	}

	if err := s.repository.Create(pelatih); err != nil {
		return nil, errors.New("gagal menyimpan data pelatih")
	}

	// Create mappings to all ekstrakurikuler
	for _, ekskulID := range req.EkstrakurikulerIDs {
		mapping := &models.PelatihEkstrakurikuler{
			PelatihID:         pelatih.ID,
			EkstrakurikulerID: ekskulID,
			Status:            "active",
			CreatedByID:       &userID,
		}

		if err := s.repository.CreatePelatihEkstrakurikuler(mapping); err != nil {
			return nil, errors.New("gagal menyimpan mapping pelatih ke ekstrakurikuler")
		}
	}

	// Handle role assignment - role_ids is required
	if len(req.RoleIDs) == 0 {
		return nil, errors.New("role_ids wajib diisi, minimal 1 role")
	}

	// Validate all role IDs exist and belong to SIEKSA (SystemID=2)
	for _, roleID := range req.RoleIDs {
		role, err := s.repository.GetRoleByID(roleID)
		if err != nil {
			return nil, errors.New("role tidak ditemukan")
		}
		if role.SystemID == nil || *role.SystemID != 2 {
			return nil, errors.New("role harus memiliki system_id = 2 (SIEKSA)")
		}
	}

	// Create role mappings
	for _, roleID := range req.RoleIDs {
		roleMapping := &models.PelatihRole{
			PelatihID: pelatih.ID,
			RoleID:    roleID,
		}

		if err := s.repository.CreatePelatihRole(roleMapping); err != nil {
			return nil, errors.New("gagal menyimpan role pelatih")
		}
	}

	// Get pelatih with relations
	pelatihWithRelations, err := s.repository.GetByID(pelatih.ID)
	if err != nil {
		return nil, errors.New("gagal mengambil data pelatih")
	}

	return s.mapToResponse(pelatihWithRelations), nil
}

func (s *PelatihServiceImpl) GetByID(id uint) (*dtos.PelatihResponse, error) {
	pelatih, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("pelatih tidak ditemukan")
	}

	return s.mapToResponse(pelatih), nil
}

func (s *PelatihServiceImpl) Update(req *dtos.PelatihUpdateRequest, fotoProfil *multipart.FileHeader, sertifikatFiles []*multipart.FileHeader, userID uint) (*dtos.PelatihResponse, error) {
	// Get existing pelatih
	existing, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("pelatih tidak ditemukan")
	}

	// Update fields if provided
	if req.Nama != nil {
		existing.Nama = *req.Nama
	}
	if req.Username != nil {
		// Check username uniqueness if changed
		if existing.Username == nil || *existing.Username != *req.Username {
			existingUser, _ := s.repository.GetByUsername(*req.Username)
			if existingUser != nil && existingUser.ID != req.ID {
				return nil, errors.New("username sudah digunakan")
			}
		}
		existing.Username = req.Username
	}
	if req.Password != nil && *req.Password != "" {
		// Hash new password
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("gagal mengenkripsi password")
		}
		existing.Password = string(hashed)
	}
	if req.Telepon != nil {
		existing.Telepon = *req.Telepon
	}
	if req.Alamat != nil {
		existing.Alamat = *req.Alamat
	}
	if req.Keahlian != nil {
		existing.Keahlian = *req.Keahlian
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}

	// Handle foto_profil upload
	if fotoProfil != nil {
		// Delete old foto_profil if exists
		if existing.FotoProfil != nil && *existing.FotoProfil != "" {
			_ = s.r2Storage.DeleteFile(*existing.FotoProfil)
		}

		// Upload new foto_profil
		fotoURL, err := s.r2Storage.UploadFile(fotoProfil, "pelatih/foto_profil")
		if err != nil {
			return nil, errors.New("gagal mengupload foto profil: " + err.Error())
		}
		existing.FotoProfil = &fotoURL
	}

	// Handle sertifikat uploads (multiple files)
	var currentSertifikat []string
	
	// Get current sertifikat list from database (stored as paths)
	if existing.Sertifikat != nil && *existing.Sertifikat != "" {
		if err := json.Unmarshal([]byte(*existing.Sertifikat), &currentSertifikat); err == nil {
			// Remove sertifikat that should be deleted
			if len(req.SertifikatToDelete) > 0 {
				// Convert full URLs to paths for comparison
				deletePathMap := make(map[string]bool)
				for _, url := range req.SertifikatToDelete {
					// Extract path from URL (remove domain prefix)
					// URL: https://sieksa-storage.sdnsukapura01dev.my.id/pelatih/sertifikat/123.pdf
					// Path: pelatih/sertifikat/123.pdf
					path := s.extractPathFromURL(url)
					deletePathMap[path] = true
				}

				// Filter out deleted files and delete from R2
				var remainingSertifikat []string
				for _, path := range currentSertifikat {
					if deletePathMap[path] {
						// Delete from R2 using path
						_ = s.r2Storage.DeleteFile(path)
					} else {
						remainingSertifikat = append(remainingSertifikat, path)
					}
				}
				currentSertifikat = remainingSertifikat
			}
		}
	}

	// Upload new sertifikat files if provided
	if len(sertifikatFiles) > 0 {
		for _, file := range sertifikatFiles {
			fileURL, err := s.r2Storage.UploadFile(file, "pelatih/sertifikat")
			if err != nil {
				return nil, errors.New("gagal mengupload sertifikat: " + err.Error())
			}
			currentSertifikat = append(currentSertifikat, fileURL)
		}
	}

	// Save updated sertifikat list to database
	if len(currentSertifikat) > 0 || len(req.SertifikatToDelete) > 0 || len(sertifikatFiles) > 0 {
		sertifikatJSON, err := json.Marshal(currentSertifikat)
		if err != nil {
			return nil, errors.New("gagal menyimpan data sertifikat")
		}
		sertifikatStr := string(sertifikatJSON)
		existing.Sertifikat = &sertifikatStr
	}

	existing.UpdatedByID = &userID

	// Update pelatih
	if err := s.repository.Update(existing); err != nil {
		return nil, errors.New("gagal mengupdate data pelatih")
	}

	// Update ekstrakurikuler assignments if provided
	if len(req.EkstrakurikulerIDs) > 0 {
		// Validate all ekstrakurikuler exist and are active
		for _, ekskulID := range req.EkstrakurikulerIDs {
			ekskul, err := s.ekstrakurikulerRepo.GetByID(ekskulID)
			if err != nil {
				return nil, errors.New("ekstrakurikuler tidak ditemukan")
			}
			if ekskul.Status != "active" {
				return nil, errors.New("ekstrakurikuler tidak aktif")
			}
		}

		// Get current active ekstrakurikuler mappings
		currentEkskulIDs, err := s.repository.GetActivePelatihEkstrakurikuler(req.ID)
		if err != nil {
			return nil, errors.New("gagal mengambil data ekstrakurikuler saat ini")
		}

		// Convert to maps for easy lookup
		currentMap := make(map[uint]bool)
		for _, id := range currentEkskulIDs {
			currentMap[id] = true
		}
		
		newMap := make(map[uint]bool)
		for _, id := range req.EkstrakurikulerIDs {
			newMap[id] = true
		}

		// Soft delete ekstrakurikuler that are removed
		for _, ekskulID := range currentEkskulIDs {
			if !newMap[ekskulID] {
				if err := s.repository.SoftDeletePelatihEkstrakurikuler(req.ID, ekskulID); err != nil {
					return nil, errors.New("gagal menghapus mapping ekstrakurikuler")
				}
			}
		}

		// Add or restore new ekstrakurikuler
		for _, ekskulID := range req.EkstrakurikulerIDs {
			if !currentMap[ekskulID] {
				// Try to restore first (if previously deleted)
				if err := s.repository.RestorePelatihEkstrakurikuler(req.ID, ekskulID); err != nil {
					// If restore fails, create new
					mapping := &models.PelatihEkstrakurikuler{
						PelatihID:         req.ID,
						EkstrakurikulerID: ekskulID,
						Status:            "active",
						CreatedByID:       &userID,
					}

					if err := s.repository.CreatePelatihEkstrakurikuler(mapping); err != nil {
						return nil, errors.New("gagal menyimpan mapping pelatih ke ekstrakurikuler")
					}
				}
			}
		}
	}

	// Update role assignments if provided
	if len(req.RoleIDs) > 0 {
		// Validate all roles exist and belong to SIEKSA (SystemID=2)
		for _, roleID := range req.RoleIDs {
			role, err := s.repository.GetRoleByID(roleID)
			if err != nil {
				return nil, errors.New("role tidak ditemukan")
			}
			if role.SystemID == nil || *role.SystemID != 2 {
				return nil, errors.New("role harus memiliki system_id = 2 (SIEKSA)")
			}
		}

		// Delete all existing role mappings
		if err := s.repository.DeleteAllPelatihRoles(req.ID); err != nil {
			return nil, errors.New("gagal menghapus role lama")
		}

		// Create new role mappings
		for _, roleID := range req.RoleIDs {
			roleMapping := &models.PelatihRole{
				PelatihID: req.ID,
				RoleID:    roleID,
			}

			if err := s.repository.CreatePelatihRole(roleMapping); err != nil {
				return nil, errors.New("gagal menyimpan role pelatih")
			}
		}
	}

	// Get updated pelatih with relations
	updatedPelatih, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("gagal mengambil data pelatih")
	}

	return s.mapToResponse(updatedPelatih), nil
}

func (s *PelatihServiceImpl) Delete(id uint) error {
	// Check if pelatih exists
	_, err := s.repository.GetByID(id)
	if err != nil {
		return errors.New("pelatih tidak ditemukan")
	}

	// Delete all ekstrakurikuler mappings first
	if err := s.repository.DeleteAllPelatihEkstrakurikuler(id); err != nil {
		return errors.New("gagal menghapus mapping ekstrakurikuler")
	}

	// Delete all role mappings
	if err := s.repository.DeleteAllPelatihRoles(id); err != nil {
		return errors.New("gagal menghapus role pelatih")
	}

	// Delete pelatih
	if err := s.repository.Delete(id); err != nil {
		return errors.New("gagal menghapus data pelatih")
	}

	return nil
}

func (s *PelatihServiceImpl) GetAllWithFilter(params repositories.GetPelatihParams) (*dtos.PelatihListWithPaginationResponse, error) {
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
	responses := make([]dtos.PelatihResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.PelatihListWithPaginationResponse{
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

// mapToResponse maps model to DTO response
func (s *PelatihServiceImpl) mapToResponse(pelatih *models.Pelatih) *dtos.PelatihResponse {
	response := &dtos.PelatihResponse{
		ID:        pelatih.ID,
		Nama:      pelatih.Nama,
		Username:  pelatih.Username,
		Telepon:   pelatih.Telepon,
		Alamat:    pelatih.Alamat,
		Keahlian:  pelatih.Keahlian,
		Status:    pelatih.Status,
		CreatedAt: pelatih.CreatedAt,
		UpdatedAt: pelatih.UpdatedAt,
	}

	// Convert foto_profil path to full URL
	if pelatih.FotoProfil != nil && *pelatih.FotoProfil != "" {
		fullURL := s.r2Storage.GetPublicURL(*pelatih.FotoProfil)
		response.FotoProfil = &fullURL
	}

	// Parse sertifikat JSON and convert paths to full URLs
	if pelatih.Sertifikat != nil && *pelatih.Sertifikat != "" {
		var sertifikatPaths []string
		if err := json.Unmarshal([]byte(*pelatih.Sertifikat), &sertifikatPaths); err == nil {
			// Convert each path to full URL
			sertifikatURLs := make([]string, len(sertifikatPaths))
			for i, path := range sertifikatPaths {
				sertifikatURLs[i] = s.r2Storage.GetPublicURL(path)
			}
			response.Sertifikat = sertifikatURLs
		}
	}

	// Map ekstrakurikuler
	if len(pelatih.Ekstrakurikuler) > 0 {
		response.Ekstrakurikuler = make([]dtos.EkstrakurikulerResponse, len(pelatih.Ekstrakurikuler))
		for i, ekskul := range pelatih.Ekstrakurikuler {
			response.Ekstrakurikuler[i] = dtos.EkstrakurikulerResponse{
				ID:       ekskul.ID,
				Name:     ekskul.Name,
				Kategori: ekskul.Kategori,
				Status:   ekskul.Status,
			}
		}
	}

	// Map roles
	if len(pelatih.Roles) > 0 {
		response.Roles = make([]dtos.RoleResponse, len(pelatih.Roles))
		for i, role := range pelatih.Roles {
			var system *dtos.SystemResponse
			if role.System != nil {
				system = &dtos.SystemResponse{
					ID:          role.System.ID,
					Nama:        role.System.Nama,
					Description: role.System.Description,
				}
			}

			response.Roles[i] = dtos.RoleResponse{
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
	}

	return response
}


// extractPathFromURL extracts the file path from a full URL
// Example: "https://sieksa-storage.sdnsukapura01dev.my.id/pelatih/sertifikat/123.pdf" -> "pelatih/sertifikat/123.pdf"
func (s *PelatihServiceImpl) extractPathFromURL(urlStr string) string {
	// If already a path (no http/https), return as is
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return urlStr
	}
	
	// Remove the protocol and domain, keep only the path
	// Split by "/" and skip protocol and domain parts
	parts := strings.Split(urlStr, "/")
	if len(parts) > 3 {
		// parts[0] = "https:"
		// parts[1] = ""
		// parts[2] = "sieksa-storage.sdnsukapura01dev.my.id"
		// parts[3+] = actual path
		return strings.Join(parts[3:], "/")
	}
	
	return urlStr
}
