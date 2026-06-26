package handler

import (
        "net/http"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
)

type DivisiHandler struct {
        repo *repository.DivisiRepository
}

func NewDivisiHandler(repo *repository.DivisiRepository) *DivisiHandler {
        return &DivisiHandler{repo: repo}
}

// GET /api/divisi — daftar divisi aktif (untuk dropdown HRD)
func (h *DivisiHandler) GetAktif(c *gin.Context) {
        list, err := h.repo.FindAll(c.Request.Context(), true)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// GET /api/admin/divisi — semua divisi termasuk nonaktif (untuk admin)
func (h *DivisiHandler) GetAll(c *gin.Context) {
        list, err := h.repo.FindAll(c.Request.Context(), false)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// divisiRequest — body untuk Create dan Update
type divisiRequest struct {
        Nama       string   `json:"nama" binding:"required"`
        Urutan     int      `json:"urutan"`
        GeoLat     *float64 `json:"geo_lat"`
        GeoLng     *float64 `json:"geo_lng"`
        GeoRadius  *int     `json:"geo_radius"`
        NamaLokasi *string  `json:"nama_lokasi"`
}

func validateGeo(req divisiRequest) bool {
        // lat dan lng harus ada berdua atau tidak sama sekali
        return (req.GeoLat == nil) == (req.GeoLng == nil)
}

// POST /api/admin/divisi — buat divisi baru
func (h *DivisiHandler) Create(c *gin.Context) {
        var req divisiRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Nama divisi wajib diisi"})
                return
        }
        if !validateGeo(req) {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Latitude dan Longitude geofencing harus diisi bersamaan"})
                return
        }

        d, err := h.repo.Create(c.Request.Context(), req.Nama, req.Urutan, req.GeoLat, req.GeoLng, req.GeoRadius, req.NamaLokasi)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "create_failed", Message: "Nama divisi sudah ada"})
                return
        }
        c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Divisi berhasil ditambahkan", Data: d})
}

// PATCH /api/admin/divisi/:id — update nama/urutan/geofencing
func (h *DivisiHandler) Update(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        var req divisiRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Nama divisi wajib diisi"})
                return
        }
        if !validateGeo(req) {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Latitude dan Longitude geofencing harus diisi bersamaan"})
                return
        }
        if err := h.repo.Update(c.Request.Context(), id, req.Nama, req.Urutan, req.GeoLat, req.GeoLng, req.GeoRadius, req.NamaLokasi); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "update_failed", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Divisi berhasil diperbarui"})
}

// PATCH /api/admin/divisi/:id/toggle — aktif/nonaktif
func (h *DivisiHandler) Toggle(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        var req struct {
                IsActive bool `json:"is_active"`
        }
        c.ShouldBindJSON(&req)

        if err := h.repo.ToggleActive(c.Request.Context(), id, req.IsActive); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        status := "dinonaktifkan"
        if req.IsActive {
                status = "diaktifkan"
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Divisi berhasil " + status})
}

// DELETE /api/admin/divisi/:id — hapus divisi
func (h *DivisiHandler) Delete(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        if err := h.repo.Delete(c.Request.Context(), id); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "delete_failed", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Divisi berhasil dihapus"})
}
