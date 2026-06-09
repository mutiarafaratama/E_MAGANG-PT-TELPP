package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	JWTSecret     string
	JWTExpiry     string
	RefreshExpiry string
	Port          string
	UploadDir     string
	MaxUploadSize int64
	AppEnv        string
	FrontendURL   string
}

var App *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	maxSize, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "104857600"), 10, 64)
	if err != nil {
		maxSize = 104857600 // 100MB default
	}

	App = &Config{
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiry:     getEnv("JWT_EXPIRY", "24h"),
		RefreshExpiry: getEnv("REFRESH_EXPIRY", "168h"),
		Port:          getEnv("PORT", "8080"),
		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSize: maxSize,
		AppEnv:        getEnv("APP_ENV", "development"),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:5173"),
	}

	if App.DatabaseURL == "" {
		log.Fatal("❌ DATABASE_URL wajib diisi di file .env")
	}
	if App.JWTSecret == "" {
		log.Fatal("❌ JWT_SECRET wajib diisi di file .env")
	}

	// Buat folder uploads jika belum ada
	if err := os.MkdirAll(App.UploadDir, 0755); err != nil {
		log.Fatalf("❌ Gagal membuat folder uploads: %v", err)
	}

	log.Printf("✅ Konfigurasi berhasil dimuat (env: %s)", App.AppEnv)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
