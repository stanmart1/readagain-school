package services

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

type ReportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

func (s *ReportService) GeneratePDF(reportType string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	switch reportType {
	case "student_reading":
		return s.generateStudentReadingReport(pdf)
	case "book_popularity":
		return s.generateBookPopularityReport(pdf)
	case "reading_completion":
		return s.generateReadingCompletionReport(pdf)
	case "library_usage":
		return s.generateLibraryUsageReport(pdf)
	default:
		return nil, fmt.Errorf("unknown report type: %s", reportType)
	}
}

func (s *ReportService) generateStudentReadingReport(pdf *gofpdf.Fpdf) ([]byte, error) {
	pdf.Cell(0, 10, "Student Reading Report")
	pdf.Ln(15)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04")))
	pdf.Ln(10)

	type StudentReading struct {
		Name        string
		Email       string
		ClassLevel  string
		SchoolName  string
		BooksRead   int64
		ReadingTime int64
	}

	var students []StudentReading
	s.db.Raw(`
		SELECT u.name, u.email, u.class_level, u.school_name,
		       COUNT(DISTINCT ul.book_id) as books_read,
		       COALESCE(SUM(rs.duration), 0) as reading_time
		FROM users u
		LEFT JOIN user_libraries ul ON u.id = ul.user_id
		LEFT JOIN reading_sessions rs ON u.id = rs.user_id
		WHERE u.role_id = (SELECT id FROM roles WHERE name = 'student')
		GROUP BY u.id, u.name, u.email, u.class_level, u.school_name
		ORDER BY books_read DESC
		LIMIT 50
	`).Scan(&students)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(50, 8, "Student")
	pdf.Cell(40, 8, "Class")
	pdf.Cell(30, 8, "Books")
	pdf.Cell(40, 8, "Reading Time (min)")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 9)
	for _, s := range students {
		pdf.Cell(50, 7, s.Name)
		pdf.Cell(40, 7, s.ClassLevel)
		pdf.Cell(30, 7, fmt.Sprintf("%d", s.BooksRead))
		pdf.Cell(40, 7, fmt.Sprintf("%d", s.ReadingTime/60))
		pdf.Ln(7)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

func (s *ReportService) generateBookPopularityReport(pdf *gofpdf.Fpdf) ([]byte, error) {
	pdf.Cell(0, 10, "Book Popularity Report")
	pdf.Ln(15)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04")))
	pdf.Ln(10)

	type PopularBook struct {
		Title        string
		Author       string
		LibraryCount int
		Views        int
		Rating       float64
	}

	var books []PopularBook
	s.db.Raw(`
		SELECT b.title, b.author, b.library_count, b.views, b.rating
		FROM books b
		ORDER BY b.library_count DESC
		LIMIT 50
	`).Scan(&books)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(70, 8, "Book Title")
	pdf.Cell(40, 8, "Author")
	pdf.Cell(30, 8, "In Libraries")
	pdf.Cell(25, 8, "Views")
	pdf.Cell(20, 8, "Rating")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 9)
	for _, b := range books {
		pdf.Cell(70, 7, b.Title)
		pdf.Cell(40, 7, b.Author)
		pdf.Cell(30, 7, fmt.Sprintf("%d", b.LibraryCount))
		pdf.Cell(25, 7, fmt.Sprintf("%d", b.Views))
		pdf.Cell(20, 7, fmt.Sprintf("%.1f", b.Rating))
		pdf.Ln(7)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

func (s *ReportService) generateReadingCompletionReport(pdf *gofpdf.Fpdf) ([]byte, error) {
	pdf.Cell(0, 10, "Reading Completion Report")
	pdf.Ln(15)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04")))
	pdf.Ln(10)

	type CompletionStats struct {
		BookTitle      string
		TotalReaders   int64
		CompletedCount int64
		AvgProgress    float64
	}

	var stats []CompletionStats
	s.db.Raw(`
		SELECT b.title as book_title,
		       COUNT(ul.id) as total_readers,
		       SUM(CASE WHEN ul.progress = 100 THEN 1 ELSE 0 END) as completed_count,
		       AVG(ul.progress) as avg_progress
		FROM books b
		LEFT JOIN user_libraries ul ON b.id = ul.book_id
		GROUP BY b.id, b.title
		HAVING COUNT(ul.id) > 0
		ORDER BY completed_count DESC
		LIMIT 50
	`).Scan(&stats)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(80, 8, "Book Title")
	pdf.Cell(30, 8, "Readers")
	pdf.Cell(30, 8, "Completed")
	pdf.Cell(35, 8, "Avg Progress")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 9)
	for _, s := range stats {
		pdf.Cell(80, 7, s.BookTitle)
		pdf.Cell(30, 7, fmt.Sprintf("%d", s.TotalReaders))
		pdf.Cell(30, 7, fmt.Sprintf("%d", s.CompletedCount))
		pdf.Cell(35, 7, fmt.Sprintf("%.1f%%", s.AvgProgress))
		pdf.Ln(7)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

func (s *ReportService) generateLibraryUsageReport(pdf *gofpdf.Fpdf) ([]byte, error) {
	pdf.Cell(0, 10, "Library Usage Report")
	pdf.Ln(15)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04")))
	pdf.Ln(10)

	type UsageStats struct {
		Date           string
		BooksAdded     int64
		ActiveReaders  int64
		TotalSessions  int64
		TotalDuration  int64
	}

	var stats []UsageStats
	s.db.Raw(`
		SELECT DATE(ul.created_at) as date,
		       COUNT(ul.id) as books_added,
		       COUNT(DISTINCT ul.user_id) as active_readers,
		       0 as total_sessions,
		       0 as total_duration
		FROM user_libraries ul
		WHERE ul.created_at >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(ul.created_at)
		ORDER BY date DESC
	`).Scan(&stats)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 8, "Date")
	pdf.Cell(35, 8, "Books Added")
	pdf.Cell(40, 8, "Active Readers")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 9)
	for _, s := range stats {
		pdf.Cell(40, 7, s.Date)
		pdf.Cell(35, 7, fmt.Sprintf("%d", s.BooksAdded))
		pdf.Cell(40, 7, fmt.Sprintf("%d", s.ActiveReaders))
		pdf.Ln(7)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
