package repository

import (
	"fmt"

	"backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll() ([]models.Product, error)
	FindByID(id uuid.UUID) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	UpdateRating(id uuid.UUID, rating float64, reviews int) error
	Delete(id uuid.UUID) error
	DeductStock(items []models.OrderItem) error // called after confirmed payment
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) UpdateRating(id uuid.UUID, rating float64, reviews int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Updates(map[string]interface{}{
		"rating":  rating,
		"reviews": reviews,
	}).Error
}

// DeductStock atomically reduces stock_quantity for each order item after payment is confirmed.
func (r *productRepository) DeductStock(items []models.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			result := tx.Model(&models.Product{}).Where("id = ? AND stock_quantity >= ?", item.ProductID, item.Quantity).
				UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", item.Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("insufficient stock for product %s", item.ProductID)
			}
		}
		return nil
	})
}
