package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend/internal/repository"
)

type LeadHandler struct {
	repo repository.LeadRepository
}

func NewLeadHandler(repo repository.LeadRepository) *LeadHandler {
	return &LeadHandler{repo: repo}
}

type createLeadRequest struct {
	Email string `json:"email"`
}

// CaptureLead godoc
// @Summary Capture a visitor email for marketing
// @Tags leads
// @Accept json
// @Produce json
// @Param request body createLeadRequest true "Email"
// @Success 200 {object} map[string]string
// @Router /leads [post]
func (h *LeadHandler) CaptureLead(w http.ResponseWriter, r *http.Request) {
	var req createLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" || !strings.Contains(email, "@") {
		http.Error(w, "Valid email is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(email); err != nil {
		http.Error(w, "Failed to save email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Thank you! You're on the list."})
}
