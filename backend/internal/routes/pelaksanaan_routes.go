package routes

import (
        "github.com/gin-gonic/gin"
        "github.com/telpp/emagang/internal/handler"
)

// RegisterPelaksanaanRoutes mendaftarkan route pelaksanaan magang.
//
//      peserta : /api  (role peserta)
//      hrd     : /api  (role hrd/admin)
func RegisterPelaksanaanRoutes(peserta, hrd *gin.RouterGroup, h *handler.PelaksanaanHandler) {
        // Peserta
        peserta.GET("/pelaksanaan/saya", h.GetMy)
        peserta.GET("/pelaksanaan/:id/sertifikat/download", h.DownloadSertifikat)

        // HRD / Admin
        hrd.POST("/pelaksanaan/pengajuan/:pengajuan_id", h.SetJadwal)
        hrd.GET("/pelaksanaan", h.GetAll)
        hrd.PATCH("/pelaksanaan/:id/status", h.UpdateStatus)
        hrd.PATCH("/pelaksanaan/:id/pembimbing", h.UpdatePembimbing)
        hrd.PATCH("/pelaksanaan/:id/wa-pembimbing", h.UpdateWAPembimbing)
        hrd.PATCH("/pelaksanaan/:id/nilai", h.SetNilai)
        hrd.POST("/pelaksanaan/:id/sertifikat", h.UploadSertifikat)
        hrd.GET("/sertifikat", h.ListSertifikatHRD)
}
