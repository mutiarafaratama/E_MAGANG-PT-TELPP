package handler

import (
        "fmt"
        "net/http"
        "os"
        "path/filepath"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/telpp/emagang/internal/config"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
        "github.com/telpp/emagang/internal/service"
)

type DokumenHandler struct {
        repo          *repository.DokumenRepository
        pengajuanRepo *repository.PengajuanRepository
        notifSvc      *service.NotifikasiService
}

func NewDokumenHandler(
        repo *repository.DokumenRepository,
        pengajuanRepo *repository.PengajuanRepository,
        notifSvc *service.NotifikasiService,
) *DokumenHandler {
        return &DokumenHandler{
                repo:          repo,
                pengajuanRepo: pengajuanRepo,
                notifSvc:      notifSvc,
        }
}

// POST /api/dokumen/upload — upload dokumen (peserta atau HRD, harus login)
func (h *DokumenHandler) Upload(c *gin.Context) {
        userID := middleware.GetUserID(c)

        pengajuanIDStr := c.PostForm("pengajuan_id")
        jenis := c.PostForm("jenis")

        if jenis == "" {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "jenis dokumen wajib diisi"})
                return
        }

        file, header, err := c.Request.FormFile("file")
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "upload_error", Message: "File tidak ditemukan"})
                return
        }
        defer file.Close()

        if header.Size > config.App.MaxUploadSize {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{
                        Error:   "file_too_large",
                        Message: fmt.Sprintf("Ukuran file maksimal %dMB", config.App.MaxUploadSize/1024/1024),
                })
                return
        }

        // Buat folder per user
        uploadPath := filepath.Join(config.App.UploadDir, userID.String())
        if err := os.MkdirAll(uploadPath, 0755); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal membuat folder upload"})
                return
        }

        // Nama file unik
        ext := filepath.Ext(header.Filename)
        uniqueName := fmt.Sprintf("%s_%s_%d%s", jenis, userID.String()[:8], time.Now().Unix(), ext)
        savePath := filepath.Join(uploadPath, uniqueName)

        if err := c.SaveUploadedFile(header, savePath); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan file"})
                return
        }

        d := &models.Dokumen{
                UserID:      &userID,
                Jenis:       jenis,
                NamaFile:    header.Filename,
                PathFile:    savePath,
                UkuranBytes: header.Size,
        }

        ct := header.Header.Get("Content-Type")
        if ct != "" {
                d.MimeType = ct
        }

        if pengajuanIDStr != "" {
                pid, err := uuid.Parse(pengajuanIDStr)
                if err == nil {
                        d.PengajuanID = &pid
                }
        }

        if err := h.repo.Save(c.Request.Context(), d); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan data dokumen"})
                return
        }

        // Kirim notifikasi ke peserta saat HRD upload surat balasan
        if jenis == string(models.DokumenSuratBalas) && d.PengajuanID != nil {
                pengajuan, err := h.pengajuanRepo.FindByID(c.Request.Context(), *d.PengajuanID)
                if err == nil && pengajuan != nil && pengajuan.UserID != nil && *pengajuan.UserID != userID {
                        h.notifSvc.KirimKeUser(
                                c.Request.Context(),
                                *pengajuan.UserID,
                                models.RolePeserta,
                                "📄 Surat Balasan Magang Tersedia",
                                "Tim HRD telah mengupload surat balasan untuk pengajuan magang Anda. Silakan cek di halaman Pengajuan.",
                                string(models.NotifDokumen),
                                d.PengajuanID,
                        )
                }
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: "Dokumen berhasil diupload",
                Data:    d,
        })
}

