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

	// Check DeWayne user
	var user struct {
		ID       uint
		Email    string
		Username string
		RoleID   uint
		RoleName string
	}
	db.Raw(`
		SELECT u.id, u.email, u.username, u.role_id, r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = 'dewayne.frazier@aun.edu.ng'
	`).Scan(&user)

	fmt.Printf("\n=== USER INFO ===\n")
	fmt.Printf("ID: %d\n", user.ID)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Role ID: %d\n", user.RoleID)
	fmt.Printf("Role Name: %s\n", user.RoleName)

	// Check permissions count for school_admin role
	var permCount int64
	db.Raw(`
		SELECT COUNT(*) 
		FROM role_permissions rp
		JOIN roles r ON r.id = rp.role_id
		WHERE r.name = 'school_admin'
	`).Scan(&permCount)

	fmt.Printf("\n=== SCHOOL_ADMIN PERMISSIONS ===\n")
	fmt.Printf("Total permissions: %d\n", permCount)

	// Show first 10 permissions
	var permissions []string
	db.Raw(`
		SELECT p.name 
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON r.id = rp.role_id
		WHERE r.name = 'school_admin'
		ORDER BY p.name
		LIMIT 10
	`).Scan(&permissions)

	fmt.Println("\nFirst 10 permissions:")
	for _, p := range permissions {
		fmt.Printf("  - %s\n", p)
	}
}
