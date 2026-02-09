package seeders

import (
	"fmt"

	"gorm.io/gorm"
)

// RoleSeeder handles seeding data for Role
type RoleSeeder struct {
	db *gorm.DB
}

// NewRoleSeeders creates a new Role seeder
func NewRoleSeeders(db *gorm.DB) *RoleSeeder {
	return &RoleSeeder{db: db}
}

// Role represents a role record
type Role struct {
	ID          uint
	Name        string
	Description string
	System      string
	Status      string
}

// Run executes the seeder
func (s *RoleSeeder) Run() error {
	fmt.Println("Seeding roles...")

	roles := []Role{
		{Name: "Administrator", Description: "Administrator", System: "PINTU", Status: "active"},
		{Name: "Kepala Sekolah", Description: "Kepala Sekolah", System: "PINTU", Status: "active"},
	}

	// Check if roles already exist
	var count int64
	s.db.Table("roles").Count(&count)
	if count > 0 {
		fmt.Println("✓ Roles already seeded")
		return nil
	}

	// Insert roles
	for _, role := range roles {
		result := s.db.Table("roles").Create(map[string]interface{}{
			"name":        role.Name,
			"description": role.Description,
			"system":      role.System,
			"status":      role.Status,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to seed role %s: %w", role.Name, result.Error)
		}
	}

	fmt.Printf("✓ %d roles seeded successfully\n", len(roles))
	return nil
}
