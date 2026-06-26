package routes

import (
        "github.com/gin-gonic/gin"
        "github.com/telpp/emagang/internal/handler"
)

// RegisterDokumenRoutes mendaftarkan route dokumen.
//
//      public  : /api  (tanpa auth)
//      peserta : /api  (role peserta)
//      api     : /api  (semua role login — shared)
//      hrd     : /api  (role hrd/admin)
func RegisterDokumenRoutes(public, peserta, api, hrd *gin.RouterGroup, h *handler.DokumenHandler) {
        // Publik — upload dari form daftar tanpa login
        public.POST("/dokumen/upload-publik", h.UploadPublik)

        // Shared — upload (peserta dan hrd/admin bisa upload, dikontrol di handler)
        api.POST("/dokumen/upload", h.Upload)

        // Shared — download (akses dikontrol di dalam handler berdasarkan role)
        api.GET("/dokumen/:id/download", h.Download)

        // Semua role login — ownership peserta dicek di dalam handler
        api.GET("/dokumen/pengajuan/:id", h.GetByPengajuan)
}
