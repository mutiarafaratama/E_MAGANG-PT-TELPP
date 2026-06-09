package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/telpp/emagang/internal/config"
	"github.com/telpp/emagang/internal/database"
	"github.com/telpp/emagang/internal/router"
	ws "github.com/telpp/emagang/internal/websocket"
)

func main() {
	// 1. Load konfigurasi dari .env
	config.Load()

	// 2. Koneksi ke PostgreSQL
	database.Connect()
	defer database.Close()

	// 3. Jalankan WebSocket hub di goroutine terpisah
	go ws.GlobalHub.Run()

	// 4. Setup semua routes
	r := router.Setup()

	// 5. Jalankan server HTTP
	srv := &http.Server{
		Addr:         ":" + config.App.Port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Jalankan server di goroutine
	go func() {
		log.Printf("🚀 e-Magang TELPP API berjalan di http://localhost:%s", config.App.Port)
		log.Printf("📋 Environment: %s", config.App.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Graceful shutdown: tunggu sinyal OS (Ctrl+C atau SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Server sedang shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Gagal shutdown: %v", err)
	}

	log.Println("✅ Server berhasil shutdown")
}
