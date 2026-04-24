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
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
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

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	res := mapToProductResponse(product)
	return &res, nil
}

func (s *productService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
