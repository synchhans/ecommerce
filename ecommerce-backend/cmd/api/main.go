package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/synchhans/ecommerce-backend/internal/config"
	"github.com/synchhans/ecommerce-backend/internal/module/address"
	"github.com/synchhans/ecommerce-backend/internal/module/cart"
	"github.com/synchhans/ecommerce-backend/internal/module/catalog"
	"github.com/synchhans/ecommerce-backend/internal/module/inventory"
	"github.com/synchhans/ecommerce-backend/internal/module/order"
	"github.com/synchhans/ecommerce-backend/internal/module/payment"
	"github.com/synchhans/ecommerce-backend/internal/module/user"
	"github.com/synchhans/ecommerce-backend/internal/platform/database"
	httpx "github.com/synchhans/ecommerce-backend/internal/platform/http"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// Database
	pg, err := database.New(ctx, database.Config{
		DSN: cfg.DatabaseDSN,
	})
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer pg.Close()

	// Router
	router := httpx.NewRouter()
	r := router.Mux()

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// ======================
	// Wire dependencies
	// ======================

	// Catalog
	catalogHandler := catalog.NewHandler(
		catalog.NewService(
			catalog.NewPostgresRepository(pg.Pool),
		),
	)

	// Cart
	cartHandler := cart.NewHandler(
		cart.NewService(
			cart.NewPostgresRepository(pg.Pool),
		),
	)

	// Order
	orderHandler := order.NewHandler(
		order.NewService(
			order.NewPostgresRepository(pg.Pool),
		),
	)

	// Payment
	paymentHandler := payment.NewHandler(
		payment.NewService(
			payment.NewPostgresRepository(pg.Pool),
		),
	)

	// Inventory
	inventoryHandler := inventory.NewHandler(
		inventory.NewService(
			inventory.NewPostgresRepository(pg.Pool),
		),
	)

	// User
	userHandler := user.NewHandler(
		user.NewService(
			user.NewPostgresRepository(pg.Pool),
			cfg.JWTSecret,
		),
		cfg.JWTSecret,
	)

	// Address (protected)
	addressHandler := address.NewHandler(
		address.NewService(
			address.NewPostgresRepository(pg.Pool),
		),
	)

	// ======================
	// Routes
	// ======================

	r.Route("/v1", func(v1 chi.Router) {
		// Public
		catalogHandler.Routes(v1)
		cartHandler.Routes(v1)
		orderHandler.Routes(v1)
		paymentHandler.Routes(v1)
		inventoryHandler.Routes(v1)
		userHandler.Routes(v1)

		// Protected
		v1.Group(func(pr chi.Router) {
			pr.Use(httpx.AuthMiddleware([]byte(cfg.JWTSecret)))
			addressHandler.Routes(pr)
		})
	})

	// ======================
	// HTTP Server
	// ======================

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("API listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
