package services

import (
	"errors"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"pintu-backend/src/utils"

	"golang.org/x/crypto/bcrypt"
)

// LoginService handles business logic for authentication
type LoginService interface {
	Login(req *dtos.LoginRequest) (*dtos.LoginResponse, error)
	LoginStudent(req *dtos.LoginRequest) (*dtos.StudentLoginResponse, error)
	LoginSieksa(req *dtos.LoginRequest) (*dtos.LoginResponse, error)
	LoginStudentSieksa(req *dtos.LoginRequest) (*dtos.StudentLoginResponse, error)
}

type LoginServiceImpl struct {
	repository repositories.LoginRepository
}

// NewLoginService creates a new Login service
func NewLoginService(repository repositories.LoginRepository) LoginService {
	return &LoginServiceImpl{repository: repository}
}

// Login authenticates user and returns JWT token
func (s *LoginServiceImpl) Login(req *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	// Try to get user from users table first
	user, userErr := s.repository.GetByUsername(req.Username)
	
	// If user not found in users table, try kepegawaian table
	if userErr != nil {
		kepegawaian, kepErr := s.repository.GetKepegawaianByUsername(req.Username)
		if kepErr != nil {
			return nil, errors.New("username atau password salah")
		}
		
		// Verify password for kepegawaian
		if err := bcrypt.CompareHashAndPassword([]byte(kepegawaian.Password), []byte(req.Password)); err != nil {
			return nil, errors.New("username atau password salah")
		}

		// Check if kepegawaian is active
		if kepegawaian.Status != "active" {
			return nil, errors.New("user tidak aktif")
		}

		// Filter roles that have system_id = 1 (PINTU)
		var pintuRoles []models.Role
		for _, role := range kepegawaian.Roles {
			if role.SystemID != nil && *role.SystemID == 1 {
				pintuRoles = append(pintuRoles, role)
			}
		}

		// Check if kepegawaian has at least one PINTU role
		if len(pintuRoles) == 0 {
			return nil, errors.New("anda tidak memiliki akses ke sistem PINTU")
		}

		// Use first PINTU role for token
		roleID := pintuRoles[0].ID

		// Generate JWT token for kepegawaian
		token, err := utils.GenerateToken(kepegawaian.ID, kepegawaian.Username, kepegawaian.Nama, &roleID, kepegawaian.Status, "PINTU")
		if err != nil {
			return nil, errors.New("gagal membuat token")
		}

		// Map PINTU roles and extract permissions
		roles := make([]dtos.RoleResponse, len(pintuRoles))
		permissionMap := make(map[string]bool)
		
		for i, role := range pintuRoles {
			var system *dtos.SystemResponse
			if role.System != nil {
				system = &dtos.SystemResponse{
					ID:          role.System.ID,
					Nama:        role.System.Nama,
					Description: role.System.Description,
				}
			}

			roles[i] = dtos.RoleResponse{
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

			// Collect unique permissions from PINTU roles
			for _, permission := range role.Permissions {
				permissionMap[permission.Name] = true
			}
		}

		// Convert permission map to slice
		permissions := make([]string, 0, len(permissionMap))
		for permName := range permissionMap {
			permissions = append(permissions, permName)
		}

		// Prepare rombel_bidang_studi (handle JSONB data)
		var rombelBidangStudi interface{} = nil
		if len(kepegawaian.RombelBidangStudi) > 0 {
			rombelBidangStudi = kepegawaian.RombelBidangStudi
		}

		// Set jabatan to pointer
		var jabatan *string = nil
		if kepegawaian.Jabatan != "" {
			jabatan = &kepegawaian.Jabatan
		}

		response := &dtos.LoginResponse{
			Token:       token,
			ExpiresAt:   time.Now().Add(24 * time.Hour),
			Permissions: permissions,
			User: dtos.UserLoginResponse{
				ID:                kepegawaian.ID,
				Nama:              kepegawaian.Nama,
				Username:          kepegawaian.Username,
				Status:            kepegawaian.Status,
				Roles:             roles,
				CreatedAt:         kepegawaian.CreatedAt,
				Jabatan:           jabatan,
				RombelGuruKelasID: kepegawaian.RombelGuruKelasID,
				BidangStudiID:     kepegawaian.BidangStudiID,
				RombelBidangStudi: rombelBidangStudi,
			},
		}

		return response, nil
	}

	// User found in users table, proceed with normal user login
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("username atau password salah")
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("user tidak aktif")
	}

	// Filter roles that have system_id = 1 (PINTU)
	var pintuRoles []models.Role
	for _, role := range user.Roles {
		if role.SystemID != nil && *role.SystemID == 1 {
			pintuRoles = append(pintuRoles, role)
		}
	}

	// Check if user has at least one PINTU role
	if len(pintuRoles) == 0 {
		return nil, errors.New("anda tidak memiliki akses ke sistem PINTU")
	}

	// Use first PINTU role for token
	roleID := pintuRoles[0].ID

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Nama, &roleID, user.Status, "PINTU")
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	// Map PINTU roles and extract permissions
	roles := make([]dtos.RoleResponse, len(pintuRoles))
	permissionMap := make(map[string]bool) // To track unique permissions
	
	for i, role := range pintuRoles {
		var system *dtos.SystemResponse
		if role.System != nil {
			system = &dtos.SystemResponse{
				ID:          role.System.ID,
				Nama:        role.System.Nama,
				Description: role.System.Description,
			}
		}

		roles[i] = dtos.RoleResponse{
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

		// Collect unique permissions from PINTU roles
		for _, permission := range role.Permissions {
			permissionMap[permission.Name] = true
		}
	}

	// Convert permission map to slice
	permissions := make([]string, 0, len(permissionMap))
	for permName := range permissionMap {
		permissions = append(permissions, permName)
	}

	response := &dtos.LoginResponse{
		Token:       token,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		Permissions: permissions,
		User: dtos.UserLoginResponse{
			ID:        user.ID,
			Nama:      user.Nama,
			Username:  user.Username,
			Status:    user.Status,
			Roles:     roles,
			CreatedAt: user.CreatedAt,
		},
	}

	return response, nil
}

