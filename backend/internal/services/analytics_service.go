package services

import (
	"time"

	"gorm.io/gorm"
)

type AnalyticsService struct {
	db *gorm.DB
}

func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

type DashboardOverview struct {
	TotalUsers       int64   `json:"total_users"`
	ActiveUsers      int64   `json:"active_users"`
	TotalOrders      int64   `json:"total_orders"`
	TotalRevenue     float64 `json:"total_revenue"`
	TotalBooks       int64   `json:"total_books"`
	TotalBooksRead   int64   `json:"total_books_read"`
	TotalReadingTime int64   `json:"total_reading_time"`
	NewUsersToday    int64   `json:"new_users_today"`
	OrdersToday      int64   `json:"orders_today"`
	RevenueToday     float64 `json:"revenue_today"`
}

type SalesStats struct {
	TotalOrders      int64   `json:"total_orders"`
	CompletedOrders  int64   `json:"completed_orders"`
	PendingOrders    int64   `json:"pending_orders"`
	CancelledOrders  int64   `json:"cancelled_orders"`
	TotalRevenue     float64 `json:"total_revenue"`
	AverageOrderValue float64 `json:"average_order_value"`
}

type UserStats struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	InactiveUsers int64 `json:"inactive_users"`
	NewThisMonth  int64 `json:"new_this_month"`
	NewThisWeek   int64 `json:"new_this_week"`
}

type ReadingStats struct {
	TotalBooksRead     int64 `json:"total_books_read"`
	TotalReadingTime   int64 `json:"total_reading_time"`
	TotalSessions      int64 `json:"total_sessions"`
	AverageSessionTime int64 `json:"average_session_time"`
	ActiveReaders      int64 `json:"active_readers"`
}

type RevenueReport struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

type GrowthMetrics struct {
	UsersGrowth   float64 `json:"users_growth"`
	RevenueGrowth float64 `json:"revenue_growth"`
	OrdersGrowth  float64 `json:"orders_growth"`
}

func (s *AnalyticsService) GetDashboardOverview() (*DashboardOverview, error) {
	var overview DashboardOverview
	today := time.Now().Truncate(24 * time.Hour)

	s.db.Model(&struct{ ID uint }{}).Table("users").Count(&overview.TotalUsers)
	s.db.Model(&struct{ ID uint }{}).Table("users").Where("is_active = ?", true).Count(&overview.ActiveUsers)
	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Count(&overview.TotalOrders) // Reusing field for books in libraries
	s.db.Model(&struct{ ID uint }{}).Table("books").Count(&overview.TotalBooks)
	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Where("progress = ?", 100).Count(&overview.TotalBooksRead)
	s.db.Model(&struct{ ID uint }{}).Table("reading_sessions").Select("COALESCE(SUM(duration), 0)").Scan(&overview.TotalReadingTime)
	s.db.Model(&struct{ ID uint }{}).Table("users").Where("created_at >= ?", today).Count(&overview.NewUsersToday)
	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Where("created_at >= ?", today).Count(&overview.OrdersToday) // Reusing field for books added today

	// Count active readers (users with reading sessions in last 7 days)
	weekAgo := today.AddDate(0, 0, -7)
	s.db.Raw("SELECT COUNT(DISTINCT user_id) FROM reading_sessions WHERE created_at >= ?", weekAgo).Scan(&overview.TotalRevenue) // Reusing field
	s.db.Raw("SELECT COUNT(DISTINCT user_id) FROM reading_sessions WHERE created_at >= ?", today).Scan(&overview.RevenueToday) // Reusing field

	return &overview, nil
}

func (s *AnalyticsService) GetSalesStats(startDate, endDate *time.Time) (*SalesStats, error) {
	var stats SalesStats
	query := s.db.Model(&struct{ ID uint }{}).Table("user_libraries")

	if startDate != nil {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", endDate)
	}

	query.Count(&stats.TotalOrders) // Books in libraries
	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Where("progress = ?", 100).Count(&stats.CompletedOrders) // Completed books
	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Where("progress > ? AND progress < ?", 0, 100).Count(&stats.PendingOrders) // In progress
	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Where("progress = ?", 0).Count(&stats.CancelledOrders) // Not started

	// No revenue for school platform
	stats.TotalRevenue = 0
	stats.AverageOrderValue = 0

	return &stats, nil
}

