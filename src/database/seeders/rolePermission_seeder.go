package seeders

import (
	"fmt"

	"gorm.io/gorm"
)

// RolePermissionSeeder handles seeding data for RolePermission
type RolePermissionSeeder struct {
	db *gorm.DB
}

// NewRolePermissionSeeders creates a new RolePermission seeder
func NewRolePermissionSeeders(db *gorm.DB) *RolePermissionSeeder {
	return &RolePermissionSeeder{db: db}
}

// RolePermissionMap represents role-permission mapping
type RolePermissionMap struct {
	RoleName        string
	PermissionNames []string
}

// Run executes the seeder
func (s *RolePermissionSeeder) Run() error {
	fmt.Println("Seeding role permissions...")

	// Define role-permission mappings
	mappings := []RolePermissionMap{
		{
			RoleName: "Administrator",
			PermissionNames: []string{
				// Informasi Sekolah
				"CREATE_INFORMASI_SEKOLAH", "READ_INFORMASI_SEKOLAH", "UPDATE_INFORMASI_SEKOLAH", "DELETE_INFORMASI_SEKOLAH",
				// Media Publikasi
				"CREATE_MEDIA_PUBLIKASI", "READ_MEDIA_PUBLIKASI", "UPDATE_MEDIA_PUBLIKASI", "DELETE_MEDIA_PUBLIKASI",
				// Kepegawaian
				"CREATE_KEPEGAWAIAN", "READ_KEPEGAWAIAN", "UPDATE_KEPEGAWAIAN", "DELETE_KEPEGAWAIAN",
				// Peserta Didik
				"CREATE_PESERTA_DIDIK", "READ_PESERTA_DIDIK", "UPDATE_PESERTA_DIDIK", "DELETE_PESERTA_DIDIK",
				// Mutasi Siswa
				"CREATE_MUTASI_SISWA", "READ_MUTASI_SISWA", "UPDATE_MUTASI_SISWA", "DELETE_MUTASI_SISWA",
				// Kritik Saran
				"CREATE_KRITIK_SARAN", "READ_KRITIK_SARAN", "UPDATE_KRITIK_SARAN", "DELETE_KRITIK_SARAN",
				// Pertanyaan dan Pengaduan
				"CREATE_PERTANYAAN_PENGADUAN", "READ_PERTANYAAN_PENGADUAN", "UPDATE_PERTANYAAN_PENGADUAN", "DELETE_PERTANYAAN_PENGADUAN",
			},
		},
		{
			RoleName: "Kepala Sekolah",
			PermissionNames: []string{
				// Informasi Sekolah - Read Only
				"READ_INFORMASI_SEKOLAH",
				// Media - Read Only
				"READ_MEDIA_PUBLIKASI",
				// Kepegawaian - Read Only
				"READ_KEPEGAWAIAN",
				// Peserta Didik - Read Only
				"READ_PESERTA_DIDIK",
				// Mutasi Siswa - Read Only
				"READ_MUTASI_SISWA",
				// Kritik Saran - Read Only
				"READ_KRITIK_SARAN",
				// Pertanyaan dan Pengaduan - Read Only
				"READ_PERTANYAAN_PENGADUAN",
			},
		},
	}

	// Check if role_permissions already exist
	var count int64
	s.db.Table("role_permissions").Count(&count)
	if count > 0 {
		fmt.Println("✓ Role permissions already seeded")
		return nil
	}

	// Assign permissions to roles
	totalAssigned := 0
	for _, mapping := range mappings {
		// Get role ID
		var roleID uint
		if err := s.db.Table("roles").Where("name = ?", mapping.RoleName).Select("id").Row().Scan(&roleID); err != nil {
			return fmt.Errorf("failed to find role %s: %w", mapping.RoleName, err)
		}

		// Assign each permission to this role
		for _, permissionName := range mapping.PermissionNames {
			// Get permission ID
			var permissionID uint
			if err := s.db.Table("permissions").Where("name = ?", permissionName).Select("id").Row().Scan(&permissionID); err != nil {
				return fmt.Errorf("failed to find permission %s: %w", permissionName, err)
			}

			// Create role-permission association
			result := s.db.Table("role_permissions").Create(map[string]interface{}{
				"role_id":       roleID,
				"permission_id": permissionID,
			})
			if result.Error != nil {
				return fmt.Errorf("failed to assign permission %s to role %s: %w", permissionName, mapping.RoleName, result.Error)
			}
			totalAssigned++
		}
	}

	fmt.Printf("✓ %d role-permission associations seeded successfully\n", totalAssigned)
	return nil
}
