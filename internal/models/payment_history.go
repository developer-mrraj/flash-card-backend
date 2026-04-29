package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentHistory struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	OrderID           uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	Order             Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	RazorpayOrderID   string    `gorm:"type:varchar(255);not null" json:"razorpay_order_id"`
	RazorpayPaymentID string    `gorm:"type:varchar(255);unique;not null" json:"razorpay_payment_id"`
	Amount            int64     `gorm:"not null" json:"amount"`
	Status            string    `gorm:"type:varchar(50);default:'success'" json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

func (p *PaymentHistory) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}

func (PaymentHistory) TableName() string {
	return "payment_history"
}