func (s *AnalyticsService) GetUserStats() (*UserStats, error) {
	var stats UserStats
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)
	monthAgo := now.AddDate(0, -1, 0)

	s.db.Model(&struct{ ID uint }{}).Table("users").Count(&stats.TotalUsers)
	s.db.Model(&struct{ ID uint }{}).Table("users").Where("is_active = ?", true).Count(&stats.ActiveUsers)
	s.db.Model(&struct{ ID uint }{}).Table("users").Where("is_active = ?", false).Count(&stats.InactiveUsers)
	s.db.Model(&struct{ ID uint }{}).Table("users").Where("created_at >= ?", monthAgo).Count(&stats.NewThisMonth)
	s.db.Model(&struct{ ID uint }{}).Table("users").Where("created_at >= ?", weekAgo).Count(&stats.NewThisWeek)

	return &stats, nil
}

func (s *AnalyticsService) GetReadingStats() (*ReadingStats, error) {
	var stats ReadingStats

	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Where("progress = ?", 100).Count(&stats.TotalBooksRead)
	s.db.Model(&struct{ ID uint }{}).Table("reading_sessions").Select("COALESCE(SUM(duration), 0)").Scan(&stats.TotalReadingTime)
	s.db.Model(&struct{ ID uint }{}).Table("reading_sessions").Count(&stats.TotalSessions)
	s.db.Model(&struct{ ID uint }{}).Table("reading_sessions").Distinct("user_id").Count(&stats.ActiveReaders)

	if stats.TotalSessions > 0 {
		stats.AverageSessionTime = stats.TotalReadingTime / stats.TotalSessions
	}

	return &stats, nil
}

