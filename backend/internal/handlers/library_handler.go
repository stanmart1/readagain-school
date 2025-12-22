package handlers

import (
	"path/filepath"
	"strconv"

	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type LibraryHandler struct {
	libraryService *services.LibraryService
	ereaderService *services.EReaderService
}

func NewLibraryHandler(libraryService *services.LibraryService, ereaderService *services.EReaderService) *LibraryHandler {
	return &LibraryHandler{
		libraryService: libraryService,
		ereaderService: ereaderService,
	}
}

func (h *LibraryHandler) GetLibrary(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")

	library, meta, err := h.libraryService.GetUserLibrary(userID, page, limit, search)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get library: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve library"})
	}

	return c.JSON(fiber.Map{
		"library":    library,
		"pagination": meta,
	})
}

func (h *LibraryHandler) AccessBook(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	book, err := h.ereaderService.ValidateBookAccess(userID, uint(bookID))
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	if book.FilePath == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book file not found"})
	}

	filePath := filepath.Join("./uploads", book.FilePath)
	return c.SendFile(filePath)
}

func (h *LibraryHandler) UpdateProgress(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	var input struct {
		CurrentPage int     `json:"current_page" validate:"required,gte=0"`
		TotalPages  int     `json:"total_pages" validate:"required,gt=0"`
		Progress    float64 `json:"progress" validate:"required,gte=0,lte=100"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.libraryService.UpdateProgress(userID, uint(bookID), input.CurrentPage, input.TotalPages, input.Progress); err != nil {
		utils.ErrorLogger.Printf("Failed to update progress: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d updated progress for book %d", userID, bookID)
	return c.JSON(fiber.Map{"message": "Progress updated successfully"})
}

func (h *LibraryHandler) GetStatistics(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	stats, err := h.libraryService.GetReadingStatistics(userID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get statistics: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve statistics"})
	}

	return c.JSON(fiber.Map{"statistics": stats})
}

func (h *LibraryHandler) GetDashboardStats(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	stats, err := h.libraryService.GetDashboardStats(userID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get dashboard stats: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve stats"})
	}

	return c.JSON(fiber.Map{"stats": stats})
}

func (h *LibraryHandler) GetUserAnalytics(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	analytics, err := h.libraryService.GetUserAnalytics(userID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get user analytics: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve analytics"})
	}

	return c.JSON(fiber.Map{"analytics": analytics})
}

func (h *LibraryHandler) GetBookmarks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	bookmarks, err := h.ereaderService.GetBookmarks(userID, uint(bookID))
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get bookmarks: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve bookmarks"})
	}

	return c.JSON(fiber.Map{"bookmarks": bookmarks})
}

func (h *LibraryHandler) CreateBookmark(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	var input struct {
		Page     int    `json:"page" validate:"required,gte=0"`
		Location string `json:"location"`
		Note     string `json:"note"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	bookmark, err := h.ereaderService.CreateBookmark(userID, uint(bookID), input.Page, input.Location, input.Note)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create bookmark: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d created bookmark for book %d", userID, bookID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"bookmark": bookmark})
}

func (h *LibraryHandler) DeleteBookmark(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	bookmarkID, err := strconv.ParseUint(c.Params("bookmarkId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid bookmark ID"})
	}

	if err := h.ereaderService.DeleteBookmark(uint(bookmarkID), userID); err != nil {
		utils.ErrorLogger.Printf("Failed to delete bookmark: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d deleted bookmark %d", userID, bookmarkID)
	return c.JSON(fiber.Map{"message": "Bookmark deleted successfully"})
}

func (h *LibraryHandler) GetNotes(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	notes, err := h.ereaderService.GetNotes(userID, uint(bookID))
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get notes: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve notes"})
	}

	return c.JSON(fiber.Map{"notes": notes})
}

func (h *LibraryHandler) CreateNote(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	var input struct {
		Page      int    `json:"page" validate:"required,gte=0"`
		Content   string `json:"content" validate:"required"`
		Highlight string `json:"highlight"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	note, err := h.ereaderService.CreateNote(userID, uint(bookID), input.Page, input.Content, input.Highlight)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create note: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d created note for book %d", userID, bookID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"note": note})
}

func (h *LibraryHandler) UpdateNote(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	noteID, err := strconv.ParseUint(c.Params("noteId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid note ID"})
	}

	var input struct {
		Content string `json:"content" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	note, err := h.ereaderService.UpdateNote(uint(noteID), userID, input.Content)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to update note: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d updated note %d", userID, noteID)
	return c.JSON(fiber.Map{"note": note})
}

func (h *LibraryHandler) DeleteNote(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	noteID, err := strconv.ParseUint(c.Params("noteId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid note ID"})
	}

	if err := h.ereaderService.DeleteNote(uint(noteID), userID); err != nil {
		utils.ErrorLogger.Printf("Failed to delete note: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d deleted note %d", userID, noteID)
	return c.JSON(fiber.Map{"message": "Note deleted successfully"})
}
