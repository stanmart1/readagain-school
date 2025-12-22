package handlers

import (
	"strconv"

	"readagain/internal/middleware"
	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	bookService    *services.BookService
	storageService *services.StorageService
}

func NewBookHandler(bookService *services.BookService, storageService *services.StorageService) *BookHandler {
	return &BookHandler{
		bookService:    bookService,
		storageService: storageService,
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
	minPrice, _ := strconv.ParseFloat(c.Query("min_price", "0"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price", "0"), 64)
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
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
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
	authorID, _ := strconv.ParseUint(c.FormValue("author_id"), 10, 32)
	categoryID, _ := strconv.ParseUint(c.FormValue("category_id"), 10, 32)
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
	pageCount, _ := strconv.Atoi(c.FormValue("page_count"))

	title := c.FormValue("title")
	description := c.FormValue("description")
	isbn := c.FormValue("isbn")
	status := c.FormValue("status")

	if title == "" || authorID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title and author_id are required"})
	}

	var coverImage, fileURL string
	var fileSize int64

	coverFile, err := c.FormFile("cover_image")
	if err == nil {
		coverImage, err = h.storageService.UploadBookCover(coverFile)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	bookFile, err := c.FormFile("book_file")
	if err == nil {
		fileURL, err = h.storageService.UploadBookFile(bookFile)
		if err != nil {
			if coverImage != "" {
				h.storageService.DeleteFile(coverImage)
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		fileSize = bookFile.Size
	}

	book, err := h.bookService.CreateBook(
		uint(authorID),
		title,
		description,
		isbn,
		uint(categoryID),
		price,
		coverImage,
		fileURL,
		fileSize,
		pageCount,
		status,
	)

	if err != nil {
		if coverImage != "" {
			h.storageService.DeleteFile(coverImage)
		}
		if fileURL != "" {
			h.storageService.DeleteFile(fileURL)
		}
		utils.ErrorLogger.Printf("Failed to create book: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Created book: %s", title)
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
	if price := c.FormValue("price"); price != "" {
		p, _ := strconv.ParseFloat(price, 64)
		updates["price"] = p
	}
	if pageCount := c.FormValue("page_count"); pageCount != "" {
		pc, _ := strconv.Atoi(pageCount)
		updates["page_count"] = pc
	}
	if status := c.FormValue("status"); status != "" {
		updates["status"] = status
	}

	coverFile, err := c.FormFile("cover_image")
	if err == nil {
		coverImage, err := h.storageService.UploadBookCover(coverFile)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		updates["cover_image"] = coverImage
	}

	bookFile, err := c.FormFile("book_file")
	if err == nil {
		fileURL, err := h.storageService.UploadBookFile(bookFile)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		updates["file_url"] = fileURL
		updates["file_size"] = bookFile.Size
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
