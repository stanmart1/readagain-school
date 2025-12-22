package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

func (s *RoleService) GetDB() *gorm.DB {
	return s.db
}

func (s *RoleService) ListRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := s.db.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch roles", err)
	}
	return roles, nil
}

func (s *RoleService) GetRoleByID(roleID uint) (*models.Role, error) {
	var role models.Role
	if err := s.db.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, utils.NewNotFoundError("Role not found")
	}
	return &role, nil
}

func (s *RoleService) CreateRole(name, description string, permissionIDs []uint) (*models.Role, error) {
	var existingRole models.Role
	if err := s.db.Where("name = ?", name).First(&existingRole).Error; err == nil {
		return nil, utils.NewBadRequestError("Role with this name already exists")
	}

	role := models.Role{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(&role).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create role", err)
	}

	if len(permissionIDs) > 0 {
		var permissions []models.Permission
		if err := s.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return nil, utils.NewInternalServerError("Failed to fetch permissions", err)
		}

		if err := s.db.Model(&role).Association("Permissions").Append(permissions); err != nil {
			return nil, utils.NewInternalServerError("Failed to assign permissions", err)
		}
	}

	if err := s.db.Preload("Permissions").First(&role, role.ID).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch created role", err)
	}

	return &role, nil
}

func (s *RoleService) UpdateRole(roleID uint, name, description string, permissionIDs []uint) (*models.Role, error) {
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return nil, utils.NewNotFoundError("Role not found")
	}

	if name != "" && name != role.Name {
		var existingRole models.Role
		if err := s.db.Where("name = ? AND id != ?", name, roleID).First(&existingRole).Error; err == nil {
			return nil, utils.NewBadRequestError("Role with this name already exists")
		}
		role.Name = name
	}

	if description != "" {
		role.Description = description
	}

	if err := s.db.Save(&role).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update role", err)
	}

	if permissionIDs != nil {
		var permissions []models.Permission
		if err := s.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return nil, utils.NewInternalServerError("Failed to fetch permissions", err)
		}

		if err := s.db.Model(&role).Association("Permissions").Replace(permissions); err != nil {
			return nil, utils.NewInternalServerError("Failed to update permissions", err)
		}
	}

	if err := s.db.Preload("Permissions").First(&role, role.ID).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch updated role", err)
	}

	return &role, nil
}

func (s *RoleService) DeleteRole(roleID uint) error {
	if roleID <= 3 {
		return utils.NewBadRequestError("Cannot delete system roles (Admin, Author, User)")
	}

	var usersCount int64
	if err := s.db.Model(&models.User{}).Where("role_id = ?", roleID).Count(&usersCount).Error; err != nil {
		return utils.NewInternalServerError("Failed to check role usage", err)
	}

	if usersCount > 0 {
		return utils.NewBadRequestError("Cannot delete role that is assigned to users")
	}

	if err := s.db.Delete(&models.Role{}, roleID).Error; err != nil {
		return utils.NewInternalServerError("Failed to delete role", err)
	}

	return nil
}

func (s *RoleService) ListPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.Find(&permissions).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch permissions", err)
	}
	return permissions, nil
}

func (s *RoleService) AddPermission(roleID, permissionID uint) error {
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return utils.NewNotFoundError("Role not found")
	}

	var permission models.Permission
	if err := s.db.First(&permission, permissionID).Error; err != nil {
		return utils.NewNotFoundError("Permission not found")
	}

	if err := s.db.Model(&role).Association("Permissions").Append(&permission); err != nil {
		return utils.NewInternalServerError("Failed to add permission", err)
	}

	return nil
}

func (s *RoleService) RemovePermission(roleID, permissionID uint) error {
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return utils.NewNotFoundError("Role not found")
	}

	var permission models.Permission
	if err := s.db.First(&permission, permissionID).Error; err != nil {
		return utils.NewNotFoundError("Permission not found")
	}

	if err := s.db.Model(&role).Association("Permissions").Delete(&permission); err != nil {
		return utils.NewInternalServerError("Failed to remove permission", err)
	}

	return nil
}
