package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/handler"
)

// RegisterPengajuanRoutes mendaftarkan route pengajuan magang.
//
//	public  : /api            (tanpa auth)
//	peserta : /api            (role peserta)
//	api     : /api            (semua role login — shared)
//	hrd     : /api            (role hrd/admin)
func RegisterPengajuanRoutes(public, peserta, api, hrd *gin.RouterGroup, h *handler.PengajuanHandler) {
	// Publik — form daftar tanpa login
	public.POST("/pengajuan/publik", h.SubmitPublik)

	// Peserta
	peserta.POST("/pengajuan", h.Submit)
	peserta.GET("/pengajuan/saya", h.GetMy)

	// Shared (semua role yang login)
	api.GET("/pengajuan/:id/history", h.GetHistory)

	// HRD / Admin
	hrd.GET("/pengajuan", h.GetAll)
	hrd.GET("/pengajuan/:id", h.GetDetail)
	hrd.PATCH("/pengajuan/:id/status", h.UpdateStatus)
	hrd.POST("/pengajuan/:id/kirim-akun", h.KirimAkun)
	hrd.DELETE("/pengajuan/:id", h.Hapus)
}
