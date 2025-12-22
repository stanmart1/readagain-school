package handlers

import (
	"strconv"

	"readagain/internal/middleware"
	"readagain/internal/models"
	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	roleService *services.RoleService
}

func NewRoleHandler(roleService *services.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) ListRoles(c *fiber.Ctx) error {
	roles, err := h.roleService.ListRoles()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list roles: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve roles"})
	}

	return c.JSON(fiber.Map{"roles": roles})
}

func (h *RoleHandler) GetRole(c *fiber.Ctx) error {
	roleID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role ID"})
	}

	role, err := h.roleService.GetRoleByID(uint(roleID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Role not found"})
	}

	return c.JSON(fiber.Map{"role": role})
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var input struct {
		Name          string `json:"name" validate:"required,min=3"`
		Description   string `json:"description"`
		PermissionIDs []uint `json:"permission_ids"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	role, err := h.roleService.CreateRole(input.Name, input.Description, input.PermissionIDs)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create role: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Admin created role: %s", input.Name)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"role": role})
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	roleID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role ID"})
	}

	var input struct {
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		PermissionIDs *[]uint `json:"permission_ids"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var permissionIDs []uint
	if input.PermissionIDs != nil {
		permissionIDs = *input.PermissionIDs
	}

	role, err := h.roleService.UpdateRole(uint(roleID), input.Name, input.Description, permissionIDs)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to update role %d: %v", roleID, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Admin updated role %d", roleID)
	return c.JSON(fiber.Map{"role": role})
}

func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	roleID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role ID"})
	}

	if err := h.roleService.DeleteRole(uint(roleID)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete role %d: %v", roleID, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	middleware.LogAudit(c, "delete_role", "role", uint(roleID), "", "")
	utils.InfoLogger.Printf("Admin deleted role %d", roleID)
	return c.JSON(fiber.Map{"message": "Role deleted successfully"})
}

func (h *RoleHandler) ListPermissions(c *fiber.Ctx) error {
	permissions, err := h.roleService.ListPermissions()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list permissions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve permissions"})
	}

	return c.JSON(fiber.Map{"permissions": permissions})
}

func (h *RoleHandler) GetUserPermissions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user models.User
	if err := h.roleService.GetDB().Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user permissions"})
	}

	permissions := []string{}
	if user.Role != nil {
		for _, perm := range user.Role.Permissions {
			permissions = append(permissions, perm.Name)
		}
	}

	return c.JSON(fiber.Map{"permissions": permissions})
}

func (h *RoleHandler) GetRolePermissions(c *fiber.Ctx) error {
	roleID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role ID"})
	}

	var role models.Role
	if err := h.roleService.GetDB().Preload("Permissions").First(&role, roleID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Role not found"})
	}

	return c.JSON(fiber.Map{"permissions": role.Permissions})
}

func (h *RoleHandler) AddPermission(c *fiber.Ctx) error {
	roleID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role ID"})
	}

	var req struct {
		PermissionID uint `json:"permission_id" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.roleService.AddPermission(uint(roleID), req.PermissionID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	middleware.LogAudit(c, "add_permission_to_role", "role", uint(roleID), "", "")
	return c.JSON(fiber.Map{"message": "Permission added successfully"})
}

func (h *RoleHandler) RemovePermission(c *fiber.Ctx) error {
	roleID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role ID"})
	}

	permissionID, err := strconv.ParseUint(c.Params("permissionId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid permission ID"})
	}

	if err := h.roleService.RemovePermission(uint(roleID), uint(permissionID)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	middleware.LogAudit(c, "remove_permission_from_role", "role", uint(roleID), "", "")
	return c.JSON(fiber.Map{"message": "Permission removed successfully"})
}
