package handler

import (
        "fmt"
        "io"
        "net/http"
        "os"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
)

type LaporanHandler struct {
        repo        *repository.LaporanRepository
        pelRepo     *repository.PelaksanaanRepository
}

func NewLaporanHandler(repo *repository.LaporanRepository, pelRepo *repository.PelaksanaanRepository) *LaporanHandler {
        return &LaporanHandler{repo: repo, pelRepo: pelRepo}
}

// POST /api/laporan/upload — peserta upload laporan
func (h *LaporanHandler) Upload(c *gin.Context) {
        userID := middleware.GetUserID(c)

        // Ambil pelaksanaan milik peserta
        pel, err := h.pelRepo.FindByUserID(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "not_found", Message: "Data pelaksanaan tidak ditemukan"})
                return
        }

        // Boleh upload sejak jadwal magang aktif hingga fase upload_laporan
        if pel.Status != models.StatusAktif && pel.Status != models.StatusUploadLaporan {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{
                        Error:   "not_allowed",
                        Message: fmt.Sprintf("Upload laporan hanya bisa saat magang sedang aktif atau fase upload laporan, status saat ini: %s", pel.Status),
                })
                return
        }

        fh, err := c.FormFile("file")
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "no_file", Message: "File tidak ditemukan dalam request"})
                return
        }

        // Validasi ukuran (20 MB)
        if fh.Size > 20*1024*1024 {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "file_too_large", Message: "Ukuran file maksimal 20 MB"})
                return
        }

        laporan, err := h.repo.Upload(c.Request.Context(), pel.ID, fh)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "upload_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: "Laporan berhasil diunggah",
                Data:    laporan,
        })
}

// GET /api/laporan/saya — peserta lihat daftar laporan miliknya
func (h *LaporanHandler) GetMy(c *gin.Context) {
        userID := middleware.GetUserID(c)

        pel, err := h.pelRepo.FindByUserID(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusOK, models.SuccessResponse{Data: []interface{}{}})
                return
        }

        list, err := h.repo.FindByPelaksanaan(c.Request.Context(), pel.ID)
        if err != nil || list == nil {
                c.JSON(http.StatusOK, models.SuccessResponse{Data: []interface{}{}})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// GET /api/laporan/pelaksanaan/:pelaksanaan_id — HRD lihat laporan peserta tertentu
func (h *LaporanHandler) GetByPelaksanaan(c *gin.Context) {
        pelID, err := uuid.Parse(c.Param("pelaksanaan_id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        list, err := h.repo.FindByPelaksanaan(c.Request.Context(), pelID)
        if err != nil || list == nil {
                c.JSON(http.StatusOK, models.SuccessResponse{Data: []interface{}{}})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// GET /api/laporan — HRD lihat semua laporan, filter ?status=menunggu_review|disetujui|revisi
func (h *LaporanHandler) GetAll(c *gin.Context) {
        page   := queryInt(c, "page", 1)
        limit  := queryInt(c, "limit", 20)
        status := c.Query("status") // kosong = semua

        list, total, err := h.repo.FindAll(c.Request.Context(), status, page, limit)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        totalPages := (total + limit - 1) / limit
        if totalPages == 0 {
                totalPages = 1
        }

        c.JSON(http.StatusOK, models.PaginatedResponse{
                Data: list, Total: total, Page: page, Limit: limit, TotalPages: totalPages,
        })
}

// PATCH /api/laporan/:id/review — HRD acc atau minta revisi
func (h *LaporanHandler) Review(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        var req models.ReviewLaporanRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        if req.Status != models.StatusLaporanDisetujui && req.Status != models.StatusLaporanRevisi {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_status", Message: "Status harus 'disetujui' atau 'revisi'"})
                return
        }

        reviewerID := middleware.GetUserID(c)

        // Ambil laporan untuk tahu pelaksanaan_id
        laporan, err := h.repo.FindByID(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Laporan tidak ditemukan"})
                return
        }

        if err := h.repo.Review(c.Request.Context(), id, req.Status, req.CatatanHRD, reviewerID); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        // Jika disetujui → update status pelaksanaan ke penilaian
        if req.Status == models.StatusLaporanDisetujui {
                _ = h.pelRepo.UpdateStatus(c.Request.Context(), laporan.PelaksanaanID, models.StatusPenilaian)
        }

        msg := "Laporan berhasil disetujui"
        if req.Status == models.StatusLaporanRevisi {
                msg = "Permintaan revisi berhasil dikirim ke peserta"
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Message: msg})
}

// GET /api/laporan/:id/download — download file laporan
func (h *LaporanHandler) Download(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        laporan, err := h.repo.FindByID(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Laporan tidak ditemukan"})
                return
        }

        f, err := os.Open(laporan.PathFile)
        if err != nil {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "file_not_found", Message: "File tidak ditemukan di server"})
                return
        }
        defer f.Close()

        contentType := "application/octet-stream"
        if laporan.MimeType != nil && *laporan.MimeType != "" {
                contentType = *laporan.MimeType
        }

        c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, laporan.NamaFile))
        c.Header("Content-Type", contentType)
        io.Copy(c.Writer, f)
}
