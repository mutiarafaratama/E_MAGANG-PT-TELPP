package handler

import (
	"context"
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

type IzinSakitHandler struct {
	repo            *repository.IzinSakitRepository
	pelaksanaanRepo *repository.PelaksanaanRepository
	notifSvc        *service.NotifikasiService
}

func NewIzinSakitHandler(
	repo *repository.IzinSakitRepository,
	pelaksanaanRepo *repository.PelaksanaanRepository,
	notifSvc *service.NotifikasiService,
) *IzinSakitHandler {
	return &IzinSakitHandler{repo: repo, pelaksanaanRepo: pelaksanaanRepo, notifSvc: notifSvc}
}

// POST /api/izin-sakit — peserta ajukan izin atau sakit (multipart/form-data)
func (h *IzinSakitHandler) Ajukan(c *gin.Context) {
	jenis := c.PostForm("jenis")
	alasan := c.PostForm("alasan")
	tanggalStr := c.PostForm("tanggal")

	if jenis != "izin" && jenis != "sakit" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Jenis harus 'izin' atau 'sakit'"})
		return
	}
	if len(alasan) < 5 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Alasan terlalu singkat (min 5 karakter)"})
		return
	}

	if _, err := time.Parse("2006-01-02", tanggalStr); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Format tanggal tidak valid (YYYY-MM-DD)"})
		return
	}

	userID := middleware.GetUserID(c)
	pelaksanaan, err := h.pelaksanaanRepo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "not_found", Message: "Data magang tidak ditemukan"})
		return
	}
	if pelaksanaan.Status != models.StatusAktif {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_status", Message: "Magang belum aktif"})
		return
	}

	item := &models.IzinSakitRequest{
		PelaksanaanID: pelaksanaan.ID,
		UserID:        userID,
		Tanggal:       tanggalStr,
		Jenis:         jenis,
		Alasan:        alasan,
	}

	// Tentukan folder berdasarkan jenis
	subDir := "surat_sakit"
	if jenis == "izin" {
		subDir = "surat_izin"
	}

	file, header, ferr := c.Request.FormFile("bukti")
	if ferr != nil {
		msg := "Surat izin wajib dilampirkan"
		if jenis == "sakit" {
			msg = "Surat sakit wajib dilampirkan"
		}
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "file_required",
			Message: msg,
		})
		return
	}
	if ferr == nil {
		defer file.Close()

		if header.Size > config.App.MaxUploadSize {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "file_too_large",
				Message: fmt.Sprintf("Ukuran file maksimal %dMB", config.App.MaxUploadSize/1024/1024),
			})
			return
		}

		ext := filepath.Ext(header.Filename)
		allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".pdf": true}
		if !allowed[ext] {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "invalid_file_type",
				Message: "Tipe file tidak didukung. Gunakan JPG, PNG, atau PDF",
			})
			return
		}

		uploadDir := config.App.UploadDir
		if err := os.MkdirAll(filepath.Join(uploadDir, subDir), 0755); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan file"})
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		destPath := filepath.Join(uploadDir, subDir, filename)
		if err := c.SaveUploadedFile(header, destPath); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan file"})
			return
		}

		relativePath := fmt.Sprintf("%s/%s", subDir, filename)
		item.BuktiPath = &relativePath
	}

	if err := h.repo.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "duplicate", Message: "Sudah ada pengajuan untuk tanggal tersebut"})
		return
	}

	// Push notif ke semua HRD aktif
	if h.notifSvc != nil {
		label := "Izin"
		if jenis == "sakit" {
			label = "Sakit"
		}
		reqID := item.ID
		go func() {
			bgCtx := context.Background()
			rows, qerr := repository.GetDB().Query(bgCtx,
				`SELECT id FROM users WHERE role IN ('hrd','admin') AND is_active = true`)
			if qerr != nil {
				return
			}
			defer rows.Close()
			for rows.Next() {
				var hrdID uuid.UUID
				if rows.Scan(&hrdID) == nil {
					h.notifSvc.KirimKeUser(bgCtx, hrdID, models.RoleHRD,
						fmt.Sprintf("Pengajuan %s Masuk", label),
						fmt.Sprintf("Peserta mengajukan %s untuk tanggal %s", label, tanggalStr),
						string(models.NotifAbsensi), &reqID,
					)
				}
			}
		}()
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Pengajuan berhasil dikirim", Data: item})
}

// GET /api/izin-sakit/saya — peserta lihat riwayat izin/sakit miliknya
func (h *IzinSakitHandler) GetSaya(c *gin.Context) {
	userID := middleware.GetUserID(c)
	list, err := h.repo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// GET /api/izin-sakit — HRD lihat semua request
func (h *IzinSakitHandler) GetAll(c *gin.Context) {
	list, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// PATCH /api/izin-sakit/:id/approve — HRD setujui, otomatis insert absensi
func (h *IzinSakitHandler) Approve(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
		return
	}
	approvedBy := middleware.GetUserID(c)
	if err := h.repo.ApproveAndInsertAbsensi(c.Request.Context(), id, approvedBy); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "failed", Message: "Gagal menyetujui atau request tidak ditemukan / sudah diproses"})
		return
	}

	// Notif ke peserta — cari user_id dari request
	if h.notifSvc != nil {
		go func() {
			bgCtx := context.Background()
			pesertaID, jenis, ferr := h.repo.FindUserAndJenis(bgCtx, id)
			if ferr != nil {
				return
			}
			label := "Izin"
			if jenis == "sakit" {
				label = "Sakit"
			}
			h.notifSvc.KirimKeUser(bgCtx, pesertaID, models.RolePeserta,
				fmt.Sprintf("Pengajuan %s Disetujui", label),
				fmt.Sprintf("HRD telah menyetujui pengajuan %s Anda, absensi dicatat otomatis.", label),
				string(models.NotifAbsensi), &id,
			)
		}()
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Disetujui — absensi telah dicatat otomatis"})
}

// PATCH /api/izin-sakit/:id/tolak — HRD tolak
func (h *IzinSakitHandler) Tolak(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
		return
	}
	var req models.TolakIzinSakitRequest
	c.ShouldBindJSON(&req) //nolint — catatan opsional

	rejectedBy := middleware.GetUserID(c)
	if err := h.repo.Tolak(c.Request.Context(), id, rejectedBy, req.CatatanHRD); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "failed", Message: "Gagal menolak atau request tidak ditemukan / sudah diproses"})
		return
	}

	// Notif ke peserta
	if h.notifSvc != nil {
		go func() {
			bgCtx := context.Background()
			pesertaID, jenis, ferr := h.repo.FindUserAndJenis(bgCtx, id)
			if ferr != nil {
				return
			}
			label := "Izin"
			if jenis == "sakit" {
				label = "Sakit"
			}
			catatan := req.CatatanHRD
			if catatan == "" {
				catatan = "Silakan hubungi HRD untuk informasi lebih lanjut."
			}
			h.notifSvc.KirimKeUser(bgCtx, pesertaID, models.RolePeserta,
				fmt.Sprintf("Pengajuan %s Ditolak", label),
				fmt.Sprintf("Pengajuan %s Anda tidak disetujui. %s", label, catatan),
				string(models.NotifAbsensi), &id,
			)
		}()
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Pengajuan ditolak"})
}
