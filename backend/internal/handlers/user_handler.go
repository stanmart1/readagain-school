package handlers

import (
	"fmt"
	"strconv"

	"readagain/internal/middleware"
	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get user profile: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"user": user})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		FirstName   *string `json:"first_name"`
		LastName    *string `json:"last_name"`
		PhoneNumber *string `json:"phone_number"`
		Bio         *string `json:"bio"`
		AvatarURL   *string `json:"avatar_url"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	updates := make(map[string]interface{})
	if input.FirstName != nil {
		updates["first_name"] = *input.FirstName
	}
	if input.LastName != nil {
		updates["last_name"] = *input.LastName
	}
	if input.PhoneNumber != nil {
		updates["phone_number"] = *input.PhoneNumber
	}
	if input.Bio != nil {
		updates["bio"] = *input.Bio
	}
	if input.AvatarURL != nil {
		updates["avatar_url"] = *input.AvatarURL
	}

	user, err := h.userService.UpdateUser(userID, updates)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to update user profile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
	}

	utils.InfoLogger.Printf("User %d updated profile", userID)
	return c.JSON(fiber.Map{"user": user})
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.userService.ChangePassword(userID, input.CurrentPassword, input.NewPassword); err != nil {
		utils.ErrorLogger.Printf("Failed to change password for user %d: %v", userID, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d changed password", userID)
	return c.JSON(fiber.Map{"message": "Password changed successfully"})
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")

	users, meta, err := h.userService.ListUsers(page, limit, search)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}

	return c.JSON(fiber.Map{
		"users":      users,
		"pagination": meta,
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.userService.GetUserByID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"user": user})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := h.userService.UpdateUser(uint(userID), input)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to update user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	utils.InfoLogger.Printf("Admin updated user %d", userID)
	return c.JSON(fiber.Map{"user": user})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	if err := h.userService.DeleteUser(uint(userID)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	middleware.LogAudit(c, "delete_user", "user", uint(userID), "", "")
	utils.InfoLogger.Printf("Admin deleted user %d", userID)
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

func (h *UserHandler) AssignRole(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var input struct {
		RoleID uint `json:"role_id" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.userService.AssignRole(uint(userID), input.RoleID); err != nil {
		utils.ErrorLogger.Printf("Failed to assign role to user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to assign role"})
	}

	utils.InfoLogger.Printf("Admin assigned role %d to user %d", input.RoleID, userID)
	return c.JSON(fiber.Map{"message": "Role assigned successfully"})
}

func (h *UserHandler) RemoveRole(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var input struct {
		RoleID uint `json:"role_id" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.userService.RemoveRole(uint(userID), input.RoleID); err != nil {
		utils.ErrorLogger.Printf("Failed to remove role from user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove role"})
	}

	utils.InfoLogger.Printf("Admin removed role %d from user %d", input.RoleID, userID)
	return c.JSON(fiber.Map{"message": "Role removed successfully"})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var input struct {
		Email     string `json:"email" validate:"required,email"`
		Username  string `json:"username" validate:"required,min=3"`
		Password  string `json:"password" validate:"required,min=8"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		RoleID    uint   `json:"role_id"`
		IsActive  *bool  `json:"is_active"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	roleID := input.RoleID
	if roleID == 0 {
		roleID = 3
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	user, err := h.userService.CreateUserByAdmin(
		input.Email,
		input.Username,
		input.Password,
		input.FirstName,
		input.LastName,
		roleID,
		isActive,
	)

	if err != nil {
		utils.ErrorLogger.Printf("Failed to create user: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Admin created user: %s", input.Email)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"user": user})
}

func (h *UserHandler) ToggleStatus(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var input struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := h.userService.ToggleUserStatus(uint(userID), input.IsActive)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to toggle user status: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update status"})
	}

	utils.InfoLogger.Printf("Admin toggled user %d status to %v", userID, input.IsActive)
	return c.JSON(fiber.Map{"user": user})
}

func (h *UserHandler) AdminResetPassword(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var input struct {
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.userService.AdminResetPassword(uint(userID), input.NewPassword); err != nil {
		utils.ErrorLogger.Printf("Failed to reset password for user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset password"})
	}

	utils.InfoLogger.Printf("Admin reset password for user %d", userID)
	return c.JSON(fiber.Map{"message": "Password reset successfully"})
}

func (h *UserHandler) BulkActivate(c *fiber.Ctx) error {
	var input struct {
		UserIDs []uint `json:"user_ids" validate:"required,min=1"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.userService.BulkActivate(input.UserIDs); err != nil {
		utils.ErrorLogger.Printf("Failed to bulk activate users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Admin activated %d users", len(input.UserIDs))
	return c.JSON(fiber.Map{"message": "Users activated successfully", "count": len(input.UserIDs)})
}

func (h *UserHandler) BulkDeactivate(c *fiber.Ctx) error {
	var input struct {
		UserIDs []uint `json:"user_ids" validate:"required,min=1"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.userService.BulkDeactivate(input.UserIDs); err != nil {
		utils.ErrorLogger.Printf("Failed to bulk deactivate users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Admin deactivated %d users", len(input.UserIDs))
	return c.JSON(fiber.Map{"message": "Users deactivated successfully", "count": len(input.UserIDs)})
}

func (h *UserHandler) BulkDelete(c *fiber.Ctx) error {
	var input struct {
		UserIDs []uint `json:"user_ids" validate:"required,min=1"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.userService.BulkDelete(input.UserIDs); err != nil {
		utils.ErrorLogger.Printf("Failed to bulk delete users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	middleware.LogAudit(c, "bulk_delete_users", "user", 0, "", fmt.Sprintf("%v", input.UserIDs))
	utils.InfoLogger.Printf("Admin deleted %d users", len(input.UserIDs))
	return c.JSON(fiber.Map{"message": "Users deleted successfully", "count": len(input.UserIDs)})
}


