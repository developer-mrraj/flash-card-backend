package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID          *uuid.UUID  `gorm:"type:uuid;primary_key;"         json:"id"`
	UserID      *uuid.UUID  `gorm:"type:uuid"                      json:"user_id,omitempty"` // nullable for guest orders
	User        User        `gorm:"foreignKey:UserID"              json:"user,omitempty"`
	TotalAmount int64       `gorm:"not null"                       json:"total_amount"` // in paise
	Status            string      `gorm:"type:varchar(50);default:'pending'" json:"status"`
	RazorpayOrderID   string      `gorm:"type:varchar(255)"              json:"razorpay_order_id,omitempty"`
	RazorpayPaymentID string      `gorm:"type:varchar(255)"              json:"razorpay_payment_id,omitempty"`
	// Guest order fields (populated when user is not logged in)
	GuestName  string `gorm:"type:varchar(255)" json:"guest_name,omitempty"`
	GuestEmail string `gorm:"type:varchar(255)" json:"guest_email,omitempty"`
	GuestPhone string `gorm:"type:varchar(20)"  json:"guest_phone,omitempty"`
	Items             []OrderItem `gorm:"foreignKey:OrderID"             json:"items"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null"     json:"order_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"     json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID"   json:"product,omitempty"`
	Quantity  int       `gorm:"not null"               json:"quantity"`
	UnitPrice int64     `gorm:"not null"               json:"unit_price"` // Price at time of purchase in paise
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New()
	o.ID = &id
	return
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	oi.ID = uuid.New()
	return
}
