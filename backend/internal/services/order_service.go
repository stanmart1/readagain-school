package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type OrderService struct {
	db                 *gorm.DB
	emailService       *EmailService
	achievementService *AchievementService
}

func NewOrderService(db *gorm.DB, emailService *EmailService, achievementService *AchievementService) *OrderService {
	return &OrderService{
		db:                 db,
		emailService:       emailService,
		achievementService: achievementService,
	}
}

func (s *OrderService) CreateOrderFromCart(userID uint, paymentMethod string) (*models.Order, error) {
	var cartItems []models.Cart
	if err := s.db.Where("user_id = ?", userID).Preload("Book").Find(&cartItems).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch cart", err)
	}

	if len(cartItems) == 0 {
		return nil, utils.NewBadRequestError("Cart is empty")
	}

	var subtotal float64
	for _, item := range cartItems {
		if item.Book != nil {
			subtotal += item.Book.Price
		}
	}

	orderNumber := s.generateOrderNumber()

	order := models.Order{
		UserID:        userID,
		AuthorID:      1,
		OrderNumber:   orderNumber,
		Subtotal:      subtotal,
		TotalAmount:   subtotal,
		Status:        "pending",
		PaymentMethod: paymentMethod,
		Notes:         orderNumber,
	}

	if err := s.db.Create(&order).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create order", err)
	}

	for _, cartItem := range cartItems {
		orderItem := models.OrderItem{
			OrderID: order.ID,
			BookID:  cartItem.BookID,
			Price:   cartItem.Book.Price,
		}
		if err := s.db.Create(&orderItem).Error; err != nil {
			s.db.Delete(&order)
			return nil, utils.NewInternalServerError("Failed to create order items", err)
		}
	}

	if err := s.db.Preload("Items.Book").Preload("User").First(&order, order.ID).Error; err != nil {
		return nil, utils.NewNotFoundError("Order not found")
	}

	return &order, nil
}

func (s *OrderService) GetUserOrders(userID uint, page, limit int, status string) ([]models.Order, *utils.PaginationMeta, error) {
	params := utils.GetPaginationParams(page, limit)

	query := s.db.Model(&models.Order{}).Where("user_id = ?", userID).Preload("Items.Book")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to count orders", err)
	}

	var orders []models.Order
	if err := query.Scopes(utils.Paginate(params)).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to fetch orders", err)
	}

	meta := utils.GetPaginationMeta(params.Page, params.Limit, total)
	return orders, &meta, nil
}

func (s *OrderService) GetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order
	if err := s.db.Preload("Items.Book.Author").Preload("Items.Book.Category").Preload("User").First(&order, orderID).Error; err != nil {
		return nil, utils.NewNotFoundError("Order not found")
	}
	return &order, nil
}

func (s *OrderService) GetOrderByReference(reference string) (*models.Order, error) {
	var order models.Order
	if err := s.db.Where("notes LIKE ?", reference+"%").Preload("Items").First(&order).Error; err != nil {
		return nil, utils.NewNotFoundError("Order not found")
	}
	return &order, nil
}

func (s *OrderService) UpdateOrderStatus(orderID uint, status string) (*models.Order, error) {
	var order models.Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return nil, utils.NewNotFoundError("Order not found")
	}

	order.Status = status
	if err := s.db.Save(&order).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update order", err)
	}

	return &order, nil
}

func (s *OrderService) CompleteOrder(orderID uint, transactionID string) error {
	var order models.Order
	if err := s.db.Preload("Items").First(&order, orderID).Error; err != nil {
		return utils.NewNotFoundError("Order not found")
	}

	order.Status = "completed"
	order.Notes = fmt.Sprintf("%s | Transaction: %s", order.Notes, transactionID)

	if err := s.db.Save(&order).Error; err != nil {
		return utils.NewInternalServerError("Failed to complete order", err)
	}

	for _, item := range order.Items {
		var existingLibrary models.UserLibrary
		if err := s.db.Where("user_id = ? AND book_id = ?", order.UserID, item.BookID).First(&existingLibrary).Error; err != nil {
			library := models.UserLibrary{
				UserID: order.UserID,
				BookID: item.BookID,
			}
			s.db.Create(&library)

			s.db.Model(&models.Book{}).Where("id = ?", item.BookID).UpdateColumn("download_count", gorm.Expr("download_count + ?", 1))
		}
	}

	s.db.Where("user_id = ?", order.UserID).Delete(&models.Cart{})

	var user models.User
	if err := s.db.First(&user, order.UserID).Error; err == nil {
		var books []OrderBook
		for _, item := range order.Items {
			if item.Book != nil {
				authorName := "Unknown Author"
				if item.Book.Author != nil && item.Book.Author.User != nil {
					authorName = fmt.Sprintf("%s %s", item.Book.Author.User.FirstName, item.Book.Author.User.LastName)
				}
				books = append(books, OrderBook{
					Title:  item.Book.Title,
					Author: authorName,
					Price:  fmt.Sprintf("%.2f", item.Price),
				})
			}
		}

		name := user.FirstName
		if name == "" {
			name = user.Username
		}

		go s.emailService.SendOrderConfirmation(
			user.Email,
			name,
			order.OrderNumber,
			books,
			fmt.Sprintf("%.2f", order.TotalAmount),
		)
	}

	go s.achievementService.CheckAndUnlockAchievements(order.UserID)

	return nil
}

func (s *OrderService) CancelOrder(orderID, userID uint) error {
	var order models.Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return utils.NewNotFoundError("Order not found")
	}

	if order.UserID != userID {
		return utils.NewForbiddenError("You can only cancel your own orders")
	}

	if order.Status == "completed" {
		return utils.NewBadRequestError("Cannot cancel completed order")
	}

	if order.Status == "cancelled" {
		return utils.NewBadRequestError("Order already cancelled")
	}

	order.Status = "cancelled"

	if err := s.db.Save(&order).Error; err != nil {
		return utils.NewInternalServerError("Failed to cancel order", err)
	}

	return nil
}

func (s *OrderService) GetAllOrders(page, limit int, status, paymentMethod string) ([]models.Order, *utils.PaginationMeta, error) {
	params := utils.GetPaginationParams(page, limit)

	query := s.db.Model(&models.Order{}).Preload("Items.Book").Preload("User")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if paymentMethod != "" {
		query = query.Where("payment_method = ?", paymentMethod)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to count orders", err)
	}

	var orders []models.Order
	if err := query.Scopes(utils.Paginate(params)).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to fetch orders", err)
	}

	meta := utils.GetPaginationMeta(params.Page, params.Limit, total)
	return orders, &meta, nil
}

func (s *OrderService) GetOrderStatistics() (map[string]interface{}, error) {
	var totalOrders int64
	s.db.Model(&models.Order{}).Count(&totalOrders)

	var completedOrders int64
	s.db.Model(&models.Order{}).Where("status = ?", "completed").Count(&completedOrders)

	var pendingOrders int64
	s.db.Model(&models.Order{}).Where("status = ?", "pending").Count(&pendingOrders)

	var totalRevenue float64
	s.db.Model(&models.Order{}).Where("status = ?", "completed").Select("COALESCE(SUM(total), 0)").Scan(&totalRevenue)

	return map[string]interface{}{
		"total_orders":     totalOrders,
		"completed_orders": completedOrders,
		"pending_orders":   pendingOrders,
		"total_revenue":    totalRevenue,
	}, nil
}

func (s *OrderService) generateOrderNumber() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("ORD-%d", timestamp)
}
