package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"github.com/razorpay/razorpay-go"
)

type PaymentService interface {
	GenerateRazorpayOrder(amount int64, receipt string) (string, error)
	VerifyPaymentSignature(orderID, paymentID, signature string) error
	VerifyWebhookSignature(body []byte, signature string) error
	ProcessPaymentSuccess(paymentID, orderID string, amount int64) error
}

type paymentService struct {
	cfg          *config.Config
	paymentRepo  repository.PaymentRepository
	orderRepo    repository.OrderRepository
	razorClient  *razorpay.Client
}

func NewPaymentService(cfg *config.Config, paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository) PaymentService {
	client := razorpay.NewClient(cfg.RazorpayKeyID, cfg.RazorpayKeySecret)
	return &paymentService{
		cfg:         cfg,
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		razorClient: client,
	}
}

func (s *paymentService) GenerateRazorpayOrder(amount int64, receipt string) (string, error) {
	data := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  receipt,
	}
	body, err := s.razorClient.Order.Create(data, nil)
	if err != nil {
		return "", err
	}

	orderID, ok := body["id"].(string)
	if !ok {
		return "", errors.New("failed to get order ID from razorpay")
	}

	return orderID, nil
}

func (s *paymentService) VerifyPaymentSignature(orderID, paymentID, signature string) error {
	data := orderID + "|" + paymentID
	h := hmac.New(sha256.New, []byte(s.cfg.RazorpayKeySecret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if expectedSignature != signature {
		return errors.New("invalid signature")
	}
	return nil
}

func (s *paymentService) VerifyWebhookSignature(body []byte, signature string) error {
	h := hmac.New(sha256.New, []byte(s.cfg.RazorpayWebhookSecret))
	h.Write(body)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if expectedSignature != signature {
		return errors.New("invalid webhook signature")
	}
	return nil
}

func (s *paymentService) ProcessPaymentSuccess(paymentID, orderID string, amount int64) error {
	// 1. Idempotency Check
	existingPayment, err := s.paymentRepo.FindByPaymentID(paymentID)
	if err == nil && existingPayment != nil {
		// Already processed
		return nil
	}

	// 2. Find Order
	order, err := s.orderRepo.FindByRazorpayOrderID(orderID)
	if err != nil {
		return fmt.Errorf("order not found for razorpay_order_id: %s", orderID)
	}

	// Fetch exact amount from Razorpay (Optional, but recommended)
	// payment, err := s.razorClient.Payment.Fetch(paymentID, nil, nil)
	// amount_paid := payment["amount"].(float64)

	// 3. Insert Payment History
	history := &models.PaymentHistory{
		OrderID:           *order.ID,
		RazorpayOrderID:   orderID,
		RazorpayPaymentID: paymentID,
		Amount:            amount,
		Status:            "success",
	}
	if err := s.paymentRepo.CreatePaymentHistory(history); err != nil {
		return err
	}

	// 4. Update Order Status
	if err := s.orderRepo.UpdatePaymentDetails(*order.ID, orderID, paymentID, "paid"); err != nil {
		return err
	}

	return nil
}

type RazorpayWebhookPayload struct {
	Event   string `json:"event"`
	Payload struct {
		Payment struct {
			Entity struct {
				ID      string `json:"id"`
				OrderID string `json:"order_id"`
				Amount  int64  `json:"amount"`
				Status  string `json:"status"`
			} `json:"entity"`
		} `json:"payment"`
	} `json:"payload"`
}
