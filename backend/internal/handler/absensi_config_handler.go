package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/middleware"
	"github.com/telpp/emagang/internal/models"
	"github.com/telpp/emagang/internal/repository"
)

type AbsensiConfigHandler struct {
	repo *repository.AbsensiConfigRepository
}

func NewAbsensiConfigHandler(repo *repository.AbsensiConfigRepository) *AbsensiConfigHandler {
	return &AbsensiConfigHandler{repo: repo}
}

// GET /api/absensi/config — semua user terautentikasi
func (h *AbsensiConfigHandler) Get(c *gin.Context) {
	cfg, err := h.repo.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal memuat konfigurasi jam absensi"})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Data: cfg})
}

// PUT /api/admin/absensi/config — admin only
func (h *AbsensiConfigHandler) Update(c *gin.Context) {
	var req struct {
		JamMasukBuka   string `json:"jam_masuk_buka" binding:"required"`
		JamMasukTutup  string `json:"jam_masuk_tutup" binding:"required"`
		JamPulangBuka  string `json:"jam_pulang_buka" binding:"required"`
		JamPulangTutup string `json:"jam_pulang_tutup" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}

	updatedBy := middleware.GetUserID(c)
	if err := h.repo.Update(c.Request.Context(),
		req.JamMasukBuka, req.JamMasukTutup,
		req.JamPulangBuka, req.JamPulangTutup,
		updatedBy,
	); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Konfigurasi jam absensi berhasil diperbarui"})
}
