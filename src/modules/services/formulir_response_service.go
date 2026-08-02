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

type FormulirResponseService interface {
	SubmitAuthenticated(req *dtos.FormulirSubmitRequest, userID uint, ipAddress, userAgent string, files map[uint][]*multipart.FileHeader) (*dtos.FormulirSubmitResponse, error)
	SubmitPublic(req *dtos.FormulirSubmitRequest, ipAddress, userAgent string, files map[uint][]*multipart.FileHeader) (*dtos.FormulirSubmitResponse, error)
	EditResponse(req *dtos.FormulirEditResponseRequest, userID uint, ipAddress, userAgent string, files map[uint][]*multipart.FileHeader) (*dtos.FormulirSubmitResponse, error)
	GetResponsesBySlug(slug string, role *string, userID uint) (*dtos.FormulirResponsesDetailResponse, error)
	GetResponseByUser(req *dtos.FormulirResponseByUserRequest, userID uint) (*dtos.FormulirResponsesDetailResponse, error)
	DeleteResponse(req *dtos.FormulirDeleteResponseRequest, userID uint) error
	ResetResponses(formulirID uint, role *string, userID uint) error
	GetStatisticsBySlug(slug string, role *string, userID uint) (*dtos.FormulirStatisticResponse, error)
}

type FormulirResponseServiceImpl struct {
	formulirRepo repositories.FormulirRepository
	responseRepo repositories.FormulirResponseRepository
	r2Storage    *utils.R2Storage
}

// NewFormulirResponseService creates a new service
func NewFormulirResponseService(
	formulirRepo repositories.FormulirRepository,
	responseRepo repositories.FormulirResponseRepository,
	r2Storage *utils.R2Storage,
) FormulirResponseService {
	return &FormulirResponseServiceImpl{
		formulirRepo: formulirRepo,
		responseRepo: responseRepo,
		r2Storage:    r2Storage,
	}
}

// getWIBLocation returns WIB timezone location
func getWIBLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback to fixed offset if timezone data not available (e.g., in Docker)
		return time.FixedZone("WIB", 7*60*60) // UTC+7
	}
	return loc
}

