package seeders

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder handles seeding data for User
type UserSeeder struct {
	db *gorm.DB
}

// NewUserSeeders creates a new User seeder
func NewUserSeeders(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

// UserData represents user record for seeding
type UserData struct {
	Nama    string
	Username string
	Password string
	RoleIDs  []uint // Many-to-many: array of role IDs
	Status   string
}

// Run executes the seeder
func (s *UserSeeder) Run() error {
	fmt.Println("Seeding users...")

	users := []UserData{
		{
			Nama:    "Administrator",
			Username: "admin",
			Password: "admin01",
			RoleIDs: []uint{1, 3}, // Role ID 1 = Administrator
			Status:   "active",
		},
		{
			Nama:    "Kepala Sekolah",
			Username: "kepsek",
			Password: "kepsek01",
			RoleIDs: []uint{2, 4}, // Role ID 2 = Kepala Sekolah
			Status:   "active",
		},
	}

	// Check if users already exist
	var count int64
	s.db.Table("users").Count(&count)
	if count > 0 {
		fmt.Println("✓ Users already seeded")
		return nil
	}

	// Insert users
	for _, user := range users {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", user.Username, err)
		}

		// Create user
		result := s.db.Table("users").Create(map[string]interface{}{
			"nama":     user.Nama,
			"username": user.Username,
			"password": string(hashedPassword),
			"status":   user.Status,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to seed user %s: %w", user.Username, result.Error)
		}

		// Get user ID
		var userID uint
		s.db.Table("users").Where("username = ?", user.Username).Select("id").Row().Scan(&userID)

		// Assign roles to user (many-to-many)
		for _, roleID := range user.RoleIDs {
			// Insert into user_roles pivot table
			if err := s.db.Table("user_roles").Create(map[string]interface{}{
				"user_id": userID,
				"role_id": roleID,
			}).Error; err != nil {
				return fmt.Errorf("failed to assign role %d to user %s: %w", roleID, user.Username, err)
			}
		}
	}

	fmt.Printf("✓ %d users seeded successfully\n", len(users))
	fmt.Println("⚠️  Remember to change default passwords in production!")
	return nil
}
