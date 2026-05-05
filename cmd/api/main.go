package main

import (
	"log"
	"net/http"

	_ "backend/docs" // Import swagger docs
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/routes"
	"backend/internal/service"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// @title           SJ Flash Cards Backend API
// @version         1.0
// @description     This is the API server for the signup and ordering system.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      api.yourdomain.com
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Connect to database & migrate (using golang-migrate)
	database.Connect(cfg)
	db := database.DB

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	promoRepo := repository.NewPromoRepository(db)
	leadRepo := repository.NewLeadRepository(db)
	featuredCollectionRepo := repository.NewFeaturedCollectionRepository(db)

	// 4. Initialize Services
	authService := service.NewAuthService(userRepo, cfg)
	productService := service.NewProductService(productRepo, reviewRepo)
	paymentService := service.NewPaymentService(cfg, paymentRepo, orderRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, paymentService)
	homeService := service.NewHomeService(productRepo, featuredCollectionRepo)

	// 5. Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	addressHandler := handlers.NewAddressHandler(addressRepo)
	promoHandler := handlers.NewPromoHandler(promoRepo)
	leadHandler := handlers.NewLeadHandler(leadRepo)
	homeHandler := handlers.NewHomeHandler(homeService)

	// 6. Setup Router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 7. Register Routes
	routes.RegisterRoutes(r, cfg, authHandler, productHandler, orderHandler, paymentHandler, addressHandler, promoHandler, leadHandler, homeHandler)

	// 8. Start Server
	log.Printf("Starting server on port %s...", cfg.Port)
	log.Printf("Swagger documentation available at /swagger/index.html")
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
