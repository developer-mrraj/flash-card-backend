package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type FeaturedCollectionRepository interface {
	FindAll() ([]models.FeaturedCollection, error)
}

type featuredCollectionRepository struct {
	db *gorm.DB
}

func NewFeaturedCollectionRepository(db *gorm.DB) FeaturedCollectionRepository {
	return &featuredCollectionRepository{db: db}
}

func (r *featuredCollectionRepository) FindAll() ([]models.FeaturedCollection, error) {
	var collections []models.FeaturedCollection
	if err := r.db.Find(&collections).Error; err != nil {
		return nil, err
	}
	return collections, nil
}
