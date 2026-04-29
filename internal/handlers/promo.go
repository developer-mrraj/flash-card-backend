package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend/internal/repository"
)

type PromoHandler struct {
	repo repository.PromoRepository
}

func NewPromoHandler(repo repository.PromoRepository) *PromoHandler {
	return &PromoHandler{repo: repo}
}

type validatePromoRequest struct {
	Code string `json:"code"`
}

type validatePromoResponse struct {
	Valid               bool   `json:"valid"`
	Code                string `json:"code"`
	DiscountPercentage  int    `json:"discount_percentage"`
	Message             string `json:"message"`
}

// ValidatePromo godoc
// @Summary Validate a promo code
// @Tags promo
// @Accept json
// @Produce json
// @Param request body validatePromoRequest true "Promo Code"
// @Success 200 {object} validatePromoResponse
// @Router /promo/validate [post]
func (h *PromoHandler) ValidatePromo(w http.ResponseWriter, r *http.Request) {
	var req validatePromoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Code == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(validatePromoResponse{Valid: false, Message: "Invalid request"})
		return
	}

	promo, err := h.repo.FindByCode(req.Code)
	if err != nil || promo == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(validatePromoResponse{Valid: false, Message: "Invalid or expired promo code"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(validatePromoResponse{
		Valid:              true,
		Code:               promo.Code,
		DiscountPercentage: promo.DiscountPercentage,
		Message:            fmt.Sprintf("🎉 Promo applied! You save %d%% on your order.", promo.DiscountPercentage),
	})
}
