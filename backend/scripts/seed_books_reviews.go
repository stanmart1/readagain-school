package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"readagain/internal/config"
	"readagain/internal/database"
	"readagain/internal/models"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Load config
	cfg := config.Load()
	
	// Connect to database
	if err := database.Connect(cfg.Database.URL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := database.DB

	// Get existing authors and categories
	var authors []models.Author
	db.Find(&authors)
	if len(authors) == 0 {
		log.Fatal("No authors found. Please create authors first.")
	}

	var categories []models.Category
	db.Find(&categories)
	if len(categories) == 0 {
		log.Fatal("No categories found. Please create categories first.")
	}

	// Get student users for reviews
	var students []models.User
	db.Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", "student").
		Find(&students)
	if len(students) == 0 {
		log.Fatal("No student users found. Please create student users first.")
	}

	// Sample book data
	bookData := []struct {
		Title       string
		Subtitle    string
		Description string
		ISBN        string
		Pages       int
		Language    string
		AgeRange    string
	}{
		{
			Title:       "The Adventures of Tom Sawyer",
			Subtitle:    "A Classic American Novel",
			Description: "Follow Tom Sawyer's mischievous adventures along the Mississippi River in this timeless tale of childhood, friendship, and growing up in 19th century America.",
			ISBN:        "978-0-14-303957-0",
			Pages:       274,
			Language:    "English",
			AgeRange:    "8-12",
		},
		{
			Title:       "Charlotte's Web",
			Subtitle:    "A Story of Friendship",
			Description: "The beloved story of a pig named Wilbur and his friendship with a clever spider named Charlotte who saves his life through her web messages.",
			ISBN:        "978-0-06-440055-8",
			Pages:       192,
			Language:    "English",
			AgeRange:    "6-10",
		},
		{
			Title:       "Harry Potter and the Sorcerer's Stone",
			Subtitle:    "Book 1",
			Description: "Harry Potter discovers he's a wizard on his 11th birthday and begins his magical education at Hogwarts School of Witchcraft and Wizardry.",
			ISBN:        "978-0-439-70818-8",
			Pages:       309,
			Language:    "English",
			AgeRange:    "8-12",
		},
		{
			Title:       "The Lion, the Witch and the Wardrobe",
			Subtitle:    "The Chronicles of Narnia",
			Description: "Four children discover a magical wardrobe that leads to the enchanted land of Narnia, where they must help defeat the White Witch.",
			ISBN:        "978-0-06-447104-7",
			Pages:       206,
			Language:    "English",
			AgeRange:    "8-12",
		},
		{
			Title:       "Wonder",
			Subtitle:    "A Story of Kindness",
			Description: "Auggie Pullman was born with facial differences that prevented him from going to a mainstream school. Now he's starting fifth grade at a regular school.",
			ISBN:        "978-0-375-86902-0",
			Pages:       310,
			Language:    "English",
			AgeRange:    "8-12",
		},
		{
			Title:       "The Giver",
			Subtitle:    "A Dystopian Classic",
			Description: "In a seemingly perfect community without war, pain, or suffering, a young boy discovers the dark secrets behind his world's facade.",
			ISBN:        "978-0-544-33626-0",
			Pages:       240,
			Language:    "English",
			AgeRange:    "10-14",
		},
		{
			Title:       "Matilda",
			Subtitle:    "A Roald Dahl Classic",
			Description: "Matilda is a brilliant child with a magical mind, but her parents and headmistress think she's nothing but trouble. She'll show them!",
			ISBN:        "978-0-14-241037-7",
			Pages:       240,
			Language:    "English",
			AgeRange:    "7-11",
		},
		{
			Title:       "Percy Jackson and the Lightning Thief",
			Subtitle:    "Book 1",
			Description: "Percy Jackson discovers he's a demigod and must prevent a war between the Greek gods by finding Zeus's stolen lightning bolt.",
			ISBN:        "978-0-7868-5629-9",
			Pages:       377,
			Language:    "English",
			AgeRange:    "9-12",
		},
		{
			Title:       "The Secret Garden",
			Subtitle:    "A Timeless Classic",
			Description: "A young orphan discovers a hidden, neglected garden and brings it back to life, transforming herself and those around her in the process.",
			ISBN:        "978-0-14-036671-0",
			Pages:       331,
			Language:    "English",
			AgeRange:    "8-12",
		},
		{
			Title:       "Diary of a Wimpy Kid",
			Subtitle:    "Book 1",
			Description: "Greg Heffley chronicles his middle school experiences in his diary, dealing with bullies, friends, and family in this hilarious illustrated novel.",
			ISBN:        "978-0-8109-9313-6",
			Pages:       217,
			Language:    "English",
			AgeRange:    "8-12",
		},
	}

	rand.Seed(time.Now().UnixNano())

	// Create books
	var createdBooks []models.Book
	for _, data := range bookData {
		categoryID := categories[rand.Intn(len(categories))].ID
		pubDate := time.Now().AddDate(-rand.Intn(10), 0, 0)
		
		book := models.Book{
			AuthorID:         authors[rand.Intn(len(authors))].ID,
			CategoryID:       &categoryID,
			Title:            data.Title,
			Subtitle:         data.Subtitle,
			Description:      data.Description,
			ShortDescription: data.Description[:min(len(data.Description), 150)],
			ISBN:             data.ISBN,
			Pages:            data.Pages,
			Language:         data.Language,
			Publisher:        "Educational Press",
			PublicationDate:  &pubDate,
			IsFeatured:       rand.Float32() < 0.3,
			IsBestseller:     rand.Float32() < 0.2,
			IsNewRelease:     rand.Float32() < 0.15,
			Status:           "published",
			IsActive:         true,
			LibraryCount:     rand.Intn(50),
			ViewCount:        rand.Intn(1000),
		}

		if err := db.Create(&book).Error; err != nil {
			log.Printf("Failed to create book %s: %v", data.Title, err)
			continue
		}

		createdBooks = append(createdBooks, book)
		fmt.Printf("Created book: %s\n", book.Title)
	}

	// Create reviews for books
	reviewComments := []string{
		"Amazing book! My child loved every page.",
		"Great story with valuable life lessons.",
		"Couldn't put it down! Highly recommend.",
		"Perfect for the age group. Very engaging.",
		"Beautiful illustrations and wonderful story.",
		"My student enjoyed reading this book.",
		"Excellent choice for young readers.",
		"Captivating from start to finish.",
		"A must-read for all children.",
		"Inspiring and educational.",
	}

	for _, book := range createdBooks {
		// Create 3-7 reviews per book
		numReviews := 3 + rand.Intn(5)
		for i := 0; i < numReviews; i++ {
			student := students[rand.Intn(len(students))]
			
			review := models.Review{
				UserID:  student.ID,
				BookID:  book.ID,
				Rating:  3 + rand.Intn(3), // Rating between 3-5
				Comment: reviewComments[rand.Intn(len(reviewComments))],
			}

			if err := db.Create(&review).Error; err != nil {
				log.Printf("Failed to create review for book %s: %v", book.Title, err)
				continue
			}
		}
		fmt.Printf("Created reviews for: %s\n", book.Title)
	}

	// Update book ratings
	for _, book := range createdBooks {
		var avgRating float64
		db.Model(&models.Review{}).Where("book_id = ?", book.ID).Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)
		db.Model(&book).Update("rating", avgRating)
	}

	fmt.Printf("\n✅ Successfully seeded %d books with reviews!\n", len(createdBooks))
}
