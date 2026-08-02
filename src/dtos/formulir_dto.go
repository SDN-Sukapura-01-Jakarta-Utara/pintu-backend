package dtos

import "time"

// FormulirPertanyaanRequest represents a question in the form
type FormulirPertanyaanRequest struct {
	ID              *uint                  `json:"id" binding:"omitempty"` // ID for update, null for create new
	Urutan          int                    `json:"urutan" binding:"required"`
	Label           string                 `json:"label" binding:"required"`
	Placeholder     string                 `json:"placeholder" binding:"omitempty"`
	Tipe            string                 `json:"tipe" binding:"required,oneof=text textarea number email phone radio checkbox select file date time datetime"`
	IsRequired      bool                   `json:"is_required" binding:"omitempty"`
	Options         []string               `json:"options" binding:"omitempty"`
	ValidationRules map[string]interface{} `json:"validation_rules" binding:"omitempty"`
	FileConfig      map[string]interface{} `json:"file_config" binding:"omitempty"`
	Link            string                 `json:"link" binding:"omitempty"` // URL eksternal (YouTube, artikel, dll)
}

// FormulirUpdateRequest represents the request payload for updating Formulir
type FormulirUpdateRequest struct {
	ID                     uint                        `json:"id" binding:"required"`
	Judul                  string                      `json:"judul" binding:"omitempty"`
	Deskripsi              string                      `json:"deskripsi" binding:"omitempty"`
	Role                   *string                     `json:"role" binding:"omitempty"` // "admin", "pendidik", "tendik", "murid", "orang_tua", etc. (null = users table)
	IsActive               *bool                       `json:"is_active" binding:"omitempty"`
	MaxResponses           *int                        `json:"max_responses" binding:"omitempty"`
	StartDate              string                      `json:"start_date" binding:"omitempty"`
	EndDate                string                      `json:"end_date" binding:"omitempty"`
	AccessType             string                      `json:"access_type" binding:"omitempty,oneof=public authenticated"`
	TargetUserTypes        []string                    `json:"target_user_types" binding:"omitempty"`
	RombelIDs              []int                       `json:"rombel_ids" binding:"omitempty"` // Array of rombel IDs for murid filtering
	AllowMultipleResponses *bool                       `json:"allow_multiple_responses" binding:"omitempty"`
	Pertanyaan             []FormulirPertanyaanRequest `json:"pertanyaan" binding:"omitempty"`
	DokumenToDelete        []int                       `json:"dokumen_to_delete" binding:"omitempty"` // Array of urutan to delete dokumen
}

// FormulirCreateRequest represents the request payload for creating Formulir
type FormulirCreateRequest struct {
	Judul                  string                      `json:"judul" binding:"required"`
	Deskripsi              string                      `json:"deskripsi" binding:"omitempty"`
	Role                   *string                     `json:"role" binding:"omitempty"` // "admin", "pendidik", "tendik", "murid", "orang_tua", etc. (null = users table)
	IsActive               bool                        `json:"is_active" binding:"omitempty"`
	MaxResponses           *int                        `json:"max_responses" binding:"omitempty"`
	StartDate              string                      `json:"start_date" binding:"omitempty"` // Format: YYYY-MM-DD HH:mm:ss or YYYY-MM-DD
	EndDate                string                      `json:"end_date" binding:"omitempty"`   // Format: YYYY-MM-DD HH:mm:ss or YYYY-MM-DD
	AccessType             string                      `json:"access_type" binding:"omitempty,oneof=public authenticated"` // "public" or "authenticated"
	TargetUserTypes        []string                    `json:"target_user_types" binding:"omitempty"` // ["pendidik", "tendik", "murid", "orang_tua", "admin"]
	RombelIDs              []int                       `json:"rombel_ids" binding:"omitempty"` // Array of rombel IDs for murid filtering (NULL = all, [] = no restriction, [1,2,3] = specific)
	AllowMultipleResponses bool                        `json:"allow_multiple_responses" binding:"omitempty"`
	Pertanyaan             []FormulirPertanyaanRequest `json:"pertanyaan" binding:"required,dive"`
}

