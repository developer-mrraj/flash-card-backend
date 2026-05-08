package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AddressHandler struct {
	repo repository.AddressRepository
}

func NewAddressHandler(repo repository.AddressRepository) *AddressHandler {
	return &AddressHandler{repo: repo}
}

type createAddressRequest struct {
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Pincode      string `json:"pincode"`
}

// getUserID extracts the authenticated user's UUID from the request context.
func getUserIDFromCtx(r *http.Request) (uuid.UUID, bool) {
	val := r.Context().Value(middleware.UserContextKey)
	if val == nil {
		log.Printf("getUserIDFromCtx: UserContextKey not found in context")
		return uuid.Nil, false
	}
	claims, ok := val.(*utils.JWTClaims)
	if !ok || claims == nil {
		log.Printf("getUserIDFromCtx: Failed to type assert context value to *utils.JWTClaims. Got: %T", val)
		return uuid.Nil, false
	}
	return claims.UserID, true
}

// GetMyAddresses godoc
// @Summary Get saved addresses
// @Tags addresses
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.UserAddress
// @Router /addresses [get]
func (h *AddressHandler) GetMyAddresses(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromCtx(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	addresses, err := h.repo.FindByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch addresses", http.StatusInternalServerError)
		return
	}

	if addresses == nil {
		addresses = []models.UserAddress{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addresses)
}

// AddAddress godoc
// @Summary Add a new delivery address
// @Tags addresses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 201 {object} models.UserAddress
// @Router /addresses [post]
func (h *AddressHandler) AddAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromCtx(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req createAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	address := &models.UserAddress{
		UserID:       userID,
		FullName:     req.FullName,
		Phone:        req.Phone,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		Pincode:      req.Pincode,
	}

	if err := h.repo.Create(address); err != nil {
		http.Error(w, "Failed to save address", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(address)
}

// SetDefaultAddress godoc
// @Summary Set an address as default
// @Tags addresses
// @Security BearerAuth
// @Router /addresses/{id}/default [put]
func (h *AddressHandler) SetDefault(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromCtx(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	addressID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.SetDefault(userID, addressID); err != nil {
		http.Error(w, "Failed to update default", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Default address updated"}`))
}

// DeleteAddress godoc
// @Summary Delete an address
// @Tags addresses
// @Security BearerAuth
// @Router /addresses/{id} [delete]
func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromCtx(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	addressID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(userID, addressID); err != nil {
		http.Error(w, "Failed to delete address", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Address deleted"}`))
}
