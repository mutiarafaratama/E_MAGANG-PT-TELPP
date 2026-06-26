package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/handler"
)

// RegisterPerpanjanganRoutes mendaftarkan route perpanjangan durasi magang.
//
//	peserta : POST /api/perpanjangan         — ajukan perpanjangan
//	peserta : GET  /api/perpanjangan/saya    — lihat status pengajuan sendiri
//	hrd     : GET  /api/perpanjangan         — semua pengajuan
//	hrd     : POST /api/perpanjangan/langsung — perpanjang langsung tanpa perlu approval
//	hrd     : PATCH /api/perpanjangan/:id/approve — setujui pengajuan peserta
//	hrd     : PATCH /api/perpanjangan/:id/tolak   — tolak pengajuan peserta
func RegisterPerpanjanganRoutes(peserta, hrd *gin.RouterGroup, h *handler.PerpanjanganHandler) {
	peserta.POST("/perpanjangan", h.AjukanPeserta)
	peserta.GET("/perpanjangan/saya", h.GetSaya)

	hrd.GET("/perpanjangan", h.GetAll)
	hrd.POST("/perpanjangan/langsung", h.AjukanHRD)
	hrd.PATCH("/perpanjangan/:id/approve", h.Approve)
	hrd.PATCH("/perpanjangan/:id/tolak", h.Tolak)
}
