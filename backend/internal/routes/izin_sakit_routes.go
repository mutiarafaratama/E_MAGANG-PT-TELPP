package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/handler"
)

// RegisterIzinSakitRoutes mendaftarkan route pengajuan izin dan sakit.
//
//	peserta : /api  (role peserta)
//	hrd     : /api  (role hrd/admin)
func RegisterIzinSakitRoutes(peserta, hrd *gin.RouterGroup, h *handler.IzinSakitHandler) {
	// Peserta — ajukan & lihat riwayat milik sendiri
	peserta.POST("/izin-sakit", h.Ajukan)
	peserta.GET("/izin-sakit/saya", h.GetSaya)

	// HRD / Admin — lihat semua, proses approve/tolak
	hrd.GET("/izin-sakit", h.GetAll)
	hrd.PATCH("/izin-sakit/:id/approve", h.Approve)
	hrd.PATCH("/izin-sakit/:id/tolak", h.Tolak)
}
