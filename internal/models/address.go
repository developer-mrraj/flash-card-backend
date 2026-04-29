package models

import (
	"time"

	"github.com/google/uuid"
)

type UserAddress struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"                              json:"user_id"`
	FullName  string    `gorm:"not null"                                        json:"full_name"`
	Phone     string    `gorm:"not null"                                        json:"phone"`
	Address   string    `gorm:"not null"                                        json:"address"`
	City      string    `gorm:"not null"                                        json:"city"`
	State     string    `gorm:"not null"                                        json:"state"`
	Pincode   string    `gorm:"not null"                                        json:"pincode"`
	IsDefault bool      `gorm:"default:false"                                   json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
