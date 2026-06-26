package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/handler"
)

// RegisterAuthRoutes mendaftarkan route autentikasi.
//
//	public : /api  (tanpa middleware auth)
//	api    : /api  (dengan middleware auth)
func RegisterAuthRoutes(public, api *gin.RouterGroup, h *handler.AuthHandler) {
	public.POST("/auth/login", h.Login)

	api.GET("/auth/me", h.Me)
	api.POST("/auth/change-password", h.ChangePassword)
}
