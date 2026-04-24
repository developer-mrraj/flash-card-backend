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
			if err := tx.Save(&product).Error; err != nil {
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
	if err := r.db.Preload("Items.Product").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
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
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}
