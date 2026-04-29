package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PromoCode struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	Code               string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	DiscountPercentage int        `gorm:"not null" json:"discount_percentage"` // e.g. 20 means 20%
	IsActive           bool       `gorm:"default:true" json:"is_active"`
	ExpiresAt          *time.Time `gorm:"default:null" json:"expires_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

func (p *PromoCode) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
