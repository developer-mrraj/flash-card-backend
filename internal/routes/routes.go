package routes

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func TestDB(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Raw(`SELECT table_name FROM information_schema.tables WHERE table_schema='public'`).Rows()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var tables []string
		for rows.Next() {
			var name string
			rows.Scan(&name)
			tables = append(tables, name)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tables)
	}
}

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
	bannerHandler *handlers.BannerHandler,
) {

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

	// Serve static images robustly at /images/*
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static", "images"))
	r.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(filesDir)))

	// Test Database Connection Route
	r.Get("/test-db", TestDB(database.DB))

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

		// Public Orders Routes (Guest Checkout & Order Success Page)
		r.Post("/orders/guest", orderHandler.PlaceGuestOrder)
		r.Get("/orders/{id}/public", orderHandler.GetPublicOrder)

		// Protected Customer Routes (Orders)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth(cfg))

			r.Route("/orders", func(r chi.Router) {
				r.Post("/", orderHandler.Create)
				r.Get("/", orderHandler.ListMyOrders)
				r.Get("/{id}", orderHandler.GetMyOrder)
			})

			r.Post("/products/{id}/reviews", productHandler.AddReview)

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
		r.Post("/verify-payment", paymentHandler.VerifyPayment)

		// Public Promo, Lead & Banner Routes (no auth required)
		r.Post("/promo/validate", promoHandler.ValidatePromo)
		r.Post("/leads", leadHandler.CaptureLead)
		r.Get("/banners", bannerHandler.ListActiveBanners)
		r.Get("/banners/{slot}", bannerHandler.GetBannerBySlot)

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

				// Admin Banners
				r.Route("/banners", func(r chi.Router) {
					r.Get("/", bannerHandler.AdminListBanners)
					r.Post("/", bannerHandler.AdminCreateBanner)
					r.Put("/{id}", bannerHandler.AdminUpdateBanner)
					r.Delete("/{id}", bannerHandler.AdminDeleteBanner)
				})
			})
		})
	})
}
