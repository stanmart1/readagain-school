package handlers

import (
	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AchievementHandler struct {
	achievementService *services.AchievementService
}

func NewAchievementHandler(achievementService *services.AchievementService) *AchievementHandler {
	return &AchievementHandler{achievementService: achievementService}
}

func (h *AchievementHandler) GetAllAchievements(c *fiber.Ctx) error {
	achievements, err := h.achievementService.GetAllAchievements()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get achievements: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve achievements"})
	}

	return c.JSON(fiber.Map{"achievements": achievements})
}

func (h *AchievementHandler) GetUserAchievements(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	userAchievements, err := h.achievementService.GetUserAchievements(userID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get user achievements: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user achievements"})
	}

	return c.JSON(fiber.Map{"achievements": userAchievements})
}

func (h *AchievementHandler) CheckAchievements(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	newAchievements, err := h.achievementService.CheckAndUnlockAchievements(userID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to check achievements: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check achievements"})
	}

	return c.JSON(fiber.Map{
		"new_achievements": newAchievements,
		"count":            len(newAchievements),
	})
}

func (h *AchievementHandler) CreateAchievement(c *fiber.Ctx) error {
	var input struct {
		Name        string `json:"name" validate:"required,min=3"`
		Description string `json:"description" validate:"required"`
		Icon        string `json:"icon"`
		Type        string `json:"type" validate:"required,oneof=books_purchased books_completed reading_minutes reading_sessions"`
		Target      int    `json:"target" validate:"required,gt=0"`
		Points      int    `json:"points" validate:"required,gte=0"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	achievement, err := h.achievementService.CreateAchievement(
		input.Name,
		input.Description,
		input.Icon,
		input.Type,
		input.Target,
		input.Points,
	)

	if err != nil {
		utils.ErrorLogger.Printf("Failed to create achievement: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Admin created achievement: %s", input.Name)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"achievement": achievement})
}

func (h *AchievementHandler) UpdateAchievement(c *fiber.Ctx) error {
	achievementID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid achievement ID"})
	}

	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	achievement, err := h.achievementService.UpdateAchievement(uint(achievementID), input)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to update achievement: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Admin updated achievement %d", achievementID)
	return c.JSON(fiber.Map{"achievement": achievement})
}

func (h *AchievementHandler) DeleteAchievement(c *fiber.Ctx) error {
	achievementID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid achievement ID"})
	}

	if err := h.achievementService.DeleteAchievement(uint(achievementID)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete achievement: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Admin deleted achievement %d", achievementID)
	return c.JSON(fiber.Map{"message": "Achievement deleted successfully"})
}