func (s *AnalyticsService) GetRevenueReport(days int) ([]RevenueReport, error) {
	var reports []RevenueReport
	startDate := time.Now().AddDate(0, 0, -days)

	rows, err := s.db.Raw(`
		SELECT 
			DATE(created_at) as date,
			COALESCE(SUM(total_amount), 0) as revenue,
			COUNT(*) as orders
		FROM orders
		WHERE created_at >= ? AND status = 'completed'
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`, startDate).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var report RevenueReport
		if err := rows.Scan(&report.Date, &report.Revenue, &report.Orders); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func (s *AnalyticsService) GetGrowthMetrics() (*GrowthMetrics, error) {
	var metrics GrowthMetrics
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	twoMonthsAgo := now.AddDate(0, -2, 0)

	var currentUsers, previousUsers int64
	var currentLibraries, previousLibraries int64
	var currentReaders, previousReaders int64

	s.db.Model(&struct{ ID uint }{}).Table("users").Where("created_at >= ?", lastMonth).Count(&currentUsers)
	s.db.Model(&struct{ ID uint }{}).Table("users").Where("created_at >= ? AND created_at < ?", twoMonthsAgo, lastMonth).Count(&previousUsers)

	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Where("created_at >= ?", lastMonth).Count(&currentLibraries)
	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Where("created_at >= ? AND created_at < ?", twoMonthsAgo, lastMonth).Count(&previousLibraries)

	s.db.Raw("SELECT COUNT(DISTINCT user_id) FROM reading_sessions WHERE created_at >= ?", lastMonth).Scan(&currentReaders)
	s.db.Raw("SELECT COUNT(DISTINCT user_id) FROM reading_sessions WHERE created_at >= ? AND created_at < ?", twoMonthsAgo, lastMonth).Scan(&previousReaders)

	if previousUsers > 0 {
		metrics.UsersGrowth = ((float64(currentUsers) - float64(previousUsers)) / float64(previousUsers)) * 100
	}
	if previousLibraries > 0 {
		metrics.RevenueGrowth = ((float64(currentLibraries) - float64(previousLibraries)) / float64(previousLibraries)) * 100 // Reusing field for library growth
	}
	if previousReaders > 0 {
		metrics.OrdersGrowth = ((float64(currentReaders) - float64(previousReaders)) / float64(previousReaders)) * 100 // Reusing field for reader growth
	}

	return &metrics, nil
}

func (s *AnalyticsService) GetEnhancedOverview() (map[string]interface{}, error) {
	overview, err := s.GetDashboardOverview()
	if err != nil {
		return nil, err
	}

	growth, _ := s.GetGrowthMetrics()

	var userGrowth []map[string]interface{}
	rows, _ := s.db.Raw(`
		SELECT 
			TO_CHAR(created_at, 'Mon YYYY') as month,
			COUNT(*) as users
		FROM users
		WHERE created_at >= NOW() - INTERVAL '6 months'
		GROUP BY TO_CHAR(created_at, 'Mon YYYY'), DATE_TRUNC('month', created_at)
		ORDER BY DATE_TRUNC('month', created_at)
	`).Rows()
	defer rows.Close()

	for rows.Next() {
		var month string
		var users int64
		rows.Scan(&month, &users)
		userGrowth = append(userGrowth, map[string]interface{}{
			"month": month,
			"users": users,
		})
	}

	return map[string]interface{}{
		"overview": map[string]interface{}{
			"total_users":        overview.TotalUsers,
			"total_books":        overview.TotalBooks,
			"books_in_libraries": overview.TotalOrders,    // Reused field
			"active_readers":     overview.TotalRevenue,   // Reused field
			"user_growth":        growth.UsersGrowth,
			"library_growth":     growth.RevenueGrowth,    // Reused field
			"reader_growth":      growth.OrdersGrowth,     // Reused field
			"book_growth":        0.0,
		},
		"user_growth":       userGrowth,
		"daily_activity":    []interface{}{},
		"recent_activities": []interface{}{},
	}, nil
}

func (s *AnalyticsService) GetReadingAnalyticsByPeriod(period string) (map[string]interface{}, error) {
	// Calculate date range based on period
	now := time.Now()
	var startDate time.Time
	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
	case "year":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, -1, 0)
	}

	// Class/Grade stats
	type ClassStats struct {
		ClassLevel      string  `json:"class_level"`
		StudentCount    int64   `json:"student_count"`
		AvgReadingTime  float64 `json:"avg_reading_time"`
		AvgCompletion   float64 `json:"avg_completion"`
		BooksCompleted  int64   `json:"books_completed"`
	}
	var classStats []ClassStats
	s.db.Raw(`
		SELECT u.class_level,
		       COUNT(DISTINCT u.id) as student_count,
		       COALESCE(AVG(rs.duration), 0) as avg_reading_time,
		       COALESCE(AVG(ul.progress), 0) as avg_completion,
		       COUNT(CASE WHEN ul.progress = 100 THEN 1 END) as books_completed
		FROM users u
		LEFT JOIN reading_sessions rs ON u.id = rs.user_id AND rs.created_at >= ?
		LEFT JOIN user_libraries ul ON u.id = ul.user_id AND ul.created_at >= ?
		WHERE u.role_id = (SELECT id FROM roles WHERE name = 'student')
		GROUP BY u.class_level
		ORDER BY u.class_level
	`, startDate, startDate).Scan(&classStats)

	// Struggling readers (< 30% avg completion)
	type StrugglingReader struct {
		Name           string  `json:"name"`
		Email          string  `json:"email"`
		ClassLevel     string  `json:"class_level"`
		AvgCompletion  float64 `json:"avg_completion"`
		BooksStarted   int64   `json:"books_started"`
	}
	var strugglingReaders []StrugglingReader
	s.db.Raw(`
		SELECT CONCAT(u.first_name, ' ', u.last_name) as name, u.email, u.class_level,
		       COALESCE(AVG(ul.progress), 0) as avg_completion,
		       COUNT(ul.id) as books_started
		FROM users u
		LEFT JOIN user_libraries ul ON u.id = ul.user_id AND ul.created_at >= ?
		WHERE u.role_id = (SELECT id FROM roles WHERE name = 'student')
		GROUP BY u.id, u.first_name, u.last_name, u.email, u.class_level
		HAVING COALESCE(AVG(ul.progress), 0) < 30 AND COUNT(ul.id) > 0
		ORDER BY avg_completion ASC
		LIMIT 20
	`, startDate).Scan(&strugglingReaders)

	// Most/Least read books by grade
	type BookByGrade struct {
		BookTitle      string `json:"book_title"`
		Author         string `json:"author"`
		ClassLevel     string `json:"class_level"`
		ReaderCount    int64  `json:"reader_count"`
		AvgCompletion  float64 `json:"avg_completion"`
	}
	var mostReadBooks []BookByGrade
	s.db.Raw(`
		SELECT b.title as book_title, a.business_name as author, u.class_level,
		       COUNT(DISTINCT ul.user_id) as reader_count,
		       COALESCE(AVG(ul.progress), 0) as avg_completion
		FROM books b
		LEFT JOIN authors a ON b.author_id = a.id
		JOIN user_libraries ul ON b.id = ul.book_id AND ul.created_at >= ?
		JOIN users u ON ul.user_id = u.id
		WHERE u.role_id = (SELECT id FROM roles WHERE name = 'student')
		GROUP BY b.id, b.title, a.business_name, u.class_level
		ORDER BY reader_count DESC
		LIMIT 50
	`, startDate).Scan(&mostReadBooks)

	// Top readers with streaks
	type TopReader struct {
		Name           string `json:"name"`
		ClassLevel     string `json:"class_level"`
		BooksCompleted int64  `json:"books_completed"`
		ReadingTime    int64  `json:"reading_time"`
		CurrentStreak  int    `json:"current_streak"`
	}
	var topReaders []TopReader
	s.db.Raw(`
		SELECT CONCAT(u.first_name, ' ', u.last_name) as name, u.class_level,
		       COUNT(CASE WHEN ul.progress = 100 THEN 1 END) as books_completed,
		       COALESCE(SUM(rs.duration), 0) as reading_time,
		       0 as current_streak
		FROM users u
		LEFT JOIN user_libraries ul ON u.id = ul.user_id AND ul.created_at >= ?
		LEFT JOIN reading_sessions rs ON u.id = rs.user_id AND rs.created_at >= ?
		WHERE u.role_id = (SELECT id FROM roles WHERE name = 'student')
		GROUP BY u.id, u.first_name, u.last_name, u.class_level
		ORDER BY books_completed DESC, reading_time DESC
		LIMIT 20
	`, startDate, startDate).Scan(&topReaders)

	// Active readers with recent library activity
	type ActiveReader struct {
		UserID         uint      `json:"user_id"`
		LibraryID      uint      `json:"library_id"`
		Name           string    `json:"name"`
		Email          string    `json:"email"`
		ClassLevel     string    `json:"class_level"`
		SessionCount   int64     `json:"session_count"`
		TotalTime      int64     `json:"total_time"`
		LastSession    time.Time `json:"last_session"`
	}
	var activeReaders []ActiveReader
	s.db.Raw(`
		SELECT u.id as user_id, MAX(ul.id) as library_id,
		       CONCAT(u.first_name, ' ', u.last_name) as name, u.email, u.class_level,
		       COUNT(DISTINCT ul.id) as session_count,
		       COALESCE(SUM(rs.duration), 0) as total_time,
		       COALESCE(MAX(ul.last_read_at), MAX(ul.updated_at)) as last_session
		FROM users u
		JOIN user_libraries ul ON u.id = ul.user_id AND ul.created_at >= ?
		LEFT JOIN reading_sessions rs ON u.id = rs.user_id AND ul.book_id = rs.book_id
		GROUP BY u.id, u.first_name, u.last_name, u.email, u.class_level
		ORDER BY last_session DESC
		LIMIT 50
	`, startDate).Scan(&activeReaders)

	// Overall stats
	var totalActiveReaders, totalBooksCompleted int64
	var avgReadingTime, avgCompletionRate float64
	s.db.Raw(`
		SELECT COUNT(DISTINCT u.id) as total_readers
		FROM users u
		JOIN user_libraries ul ON u.id = ul.user_id AND ul.created_at >= ?
	`, startDate).Scan(&totalActiveReaders)
	
	s.db.Raw(`SELECT COUNT(*) FROM user_libraries WHERE progress = 100 AND created_at >= ?`, startDate).Scan(&totalBooksCompleted)
	s.db.Raw(`SELECT COALESCE(AVG(duration), 0) FROM reading_sessions WHERE created_at >= ?`, startDate).Scan(&avgReadingTime)
	s.db.Raw(`SELECT COALESCE(AVG(progress), 0) FROM user_libraries WHERE created_at >= ?`, startDate).Scan(&avgCompletionRate)

	return map[string]interface{}{
		"overview": map[string]interface{}{
			"total_active_readers": totalActiveReaders,
			"total_books_completed": totalBooksCompleted,
			"avg_reading_time": avgReadingTime / 60, // Convert to hours
			"avg_completion_rate": avgCompletionRate,
		},
		"class_stats": classStats,
		"struggling_readers": strugglingReaders,
		"most_read_books": mostReadBooks,
		"top_readers": topReaders,
		"active_readers": activeReaders,
	}, nil
}

