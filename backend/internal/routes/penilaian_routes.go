package routes

import (
        "github.com/gin-gonic/gin"
        "github.com/telpp/emagang/internal/handler"
)

// api    = JWT required, semua role
// hrd    = JWT required, role HRD/Admin
func RegisterPenilaianRoutes(api, hrd *gin.RouterGroup, h *handler.PenilaianHandler) {
        // Semua role (peserta bisa GET nilai dirinya sendiri)
        api.GET("/penilaian/:pelaksanaan_id", h.Get)

        // Khusus HRD/Admin
        hrd.GET("/penilaian", h.List)
        hrd.POST("/penilaian/:pelaksanaan_id", h.Upsert)
}
