package routes

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(
	r *chi.Mux,
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	productHandler *handlers.ProductHandler,
	orderHandler *handlers.OrderHandler,
) {

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

	r.Route("/api", func(r chi.Router) {

		// Public Auth Routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", authHandler.Signup)
			r.Post("/login", authHandler.Login)

			// Protected Auth Routes
			r.With(middleware.RequireAuth(cfg)).Get("/me", authHandler.Me)
		})

		// Public Product Routes
		r.Route("/products", func(r chi.Router) {
			r.Get("/", productHandler.List)
			r.Get("/{id}", productHandler.Get)
		})

		// Protected Customer Routes (Orders)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth(cfg))

			r.Route("/orders", func(r chi.Router) {
				r.Post("/", orderHandler.Create)
				r.Get("/", orderHandler.ListMyOrders)
				r.Get("/{id}", orderHandler.GetMyOrder)
			})
		})

		// Protected Admin Routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth(cfg))
			r.Use(middleware.RequireAdmin)

			r.Route("/admin", func(r chi.Router) {
				// Admin Products
				r.Route("/products", func(r chi.Router) {
					r.Post("/", productHandler.Create)
					r.Put("/{id}", productHandler.Update)
					r.Delete("/{id}", productHandler.Delete)
				})

				// Admin Orders
				r.Route("/orders", func(r chi.Router) {
					r.Get("/", orderHandler.ListAllOrders)
					r.Patch("/{id}/status", orderHandler.UpdateStatus)
				})
			})
		})
	})
}
