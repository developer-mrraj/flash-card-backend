package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"backend/internal/service"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

type VerifyPaymentRequest struct {
	RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
	RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
	RazorpaySignature string `json:"razorpay_signature" binding:"required"`
	Amount            int64  `json:"amount" binding:"required"` // amount in paise to cross-check or store
}

// VerifyPayment godoc
// @Summary Verify razorpay payment
// @Description Verify the payment signature returned by the Razorpay frontend
// @Tags payments
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body VerifyPaymentRequest true "Verify Payment Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /verify-payment [post]
func (h *PaymentHandler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	var req VerifyPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.paymentService.VerifyPaymentSignature(req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature); err != nil {
		http.Error(w, "Invalid payment signature", http.StatusBadRequest)
		return
	}

	if err := h.paymentService.ProcessPaymentSuccess(req.RazorpayPaymentID, req.RazorpayOrderID, req.Amount); err != nil {
		http.Error(w, "Failed to process payment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Payment verified successfully"})
}

// RazorpayWebhook godoc
// @Summary Razorpay server-to-server webhook
// @Description Endpoint to receive payment events from Razorpay directly
// @Tags payments
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /razorpay-webhook [post]
func (h *PaymentHandler) RazorpayWebhook(w http.ResponseWriter, r *http.Request) {
	signature := r.Header.Get("X-Razorpay-Signature")
	if signature == "" {
		http.Error(w, "Missing signature", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	if err := h.paymentService.VerifyWebhookSignature(body, signature); err != nil {
		http.Error(w, "Invalid webhook signature", http.StatusUnauthorized)
		return
	}

	var payload service.RazorpayWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if payload.Event == "payment.captured" {
		entity := payload.Payload.Payment.Entity
		// Use idempotency process
		if err := h.paymentService.ProcessPaymentSuccess(entity.ID, entity.OrderID, entity.Amount); err != nil {
			http.Error(w, "Failed to process payment webhook: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
