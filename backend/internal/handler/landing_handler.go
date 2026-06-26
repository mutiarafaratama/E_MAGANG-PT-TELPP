package handler

import (
        "fmt"
        "net/http"
        "os"
        "path/filepath"
        "strings"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/telpp/emagang/internal/config"
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
                result[item.Kunci] = item.Nilai
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

// GET /api/landing/periodes — publik, semua periode untuk tampilan jadwal
func (h *LandingHandler) GetAllPeriodePublic(c *gin.Context) {
        list, err := h.repo.GetAllPeriode(c.Request.Context())
        if err != nil {
                c.JSON(http.StatusOK, models.SuccessResponse{Data: []interface{}{}})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
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

        kuota := 0
        if req.Kuota != nil {
                kuota = *req.Kuota
        }
        p := &models.PeriodeMagang{
                Nama:         req.Nama,
                TanggalBuka:  buka,
                TanggalTutup: tutup,
                Kuota:        kuota,
        }

        if err := h.repo.CreatePeriode(c.Request.Context(), p); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Periode berhasil dibuat", Data: p})
}

// PUT /api/admin/periode/:id — Admin edit periode
func (h *LandingHandler) UpdatePeriode(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        var req struct {
                Nama         string `json:"nama" binding:"required"`
                TanggalBuka  string `json:"tanggal_buka" binding:"required"`
                TanggalTutup string `json:"tanggal_tutup" binding:"required"`
                Kuota        *int   `json:"kuota"`
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
        if err := h.repo.UpdatePeriode(c.Request.Context(), id, req.Nama, buka, tutup, req.Kuota); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Periode berhasil diperbarui"})
}

// DELETE /api/admin/periode/:id — Admin hapus periode
func (h *LandingHandler) DeletePeriode(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        if err := h.repo.DeletePeriode(c.Request.Context(), id); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Periode berhasil dihapus"})
}

// PATCH /api/admin/periode/:id/aktif — Admin aktifkan/nonaktifkan periode
func (h *LandingHandler) SetPeriodeAktif(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        var req struct {
                IsActive bool `json:"is_active"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        if err := h.repo.SetPeriodeActive(c.Request.Context(), id, req.IsActive); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal mengubah status periode: " + err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Status periode berhasil diperbarui"})
}

// GET /api/landing/alur — publik, ambil semua item alur pendaftaran
func (h *LandingHandler) GetAlur(c *gin.Context) {
        list, err := h.repo.GetAllAlur(c.Request.Context())
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        if list == nil {
                list = []models.AlurItem{}
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// POST /api/admin/alur — Admin tambah item alur
func (h *LandingHandler) CreateAlur(c *gin.Context) {
        var req models.AlurItem
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }
        if err := h.repo.CreateAlur(c.Request.Context(), &req); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Item alur berhasil ditambahkan", Data: req})
}

// PUT /api/admin/alur/:id — Admin update item alur
func (h *LandingHandler) UpdateAlur(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        var req models.AlurItem
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }
        req.ID = id
        if err := h.repo.UpdateAlur(c.Request.Context(), &req); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Item alur berhasil diperbarui", Data: req})
}

// DELETE /api/admin/alur/:id — Admin hapus item alur
func (h *LandingHandler) DeleteAlur(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        h.repo.DeleteAlur(c.Request.Context(), id)
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Item alur berhasil dihapus"})
}

// POST /api/admin/hero/upload — Admin upload background hero
func (h *LandingHandler) UploadHeroBg(c *gin.Context) {
        header, err := c.FormFile("gambar")
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "no_file", Message: "File gambar tidak ditemukan (field: gambar)"})
                return
        }

        ext := strings.ToLower(filepath.Ext(header.Filename))
        allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
        if !allowed[ext] {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_type", Message: "Format tidak didukung. Gunakan jpg/png/webp"})
                return
        }

        if header.Size > 10*1024*1024 {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "file_too_large", Message: "Ukuran gambar maksimal 10 MB"})
                return
        }

        uploadDir := filepath.Join(config.App.UploadDir, "hero")
        if err := os.MkdirAll(uploadDir, 0755); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal membuat folder upload"})
                return
        }

        filename := fmt.Sprintf("%d_hero%s", time.Now().UnixMilli(), ext)
        savePath := filepath.Join(uploadDir, filename)
        if err := c.SaveUploadedFile(header, savePath); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan gambar"})
                return
        }

        url := "/uploads/hero/" + filename

        // Simpan URL ke landing_content
        updatedBy := middleware.GetUserID(c)
        h.repo.UpsertContent(c.Request.Context(), "hero_bg_url", url, "image", updatedBy)

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Gambar hero berhasil diupload", Data: map[string]string{"url": url}})
}

// POST /api/admin/alur/upload — Admin upload gambar untuk item alur
func (h *LandingHandler) UploadGambarAlur(c *gin.Context) {
        header, err := c.FormFile("gambar")
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "no_file", Message: "File gambar tidak ditemukan (field: gambar)"})
                return
        }

        ext := strings.ToLower(filepath.Ext(header.Filename))
        allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
        if !allowed[ext] {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_type", Message: "Format tidak didukung. Gunakan jpg/png/webp/gif"})
                return
        }

        if header.Size > 5*1024*1024 {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "file_too_large", Message: "Ukuran gambar maksimal 5 MB"})
                return
        }

        uploadDir := filepath.Join(config.App.UploadDir, "alur")
        if err := os.MkdirAll(uploadDir, 0755); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal membuat folder upload"})
                return
        }

        filename := fmt.Sprintf("%d_alur%s", time.Now().UnixMilli(), ext)
        savePath := filepath.Join(uploadDir, filename)
        if err := c.SaveUploadedFile(header, savePath); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan gambar"})
                return
        }

        url := "/uploads/alur/" + filename
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Gambar berhasil diupload", Data: map[string]string{"url": url}})
}

// POST /api/admin/jadwal/bg — Update background settings (type, color, opacity)
func (h *LandingHandler) UpdateJadwalBg(c *gin.Context) {
        var req struct {
                Type    string `json:"type"`
                Color   string `json:"color"`
                Opacity string `json:"opacity"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad_request", Message: "Format JSON tidak valid"})
                return
        }
        updatedBy := middleware.GetUserID(c)
        ctx := c.Request.Context()
        validTypes := map[string]bool{"none": true, "color": true, "image": true}
        if req.Type != "" && validTypes[req.Type] {
                h.repo.UpsertContent(ctx, "jadwal_bg_type", req.Type, "text", updatedBy)
        }
        if req.Color != "" {
                h.repo.UpsertContent(ctx, "jadwal_bg_color", req.Color, "text", updatedBy)
        }
        if req.Opacity != "" {
                h.repo.UpsertContent(ctx, "jadwal_bg_opacity", req.Opacity, "text", updatedBy)
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Background jadwal berhasil diupdate"})
}

// POST /api/admin/jadwal/upload-bg — Upload background image for jadwal section
func (h *LandingHandler) UploadJadwalBg(c *gin.Context) {
        header, err := c.FormFile("gambar")
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "no_file", Message: "File gambar tidak ditemukan (field: gambar)"})
                return
        }
        ext := strings.ToLower(filepath.Ext(header.Filename))
        allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
        if !allowed[ext] {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_type", Message: "Format tidak didukung. Gunakan jpg/png/webp"})
                return
        }
        if header.Size > 10*1024*1024 {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "file_too_large", Message: "Ukuran gambar maksimal 10 MB"})
                return
        }
        uploadDir := filepath.Join(config.App.UploadDir, "jadwal")
        if err := os.MkdirAll(uploadDir, 0755); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal membuat folder upload"})
                return
        }
        filename := fmt.Sprintf("%d_jadwal%s", time.Now().UnixMilli(), ext)
        savePath := filepath.Join(uploadDir, filename)
        if err := c.SaveUploadedFile(header, savePath); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan gambar"})
                return
        }
        url := "/uploads/jadwal/" + filename
        updatedBy := middleware.GetUserID(c)
        h.repo.UpsertContent(c.Request.Context(), "jadwal_bg_url", url, "image", updatedBy)
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Gambar jadwal berhasil diupload", Data: map[string]string{"url": url}})
}
