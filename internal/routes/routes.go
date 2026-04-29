package routes

import (
	"net/http"
	"os"
	"path/filepath"

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
	paymentHandler *handlers.PaymentHandler,
	addressHandler *handlers.AddressHandler,
	promoHandler *handlers.PromoHandler,
	leadHandler *handlers.LeadHandler,
	homeHandler *handlers.HomeHandler,
) {

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

	// Serve static images robustly at /images/*
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static", "images"))
	r.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(filesDir)))

	r.Route("/api", func(r chi.Router) {

		r.Get("/home", homeHandler.GetHome)

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
			r.Get("/{id}/reviews", productHandler.GetReviews)
		})

		// Protected Customer Routes (Orders)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth(cfg))

			r.Route("/orders", func(r chi.Router) {
				r.Post("/", orderHandler.Create)
				r.Get("/", orderHandler.ListMyOrders)
				r.Get("/{id}", orderHandler.GetMyOrder)
			})

			r.Post("/products/{id}/reviews", productHandler.AddReview)

			r.Post("/verify-payment", paymentHandler.VerifyPayment)

			// Address Routes
			r.Route("/addresses", func(r chi.Router) {
				r.Get("/", addressHandler.GetMyAddresses)
				r.Post("/", addressHandler.AddAddress)
				r.Put("/{id}/default", addressHandler.SetDefault)
				r.Delete("/{id}", addressHandler.DeleteAddress)
			})
		})

		// Public Webhook Route
		r.Post("/razorpay-webhook", paymentHandler.RazorpayWebhook)

		// Public Promo & Lead Routes (no auth required)
		r.Post("/promo/validate", promoHandler.ValidatePromo)
		r.Post("/leads", leadHandler.CaptureLead)

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
