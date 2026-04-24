package handlers

import (
	"encoding/json"
	"net/http"

	"backend/internal/dto"
	"backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// List godoc
// @Summary List all products
// @Description Get a list of all available products
// @Tags products
// @Produce  json
// @Success 200 {array} dto.ProductResponse
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	res, err := h.productService.GetAll()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Get godoc
// @Summary Get a product by ID
// @Description Get detailed information about a specific product
// @Tags products
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	res, err := h.productService.GetByID(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Create godoc
// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body dto.ProductRequest true "Product Request"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	res, err := h.productService.Create(req)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// Update godoc
// @Summary Update an existing product
// @Description Update product details (Admin only)
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body dto.ProductRequest true "Product Request"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	res, err := h.productService.Update(id, req)
	if err != nil {
		http.Error(w, "Failed to update product or not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Delete godoc
// @Summary Delete a product
// @Description Remove a product from the catalog (Admin only)
// @Tags admin
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.productService.Delete(id); err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