// POST /api/dokumen/upload-publik — upload dari form publik (tanpa login)
func (h *DokumenHandler) UploadPublik(c *gin.Context) {
        pengajuanIDStr := c.PostForm("pengajuan_id")
        jenis := c.PostForm("jenis")

        if jenis == "" {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "jenis dokumen wajib diisi"})
                return
        }
        if pengajuanIDStr == "" {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "pengajuan_id wajib diisi"})
                return
        }

        pid, err := uuid.Parse(pengajuanIDStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "pengajuan_id tidak valid"})
                return
        }

        file, header, err := c.Request.FormFile("file")
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "upload_error", Message: "File tidak ditemukan"})
                return
        }
        defer file.Close()

        if header.Size > config.App.MaxUploadSize {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{
                        Error:   "file_too_large",
                        Message: fmt.Sprintf("Ukuran file maksimal %dMB", config.App.MaxUploadSize/1024/1024),
                })
                return
        }

        // Folder publik dipisah dari folder per-user
        uploadPath := filepath.Join(config.App.UploadDir, "publik")
        if err := os.MkdirAll(uploadPath, 0755); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal membuat folder upload"})
                return
        }

        ext := filepath.Ext(header.Filename)
        uniqueName := fmt.Sprintf("%s_%s_%d%s", jenis, pid.String()[:8], time.Now().Unix(), ext)
        savePath := filepath.Join(uploadPath, uniqueName)

        if err := c.SaveUploadedFile(header, savePath); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan file"})
                return
        }

        // user_id sengaja tidak diisi — SavePublik tidak menyertakan kolom user_id
        d := &models.Dokumen{
                PengajuanID: &pid,
                Jenis:       jenis,
                NamaFile:    header.Filename,
                PathFile:    savePath,
                UkuranBytes: header.Size,
        }

        ct := header.Header.Get("Content-Type")
        if ct != "" {
                d.MimeType = ct
        }

        if err := h.repo.SavePublik(c.Request.Context(), d); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan data dokumen: " + err.Error()})
                return
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: "Dokumen berhasil diupload",
                Data:    d,
        })
}

// GET /api/dokumen/pengajuan/:id — daftar dokumen per pengajuan (HRD/Admin)
func (h *DokumenHandler) GetByPengajuan(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        // Peserta hanya boleh lihat dokumen pengajuan milik sendiri
        role := middleware.GetUserRole(c)
        if role == models.RolePeserta {
                userID := middleware.GetUserID(c)
                pengajuan, err := h.pengajuanRepo.FindByID(c.Request.Context(), id)
                if err != nil || pengajuan == nil {
                        c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Pengajuan tidak ditemukan"})
                        return
                }
                if pengajuan.UserID == nil || *pengajuan.UserID != userID {
                        c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "forbidden", Message: "Akses ditolak"})
                        return
                }
        }

        docs, err := h.repo.FindByPengajuanID(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Data: docs})
}

// GET /api/admin/dokumen — list semua dokumen (admin only)
func (h *DokumenHandler) AdminList(c *gin.Context) {
        jenis  := c.Query("jenis")
        search := c.Query("search")
        page   := queryInt(c, "page", 1)
        limit  := queryInt(c, "limit", 20)

        list, total, err := h.repo.FindAll(c.Request.Context(), jenis, search, page, limit)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        if list == nil {
                list = []models.DokumenWithUser{}
        }
        c.JSON(http.StatusOK, models.PaginatedResponse{
                Data: list, Total: total, Page: page, Limit: limit,
        })
}

// GET /api/admin/dokumen/:id/download — download dokumen oleh admin
func (h *DokumenHandler) AdminDownload(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        doc, err := h.repo.FindByID(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Dokumen tidak ditemukan"})
                return
        }

        c.FileAttachment(doc.PathFile, doc.NamaFile)
}

// DELETE /api/admin/dokumen/:id — hapus dokumen dari DB + filesystem (admin only)
func (h *DokumenHandler) AdminHapus(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        doc, err := h.repo.FindByID(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Dokumen tidak ditemukan"})
                return
        }

        if err := h.repo.DeleteByID(c.Request.Context(), id); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menghapus dokumen"})
                return
        }

        // Hapus file fisik (best-effort, tidak gagalkan response jika file sudah tidak ada)
        os.Remove(doc.PathFile)

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Dokumen berhasil dihapus"})
}

// GET /api/dokumen/:id/download — download file
func (h *DokumenHandler) Download(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        doc, err := h.repo.FindByID(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Dokumen tidak ditemukan"})
                return
        }

        // Peserta hanya boleh download dokumen milik sendiri atau dari pengajuannya
        userID := middleware.GetUserID(c)
        role := middleware.GetUserRole(c)
        if role == models.RolePeserta {
                allowed := false
                if doc.UserID != nil && *doc.UserID == userID {
                        allowed = true
                } else if doc.PengajuanID != nil {
                        // Dokumen diupload via form publik (user_id null) — cek via pengajuan
                        pengajuan, err2 := h.pengajuanRepo.FindByID(c.Request.Context(), *doc.PengajuanID)
                        if err2 == nil && pengajuan != nil && pengajuan.UserID != nil && *pengajuan.UserID == userID {
                                allowed = true
                        }
                }
                if !allowed {
                        c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "forbidden", Message: "Akses ditolak"})
                        return
                }
        }

        c.FileAttachment(doc.PathFile, doc.NamaFile)
}
