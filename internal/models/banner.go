package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Banner represents a promotional advertisement banner stored in the DB.
// Each banner is tied to a "slot" (placement position on the frontend).
type Banner struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"      json:"id"`
	Slot      string    `gorm:"type:varchar(100);not null;index" json:"slot"`       // e.g. home_mid, categories_top, home_pre_footer
	Title     string    `gorm:"type:varchar(255);not null"   json:"title"`
	Subtitle  string    `gorm:"type:varchar(500)"            json:"subtitle"`
	CtaText   string    `gorm:"type:varchar(100)"            json:"cta_text"`
	CtaLink   string    `gorm:"type:varchar(500)"            json:"cta_link"`
	ImageURL  string    `gorm:"type:varchar(1000)"           json:"image_url"`
	BgColor   string    `gorm:"type:varchar(255);default:'#0f0f23'" json:"bg_color"`
	TextColor string    `gorm:"type:varchar(50);default:'#ffffff'" json:"text_color"`
	BadgeText string    `gorm:"type:varchar(100)"            json:"badge_text"`
	IsActive  bool      `gorm:"default:true"                 json:"is_active"`
	SortOrder int       `gorm:"default:0"                    json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Banner) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
