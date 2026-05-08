package repository

import (
	"backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderTransaction(order *models.Order, updates []models.Product) error
	FindByUserID(userID uuid.UUID) ([]models.Order, error)
	FindByIDAndUserID(id, userID uuid.UUID) (*models.Order, error)
	FindAll() ([]models.Order, error)
	FindByID(id uuid.UUID) (*models.Order, error)
	UpdateStatus(id uuid.UUID, status string) error
	UpdatePaymentDetails(id uuid.UUID, rzpOrderID string, rzpPaymentID string, status string) error
	FindByRazorpayOrderID(rzpOrderID string) (*models.Order, error)
	LinkGuestOrdersByEmail(email string, userID uuid.UUID) error
	LinkGuestOrdersByPhone(phone string, userID uuid.UUID) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrderTransaction(order *models.Order, updates []models.Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, product := range updates {
			if err := tx.Model(&models.Product{}).Where("id = ?", product.ID).Update("stock_quantity", product.StockQuantity).Error; err != nil {
				return err
			}
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *orderRepository) FindByUserID(userID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Preload("Items.Product").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) LinkGuestOrdersByEmail(email string, userID uuid.UUID) error {
	return r.db.Model(&models.Order{}).
		Where("user_id IS NULL AND guest_email = ?", email).
		Update("user_id", userID).Error
}

func (r *orderRepository) LinkGuestOrdersByPhone(phone string, userID uuid.UUID) error {
	return r.db.Model(&models.Order{}).
		Where("user_id IS NULL AND guest_phone = ?", phone).
		Update("user_id", userID).Error
}

func (r *orderRepository) FindByIDAndUserID(id, userID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("Items.Product").Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Preload("User").Preload("Items.Product").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) FindByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("Items.Product").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) UpdatePaymentDetails(id uuid.UUID, rzpOrderID string, rzpPaymentID string, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Updates(map[string]interface{}{
		"razorpay_order_id":   rzpOrderID,
		"razorpay_payment_id": rzpPaymentID,
		"status":              status,
	}).Error
}

func (r *orderRepository) FindByRazorpayOrderID(rzpOrderID string) (*models.Order, error) {
	var order models.Order
	if err := r.db.Where("razorpay_order_id = ?", rzpOrderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
