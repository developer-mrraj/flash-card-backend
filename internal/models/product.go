package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;"     json:"id"`
	Name          string         `gorm:"type:varchar(255);not null"  json:"name"`
	Description   string         `gorm:"type:text"                   json:"description"`
	Price         int64          `gorm:"not null"                    json:"price"`          // in paise / cents
	StockQuantity int            `gorm:"not null;default:0"          json:"stock_quantity"`
	// Deck-specific fields
	Slug          string         `gorm:"type:varchar(255);uniqueIndex" json:"slug"`
	Title         string         `gorm:"type:varchar(500)"           json:"title"`
	Badge         string         `gorm:"type:varchar(100)"           json:"badge"`
	BadgeClass    string         `gorm:"type:varchar(100)"           json:"badge_class"`
	Rating        float64        `gorm:"type:numeric(3,1);default:0" json:"rating"`
	Reviews       int            `gorm:"default:0"                   json:"reviews"`
	OriginalPrice *int64         `gorm:"default:null"                json:"original_price,omitempty"`
	Discount      string         `gorm:"type:varchar(50)"            json:"discount,omitempty"`
	MainImage     string         `gorm:"type:varchar(500)"           json:"main_image"`
	GalleryImages pq.StringArray    `gorm:"type:text[]"                 json:"gallery_images"`
	CardsCount    int               `gorm:"default:0"                   json:"cards_count"`
	Features      []byte            `gorm:"type:jsonb;default:'[]'"     json:"features"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
