package service

import (
	"errors"

	"backend/internal/dto"
	"backend/internal/models"
	"backend/internal/repository"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(userID uuid.UUID, req dto.CreateOrderRequest) error
	GetMyOrders(userID uuid.UUID) ([]dto.OrderResponse, error)
	GetMyOrder(userID, orderID uuid.UUID) (*dto.OrderResponse, error)
	GetAllOrders() ([]dto.OrderResponse, error)
	UpdateOrderStatus(orderID uuid.UUID, req dto.UpdateOrderStatusRequest) error
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) OrderService {
	return &orderService{orderRepo: orderRepo, productRepo: productRepo}
}

func mapToOrderResponse(o *models.Order) dto.OrderResponse {
	res := dto.OrderResponse{
		ID:          o.ID,
		UserID:      o.UserID,
		TotalAmount: o.TotalAmount,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}

	if o.User.ID != uuid.Nil {
		res.User = &dto.UserResponse{
			ID:        o.User.ID,
			Name:      o.User.Name,
			Email:     o.User.Email,
			Role:      o.User.Role,
			CreatedAt: o.User.CreatedAt,
			UpdatedAt: o.User.UpdatedAt,
		}
	}

	for _, item := range o.Items {
		itemRes := dto.OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
		if item.Product.ID != uuid.Nil {
			pRes := dto.ProductResponse{
				ID:            item.Product.ID,
				Name:          item.Product.Name,
				Description:   item.Product.Description,
				Price:         item.Product.Price,
				StockQuantity: item.Product.StockQuantity,
				CreatedAt:     item.Product.CreatedAt,
				UpdatedAt:     item.Product.UpdatedAt,
			}
			itemRes.Product = &pRes
		}
		res.Items = append(res.Items, itemRes)
	}

	return res
}

func (s *orderService) CreateOrder(userID uuid.UUID, req dto.CreateOrderRequest) error {
	var totalAmount int64
	var orderItems []models.OrderItem
	var productsToUpdate []models.Product

	for _, itemReq := range req.Items {
		product, err := s.productRepo.FindByID(itemReq.ProductID)
		if err != nil {
			return errors.New("product not found")
		}

		if product.StockQuantity < itemReq.Quantity {
			return errors.New("not enough stock")
		}

		product.StockQuantity -= itemReq.Quantity
		productsToUpdate = append(productsToUpdate, *product)

		itemTotal := product.Price * int64(itemReq.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ID,
			Quantity:  itemReq.Quantity,
			UnitPrice: product.Price,
		})
	}

	order := &models.Order{
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      "pending",
		Items:       orderItems,
	}

	return s.orderRepo.CreateOrderTransaction(order, productsToUpdate)
}

func (s *orderService) GetMyOrders(userID uuid.UUID) ([]dto.OrderResponse, error) {
	orders, err := s.orderRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var res []dto.OrderResponse
	for _, o := range orders {
		res = append(res, mapToOrderResponse(&o))
	}
	return res, nil
}

func (s *orderService) GetMyOrder(userID, orderID uuid.UUID) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByIDAndUserID(orderID, userID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	res := mapToOrderResponse(order)
	return &res, nil
}

func (s *orderService) GetAllOrders() ([]dto.OrderResponse, error) {
	orders, err := s.orderRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []dto.OrderResponse
	for _, o := range orders {
		res = append(res, mapToOrderResponse(&o))
	}
	return res, nil
}

func (s *orderService) UpdateOrderStatus(orderID uuid.UUID, req dto.UpdateOrderStatusRequest) error {
	if err := s.orderRepo.UpdateStatus(orderID, req.Status); err != nil {
		return err
	}
	return nil
}
