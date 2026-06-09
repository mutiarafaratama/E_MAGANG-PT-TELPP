package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/telpp/emagang/internal/config"
)

var DB *pgxpool.Pool

func Connect() {
	poolConfig, err := pgxpool.ParseConfig(config.App.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Gagal parse DATABASE_URL: %v", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke database: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ Database tidak bisa diping: %v", err)
	}

	DB = pool
	log.Println("✅ Koneksi database PostgreSQL berhasil")
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("🔌 Koneksi database ditutup")
	}
}
