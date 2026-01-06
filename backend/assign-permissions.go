package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// Get role IDs
	var platformAdminID, schoolAdminID uint
	db.Raw("SELECT id FROM roles WHERE name = 'platform_admin'").Scan(&platformAdminID)
	db.Raw("SELECT id FROM roles WHERE name = 'school_admin'").Scan(&schoolAdminID)

	// Get all permission IDs
	var allPermissions []uint
	db.Raw("SELECT id FROM permissions").Scan(&allPermissions)

	// Assign ALL permissions to platform_admin
	for _, permID := range allPermissions {
		db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", platformAdminID, permID)
	}
	fmt.Printf("✓ Assigned %d permissions to platform_admin\n", len(allPermissions))

	// Get permissions for school_admin (exclude roles.*, permissions.*, settings.manage)
	var schoolAdminPermissions []uint
	db.Raw(`
		SELECT id FROM permissions 
		WHERE name NOT LIKE 'roles.%' 
		AND name NOT LIKE 'permissions.%'
		AND name != 'settings.manage'
	`).Scan(&schoolAdminPermissions)

	// Assign to school_admin
	for _, permID := range schoolAdminPermissions {
		db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", schoolAdminID, permID)
	}
	fmt.Printf("✓ Assigned %d permissions to school_admin\n", len(schoolAdminPermissions))
}
