package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/handler"
)

// RegisterDivisiRoutes mendaftarkan route divisi.
//
//	api   : /api        (semua user login — untuk dropdown HRD)
//	admin : /api/admin  (admin only — CRUD)
func RegisterDivisiRoutes(api, admin *gin.RouterGroup, h *handler.DivisiHandler) {
	// Semua user login — dropdown list divisi aktif
	api.GET("/divisi", h.GetAktif)

	// Admin only
	admin.GET("/divisi", h.GetAll)
	admin.POST("/divisi", h.Create)
	admin.PATCH("/divisi/:id", h.Update)
	admin.PATCH("/divisi/:id/toggle", h.Toggle)
	admin.DELETE("/divisi/:id", h.Delete)
}