// FormulirPertanyaanResponse represents a question response
type FormulirPertanyaanResponse struct {
	ID              uint                   `json:"id"`
	FormulirID      uint                   `json:"formulir_id"`
	Urutan          int                    `json:"urutan"`
	Label           string                 `json:"label"`
	Placeholder     string                 `json:"placeholder"`
	Tipe            string                 `json:"tipe"`
	IsRequired      bool                   `json:"is_required"`
	Options         []string               `json:"options"`
	ValidationRules map[string]interface{} `json:"validation_rules"`
	FileConfig      map[string]interface{} `json:"file_config"`
	Dokumen         string                 `json:"dokumen"`
	Link            string                 `json:"link"`
}

// FormulirResponse represents the response payload for Formulir
type FormulirResponse struct {
	ID                     uint                         `json:"id"`
	Judul                  string                       `json:"judul"`
	Slug                   string                       `json:"slug"`
	Deskripsi              string                       `json:"deskripsi"`
	CreatedByUserID        uint                         `json:"created_by_user_id"`
	IsActive               bool                         `json:"is_active"`
	MaxResponses           *int                         `json:"max_responses"`
	StartDate              *time.Time                   `json:"start_date"`
	EndDate                *time.Time                   `json:"end_date"`
	AccessType             string                       `json:"access_type"`
	TargetUserTypes        []string                     `json:"target_user_types"`
	RombelIDs              []int                        `json:"rombel_ids"`
	AllowMultipleResponses bool                         `json:"allow_multiple_responses"`
	CreatedAt              time.Time                    `json:"created_at"`
	UpdatedAt              time.Time                    `json:"updated_at"`
	Pertanyaan             []FormulirPertanyaanResponse `json:"pertanyaan"`
	PublicURL              string                       `json:"public_url"` // Full URL to access the form
}

