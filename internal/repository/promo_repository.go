package repository

import (
	"backend/internal/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PromoRepository interface {
	FindByCode(code string) (*models.PromoCode, error)
}

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) PromoRepository {
	return &promoRepository{db: db}
}

func (r *promoRepository) FindByCode(code string) (*models.PromoCode, error) {
	var promo models.PromoCode
	if err := r.db.Where("code = ? AND is_active = true", strings.ToUpper(code)).First(&promo).Error; err != nil {
		return nil, err
	}
	// Check expiry
	if promo.ExpiresAt != nil && promo.ExpiresAt.Before(time.Now()) {
		return nil, gorm.ErrRecordNotFound
	}
	return &promo, nil
}
