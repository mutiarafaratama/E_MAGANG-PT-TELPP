package handler

import (
        "io"
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
        repo          *repository.PelaksanaanRepository
        pengajuanRepo *repository.PengajuanRepository
        sertifikatSvc *service.SertifikatService
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

        pengajuan, err := h.pengajuanRepo.FindByID(c.Request.Context(), pengajuanID)
        if err != nil {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Pengajuan tidak ditemukan"})
                return
        }

        mulai, err := time.Parse("2006-01-02", req.TanggalMulai)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_date", Message: "Format tanggal_mulai tidak valid (YYYY-MM-DD)"})
                return
        }
        selesai, err := time.Parse("2006-01-02", req.TanggalSelesai)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_date", Message: "Format tanggal_selesai tidak valid (YYYY-MM-DD)"})
                return
        }

        var userIDVal uuid.UUID
        if pengajuan.UserID != nil {
                userIDVal = *pengajuan.UserID
        }
        p := &models.PelaksanaanMagang{
                PengajuanID:    pengajuanID,
                UserID:         userIDVal,
                TanggalMulai:   mulai,
                TanggalSelesai: selesai,
                Divisi:         req.Divisi,
                PembimbingNama: &req.Pembimbing,
        }
        if err := h.repo.Create(c.Request.Context(), p); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Jadwal magang berhasil ditetapkan", Data: p})
}

// GET /api/pelaksanaan/saya — peserta ambil data pelaksanaan sendiri
func (h *PelaksanaanHandler) GetMy(c *gin.Context) {
        userID := middleware.GetUserID(c)
        p, err := h.repo.FindByUserID(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Data pelaksanaan tidak ditemukan"})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: p})
}

// GET /api/pelaksanaan — HRD ambil semua pelaksanaan (dengan paginasi & filter status)
func (h *PelaksanaanHandler) GetAll(c *gin.Context) {
        status := c.Query("status")
        page := 1
        limit := 20
        list, total, err := h.repo.FindAll(c.Request.Context(), status, page, limit)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{
                Data: map[string]any{"list": list, "total": total, "page": page, "limit": limit},
        })
}

// PATCH /api/pelaksanaan/:id/status — HRD update status pelaksanaan
func (h *PelaksanaanHandler) UpdateStatus(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        var req struct {
                Status string `json:"status" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }
        if err := h.repo.UpdateStatus(c.Request.Context(), id, models.StatusPelaksanaan(req.Status)); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Status berhasil diperbarui"})
}

// PATCH /api/pelaksanaan/:id/nilai — HRD set nilai akhir
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
        dinilaiOleh := middleware.GetUserID(c)
        catatanNilai := ""
        if req.CatatanNilai != nil {
                catatanNilai = *req.CatatanNilai
        }
        if err := h.repo.SetNilai(c.Request.Context(), id, req.Nilai, catatanNilai, dinilaiOleh); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Nilai berhasil disimpan"})
}

// PATCH /api/pelaksanaan/:id/pembimbing — HRD set pembimbing lapangan
func (h *PelaksanaanHandler) UpdatePembimbing(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        var req struct {
                Pembimbing string `json:"pembimbing" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }
        if err := h.repo.UpdatePembimbing(c.Request.Context(), id, req.Pembimbing); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Pembimbing berhasil diperbarui"})
}

// POST /api/pelaksanaan/:id/sertifikat — HRD upload file PDF sertifikat untuk peserta
func (h *PelaksanaanHandler) UploadSertifikat(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        file, _, err := c.Request.FormFile("file")
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "file_required", Message: "File PDF wajib dilampirkan"})
                return
        }
        defer file.Close()

        data, err := io.ReadAll(file)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "read_error", Message: "Gagal membaca file"})
                return
        }

        if err := h.sertifikatSvc.Upload(c.Request.Context(), id, data); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "upload_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Sertifikat berhasil diunggah dan dikirim ke peserta"})
}

// GET /api/sertifikat — HRD list peserta yang siap menerima sertifikat
func (h *PelaksanaanHandler) ListSertifikatHRD(c *gin.Context) {
        rows, err := database.DB.Query(c.Request.Context(), `
                SELECT pm.id, pm.status, pm.tanggal_mulai, pm.tanggal_selesai,
                       COALESCE(pm.divisi, '') AS divisi,
                       COALESCE(pm.pembimbing, u2.nama_lengkap, '') AS pembimbing,
                       pm.sertifikat_generated, pm.sertifikat_generated_at,
                       pm.nilai,
                       u.nama_lengkap, u.email
                FROM pelaksanaan_magang pm
                JOIN users u ON u.id = pm.user_id
                LEFT JOIN users u2 ON u2.id = pm.pembimbing_id
                WHERE pm.status IN ('penilaian', 'selesai') AND pm.nilai IS NOT NULL
                ORDER BY pm.sertifikat_generated ASC, pm.updated_at DESC`)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        defer rows.Close()

        var list []models.SertifikatListItem
        for rows.Next() {
                var item models.SertifikatListItem
                var generatedAt *time.Time
                if err := rows.Scan(
                        &item.PelaksanaanID, &item.Status,
                        &item.TanggalMulai, &item.TanggalSelesai,
                        &item.Divisi, &item.Pembimbing,
                        &item.SertifikatGenerated, &generatedAt,
                        &item.Nilai,
                        &item.NamaLengkap, &item.Email,
                ); err != nil {
                        continue
                }
                item.SertifikatGeneratedAt = generatedAt
                list = append(list, item)
        }
        if list == nil {
                list = []models.SertifikatListItem{}
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// GET /api/pelaksanaan/:id/sertifikat/download — peserta download sertifikat
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
