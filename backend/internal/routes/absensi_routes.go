package routes

import (
        "github.com/gin-gonic/gin"
        "github.com/telpp/emagang/internal/handler"
)

// RegisterAbsensiRoutes mendaftarkan route absensi harian.
//
//      peserta : /api  (role peserta)
//      hrd     : /api  (role hrd/admin)
func RegisterAbsensiRoutes(peserta, hrd *gin.RouterGroup, h *handler.AbsensiHandler) {
        // Peserta — self-service absensi
        peserta.GET("/absensi/geo-status", h.GeoStatus)
        peserta.POST("/absensi/checkin", h.CheckIn)
        peserta.PATCH("/absensi/checkout", h.CheckOut)
        peserta.GET("/absensi/saya", h.GetMy)
        peserta.GET("/absensi/saya/pdf", h.DownloadPDF)

        // HRD / Admin — rekap & approve
        hrd.GET("/absensi/rekap", h.GetRekapHRD)
        hrd.GET("/absensi/pelaksanaan/:id", h.GetByPelaksanaan)
        hrd.GET("/absensi/pelaksanaan/:id/pdf", h.GetByPelaksanaanPDF)
        hrd.PATCH("/absensi/:id/approve", h.Approve)
}

// RegisterAbsensiConfigRoutes mendaftarkan route konfigurasi jam absensi.
//
//      api   : /api  (semua user terautentikasi)
//      admin : /api/admin  (admin only)
func RegisterAbsensiConfigRoutes(api, admin *gin.RouterGroup, h *handler.AbsensiConfigHandler) {
        api.GET("/absensi/config", h.Get)
        admin.PUT("/absensi/config", h.Update)
}
