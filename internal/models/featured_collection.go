package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeaturedCollection struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Image     string    `gorm:"type:varchar(500);not null" json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (fc *FeaturedCollection) BeforeCreate(tx *gorm.DB) (err error) {
	fc.ID = uuid.New()
	return
}
