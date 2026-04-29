package repository

import (
	"backend/internal/models"
	"strings"

	"gorm.io/gorm"
)

type LeadRepository interface {
	Create(email string) error
}

type leadRepository struct {
	db *gorm.DB
}

func NewLeadRepository(db *gorm.DB) LeadRepository {
	return &leadRepository{db: db}
}

func (r *leadRepository) Create(email string) error {
	lead := &models.Lead{Email: strings.ToLower(strings.TrimSpace(email))}
	// On conflict (duplicate email) do nothing - not an error
	result := r.db.Where(models.Lead{Email: lead.Email}).FirstOrCreate(lead)
	return result.Error
}
