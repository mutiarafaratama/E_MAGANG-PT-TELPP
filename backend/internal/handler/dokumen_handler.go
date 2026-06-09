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
)

type DokumenHandler struct {
	repo *repository.DokumenRepository
}

func NewDokumenHandler(repo *repository.DokumenRepository) *DokumenHandler {
	return &DokumenHandler{repo: repo}
}

// POST /api/dokumen/upload — upload dokumen
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

	// Simpan file
	if err := c.SaveUploadedFile(header, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan file"})
		return
	}

	// Simpan ke DB
	d := &models.Dokumen{
		UserID:      userID,
		Jenis:       models.JenisDokumen(jenis),
		NamaFile:    header.Filename,
		PathFile:    savePath,
		UkuranBytes: &header.Size,
	}

	ct := header.Header.Get("Content-Type")
	if ct != "" {
		d.MimeType = &ct
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

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Dokumen berhasil diupload",
		Data:    d,
	})
}

// GET /api/dokumen/pengajuan/:id — daftar dokumen per pengajuan
func (h *DokumenHandler) GetByPengajuan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
		return
	}

	docs, err := h.repo.FindByPengajuanID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Data: docs})
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

	// Cek akses: peserta hanya bisa download milik sendiri
	userID := middleware.GetUserID(c)
	role := middleware.GetUserRole(c)
	if role == models.RolePeserta && doc.UserID != userID {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "forbidden", Message: "Akses ditolak"})
		return
	}

	c.FileAttachment(doc.PathFile, doc.NamaFile)
}
