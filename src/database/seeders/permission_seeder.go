package seeders

import (
	"fmt"

	"gorm.io/gorm"
)

// PermissionSeeder handles seeding data for Permission
type PermissionSeeder struct {
	db *gorm.DB
}

// NewPermissionSeeders creates a new Permission seeder
func NewPermissionSeeders(db *gorm.DB) *PermissionSeeder {
	return &PermissionSeeder{db: db}
}

// Permission represents a permission record
type Permission struct {
	ID          uint
	Name        string
	Description string
	GroupName   string
	System      string
	Status      string
}

// Run executes the seeder
func (s *PermissionSeeder) Run() error {
	fmt.Println("Seeding permissions...")

	permissions := []Permission{
		// Informasi Sekolah
		{Name: "CREATE_INFORMASI_SEKOLAH", Description: "Create school information", GroupName: "INFORMASI_SEKOLAH", System: "PINTU", Status: "active"},
		{Name: "READ_INFORMASI_SEKOLAH", Description: "Read school information", GroupName: "INFORMASI_SEKOLAH", System: "PINTU", Status: "active"},
		{Name: "UPDATE_INFORMASI_SEKOLAH", Description: "Update school information", GroupName: "INFORMASI_SEKOLAH", System: "PINTU", Status: "active"},
		{Name: "DELETE_INFORMASI_SEKOLAH", Description: "Delete school information", GroupName: "INFORMASI_SEKOLAH", System: "PINTU", Status: "active"},

		// Media
		{Name: "CREATE_MEDIA_PUBLIKASI", Description: "Create media & publikasi", GroupName: "MEDIA_PUBLIKASI", System: "PINTU", Status: "active"},
		{Name: "READ_MEDIA_PUBLIKASI", Description: "Read media & publikasi", GroupName: "MEDIA_PUBLIKASI", System: "PINTU", Status: "active"},
		{Name: "UPDATE_MEDIA_PUBLIKASI", Description: "Update media & publikasi", GroupName: "MEDIA_PUBLIKASI", System: "PINTU", Status: "active"},
		{Name: "DELETE_MEDIA_PUBLIKASI", Description: "Delete media & publikasi", GroupName: "MEDIA_PUBLIKASI", System: "PINTU", Status: "active"},

		// Kepegawaian
		{Name: "CREATE_KEPEGAWAIAN", Description: "Create employee data", GroupName: "KEPEGAWAIAN", System: "PINTU", Status: "active"},
		{Name: "READ_KEPEGAWAIAN", Description: "Read employee data", GroupName: "KEPEGAWAIAN", System: "PINTU", Status: "active"},
		{Name: "UPDATE_KEPEGAWAIAN", Description: "Update employee data", GroupName: "KEPEGAWAIAN", System: "PINTU", Status: "active"},
		{Name: "DELETE_KEPEGAWAIAN", Description: "Delete employee data", GroupName: "KEPEGAWAIAN", System: "PINTU", Status: "active"},

		// Peserta Didik
		{Name: "CREATE_PESERTA_DIDIK", Description: "Create student data", GroupName: "PESERTA_DIDIK", System: "PINTU", Status: "active"},
		{Name: "READ_PESERTA_DIDIK", Description: "Read student data", GroupName: "PESERTA_DIDIK", System: "PINTU", Status: "active"},
		{Name: "UPDATE_PESERTA_DIDIK", Description: "Update student data", GroupName: "PESERTA_DIDIK", System: "PINTU", Status: "active"},
		{Name: "DELETE_PESERTA_DIDIK", Description: "Delete student data", GroupName: "PESERTA_DIDIK", System: "PINTU", Status: "active"},

		// Mutasi Siswa
		{Name: "CREATE_MUTASI_SISWA", Description: "Create student mutation", GroupName: "MUTASI_SISWA", System: "PINTU", Status: "active"},
		{Name: "READ_MUTASI_SISWA", Description: "Read student mutation", GroupName: "MUTASI_SISWA", System: "PINTU", Status: "active"},
		{Name: "UPDATE_MUTASI_SISWA", Description: "Update student mutation", GroupName: "MUTASI_SISWA", System: "PINTU", Status: "active"},
		{Name: "DELETE_MUTASI_SISWA", Description: "Delete student mutation", GroupName: "MUTASI_SISWA", System: "PINTU", Status: "active"},

		// Kritik Saran
		{Name: "CREATE_KRITIK_SARAN", Description: "Create criticism and suggestion", GroupName: "KRITIK_SARAN", System: "PINTU", Status: "active"},
		{Name: "READ_KRITIK_SARAN", Description: "Read criticism and suggestion", GroupName: "KRITIK_SARAN", System: "PINTU", Status: "active"},
		{Name: "UPDATE_KRITIK_SARAN", Description: "Update criticism and suggestion", GroupName: "KRITIK_SARAN", System: "PINTU", Status: "active"},
		{Name: "DELETE_KRITIK_SARAN", Description: "Delete criticism and suggestion", GroupName: "KRITIK_SARAN", System: "PINTU", Status: "active"},

		// Pertanyaan dan Pengaduan
		{Name: "CREATE_PERTANYAAN_PENGADUAN", Description: "Create question and complaint", GroupName: "PERTANYAAN_PENGADUAN", System: "PINTU", Status: "active"},
		{Name: "READ_PERTANYAAN_PENGADUAN", Description: "Read question and complaint", GroupName: "PERTANYAAN_PENGADUAN", System: "PINTU", Status: "active"},
		{Name: "UPDATE_PERTANYAAN_PENGADUAN", Description: "Update question and complaint", GroupName: "PERTANYAAN_PENGADUAN", System: "PINTU", Status: "active"},
		{Name: "DELETE_PERTANYAAN_PENGADUAN", Description: "Delete question and complaint", GroupName: "PERTANYAAN_PENGADUAN", System: "PINTU", Status: "active"},
	}

	// Check if permissions already exist
	var count int64
	s.db.Table("permissions").Count(&count)
	if count > 0 {
		fmt.Println("✓ Permissions already seeded")
		return nil
	}

	// Insert permissions
	for _, permission := range permissions {
		result := s.db.Table("permissions").Create(map[string]interface{}{
			"name":        permission.Name,
			"description": permission.Description,
			"group_name":  permission.GroupName,
			"system":      permission.System,
			"status":      permission.Status,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to seed permission %s: %w", permission.Name, result.Error)
		}
	}

	fmt.Printf("✓ %d permissions seeded successfully\n", len(permissions))
	return nil
}
