package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"sync"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

type KepegawaianService interface {
	Create(req *dtos.KepegawaianCreateRequest, userID uint) (*dtos.KepegawaianResponse, error)
	GetByID(id uint) (*dtos.KepegawaianResponse, error)
	GetByNIP(nip string) (*dtos.KepegawaianResponse, error)
	GetAll(limit int, offset int) (*dtos.KepegawaianListResponse, error)
	GetAllWithFilter(params repositories.GetKepegawaianParams) (*dtos.KepegawaianListWithPaginationResponse, error)
	Update(id uint, foto *multipart.FileHeader, docs map[string][]*multipart.FileHeader, req *dtos.KepegawaianUpdateRequest, userID uint) (*dtos.KepegawaianResponse, error)
	Delete(id uint) error
}

type KepegawaianServiceImpl struct {
	repository repositories.KepegawaianRepository
	r2Storage  *utils.R2Storage
}

// NewKepegawaianService creates a new Kepegawaian service
func NewKepegawaianService(repository repositories.KepegawaianRepository, r2Storage *utils.R2Storage) KepegawaianService {
	return &KepegawaianServiceImpl{
		repository: repository,
		r2Storage:  r2Storage,
	}
}

// Create creates a new Kepegawaian (without file uploads)
func (s *KepegawaianServiceImpl) Create(req *dtos.KepegawaianCreateRequest, userID uint) (*dtos.KepegawaianResponse, error) {
	// Check if NIP already exists
	existing, _ := s.repository.GetByNIP(req.NIP)
	if existing != nil {
		return nil, errors.New("NIP already exists")
	}

	// Check if username already exists
	existing, _ = s.repository.GetByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Parse rombel_bidang_studi
	rombelBidangStudiJSON, _ := json.Marshal(req.RombelBidangStudi)

	// Create kepegawaian record
	data := &models.Kepegawaian{
		Nama:              req.Nama,
		Username:          req.Username,
		Password:          string(hashedPassword),
		NIP:               req.NIP,
		NKKI:              req.NKKI,
		Kategori:          req.Kategori,
		Jabatan:           req.Jabatan,
		BidangStudiID:     req.BidangStudiID,
		RombelGuruKelasID: req.RombelGuruKelasID,
		RombelBidangStudi: rombelBidangStudiJSON,
		Status:            status,
		CreatedByID:       &userID,
	}

	if err := s.repository.Create(data); err != nil {
		return nil, err
	}

	// Assign roles if provided
	if len(req.RoleIDs) > 0 {
		if err := s.repository.AssignRoles(data.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	return s.mapToResponse(data), nil
}

// GetByID retrieves Kepegawaian by ID
func (s *KepegawaianServiceImpl) GetByID(id uint) (*dtos.KepegawaianResponse, error) {
	data, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetByNIP retrieves Kepegawaian by NIP
func (s *KepegawaianServiceImpl) GetByNIP(nip string) (*dtos.KepegawaianResponse, error) {
	data, err := s.repository.GetByNIP(nip)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(data), nil
}

// GetAll retrieves all Kepegawaian
func (s *KepegawaianServiceImpl) GetAll(limit int, offset int) (*dtos.KepegawaianListResponse, error) {
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
	responses := make([]dtos.KepegawaianResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	return &dtos.KepegawaianListResponse{
		Data:   responses,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// GetAllWithFilter retrieves Kepegawaian with filters and pagination
func (s *KepegawaianServiceImpl) GetAllWithFilter(params repositories.GetKepegawaianParams) (*dtos.KepegawaianListWithPaginationResponse, error) {
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
	responses := make([]dtos.KepegawaianResponse, len(data))
	for i, item := range data {
		responses[i] = *s.mapToResponse(&item)
	}

	totalPages := (int(total) + params.Limit - 1) / params.Limit

	return &dtos.KepegawaianListWithPaginationResponse{
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

// Update updates Kepegawaian
func (s *KepegawaianServiceImpl) Update(id uint, foto *multipart.FileHeader, docs map[string][]*multipart.FileHeader, req *dtos.KepegawaianUpdateRequest, userID uint) (*dtos.KepegawaianResponse, error) {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("kepegawaian not found")
	}

	oldFoto := existing.Foto

	// Update basic fields if provided
	if req.Nama != "" {
		existing.Nama = req.Nama
	}
	if req.Username != "" {
		// Check if username already exists for other users
		existingUser, _ := s.repository.GetByUsername(req.Username)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("username already exists")
		}
		existing.Username = req.Username
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		existing.Password = string(hashedPassword)
	}
	if req.NIP != "" {
		// Check if NIP already exists for other users
		existingNIP, _ := s.repository.GetByNIP(req.NIP)
		if existingNIP != nil && existingNIP.ID != id {
			return nil, errors.New("NIP already exists")
		}
		existing.NIP = req.NIP
	}
	if req.NKKI != "" {
		existing.NKKI = req.NKKI
	}
	if req.Kategori != "" {
		existing.Kategori = req.Kategori
	}
	if req.Jabatan != "" {
		existing.Jabatan = req.Jabatan
	}
	if req.BidangStudiID != nil {
		existing.BidangStudiID = req.BidangStudiID
	}
	if req.RombelGuruKelasID != nil {
		existing.RombelGuruKelasID = req.RombelGuruKelasID
	}
	if len(req.RombelBidangStudi) > 0 {
		rombelBidangStudiJSON, _ := json.Marshal(req.RombelBidangStudi)
		existing.RombelBidangStudi = rombelBidangStudiJSON
	}
	if req.Status != "" {
		existing.Status = req.Status
	}

	// Update foto if provided
	if foto != nil {
		if foto.Size > 5*1024*1024 { // 5MB
			return nil, errors.New("foto size must not exceed 5MB")
		}

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
		}
		contentType := foto.Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			return nil, errors.New("only image files are allowed for foto (jpeg, png, gif, webp)")
		}

		newFileKey, err := s.r2Storage.UploadFile(foto, "kepegawaian/foto")
		if err != nil {
			return nil, err
		}

		// Delete old foto if exists
		if oldFoto != "" {
			_ = s.r2Storage.DeleteFile(oldFoto)
		}

		existing.Foto = newFileKey
	}

	// Delete documents if specified
	if len(req.SertifikatLainnyaToDelete) > 0 {
		s.deleteDocumentsFromJSONB(&existing.SertifikatLainnya, req.SertifikatLainnyaToDelete)
	}
	if len(req.DokumenLainnyaToDelete) > 0 {
		s.deleteDocumentsFromJSONB(&existing.DokumenLainnya, req.DokumenLainnyaToDelete)
	}

	// Update documents if provided (parallel)
	if len(docs) > 0 {
		uploadResults := s.uploadDocumentsParallel(docs)
		
		for docType, result := range uploadResults {
			if result.err != nil {
				return nil, result.err
			}

			if docType != "sertifikat_lainnya" && docType != "dokumen_lainnya" {
				// Single file result
				switch docType {
				case "kk":
					if existing.KK != "" {
						_ = s.r2Storage.DeleteFile(existing.KK)
					}
					existing.KK = result.fileKey
				case "akta_lahir":
					if existing.AktaLahir != "" {
						_ = s.r2Storage.DeleteFile(existing.AktaLahir)
					}
					existing.AktaLahir = result.fileKey
				case "ktp":
					if existing.KTP != "" {
						_ = s.r2Storage.DeleteFile(existing.KTP)
					}
					existing.KTP = result.fileKey
				case "ijazah_sd":
					if existing.IjazahSD != "" {
						_ = s.r2Storage.DeleteFile(existing.IjazahSD)
					}
					existing.IjazahSD = result.fileKey
				case "ijazah_smp":
					if existing.IjazahSMP != "" {
						_ = s.r2Storage.DeleteFile(existing.IjazahSMP)
					}
					existing.IjazahSMP = result.fileKey
				case "ijazah_sma":
					if existing.IjazahSMA != "" {
						_ = s.r2Storage.DeleteFile(existing.IjazahSMA)
					}
					existing.IjazahSMA = result.fileKey
				case "ijazah_s1":
					if existing.IjazahS1 != "" {
						_ = s.r2Storage.DeleteFile(existing.IjazahS1)
					}
					existing.IjazahS1 = result.fileKey
				case "ijazah_s2":
					if existing.IjazahS2 != "" {
						_ = s.r2Storage.DeleteFile(existing.IjazahS2)
					}
					existing.IjazahS2 = result.fileKey
				case "ijazah_s3":
					if existing.IjazahS3 != "" {
						_ = s.r2Storage.DeleteFile(existing.IjazahS3)
					}
					existing.IjazahS3 = result.fileKey
				case "sertifikat_pendidik":
					if existing.SertifikatPendidik != "" {
						_ = s.r2Storage.DeleteFile(existing.SertifikatPendidik)
					}
					existing.SertifikatPendidik = result.fileKey
				case "sk":
					if existing.SK != "" {
						_ = s.r2Storage.DeleteFile(existing.SK)
					}
					existing.SK = result.fileKey
				}
			} else {
				// Multiple file result
				fileKeysJSON, _ := json.Marshal(result.fileKeys)
				if docType == "sertifikat_lainnya" {
					existing.SertifikatLainnya = fileKeysJSON
				} else if docType == "dokumen_lainnya" {
					existing.DokumenLainnya = fileKeysJSON
				}
			}
		}
	}

	existing.UpdatedByID = &userID

	if err := s.repository.Update(existing); err != nil {
		return nil, err
	}

	// Assign roles if provided
	if len(req.RoleIDs) > 0 {
		if err := s.repository.AssignRoles(id, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	return s.mapToResponse(existing), nil
}

// Delete deletes Kepegawaian by ID
func (s *KepegawaianServiceImpl) Delete(id uint) error {
	// Get existing data
	existing, err := s.repository.GetByID(id)
	if err != nil {
		return errors.New("kepegawaian not found")
	}

	// Delete foto from R2
	if existing.Foto != "" {
		_ = s.r2Storage.DeleteFile(existing.Foto)
	}

	// Delete all document files from R2
	s.deleteDocumentFiles(existing)

	// Delete from database
	return s.repository.Delete(id)
}

// Helper function to delete all document files
func (s *KepegawaianServiceImpl) deleteDocumentFiles(data *models.Kepegawaian) {
	// Delete single file documents
	files := []string{data.KK, data.AktaLahir, data.KTP, data.IjazahSD, data.IjazahSMP,
		data.IjazahSMA, data.IjazahS1, data.IjazahS2, data.IjazahS3, data.SertifikatPendidik, data.SK}
	for _, file := range files {
		if file != "" {
			_ = s.r2Storage.DeleteFile(file)
		}
	}

	// Delete multiple file documents
	var certFiles []string
	json.Unmarshal(data.SertifikatLainnya, &certFiles)
	for _, file := range certFiles {
		if file != "" {
			_ = s.r2Storage.DeleteFile(file)
		}
	}

	var dokFiles []string
	json.Unmarshal(data.DokumenLainnya, &dokFiles)
	for _, file := range dokFiles {
		if file != "" {
			_ = s.r2Storage.DeleteFile(file)
		}
	}
}

// Helper function to delete documents from JSONB field
func (s *KepegawaianServiceImpl) deleteDocumentsFromJSONB(jsonbField *datatypes.JSON, filesToDelete []string) {
	var fileKeys []string
	json.Unmarshal(*jsonbField, &fileKeys)

	deleteMap := make(map[string]bool)
	for _, fileKey := range filesToDelete {
		deleteMap[fileKey] = true
	}

	var remainingFiles []string
	for _, file := range fileKeys {
		if !deleteMap[file] {
			remainingFiles = append(remainingFiles, file)
			continue
		}
		// Delete from R2
		_ = s.r2Storage.DeleteFile(file)
	}

	updatedJSON, _ := json.Marshal(remainingFiles)
	*jsonbField = updatedJSON
}

// mapToResponse maps model to DTO response
func (s *KepegawaianServiceImpl) mapToResponse(data *models.Kepegawaian) *dtos.KepegawaianResponse {
	// Map rombel_bidang_studi
	var rombelBidangStudi []uint
	json.Unmarshal(data.RombelBidangStudi, &rombelBidangStudi)

	// Map sertifikat_lainnya
	var sertifikatLainnya []string
	json.Unmarshal(data.SertifikatLainnya, &sertifikatLainnya)

	// Map dokumen_lainnya
	var dokumenLainnya []string
	json.Unmarshal(data.DokumenLainnya, &dokumenLainnya)

	return &dtos.KepegawaianResponse{
		ID:                    data.ID,
		Nama:                  data.Nama,
		Username:              data.Username,
		NIP:                   data.NIP,
		NKKI:                  data.NKKI,
		Foto:                  s.stringOrNil(s.r2Storage.GetPublicURL(data.Foto)),
		Kategori:              data.Kategori,
		Jabatan:               data.Jabatan,
		BidangStudiID:         data.BidangStudiID,
		RombelGuruKelasID:     data.RombelGuruKelasID,
		RombelBidangStudi:     rombelBidangStudi,
		KK:                    s.stringOrNil(s.r2Storage.GetPublicURL(data.KK)),
		AktaLahir:             s.stringOrNil(s.r2Storage.GetPublicURL(data.AktaLahir)),
		KTP:                   s.stringOrNil(s.r2Storage.GetPublicURL(data.KTP)),
		IjazahSD:              s.stringOrNil(s.r2Storage.GetPublicURL(data.IjazahSD)),
		IjazahSMP:             s.stringOrNil(s.r2Storage.GetPublicURL(data.IjazahSMP)),
		IjazahSMA:             s.stringOrNil(s.r2Storage.GetPublicURL(data.IjazahSMA)),
		IjazahS1:              s.stringOrNil(s.r2Storage.GetPublicURL(data.IjazahS1)),
		IjazahS2:              s.stringOrNil(s.r2Storage.GetPublicURL(data.IjazahS2)),
		IjazahS3:              s.stringOrNil(s.r2Storage.GetPublicURL(data.IjazahS3)),
		SertifikatPendidik:    s.stringOrNil(s.r2Storage.GetPublicURL(data.SertifikatPendidik)),
		SertifikatLainnya:     s.mapURLsToPublic(sertifikatLainnya),
		SK:                    s.stringOrNil(s.r2Storage.GetPublicURL(data.SK)),
		DokumenLainnya:        s.mapURLsToPublic(dokumenLainnya),
		Status:                data.Status,
		CreatedAt:             data.CreatedAt,
		UpdatedAt:             data.UpdatedAt,
		CreatedByID:           data.CreatedByID,
		UpdatedByID:           data.UpdatedByID,
	}
}

// Helper function to convert empty string to nil pointer
func (s *KepegawaianServiceImpl) stringOrNil(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

// Helper function to map storage URLs to public URLs
func (s *KepegawaianServiceImpl) mapURLsToPublic(urls []string) []string {
	var publicURLs []string
	for _, url := range urls {
		if url != "" {
			publicURLs = append(publicURLs, s.r2Storage.GetPublicURL(url))
		}
	}
	return publicURLs
}

// Helper function to get document folder path based on document type
func (s *KepegawaianServiceImpl) getDocumentFolderPath(docType string) string {
	switch docType {
	case "kk":
		return "kepegawaian/kk"
	case "akta_lahir":
		return "kepegawaian/akta-lahir"
	case "ktp":
		return "kepegawaian/ktp"
	case "ijazah_sd":
		return "kepegawaian/ijazah-sd"
	case "ijazah_smp":
		return "kepegawaian/ijazah-smp"
	case "ijazah_sma":
		return "kepegawaian/ijazah-sma"
	case "ijazah_s1":
		return "kepegawaian/ijazah-s1"
	case "ijazah_s2":
		return "kepegawaian/ijazah-s2"
	case "ijazah_s3":
		return "kepegawaian/ijazah-s3"
	case "sertifikat_pendidik":
		return "kepegawaian/sertifikat-pendidik"
	case "sertifikat_lainnya":
		return "kepegawaian/sertifikat-lainnya"
	case "sk":
		return "kepegawaian/sk"
	case "dokumen_lainnya":
		return "kepegawaian/dokumen-lainnya"
	default:
		return "kepegawaian"
	}
}

// uploadResult holds upload result for each document type
type uploadResult struct {
	fileKey  string
	fileKeys []string
	err      error
}

// uploadDocumentsParallel uploads multiple documents in parallel (max 5 concurrent)
func (s *KepegawaianServiceImpl) uploadDocumentsParallel(docs map[string][]*multipart.FileHeader) map[string]uploadResult {
	results := make(map[string]uploadResult)
	var wg sync.WaitGroup
	resultsMutex := sync.Mutex{}
	
	// Semaphore to limit concurrent uploads (max 5)
	semaphore := make(chan struct{}, 5)
	
	for docType, files := range docs {
		if len(files) == 0 {
			continue
		}
		
		wg.Add(1)
		go func(docType string, files []*multipart.FileHeader) {
			defer wg.Done()
			
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			folderPath := s.getDocumentFolderPath(docType)
			result := uploadResult{}
			
			// For single file documents
			if docType != "sertifikat_lainnya" && docType != "dokumen_lainnya" {
				if len(files) > 0 && files[0] != nil {
					file := files[0]
					
					// Validate file size
					if file.Size > 10*1024*1024 {
						result.err = fmt.Errorf("%s size must not exceed 10MB", docType)
					} else {
						// Upload file to R2
						fileKey, err := s.r2Storage.UploadFile(file, folderPath)
						if err != nil {
							result.err = err
						} else {
							result.fileKey = fileKey
						}
					}
				}
			} else {
				// For multiple file documents
				var fileKeys []string
				for _, file := range files {
					if file == nil {
						continue
					}
					
					if file.Size > 10*1024*1024 {
						result.err = fmt.Errorf("%s size must not exceed 10MB per file", docType)
						break
					}
					
					fileKey, err := s.r2Storage.UploadFile(file, folderPath)
					if err != nil {
						result.err = err
						break
					}
					fileKeys = append(fileKeys, fileKey)
				}
				result.fileKeys = fileKeys
			}
			
			// Store result safely
			resultsMutex.Lock()
			results[docType] = result
			resultsMutex.Unlock()
		}(docType, files)
	}
	
	wg.Wait()
	return results
}
