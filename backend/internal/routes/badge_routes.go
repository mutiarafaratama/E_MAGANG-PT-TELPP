package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/handler"
)

func RegisterBadgeRoutes(peserta, hrd *gin.RouterGroup, h *handler.BadgeHandler) {
	peserta.GET("/badge/peserta", h.GetPeserta)
	hrd.GET("/badge/hrd", h.GetHRD)
}
