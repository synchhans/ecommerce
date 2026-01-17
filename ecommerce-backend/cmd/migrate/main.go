package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/synchhans/ecommerce-backend/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate up")
	}

	if os.Args[1] != "up" {
		log.Fatal("only 'up' command is supported")
	}

	cfg := config.Load()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("db connect error:", err)
	}
	defer conn.Close(ctx)

	ensureSchemaTable(ctx, conn)
	applyMigrations(ctx, conn)
}

func ensureSchemaTable(ctx context.Context, conn *pgx.Conn) {
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		log.Fatal("ensure schema_migrations failed:", err)
	}
}

func applyMigrations(ctx context.Context, conn *pgx.Conn) {
	files, err := os.ReadDir("migrations")
	if err != nil {
		log.Fatal("read migrations dir failed:", err)
	}

	var migrations []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migrations = append(migrations, f.Name())
		}
	}

	sort.Strings(migrations)

	for _, m := range migrations {
		if alreadyApplied(ctx, conn, m) {
			continue
		}

		log.Println("applying", m)

		sqlFile := filepath.Join("migrations", m)
		sql, err := os.ReadFile(sqlFile)
		if err != nil {
			log.Fatal(err)
		}

		tx, err := conn.Begin(ctx)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			tx.Rollback(ctx)
			log.Fatalf("migration %s failed: %v", m, err)
		}

		if _, err := tx.Exec(ctx,
			"INSERT INTO schema_migrations (version) VALUES ($1)",
			m,
		); err != nil {
			tx.Rollback(ctx)
			log.Fatal(err)
		}

		if err := tx.Commit(ctx); err != nil {
			log.Fatal(err)
		}

		fmt.Println("✓", m)
	}
}

func alreadyApplied(ctx context.Context, conn *pgx.Conn, version string) bool {
	var exists bool
	err := conn.QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version=$1)",
		version,
	).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}