// FormulirListResponse represents the response for list view (without pertanyaan)
type FormulirListResponse struct {
	ID                     uint       `json:"id"`
	Judul                  string     `json:"judul"`
	Slug                   string     `json:"slug"`
	CreatedByUserID        uint       `json:"created_by_user_id"`
	CreatedBy              *UserBasic `json:"created_by"` // User detail
	IsActive               bool       `json:"is_active"`
	MaxResponses           *int       `json:"max_responses"`
	StartDate              *time.Time `json:"start_date"`
	EndDate                *time.Time `json:"end_date"`
	AccessType             string     `json:"access_type"`
	TargetUserTypes        []string   `json:"target_user_types"`
	RombelIDs              []int      `json:"rombel_ids"`
	AllowMultipleResponses bool       `json:"allow_multiple_responses"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	PublicURL              string     `json:"public_url"`
}

// UserBasic represents basic user info
type UserBasic struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

// FormulirSearchFilter represents search filters
type FormulirSearchFilter struct {
	Judul           string   `json:"judul"`              // Filter by title (partial match)
	CreatedByID     *uint    `json:"created_by_id"`      // Filter by creator, null = all
	Role            *string  `json:"role"`               // Filter by role: null = all, "pendidik"/"tendik"/"murid" = filter by created_by_role and validate user
	StartDate       string   `json:"start_date"`         // Filter created_at >= start_date (YYYY-MM-DD)
	EndDate         string   `json:"end_date"`           // Filter created_at <= end_date (YYYY-MM-DD)
	AccessType      *string  `json:"access_type"`        // Filter by access_type: null = show all (public & authenticated), "" = show all, "public"/"authenticated" = filter by specific type
	TargetUserTypes []string `json:"target_user_types"`  // Filter by target_user_types, empty = all
	RombelID        *int     `json:"rombel_id"`          // Filter by rombel_id: null = all, specific ID = filter forms where rombel_ids is NULL or contains this ID
}

// FormulirPaginationRequest represents pagination parameters
type FormulirPaginationRequest struct {
	Limit int `json:"limit"` // Items per page
	Page  int `json:"page"`  // Page number (starts from 1)
}

// FormulirGetAllRequest represents the request for getting all formulir with filters
type FormulirGetAllRequest struct {
	Search     FormulirSearchFilter      `json:"search"`
	Pagination FormulirPaginationRequest `json:"pagination"`
}

// FormulirListWithPaginationResponse represents paginated list response
type FormulirListWithPaginationResponse struct {
	Data       []FormulirListResponse `json:"data"`
	Pagination PaginationInfo         `json:"pagination"`
}

// FormulirGetBySlugRequest represents request to get formulir by slug
type FormulirGetBySlugRequest struct {
	Slug string `json:"slug" binding:"required"`
}


// FormulirJawabanRequest represents an answer to a question
type FormulirJawabanRequest struct {
	PertanyaanID uint        `json:"pertanyaan_id" binding:"required"`
	JawabanText  *string     `json:"jawaban_text"`
	JawabanJSON  interface{} `json:"jawaban_json"` // For checkbox, radio with multiple values
}

// FormulirJawabanEditRequest represents an answer update in edit response
type FormulirJawabanEditRequest struct {
	ID           *uint       `json:"id"` // Answer ID, null for new answer
	PertanyaanID uint        `json:"pertanyaan_id" binding:"required"`
	JawabanText  *string     `json:"jawaban_text"`
	JawabanJSON  interface{} `json:"jawaban_json"` // For checkbox, radio with multiple values
}

// FormulirSubmitRequest represents the request to submit form response
type FormulirSubmitRequest struct {
	FormulirID      uint                       `json:"formulir_id" binding:"required"`
	SubmittedAsRole *string                    `json:"submitted_as_role"` // "pendidik", "tendik", "murid", "orang_tua", "admin"
	Jawaban         []FormulirJawabanRequest   `json:"jawaban" binding:"required,dive"`
}

// FormulirEditResponseRequest represents the request to edit form response
type FormulirEditResponseRequest struct {
	ResponseID    uint                          `json:"response_id" binding:"required"`
	Role          *string                       `json:"role"` // "pendidik", "tendik", "murid" - for authorization check against submitted_as_role
	Jawaban       []FormulirJawabanEditRequest  `json:"jawaban" binding:"required,dive"`
	FilesToDelete []uint                        `json:"files_to_delete"` // Array of pertanyaan_id to delete files
}

// FormulirSubmitResponse represents the response after submitting
type FormulirSubmitResponse struct {
	ID          uint      `json:"id"`
	FormulirID  uint      `json:"formulir_id"`
	SubmittedAt time.Time `json:"submitted_at"`
	Message     string    `json:"message"`
}

// FormulirResponseDetailRequest represents request to get responses by slug
type FormulirResponseDetailRequest struct {
	Slug string  `json:"slug" binding:"required"`
	Role *string `json:"role"` // "admin", "pendidik", "tendik", "murid" - for authorization check
}

// FormulirResponseByUserRequest represents request to get response by user
type FormulirResponseByUserRequest struct {
	Slug string  `json:"slug" binding:"required"`
	Role *string `json:"role"` // "pendidik", "tendik", "murid" - optional, for role-specific filtering
}

// FormulirDeleteResponseRequest represents request to delete response
type FormulirDeleteResponseRequest struct {
	ResponseID uint    `json:"response_id" binding:"required"`
	Role       *string `json:"role"` // "pendidik", "tendik", "murid" - for authorization check
}

// FormulirDeleteRequest represents request to delete formulir
type FormulirDeleteRequest struct {
	FormulirID uint    `json:"formulir_id" binding:"required"`
	Role       *string `json:"role"` // "admin", "pendidik", "tendik", "murid" - for authorization check
}

// FormulirResetResponseRequest represents request to reset/delete all responses for a form
type FormulirResetResponseRequest struct {
	FormulirID uint    `json:"formulir_id" binding:"required"`
	Role       *string `json:"role"` // "admin", "pendidik", "tendik", "murid" - for authorization check
}

// FormulirResponseAnswerDetail represents answer detail for a question
type FormulirResponseAnswerDetail struct {
	ID           uint        `json:"id"` // Answer ID from form_response_answers table
	PertanyaanID uint        `json:"pertanyaan_id"`
	Label        string      `json:"label"` // Question label for easy reference
	JawabanText  *string     `json:"jawaban_text"`
	JawabanJSON  interface{} `json:"jawaban_json"`
}

// FormulirResponseRowDetail represents one submission with all answers
type FormulirResponseRowDetail struct {
	ResponseID        uint                           `json:"response_id"`
	SubmittedByUserID *uint                          `json:"submitted_by_user_id"`
	SubmittedAsRole   *string                        `json:"submitted_as_role"`
	SubmittedBy       *UserBasic                     `json:"submitted_by"`
	IPAddress         *string                        `json:"ip_address"`
	UserAgent         *string                        `json:"user_agent"`
	SubmittedAt       time.Time                      `json:"submitted_at"`
	Jawaban           []FormulirResponseAnswerDetail `json:"jawaban"` // All answers for this response
}

// FormulirResponsesDetailResponse represents the full response view (like Google Forms)
type FormulirResponsesDetailResponse struct {
	Formulir  FormulirResponse            `json:"formulir"`  // Form details with questions
	Responses []FormulirResponseRowDetail `json:"responses"` // All submissions
	TotalResponses int64                  `json:"total_responses"`
}

// FormulirStatisticRequest represents request to get form statistics
type FormulirStatisticRequest struct {
	Slug string  `json:"slug" binding:"required"`
	Role *string `json:"role"` // "admin", "pendidik", "tendik", "murid" - for authorization check
}

// QuestionStatistic represents statistics for a single question
type QuestionStatistic struct {
	PertanyaanID uint                   `json:"pertanyaan_id"`
	Label        string                 `json:"label"`
	Tipe         string                 `json:"tipe"`
	TotalAnswers int                    `json:"total_answers"`
	Statistics   map[string]interface{} `json:"statistics"` // Different structure based on question type
}

// FormulirStatisticResponse represents the complete statistics response
type FormulirStatisticResponse struct {
	FormulirID     uint                `json:"formulir_id"`
	Judul          string              `json:"judul"`
	Slug           string              `json:"slug"`
	TotalResponses int64               `json:"total_responses"`
	Questions      []QuestionStatistic `json:"questions"`
}

// OptionCount represents count for each option (for radio, checkbox, select)
type OptionCount struct {
	Option string `json:"option"`
	Count  int    `json:"count"`
}

// NumericStats represents statistics for numeric questions
type NumericStats struct {
	Min     *float64 `json:"min"`
	Max     *float64 `json:"max"`
	Average *float64 `json:"average"`
	Sum     *float64 `json:"sum"`
}

// TextStats represents statistics for text questions
type TextStats struct {
	Responses []string `json:"responses"` // List of all text responses
}

// FileStats represents statistics for file upload questions
type FileStats struct {
	TotalUploaded int      `json:"total_uploaded"`
	FileURLs      []string `json:"file_urls"`
}

// DateStats represents statistics for date questions
type DateStats struct {
	Dates []string `json:"dates"` // List of all dates
}

// TimeStats represents statistics for time questions
type TimeStats struct {
	Times []string `json:"times"` // List of all times
}

// DateTimeStats represents statistics for datetime questions
type DateTimeStats struct {
	DateTimes []string `json:"datetimes"` // List of all datetimes
}

// FormulirGetByUserRequest represents request to get formulir filtered by user context
type FormulirGetByUserRequest struct {
	Role      string                    `json:"role" binding:"required,oneof=pendidik tendik murid"` // User role: pendidik, tendik, murid
	RombelID  *int                      `json:"rombel_id"` // Required only for murid role
	StartDate string                    `json:"start_date"` // Filter by start_date (YYYY-MM-DD)
	EndDate   string                    `json:"end_date"`   // Filter by end_date (YYYY-MM-DD)
	Judul     string                    `json:"judul"`      // Filter by title (partial match)
	Pagination FormulirPaginationRequest `json:"pagination"`
}
