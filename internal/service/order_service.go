package service

import (
	"errors"
	"log"
	"strings"

	"backend/internal/dto"
	"backend/internal/models"
	"backend/internal/repository"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(userID uuid.UUID, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
	GetMyOrders(userID uuid.UUID) ([]dto.OrderResponse, error)
	GetMyOrder(userID, orderID uuid.UUID) (*dto.OrderResponse, error)
	GetAllOrders() ([]dto.OrderResponse, error)
	UpdateOrderStatus(orderID uuid.UUID, req dto.UpdateOrderStatusRequest) error
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	paymentSvc  PaymentService
	promoRepo   repository.PromoRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, paymentSvc PaymentService, promoRepo repository.PromoRepository) OrderService {
	return &orderService{orderRepo: orderRepo, productRepo: productRepo, paymentSvc: paymentSvc, promoRepo: promoRepo}
}

func mapToOrderResponse(o *models.Order) dto.OrderResponse {
	res := dto.OrderResponse{
		ID:              o.ID,
		UserID:          o.UserID,
		TotalAmount:     o.TotalAmount,
		Status:          o.Status,
		RazorpayOrderID: o.RazorpayOrderID,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
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

func (s *orderService) CreateOrder(userID uuid.UUID, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	var totalAmount int64
	var orderItems []models.OrderItem

	// Validate stock and calculate total — do NOT deduct stock yet.
	// Stock is only deducted in ProcessPaymentSuccess after confirmed payment.
	for _, itemReq := range req.Items {
		product, err := s.productRepo.FindByID(itemReq.ProductID)
		if err != nil {
			return nil, errors.New("product not found")
		}

		if product.StockQuantity < itemReq.Quantity {
			return nil, errors.New("not enough stock")
		}

		itemTotal := product.Price * int64(itemReq.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ID,
			Quantity:  itemReq.Quantity,
			UnitPrice: product.Price,
		})
	}

	// Apply promo discount server-side (validates code, calculates discount)
	var discountAmount int64
	if code := strings.TrimSpace(req.PromoCode); code != "" {
		promo, err := s.promoRepo.FindByCode(code)
		if err == nil && promo != nil {
			discountAmount = totalAmount * int64(promo.DiscountPercentage) / 100
			totalAmount -= discountAmount
			log.Printf("[Promo] Applied %s: -%d paise (%.0f%% of original)", promo.Code, discountAmount, float64(promo.DiscountPercentage))
		} else {
			log.Printf("[Promo] Code '%s' not found or expired, skipping discount", code)
		}
	}

	order := &models.Order{
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      "pending",
		Items:       orderItems,
	}

	// Create the order record only — stock is NOT deducted here
	if err := s.orderRepo.CreateOrderTransaction(order, nil); err != nil {
		return nil, err
	}

	// Generate Razorpay Order with the discounted totalAmount
	receipt := order.ID.String()
	rzpOrderID, rzpErr := s.paymentSvc.GenerateRazorpayOrder(order.TotalAmount, receipt)
	if rzpErr != nil {
		log.Printf("[Razorpay] WARN: Failed to create Razorpay order for order %s: %v", order.ID, rzpErr)
	} else if rzpOrderID != "" {
		s.orderRepo.UpdatePaymentDetails(order.ID, rzpOrderID, "", "pending")
		order.RazorpayOrderID = rzpOrderID
		log.Printf("[Razorpay] Order created: %s for order %s (amount: %d paise)", rzpOrderID, order.ID, order.TotalAmount)
	}

	res := mapToOrderResponse(order)
	res.DiscountAmount = discountAmount
	return &res, nil
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
