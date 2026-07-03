package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/models"
)

// ── Security Headers ─────────────────────────────────────────────────────────
// Pasang di semua response API untuk mencegah serangan umum (XSS, clickjacking,
// MIME-sniffing, dll).

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(self)")
		c.Header("X-Download-Options", "noopen")
		c.Next()
	}
}

// ── Rate Limiter (Token Bucket, per IP) ──────────────────────────────────────
// Dipakai untuk membatasi percobaan login agar tidak bisa di-brute-force.

type ipLimiter struct {
	mu      sync.Mutex
	clients map[string]*bucket
}

type bucket struct {
	tokens    float64
	maxTokens float64
	rate      float64 // token per detik
	lastCheck time.Time
}

func newIPLimiter() *ipLimiter {
	il := &ipLimiter{clients: make(map[string]*bucket)}
	go il.cleanup()
	return il
}

// cleanup hapus entri IP yang sudah lama tidak aktif (hemat memory)
func (il *ipLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		il.mu.Lock()
		cutoff := time.Now().Add(-30 * time.Minute)
		for ip, b := range il.clients {
			if b.lastCheck.Before(cutoff) {
				delete(il.clients, ip)
			}
		}
		il.mu.Unlock()
	}
}

func (il *ipLimiter) allow(ip string, maxReq int, window time.Duration) bool {
	il.mu.Lock()
	defer il.mu.Unlock()

	now := time.Now()
	b, ok := il.clients[ip]
	if !ok {
		il.clients[ip] = &bucket{
			tokens:    float64(maxReq) - 1,
			maxTokens: float64(maxReq),
			rate:      float64(maxReq) / window.Seconds(),
			lastCheck: now,
		}
		return true
	}

	// Isi ulang token sesuai waktu yang berlalu
	elapsed := now.Sub(b.lastCheck).Seconds()
	b.tokens += elapsed * b.rate
	if b.tokens > b.maxTokens {
		b.tokens = b.maxTokens
	}
	b.lastCheck = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

var loginLimiter = newIPLimiter()

// RateLimitLogin — maks 10 percobaan per 15 menit per IP address.
// Jika melewati batas → 429 Too Many Requests.
func RateLimitLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !loginLimiter.allow(c.ClientIP(), 10, 15*time.Minute) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, models.ErrorResponse{
				Error:   "rate_limit_exceeded",
				Message: "Terlalu banyak percobaan login. Coba lagi dalam beberapa menit.",
			})
			return
		}
		c.Next()
	}
}
