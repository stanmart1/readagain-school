package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"readagain/internal/config"
	"readagain/internal/database"
	"readagain/internal/handlers"
	"readagain/internal/middleware"
	"readagain/internal/services"
	"readagain/internal/utils"
)

func main() {
	utils.InitLogger()
	cfg := config.Load()

	if err := database.Connect(cfg.Database.URL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "ReadAgain API v1.0.0",
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(logger.New())
	
	allowOrigins := "*"
	if cfg.Server.Env != "development" {
		allowOrigins = "https://readagain.com"
	}
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	app.Static("/uploads", "./uploads")

	app.Get("/", handlers.GetRoot)
	app.Get("/health", handlers.GetHealth)

	emailService := services.NewEmailService(cfg.Email.ResendAPIKey, cfg.Email.FromEmail, cfg.Email.FromName, cfg.Email.AppURL)
	achievementService := services.NewAchievementService(database.DB)
	authService := services.NewAuthService(database.DB, cfg, emailService)
	userService := services.NewUserService(database.DB)
	roleService := services.NewRoleService(database.DB)
	categoryService := services.NewCategoryService(database.DB)
	authorService := services.NewAuthorService(database.DB)
	bookService := services.NewBookService(database.DB)
	storageService := services.NewStorageService("./uploads")
	cartService := services.NewCartService(database.DB)
	orderService := services.NewOrderService(database.DB, emailService, achievementService)
	bankTransferService := services.NewBankTransferService(database.DB)
	paymentService := services.NewPaymentService(cfg.Payment.PaystackSecretKey, cfg.Payment.FlutterwaveSecretKey, bankTransferService)
	libraryService := services.NewLibraryService(database.DB)
	ereaderService := services.NewEReaderService(database.DB)
	sessionService := services.NewReadingSessionService(database.DB)
	goalService := services.NewReadingGoalService(database.DB)
	blogService := services.NewBlogService(database.DB)
	faqService := services.NewFAQService(database.DB)
	testimonialService := services.NewTestimonialService(database.DB)
	contactService := services.NewContactService(database.DB)
	settingsService := services.NewSettingsService(database.DB)
	analyticsService := services.NewAnalyticsService(database.DB)
	notificationService := services.NewNotificationService(database.DB)
	auditService := services.NewAuditService(database.DB)
	reviewService := services.NewReviewService(database.DB)
	aboutService := services.NewAboutService(database.DB)
	activityService := services.NewActivityService(database.DB)
	wishlistService := services.NewWishlistService(database.DB)
	workService := services.NewWorkService(database.DB)

	achievementService.SeedAchievements()

	handlers.SetupRoutes(app, authService, userService, roleService, categoryService, authorService, bookService, storageService, cartService, orderService, paymentService, libraryService, ereaderService, sessionService, goalService, achievementService, blogService, faqService, testimonialService, contactService, settingsService, analyticsService, notificationService, auditService, reviewService, aboutService, activityService, wishlistService, workService)

	utils.InfoLogger.Printf("🚀 Server starting on port %s", cfg.Server.Port)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
