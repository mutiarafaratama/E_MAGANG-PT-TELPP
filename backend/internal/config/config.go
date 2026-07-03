package config

import (
        "log"
        "os"
        "strconv"

        "github.com/joho/godotenv"
)

type Config struct {
        DatabaseURL                string
        JWTSecret                  string
        JWTExpiry                  string
        RefreshExpiry              string
        Port                       string
        UploadDir                  string
        MaxUploadSize              int64
        AppEnv                     string
        FrontendURL                string
        FirebaseProjectID          string
        FirebaseServiceAccountJSON string
}

var App *Config

func Load() {
        if err := godotenv.Load(); err != nil {
                log.Println("⚠️  File .env tidak ditemukan, menggunakan environment variables sistem")
        }

        maxSize, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "104857600"), 10, 64)
        if err != nil {
                maxSize = 104857600
        }

        App = &Config{
                DatabaseURL:   getEnv("DATABASE_URL", ""),
                JWTSecret:     getEnv("JWT_SECRET", ""),
                JWTExpiry:     getEnv("JWT_EXPIRY", "24h"),
                RefreshExpiry: getEnv("REFRESH_EXPIRY", "168h"),
                Port:          getEnv("GO_PORT", "8080"),
                UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
                MaxUploadSize: maxSize,
                AppEnv:        getEnv("APP_ENV", "development"),
                FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:5174"),
                FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", ""),
                FirebaseServiceAccountJSON: loadFirebaseServiceAccount(),
        }

        if App.DatabaseURL == "" {
                log.Fatal("❌ DATABASE_URL wajib diisi di file .env")
        }
        if App.JWTSecret == "" {
                log.Fatal("❌ JWT_SECRET wajib diisi di file .env")
        }

        if err := os.MkdirAll(App.UploadDir, 0755); err != nil {
                log.Fatalf("❌ Gagal membuat folder uploads: %v", err)
        }

        if App.FirebaseProjectID != "" && App.FirebaseServiceAccountJSON != "" {
                log.Printf("✅ Firebase FCM dikonfigurasi untuk project: %s", App.FirebaseProjectID)
        } else {
                log.Println("⚠️  Firebase FCM tidak dikonfigurasi — push notif dinonaktifkan")
        }

        log.Printf("✅ Konfigurasi berhasil dimuat (env: %s)", App.AppEnv)
}

// loadFirebaseServiceAccount mencoba memuat service account JSON dari:
// 1. FIREBASE_SERVICE_ACCOUNT_PATH — path ke file serviceAccountKey.json (untuk lokal)
// 2. FIREBASE_SERVICE_ACCOUNT_JSON — isi JSON langsung sebagai env var (untuk Replit/production)
func loadFirebaseServiceAccount() string {
        // Prioritas 1: dari file path (cocok untuk lokal dengan serviceAccountKey.json)
        if path := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH"); path != "" {
                data, err := os.ReadFile(path)
                if err != nil {
                        log.Printf("⚠️  Gagal baca FIREBASE_SERVICE_ACCOUNT_PATH (%s): %v", path, err)
                } else {
                        log.Printf("✅ Service account Firebase dimuat dari file: %s", path)
                        return string(data)
                }
        }

        // Prioritas 2: dari JSON string langsung (cocok untuk Replit secrets)
        if jsonStr := os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON"); jsonStr != "" {
                log.Println("✅ Service account Firebase dimuat dari environment variable")
                return jsonStr
        }

        return ""
}

func getEnv(key, defaultValue string) string {
        if value := os.Getenv(key); value != "" {
                return value
        }
        return defaultValue
}
