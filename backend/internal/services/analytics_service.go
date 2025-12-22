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
	s.db.Model(&struct{ ID uint }{}).Table("orders").Count(&overview.TotalOrders)
	s.db.Model(&struct{ ID uint }{}).Table("books").Count(&overview.TotalBooks)
	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Where("progress = ?", 100).Count(&overview.TotalBooksRead)
	s.db.Model(&struct{ ID uint }{}).Table("reading_sessions").Select("COALESCE(SUM(duration), 0)").Scan(&overview.TotalReadingTime)
	s.db.Model(&struct{ ID uint }{}).Table("users").Where("created_at >= ?", today).Count(&overview.NewUsersToday)
	s.db.Model(&struct{ ID uint }{}).Table("orders").Where("created_at >= ?", today).Count(&overview.OrdersToday)

	s.db.Model(&struct{ ID uint }{}).Table("orders").Select("COALESCE(SUM(total_amount), 0)").Scan(&overview.TotalRevenue)
	s.db.Model(&struct{ ID uint }{}).Table("orders").Where("created_at >= ?", today).Select("COALESCE(SUM(total_amount), 0)").Scan(&overview.RevenueToday)

	return &overview, nil
}

func (s *AnalyticsService) GetSalesStats(startDate, endDate *time.Time) (*SalesStats, error) {
	var stats SalesStats
	query := s.db.Model(&struct{ ID uint }{}).Table("orders")

	if startDate != nil {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", endDate)
	}

	query.Count(&stats.TotalOrders)
	s.db.Model(&struct{ ID uint }{}).Table("orders").Where("status = ?", "completed").Count(&stats.CompletedOrders)
	s.db.Model(&struct{ ID uint }{}).Table("orders").Where("status = ?", "pending").Count(&stats.PendingOrders)
	s.db.Model(&struct{ ID uint }{}).Table("orders").Where("status = ?", "cancelled").Count(&stats.CancelledOrders)

	query.Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TotalRevenue)

	if stats.TotalOrders > 0 {
		stats.AverageOrderValue = stats.TotalRevenue / float64(stats.TotalOrders)
	}

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
	var currentRevenue, previousRevenue float64
	var currentOrders, previousOrders int64

	s.db.Model(&struct{ ID uint }{}).Table("users").Where("created_at >= ?", lastMonth).Count(&currentUsers)
	s.db.Model(&struct{ ID uint }{}).Table("users").Where("created_at >= ? AND created_at < ?", twoMonthsAgo, lastMonth).Count(&previousUsers)

	s.db.Model(&struct{ ID uint }{}).Table("orders").Where("created_at >= ?", lastMonth).Select("COALESCE(SUM(total_amount), 0)").Scan(&currentRevenue)
	s.db.Model(&struct{ ID uint }{}).Table("orders").Where("created_at >= ? AND created_at < ?", twoMonthsAgo, lastMonth).Select("COALESCE(SUM(total_amount), 0)").Scan(&previousRevenue)

	s.db.Model(&struct{ ID uint }{}).Table("orders").Where("created_at >= ?", lastMonth).Count(&currentOrders)
	s.db.Model(&struct{ ID uint }{}).Table("orders").Where("created_at >= ? AND created_at < ?", twoMonthsAgo, lastMonth).Count(&previousOrders)

	if previousUsers > 0 {
		metrics.UsersGrowth = ((float64(currentUsers) - float64(previousUsers)) / float64(previousUsers)) * 100
	}
	if previousRevenue > 0 {
		metrics.RevenueGrowth = ((currentRevenue - previousRevenue) / previousRevenue) * 100
	}
	if previousOrders > 0 {
		metrics.OrdersGrowth = ((float64(currentOrders) - float64(previousOrders)) / float64(previousOrders)) * 100
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
			"total_users":    overview.TotalUsers,
			"total_books":    overview.TotalBooks,
			"total_orders":   overview.TotalOrders,
			"total_revenue":  overview.TotalRevenue,
			"user_growth":    growth.UsersGrowth,
			"revenue_growth": growth.RevenueGrowth,
			"order_growth":   growth.OrdersGrowth,
			"book_growth":    0.0,
		},
		"user_growth":      userGrowth,
		"daily_activity":   []interface{}{},
		"recent_activities": []interface{}{},
	}, nil
}

func (s *AnalyticsService) GetReadingAnalyticsByPeriod(period string) (map[string]interface{}, error) {
	var totalSessions, booksStarted int64
	var totalPages int
	var avgSessionTime float64

	s.db.Model(&struct{ ID uint }{}).Table("reading_sessions").Count(&totalSessions)
	s.db.Model(&struct{ ID uint }{}).Table("user_libraries").Where("progress > ?", 0).Count(&booksStarted)
	s.db.Model(&struct{ ID uint }{}).Table("reading_sessions").Select("COALESCE(AVG(duration), 0)").Scan(&avgSessionTime)
	s.db.Model(&struct{ ID uint }{}).Table("reading_sessions").Select("COALESCE(SUM(pages_read), 0)").Scan(&totalPages)

	var currentlyReading []map[string]interface{}
	rows, _ := s.db.Raw(`
		SELECT b.title, ul.progress
		FROM user_libraries ul
		JOIN books b ON b.id = ul.book_id
		WHERE ul.progress > 0 AND ul.progress < 100
		LIMIT 10
	`).Rows()
	defer rows.Close()

	for rows.Next() {
		var title string
		var progress float64
		rows.Scan(&title, &progress)
		currentlyReading = append(currentlyReading, map[string]interface{}{
			"title":    title,
			"progress": progress,
		})
	}

	return map[string]interface{}{
		"stats": map[string]interface{}{
			"totalSessions":      totalSessions,
			"booksStarted":       booksStarted,
			"averageSessionTime": avgSessionTime,
			"totalPages":         totalPages,
		},
		"weeklyData":       []interface{}{},
		"currentlyReading": currentlyReading,
	}, nil
}

func (s *AnalyticsService) GetReportsData() (map[string]interface{}, error) {
	sales, _ := s.GetSalesStats(nil, nil)
	users, _ := s.GetUserStats()
	reading, _ := s.GetReadingStats()

	return map[string]interface{}{
		"sales":   sales,
		"users":   users,
		"reading": reading,
	}, nil
}

