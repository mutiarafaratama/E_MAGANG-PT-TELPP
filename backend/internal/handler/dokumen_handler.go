package handler

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
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

// ── Validasi tipe file via magic bytes ───────────────────────────────────────
// Membaca 512 byte pertama file untuk menentukan tipe sesungguhnya.
// Tidak percaya pada ekstensi atau Content-Type header dari client.

type allowedMime struct {
	mime   string
	magic  []byte
	offset int
}

var allowedDocumentTypes = []allowedMime{
	{mime: "application/pdf", magic: []byte{0x25, 0x50, 0x44, 0x46}, offset: 0},         // %PDF
	{mime: "image/jpeg", magic: []byte{0xFF, 0xD8, 0xFF}, offset: 0},                     // JPEG
	{mime: "image/png", magic: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, offset: 0}, // PNG
	{mime: "application/msword", magic: []byte{0xD0, 0xCF, 0x11, 0xE0}, offset: 0},       // DOC
	{mime: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		magic: []byte{0x50, 0x4B, 0x03, 0x04}, offset: 0}, // DOCX (ZIP signature)
}

var allowedExtensions = map[string]bool{
	".pdf":  true,
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".doc":  true,
	".docx": true,
}

// validateFileType membaca magic bytes file dan memverifikasi tipe file.
// Mengembalikan detected MIME type atau error jika tidak diizinkan.
func validateFileType(fh *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("ekstensi file tidak diizinkan: %s (diizinkan: PDF, JPG, PNG, DOC, DOCX)", ext)
	}

	f, err := fh.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file")
	}
	defer f.Close()

	header := make([]byte, 512)
	n, err := f.Read(header)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("gagal membaca file")
	}
	header = header[:n]

	for _, t := range allowedDocumentTypes {
		if len(header) >= t.offset+len(t.magic) {
			if bytes.Equal(header[t.offset:t.offset+len(t.magic)], t.magic) {
				return t.mime, nil
			}
		}
	}

	return "", fmt.Errorf("tipe file tidak dikenali atau tidak diizinkan. Upload PDF, JPG, PNG, atau DOC/DOCX")
}

// safeInUploadDir memastikan path hasil resolve ada di dalam uploadDir
// (mencegah path traversal seperti ../../etc/passwd)
func safeInUploadDir(uploadDir, targetPath string) bool {
	absUpload, err1 := filepath.Abs(uploadDir)
	absTarget, err2 := filepath.Abs(targetPath)
	if err1 != nil || err2 != nil {
		return false
	}
	return strings.HasPrefix(absTarget, absUpload+string(filepath.Separator)) ||
		absTarget == absUpload
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

	// Validasi tipe file via magic bytes
	detectedMime, err := validateFileType(header)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_file_type", Message: err.Error()})
		return
	}

	// Buat folder per user
	uploadPath := filepath.Join(config.App.UploadDir, userID.String())
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal membuat folder upload"})
		return
	}

	// Nama file unik (tanpa karakter berbahaya)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	uniqueName := fmt.Sprintf("%s_%s_%d%s", jenis, userID.String()[:8], time.Now().Unix(), ext)
	savePath := filepath.Join(uploadPath, uniqueName)

	// Path traversal guard
	if !safeInUploadDir(config.App.UploadDir, savePath) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_path", Message: "Path file tidak valid"})
		return
	}

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
		MimeType:    detectedMime,
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

	// Validasi tipe file via magic bytes
	detectedMime, err := validateFileType(header)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_file_type", Message: err.Error()})
		return
	}

	// Folder publik dipisah dari folder per-user
	uploadPath := filepath.Join(config.App.UploadDir, "publik")
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal membuat folder upload"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	uniqueName := fmt.Sprintf("%s_%s_%d%s", jenis, pid.String()[:8], time.Now().Unix(), ext)
	savePath := filepath.Join(uploadPath, uniqueName)

	// Path traversal guard
	if !safeInUploadDir(config.App.UploadDir, savePath) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_path", Message: "Path file tidak valid"})
		return
	}

	if err := c.SaveUploadedFile(header, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan file"})
		return
	}

	d := &models.Dokumen{
		PengajuanID: &pid,
		Jenis:       jenis,
		NamaFile:    header.Filename,
		PathFile:    savePath,
		UkuranBytes: header.Size,
		MimeType:    detectedMime,
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

	// Path traversal guard
	if !safeInUploadDir(config.App.UploadDir, doc.PathFile) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "forbidden", Message: "Akses file ditolak"})
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

	os.Remove(doc.PathFile)

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Dokumen berhasil dihapus"})
}

// GET /api/dokumen/:id/download — download file (dengan cek kepemilikan)
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

	// Path traversal guard
	if !safeInUploadDir(config.App.UploadDir, doc.PathFile) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "forbidden", Message: "Akses file ditolak"})
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