// SubmitAuthenticated handles authenticated form submission with validations and file uploads
func (s *FormulirResponseServiceImpl) SubmitAuthenticated(req *dtos.FormulirSubmitRequest, userID uint, ipAddress, userAgent string, files map[uint][]*multipart.FileHeader) (*dtos.FormulirSubmitResponse, error) {
	// Get formulir details
	formulir, err := s.formulirRepo.GetByID(req.FormulirID)
	if err != nil {
		return nil, errors.New("formulir not found")
	}

	// Validation 1: Check if form is active
	if !formulir.IsActive {
		return nil, errors.New("formulir is not active")
	}

	// Get current time in WIB (Asia/Jakarta, UTC+7)
	wibLocation := getWIBLocation()
	now := time.Now().In(wibLocation)

	// Validation 2: Check start date
	if formulir.StartDate != nil {
		// Treat start_date as WIB time (since user inputs WIB time)
		startDateWIB := time.Date(
			formulir.StartDate.Year(), formulir.StartDate.Month(), formulir.StartDate.Day(),
			formulir.StartDate.Hour(), formulir.StartDate.Minute(), formulir.StartDate.Second(),
			formulir.StartDate.Nanosecond(), wibLocation,
		)
		
		if now.Before(startDateWIB) {
			return nil, fmt.Errorf("formulir belum dibuka. Mulai: %s WIB", startDateWIB.Format("2006-01-02 15:04:05"))
		}
	}

	// Validation 3: Check end date
	if formulir.EndDate != nil {
		// Treat end_date as WIB time (since user inputs WIB time)
		endDateWIB := time.Date(
			formulir.EndDate.Year(), formulir.EndDate.Month(), formulir.EndDate.Day(),
			formulir.EndDate.Hour(), formulir.EndDate.Minute(), formulir.EndDate.Second(),
			formulir.EndDate.Nanosecond(), wibLocation,
		)
		
		if now.After(endDateWIB) {
			return nil, fmt.Errorf("formulir sudah ditutup. Berakhir: %s WIB", endDateWIB.Format("2006-01-02 15:04:05"))
		}
	}

	// Validation 4: Check max responses
	if formulir.MaxResponses != nil {
		currentCount, err := s.responseRepo.CountResponsesByFormulirID(req.FormulirID)
		if err != nil {
			return nil, err
		}
		if currentCount >= int64(*formulir.MaxResponses) {
			return nil, fmt.Errorf("formulir sudah mencapai batas maksimal respons (%d)", *formulir.MaxResponses)
		}
	}

	// Validation 5: Check if user already submitted (if allow_multiple_responses = false)
	if !formulir.AllowMultipleResponses {
		existingResponse, _ := s.responseRepo.GetByFormulirIDAndUserID(req.FormulirID, userID)
		if existingResponse != nil {
			return nil, errors.New("Anda sudah pernah mengisi formulir ini")
		}
	}

	// Validation 6: Check all required questions are answered
	requiredQuestions := make(map[uint]bool)
	questionDetails := make(map[uint]*models.FormulirPertanyaan) // Store question details
	for _, p := range formulir.Pertanyaan {
		questionDetails[p.ID] = &p
		if p.IsRequired {
			requiredQuestions[p.ID] = false
		}
	}

	answeredQuestions := make(map[uint]bool)
	for _, j := range req.Jawaban {
		answeredQuestions[j.PertanyaanID] = true
		// Mark required questions as answered
		if _, exists := requiredQuestions[j.PertanyaanID]; exists {
			requiredQuestions[j.PertanyaanID] = true
		}
	}

	// Check if all required questions are answered
	for qID, answered := range requiredQuestions {
		if !answered {
			question := questionDetails[qID]
			return nil, fmt.Errorf("pertanyaan '%s' wajib diisi", question.Label)
		}
	}

	// Validation 7: Validate question IDs exist in form
	validQuestionIDs := make(map[uint]bool)
	for _, p := range formulir.Pertanyaan {
		validQuestionIDs[p.ID] = true
	}

	for _, j := range req.Jawaban {
		if !validQuestionIDs[j.PertanyaanID] {
			return nil, fmt.Errorf("pertanyaan ID %d tidak ditemukan dalam formulir ini", j.PertanyaanID)
		}
	}

	// Start transaction
	tx := s.responseRepo.BeginTransaction()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Track uploaded files for cleanup on error
	uploadedFiles := []string{}

	// Create response record with user ID
	response := &models.FormulirResponse{
		FormulirID:        req.FormulirID,
		SubmittedByUserID: &userID,
		SubmittedAsRole:   req.SubmittedAsRole,
		IPAddress:         &ipAddress,
		UserAgent:         &userAgent,
		SubmittedAt:       now,
	}

	if err := tx.Create(response).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Handle file uploads for file-type questions
	fileURLMap := make(map[uint]string) // pertanyaan_id -> file_url
	
	for pertanyaanID, fileHeaders := range files {
		if len(fileHeaders) > 0 {
			fileHeader := fileHeaders[0] // Take first file

			// Validate file size (10MB max)
			if fileHeader.Size > 10*1024*1024 {
				tx.Rollback()
				// Cleanup uploaded files
				for _, url := range uploadedFiles {
					s.r2Storage.DeleteFile(url)
				}
				return nil, fmt.Errorf("file untuk pertanyaan ID %d terlalu besar (max 10MB)", pertanyaanID)
			}

			// Upload to R2: formulir/formulir-{id}/formulir-response/response-{response_id}-{filename}
			folder := fmt.Sprintf("formulir/formulir-%d/formulir-response", req.FormulirID)
			fileURL, err := s.r2Storage.UploadFile(fileHeader, folder)
			if err != nil {
				tx.Rollback()
				// Cleanup uploaded files
				for _, url := range uploadedFiles {
					s.r2Storage.DeleteFile(url)
				}
				return nil, fmt.Errorf("gagal upload file untuk pertanyaan ID %d: %v", pertanyaanID, err)
			}

			uploadedFiles = append(uploadedFiles, fileURL)
			fileURLMap[pertanyaanID] = fileURL
		}
	}

	// Create jawaban records
	for _, j := range req.Jawaban {
		var jawabanJSON []byte
		if j.JawabanJSON != nil {
			jawabanJSON, _ = json.Marshal(j.JawabanJSON)
		}

		// If this is a file question and file was uploaded, use file URL
		jawabanText := j.JawabanText
		if fileURL, exists := fileURLMap[j.PertanyaanID]; exists {
			jawabanText = &fileURL
		}

		jawaban := &models.FormulirResponseJawaban{
			ResponseID:   response.ID,
			PertanyaanID: j.PertanyaanID,
			JawabanText:  jawabanText,
			JawabanJSON:  jawabanJSON,
		}

		if err := tx.Create(jawaban).Error; err != nil {
			tx.Rollback()
			// Cleanup uploaded files
			for _, url := range uploadedFiles {
				s.r2Storage.DeleteFile(url)
			}
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		// Cleanup uploaded files
		for _, url := range uploadedFiles {
			s.r2Storage.DeleteFile(url)
		}
		return nil, err
	}

	return &dtos.FormulirSubmitResponse{
		ID:          response.ID,
		FormulirID:  response.FormulirID,
		SubmittedAt: response.SubmittedAt,
		Message:     "Formulir berhasil disubmit",
	}, nil
}

// SubmitPublic handles public form submission (no auth required) with file uploads
func (s *FormulirResponseServiceImpl) SubmitPublic(req *dtos.FormulirSubmitRequest, ipAddress, userAgent string, files map[uint][]*multipart.FileHeader) (*dtos.FormulirSubmitResponse, error) {
	// Get formulir details
	formulir, err := s.formulirRepo.GetByID(req.FormulirID)
	if err != nil {
		return nil, errors.New("formulir not found")
	}

	// Validation 1: Check if form is active
	if !formulir.IsActive {
		return nil, errors.New("formulir is not active")
	}

	// Validation 2: Check if form is public
	if formulir.AccessType != "public" {
		return nil, errors.New("formulir ini memerlukan login")
	}

	// Get current time in WIB (Asia/Jakarta, UTC+7)
	wibLocation := getWIBLocation()
	now := time.Now().In(wibLocation)

	// Validation 3: Check start date
	if formulir.StartDate != nil {
		startDateWIB := time.Date(
			formulir.StartDate.Year(), formulir.StartDate.Month(), formulir.StartDate.Day(),
			formulir.StartDate.Hour(), formulir.StartDate.Minute(), formulir.StartDate.Second(),
			formulir.StartDate.Nanosecond(), wibLocation,
		)
		
		if now.Before(startDateWIB) {
			return nil, fmt.Errorf("formulir belum dibuka. Mulai: %s WIB", startDateWIB.Format("2006-01-02 15:04:05"))
		}
	}

	// Validation 4: Check end date
	if formulir.EndDate != nil {
		endDateWIB := time.Date(
			formulir.EndDate.Year(), formulir.EndDate.Month(), formulir.EndDate.Day(),
			formulir.EndDate.Hour(), formulir.EndDate.Minute(), formulir.EndDate.Second(),
			formulir.EndDate.Nanosecond(), wibLocation,
		)
		
		if now.After(endDateWIB) {
			return nil, fmt.Errorf("formulir sudah ditutup. Berakhir: %s WIB", endDateWIB.Format("2006-01-02 15:04:05"))
		}
	}

	// Validation 5: Check max responses
	if formulir.MaxResponses != nil {
		currentCount, err := s.responseRepo.CountResponsesByFormulirID(req.FormulirID)
		if err != nil {
			return nil, err
		}
		if currentCount >= int64(*formulir.MaxResponses) {
			return nil, fmt.Errorf("formulir sudah mencapai batas maksimal respons (%d)", *formulir.MaxResponses)
		}
	}

	// Validation 6: Check if IP already submitted (only if allow_multiple_responses = false)
	if !formulir.AllowMultipleResponses {
		existingResponse, _ := s.responseRepo.GetByFormulirIDAndIPAddress(req.FormulirID, ipAddress)
		if existingResponse != nil {
			return nil, errors.New("IP address Anda sudah pernah mengisi formulir ini")
		}
	}

	// Validation 7: Check all required questions are answered
	requiredQuestions := make(map[uint]bool)
	questionDetails := make(map[uint]*models.FormulirPertanyaan) // Store question details
	for _, p := range formulir.Pertanyaan {
		questionDetails[p.ID] = &p
		if p.IsRequired {
			requiredQuestions[p.ID] = false
		}
	}

	for _, j := range req.Jawaban {
		if _, exists := requiredQuestions[j.PertanyaanID]; exists {
			requiredQuestions[j.PertanyaanID] = true
		}
	}

	for qID, answered := range requiredQuestions {
		if !answered {
			question := questionDetails[qID]
			return nil, fmt.Errorf("pertanyaan '%s' wajib diisi", question.Label)
		}
	}

	// Validation 8: Validate question IDs exist in form
	validQuestionIDs := make(map[uint]bool)
	for _, p := range formulir.Pertanyaan {
		validQuestionIDs[p.ID] = true
	}

	for _, j := range req.Jawaban {
		if !validQuestionIDs[j.PertanyaanID] {
			return nil, fmt.Errorf("pertanyaan ID %d tidak ditemukan dalam formulir ini", j.PertanyaanID)
		}
	}

	// Start transaction
	tx := s.responseRepo.BeginTransaction()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Track uploaded files for cleanup on error
	uploadedFiles := []string{}

	// Create response record WITHOUT user ID (public submission)
	response := &models.FormulirResponse{
		FormulirID:        req.FormulirID,
		SubmittedByUserID: nil, // Public form, no user ID
		IPAddress:         &ipAddress,
		UserAgent:         &userAgent,
		SubmittedAt:       now,
	}

	if err := tx.Create(response).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Handle file uploads for file-type questions
	fileURLMap := make(map[uint]string) // pertanyaan_id -> file_url
	
	for pertanyaanID, fileHeaders := range files {
		if len(fileHeaders) > 0 {
			fileHeader := fileHeaders[0] // Take first file

			// Validate file size (10MB max)
			if fileHeader.Size > 10*1024*1024 {
				tx.Rollback()
				// Cleanup uploaded files
				for _, url := range uploadedFiles {
					s.r2Storage.DeleteFile(url)
				}
				return nil, fmt.Errorf("file untuk pertanyaan ID %d terlalu besar (max 10MB)", pertanyaanID)
			}

			// Upload to R2: formulir/formulir-{id}/formulir-response/response-{response_id}-{filename}
			folder := fmt.Sprintf("formulir/formulir-%d/formulir-response", req.FormulirID)
			fileURL, err := s.r2Storage.UploadFile(fileHeader, folder)
			if err != nil {
				tx.Rollback()
				// Cleanup uploaded files
				for _, url := range uploadedFiles {
					s.r2Storage.DeleteFile(url)
				}
				return nil, fmt.Errorf("gagal upload file untuk pertanyaan ID %d: %v", pertanyaanID, err)
			}

			uploadedFiles = append(uploadedFiles, fileURL)
			fileURLMap[pertanyaanID] = fileURL
		}
	}

	// Create jawaban records
	for _, j := range req.Jawaban {
		var jawabanJSON []byte
		if j.JawabanJSON != nil {
			jawabanJSON, _ = json.Marshal(j.JawabanJSON)
		}

		// If this is a file question and file was uploaded, use file URL
		jawabanText := j.JawabanText
		if fileURL, exists := fileURLMap[j.PertanyaanID]; exists {
			jawabanText = &fileURL
		}

		jawaban := &models.FormulirResponseJawaban{
			ResponseID:   response.ID,
			PertanyaanID: j.PertanyaanID,
			JawabanText:  jawabanText,
			JawabanJSON:  jawabanJSON,
		}

		if err := tx.Create(jawaban).Error; err != nil {
			tx.Rollback()
			// Cleanup uploaded files
			for _, url := range uploadedFiles {
				s.r2Storage.DeleteFile(url)
			}
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		// Cleanup uploaded files
		for _, url := range uploadedFiles {
			s.r2Storage.DeleteFile(url)
		}
		return nil, err
	}

	return &dtos.FormulirSubmitResponse{
		ID:          response.ID,
		FormulirID:  response.FormulirID,
		SubmittedAt: response.SubmittedAt,
		Message:     "Formulir berhasil disubmit",
	}, nil
}

// EditResponse handles editing existing form response with validations and file uploads
func (s *FormulirResponseServiceImpl) EditResponse(req *dtos.FormulirEditResponseRequest, userID uint, ipAddress, userAgent string, files map[uint][]*multipart.FileHeader) (*dtos.FormulirSubmitResponse, error) {
	// Get existing response
	existingResponse, err := s.responseRepo.GetResponseByID(req.ResponseID)
	if err != nil {
		return nil, errors.New("response not found")
	}

	// AUTHORIZATION CHECK with role-based logic
	authorized := false
	authErrorMsg := "unauthorized: you cannot edit this response"

	// Rule 1: Admin role can edit anything
	if req.Role != nil && *req.Role == "admin" {
		authorized = true
	} else if req.Role != nil {
		// Rule 2: Role-based authorization for non-admin
		requestRole := *req.Role
		
		// Check if submitted_as_role exists and matches request role
		if existingResponse.SubmittedAsRole == nil {
			return nil, errors.New("this response has no submitted_as_role, cannot use role-based authorization")
		}
		
		submittedAsRole := *existingResponse.SubmittedAsRole
		if requestRole != submittedAsRole {
			return nil, fmt.Errorf("unauthorized: role mismatch. Response was submitted as '%s' but you're trying to edit as '%s'", submittedAsRole, requestRole)
		}

		// Check if submitted_by_user_id matches current user
		if existingResponse.SubmittedByUserID == nil {
			return nil, errors.New("cannot edit public response with role-based authorization")
		}
		
		if *existingResponse.SubmittedByUserID != userID {
			return nil, errors.New("unauthorized: this response was not submitted by your user account")
		}

		// Get user details to check username
		if existingResponse.SubmittedBy == nil {
			return nil, errors.New("cannot verify user details for authorization")
		}

		username := existingResponse.SubmittedBy.Username

		if requestRole == "pendidik" || requestRole == "tendik" {
			// Check in kepegawaian table
			kepegawaian, err := s.responseRepo.GetKepegawaianByUsername(username)
			if err == nil && kepegawaian != nil {
				// Username matches, user is authorized
				authorized = true
			} else {
				authErrorMsg = fmt.Sprintf("unauthorized: user not found in kepegawaian table for role %s", requestRole)
			}
		} else if requestRole == "murid" {
			// Check in peserta_didik table
			pesertaDidik, err := s.responseRepo.GetPesertaDidikByUsername(username)
			if err == nil && pesertaDidik != nil {
				// Username matches, user is authorized
				authorized = true
			} else {
				authErrorMsg = "unauthorized: user not found in peserta_didik table for role murid"
			}
		} else {
			// Other roles just check user_id match (orang_tua, etc.)
			authorized = true
		}
	} else {
		// Rule 3: Default - Only the user who submitted can edit (check by submitted_by_user_id)
		if existingResponse.SubmittedByUserID == nil {
			return nil, errors.New("cannot edit public response without specifying role")
		}
		if *existingResponse.SubmittedByUserID == userID {
			authorized = true
		} else {
			authErrorMsg = "unauthorized: you can only edit your own response"
		}
	}

	if !authorized {
		return nil, errors.New(authErrorMsg)
	}

	// Get formulir details
	formulir, err := s.formulirRepo.GetByID(existingResponse.FormulirID)
	if err != nil {
		return nil, errors.New("formulir not found")
	}

	// Validation 1: Check if form is active
	if !formulir.IsActive {
		return nil, errors.New("formulir is not active")
	}

	// Get current time in WIB
	wibLocation := getWIBLocation()
	now := time.Now().In(wibLocation)

	// Validation 2: Check if still within edit period (before end_date)
	if formulir.EndDate != nil {
		endDateWIB := time.Date(
			formulir.EndDate.Year(), formulir.EndDate.Month(), formulir.EndDate.Day(),
			formulir.EndDate.Hour(), formulir.EndDate.Minute(), formulir.EndDate.Second(),
			formulir.EndDate.Nanosecond(), wibLocation,
		)
		
		if now.After(endDateWIB) {
			return nil, fmt.Errorf("formulir sudah ditutup. Tidak dapat mengedit response setelah: %s WIB", endDateWIB.Format("2006-01-02 15:04:05"))
		}
	}

	// Build question details map
	questionDetails := make(map[uint]*models.FormulirPertanyaan)
	requiredQuestions := make(map[uint]bool)
	for i := range formulir.Pertanyaan {
		p := &formulir.Pertanyaan[i]
		questionDetails[p.ID] = p
		if p.IsRequired {
			requiredQuestions[p.ID] = false
		}
	}

	// Get existing jawaban to check which answers already exist
	existingJawabanCheck, err := s.responseRepo.GetJawabanByResponseID(req.ResponseID)
	if err != nil {
		return nil, err
	}
	existingAnswersByPertanyaan := make(map[uint]bool)
	for _, ej := range existingJawabanCheck {
		existingAnswersByPertanyaan[ej.PertanyaanID] = true
	}

	// Validation 3: Check all required questions are answered
	// For EDIT: Only validate NEW answers (without ID), skip validation for UPDATE (with ID)
	for _, j := range req.Jawaban {
		if _, exists := requiredQuestions[j.PertanyaanID]; exists {
			// If jawaban has ID (UPDATE existing), mark as answered
			if j.ID != nil && *j.ID > 0 {
				requiredQuestions[j.PertanyaanID] = true
			} else {
				// New answer (CREATE), check if it has value
				if j.JawabanText != nil && *j.JawabanText != "" {
					requiredQuestions[j.PertanyaanID] = true
				} else if j.JawabanJSON != nil {
					requiredQuestions[j.PertanyaanID] = true
				}
			}
		}
	}

	// Also check if files are being uploaded for required file questions
	for pertanyaanID := range files {
		if _, exists := requiredQuestions[pertanyaanID]; exists {
			requiredQuestions[pertanyaanID] = true
		}
	}

	// Check for required questions that are not in request but exist in old response
	for qID, answered := range requiredQuestions {
		if !answered {
			// Check if this question was answered before (exists in old response)
			if existingAnswersByPertanyaan[qID] {
				// Check if it's being deleted (not in new request)
				inNewRequest := false
				for _, j := range req.Jawaban {
					if j.PertanyaanID == qID {
						inNewRequest = true
						break
					}
				}
				// If not in new request, it will be deleted, so validate
				if !inNewRequest {
					question := questionDetails[qID]
					return nil, fmt.Errorf("pertanyaan '%s' wajib diisi", question.Label)
				}
			} else {
				// Question was never answered and is not in new request
				question := questionDetails[qID]
				return nil, fmt.Errorf("pertanyaan '%s' wajib diisi", question.Label)
			}
		}
	}

	// Validation 4: Validate question IDs exist in form
	validQuestionIDs := make(map[uint]bool)
	for _, p := range formulir.Pertanyaan {
		validQuestionIDs[p.ID] = true
	}

	for _, j := range req.Jawaban {
		if !validQuestionIDs[j.PertanyaanID] {
			return nil, fmt.Errorf("pertanyaan ID %d tidak ditemukan dalam formulir ini", j.PertanyaanID)
		}
	}

	// Start transaction
	tx := s.responseRepo.BeginTransaction()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Track uploaded files for cleanup on error
	uploadedFiles := []string{}
	deletedFiles := []string{}

	// Get existing jawaban
	existingJawaban, err := s.responseRepo.GetJawabanByResponseID(req.ResponseID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Map existing jawaban by pertanyaan_id and by ID
	existingJawabanMap := make(map[uint]*models.FormulirResponseJawaban)
	existingJawabanByID := make(map[uint]*models.FormulirResponseJawaban)
	for i := range existingJawaban {
		existingJawabanMap[existingJawaban[i].PertanyaanID] = &existingJawaban[i]
		existingJawabanByID[existingJawaban[i].ID] = &existingJawaban[i]
	}

	// Handle file deletions from files_to_delete array
	filesToDeleteMap := make(map[uint]bool)
	for _, pertanyaanID := range req.FilesToDelete {
		filesToDeleteMap[pertanyaanID] = true
		
		// Delete old file from R2 if exists
		if oldJawaban, exists := existingJawabanMap[pertanyaanID]; exists {
			if oldJawaban.JawabanText != nil && *oldJawaban.JawabanText != "" {
				err := s.r2Storage.DeleteFile(*oldJawaban.JawabanText)
				if err == nil {
					deletedFiles = append(deletedFiles, *oldJawaban.JawabanText)
				}
			}
		}
	}

	// Handle file uploads for file-type questions
	fileURLMap := make(map[uint]string) // pertanyaan_id -> file_url
	
	for pertanyaanID, fileHeaders := range files {
		if len(fileHeaders) > 0 {
			fileHeader := fileHeaders[0] // Take first file

			// Validate file size (10MB max)
			if fileHeader.Size > 10*1024*1024 {
				tx.Rollback()
				// Cleanup uploaded files
				for _, url := range uploadedFiles {
					s.r2Storage.DeleteFile(url)
				}
				return nil, fmt.Errorf("file untuk pertanyaan ID %d terlalu besar (max 10MB)", pertanyaanID)
			}

			// Delete old file from R2 if exists
			if oldJawaban, exists := existingJawabanMap[pertanyaanID]; exists {
				if oldJawaban.JawabanText != nil && *oldJawaban.JawabanText != "" {
					err := s.r2Storage.DeleteFile(*oldJawaban.JawabanText)
					if err == nil {
						deletedFiles = append(deletedFiles, *oldJawaban.JawabanText)
					}
				}
			}

			// Upload new file to R2
			folder := fmt.Sprintf("formulir/formulir-%d/formulir-response", existingResponse.FormulirID)
			fileURL, err := s.r2Storage.UploadFile(fileHeader, folder)
			if err != nil {
				tx.Rollback()
				// Cleanup uploaded files
				for _, url := range uploadedFiles {
					s.r2Storage.DeleteFile(url)
				}
				return nil, fmt.Errorf("gagal upload file untuk pertanyaan ID %d: %v", pertanyaanID, err)
			}

			uploadedFiles = append(uploadedFiles, fileURL)
			fileURLMap[pertanyaanID] = fileURL
		}
	}

	// Map new jawaban by pertanyaan_id and by ID
	newJawabanMap := make(map[uint]*dtos.FormulirJawabanEditRequest)
	newJawabanByID := make(map[uint]*dtos.FormulirJawabanEditRequest)
	for i := range req.Jawaban {
		newJawabanMap[req.Jawaban[i].PertanyaanID] = &req.Jawaban[i]
		if req.Jawaban[i].ID != nil {
			newJawabanByID[*req.Jawaban[i].ID] = &req.Jawaban[i]
		}
	}

	// Update or create jawaban
	for _, newJawaban := range req.Jawaban {
		var jawabanJSON []byte
		if newJawaban.JawabanJSON != nil {
			jawabanJSON, _ = json.Marshal(newJawaban.JawabanJSON)
		}

		// Determine jawaban_text
		jawabanText := newJawaban.JawabanText
		
		// If file was uploaded for this question, use file URL
		if fileURL, exists := fileURLMap[newJawaban.PertanyaanID]; exists {
			jawabanText = &fileURL
		}
		
		// If file was deleted for this question, set to empty
		if filesToDeleteMap[newJawaban.PertanyaanID] {
			emptyStr := ""
			jawabanText = &emptyStr
		}

		// Check if jawaban has ID (UPDATE) or no ID (CREATE)
		if newJawaban.ID != nil && *newJawaban.ID > 0 {
			// UPDATE existing jawaban by ID
			if existingJaw, exists := existingJawabanByID[*newJawaban.ID]; exists {
				existingJaw.JawabanText = jawabanText
				existingJaw.JawabanJSON = jawabanJSON

				if err := tx.Save(existingJaw).Error; err != nil {
					tx.Rollback()
					// Cleanup uploaded files
					for _, url := range uploadedFiles {
						s.r2Storage.DeleteFile(url)
					}
					return nil, err
				}
			} else {
				tx.Rollback()
				// Cleanup uploaded files
				for _, url := range uploadedFiles {
					s.r2Storage.DeleteFile(url)
				}
				return nil, fmt.Errorf("jawaban ID %d not found in this response", *newJawaban.ID)
			}
		} else {
			// CREATE new jawaban
			newJaw := &models.FormulirResponseJawaban{
				ResponseID:   req.ResponseID,
				PertanyaanID: newJawaban.PertanyaanID,
				JawabanText:  jawabanText,
				JawabanJSON:  jawabanJSON,
			}

			if err := tx.Create(newJaw).Error; err != nil {
				tx.Rollback()
				// Cleanup uploaded files
				for _, url := range uploadedFiles {
					s.r2Storage.DeleteFile(url)
				}
				return nil, err
			}
		}
	}

	// Delete jawaban that are not in the new request (based on ID)
	for jawabanID, existingJaw := range existingJawabanByID {
		if _, exists := newJawabanByID[jawabanID]; !exists {
			// This jawaban ID is not in request, delete it
			// Delete file from R2 if exists
			if existingJaw.JawabanText != nil && *existingJaw.JawabanText != "" {
				s.r2Storage.DeleteFile(*existingJaw.JawabanText)
			}

			// Delete jawaban record
			if err := tx.Delete(existingJaw).Error; err != nil {
				tx.Rollback()
				// Cleanup uploaded files
				for _, url := range uploadedFiles {
					s.r2Storage.DeleteFile(url)
				}
				return nil, err
			}
		}
	}

	// Update response metadata (do NOT update submitted_as_role, keep original)
	existingResponse.IPAddress = &ipAddress
	existingResponse.UserAgent = &userAgent
	existingResponse.SubmittedAt = now

	if err := tx.Save(existingResponse).Error; err != nil {
		tx.Rollback()
		// Cleanup uploaded files
		for _, url := range uploadedFiles {
			s.r2Storage.DeleteFile(url)
		}
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		// Cleanup uploaded files
		for _, url := range uploadedFiles {
			s.r2Storage.DeleteFile(url)
		}
		return nil, err
	}

	return &dtos.FormulirSubmitResponse{
		ID:          existingResponse.ID,
		FormulirID:  existingResponse.FormulirID,
		SubmittedAt: existingResponse.SubmittedAt,
		Message:     "Response berhasil diupdate",
	}, nil
}

// GetResponsesBySlug gets all responses for a form by slug (Google Forms style)
func (s *FormulirResponseServiceImpl) GetResponsesBySlug(slug string, role *string, userID uint) (*dtos.FormulirResponsesDetailResponse, error) {
	// Get formulir by slug
	formulir, err := s.formulirRepo.GetBySlug(slug)
	if err != nil {
		return nil, errors.New("formulir not found")
	}

	// If role is provided, validate form ownership (except for admin)
	if role != nil && *role != "" {
		requestRole := *role
		
		// Admin can access any form
		if requestRole != "admin" {
			// Validate: form created_by_role must match request role
			if formulir.CreatedByRole == nil {
				return nil, errors.New("form was not created with a role (created_by_role is NULL)")
			}
			
			if *formulir.CreatedByRole != requestRole {
				return nil, fmt.Errorf("form was created by role '%s', not '%s'", *formulir.CreatedByRole, requestRole)
			}
			
			// Validate: form created_by_user_id must match current user
			if formulir.CreatedByUserID != userID {
				return nil, errors.New("you are not the owner of this form")
			}
			
			// Validate: user exists in respective table
			user, err := s.responseRepo.GetUserByID(userID)
			if err != nil {
				return nil, errors.New("user not found")
			}
			
			username := user.Username
			
			if requestRole == "pendidik" || requestRole == "tendik" {
				// Check in kepegawaian table
				kepegawaian, err := s.responseRepo.GetKepegawaianByUsername(username)
				if err != nil || kepegawaian == nil {
					return nil, fmt.Errorf("user not found in kepegawaian table for role %s", requestRole)
				}
			} else if requestRole == "murid" {
				// Check in peserta_didik table
				pesertaDidik, err := s.responseRepo.GetPesertaDidikByUsername(username)
				if err != nil || pesertaDidik == nil {
					return nil, errors.New("user not found in peserta_didik table for role murid")
				}
			}
			// Other roles (orang_tua, etc.) don't need special validation
		}
		// If role == "admin", skip all validation checks
	}

	// Get all responses for this formulir
	responses, err := s.responseRepo.GetAllResponsesByFormulirID(formulir.ID)
	if err != nil {
		return nil, err
	}

	// Count total responses
	totalResponses := int64(len(responses))

	// Map formulir to response DTO
	var targetUserTypes []string
	if formulir.TargetUserTypes != nil {
		_ = json.Unmarshal(formulir.TargetUserTypes, &targetUserTypes)
	}

	// Map pertanyaan
	pertanyaanMap := make(map[uint]*models.FormulirPertanyaan)
	pertanyaanResponses := make([]dtos.FormulirPertanyaanResponse, 0)
	
	for i := range formulir.Pertanyaan {
		p := &formulir.Pertanyaan[i]
		pertanyaanMap[p.ID] = p

		// Parse options
		var options []string
		if p.Options != nil {
			_ = json.Unmarshal(p.Options, &options)
		}

		// Parse validation_rules
		var validationRules map[string]interface{}
		if p.ValidationRules != nil {
			_ = json.Unmarshal(p.ValidationRules, &validationRules)
		}

		// Parse file_config
		var fileConfig map[string]interface{}
		if p.FileConfig != nil {
			_ = json.Unmarshal(p.FileConfig, &fileConfig)
		}

		placeholder := ""
		if p.Placeholder != nil {
			placeholder = *p.Placeholder
		}

		dokumen := ""
		if p.Dokumen != nil && *p.Dokumen != "" {
			dokumen = s.r2Storage.GetPublicURL(*p.Dokumen)
		}

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

	formulirResponse := dtos.FormulirResponse{
		ID:                     formulir.ID,
		Judul:                  formulir.Judul,
		Slug:                   formulir.Slug,
		Deskripsi:              formulir.Deskripsi,
		CreatedByUserID:        formulir.CreatedByUserID,
		IsActive:               formulir.IsActive,
		MaxResponses:           formulir.MaxResponses,
		StartDate:              formulir.StartDate,
		EndDate:                formulir.EndDate,
		AccessType:             formulir.AccessType,
		TargetUserTypes:        targetUserTypes,
		AllowMultipleResponses: formulir.AllowMultipleResponses,
		CreatedAt:              formulir.CreatedAt,
		UpdatedAt:              formulir.UpdatedAt,
		Pertanyaan:             pertanyaanResponses,
		PublicURL:              generatePublicURL(formulir.Slug, formulir.AccessType),
	}

	// Map responses to row details
	responseRows := make([]dtos.FormulirResponseRowDetail, 0)
	
	for _, resp := range responses {
		// Map user basic info with role-based name fetching
		var submittedBy *dtos.UserBasic
		if resp.SubmittedBy != nil {
			fullName := resp.SubmittedBy.Nama // Default fallback
			
			// Fetch name based on submitted_as_role
			if resp.SubmittedAsRole != nil {
				role := *resp.SubmittedAsRole
				
				// For pendidik or tendik, fetch from kepegawaian table
				if role == "pendidik" || role == "tendik" {
					kepegawaian, err := s.responseRepo.GetKepegawaianByUsername(resp.SubmittedBy.Username)
					if err == nil && kepegawaian != nil {
						fullName = kepegawaian.Nama
					} else {
						// Log: kepegawaian not found for this username, using fallback from users table
						fmt.Printf("Warning: User '%s' with role '%s' not found in kepegawaian table. Using fallback name from users table.\n", resp.SubmittedBy.Username, role)
					}
				}
				
				// For murid, fetch from peserta_didik table
				if role == "murid" {
					pesertaDidik, err := s.responseRepo.GetPesertaDidikByUsername(resp.SubmittedBy.Username)
					if err == nil && pesertaDidik != nil {
						fullName = pesertaDidik.Nama
					} else {
						// Log: peserta_didik not found for this username, using fallback from users table
						fmt.Printf("Warning: User '%s' with role 'murid' not found in peserta_didik table. Using fallback name from users table.\n", resp.SubmittedBy.Username)
					}
				}
				
				// For orang_tua and admin, use nama from users table (already set as default)
			}
			
			submittedBy = &dtos.UserBasic{
				ID:       resp.SubmittedBy.ID,
				Username: resp.SubmittedBy.Username,
				FullName: fullName,
			}
		}

		// Map jawaban with question labels and answer IDs
		jawaban := make([]dtos.FormulirResponseAnswerDetail, 0)
		for _, j := range resp.Jawaban {
			var jawabanJSON interface{}
			if len(j.JawabanJSON) > 0 {
				_ = json.Unmarshal(j.JawabanJSON, &jawabanJSON)
			}

			// Get question label
			label := ""
			if pertanyaan, exists := pertanyaanMap[j.PertanyaanID]; exists {
				label = pertanyaan.Label
			}

			jawaban = append(jawaban, dtos.FormulirResponseAnswerDetail{
				ID:           j.ID, // Include answer ID for edit operations
				PertanyaanID: j.PertanyaanID,
				Label:        label,
				JawabanText:  j.JawabanText,
				JawabanJSON:  jawabanJSON,
			})
		}

		responseRows = append(responseRows, dtos.FormulirResponseRowDetail{
			ResponseID:        resp.ID,
			SubmittedByUserID: resp.SubmittedByUserID,
			SubmittedAsRole:   resp.SubmittedAsRole,
			SubmittedBy:       submittedBy,
			IPAddress:         resp.IPAddress,
			UserAgent:         resp.UserAgent,
			SubmittedAt:       resp.SubmittedAt,
			Jawaban:           jawaban,
		})
	}

	return &dtos.FormulirResponsesDetailResponse{
		Formulir:       formulirResponse,
		Responses:      responseRows,
		TotalResponses: totalResponses,
	}, nil
}

// GetResponseByUser gets a specific user's response for a form by slug
func (s *FormulirResponseServiceImpl) GetResponseByUser(req *dtos.FormulirResponseByUserRequest, userID uint) (*dtos.FormulirResponsesDetailResponse, error) {
	// Get formulir by slug
	formulir, err := s.formulirRepo.GetBySlug(req.Slug)
	if err != nil {
		return nil, errors.New("formulir not found")
	}

	// If role is provided, validate form ownership
	if req.Role != nil && *req.Role != "" {
		role := *req.Role
		
		// Validate: form created_by_role must match request role
		if formulir.CreatedByRole == nil {
			return nil, errors.New("form was not created with a role (created_by_role is NULL)")
		}
		
		if *formulir.CreatedByRole != role {
			return nil, fmt.Errorf("form was created by role '%s', not '%s'", *formulir.CreatedByRole, role)
		}
		
		// Validate: form created_by_user_id must match current user
		if formulir.CreatedByUserID != userID {
			return nil, errors.New("you are not the owner of this form")
		}
		
		// Validate: user exists in respective table
		user, err := s.responseRepo.GetUserByID(userID)
		if err != nil {
			return nil, errors.New("user not found")
		}
		
		username := user.Username
		
		if role == "pendidik" || role == "tendik" {
			// Check in kepegawaian table
			kepegawaian, err := s.responseRepo.GetKepegawaianByUsername(username)
			if err != nil || kepegawaian == nil {
				return nil, fmt.Errorf("user not found in kepegawaian table for role %s", role)
			}
		} else if role == "murid" {
			// Check in peserta_didik table
			pesertaDidik, err := s.responseRepo.GetPesertaDidikByUsername(username)
			if err != nil || pesertaDidik == nil {
				return nil, errors.New("user not found in peserta_didik table for role murid")
			}
		} else {
			return nil, fmt.Errorf("invalid role '%s'. Valid roles: pendidik, tendik, murid", role)
		}
	}

	// Get user's response for this formulir
	response, err := s.responseRepo.GetResponseByFormulirIDAndUserID(formulir.ID, userID)
	if err != nil {
		return nil, errors.New("response not found for this user")
	}

	// If role is specified, validate that submitted_as_role matches
	if req.Role != nil {
		if response.SubmittedAsRole == nil {
			return nil, errors.New("this response has no submitted_as_role")
		}
		if *response.SubmittedAsRole != *req.Role {
			return nil, fmt.Errorf("response was submitted as '%s', not as '%s'", *response.SubmittedAsRole, *req.Role)
		}
	}

	// Map formulir to response DTO
	var targetUserTypes []string
	if formulir.TargetUserTypes != nil {
		_ = json.Unmarshal(formulir.TargetUserTypes, &targetUserTypes)
	}

	// Map pertanyaan
	pertanyaanMap := make(map[uint]*models.FormulirPertanyaan)
	pertanyaanResponses := make([]dtos.FormulirPertanyaanResponse, 0)
	
	for i := range formulir.Pertanyaan {
		p := &formulir.Pertanyaan[i]
		pertanyaanMap[p.ID] = p

		// Parse options
		var options []string
		if p.Options != nil {
			_ = json.Unmarshal(p.Options, &options)
		}

		// Parse validation_rules
		var validationRules map[string]interface{}
		if p.ValidationRules != nil {
			_ = json.Unmarshal(p.ValidationRules, &validationRules)
		}

		// Parse file_config
		var fileConfig map[string]interface{}
		if p.FileConfig != nil {
			_ = json.Unmarshal(p.FileConfig, &fileConfig)
		}

		placeholder := ""
		if p.Placeholder != nil {
			placeholder = *p.Placeholder
		}

		dokumen := ""
		if p.Dokumen != nil && *p.Dokumen != "" {
			dokumen = s.r2Storage.GetPublicURL(*p.Dokumen)
		}

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

	formulirResponse := dtos.FormulirResponse{
		ID:                     formulir.ID,
		Judul:                  formulir.Judul,
		Slug:                   formulir.Slug,
		Deskripsi:              formulir.Deskripsi,
		CreatedByUserID:        formulir.CreatedByUserID,
		IsActive:               formulir.IsActive,
		MaxResponses:           formulir.MaxResponses,
		StartDate:              formulir.StartDate,
		EndDate:                formulir.EndDate,
		AccessType:             formulir.AccessType,
		TargetUserTypes:        targetUserTypes,
		AllowMultipleResponses: formulir.AllowMultipleResponses,
		CreatedAt:              formulir.CreatedAt,
		UpdatedAt:              formulir.UpdatedAt,
		Pertanyaan:             pertanyaanResponses,
		PublicURL:              generatePublicURL(formulir.Slug, formulir.AccessType),
	}

	// Map user basic info with role-based name fetching
	var submittedBy *dtos.UserBasic
	if response.SubmittedBy != nil {
		fullName := response.SubmittedBy.Nama // Default fallback
		
		// Fetch name based on submitted_as_role
		if response.SubmittedAsRole != nil {
			role := *response.SubmittedAsRole
			
			// For pendidik or tendik, fetch from kepegawaian table
			if role == "pendidik" || role == "tendik" {
				kepegawaian, err := s.responseRepo.GetKepegawaianByUsername(response.SubmittedBy.Username)
				if err == nil && kepegawaian != nil {
					fullName = kepegawaian.Nama
				} else {
					fmt.Printf("Warning: User '%s' with role '%s' not found in kepegawaian table. Using fallback name from users table.\n", response.SubmittedBy.Username, role)
				}
			}
			
			// For murid, fetch from peserta_didik table
			if role == "murid" {
				pesertaDidik, err := s.responseRepo.GetPesertaDidikByUsername(response.SubmittedBy.Username)
				if err == nil && pesertaDidik != nil {
					fullName = pesertaDidik.Nama
				} else {
					fmt.Printf("Warning: User '%s' with role 'murid' not found in peserta_didik table. Using fallback name from users table.\n", response.SubmittedBy.Username)
				}
			}
		}
		
		submittedBy = &dtos.UserBasic{
			ID:       response.SubmittedBy.ID,
			Username: response.SubmittedBy.Username,
			FullName: fullName,
		}
	}

	// Map jawaban with question labels and answer IDs
	jawaban := make([]dtos.FormulirResponseAnswerDetail, 0)
	for _, j := range response.Jawaban {
		var jawabanJSON interface{}
		if len(j.JawabanJSON) > 0 {
			_ = json.Unmarshal(j.JawabanJSON, &jawabanJSON)
		}

		// Get question label
		label := ""
		if pertanyaan, exists := pertanyaanMap[j.PertanyaanID]; exists {
			label = pertanyaan.Label
		}

		jawaban = append(jawaban, dtos.FormulirResponseAnswerDetail{
			ID:           j.ID,
			PertanyaanID: j.PertanyaanID,
			Label:        label,
			JawabanText:  j.JawabanText,
			JawabanJSON:  jawabanJSON,
		})
	}

	responseRow := dtos.FormulirResponseRowDetail{
		ResponseID:        response.ID,
		SubmittedByUserID: response.SubmittedByUserID,
		SubmittedAsRole:   response.SubmittedAsRole,
		SubmittedBy:       submittedBy,
		IPAddress:         response.IPAddress,
		UserAgent:         response.UserAgent,
		SubmittedAt:       response.SubmittedAt,
		Jawaban:           jawaban,
	}

	return &dtos.FormulirResponsesDetailResponse{
		Formulir:       formulirResponse,
		Responses:      []dtos.FormulirResponseRowDetail{responseRow},
		TotalResponses: 1,
	}, nil
}

// DeleteResponse deletes a response and all its answers with file cleanup
func (s *FormulirResponseServiceImpl) DeleteResponse(req *dtos.FormulirDeleteResponseRequest, userID uint) error {
	// Get existing response with jawaban
	existingResponse, err := s.responseRepo.GetResponseByID(req.ResponseID)
	if err != nil {
		return errors.New("response not found")
	}

	// Get the formulir to check ownership
	formulir, err := s.formulirRepo.GetByID(existingResponse.FormulirID)
	if err != nil {
		return errors.New("formulir not found")
	}

	// AUTHORIZATION CHECK with role-based logic and form ownership
	authorized := false
	authErrorMsg := "unauthorized: you cannot delete this response"

	// Rule 1: Admin role can delete anything
	if req.Role != nil && *req.Role == "admin" {
		authorized = true
	} else if s.isFormOwner(formulir, userID, existingResponse.SubmittedBy) {
		// Rule 2: Form owner can delete any response in their form
		authorized = true
	} else if req.Role != nil {
		// Rule 3: Role-based authorization for response owner
		requestRole := *req.Role
		
		// Check if submitted_as_role exists and matches request role
		if existingResponse.SubmittedAsRole == nil {
			return errors.New("this response has no submitted_as_role, cannot use role-based authorization")
		}
		
		submittedAsRole := *existingResponse.SubmittedAsRole
		if requestRole != submittedAsRole {
			return fmt.Errorf("unauthorized: role mismatch. Response was submitted as '%s' but you're trying to delete as '%s'", submittedAsRole, requestRole)
		}

		// Check if submitted_by_user_id matches current user
		if existingResponse.SubmittedByUserID == nil {
			return errors.New("cannot delete public response with role-based authorization")
		}
		
		if *existingResponse.SubmittedByUserID != userID {
			return errors.New("unauthorized: this response was not submitted by your user account")
		}

		// Get user details to check username
		if existingResponse.SubmittedBy == nil {
			return errors.New("cannot verify user details for authorization")
		}

		username := existingResponse.SubmittedBy.Username

		if requestRole == "pendidik" || requestRole == "tendik" {
			// Check in kepegawaian table
			kepegawaian, err := s.responseRepo.GetKepegawaianByUsername(username)
			if err == nil && kepegawaian != nil {
				// Username matches, user is authorized
				authorized = true
			} else {
				authErrorMsg = fmt.Sprintf("unauthorized: user not found in kepegawaian table for role %s", requestRole)
			}
		} else if requestRole == "murid" {
			// Check in peserta_didik table
			pesertaDidik, err := s.responseRepo.GetPesertaDidikByUsername(username)
			if err == nil && pesertaDidik != nil {
				// Username matches, user is authorized
				authorized = true
			} else {
				authErrorMsg = "unauthorized: user not found in peserta_didik table for role murid"
			}
		} else {
			// Other roles just check user_id match
			authorized = true
		}
	} else {
		// Rule 4: Default - Only the user who submitted can delete (check by submitted_by_user_id)
		if existingResponse.SubmittedByUserID == nil {
			return errors.New("cannot delete public response without specifying role")
		}
		if *existingResponse.SubmittedByUserID == userID {
			authorized = true
		} else {
			authErrorMsg = "unauthorized: you can only delete your own response"
		}
	}

	if !authorized {
		return errors.New(authErrorMsg)
	}

	// Get all jawaban to delete files
	jawaban, err := s.responseRepo.GetJawabanByResponseID(req.ResponseID)
	if err != nil {
		return err
	}

	// Start transaction
	tx := s.responseRepo.BeginTransaction()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Track deleted files
	deletedFiles := []string{}

	// Delete all jawaban and their files
	for _, j := range jawaban {
		// Delete file from R2 if exists
		if j.JawabanText != nil && *j.JawabanText != "" {
			// Check if it's a file URL (contains "formulir/")
			if len(*j.JawabanText) > 0 {
				err := s.r2Storage.DeleteFile(*j.JawabanText)
				if err == nil {
					deletedFiles = append(deletedFiles, *j.JawabanText)
				}
				// Ignore error if file doesn't exist in R2
			}
		}

		// Delete jawaban record
		if err := tx.Delete(&j).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete jawaban ID %d: %v", j.ID, err)
		}
	}

	// Delete response record
	if err := tx.Delete(existingResponse).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete response: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// isFormOwner checks if the user is the owner of the form based on created_by_role
func (s *FormulirResponseServiceImpl) isFormOwner(formulir *models.Formulir, userID uint, submittedByUser *models.User) bool {
	// Check if created_by_user_id matches
	if formulir.CreatedByUserID != userID {
		return false
	}

	// If created_by_role is NULL, owner is from users table (already checked created_by_user_id)
	if formulir.CreatedByRole == nil {
		return true
	}

	// Role-based ownership check
	// We need to get the username of the current user to check in respective tables
	// Since we only have userID, we need to fetch user details
	// For simplicity, we'll use the submittedByUser if available and matches userID
	// Otherwise, we'll need to query the user by ID
	
	var username string
	if submittedByUser != nil && submittedByUser.ID == userID {
		username = submittedByUser.Username
	} else {
		// Need to get username by userID - we'll query directly from repository
		user, err := s.responseRepo.GetUserByID(userID)
		if err != nil {
			return false
		}
		username = user.Username
	}

	role := *formulir.CreatedByRole

	if role == "pendidik" || role == "tendik" {
		// Check in kepegawaian table
		kepegawaian, err := s.responseRepo.GetKepegawaianByUsername(username)
		if err == nil && kepegawaian != nil {
			return true
		}
	} else if role == "murid" {
		// Check in peserta_didik table
		pesertaDidik, err := s.responseRepo.GetPesertaDidikByUsername(username)
		if err == nil && pesertaDidik != nil {
			return true
		}
	}

	return false
}

// ResetResponses deletes all responses and their answers for a formulir (with file cleanup)
func (s *FormulirResponseServiceImpl) ResetResponses(formulirID uint, role *string, userID uint) error {
	// Get formulir to check ownership
	formulir, err := s.formulirRepo.GetByID(formulirID)
	if err != nil {
		return errors.New("formulir not found")
	}

	// AUTHORIZATION CHECK
	authorized := false
	authErrorMsg := "unauthorized: you cannot reset responses for this formulir"

	// Rule 1: Admin role can reset anything
	if role != nil && *role == "admin" {
		authorized = true
	} else {
		// Rule 2: Form owner can reset responses
		// Check if created_by_user_id matches
		if formulir.CreatedByUserID == userID {
			// If created_by_role is NULL, owner is from users table
			if formulir.CreatedByRole == nil {
				authorized = true
			} else {
				// Validate user exists in respective table based on created_by_role
				user, err := s.responseRepo.GetUserByID(userID)
				if err != nil {
					return errors.New("user not found")
				}

				username := user.Username
				createdByRole := *formulir.CreatedByRole

				if createdByRole == "pendidik" || createdByRole == "tendik" {
					// Check in kepegawaian table
					kepegawaian, err := s.responseRepo.GetKepegawaianByUsername(username)
					if err == nil && kepegawaian != nil {
						authorized = true
					} else {
						authErrorMsg = fmt.Sprintf("unauthorized: user not found in kepegawaian table for role %s", createdByRole)
					}
				} else if createdByRole == "murid" {
					// Check in peserta_didik table
					pesertaDidik, err := s.responseRepo.GetPesertaDidikByUsername(username)
					if err == nil && pesertaDidik != nil {
						authorized = true
					} else {
						authErrorMsg = "unauthorized: user not found in peserta_didik table for role murid"
					}
				} else {
					// Other roles just check user_id match
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

	// Get all responses for this formulir
	responses, err := s.responseRepo.GetAllResponsesByFormulirID(formulirID)
	if err != nil {
		return fmt.Errorf("failed to get responses: %v", err)
	}

	// Start transaction
	tx := s.responseRepo.BeginTransaction()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Track files to delete
	filesToDelete := []string{}

	// Delete all response answers and collect file URLs
	for _, response := range responses {
		// Get jawaban for this response
		jawaban, err := s.responseRepo.GetJawabanByResponseID(response.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get jawaban for response %d: %v", response.ID, err)
		}

		// Collect file URLs from jawaban
		for _, j := range jawaban {
			if j.JawabanText != nil && *j.JawabanText != "" {
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

	// Delete all responses
	if err := tx.Where("form_id = ?", formulirID).Delete(&models.FormulirResponse{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete responses: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Delete all files from R2 (after successful transaction)
	for _, fileKey := range filesToDelete {
		_ = s.r2Storage.DeleteFile(fileKey)
		// Ignore errors for file deletion - file might not exist
	}

	return nil
}

// GetStatisticsBySlug generates statistics for all questions in a form
func (s *FormulirResponseServiceImpl) GetStatisticsBySlug(slug string, role *string, userID uint) (*dtos.FormulirStatisticResponse, error) {
	// Get formulir by slug
	formulir, err := s.formulirRepo.GetBySlug(slug)
	if err != nil {
		return nil, errors.New("formulir not found")
	}

	// Authorization check (same as GetResponsesBySlug)
	if role != nil && *role != "" {
		requestRole := *role
		
		// Admin can access any form
		if requestRole != "admin" {
			// Validate form ownership for non-admin
			if formulir.CreatedByRole == nil {
				return nil, errors.New("form was not created with a role (created_by_role is NULL)")
			}
			
			if *formulir.CreatedByRole != requestRole {
				return nil, fmt.Errorf("form was created by role '%s', not '%s'", *formulir.CreatedByRole, requestRole)
			}
			
			if formulir.CreatedByUserID != userID {
				return nil, errors.New("you are not the owner of this form")
			}
			
			user, err := s.responseRepo.GetUserByID(userID)
			if err != nil {
				return nil, errors.New("user not found")
			}
			
			username := user.Username
			
			if requestRole == "pendidik" || requestRole == "tendik" {
				kepegawaian, err := s.responseRepo.GetKepegawaianByUsername(username)
				if err != nil || kepegawaian == nil {
					return nil, fmt.Errorf("user not found in kepegawaian table for role %s", requestRole)
				}
			} else if requestRole == "murid" {
				pesertaDidik, err := s.responseRepo.GetPesertaDidikByUsername(username)
				if err != nil || pesertaDidik == nil {
					return nil, errors.New("user not found in peserta_didik table for role murid")
				}
			}
		}
	}

	// Get all responses
	responses, err := s.responseRepo.GetAllResponsesByFormulirID(formulir.ID)
	if err != nil {
		return nil, err
	}

	totalResponses := int64(len(responses))

	// Get all jawaban for all responses
	allJawaban := make(map[uint][]models.FormulirResponseJawaban) // pertanyaan_id -> []jawaban
	for _, response := range responses {
		jawaban, err := s.responseRepo.GetJawabanByResponseID(response.ID)
		if err != nil {
			continue
		}
		
		for _, j := range jawaban {
			allJawaban[j.PertanyaanID] = append(allJawaban[j.PertanyaanID], j)
		}
	}

	// Generate statistics for each question
	questionStats := make([]dtos.QuestionStatistic, 0)
	
	for _, pertanyaan := range formulir.Pertanyaan {
		jawaban := allJawaban[pertanyaan.ID]
		totalAnswers := len(jawaban)
		
		stat := dtos.QuestionStatistic{
			PertanyaanID: pertanyaan.ID,
			Label:        pertanyaan.Label,
			Tipe:         pertanyaan.Tipe,
			TotalAnswers: totalAnswers,
			Statistics:   make(map[string]interface{}),
		}

		// Generate statistics based on question type
		switch pertanyaan.Tipe {
		case "radio", "select":
			// Count occurrences of each option
			optionCounts := make(map[string]int)
			for _, j := range jawaban {
				if j.JawabanText != nil && *j.JawabanText != "" {
					optionCounts[*j.JawabanText]++
				}
			}
			
			// Convert to array for easier frontend consumption
			options := []dtos.OptionCount{}
			for option, count := range optionCounts {
				options = append(options, dtos.OptionCount{
					Option: option,
					Count:  count,
				})
			}
			
			stat.Statistics["options"] = options
			stat.Statistics["type"] = "single_choice"

		case "checkbox":
			// Count occurrences of each option (multiple selections)
			optionCounts := make(map[string]int)
			for _, j := range jawaban {
				if len(j.JawabanJSON) > 0 {
					var selections []string
					if err := json.Unmarshal(j.JawabanJSON, &selections); err == nil {
						for _, selection := range selections {
							optionCounts[selection]++
						}
					}
				}
			}
			
			options := []dtos.OptionCount{}
			for option, count := range optionCounts {
				options = append(options, dtos.OptionCount{
					Option: option,
					Count:  count,
				})
			}
			
			stat.Statistics["options"] = options
			stat.Statistics["type"] = "multiple_choice"

		case "number":
			// Calculate min, max, average, sum
			var numbers []float64
			for _, j := range jawaban {
				if j.JawabanText != nil && *j.JawabanText != "" {
					var num float64
					if _, err := fmt.Sscanf(*j.JawabanText, "%f", &num); err == nil {
						numbers = append(numbers, num)
					}
				}
			}
			
			if len(numbers) > 0 {
				min := numbers[0]
				max := numbers[0]
				sum := 0.0
				
				for _, num := range numbers {
					if num < min {
						min = num
					}
					if num > max {
						max = num
					}
					sum += num
				}
				
				avg := sum / float64(len(numbers))
				
				stat.Statistics["min"] = min
				stat.Statistics["max"] = max
				stat.Statistics["average"] = avg
				stat.Statistics["sum"] = sum
				stat.Statistics["type"] = "numeric"
			} else {
				stat.Statistics["type"] = "numeric"
				stat.Statistics["min"] = nil
				stat.Statistics["max"] = nil
				stat.Statistics["average"] = nil
				stat.Statistics["sum"] = nil
			}

		case "text", "email", "phone":
			// List all text responses (limited to show sample)
			responses := []string{}
			for i, j := range jawaban {
				if i >= 100 { // Limit to 100 responses for performance
					break
				}
				if j.JawabanText != nil && *j.JawabanText != "" {
					responses = append(responses, *j.JawabanText)
				}
			}
			
			stat.Statistics["responses"] = responses
			stat.Statistics["type"] = "text"
			stat.Statistics["total_count"] = len(jawaban)

		case "textarea":
			// List all textarea responses (limited)
			responses := []string{}
			for i, j := range jawaban {
				if i >= 50 { // Limit to 50 for long texts
					break
				}
				if j.JawabanText != nil && *j.JawabanText != "" {
					responses = append(responses, *j.JawabanText)
				}
			}
			
			stat.Statistics["responses"] = responses
			stat.Statistics["type"] = "long_text"
			stat.Statistics["total_count"] = len(jawaban)

		case "file":
			// Count uploaded files and list URLs
			fileURLs := []string{}
			uploadedCount := 0
			for _, j := range jawaban {
				if j.JawabanText != nil && *j.JawabanText != "" {
					uploadedCount++
					fileURLs = append(fileURLs, s.r2Storage.GetPublicURL(*j.JawabanText))
				}
			}
			
			stat.Statistics["total_uploaded"] = uploadedCount
			stat.Statistics["file_urls"] = fileURLs
			stat.Statistics["type"] = "file"

		case "date":
			// List all dates
			dates := []string{}
			dateCounts := make(map[string]int)
			for _, j := range jawaban {
				if j.JawabanText != nil && *j.JawabanText != "" {
					dates = append(dates, *j.JawabanText)
					dateCounts[*j.JawabanText]++
				}
			}
			
			stat.Statistics["dates"] = dates
			stat.Statistics["date_distribution"] = dateCounts
			stat.Statistics["type"] = "date"

		case "time":
			// List all times
			times := []string{}
			timeCounts := make(map[string]int)
			for _, j := range jawaban {
				if j.JawabanText != nil && *j.JawabanText != "" {
					times = append(times, *j.JawabanText)
					timeCounts[*j.JawabanText]++
				}
			}
			
			stat.Statistics["times"] = times
			stat.Statistics["time_distribution"] = timeCounts
			stat.Statistics["type"] = "time"

		case "datetime":
			// List all datetimes
			datetimes := []string{}
			for _, j := range jawaban {
				if j.JawabanText != nil && *j.JawabanText != "" {
					datetimes = append(datetimes, *j.JawabanText)
				}
			}
			
			stat.Statistics["datetimes"] = datetimes
			stat.Statistics["type"] = "datetime"
		}

		questionStats = append(questionStats, stat)
	}

	return &dtos.FormulirStatisticResponse{
		FormulirID:     formulir.ID,
		Judul:          formulir.Judul,
		Slug:           formulir.Slug,
		TotalResponses: totalResponses,
		Questions:      questionStats,
	}, nil
}
