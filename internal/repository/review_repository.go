package repository

import (
	"backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	CreateOrUpdate(review *models.Review) error
	FindByProductID(productID uuid.UUID) ([]models.Review, error)
	GetAverageRating(productID uuid.UUID) (float64, int, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) CreateOrUpdate(review *models.Review) error {
	// Upsert review (user can only have 1 review per product)
	return r.db.Where(models.Review{UserID: review.UserID, ProductID: review.ProductID}).
		Assign(models.Review{Rating: review.Rating, Comment: review.Comment}).
		FirstOrCreate(review).Error
}

func (r *reviewRepository) FindByProductID(productID uuid.UUID) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.db.Preload("User").Where("product_id = ?", productID).Order("created_at DESC").Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) GetAverageRating(productID uuid.UUID) (float64, int, error) {
	var result struct {
		AvgRating float64
		Count     int
	}
	err := r.db.Model(&models.Review{}).
		Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as count").
		Where("product_id = ?", productID).
		Scan(&result).Error

	return result.AvgRating, result.Count, err
}