func (s *AnalyticsService) GetReportsData() (map[string]interface{}, error) {
	// Get reading statistics
	var totalBooksInLibraries int64
	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Count(&totalBooksInLibraries)

	var completedBooks int64
	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Where("progress = ?", 100).Count(&completedBooks)

	var activeReaders int64
	weekAgo := time.Now().AddDate(0, 0, -7)
	s.db.Raw("SELECT COUNT(DISTINCT user_id) FROM reading_sessions WHERE created_at >= ?", weekAgo).Scan(&activeReaders)

	// Get popular books
	type PopularBook struct {
		BookID       uint    `json:"book_id"`
		Title        string  `json:"title"`
		LibraryCount int     `json:"library_count"`
		Views        int     `json:"views"`
		Rating       float64 `json:"rating"`
	}
	var popularBooks []PopularBook
	s.db.Raw(`
		SELECT 
			b.id as book_id, 
			b.title, 
			b.library_count,
			b.view_count as views,
			COALESCE(ROUND(AVG(r.rating)::numeric, 1), 0) as rating
		FROM books b
		LEFT JOIN reviews r ON b.id = r.book_id
		GROUP BY b.id, b.title, b.library_count, b.view_count
		ORDER BY b.library_count DESC, b.view_count DESC
		LIMIT 10
	`).Scan(&popularBooks)

	// Available report types
	reports := []map[string]interface{}{
		{
			"type":          "student_reading",
			"title":         "Student Reading Report",
			"description":   "Detailed reading activity by student, class, and grade",
			"status":        "ready",
			"lastGenerated": time.Now().Format("2006-01-02"),
		},
		{
			"type":          "book_popularity",
			"title":         "Book Popularity Report",
			"description":   "Most popular books by library additions and completion rates",
			"status":        "ready",
			"lastGenerated": time.Now().Format("2006-01-02"),
		},
		{
			"type":          "reading_completion",
			"title":         "Reading Completion Report",
			"description":   "Book completion rates and reading progress statistics",
			"status":        "ready",
			"lastGenerated": time.Now().Format("2006-01-02"),
		},
		{
			"type":          "library_usage",
			"title":         "Library Usage Report",
			"description":   "Library access patterns and usage statistics",
			"status":        "ready",
			"lastGenerated": time.Now().Format("2006-01-02"),
		},
	}

	completionRate := 0.0
	if totalBooksInLibraries > 0 {
		completionRate = float64(completedBooks) / float64(totalBooksInLibraries) * 100
	}

	return map[string]interface{}{
		"summary": map[string]interface{}{
			"total_books_in_libraries": totalBooksInLibraries,
			"completed_books":          completedBooks,
			"active_readers":           activeReaders,
			"completion_rate":          completionRate,
		},
		"popularBooks": popularBooks,
		"reports":      reports,
	}, nil
}

