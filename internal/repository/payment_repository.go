package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePaymentHistory(history *models.PaymentHistory) error
	FindByPaymentID(paymentID string) (*models.PaymentHistory, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreatePaymentHistory(history *models.PaymentHistory) error {
	return r.db.Create(history).Error
}

func (r *paymentRepository) FindByPaymentID(paymentID string) (*models.PaymentHistory, error) {
	var history models.PaymentHistory
	if err := r.db.Where("razorpay_payment_id = ?", paymentID).First(&history).Error; err != nil {
		return nil, err
	}
	return &history, nil
}
