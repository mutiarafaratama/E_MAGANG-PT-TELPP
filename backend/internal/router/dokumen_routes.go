package router

import (
	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/handler"
	"github.com/telpp/emagang/internal/middleware"
	"github.com/telpp/emagang/internal/models"
)

// SetupDokumenRoutes mendaftarkan semua routes terkait dokumen
func SetupDokumenRoutes(api *gin.RouterGroup, dokumenH *handler.DokumenHandler) {
	// --------------------------------------------------------
	// PESERTA ROUTES - Dokumen
	// --------------------------------------------------------
	peserta := api.Group("/peserta")
	peserta.Use(middleware.RoleRequired(models.RolePeserta))
	{
		// Dokumen upload (peserta upload berkas pengajuan)
		peserta.POST("/dokumen/upload", dokumenH.Upload)
	}

	// --------------------------------------------------------
	// HRD ROUTES - Dokumen
	// --------------------------------------------------------
	hrd := api.Group("/hrd")
	hrd.Use(middleware.RoleRequired(models.RoleHRD, models.RoleAdmin))
	{
		// Dokumen
		hrd.POST("/dokumen/upload", dokumenH.Upload)
		hrd.GET("/dokumen/pengajuan/:id", dokumenH.GetByPengajuan)
	}

	// --------------------------------------------------------
	// SHARED ROUTES - Dokumen
	// --------------------------------------------------------
	api.GET("/dokumen/:id/download", dokumenH.Download)
}