// LoginStudent authenticates student and returns JWT token
func (s *LoginServiceImpl) LoginStudent(req *dtos.LoginRequest) (*dtos.StudentLoginResponse, error) {
	// Get student from peserta_didik table
	student, err := s.repository.GetPesertaDidikByUsername(req.Username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("username atau password salah")
	}

	// Check if student is active
	if student.Status != "active" {
		return nil, errors.New("akun siswa tidak aktif")
	}

	// Filter roles that have system_id = 1 (PINTU)
	var pintuRoles []models.Role
	for _, role := range student.Roles {
		if role.SystemID != nil && *role.SystemID == 1 {
			pintuRoles = append(pintuRoles, role)
		}
	}

	// Check if student has at least one PINTU role
	if len(pintuRoles) == 0 {
		return nil, errors.New("anda tidak memiliki akses ke sistem PINTU")
	}

	// Use first PINTU role for token
	roleID := pintuRoles[0].ID

	// Generate JWT token for student
	token, err := utils.GenerateToken(student.ID, student.Username, student.Nama, &roleID, student.Status, "PINTU")
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	// Map PINTU roles and extract permissions
	roles := make([]dtos.RoleResponse, len(pintuRoles))
	permissionMap := make(map[string]bool)
	
	for i, role := range pintuRoles {
		var system *dtos.SystemResponse
		if role.System != nil {
			system = &dtos.SystemResponse{
				ID:          role.System.ID,
				Nama:        role.System.Nama,
				Description: role.System.Description,
			}
		}

		roles[i] = dtos.RoleResponse{
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

		// Collect unique permissions from PINTU roles
		for _, permission := range role.Permissions {
			permissionMap[permission.Name] = true
		}
	}

	// Convert permission map to slice
	permissions := make([]string, 0, len(permissionMap))
	for permName := range permissionMap {
		permissions = append(permissions, permName)
	}

	// Get student's rombel data
	rombelData, err := s.repository.GetPesertaDidikRombelByStudentID(student.ID)
	if err != nil {
		// If error getting rombel, just return empty array
		rombelData = []models.PesertaDidikRombel{}
	}

	// Map rombel data
	rombelResponse := make([]dtos.StudentRombelLoginResponse, len(rombelData))
	for i, rombel := range rombelData {
		kelasID := uint(0)
		kelasName := ""
		if rombel.Rombel != nil {
			if rombel.Rombel.Kelas != nil {
				kelasID = rombel.Rombel.Kelas.ID
				kelasName = rombel.Rombel.Kelas.Name
			}
		}

		tahunPelajaran := ""
		if rombel.TahunPelajaran != nil {
			tahunPelajaran = rombel.TahunPelajaran.TahunPelajaran
		}

		rombelName := ""
		rombelIDValue := uint(0)
		if rombel.Rombel != nil {
			rombelName = rombel.Rombel.Name
			rombelIDValue = rombel.Rombel.ID
		}

		rombelResponse[i] = dtos.StudentRombelLoginResponse{
			ID:               rombel.ID,
			RombelID:         rombelIDValue,
			RombelName:       rombelName,
			KelasID:          kelasID,
			KelasName:        kelasName,
			TahunPelajaranID: rombel.TahunPelajaranID,
			TahunPelajaran:   tahunPelajaran,
			Status:           rombel.Status,
		}
	}

	response := &dtos.StudentLoginResponse{
		Token:       token,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		Permissions: permissions,
		Student: dtos.StudentDetailLoginResponse{
			ID:           student.ID,
			Nama:         student.Nama,
			NIS:          student.NIS,
			NISN:         student.NISN,
			JenisKelamin: student.JenisKelamin,
			TempatLahir:  student.TempatLahir,
			TanggalLahir: student.TanggalLahir,
			Username:     student.Username,
			Status:       student.Status,
			Photo:        student.Photo,
			Roles:        roles,
			Rombel:       rombelResponse,
			CreatedAt:    student.CreatedAt,
		},
	}

	return response, nil
}

// LoginSieksa authenticates user for SIEKSA application and returns JWT token
func (s *LoginServiceImpl) LoginSieksa(req *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	// Try to get user from users table first
	user, userErr := s.repository.GetByUsername(req.Username)
	
	// If user not found in users table, try kepegawaian table
	if userErr != nil {
		kepegawaian, kepErr := s.repository.GetKepegawaianByUsername(req.Username)
		if kepErr != nil {
			// If not found in kepegawaian, try pelatih table
			pelatih, pelatihErr := s.repository.GetPelatihByUsername(req.Username)
			if pelatihErr != nil {
				return nil, errors.New("username atau password salah")
			}

			// Verify password for pelatih
			if err := bcrypt.CompareHashAndPassword([]byte(pelatih.Password), []byte(req.Password)); err != nil {
				return nil, errors.New("username atau password salah")
			}

			// Check if pelatih is active
			if pelatih.Status != "active" {
				return nil, errors.New("user tidak aktif")
			}

			// Filter roles that have system_id = 2 (SIEKSA)
			var sieksaRoles []models.Role
			for _, role := range pelatih.Roles {
				if role.SystemID != nil && *role.SystemID == 2 {
					sieksaRoles = append(sieksaRoles, role)
				}
			}

			// Check if pelatih has at least one SIEKSA role
			if len(sieksaRoles) == 0 {
				return nil, errors.New("anda tidak memiliki akses ke sistem SIEKSA")
			}

			// Use first SIEKSA role for token
			roleID := sieksaRoles[0].ID

			// Generate JWT token for pelatih with SIEKSA app name
			token, err := utils.GenerateToken(pelatih.ID, *pelatih.Username, pelatih.Nama, &roleID, pelatih.Status, "SIEKSA")
			if err != nil {
				return nil, errors.New("gagal membuat token")
			}

			// Map SIEKSA roles and extract permissions
			roles := make([]dtos.RoleResponse, len(sieksaRoles))
			permissionMap := make(map[string]bool)
			
			for i, role := range sieksaRoles {
				var system *dtos.SystemResponse
				if role.System != nil {
					system = &dtos.SystemResponse{
						ID:          role.System.ID,
						Nama:        role.System.Nama,
						Description: role.System.Description,
					}
				}

				roles[i] = dtos.RoleResponse{
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

				// Collect unique permissions from SIEKSA roles
				for _, permission := range role.Permissions {
					permissionMap[permission.Name] = true
				}
			}

			// Convert permission map to slice
			permissions := make([]string, 0, len(permissionMap))
			for permName := range permissionMap {
				permissions = append(permissions, permName)
			}

			response := &dtos.LoginResponse{
				Token:       token,
				ExpiresAt:   time.Now().Add(24 * time.Hour),
				Permissions: permissions,
				User: dtos.UserLoginResponse{
					ID:        pelatih.ID,
					Nama:      pelatih.Nama,
					Username:  *pelatih.Username,
					Status:    pelatih.Status,
					Roles:     roles,
					CreatedAt: pelatih.CreatedAt,
				},
			}

			return response, nil
		}
		
		// Verify password for kepegawaian
		if err := bcrypt.CompareHashAndPassword([]byte(kepegawaian.Password), []byte(req.Password)); err != nil {
			return nil, errors.New("username atau password salah")
		}

		// Check if kepegawaian is active
		if kepegawaian.Status != "active" {
			return nil, errors.New("user tidak aktif")
		}

		// Filter roles that have system_id = 2 (SIEKSA)
		var sieksaRoles []models.Role
		for _, role := range kepegawaian.Roles {
			if role.SystemID != nil && *role.SystemID == 2 {
				sieksaRoles = append(sieksaRoles, role)
			}
		}

		// Check if kepegawaian has at least one SIEKSA role
		if len(sieksaRoles) == 0 {
			return nil, errors.New("anda tidak memiliki akses ke sistem SIEKSA")
		}

		// Use first SIEKSA role for token
		roleID := sieksaRoles[0].ID

		// Generate JWT token for kepegawaian with SIEKSA app name
		token, err := utils.GenerateToken(kepegawaian.ID, kepegawaian.Username, kepegawaian.Nama, &roleID, kepegawaian.Status, "SIEKSA")
		if err != nil {
			return nil, errors.New("gagal membuat token")
		}

		// Map SIEKSA roles and extract permissions
		roles := make([]dtos.RoleResponse, len(sieksaRoles))
		permissionMap := make(map[string]bool)
		
		for i, role := range sieksaRoles {
			var system *dtos.SystemResponse
			if role.System != nil {
				system = &dtos.SystemResponse{
					ID:          role.System.ID,
					Nama:        role.System.Nama,
					Description: role.System.Description,
				}
			}

			roles[i] = dtos.RoleResponse{
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

			// Collect unique permissions from SIEKSA roles
			for _, permission := range role.Permissions {
				permissionMap[permission.Name] = true
			}
		}

		// Convert permission map to slice
		permissions := make([]string, 0, len(permissionMap))
		for permName := range permissionMap {
			permissions = append(permissions, permName)
		}

		// Prepare rombel_bidang_studi (handle JSONB data)
		var rombelBidangStudi interface{} = nil
		if len(kepegawaian.RombelBidangStudi) > 0 {
			rombelBidangStudi = kepegawaian.RombelBidangStudi
		}

		// Set jabatan to pointer
		var jabatan *string = nil
		if kepegawaian.Jabatan != "" {
			jabatan = &kepegawaian.Jabatan
		}

		response := &dtos.LoginResponse{
			Token:       token,
			ExpiresAt:   time.Now().Add(24 * time.Hour),
			Permissions: permissions,
			User: dtos.UserLoginResponse{
				ID:                kepegawaian.ID,
				Nama:              kepegawaian.Nama,
				Username:          kepegawaian.Username,
				Status:            kepegawaian.Status,
				Roles:             roles,
				CreatedAt:         kepegawaian.CreatedAt,
				Jabatan:           jabatan,
				RombelGuruKelasID: kepegawaian.RombelGuruKelasID,
				BidangStudiID:     kepegawaian.BidangStudiID,
				RombelBidangStudi: rombelBidangStudi,
			},
		}

		return response, nil
	}

	// User found in users table, proceed with normal user login
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("username atau password salah")
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("user tidak aktif")
	}

	// Filter roles that have system_id = 2 (SIEKSA)
	var sieksaRoles []models.Role
	for _, role := range user.Roles {
		if role.SystemID != nil && *role.SystemID == 2 {
			sieksaRoles = append(sieksaRoles, role)
		}
	}

	// Check if user has at least one SIEKSA role
	if len(sieksaRoles) == 0 {
		return nil, errors.New("anda tidak memiliki akses ke sistem SIEKSA")
	}

	// Use first SIEKSA role for token
	roleID := sieksaRoles[0].ID

	// Generate JWT token with SIEKSA app name
	token, err := utils.GenerateToken(user.ID, user.Username, user.Nama, &roleID, user.Status, "SIEKSA")
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	// Map SIEKSA roles and extract permissions
	roles := make([]dtos.RoleResponse, len(sieksaRoles))
	permissionMap := make(map[string]bool) // To track unique permissions
	
	for i, role := range sieksaRoles {
		var system *dtos.SystemResponse
		if role.System != nil {
			system = &dtos.SystemResponse{
				ID:          role.System.ID,
				Nama:        role.System.Nama,
				Description: role.System.Description,
			}
		}

		roles[i] = dtos.RoleResponse{
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

		// Collect unique permissions from SIEKSA roles
		for _, permission := range role.Permissions {
			permissionMap[permission.Name] = true
		}
	}

	// Convert permission map to slice
	permissions := make([]string, 0, len(permissionMap))
	for permName := range permissionMap {
		permissions = append(permissions, permName)
	}

	response := &dtos.LoginResponse{
		Token:       token,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		Permissions: permissions,
		User: dtos.UserLoginResponse{
			ID:        user.ID,
			Nama:      user.Nama,
			Username:  user.Username,
			Status:    user.Status,
			Roles:     roles,
			CreatedAt: user.CreatedAt,
		},
	}

	return response, nil
}

// LoginStudentSieksa authenticates student for SIEKSA application and returns JWT token
func (s *LoginServiceImpl) LoginStudentSieksa(req *dtos.LoginRequest) (*dtos.StudentLoginResponse, error) {
	// Get student from peserta_didik table
	student, err := s.repository.GetPesertaDidikByUsername(req.Username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("username atau password salah")
	}

	// Check if student is active
	if student.Status != "active" {
		return nil, errors.New("akun siswa tidak aktif")
	}

	// Filter roles that have system_id = 2 (SIEKSA)
	var sieksaRoles []models.Role
	for _, role := range student.Roles {
		if role.SystemID != nil && *role.SystemID == 2 {
			sieksaRoles = append(sieksaRoles, role)
		}
	}

	// Check if student has at least one SIEKSA role
	if len(sieksaRoles) == 0 {
		return nil, errors.New("anda tidak memiliki akses ke sistem SIEKSA")
	}

	// Use first SIEKSA role for token
	roleID := sieksaRoles[0].ID

	// Generate JWT token for student with SIEKSA app name
	token, err := utils.GenerateToken(student.ID, student.Username, student.Nama, &roleID, student.Status, "SIEKSA")
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	// Map SIEKSA roles and extract permissions
	roles := make([]dtos.RoleResponse, len(sieksaRoles))
	permissionMap := make(map[string]bool)
	
	for i, role := range sieksaRoles {
		var system *dtos.SystemResponse
		if role.System != nil {
			system = &dtos.SystemResponse{
				ID:          role.System.ID,
				Nama:        role.System.Nama,
				Description: role.System.Description,
			}
		}

		roles[i] = dtos.RoleResponse{
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

		// Collect unique permissions from SIEKSA roles
		for _, permission := range role.Permissions {
			permissionMap[permission.Name] = true
		}
	}

	// Convert permission map to slice
	permissions := make([]string, 0, len(permissionMap))
	for permName := range permissionMap {
		permissions = append(permissions, permName)
	}

	// Get student's rombel data
	rombelData, err := s.repository.GetPesertaDidikRombelByStudentID(student.ID)
	if err != nil {
		// If error getting rombel, just return empty array
		rombelData = []models.PesertaDidikRombel{}
	}

	// Map rombel data
	rombelResponse := make([]dtos.StudentRombelLoginResponse, len(rombelData))
	for i, rombel := range rombelData {
		kelasID := uint(0)
		kelasName := ""
		if rombel.Rombel != nil {
			if rombel.Rombel.Kelas != nil {
				kelasID = rombel.Rombel.Kelas.ID
				kelasName = rombel.Rombel.Kelas.Name
			}
		}

		tahunPelajaran := ""
		if rombel.TahunPelajaran != nil {
			tahunPelajaran = rombel.TahunPelajaran.TahunPelajaran
		}

		rombelName := ""
		rombelIDValue := uint(0)
		if rombel.Rombel != nil {
			rombelName = rombel.Rombel.Name
			rombelIDValue = rombel.Rombel.ID
		}

		rombelResponse[i] = dtos.StudentRombelLoginResponse{
			ID:               rombel.ID,
			RombelID:         rombelIDValue,
			RombelName:       rombelName,
			KelasID:          kelasID,
			KelasName:        kelasName,
			TahunPelajaranID: rombel.TahunPelajaranID,
			TahunPelajaran:   tahunPelajaran,
			Status:           rombel.Status,
		}
	}

	response := &dtos.StudentLoginResponse{
		Token:       token,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		Permissions: permissions,
		Student: dtos.StudentDetailLoginResponse{
			ID:           student.ID,
			Nama:         student.Nama,
			NIS:          student.NIS,
			NISN:         student.NISN,
			JenisKelamin: student.JenisKelamin,
			TempatLahir:  student.TempatLahir,
			TanggalLahir: student.TanggalLahir,
			Username:     student.Username,
			Status:       student.Status,
			Photo:        student.Photo,
			Roles:        roles,
			Rombel:       rombelResponse,
			CreatedAt:    student.CreatedAt,
		},
	}

	return response, nil
}
