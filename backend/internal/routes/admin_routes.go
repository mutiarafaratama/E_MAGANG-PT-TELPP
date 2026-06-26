package routes

import (
        "github.com/gin-gonic/gin"
        "github.com/telpp/emagang/internal/handler"
)

// RegisterAdminRoutes mendaftarkan route admin.
//
//      hrd   : /api        (role hrd/admin — untuk endpoint shared seperti hrd-list)
//      admin : /api/admin  (role admin only)
func RegisterAdminRoutes(hrd, admin *gin.RouterGroup, h *handler.AdminHandler, dokumenH *handler.DokumenHandler) {
        // HRD / Admin — daftar HRD untuk keperluan assign pembimbing
        hrd.GET("/admin/hrd-list", h.GetHRDList)

        // Admin only (prefix /admin sudah ada di group)
        admin.GET("/stats", h.GetStats)
        admin.GET("/users", h.GetUsers)
        admin.POST("/users", h.CreateUser)
        admin.PATCH("/users/:id", h.UpdateUser)
        admin.PATCH("/users/:id/toggle", h.ToggleUser)
        admin.DELETE("/users/:id", h.HapusPeserta)

        // Kelola dokumen (admin only)
        admin.GET("/dokumen", dokumenH.AdminList)
        admin.GET("/dokumen/:id/download", dokumenH.AdminDownload)
        admin.DELETE("/dokumen/:id", dokumenH.AdminHapus)

        // Pantau akun yang akan dihapus otomatis
        admin.GET("/cleanup/preview", h.GetAkunJatuhTempo)
}
