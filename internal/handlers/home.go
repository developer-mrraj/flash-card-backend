package handlers

import (
	"encoding/json"
	"net/http"

	"backend/internal/service"
	_ "backend/internal/dto"
)

type HomeHandler struct {
	homeService service.HomeService
}

func NewHomeHandler(homeService service.HomeService) *HomeHandler {
	return &HomeHandler{homeService: homeService}
}

// GetHome godoc
// @Summary Get home page content
// @Description Get aggregated home page content including hero cards, featured collections, and bestselling products
// @Tags home
// @Produce  json
// @Success 200 {object} dto.HomeResponse
// @Failure 500 {object} map[string]string
// @Router /home [get]
func (h *HomeHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	res, err := h.homeService.GetHomeContent()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
