package routes

import (
        "github.com/gin-gonic/gin"
        "github.com/telpp/emagang/internal/handler"
        "github.com/telpp/emagang/internal/middleware"
)

// RegisterAuthRoutes mendaftarkan route autentikasi.
//
//      public : /api  (tanpa middleware auth)
//      api    : /api  (dengan middleware auth)
func RegisterAuthRoutes(public, api *gin.RouterGroup, h *handler.AuthHandler) {
        // Rate limit login: maks 10 percobaan per 15 menit per IP
        public.POST("/auth/login", middleware.RateLimitLogin(), h.Login)
        public.POST("/auth/forgot-password", h.ForgotPassword)
        public.POST("/auth/reset-password", h.ResetPassword)

        api.GET("/auth/me", h.Me)
        api.POST("/auth/change-password", h.ChangePassword)
}
