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

	"github.com/google/uuid"
)

type FormulirService interface {
	Create(req *dtos.FormulirCreateRequest, dokumenMap map[int]*multipart.FileHeader, userID uint) (*dtos.FormulirResponse, error)
	Update(id uint, req *dtos.FormulirUpdateRequest, dokumenMap map[int]*multipart.FileHeader, userID uint) (*dtos.FormulirResponse, error)
	GetAllWithFilter(params repositories.GetFormulirParams, userID uint) (*dtos.FormulirListWithPaginationResponse, error)
	GetByID(id uint) (*dtos.FormulirResponse, error)
	GetBySlug(slug string) (*dtos.FormulirResponse, error)
	GetFormulirByUser(req *dtos.FormulirGetByUserRequest, userID uint) (*dtos.FormulirListWithPaginationResponse, error)
	Delete(formulirID uint, role *string, userID uint) error
}

type FormulirServiceImpl struct {
	repository repositories.FormulirRepository
	r2Storage  *utils.R2Storage
}

// NewFormulirService creates a new Formulir service
func NewFormulirService(repository repositories.FormulirRepository, r2Storage *utils.R2Storage) FormulirService {
	return &FormulirServiceImpl{
		repository: repository,
		r2Storage:  r2Storage,
	}
}

// generateUUIDSlug creates a unique UUID-based slug
func generateUUIDSlug() string {
	return uuid.New().String()
}

