package router

import (
	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/handler"
	"github.com/telpp/emagang/internal/middleware"
	"github.com/telpp/emagang/internal/models"
)

// SetupPengajuanRoutes mendaftarkan semua routes terkait pengajuan
func SetupPengajuanRoutes(api *gin.RouterGroup, pengajuanH *handler.PengajuanHandler) {
	// --------------------------------------------------------
	// PESERTA ROUTES - Pengajuan
	// --------------------------------------------------------
	peserta := api.Group("/peserta")
	peserta.Use(middleware.RoleRequired(models.RolePeserta))
	{
		// Pengajuan magang
		peserta.POST("/pengajuan", pengajuanH.Submit)
		peserta.GET("/pengajuan/saya", pengajuanH.GetMy)
	}

	// --------------------------------------------------------
	// HRD ROUTES - Pengajuan
	// --------------------------------------------------------
	hrd := api.Group("/hrd")
	hrd.Use(middleware.RoleRequired(models.RoleHRD, models.RoleAdmin))
	{
		// Pengajuan
		hrd.GET("/pengajuan", pengajuanH.GetAll)
		hrd.GET("/pengajuan/:id", pengajuanH.GetDetail)
		hrd.PATCH("/pengajuan/:id/status", pengajuanH.UpdateStatus)
	}

	// --------------------------------------------------------
	// SHARED ROUTES - Pengajuan
	// --------------------------------------------------------
	api.GET("/pengajuan/:id/history", pengajuanH.GetHistory)
}
