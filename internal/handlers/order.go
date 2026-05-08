package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"backend/internal/dto"
	"backend/internal/middleware"
	"backend/internal/service"
	"backend/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// Create godoc
// @Summary Create a new order
// @Description Place a new order with items. Automatically deducts stock.
// @Tags orders
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body dto.CreateOrderRequest true "Order Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Items) == 0 {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	res, err := h.orderService.CreateOrder(claims.UserID, req)
	if err != nil {
		log.Printf("[OrderHandler] CreateOrder error for user %s: %v", claims.UserID, err)
		http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// PlaceGuestOrder godoc
// @Summary Place a new order as guest
// @Description Place a new order with items without authentication. Automatically deducts stock.
// @Tags orders
// @Accept  json
// @Produce  json
// @Param request body dto.GuestOrderRequest true "Guest Order Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/guest [post]
func (h *OrderHandler) PlaceGuestOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.GuestOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Items) == 0 {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	res, err := h.orderService.CreateGuestOrder(req)
	if err != nil {
		log.Printf("[OrderHandler] CreateGuestOrder error: %v", err)
		http.Error(w, "Failed to create guest order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// GetPublicOrder godoc
// @Summary Get order by ID (public)
// @Description Fetch order details by ID. Used by order-success page for both guests and authenticated users.
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 404 {object} map[string]string
// @Router /orders/{id}/public [get]
func (h *OrderHandler) GetPublicOrder(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	res, err := h.orderService.GetPublicOrder(orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// ListMyOrders godoc
// @Summary Get user's orders
// @Description Get a list of all orders placed by the authenticated user
// @Tags orders
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} dto.OrderResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [get]
func (h *OrderHandler) ListMyOrders(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	res, err := h.orderService.GetMyOrders(claims.UserID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GetMyOrder godoc
// @Summary Get user's specific order
// @Description Get detailed information about a specific order belonging to the user
// @Tags orders
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func (h *OrderHandler) GetMyOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idParam := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	res, err := h.orderService.GetMyOrder(claims.UserID, orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// ListAllOrders godoc
// @Summary List all orders in the system
// @Description Get a list of all orders across the system (Admin only)
// @Tags admin
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} dto.OrderResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/orders [get]
func (h *OrderHandler) ListAllOrders(w http.ResponseWriter, r *http.Request) {
	res, err := h.orderService.GetAllOrders()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// UpdateStatus godoc
// @Summary Update an order's status
// @Description Update the status (e.g. shipped, delivered) of a specific order (Admin only)
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param request body dto.UpdateOrderStatusRequest true "Update Status Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/orders/{id}/status [patch]
func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.orderService.UpdateOrderStatus(orderID, req); err != nil {
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Order status updated"})
}
