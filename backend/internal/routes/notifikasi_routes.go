package routes

import (
        "github.com/gin-gonic/gin"
        "github.com/telpp/emagang/internal/handler"
)

// RegisterNotifikasiRoutes mendaftarkan route notifikasi.
//
//      api : /api  (semua role login)
func RegisterNotifikasiRoutes(api *gin.RouterGroup, h *handler.NotifikasiHandler) {
        api.GET("/notifikasi", h.GetAll)
        api.PATCH("/notifikasi/:id/read", h.MarkRead)
        api.PATCH("/notifikasi/read-all", h.MarkAllRead)
        api.GET("/notifikasi/badge", h.GetBadge)
        api.POST("/notifikasi/fcm-token", h.SaveFCMToken)
}
