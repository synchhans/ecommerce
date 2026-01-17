package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

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
	_ = godotenv.Load()

	ctx := context.Background()

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is required")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	pg, err := database.New(ctx, database.Config{DSN: dsn})
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer pg.Close()

	router := httpx.NewRouter()
	r := router.Mux()

	// Health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Wire modules
	cRepo := catalog.NewPostgresRepository(pg.Pool)
	cSvc := catalog.NewService(cRepo)
	cHandler := catalog.NewHandler(cSvc)

	cartRepo := cart.NewPostgresRepository(pg.Pool)
	cartSvc := cart.NewService(cartRepo)
	cartHandler := cart.NewHandler(cartSvc)

	orderRepo := order.NewPostgresRepository(pg.Pool)
	orderSvc := order.NewService(orderRepo)
	orderHandler := order.NewHandler(orderSvc)

	payRepo := payment.NewPostgresRepository(pg.Pool)
	paySvc := payment.NewService(payRepo)
	payHandler := payment.NewHandler(paySvc)

	invRepo := inventory.NewPostgresRepository(pg.Pool)
	invSvc := inventory.NewService(invRepo)
	invHandler := inventory.NewHandler(invSvc)

	userRepo := user.NewPostgresRepository(pg.Pool)
	userSvc := user.NewService(userRepo, jwtSecret)
	userHandler := user.NewHandler(userSvc, jwtSecret)

	addrRepo := address.NewPostgresRepository(pg.Pool)
	addrSvc := address.NewService(addrRepo)
	addrHandler := address.NewHandler(addrSvc)

	// v1 routes
	r.Route("/v1", func(v1 chi.Router) {
		// Public
		cHandler.Routes(v1)
		cartHandler.Routes(v1)
		orderHandler.Routes(v1)
		payHandler.Routes(v1)
		invHandler.Routes(v1)
		userHandler.Routes(v1) // register/login are public; /me is protected inside handler

		// Protected group (addresses)
		v1.Group(func(pr chi.Router) {
			pr.Use(httpx.AuthMiddleware([]byte(jwtSecret)))
			addrHandler.Routes(pr)
		})
	})

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
