package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/synchhans/ecommerce-backend/internal/platform/database"
)

type Product struct {
	Name        string
	Description string
	Price       int
	Category    string
	ImageURL    string
	Stock       int
}

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is required")
	}

	ctx := context.Background()

	db, err := database.New(ctx, database.Config{DSN: dsn})
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer db.Close()

	products := []Product{
		{
			Name:        "iPhone 15 Pro",
			Description: "Apple smartphone with A17 Pro chip",
			Price:       19999000,
			Category:    "electronics",
			ImageURL:    "https://images.unsplash.com/photo-1695048133142-1a2045d3fda2",
			Stock:       20,
		},
		{
			Name:        "Mechanical Keyboard RGB",
			Description: "Hot-swappable mechanical keyboard",
			Price:       1299000,
			Category:    "electronics",
			ImageURL:    "https://images.unsplash.com/photo-1517336714731-489689fd1ca8",
			Stock:       50,
		},
		{
			Name:        "Men Oversized Hoodie",
			Description: "Premium cotton hoodie",
			Price:       349000,
			Category:    "fashion",
			ImageURL:    "https://images.unsplash.com/photo-1520975922215-36c8e0c42b97",
			Stock:       100,
		},
		{
			Name:        "Skincare Serum",
			Description: "Hydrating facial serum",
			Price:       189000,
			Category:    "beauty",
			ImageURL:    "https://images.unsplash.com/photo-1611930022073-b7a4ba5fcccd",
			Stock:       75,
		},
		{
			Name:        "Minimalist Desk Lamp",
			Description: "Warm LED desk lamp",
			Price:       499000,
			Category:    "home",
			ImageURL:    "https://images.unsplash.com/photo-1507473885765-e6ed057f782c",
			Stock:       40,
		},
	}

	for _, p := range products {
		_, err := db.Pool.Exec(ctx, `
			INSERT INTO products (
				name, description, price_cents, category, image_url, stock, created_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (name) DO NOTHING
		`,
			p.Name,
			p.Description,
			p.Price,
			p.Category,
			p.ImageURL,
			p.Stock,
			time.Now(),
		)

		if err != nil {
			log.Fatalf("insert product failed: %v", err)
		}
	}

	log.Println("✅ product seed completed")
}
