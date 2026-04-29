package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Lead struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (l *Lead) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New()
	return
}
