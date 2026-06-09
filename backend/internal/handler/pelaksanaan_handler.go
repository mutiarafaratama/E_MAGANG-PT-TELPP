package handler

import (
        "net/http"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/telpp/emagang/internal/database"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
        "github.com/telpp/emagang/internal/service"
)

type PelaksanaanHandler struct {
        repo            *repository.PelaksanaanRepository
        pengajuanRepo   *repository.PengajuanRepository
        sertifikatSvc   *service.SertifikatService
}

func NewPelaksanaanHandler(
        repo *repository.PelaksanaanRepository,
        pengajuanRepo *repository.PengajuanRepository,
        sertifikatSvc *service.SertifikatService,
) *PelaksanaanHandler {
        return &PelaksanaanHandler{repo: repo, pengajuanRepo: pengajuanRepo, sertifikatSvc: sertifikatSvc}
}

// POST /api/pelaksanaan — HRD set jadwal magang setelah diterima
func (h *PelaksanaanHandler) SetJadwal(c *gin.Context) {
        var req models.SetJadwalRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        pengajuanID, err := uuid.Parse(c.Param("pengajuan_id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID pengajuan tidak valid"})
                return
        }

        // Ambil pengajuan
        pengajuan, err := h.pengajuanRepo.FindByID(c.Request.Context(), pengajuanID)
        if err != nil || pengajuan.Status != models.StatusDiterima {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_status", Message: "Pengajuan belum berstatus diterima"})
                return
        }

        tglMulai, err := time.Parse("2006-01-02", req.TanggalMulai)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_date", Message: "Format tanggal mulai tidak valid"})
                return
        }
        tglSelesai, err := time.Parse("2006-01-02", req.TanggalSelesai)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_date", Message: "Format tanggal selesai tidak valid"})
                return
        }

        p := &models.PelaksanaanMagang{
                PengajuanID:    pengajuanID,
                UserID:         pengajuan.UserID,
                TanggalMulai:   tglMulai,
                TanggalSelesai: tglSelesai,
                Divisi:         &req.Divisi,
        }

        if req.PembimbingID != "" {
                pid, err := uuid.Parse(req.PembimbingID)
                if err == nil {
                        p.PembimbingID = &pid
                }
        }

        if err := h.repo.Create(c.Request.Context(), p); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: "Jadwal magang berhasil ditetapkan",
                Data:    p,
        })
}

// GET /api/pelaksanaan/saya — peserta lihat pelaksanaan sendiri
func (h *PelaksanaanHandler) GetMy(c *gin.Context) {
        userID := middleware.GetUserID(c)
        p, err := h.repo.FindByUserID(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusOK, models.SuccessResponse{Data: nil})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: p})
}

// GET /api/pelaksanaan — HRD/Admin lihat semua
func (h *PelaksanaanHandler) GetAll(c *gin.Context) {
        status := c.Query("status")
        page := queryInt(c, "page", 1)
        limit := queryInt(c, "limit", 10)

        list, total, err := h.repo.FindAll(c.Request.Context(), status, page, limit)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.PaginatedResponse{
                Data: list, Total: total, Page: page, Limit: limit,
        })
}

// PATCH /api/pelaksanaan/:id/nilai — HRD set nilai
func (h *PelaksanaanHandler) SetNilai(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        var req models.NilaiRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        changedBy := middleware.GetUserID(c)
        if err := h.repo.SetNilai(c.Request.Context(), id, req.Nilai, req.CatatanNilai, changedBy); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Nilai berhasil disimpan"})
}

// POST /api/pelaksanaan/:id/sertifikat — generate sertifikat
func (h *PelaksanaanHandler) GenerateSertifikat(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        path, err := h.sertifikatSvc.Generate(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "generate_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{
                Message: "Sertifikat berhasil digenerate",
                Data:    map[string]string{"path": path},
        })
}

// GET /api/pelaksanaan/:id/sertifikat/download — download sertifikat
func (h *PelaksanaanHandler) DownloadSertifikat(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        var sertPath string
        database.DB.QueryRow(c.Request.Context(),
                "SELECT sertifikat_path FROM pelaksanaan_magang WHERE id=$1 AND sertifikat_generated=true", id).
                Scan(&sertPath)

        if sertPath == "" {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Sertifikat belum tersedia"})
                return
        }

        c.FileAttachment(sertPath, "Sertifikat_Magang_TELPP.pdf")
}
