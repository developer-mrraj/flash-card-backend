package repository

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BannerRepository interface {
	FindAll() ([]models.Banner, error)
	FindActiveBySlot(slot string) (*models.Banner, error)
	FindAllActive() ([]models.Banner, error)
	Create(b *models.Banner) error
	Update(id uuid.UUID, b *models.Banner) error
	Delete(id uuid.UUID) error
}

type bannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) BannerRepository {
	return &bannerRepository{db: db}
}

func (r *bannerRepository) FindAll() ([]models.Banner, error) {
	var banners []models.Banner
	if err := r.db.Order("slot, sort_order ASC").Find(&banners).Error; err != nil {
		return nil, err
	}
	return banners, nil
}

func (r *bannerRepository) FindAllActive() ([]models.Banner, error) {
	var banners []models.Banner
	if err := r.db.Where("is_active = true").Order("slot, sort_order ASC").Find(&banners).Error; err != nil {
		return nil, err
	}
	return banners, nil
}

func (r *bannerRepository) FindActiveBySlot(slot string) (*models.Banner, error) {
	var banner models.Banner
	if err := r.db.Where("slot = ? AND is_active = true", slot).
		Order("sort_order ASC").
		First(&banner).Error; err != nil {
		return nil, err
	}
	return &banner, nil
}

func (r *bannerRepository) Create(b *models.Banner) error {
	return r.db.Create(b).Error
}

func (r *bannerRepository) Update(id uuid.UUID, b *models.Banner) error {
	return r.db.Model(&models.Banner{}).Where("id = ?", id).Updates(b).Error
}

func (r *bannerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Banner{}, "id = ?", id).Error
}