// Create creates a new Formulir with pertanyaan and dokumen uploads
func (s *FormulirServiceImpl) Create(req *dtos.FormulirCreateRequest, dokumenMap map[int]*multipart.FileHeader, userID uint) (*dtos.FormulirResponse, error) {
	// Validate pertanyaan minimal 1
	if len(req.Pertanyaan) == 0 {
		return nil, errors.New("formulir harus memiliki minimal 1 pertanyaan")
	}

	// Validate access_type and target_user_types
	accessType := req.AccessType
	if accessType == "" {
		accessType = "public" // Default
	}
	if accessType != "public" && accessType != "authenticated" {
		return nil, errors.New("access_type harus 'public' atau 'authenticated'")
	}

	// If authenticated, target_user_types should be provided
	if accessType == "authenticated" && len(req.TargetUserTypes) == 0 {
		return nil, errors.New("target_user_types wajib diisi untuk form authenticated")
	}

	// Validate target_user_types values
	validUserTypes := map[string]bool{
		"pendidik":   true,
		"tendik":     true,
		"murid":      true,
		"orang_tua":  true,
		"admin":      true,
	}
	for _, userType := range req.TargetUserTypes {
		if !validUserTypes[userType] {
			return nil, fmt.Errorf("target_user_type '%s' tidak valid. Gunakan: pendidik, tendik, murid, orang_tua, admin", userType)
		}
	}

	// Generate UUID slug (always unique)
	slug := generateUUIDSlug()

	// Parse start_date if provided  
	var startDate *time.Time
	if req.StartDate != "" {
		dateStr := req.StartDate
		if len(dateStr) == 10 { // YYYY-MM-DD
			dateStr = dateStr + " 00:00:00"
		}
		
		// Force parse in UTC using ParseInLocation
		parsed, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, time.UTC)
		if err != nil {
			return nil, errors.New("invalid start_date format (use YYYY-MM-DD or YYYY-MM-DD HH:mm:ss)")
		}
		startDate = &parsed
	}

	// Parse end_date if provided
	var endDate *time.Time
	if req.EndDate != "" {
		dateStr := req.EndDate
		if len(dateStr) == 10 { // YYYY-MM-DD
			dateStr = dateStr + " 23:59:59"
		}
		
		// Force parse in UTC using ParseInLocation
		parsed, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, time.UTC)
		if err != nil {
			return nil, errors.New("invalid end_date format (use YYYY-MM-DD or YYYY-MM-DD HH:mm:ss)")
		}
		endDate = &parsed
	}

	// Convert target_user_types to JSONB
	var targetUserTypesJSON []byte
	if len(req.TargetUserTypes) > 0 {
		targetUserTypesJSON, _ = json.Marshal(req.TargetUserTypes)
	}

	// Convert rombel_ids to JSONB
	var rombelIDsJSON []byte
	if len(req.RombelIDs) > 0 {
		rombelIDsJSON, _ = json.Marshal(req.RombelIDs)
	}

	// Create formulir record
	formulir := &models.Formulir{
		Judul:                  req.Judul,
		Slug:                   slug,
		Deskripsi:              req.Deskripsi,
		CreatedByUserID:        userID,
		CreatedByRole:          req.Role,
		IsActive:               req.IsActive,
		MaxResponses:           req.MaxResponses,
		StartDate:              startDate,
		EndDate:                endDate,
		AccessType:             accessType,
		TargetUserTypes:        targetUserTypesJSON,
		RombelIDs:              rombelIDsJSON,
		AllowMultipleResponses: req.AllowMultipleResponses,
	}

	// Default is_active to true if not provided
	if !req.IsActive {
		formulir.IsActive = true
	}

	// Start database transaction
	tx := s.repository.BeginTransaction()
	
	// Track uploaded files for cleanup on error
	uploadedDokumen := []string{}
	
	// Defer rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			// Cleanup uploaded files on panic
			for _, fileKey := range uploadedDokumen {
				_ = s.r2Storage.DeleteFile(fileKey)
			}
		}
	}()

	// Save formulir
	if err := tx.Create(formulir).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create pertanyaan with dokumen uploads
	for _, p := range req.Pertanyaan {
		// Convert options to JSON if exists
		var optionsJSON []byte
		if len(p.Options) > 0 {
			optionsJSON, _ = json.Marshal(p.Options)
		}

		// Convert validation_rules to JSON if exists
		var validationJSON []byte
		if len(p.ValidationRules) > 0 {
			validationJSON, _ = json.Marshal(p.ValidationRules)
		}

		// Convert file_config to JSON if exists
		var fileConfigJSON []byte
		if len(p.FileConfig) > 0 {
			fileConfigJSON, _ = json.Marshal(p.FileConfig)
		}

		// Handle placeholder
		var placeholder *string
		if p.Placeholder != "" {
			placeholder = &p.Placeholder
		}

		// Handle dokumen upload (matching by urutan)
		var dokumen *string
		if dokumenFile, exists := dokumenMap[p.Urutan]; exists {
			// Validate file size (max 10MB)
			if dokumenFile.Size > 10*1024*1024 {
				tx.Rollback()
				// Cleanup uploaded files on error
				for _, fileKey := range uploadedDokumen {
					_ = s.r2Storage.DeleteFile(fileKey)
				}
				return nil, errors.New("dokumen file must not exceed 10MB")
			}

			// Upload to R2 with formulir-specific folder
			folderPath := "formulir/formulir-" + fmt.Sprint(formulir.ID)
			fileKey, err := s.r2Storage.UploadFile(dokumenFile, folderPath)
			if err != nil {
				tx.Rollback()
				// Cleanup uploaded files on error
				for _, fileKey := range uploadedDokumen {
					_ = s.r2Storage.DeleteFile(fileKey)
				}
				return nil, err
			}
			dokumen = &fileKey
			uploadedDokumen = append(uploadedDokumen, fileKey)
		}

		// Handle link
		var link *string
		if p.Link != "" {
			link = &p.Link
		}

		pertanyaan := &models.FormulirPertanyaan{
			FormulirID:      formulir.ID,
			Urutan:          p.Urutan,
			Label:           p.Label,
			Placeholder:     placeholder,
			Tipe:            p.Tipe,
			IsRequired:      p.IsRequired,
			Options:         optionsJSON,
			ValidationRules: validationJSON,
			FileConfig:      fileConfigJSON,
			Dokumen:         dokumen,
			Link:            link,
		}

		if err := tx.Create(pertanyaan).Error; err != nil {
			tx.Rollback()
			// Cleanup uploaded files on error
			for _, fileKey := range uploadedDokumen {
				_ = s.r2Storage.DeleteFile(fileKey)
			}
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		// Cleanup uploaded files on error
		for _, fileKey := range uploadedDokumen {
			_ = s.r2Storage.DeleteFile(fileKey)
		}
		return nil, err
	}

	// Get created formulir with pertanyaan
	created, err := s.repository.GetByID(formulir.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(created), nil
}

// mapToResponse maps model to DTO response
func (s *FormulirServiceImpl) mapToResponse(data *models.Formulir) *dtos.FormulirResponse {
	// Parse target_user_types from JSONB
	var targetUserTypes []string
	if data.TargetUserTypes != nil {
		_ = json.Unmarshal(data.TargetUserTypes, &targetUserTypes)
	}

	// Parse rombel_ids from JSONB
	var rombelIDs []int
	if data.RombelIDs != nil {
		_ = json.Unmarshal(data.RombelIDs, &rombelIDs)
	}

	// Map pertanyaan
	pertanyaanResponses := make([]dtos.FormulirPertanyaanResponse, 0)
	for _, p := range data.Pertanyaan {
		// Parse options from JSON
		var options []string
		if p.Options != nil {
			_ = json.Unmarshal(p.Options, &options)
		}

		// Parse validation_rules from JSON
		var validationRules map[string]interface{}
		if p.ValidationRules != nil {
			_ = json.Unmarshal(p.ValidationRules, &validationRules)
		}

		// Parse file_config from JSON
		var fileConfig map[string]interface{}
		if p.FileConfig != nil {
			_ = json.Unmarshal(p.FileConfig, &fileConfig)
		}

		// Get placeholder value
		placeholder := ""
		if p.Placeholder != nil {
			placeholder = *p.Placeholder
		}

		// Get dokumen value with public URL
		dokumen := ""
		if p.Dokumen != nil && *p.Dokumen != "" {
			dokumen = s.r2Storage.GetPublicURL(*p.Dokumen)
		}

		// Get link value
		link := ""
		if p.Link != nil {
			link = *p.Link
		}

		pertanyaanResponses = append(pertanyaanResponses, dtos.FormulirPertanyaanResponse{
			ID:              p.ID,
			FormulirID:      p.FormulirID,
			Urutan:          p.Urutan,
			Label:           p.Label,
			Placeholder:     placeholder,
			Tipe:            p.Tipe,
			IsRequired:      p.IsRequired,
			Options:         options,
			ValidationRules: validationRules,
			FileConfig:      fileConfig,
			Dokumen:         dokumen,
			Link:            link,
		})
	}

	return &dtos.FormulirResponse{
		ID:                     data.ID,
		Judul:                  data.Judul,
		Slug:                   data.Slug,
		Deskripsi:              data.Deskripsi,
		CreatedByUserID:        data.CreatedByUserID,
		IsActive:               data.IsActive,
		MaxResponses:           data.MaxResponses,
		StartDate:              data.StartDate,
		EndDate:                data.EndDate,
		AccessType:             data.AccessType,
		TargetUserTypes:        targetUserTypes,
		RombelIDs:              rombelIDs,
		AllowMultipleResponses: data.AllowMultipleResponses,
		CreatedAt:              data.CreatedAt,
		UpdatedAt:              data.UpdatedAt,
		Pertanyaan:             pertanyaanResponses,
		PublicURL:              generatePublicURL(data.Slug, data.AccessType),
	}
}

// generatePublicURL generates the public URL based on access type
func generatePublicURL(slug string, accessType string) string {
	if accessType == "authenticated" {
		return fmt.Sprintf("/backoffice/forms/%s", slug)
	}
	return fmt.Sprintf("/public/forms/%s", slug)
}

// Update updates Formulir with pertanyaan and dokumen uploads
func (s *FormulirServiceImpl) Update(id uint, req *dtos.FormulirUpdateRequest, dokumenMap map[int]*multipart.FileHeader, userID uint) (*dtos.FormulirResponse, error) {
	// Get existing formulir
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("formulir not found")
	}

	// Update basic fields if provided
	if req.Judul != "" {
		existing.Judul = req.Judul
	}
	// Slug tidak bisa diubah setelah dibuat
	if req.Deskripsi != "" {
		existing.Deskripsi = req.Deskripsi
	}
	// Update created_by_role if provided (can be set or nullified)
	if req.Role != nil {
		existing.CreatedByRole = req.Role
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.MaxResponses != nil {
		existing.MaxResponses = req.MaxResponses
	}
	if req.AccessType != "" {
		if req.AccessType != "public" && req.AccessType != "authenticated" {
			return nil, errors.New("access_type harus 'public' atau 'authenticated'")
		}
		existing.AccessType = req.AccessType
		
		// Validate: if authenticated, target_user_types must be provided
		if req.AccessType == "authenticated" && len(req.TargetUserTypes) == 0 {
			return nil, errors.New("target_user_types wajib diisi untuk access_type 'authenticated'")
		}
	}
	
	// Allow updating target_user_types (can be null for public forms)
	if req.TargetUserTypes != nil {
		if len(req.TargetUserTypes) > 0 {
			targetUserTypesJSON, _ := json.Marshal(req.TargetUserTypes)
			existing.TargetUserTypes = targetUserTypesJSON
		} else {
			// Empty array = set to NULL
			existing.TargetUserTypes = nil
		}
	}

	// Allow updating rombel_ids (NULL = all rombels, [] = no restriction, [1,2,3] = specific rombels)
	if req.RombelIDs != nil {
		if len(req.RombelIDs) > 0 {
			rombelIDsJSON, _ := json.Marshal(req.RombelIDs)
			existing.RombelIDs = rombelIDsJSON
		} else {
			// Empty array = set to NULL
			existing.RombelIDs = nil
		}
	}
	
	if req.AllowMultipleResponses != nil {
		existing.AllowMultipleResponses = *req.AllowMultipleResponses
	}

	// Update start_date if provided
	if req.StartDate != "" {
		dateStr := req.StartDate
		if len(dateStr) == 10 {
			dateStr = dateStr + " 00:00:00"
		}
		
		// Force parse in UTC using ParseInLocation
		parsed, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, time.UTC)
		if err != nil {
			return nil, errors.New("invalid start_date format")
		}
		existing.StartDate = &parsed
	}

	// Update end_date if provided
	if req.EndDate != "" {
		dateStr := req.EndDate
		if len(dateStr) == 10 {
			dateStr = dateStr + " 23:59:59"
		}
		
		// Force parse in UTC using ParseInLocation
		parsed, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, time.UTC)
		if err != nil {
			return nil, errors.New("invalid end_date format")
		}
		existing.EndDate = &parsed
	}

	// Start database transaction
	tx := s.repository.BeginTransaction()
	
	// Track uploaded files for cleanup on error
	uploadedDokumen := []string{}
	
	// Defer rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			// Cleanup uploaded files on panic
			for _, fileKey := range uploadedDokumen {
				_ = s.r2Storage.DeleteFile(fileKey)
			}
		}
	}()

	// Save formulir updates
	if err := tx.Save(existing).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// If pertanyaan provided, handle UPDATE/CREATE/DELETE
	if len(req.Pertanyaan) > 0 {
		// Get existing pertanyaan
		oldPertanyaan, _ := s.repository.GetPertanyaanByFormulirID(id)
		
		// Create maps for quick lookup
		oldPertanyaanMap := make(map[uint]*models.FormulirPertanyaan)
		for i := range oldPertanyaan {
			oldPertanyaanMap[oldPertanyaan[i].ID] = &oldPertanyaan[i]
		}
		
		requestedIDs := make(map[uint]bool)
		dokumenToDeleteMap := make(map[int]bool)
		
		for _, urutan := range req.DokumenToDelete {
			dokumenToDeleteMap[urutan] = true
		}
		
		// Process each pertanyaan in request
		for _, p := range req.Pertanyaan {
			// Convert options to JSON if exists
			var optionsJSON []byte
			if len(p.Options) > 0 {
				optionsJSON, _ = json.Marshal(p.Options)
			}

			// Convert validation_rules to JSON if exists
			var validationJSON []byte
			if len(p.ValidationRules) > 0 {
				validationJSON, _ = json.Marshal(p.ValidationRules)
			}

			// Convert file_config to JSON if exists
			var fileConfigJSON []byte
			if len(p.FileConfig) > 0 {
				fileConfigJSON, _ = json.Marshal(p.FileConfig)
			}

			// Handle placeholder
			var placeholder *string
			if p.Placeholder != "" {
				placeholder = &p.Placeholder
			}

			// Handle dokumen upload
			var dokumen *string
			
			// Check if there's a new dokumen file for this urutan
			if dokumenFile, exists := dokumenMap[p.Urutan]; exists {
				// Validate file size (max 10MB)
				if dokumenFile.Size > 10*1024*1024 {
					tx.Rollback()
					for _, fileKey := range uploadedDokumen {
						_ = s.r2Storage.DeleteFile(fileKey)
					}
					return nil, errors.New("dokumen file must not exceed 10MB")
				}

				// Delete old dokumen if exists (when updating)
				if p.ID != nil {
					if oldP, exists := oldPertanyaanMap[*p.ID]; exists && oldP.Dokumen != nil && *oldP.Dokumen != "" {
						_ = s.r2Storage.DeleteFile(*oldP.Dokumen)
					}
				}

				// Upload to R2
				folderPath := "formulir/formulir-" + fmt.Sprint(id)
				fileKey, err := s.r2Storage.UploadFile(dokumenFile, folderPath)
				if err != nil {
					tx.Rollback()
					for _, fileKey := range uploadedDokumen {
						_ = s.r2Storage.DeleteFile(fileKey)
					}
					return nil, err
				}
				dokumen = &fileKey
				uploadedDokumen = append(uploadedDokumen, fileKey)
			} else {
				// Check if dokumen should be deleted
				if dokumenToDeleteMap[p.Urutan] {
					// Delete from R2 if exists
					if p.ID != nil {
						if oldP, exists := oldPertanyaanMap[*p.ID]; exists && oldP.Dokumen != nil && *oldP.Dokumen != "" {
							_ = s.r2Storage.DeleteFile(*oldP.Dokumen)
						}
					}
					dokumen = nil
				} else {
					// Keep old dokumen if exists
					if p.ID != nil {
						if oldP, exists := oldPertanyaanMap[*p.ID]; exists {
							dokumen = oldP.Dokumen
						}
					}
				}
			}

			// Handle link
			var link *string
			if p.Link != "" {
				link = &p.Link
			}

			// UPDATE existing pertanyaan or CREATE new
			if p.ID != nil && oldPertanyaanMap[*p.ID] != nil {
				// UPDATE existing pertanyaan
				requestedIDs[*p.ID] = true
				
				oldP := oldPertanyaanMap[*p.ID]
				oldP.Urutan = p.Urutan
				oldP.Label = p.Label
				oldP.Placeholder = placeholder
				oldP.Tipe = p.Tipe
				oldP.IsRequired = p.IsRequired
				oldP.Options = optionsJSON
				oldP.ValidationRules = validationJSON
				oldP.FileConfig = fileConfigJSON
				oldP.Dokumen = dokumen
				oldP.Link = link

				if err := tx.Save(oldP).Error; err != nil {
					tx.Rollback()
					for _, fileKey := range uploadedDokumen {
						_ = s.r2Storage.DeleteFile(fileKey)
					}
					return nil, err
				}
			} else {
				// CREATE new pertanyaan
				pertanyaan := &models.FormulirPertanyaan{
					FormulirID:      id,
					Urutan:          p.Urutan,
					Label:           p.Label,
					Placeholder:     placeholder,
					Tipe:            p.Tipe,
					IsRequired:      p.IsRequired,
					Options:         optionsJSON,
					ValidationRules: validationJSON,
					FileConfig:      fileConfigJSON,
					Dokumen:         dokumen,
					Link:            link,
				}

				if err := tx.Create(pertanyaan).Error; err != nil {
					tx.Rollback()
					for _, fileKey := range uploadedDokumen {
						_ = s.r2Storage.DeleteFile(fileKey)
					}
					return nil, err
				}
				
				if p.ID != nil {
					requestedIDs[*p.ID] = true
				}
			}
		}
		
		// DELETE pertanyaan that are not in the request
		for _, oldP := range oldPertanyaan {
			if !requestedIDs[oldP.ID] {
				// Delete dokumen from R2 if exists
				if oldP.Dokumen != nil && *oldP.Dokumen != "" {
					_ = s.r2Storage.DeleteFile(*oldP.Dokumen)
				}
				
				// Delete pertanyaan from database
				if err := tx.Delete(&oldP).Error; err != nil {
					tx.Rollback()
					for _, fileKey := range uploadedDokumen {
						_ = s.r2Storage.DeleteFile(fileKey)
					}
					return nil, err
				}
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		// Cleanup uploaded files on error
		for _, fileKey := range uploadedDokumen {
			_ = s.r2Storage.DeleteFile(fileKey)
		}
		return nil, err
	}

	// Get updated formulir with pertanyaan
	updated, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(updated), nil
}

// GetAllWithFilter retrieves all formulir with filters and pagination
func (s *FormulirServiceImpl) GetAllWithFilter(params repositories.GetFormulirParams, userID uint) (*dtos.FormulirListWithPaginationResponse, error) {
	// If role filter is provided
	if params.Filter.CreatedByRole != nil && *params.Filter.CreatedByRole != "" {
		role := *params.Filter.CreatedByRole
		
		// Special case: admin can see all forms
		if role == "admin" {
			// Admin sees everything, clear the role filter to show all forms
			params.Filter.CreatedByRole = nil
			params.Filter.CreatedByID = nil
		} else {
			// For non-admin roles: validate user exists in respective table
			// First try to get user from users table to get username
			user, err := s.repository.GetUserByID(userID)
			if err != nil {
				// If user not in users table, try to get username from respective table based on role
				if role == "pendidik" || role == "tendik" {
					// Get from kepegawaian table using userID
					kepegawaian, err := s.repository.GetKepegawaianByID(userID)
					if err != nil || kepegawaian == nil {
						return nil, fmt.Errorf("user not found in kepegawaian table for role %s", role)
					}
					// User exists in kepegawaian, set filter
					params.Filter.CreatedByID = &userID
				} else if role == "murid" {
					// Get from peserta_didik table using userID
					pesertaDidik, err := s.repository.GetPesertaDidikByID(userID)
					if err != nil || pesertaDidik == nil {
						return nil, errors.New("user not found in peserta_didik table for role murid")
					}
					// User exists in peserta_didik, set filter
					params.Filter.CreatedByID = &userID
				} else {
					// Other roles must exist in users table
					return nil, errors.New("user not found")
				}
			} else {
				// User found in users table, validate in respective table if needed
				username := user.Username
				
				if role == "pendidik" || role == "tendik" {
					// Check in kepegawaian table
					kepegawaian, err := s.repository.GetKepegawaianByUsername(username)
					if err != nil || kepegawaian == nil {
						return nil, fmt.Errorf("user not found in kepegawaian table for role %s", role)
					}
				} else if role == "murid" {
					// Check in peserta_didik table
					pesertaDidik, err := s.repository.GetPesertaDidikByUsername(username)
					if err != nil || pesertaDidik == nil {
						return nil, errors.New("user not found in peserta_didik table for role murid")
					}
				}
				// Other roles just check against users table (already validated above)
				
				// Set created_by_user_id to current user for filtering
				params.Filter.CreatedByID = &userID
			}
		}
	}

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

	// Get data from repository
	formulirs, total, err := s.repository.GetAllWithFilter(params)
	if err != nil {
		return nil, err
	}

	// Map to list response
	responses := make([]dtos.FormulirListResponse, len(formulirs))
	for i, f := range formulirs {
		// Parse target_user_types from JSONB
		var targetUserTypes []string
		if f.TargetUserTypes != nil {
			_ = json.Unmarshal(f.TargetUserTypes, &targetUserTypes)
		}

		// Parse rombel_ids from JSONB
		var rombelIDs []int
		if f.RombelIDs != nil {
			_ = json.Unmarshal(f.RombelIDs, &rombelIDs)
		}

		// Map user data if exists
		var createdBy *dtos.UserBasic
		if f.CreatedBy != nil {
			createdBy = &dtos.UserBasic{
				ID:       f.CreatedBy.ID,
				Username: f.CreatedBy.Username,
				FullName: f.CreatedBy.Nama,
			}
		}

		responses[i] = dtos.FormulirListResponse{
			ID:                     f.ID,
			Judul:                  f.Judul,
			Slug:                   f.Slug,
			CreatedByUserID:        f.CreatedByUserID,
			CreatedBy:              createdBy,
			IsActive:               f.IsActive,
			MaxResponses:           f.MaxResponses,
			StartDate:              f.StartDate,
			EndDate:                f.EndDate,
			AccessType:             f.AccessType,
			TargetUserTypes:        targetUserTypes,
			RombelIDs:              rombelIDs,
			AllowMultipleResponses: f.AllowMultipleResponses,
			CreatedAt:              f.CreatedAt,
			UpdatedAt:              f.UpdatedAt,
			PublicURL:              generatePublicURL(f.Slug, f.AccessType),
		}
	}

	// Calculate pagination metadata
	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.FormulirListWithPaginationResponse{
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

// GetByID retrieves formulir by ID with pertanyaan
func (s *FormulirServiceImpl) GetByID(id uint) (*dtos.FormulirResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// GetBySlug retrieves formulir by slug with pertanyaan
func (s *FormulirServiceImpl) GetBySlug(slug string) (*dtos.FormulirResponse, error) {
	data, err := s.repository.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(data), nil
}

// Delete deletes a formulir with cascade delete (responses, questions) and file cleanup
func (s *FormulirServiceImpl) Delete(formulirID uint, role *string, userID uint) error {
	// Get formulir to check ownership and get related data
	formulir, err := s.repository.GetByID(formulirID)
	if err != nil {
		return errors.New("formulir not found")
	}

	// AUTHORIZATION CHECK
	authorized := false
	authErrorMsg := "unauthorized: you cannot delete this formulir"

	// Rule 1: Admin role can delete anything
	if role != nil && *role == "admin" {
		authorized = true
	} else {
		// Rule 2: Owner can delete their own form
		if formulir.CreatedByUserID == userID {
			// If created_by_role is NULL, owner is from users table (no further check needed)
			if formulir.CreatedByRole == nil {
				authorized = true
			} else {
				// Validate user exists in respective table based on created_by_role
				user, err := s.repository.GetUserByID(userID)
				if err != nil {
					return errors.New("user not found")
				}

				username := user.Username
				createdByRole := *formulir.CreatedByRole

				if createdByRole == "pendidik" || createdByRole == "tendik" {
					// Check in kepegawaian table
					kepegawaian, err := s.repository.GetKepegawaianByUsername(username)
					if err == nil && kepegawaian != nil {
						authorized = true
					} else {
						authErrorMsg = fmt.Sprintf("unauthorized: user not found in kepegawaian table for role %s", createdByRole)
					}
				} else if createdByRole == "murid" {
					// Check in peserta_didik table
					pesertaDidik, err := s.repository.GetPesertaDidikByUsername(username)
					if err == nil && pesertaDidik != nil {
						authorized = true
					} else {
						authErrorMsg = "unauthorized: user not found in peserta_didik table for role murid"
					}
				} else {
					// Other roles (orang_tua, admin, etc.) just check user_id match
					authorized = true
				}
			}
		} else {
			authErrorMsg = "unauthorized: you are not the owner of this formulir"
		}
	}

	if !authorized {
		return errors.New(authErrorMsg)
	}

	// Start transaction
	tx := s.repository.BeginTransaction()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Track files to delete
	filesToDelete := []string{}

	// Step 1: Get all responses for this form
	responses, err := s.repository.GetAllResponsesByFormulirID(formulirID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get responses: %v", err)
	}

	// Step 2: Delete all response answers and collect file URLs
	for _, response := range responses {
		// Get jawaban for this response
		jawaban, err := s.repository.GetJawabanByResponseID(response.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get jawaban for response %d: %v", response.ID, err)
		}

		// Collect file URLs from jawaban
		for _, j := range jawaban {
			if j.JawabanText != nil && *j.JawabanText != "" {
				// Only add if it looks like a file path
				if len(*j.JawabanText) > 0 {
					filesToDelete = append(filesToDelete, *j.JawabanText)
				}
			}
		}

		// Delete all jawaban for this response
		if err := tx.Where("response_id = ?", response.ID).Delete(&models.FormulirResponseJawaban{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete jawaban for response %d: %v", response.ID, err)
		}
	}

	// Step 3: Delete all responses
	if err := tx.Where("form_id = ?", formulirID).Delete(&models.FormulirResponse{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete responses: %v", err)
	}

	// Step 4: Delete all questions and collect dokumen files
	for _, pertanyaan := range formulir.Pertanyaan {
		if pertanyaan.Dokumen != nil && *pertanyaan.Dokumen != "" {
			filesToDelete = append(filesToDelete, *pertanyaan.Dokumen)
		}
	}

	if err := tx.Where("form_id = ?", formulirID).Delete(&models.FormulirPertanyaan{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete questions: %v", err)
	}

	// Step 5: Delete the form itself
	if err := tx.Delete(formulir).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete formulir: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Step 6: Delete all files from R2 (after successful transaction)
	for _, fileKey := range filesToDelete {
		_ = s.r2Storage.DeleteFile(fileKey)
		// Ignore errors for file deletion - file might not exist
	}

	return nil
}

// GetFormulirByUser retrieves formulir filtered by user role and context
func (s *FormulirServiceImpl) GetFormulirByUser(req *dtos.FormulirGetByUserRequest, userID uint) (*dtos.FormulirListWithPaginationResponse, error) {
	// Validate role
	if req.Role != "pendidik" && req.Role != "tendik" && req.Role != "murid" {
		return nil, errors.New("role harus 'pendidik', 'tendik', atau 'murid'")
	}

	// Validate rombel_id for murid
	if req.Role == "murid" && req.RombelID == nil {
		return nil, errors.New("rombel_id wajib diisi untuk role murid")
	}

	// Validate user exists in respective table
	user, err := s.repository.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	username := user.Username

	// Check user in respective table based on role
	if req.Role == "pendidik" || req.Role == "tendik" {
		// Check in kepegawaian table
		kepegawaian, err := s.repository.GetKepegawaianByUsername(username)
		if err != nil || kepegawaian == nil {
			return nil, fmt.Errorf("user tidak ditemukan di tabel kepegawaian untuk role %s", req.Role)
		}
	} else if req.Role == "murid" {
		// Check in peserta_didik table
		pesertaDidik, err := s.repository.GetPesertaDidikByUsername(username)
		if err != nil || pesertaDidik == nil {
			return nil, errors.New("user tidak ditemukan di tabel peserta_didik untuk role murid")
		}
	}

	// Validate and set default pagination
	if req.Pagination.Page < 1 {
		req.Pagination.Page = 1
	}
	if req.Pagination.Limit < 1 {
		req.Pagination.Limit = 10
	}
	if req.Pagination.Limit > 100 {
		req.Pagination.Limit = 100
	}

	// Calculate offset
	offset := (req.Pagination.Page - 1) * req.Pagination.Limit

	// Parse date filters
	var startDate, endDate time.Time
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, errors.New("format start_date tidak valid (gunakan YYYY-MM-DD)")
		}
		startDate = parsed
	}

	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, errors.New("format end_date tidak valid (gunakan YYYY-MM-DD)")
		}
		// Set to end of day
		endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// Call repository
	formulirs, total, err := s.repository.GetFormulirByUser(req.Role, req.RombelID, startDate, endDate, req.Judul, req.Pagination.Limit, offset)
	if err != nil {
		return nil, err
	}

	// Map to list response
	responses := make([]dtos.FormulirListResponse, len(formulirs))
	for i, f := range formulirs {
		// Parse target_user_types from JSONB
		var targetUserTypes []string
		if f.TargetUserTypes != nil {
			_ = json.Unmarshal(f.TargetUserTypes, &targetUserTypes)
		}

		// Parse rombel_ids from JSONB
		var rombelIDs []int
		if f.RombelIDs != nil {
			_ = json.Unmarshal(f.RombelIDs, &rombelIDs)
		}

		// Map user data if exists
		var createdBy *dtos.UserBasic
		if f.CreatedBy != nil {
			createdBy = &dtos.UserBasic{
				ID:       f.CreatedBy.ID,
				Username: f.CreatedBy.Username,
				FullName: f.CreatedBy.Nama,
			}
		}

		responses[i] = dtos.FormulirListResponse{
			ID:                     f.ID,
			Judul:                  f.Judul,
			Slug:                   f.Slug,
			CreatedByUserID:        f.CreatedByUserID,
			CreatedBy:              createdBy,
			IsActive:               f.IsActive,
			MaxResponses:           f.MaxResponses,
			StartDate:              f.StartDate,
			EndDate:                f.EndDate,
			AccessType:             f.AccessType,
			TargetUserTypes:        targetUserTypes,
			RombelIDs:              rombelIDs,
			AllowMultipleResponses: f.AllowMultipleResponses,
			CreatedAt:              f.CreatedAt,
			UpdatedAt:              f.UpdatedAt,
			PublicURL:              generatePublicURL(f.Slug, f.AccessType),
		}
	}

	// Calculate pagination metadata
	totalPages := (int(total) + req.Pagination.Limit - 1) / req.Pagination.Limit

	return &dtos.FormulirListWithPaginationResponse{
		Data: responses,
		Pagination: dtos.PaginationInfo{
			Limit:      req.Pagination.Limit,
			Offset:     offset,
			Page:       req.Pagination.Page,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
