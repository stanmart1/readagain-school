package handlers

import (
	"strconv"

	"readagain/internal/middleware"
	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	bookService *services.BookService
}

func NewBookHandler(bookService *services.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.bookService.GetStats()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get book stats: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch book statistics"})
	}

	return c.JSON(stats)
}

func (h *BookHandler) ListBooks(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")
	categoryID, _ := strconv.ParseUint(c.Query("category_id", "0"), 10, 32)
	authorID, _ := strconv.ParseUint(c.Query("author_id", "0"), 10, 32)
	status := c.Query("status", "")
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	var isFeatured *bool
	if c.Query("is_featured") != "" {
		featured := c.Query("is_featured") == "true"
		isFeatured = &featured
	}

	filters := services.BookFilters{
		Search:     search,
		CategoryID: uint(categoryID),
		AuthorID:   uint(authorID),
		IsFeatured: isFeatured,
		Status:     status,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	books, meta, err := h.bookService.ListBooks(page, limit, filters)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve books"})
	}

	return c.JSON(fiber.Map{
		"books":      books,
		"pagination": meta,
	})
}

func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	book, err := h.bookService.GetBookByID(uint(bookID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	return c.JSON(fiber.Map{"book": book})
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	var req struct {
		Title          string `json:"title" validate:"required"`
		AuthorID       uint   `json:"author_id" validate:"required"`
		CategoryID     uint   `json:"category_id" validate:"required"`
		Description    string `json:"description"`
		ISBN           string `json:"isbn"`
		Language       string `json:"language"`
		PageCount      int    `json:"page_count"`
		Publisher      string `json:"publisher"`
		Status         string `json:"status"`
		CoverImage     string `json:"cover_image"`
		BookFile       string `json:"book_file"`
		FileSize       int64  `json:"file_size"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	book, err := h.bookService.CreateBook(
		req.AuthorID,
		req.Title,
		req.Description,
		req.ISBN,
		req.CategoryID,
		req.CoverImage,
		req.BookFile,
		req.FileSize,
		req.PageCount,
		req.Status,
	)

	if err != nil {
		utils.ErrorLogger.Printf("Failed to create book: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Created book: %s", req.Title)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"book": book})
}

func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	updates := make(map[string]interface{})

	if title := c.FormValue("title"); title != "" {
		updates["title"] = title
	}
	if description := c.FormValue("description"); description != "" {
		updates["description"] = description
	}
	if isbn := c.FormValue("isbn"); isbn != "" {
		updates["isbn"] = isbn
	}
	if categoryID := c.FormValue("category_id"); categoryID != "" {
		id, _ := strconv.ParseUint(categoryID, 10, 32)
		updates["category_id"] = uint(id)
	}
	if pageCount := c.FormValue("page_count"); pageCount != "" {
		pc, _ := strconv.Atoi(pageCount)
		updates["page_count"] = pc
	}
	if status := c.FormValue("status"); status != "" {
		updates["status"] = status
	}

	if coverImage := c.FormValue("cover_image"); coverImage != "" {
		updates["cover_image"] = coverImage
	}

	if fileURL := c.FormValue("book_file"); fileURL != "" {
		updates["file_url"] = fileURL
	}

	if fileSize := c.FormValue("file_size"); fileSize != "" {
		if size, err := strconv.ParseInt(fileSize, 10, 64); err == nil {
			updates["file_size"] = size
		}
	}

	book, err := h.bookService.UpdateBook(uint(bookID), updates)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to update book: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Updated book %d", bookID)
	return c.JSON(fiber.Map{"book": book})
}

func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	if err := h.bookService.DeleteBook(uint(bookID)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete book: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	middleware.LogAudit(c, "delete_book", "book", uint(bookID), "", "")
	utils.InfoLogger.Printf("Deleted book %d", bookID)
	return c.JSON(fiber.Map{"message": "Book deleted successfully"})
}

func (h *BookHandler) GetFeaturedBooks(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	books, err := h.bookService.GetFeaturedBooks(limit)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get featured books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve featured books"})
	}

	return c.JSON(fiber.Map{"books": books})
}

func (h *BookHandler) GetNewReleases(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	books, err := h.bookService.GetNewReleases(limit)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get new releases: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve new releases"})
	}

	return c.JSON(fiber.Map{"books": books})
}

func (h *BookHandler) GetBestsellers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	books, err := h.bookService.GetBestsellers(limit)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get bestsellers: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve bestsellers"})
	}

	return c.JSON(fiber.Map{"books": books})
}

func (h *BookHandler) ToggleFeatured(c *fiber.Ctx) error {
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	var input struct {
		IsFeatured bool `json:"is_featured"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	book, err := h.bookService.ToggleFeatured(uint(bookID), input.IsFeatured)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to toggle featured: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Toggled featured for book %d", bookID)
	return c.JSON(fiber.Map{"book": book})
}
