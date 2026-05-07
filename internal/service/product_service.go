package service

import (
	"backend/internal/dto"
	"backend/internal/models"
	"backend/internal/repository"
	"github.com/google/uuid"
)

type ProductService interface {
	GetAll() ([]dto.ProductResponse, error)
	GetByID(id uuid.UUID) (*dto.ProductResponse, error)
	Create(req dto.ProductRequest) (*dto.ProductResponse, error)
	Update(id uuid.UUID, req dto.ProductRequest) (*dto.ProductResponse, error)
	Delete(id uuid.UUID) error
	AddReview(productID uuid.UUID, userID uuid.UUID, req dto.ReviewRequest) (*dto.ReviewResponse, error)
	GetReviews(productID uuid.UUID) ([]dto.ReviewResponse, error)
}

type productService struct {
	repo       repository.ProductRepository
	reviewRepo repository.ReviewRepository
}

func NewProductService(repo repository.ProductRepository, reviewRepo repository.ReviewRepository) ProductService {
	return &productService{repo: repo, reviewRepo: reviewRepo}
}

func mapToProductResponse(p *models.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Slug:          p.Slug,
		Title:         p.Title,
		Badge:         p.Badge,
		BadgeClass:    p.BadgeClass,
		Rating:        p.Rating,
		Reviews:       p.Reviews,
		OriginalPrice: p.OriginalPrice,
		Discount:      p.Discount,
		MainImage:     p.MainImage,
		GalleryImages: []string(p.GalleryImages),
		CardsCount:    p.CardsCount,
		Features:      p.Features,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func (s *productService) GetAll() ([]dto.ProductResponse, error) {
	products, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []dto.ProductResponse
	for _, p := range products {
		res = append(res, mapToProductResponse(&p))
	}
	return res, nil
}

func (s *productService) GetByID(id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	res := mapToProductResponse(product)
	return &res, nil
}

func (s *productService) Create(req dto.ProductRequest) (*dto.ProductResponse, error) {
	product := &models.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Slug:          req.Slug,
		Title:         req.Title,
		Badge:         req.Badge,
		BadgeClass:    req.BadgeClass,
		Rating:        req.Rating,
		Reviews:       req.Reviews,
		OriginalPrice: req.OriginalPrice,
		Discount:      req.Discount,
		MainImage:     req.MainImage,
		GalleryImages: req.GalleryImages,
		CardsCount:    req.CardsCount,
		Features:      req.Features,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	res := mapToProductResponse(product)
	return &res, nil
}

func (s *productService) Update(id uuid.UUID, req dto.ProductRequest) (*dto.ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.StockQuantity = req.StockQuantity
	product.Slug = req.Slug
	product.Title = req.Title
	product.Badge = req.Badge
	product.BadgeClass = req.BadgeClass
	product.Rating = req.Rating
	product.Reviews = req.Reviews
	product.OriginalPrice = req.OriginalPrice
	product.Discount = req.Discount
	product.MainImage = req.MainImage
	product.GalleryImages = req.GalleryImages
	product.CardsCount = req.CardsCount
	if req.Features != nil {
		product.Features = req.Features
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	res := mapToProductResponse(product)
	return &res, nil
}

func (s *productService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *productService) AddReview(productID uuid.UUID, userID uuid.UUID, req dto.ReviewRequest) (*dto.ReviewResponse, error) {
	review := &models.Review{
		ProductID: productID,
		UserID:    userID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := s.reviewRepo.CreateOrUpdate(review); err != nil {
		return nil, err
	}

	// Update product rating
	avgRating, count, err := s.reviewRepo.GetAverageRating(productID)
	if err == nil {
		s.repo.UpdateRating(productID, avgRating, count)
	}

	// Re-fetch to get user details
	reviews, _ := s.reviewRepo.FindByProductID(productID)
	for _, r := range reviews {
		if r.ID == review.ID {
			review = &r
			break
		}
	}

	res := dto.ReviewResponse{
		ID:        review.ID,
		ProductID: review.ProductID,
		UserID:    review.UserID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
		User: dto.UserResponse{
			ID:   review.User.ID,
			Name: review.User.Name,
		},
	}

	return &res, nil
}

func (s *productService) GetReviews(productID uuid.UUID) ([]dto.ReviewResponse, error) {
	reviews, err := s.reviewRepo.FindByProductID(productID)
	if err != nil {
		return nil, err
	}

	var res []dto.ReviewResponse
	for _, r := range reviews {
		res = append(res, dto.ReviewResponse{
			ID:        r.ID,
			ProductID: r.ProductID,
			UserID:    r.UserID,
			Rating:    r.Rating,
			Comment:   r.Comment,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			User: dto.UserResponse{
				ID:   r.User.ID,
				Name: r.User.Name,
			},
		})
	}

	return res, nil
}
