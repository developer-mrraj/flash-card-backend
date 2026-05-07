package handlers

import (
	"encoding/json"
	"net/http"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type BannerHandler struct {
	bannerRepo repository.BannerRepository
}

func NewBannerHandler(bannerRepo repository.BannerRepository) *BannerHandler {
	return &BannerHandler{bannerRepo: bannerRepo}
}

// ListActiveBanners godoc
// @Summary Get all active banners
// @Description Returns all active banners (used by frontend to pick by slot)
// @Tags banners
// @Produce json
// @Success 200 {array} models.Banner
// @Router /banners [get]
func (h *BannerHandler) ListActiveBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := h.bannerRepo.FindAllActive()
	if err != nil {
		http.Error(w, "Failed to load banners", http.StatusInternalServerError)
		return
	}
	if banners == nil {
		banners = []models.Banner{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

// GetBannerBySlot godoc
// @Summary Get active banner for a specific slot
// @Description Returns the active banner for a given slot (home_mid, categories_top, home_pre_footer)
// @Tags banners
// @Produce json
// @Param slot path string true "Banner Slot"
// @Success 200 {object} models.Banner
// @Failure 404 {object} map[string]string
// @Router /banners/{slot} [get]
func (h *BannerHandler) GetBannerBySlot(w http.ResponseWriter, r *http.Request) {
	slot := chi.URLParam(r, "slot")
	banner, err := h.bannerRepo.FindActiveBySlot(slot)
	if err != nil {
		http.Error(w, "Banner not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banner)
}

// AdminListBanners godoc
// @Summary Admin: List all banners
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Banner
// @Router /admin/banners [get]
func (h *BannerHandler) AdminListBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := h.bannerRepo.FindAll()
	if err != nil {
		http.Error(w, "Failed to load banners", http.StatusInternalServerError)
		return
	}
	if banners == nil {
		banners = []models.Banner{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

// AdminCreateBanner godoc
// @Summary Admin: Create a banner
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.Banner true "Banner"
// @Success 201 {object} models.Banner
// @Router /admin/banners [post]
func (h *BannerHandler) AdminCreateBanner(w http.ResponseWriter, r *http.Request) {
	var b models.Banner
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.bannerRepo.Create(&b); err != nil {
		http.Error(w, "Failed to create banner", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

// AdminUpdateBanner godoc
// @Summary Admin: Update a banner
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Banner ID"
// @Param request body models.Banner true "Banner"
// @Success 200 {object} models.Banner
// @Router /admin/banners/{id} [put]
func (h *BannerHandler) AdminUpdateBanner(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid banner ID", http.StatusBadRequest)
		return
	}
	var b models.Banner
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.bannerRepo.Update(id, &b); err != nil {
		http.Error(w, "Failed to update banner", http.StatusInternalServerError)
		return
	}
	b.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

// AdminDeleteBanner godoc
// @Summary Admin: Delete a banner
// @Tags admin
// @Security BearerAuth
// @Param id path string true "Banner ID"
// @Success 204 "No Content"
// @Router /admin/banners/{id} [delete]
func (h *BannerHandler) AdminDeleteBanner(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid banner ID", http.StatusBadRequest)
		return
	}
	if err := h.bannerRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete banner", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
