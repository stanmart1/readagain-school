package handlers

import (
	"readagain/internal/middleware"
	"readagain/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	authService *services.AuthService,
	userService *services.UserService,
	roleService *services.RoleService,
	categoryService *services.CategoryService,
	authorService *services.AuthorService,
	bookService *services.BookService,
	storageService *services.StorageService,
	cartService *services.CartService,
	orderService *services.OrderService,
	paymentService *services.PaymentService,
	libraryService *services.LibraryService,
	ereaderService *services.EReaderService,
	sessionService *services.ReadingSessionService,
	goalService *services.ReadingGoalService,
	achievementService *services.AchievementService,
	blogService *services.BlogService,
	faqService *services.FAQService,
	testimonialService *services.TestimonialService,
	contactService *services.ContactService,
	settingsService *services.SettingsService,
	analyticsService *services.AnalyticsService,
	notificationService *services.NotificationService,
	auditService *services.AuditService,
	reviewService *services.ReviewService,
	aboutService *services.AboutService,
	activityService *services.ActivityService,
	wishlistService *services.WishlistService,
	workService *services.WorkService,
) {
	api := app.Group("/api/v1")

	authHandler := NewAuthHandler(authService)
	userHandler := NewUserHandler(userService)
	roleHandler := NewRoleHandler(roleService)
	categoryHandler := NewCategoryHandler(categoryService)
	authorHandler := NewAuthorHandler(authorService)
	bookHandler := NewBookHandler(bookService, storageService)
	cartHandler := NewCartHandler(cartService)
	checkoutHandler := NewCheckoutHandler(orderService, paymentService)
	orderHandler := NewOrderHandler(orderService)
	libraryHandler := NewLibraryHandler(libraryService, ereaderService)
	readingHandler := NewReadingHandler(sessionService, goalService)
	achievementHandler := NewAchievementHandler(achievementService)
	blogHandler := NewBlogHandler(blogService)
	faqHandler := NewFAQHandler(faqService)
	testimonialHandler := NewTestimonialHandler(testimonialService)
	contactHandler := NewContactHandler(contactService)
	settingsHandler := NewSettingsHandler(settingsService)
	analyticsHandler := NewAnalyticsHandler(analyticsService)
	notificationHandler := NewNotificationHandler(notificationService)
	auditHandler := NewAuditHandler(auditService)
	reviewHandler := NewReviewHandler(reviewService)
	aboutHandler := NewAboutHandler(aboutService)
	activityHandler := NewActivityHandler(activityService)
	wishlistHandler := NewWishlistHandler(wishlistService)
	workHandler := NewWorkHandler(workService)

	app.Use(middleware.AuditMiddleware(auditService))

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", middleware.AuthRequired(), authHandler.Logout)
	auth.Get("/me", middleware.AuthRequired(), authHandler.GetMe)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)

	users := api.Group("/users")
	users.Get("/profile", middleware.AuthRequired(), userHandler.GetProfile)
	users.Put("/profile", middleware.AuthRequired(), userHandler.UpdateProfile)
	users.Post("/change-password", middleware.AuthRequired(), userHandler.ChangePassword)

	users.Post("/", middleware.AdminRequired(), userHandler.CreateUser)
	users.Get("/", middleware.AdminRequired(), userHandler.ListUsers)
	users.Get("/:id", middleware.AdminRequired(), userHandler.GetUser)
	users.Put("/:id", middleware.AdminRequired(), userHandler.UpdateUser)
	users.Patch("/:id/status", middleware.AdminRequired(), userHandler.ToggleStatus)
	users.Delete("/:id", middleware.AdminRequired(), userHandler.DeleteUser)
	users.Post("/:id/roles", middleware.AdminRequired(), userHandler.AssignRole)
	users.Delete("/:id/roles", middleware.AdminRequired(), userHandler.RemoveRole)
	users.Post("/:id/reset-password", middleware.AdminRequired(), userHandler.AdminResetPassword)
	
	users.Post("/bulk/activate", middleware.AdminRequired(), userHandler.BulkActivate)
	users.Post("/bulk/deactivate", middleware.AdminRequired(), userHandler.BulkDeactivate)
	users.Post("/bulk/delete", middleware.AdminRequired(), userHandler.BulkDelete)

	roles := api.Group("/roles")
	roles.Get("/", roleHandler.ListRoles)
	roles.Post("/", middleware.AdminRequired(), roleHandler.CreateRole)
	roles.Get("/:id", roleHandler.GetRole)
	roles.Put("/:id", middleware.AdminRequired(), roleHandler.UpdateRole)
	roles.Delete("/:id", middleware.AdminRequired(), roleHandler.DeleteRole)

	permissions := api.Group("/permissions")
	permissions.Get("/", roleHandler.ListPermissions)

	// Auth permissions alias
	api.Get("/auth/permissions", middleware.AuthRequired(), roleHandler.GetUserPermissions)

	// RBAC aliases for frontend compatibility
	rbac := api.Group("/rbac")
	rbac.Get("/roles", roleHandler.ListRoles)
	rbac.Post("/roles", middleware.AdminRequired(), roleHandler.CreateRole)
	rbac.Get("/roles/:id", roleHandler.GetRole)
	rbac.Put("/roles/:id", middleware.AdminRequired(), roleHandler.UpdateRole)
	rbac.Delete("/roles/:id", middleware.AdminRequired(), roleHandler.DeleteRole)
	rbac.Get("/roles/:id/permissions", roleHandler.GetRolePermissions)
	rbac.Post("/roles/:id/permissions", middleware.AdminRequired(), roleHandler.AddPermission)
	rbac.Delete("/roles/:id/permissions/:permissionId", middleware.AdminRequired(), roleHandler.RemovePermission)
	rbac.Get("/permissions", roleHandler.ListPermissions)

	// Admin users alias
	adminUsers := api.Group("/admin/users", middleware.AdminRequired())
	adminUsers.Get("/", userHandler.ListUsers)
	adminUsers.Post("/", userHandler.CreateUser)
	adminUsers.Get("/:id", userHandler.GetUser)
	adminUsers.Put("/:id", userHandler.UpdateUser)
	adminUsers.Put("/:id/status", userHandler.ToggleStatus)
	adminUsers.Delete("/:id", userHandler.DeleteUser)

	categories := api.Group("/categories")
	categories.Get("/", categoryHandler.ListCategories)
	categories.Get("/:id", categoryHandler.GetCategory)
	categories.Post("/", middleware.AdminRequired(), categoryHandler.CreateCategory)
	categories.Put("/:id", middleware.AdminRequired(), categoryHandler.UpdateCategory)
	categories.Delete("/:id", middleware.AdminRequired(), categoryHandler.DeleteCategory)

	authors := api.Group("/authors")
	authors.Get("/", authorHandler.ListAuthors)
	authors.Get("/:id", authorHandler.GetAuthor)
	authors.Post("/", middleware.AdminRequired(), authorHandler.CreateAuthor)
	authors.Put("/:id", middleware.AdminRequired(), authorHandler.UpdateAuthor)
	authors.Delete("/:id", middleware.AdminRequired(), authorHandler.DeleteAuthor)

	// Admin authors routes
	adminAuthors := api.Group("/admin/authors", middleware.AdminRequired())
	adminAuthors.Get("/", authorHandler.ListAuthors)
	adminAuthors.Get("/stats", authorHandler.GetStats)
	adminAuthors.Post("/", authorHandler.CreateAuthor)
	adminAuthors.Get("/:id", authorHandler.GetAuthor)
	adminAuthors.Put("/:id", authorHandler.UpdateAuthor)
	adminAuthors.Delete("/:id", authorHandler.DeleteAuthor)

	books := api.Group("/books")
	books.Get("/", bookHandler.ListBooks)
	books.Get("/featured", bookHandler.GetFeaturedBooks)
	books.Get("/new-releases", bookHandler.GetNewReleases)
	books.Get("/bestsellers", bookHandler.GetBestsellers)
	books.Get("/:id", bookHandler.GetBook)
	books.Post("/", middleware.AdminRequired(), bookHandler.CreateBook)
	books.Put("/:id", middleware.AdminRequired(), bookHandler.UpdateBook)
	books.Delete("/:id", middleware.AdminRequired(), bookHandler.DeleteBook)
	books.Patch("/:id/featured", middleware.AdminRequired(), bookHandler.ToggleFeatured)

	// Admin books routes
	adminBooks := api.Group("/admin/books", middleware.AdminRequired())
	adminBooks.Get("/", bookHandler.ListBooks)
	adminBooks.Get("/stats", bookHandler.GetStats)
	adminBooks.Post("/", bookHandler.CreateBook)
	adminBooks.Get("/:id", bookHandler.GetBook)
	adminBooks.Put("/:id", bookHandler.UpdateBook)
	adminBooks.Delete("/:id", bookHandler.DeleteBook)
	adminBooks.Patch("/:id/featured", bookHandler.ToggleFeatured)

	cart := api.Group("/cart", middleware.AuthRequired())
	cart.Get("/", cartHandler.GetCart)
	cart.Get("/count", cartHandler.GetCartCount)
	cart.Post("/", cartHandler.AddToCart)
	cart.Delete("/:id", cartHandler.RemoveFromCart)
	cart.Delete("/", cartHandler.ClearCart)
	cart.Post("/merge", cartHandler.MergeGuestCart)

	checkout := api.Group("/checkout", middleware.AuthRequired())
	checkout.Post("/initialize", checkoutHandler.InitializeCheckout)
	checkout.Post("/payment", checkoutHandler.InitializePayment)
	checkout.Get("/verify/:reference", checkoutHandler.VerifyPayment)

	webhooks := api.Group("/webhooks")
	webhooks.Post("/paystack", checkoutHandler.PaystackWebhook)
	webhooks.Post("/flutterwave", checkoutHandler.FlutterwaveWebhook)

	orders := api.Group("/orders", middleware.AuthRequired())
	orders.Get("/", orderHandler.GetUserOrders)
	orders.Get("/:id", orderHandler.GetOrder)
	orders.Post("/:id/cancel", orderHandler.CancelOrder)

	adminOrders := api.Group("/admin/orders", middleware.AdminRequired())
	adminOrders.Get("/", orderHandler.GetAllOrders)
	adminOrders.Get("/stats", orderHandler.GetOrderStatistics)
	adminOrders.Get("/:id", orderHandler.GetOrderAdmin)
	adminOrders.Patch("/:id/status", orderHandler.UpdateOrderStatus)

	library := api.Group("/library", middleware.AuthRequired())
	library.Get("/", libraryHandler.GetLibrary)
	library.Get("/statistics", libraryHandler.GetStatistics)
	library.Get("/:id/access", libraryHandler.AccessBook)
	library.Put("/:id/progress", libraryHandler.UpdateProgress)
	library.Get("/:id/bookmarks", libraryHandler.GetBookmarks)
	library.Post("/:id/bookmarks", libraryHandler.CreateBookmark)
	library.Delete("/bookmarks/:bookmarkId", libraryHandler.DeleteBookmark)
	library.Get("/:id/notes", libraryHandler.GetNotes)
	library.Post("/:id/notes", libraryHandler.CreateNote)
	library.Put("/notes/:noteId", libraryHandler.UpdateNote)
	library.Delete("/notes/:noteId", libraryHandler.DeleteNote)

	reading := api.Group("/reading", middleware.AuthRequired())
	reading.Post("/sessions/start", readingHandler.StartSession)
	reading.Post("/sessions/end", readingHandler.EndSession)
	reading.Get("/sessions", readingHandler.GetSessions)
	reading.Get("/goals", readingHandler.GetGoals)
	reading.Post("/goals", readingHandler.CreateGoal)
	reading.Put("/goals/:id", readingHandler.UpdateGoal)
	reading.Delete("/goals/:id", readingHandler.DeleteGoal)

	readingGoals := api.Group("/reading-goals", middleware.AuthRequired())
	readingGoals.Get("/", readingHandler.GetGoals)
	readingGoals.Post("/", readingHandler.CreateGoal)
	readingGoals.Put("/:id", readingHandler.UpdateGoal)
	readingGoals.Delete("/:id", readingHandler.DeleteGoal)

	achievements := api.Group("/achievements")
	achievements.Get("/", achievementHandler.GetAllAchievements)
	achievements.Get("/user", middleware.AuthRequired(), achievementHandler.GetUserAchievements)
	achievements.Post("/check", middleware.AuthRequired(), achievementHandler.CheckAchievements)
	achievements.Post("/", middleware.AdminRequired(), achievementHandler.CreateAchievement)
	achievements.Put("/:id", middleware.AdminRequired(), achievementHandler.UpdateAchievement)
	achievements.Delete("/:id", middleware.AdminRequired(), achievementHandler.DeleteAchievement)

	blogs := api.Group("/blogs")
	blogs.Get("/", blogHandler.List)
	blogs.Get("/:slug", blogHandler.GetBySlug)

	adminBlogs := api.Group("/admin/blogs", middleware.AdminRequired())
	adminBlogs.Get("/", blogHandler.AdminList)
	adminBlogs.Get("/stats", blogHandler.GetStats)
	adminBlogs.Get("/:id", blogHandler.GetByID)
	adminBlogs.Post("/", blogHandler.Create)
	adminBlogs.Put("/:id", blogHandler.Update)
	adminBlogs.Delete("/:id", blogHandler.Delete)

	faqs := api.Group("/faqs")
	faqs.Get("/", faqHandler.List)
	faqs.Get("/categories", faqHandler.GetCategories)

	adminFaqs := api.Group("/admin/faqs", middleware.AdminRequired())
	adminFaqs.Get("/", faqHandler.AdminList)
	adminFaqs.Get("/:id", faqHandler.GetByID)
	adminFaqs.Post("/", faqHandler.Create)
	adminFaqs.Put("/:id", faqHandler.Update)
	adminFaqs.Delete("/:id", faqHandler.Delete)

	testimonials := api.Group("/testimonials")
	testimonials.Get("/", testimonialHandler.List)

	adminTestimonials := api.Group("/admin/testimonials", middleware.AdminRequired())
	adminTestimonials.Get("/", testimonialHandler.AdminList)
	adminTestimonials.Get("/:id", testimonialHandler.GetByID)
	adminTestimonials.Post("/", testimonialHandler.Create)
	adminTestimonials.Put("/:id", testimonialHandler.Update)
	adminTestimonials.Delete("/:id", testimonialHandler.Delete)

	contact := api.Group("/contact")
	contact.Post("/", contactHandler.Submit)

	adminContact := api.Group("/admin/contact", middleware.AdminRequired())
	adminContact.Get("/", contactHandler.List)
	adminContact.Get("/:id", contactHandler.GetByID)
	adminContact.Post("/:id/reply", contactHandler.Reply)
	adminContact.Patch("/:id/status", contactHandler.UpdateStatus)
	adminContact.Delete("/:id", contactHandler.Delete)

	settings := api.Group("/admin/settings", middleware.AdminRequired())
	settings.Get("/", settingsHandler.GetByCategory)
	settings.Get("/:key", settingsHandler.GetByKey)
	settings.Post("/", settingsHandler.Set)
	settings.Delete("/:key", settingsHandler.Delete)
	settings.Get("/email/config", settingsHandler.GetEmailSettings)
	settings.Put("/email/config", settingsHandler.UpdateEmailSettings)
	settings.Get("/payment/config", settingsHandler.GetPaymentSettings)
	settings.Put("/payment/config", settingsHandler.UpdatePaymentSettings)

	// Email gateways endpoint
	api.Get("/admin/email/gateways", middleware.AdminRequired(), settingsHandler.GetEmailSettings)
	api.Post("/admin/email/gateways/test", middleware.AdminRequired(), settingsHandler.TestEmailGateway)
	
	// Payment settings endpoint  
	api.Get("/admin/payment-settings", middleware.AdminRequired(), settingsHandler.GetPaymentSettings)

	analytics := api.Group("/admin/analytics", middleware.AdminRequired())
	analytics.Get("/dashboard", analyticsHandler.GetEnhancedOverview)
	analytics.Get("/sales", analyticsHandler.GetSalesStats)
	analytics.Get("/users", analyticsHandler.GetUserStats)
	analytics.Get("/reading", analyticsHandler.GetReadingStats)
	analytics.Get("/revenue", analyticsHandler.GetRevenueReport)
	analytics.Get("/growth", analyticsHandler.GetGrowthMetrics)

	api.Get("/analytics/reading", middleware.AdminRequired(), analyticsHandler.GetReadingAnalyticsByPeriod)

	reports := api.Group("/admin/reports", middleware.AdminRequired())
	reports.Get("/data", analyticsHandler.GetReportsData)
	reports.Post("/generate", analyticsHandler.GenerateReport)
	reports.Get("/download/:type", analyticsHandler.DownloadReport)

	notifications := api.Group("/notifications", middleware.AuthRequired())
	notifications.Get("/", notificationHandler.GetNotifications)
	notifications.Get("/unread-count", notificationHandler.GetUnreadCount)
	notifications.Get("/:id", notificationHandler.GetNotification)
	notifications.Patch("/:id/read", notificationHandler.MarkAsRead)
	notifications.Post("/mark-all-read", notificationHandler.MarkAllAsRead)
	notifications.Delete("/:id", notificationHandler.Delete)

	audit := api.Group("/admin/audit", middleware.AdminRequired())
	audit.Get("/", auditHandler.GetLogs)
	audit.Get("/:id", auditHandler.GetLog)

	api.Post("/books/:id/reviews", middleware.AuthRequired(), reviewHandler.CreateReview)
	api.Get("/books/:id/reviews", reviewHandler.GetBookReviews)
	api.Get("/reviews/featured", reviewHandler.GetFeatured)

	adminReviews := api.Group("/admin/reviews", middleware.AdminRequired())
	adminReviews.Get("/", reviewHandler.ListAllReviews)
	adminReviews.Get("/stats", reviewHandler.GetStats)
	adminReviews.Patch("/", reviewHandler.UpdateStatus)
	adminReviews.Patch("/feature", reviewHandler.ToggleFeatured)
	adminReviews.Delete("/:id", reviewHandler.Delete)

	api.Get("/about", aboutHandler.Get)
	api.Put("/admin/about", middleware.AdminRequired(), aboutHandler.Update)

	api.Get("/dashboard/activity", middleware.AuthRequired(), activityHandler.GetActivities)
	api.Get("/dashboard/stats", middleware.AuthRequired(), libraryHandler.GetDashboardStats)
	api.Get("/dashboard/reading-progress", middleware.AuthRequired(), readingHandler.GetReadingProgress)
	api.Get("/dashboard/analytics", middleware.AuthRequired(), libraryHandler.GetUserAnalytics)

	wishlist := api.Group("/wishlist", middleware.AuthRequired())
	wishlist.Get("/", wishlistHandler.GetWishlist)
	wishlist.Post("/", wishlistHandler.AddToWishlist)
	wishlist.Delete("/:id", wishlistHandler.RemoveFromWishlist)

	api.Get("/api/works", workHandler.GetAll)
	adminWorks := api.Group("/admin/works", middleware.AdminRequired())
	adminWorks.Get("/", workHandler.GetAllAdmin)
	adminWorks.Post("/", workHandler.Create)
	adminWorks.Put("/:id", workHandler.Update)
	adminWorks.Patch("/:id/toggle", workHandler.Toggle)
	adminWorks.Delete("/:id", workHandler.Delete)

	api.Get("/admin/system-settings/public", settingsHandler.GetPublic)
}
