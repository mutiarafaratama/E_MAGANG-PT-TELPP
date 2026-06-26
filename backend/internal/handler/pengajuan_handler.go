package handler

import (
        "net/http"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/service"
)

type PengajuanHandler struct {
        svc *service.PengajuanService
}

func NewPengajuanHandler(svc *service.PengajuanService) *PengajuanHandler {
        return &PengajuanHandler{svc: svc}
}

// POST /api/pengajuan/publik — form publik tanpa login
func (h *PengajuanHandler) SubmitPublik(c *gin.Context) {
        var body struct {
                Step1 models.PengajuanStep1Request `json:"step1" binding:"required"`
                Step2 models.PengajuanStep2Request `json:"step2" binding:"required"`
        }
        if err := c.ShouldBindJSON(&body); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        p, err := h.svc.SubmitPublik(c.Request.Context(), body.Step1, body.Step2)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "submit_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: "Pengajuan berhasil dikirim! Tim HRD akan menghubungi Anda via email setelah verifikasi selesai.",
                Data:    p,
        })
}

// POST /api/hrd/pengajuan/:id/kirim-akun — HRD buat akun & kirim kredensial
func (h *PengajuanHandler) KirimAkun(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        if err := h.svc.KirimAkun(c.Request.Context(), id); err != nil {
                status := http.StatusBadRequest
                if err.Error() == "pengajuan tidak ditemukan" {
                        status = http.StatusNotFound
                }
                c.JSON(status, models.ErrorResponse{Error: "kirim_akun_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{
                Message: "Akun berhasil dibuat dan kredensial telah dikirim ke email peserta",
        })
}

// POST /api/pengajuan — peserta submit pengajuan magang
func (h *PengajuanHandler) Submit(c *gin.Context) {
        var body struct {
                Step1 models.PengajuanStep1Request `json:"step1" binding:"required"`
                Step2 models.PengajuanStep2Request `json:"step2" binding:"required"`
        }
        if err := c.ShouldBindJSON(&body); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        userID := middleware.GetUserID(c)
        p, err := h.svc.Submit(c.Request.Context(), userID, body.Step1, body.Step2)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "submit_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: "Pengajuan magang berhasil dikirim",
                Data:    p,
        })
}

// GET /api/pengajuan/saya — peserta lihat pengajuan sendiri
func (h *PengajuanHandler) GetMy(c *gin.Context) {
        userID := middleware.GetUserID(c)
        p, err := h.svc.GetMyPengajuan(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: p})
}

// GET /api/pengajuan — HRD/Admin lihat semua pengajuan
func (h *PengajuanHandler) GetAll(c *gin.Context) {
        status := c.Query("status")
        search := c.Query("search")
        page := queryInt(c, "page", 1)
        limit := queryInt(c, "limit", 10)

        list, total, err := h.svc.GetAll(c.Request.Context(), status, search, page, limit)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        totalPages := total / limit
        if total%limit != 0 {
                totalPages++
        }

        c.JSON(http.StatusOK, models.PaginatedResponse{
                Data:       list,
                Total:      total,
                Page:       page,
                Limit:      limit,
                TotalPages: totalPages,
        })
}

// GET /api/pengajuan/:id — detail pengajuan
func (h *PengajuanHandler) GetDetail(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        p, err := h.svc.GetDetail(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Pengajuan tidak ditemukan"})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Data: p})
}

// PATCH /api/pengajuan/:id/status — HRD update status
func (h *PengajuanHandler) UpdateStatus(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        var req models.UpdateStatusRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        changedBy := middleware.GetUserID(c)
        if err := h.svc.UpdateStatus(c.Request.Context(), id, req.Status, req.Catatan, changedBy); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "update_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Status pengajuan berhasil diperbarui"})
}

// GET /api/pengajuan/:id/history — riwayat status
func (h *PengajuanHandler) GetHistory(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        history, err := h.svc.GetStatusHistory(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Data: history})
}

// DELETE /api/hrd/pengajuan/:id — hapus pengajuan (HRD only)
func (h *PengajuanHandler) Hapus(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        if err := h.svc.HapusPengajuan(c.Request.Context(), id); err != nil {
                status := http.StatusInternalServerError
                if err.Error() == "pengajuan tidak ditemukan" {
                        status = http.StatusNotFound
                }
                c.JSON(status, models.ErrorResponse{Error: "hapus_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Pengajuan berhasil dihapus"})
}
