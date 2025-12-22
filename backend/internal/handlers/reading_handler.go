package handlers

import (
	"strconv"
	"time"

	"readagain/internal/models"
	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type ReadingHandler struct {
	sessionService *services.ReadingSessionService
	goalService    *services.ReadingGoalService
}

func NewReadingHandler(sessionService *services.ReadingSessionService, goalService *services.ReadingGoalService) *ReadingHandler {
	return &ReadingHandler{
		sessionService: sessionService,
		goalService:    goalService,
	}
}

func (h *ReadingHandler) StartSession(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		BookID uint `json:"book_id" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	session, err := h.sessionService.StartSession(userID, input.BookID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to start session: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d started reading session for book %d", userID, input.BookID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"session": session})
}

func (h *ReadingHandler) EndSession(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		SessionID uint `json:"session_id" validate:"required"`
		Duration  int  `json:"duration" validate:"required,gte=0"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.sessionService.EndSession(input.SessionID, input.Duration); err != nil {
		utils.ErrorLogger.Printf("Failed to end session: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d ended reading session %d", userID, input.SessionID)
	return c.JSON(fiber.Map{"message": "Session ended successfully"})
}

func (h *ReadingHandler) GetSessions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	sessions, meta, err := h.sessionService.GetUserSessions(userID, page, limit)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get sessions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve sessions"})
	}

	return c.JSON(fiber.Map{
		"sessions":   sessions,
		"pagination": meta,
	})
}

func (h *ReadingHandler) GetReadingProgress(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var libraries []models.UserLibrary
	err := h.sessionService.GetDB().
		Where("user_id = ? AND progress > 0 AND progress < 100", userID).
		Preload("Book").
		Order("last_read_at DESC").
		Limit(10).
		Find(&libraries).Error

	if err != nil {
		utils.ErrorLogger.Printf("Failed to get reading progress: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve reading progress"})
	}

	progress := make([]map[string]interface{}, 0)
	for _, lib := range libraries {
		if lib.Book != nil {
			progress = append(progress, map[string]interface{}{
				"id":          lib.ID,
				"book_id":     lib.BookID,
				"title":       lib.Book.Title,
				"author":      lib.Book.Author,
				"cover_image": lib.Book.CoverImage,
				"progress":    lib.Progress,
			})
		}
	}

	return c.JSON(fiber.Map{"reading_progress": progress})
}

func (h *ReadingHandler) GetGoals(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	goals, err := h.goalService.GetUserGoals(userID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get goals: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve goals"})
	}

	return c.JSON(fiber.Map{"goals": goals})
}

func (h *ReadingHandler) CreateGoal(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		GoalType    string `json:"goal_type" validate:"required,oneof=books pages minutes"`
		TargetValue int    `json:"target_value" validate:"required,gt=0"`
		StartDate   string `json:"start_date" validate:"required"`
		EndDate     string `json:"end_date" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start date format"})
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid end date format"})
	}

	goal, err := h.goalService.CreateGoal(userID, input.GoalType, input.TargetValue, startDate, endDate)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create goal: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d created reading goal", userID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"goal": goal})
}

func (h *ReadingHandler) UpdateGoal(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	goalID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid goal ID"})
	}

	var input struct {
		TargetValue int `json:"target_value" validate:"required,gt=0"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	goal, err := h.goalService.UpdateGoal(uint(goalID), userID, input.TargetValue)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to update goal: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d updated goal %d", userID, goalID)
	return c.JSON(fiber.Map{"goal": goal})
}

func (h *ReadingHandler) DeleteGoal(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	goalID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid goal ID"})
	}

	if err := h.goalService.DeleteGoal(uint(goalID), userID); err != nil {
		utils.ErrorLogger.Printf("Failed to delete goal: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d deleted goal %d", userID, goalID)
	return c.JSON(fiber.Map{"message": "Goal deleted successfully"})
}
