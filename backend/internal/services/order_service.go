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
