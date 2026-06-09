package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telpp/emagang/internal/middleware"
	"github.com/telpp/emagang/internal/models"
	"github.com/telpp/emagang/internal/repository"
)

type LandingHandler struct {
	repo *repository.LandingRepository
}

func NewLandingHandler(repo *repository.LandingRepository) *LandingHandler {
	return &LandingHandler{repo: repo}
}

// GET /api/landing/content — publik, ambil semua konten
func (h *LandingHandler) GetContent(c *gin.Context) {
	list, err := h.repo.GetAllContent(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}
	// Convert ke map untuk kemudahan frontend
	result := make(map[string]string)
	for _, item := range list {
		if item.Nilai != nil {
			result[item.Kunci] = *item.Nilai
		}
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Data: result})
}

// PUT /api/landing/content — HRD/Admin update konten
func (h *LandingHandler) UpdateContent(c *gin.Context) {
	var req struct {
		Kunci string `json:"kunci" binding:"required"`
		Nilai string `json:"nilai" binding:"required"`
		Tipe  string `json:"tipe"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	if req.Tipe == "" {
		req.Tipe = "text"
	}

	updatedBy := middleware.GetUserID(c)
	if err := h.repo.UpsertContent(c.Request.Context(), req.Kunci, req.Nilai, req.Tipe, updatedBy); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Konten berhasil diperbarui"})
}

// GET /api/landing/faq — publik
func (h *LandingHandler) GetFAQ(c *gin.Context) {
	list, err := h.repo.GetAllFAQ(c.Request.Context(), true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// GET /api/landing/faq/all — HRD/Admin (termasuk nonaktif)
func (h *LandingHandler) GetFAQAll(c *gin.Context) {
	list, err := h.repo.GetAllFAQ(c.Request.Context(), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// POST /api/landing/faq — HRD/Admin tambah/edit FAQ
func (h *LandingHandler) UpsertFAQ(c *gin.Context) {
	var f models.FAQ
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	if err := h.repo.UpsertFAQ(c.Request.Context(), &f); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "FAQ berhasil disimpan", Data: f})
}

// DELETE /api/landing/faq/:id
func (h *LandingHandler) DeleteFAQ(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
		return
	}
	h.repo.DeleteFAQ(c.Request.Context(), id)
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "FAQ berhasil dihapus"})
}

// GET /api/landing/periode — publik
func (h *LandingHandler) GetPeriodeAktif(c *gin.Context) {
	p, err := h.repo.GetActivePeriode(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, models.SuccessResponse{Data: nil})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Data: p})
}

// GET /api/landing/periode/all — Admin
func (h *LandingHandler) GetAllPeriode(c *gin.Context) {
	list, err := h.repo.GetAllPeriode(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// POST /api/landing/periode — Admin buat periode
func (h *LandingHandler) CreatePeriode(c *gin.Context) {
	var req struct {
		Nama           string `json:"nama" binding:"required"`
		TanggalBuka    string `json:"tanggal_buka" binding:"required"`
		TanggalTutup   string `json:"tanggal_tutup" binding:"required"`
		Kuota          *int   `json:"kuota"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}

	buka, err := parseTimeStr(req.TanggalBuka)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_date", Message: "Format tanggal tidak valid"})
		return
	}
	tutup, err := parseTimeStr(req.TanggalTutup)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_date", Message: "Format tanggal tidak valid"})
		return
	}

	p := &models.PeriodeMagang{
		Nama:         req.Nama,
		TanggalBuka:  buka,
		TanggalTutup: tutup,
		Kuota:        req.Kuota,
	}

	if err := h.repo.CreatePeriode(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Periode berhasil dibuat", Data: p})
}

// PATCH /api/landing/periode/:id/aktif — Admin aktifkan/nonaktifkan periode
func (h *LandingHandler) SetPeriodeAktif(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
		return
	}
	var req struct {
		IsActive bool `json:"is_active"`
	}
	c.ShouldBindJSON(&req)

	h.repo.SetPeriodeActive(c.Request.Context(), id, req.IsActive)
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Status periode diperbarui"})
}
