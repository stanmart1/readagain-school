package main

import (
	"log"
	"readagain/internal/config"
	"readagain/internal/database"
	"readagain/internal/models"
)

func main() {
	cfg := config.Load()
	if err := database.Connect(cfg.Database.URL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get all permissions
	var allPermissions []models.Permission
	if err := database.DB.Find(&allPermissions).Error; err != nil {
		log.Fatal("Failed to fetch permissions:", err)
	}
	log.Printf("Found %d permissions", len(allPermissions))

	// Get platform_admin role
	var platformAdminRole models.Role
	if err := database.DB.Where("name = ?", "platform_admin").First(&platformAdminRole).Error; err != nil {
		log.Fatal("Failed to find platform_admin role:", err)
	}

	// Assign all permissions to platform_admin
	if err := database.DB.Model(&platformAdminRole).Association("Permissions").Replace(&allPermissions); err != nil {
		log.Fatal("Failed to assign permissions to platform_admin:", err)
	}
	log.Printf("✓ Assigned all %d permissions to platform_admin role", len(allPermissions))

	// Get school_admin role
	var schoolAdminRole models.Role
	if err := database.DB.Where("name = ?", "school_admin").First(&schoolAdminRole).Error; err != nil {
		log.Fatal("Failed to find school_admin role:", err)
	}

	// School admin permissions (manage school-level resources)
	schoolAdminPerms := []string{
		"users.view", "users.create", "users.update",
		"books.view", "books.create", "books.update",
		"library.view", "library.manage",
		"reading.view_analytics",
		"reviews.view", "reviews.moderate",
		"reports.view",
		"blog.view",
		"about.view",
		"contact.view",
		"faq.view",
	}

	var schoolAdminPermissions []models.Permission
	for _, permName := range schoolAdminPerms {
		for _, perm := range allPermissions {
			if perm.Name == permName {
				schoolAdminPermissions = append(schoolAdminPermissions, perm)
				break
			}
		}
	}

	if err := database.DB.Model(&schoolAdminRole).Association("Permissions").Replace(&schoolAdminPermissions); err != nil {
		log.Fatal("Failed to assign permissions to school_admin:", err)
	}
	log.Printf("✓ Assigned %d permissions to school_admin role", len(schoolAdminPermissions))

	// Get teacher role
	var teacherRole models.Role
	if err := database.DB.Where("name = ?", "teacher").First(&teacherRole).Error; err != nil {
		log.Fatal("Failed to find teacher role:", err)
	}

	// Teacher permissions (view and manage their classes)
	teacherPerms := []string{
		"books.view",
		"library.view",
		"reading.view_analytics",
		"reviews.view",
	}

	var teacherPermissions []models.Permission
	for _, permName := range teacherPerms {
		for _, perm := range allPermissions {
			if perm.Name == permName {
				teacherPermissions = append(teacherPermissions, perm)
				break
			}
		}
	}

	if err := database.DB.Model(&teacherRole).Association("Permissions").Replace(&teacherPermissions); err != nil {
		log.Fatal("Failed to assign permissions to teacher:", err)
	}
	log.Printf("✓ Assigned %d permissions to teacher role", len(teacherPermissions))

	// Get student role
	var studentRole models.Role
	if err := database.DB.Where("name = ?", "student").First(&studentRole).Error; err != nil {
		log.Fatal("Failed to find student role:", err)
	}

	// Student permissions (basic reading access)
	studentPerms := []string{
		"books.view",
		"library.view",
	}

	var studentPermissions []models.Permission
	for _, permName := range studentPerms {
		for _, perm := range allPermissions {
			if perm.Name == permName {
				studentPermissions = append(studentPermissions, perm)
				break
			}
		}
	}

	if err := database.DB.Model(&studentRole).Association("Permissions").Replace(&studentPermissions); err != nil {
		log.Fatal("Failed to assign permissions to student:", err)
	}
	log.Printf("✓ Assigned %d permissions to student role", len(studentPermissions))

	log.Println("\n✅ All role permissions assigned successfully!")
}
