package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID          uuid.UUID   `gorm:"type:uuid;primary_key;" json:"id"`
	UserID      uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TotalAmount int64       `gorm:"not null" json:"total_amount"` // Stored in cents
	Status      string      `gorm:"type:varchar(50);default:'pending'" json:"status"` // 'pending', 'paid', 'shipped', 'delivered', 'cancelled'
	Items       []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	UnitPrice int64     `gorm:"not null" json:"unit_price"` // Price at time of purchase in cents
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	oi.ID = uuid.New()
	return
}
