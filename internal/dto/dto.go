package dto

import (
	"time"

	"github.com/google/uuid"
)

// Auth
type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token   string `json:"token"`
	Message string `json:"message,omitempty"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Product
type ProductRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	Price         int64    `json:"price" binding:"required"`
	StockQuantity int      `json:"stock_quantity"`
	Slug          string   `json:"slug"`
	Title         string   `json:"title"`
	Badge         string   `json:"badge"`
	BadgeClass    string   `json:"badge_class"`
	Rating        float64  `json:"rating"`
	Reviews       int      `json:"reviews"`
	OriginalPrice *int64   `json:"original_price"`
	Discount      string   `json:"discount"`
	MainImage     string   `json:"main_image"`
	GalleryImages []string `json:"gallery_images"`
	CardsCount    int      `json:"cards_count"`
}

type ProductResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         int64     `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	Badge         string    `json:"badge"`
	BadgeClass    string    `json:"badge_class"`
	Rating        float64   `json:"rating"`
	Reviews       int       `json:"reviews"`
	OriginalPrice *int64    `json:"original_price,omitempty"`
	Discount      string    `json:"discount,omitempty"`
	MainImage     string    `json:"main_image"`
	GalleryImages []string  `json:"gallery_images"`
	CardsCount    int       `json:"cards_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Order
type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,dive"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type OrderResponse struct {
	ID          uuid.UUID             `json:"id"`
	UserID      uuid.UUID             `json:"user_id"`
	User        *UserResponse         `json:"user,omitempty"`
	TotalAmount int64                 `json:"total_amount"`
	Status      string                `json:"status"`
	Items       []OrderItemResponse   `json:"items"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type OrderItemResponse struct {
	ID        uuid.UUID        `json:"id"`
	ProductID uuid.UUID        `json:"product_id"`
	Product   *ProductResponse `json:"product,omitempty"`
	Quantity  int              `json:"quantity"`
	UnitPrice int64            `json:"unit_price"`
}
