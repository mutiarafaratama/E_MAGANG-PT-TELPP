package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/handler"
)

// RegisterLaporanRoutes mendaftarkan route laporan magang.
//
//	peserta : /api  (role peserta)
//	hrd     : /api  (role hrd/admin)
//	api     : /api  (semua user login)
func RegisterLaporanRoutes(peserta, hrd, api *gin.RouterGroup, h *handler.LaporanHandler) {
	// Peserta
	peserta.POST("/laporan/upload", h.Upload)
	peserta.GET("/laporan/saya", h.GetMy)

	// HRD / Admin
	hrd.GET("/laporan", h.GetAll)
	hrd.GET("/laporan/pelaksanaan/:pelaksanaan_id", h.GetByPelaksanaan)
	hrd.PATCH("/laporan/:id/review", h.Review)

	// Semua user login (download)
	api.GET("/laporan/:id/download", h.Download)
}
