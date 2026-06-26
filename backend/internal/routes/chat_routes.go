package routes

import (
        "github.com/gin-gonic/gin"
        "github.com/telpp/emagang/internal/handler"
)

// RegisterChatRoutes mendaftarkan route chat helpdesk.
//
//      peserta : /api  (role peserta)
//      api     : /api  (semua role login — shared)
//      hrd     : /api  (role hrd/admin)
func RegisterChatRoutes(peserta, api, hrd *gin.RouterGroup, h *handler.ChatHandler) {
        // Peserta — buat & lihat tiket milik sendiri
        peserta.POST("/chat/tiket", h.BuatTiket)
        peserta.GET("/chat/tiket/saya", h.GetTiketSaya)

        // Shared — baca/kirim pesan
        api.GET("/chat/tiket/:id/pesan", h.GetPesan)
        api.POST("/chat/tiket/:id/pesan", h.KirimPesan)

        // HRD / Admin — kelola tiket
        hrd.GET("/chat/tiket", h.GetAllTiket)
        hrd.PATCH("/chat/tiket/:id/ambil", h.AmbilTiket)
        hrd.PATCH("/chat/tiket/:id/status", h.UpdateStatus)
}
