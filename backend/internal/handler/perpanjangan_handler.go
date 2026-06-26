package handler

import (
        "context"
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
        "github.com/telpp/emagang/internal/service"
)

var validDurasi = map[int]bool{7: true, 14: true, 30: true, 60: true, 90: true}

type PerpanjanganHandler struct {
        repo            *repository.PerpanjanganRepository
        pelaksanaanRepo *repository.PelaksanaanRepository
        notifSvc        *service.NotifikasiService
        emailSvc        *service.EmailService
}

func NewPerpanjanganHandler(
        repo *repository.PerpanjanganRepository,
        pelaksanaanRepo *repository.PelaksanaanRepository,
        notifSvc *service.NotifikasiService,
        emailSvc *service.EmailService,
) *PerpanjanganHandler {
        return &PerpanjanganHandler{repo: repo, pelaksanaanRepo: pelaksanaanRepo, notifSvc: notifSvc, emailSvc: emailSvc}
}

// simpanSurat — helper upload surat perpanjangan, return relative path
func simpanSurat(c *gin.Context, formKey string) (*string, error) {
        file, header, ferr := c.Request.FormFile(formKey)
        if ferr != nil {
                return nil, nil // opsional
        }
        defer file.Close()

        if header.Size > config.App.MaxUploadSize {
                return nil, fmt.Errorf("ukuran file maksimal %dMB", config.App.MaxUploadSize/1024/1024)
        }
        ext := strings.ToLower(filepath.Ext(header.Filename))
        if !map[string]bool{".pdf": true, ".jpg": true, ".jpeg": true, ".png": true}[ext] {
                return nil, fmt.Errorf("tipe file tidak didukung, gunakan PDF, JPG, atau PNG")
        }

        subDir := "surat_perpanjangan"
        if err := os.MkdirAll(filepath.Join(config.App.UploadDir, subDir), 0755); err != nil {
                return nil, fmt.Errorf("gagal membuat direktori upload")
        }

        filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
        destPath := filepath.Join(config.App.UploadDir, subDir, filename)
        if err := c.SaveUploadedFile(header, destPath); err != nil {
                return nil, fmt.Errorf("gagal menyimpan file")
        }

        rel := subDir + "/" + filename
        return &rel, nil
}

// POST /api/perpanjangan — peserta ajukan perpanjangan (JSON)
func (h *PerpanjanganHandler) AjukanPeserta(c *gin.Context) {
        userID := middleware.GetUserID(c)

        var req models.AjukanPerpanjanganRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }
        if !validDurasi[req.DurasiHari] {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Durasi harus salah satu dari: 7, 14, 30, 60, atau 90 hari"})
                return
        }

        pelaksanaan, err := h.pelaksanaanRepo.FindByUserID(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "not_found", Message: "Data magang tidak ditemukan"})
                return
        }
        if pelaksanaan.Status != models.StatusAktif {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_status", Message: "Perpanjangan hanya bisa diajukan saat magang berstatus aktif"})
                return
        }
        // Cek window 5 hari: hanya bisa ajukan dalam 5 hari terakhir sebelum magang selesai
        today := time.Now().Truncate(24 * time.Hour)
        selesai := pelaksanaan.TanggalSelesai.Truncate(24 * time.Hour)
        sisaHari := int(selesai.Sub(today).Hours() / 24)
        if sisaHari > 5 {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "too_early", Message: fmt.Sprintf("Perpanjangan baru bisa diajukan dalam 5 hari terakhir sebelum magang berakhir (sisa %d hari)", sisaHari)})
                return
        }

        p, err := h.repo.CreatePeserta(c.Request.Context(), pelaksanaan.ID, userID, req.DurasiHari, req.Alasan)
        if err != nil {
                if err.Error() == "sudah_diperpanjang" {
                        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "already_extended", Message: "Magang ini sudah pernah diperpanjang, tidak bisa mengajukan lagi"})
                        return
                }
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "failed", Message: err.Error()})
                return
        }

        // Notif ke semua HRD aktif
        go func() {
                rows, qerr := repository.GetDB().Query(c.Request.Context(),
                        `SELECT id FROM users WHERE role = 'hrd' AND is_active = true`)
                if qerr != nil {
                        return
                }
                defer rows.Close()
                for rows.Next() {
                        var hrdID uuid.UUID
                        if rows.Scan(&hrdID) == nil {
                                h.notifSvc.KirimKeUser(c.Request.Context(), hrdID, models.RoleHRD,
                                        "Pengajuan Perpanjangan Magang",
                                        fmt.Sprintf("Peserta mengajukan perpanjangan %d hari", req.DurasiHari),
                                        "perpanjangan_baru", &p.ID)
                        }
                }
        }()

        c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Pengajuan perpanjangan berhasil dikirim, menunggu persetujuan HRD", Data: p})
}

// GET /api/perpanjangan/saya — peserta lihat status pengajuan perpanjangan miliknya
func (h *PerpanjanganHandler) GetSaya(c *gin.Context) {
        userID := middleware.GetUserID(c)
        pelaksanaan, err := h.pelaksanaanRepo.FindByUserID(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusOK, models.SuccessResponse{Data: nil})
                return
        }
        p, err := h.repo.FindByPelaksanaan(c.Request.Context(), pelaksanaan.ID)
        if err != nil {
                c.JSON(http.StatusOK, models.SuccessResponse{Data: nil})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: p})
}

// GET /api/perpanjangan — HRD lihat semua pengajuan
func (h *PerpanjanganHandler) GetAll(c *gin.Context) {
        list, err := h.repo.FindAll(c.Request.Context())
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        if list == nil {
                list = []models.PerpanjanganMagangDetail{}
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// POST /api/perpanjangan/langsung — HRD perpanjang langsung (multipart/form-data)
func (h *PerpanjanganHandler) AjukanHRD(c *gin.Context) {
        hrdID := middleware.GetUserID(c)

        pelaksanaanIDStr := c.PostForm("pelaksanaan_id")
        alasan := c.PostForm("alasan")
        durasiStr := c.PostForm("durasi_hari")

        pelaksanaanID, err := uuid.Parse(pelaksanaanIDStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "pelaksanaan_id tidak valid"})
                return
        }
        if len(strings.TrimSpace(alasan)) < 10 {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Alasan minimal 10 karakter"})
                return
        }

        var durasiHari int
        fmt.Sscanf(durasiStr, "%d", &durasiHari)
        if !validDurasi[durasiHari] {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Durasi harus salah satu dari: 7, 14, 30, 60, atau 90 hari"})
                return
        }

        // Validasi: status pelaksanaan harus aktif
        pelaksanaanHRD, err := h.pelaksanaanRepo.FindByID(c.Request.Context(), pelaksanaanID)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "not_found", Message: "Data magang tidak ditemukan"})
                return
        }
        if pelaksanaanHRD.Status != models.StatusAktif {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_status", Message: "Perpanjangan langsung HRD hanya bisa dilakukan saat magang berstatus aktif"})
                return
        }

        suratPath, err := simpanSurat(c, "surat")
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "file_error", Message: err.Error()})
                return
        }

        p, err := h.repo.CreateHRDDirect(c.Request.Context(), pelaksanaanID, hrdID, durasiHari, alasan, suratPath)
        if err != nil {
                if err.Error() == "sudah_diperpanjang" {
                        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "already_extended", Message: "Magang ini sudah pernah diperpanjang"})
                        return
                }
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "failed", Message: err.Error()})
                return
        }

        // Notif + email ke peserta (pakai context.Background() agar tidak expired setelah c.JSON)
        suratPathCopy := ""
        if suratPath != nil {
                suratPathCopy = *suratPath
        }
        alasanCopy := alasan
        durasiCopy := durasiHari
        pIDCopy := p.ID
        emailSvc := h.emailSvc
        notifSvc := h.notifSvc
        pelaksanaanIDCopy := pelaksanaanID
        go func() {
                ctx := context.Background()
                var pesertaID uuid.UUID
                var email, nama string
                var tanggalSelesai string
                if scanErr := repository.GetDB().QueryRow(ctx, `
                        SELECT u.id, u.email, u.nama_lengkap,
                               to_char(pm.tanggal_selesai, 'DD Month YYYY')
                        FROM pelaksanaan_magang pm
                        JOIN users u ON u.id = pm.user_id
                        WHERE pm.id = $1`, pelaksanaanIDCopy,
                ).Scan(&pesertaID, &email, &nama, &tanggalSelesai); scanErr != nil {
                        fmt.Printf("[perpanjangan] gagal fetch peserta untuk email: %v\n", scanErr)
                        return
                }
                notifSvc.KirimKeUser(ctx, pesertaID, models.RolePeserta,
                        "Magang Anda Diperpanjang",
                        fmt.Sprintf("HRD telah memperpanjang durasi magang Anda sebesar %d hari", durasiCopy),
                        "perpanjangan_disetujui", &pIDCopy)
                if emailSvc != nil {
                        if errEmail := emailSvc.KirimPerpanjangan(email, nama, durasiCopy, tanggalSelesai, alasanCopy, config.App.UploadDir, suratPathCopy); errEmail != nil {
                                fmt.Printf("[perpanjangan] gagal kirim email ke %s: %v\n", email, errEmail)
                        } else {
                                fmt.Printf("[perpanjangan] email perpanjangan terkirim ke %s\n", email)
                        }
                }
        }()

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: fmt.Sprintf("Perpanjangan %d hari berhasil diterapkan", durasiHari),
                Data:    p,
        })
}

// PATCH /api/perpanjangan/:id/approve — HRD setujui pengajuan peserta (multipart: opsional surat)
func (h *PerpanjanganHandler) Approve(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        approvedBy := middleware.GetUserID(c)
        catatan := c.PostForm("catatan_hrd")

        suratPath, err := simpanSurat(c, "surat")
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "file_error", Message: err.Error()})
                return
        }

        pelaksanaanID, err := h.repo.Approve(c.Request.Context(), id, approvedBy, suratPath, catatan)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "failed", Message: err.Error()})
                return
        }

        // Notif + email ke peserta (pakai context.Background() agar tidak expired setelah c.JSON)
        suratApprPath := ""
        if suratPath != nil {
                suratApprPath = *suratPath
        }
        catatanCopy := catatan
        idCopy := id
        emailSvcA := h.emailSvc
        notifSvcA := h.notifSvc
        go func() {
                ctx := context.Background()
                var pesertaID uuid.UUID
                var email, nama string
                var durasi int
                var tanggalSelesai, alasan string
                if scanErr := repository.GetDB().QueryRow(ctx, `
                        SELECT pm.user_id, u.email, u.nama_lengkap, pp.durasi_hari,
                               to_char(pm.tanggal_selesai, 'DD Month YYYY'),
                               COALESCE(pp.alasan, '')
                        FROM perpanjangan_magang pp
                        JOIN pelaksanaan_magang pm ON pm.id = pp.pelaksanaan_id
                        JOIN users u ON u.id = pm.user_id
                        WHERE pp.id = $1`, idCopy,
                ).Scan(&pesertaID, &email, &nama, &durasi, &tanggalSelesai, &alasan); scanErr != nil {
                        fmt.Printf("[approve] gagal fetch data untuk email: %v\n", scanErr)
                        _ = pelaksanaanID
                        return
                }
                notifSvcA.KirimKeUser(ctx, pesertaID, models.RolePeserta,
                        "Perpanjangan Magang Disetujui",
                        fmt.Sprintf("Pengajuan perpanjangan %d hari Anda telah disetujui oleh HRD", durasi),
                        "perpanjangan_disetujui", &idCopy)
                keterangan := alasan
                if catatanCopy != "" {
                        keterangan = catatanCopy
                }
                if emailSvcA != nil {
                        if errEmail := emailSvcA.KirimPerpanjangan(email, nama, durasi, tanggalSelesai, keterangan, config.App.UploadDir, suratApprPath); errEmail != nil {
                                fmt.Printf("[approve] gagal kirim email ke %s: %v\n", email, errEmail)
                        } else {
                                fmt.Printf("[approve] email perpanjangan terkirim ke %s\n", email)
                        }
                }
                _ = pelaksanaanID
        }()

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Perpanjangan disetujui dan tanggal selesai telah diperbarui"})
}

// PATCH /api/perpanjangan/:id/tolak — HRD tolak pengajuan peserta (JSON)
func (h *PerpanjanganHandler) Tolak(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        rejectedBy := middleware.GetUserID(c)

        var req models.ProsesPerpanjanganRequest
        c.ShouldBindJSON(&req) //nolint — catatan opsional

        pelaksanaanID, err := h.repo.Tolak(c.Request.Context(), id, rejectedBy, req.CatatanHRD)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "failed", Message: "Gagal menolak atau pengajuan tidak ditemukan / sudah diproses"})
                return
        }

        // Notif ke peserta
        go func() {
                var pesertaID uuid.UUID
                if scanErr := repository.GetDB().QueryRow(c.Request.Context(),
                        `SELECT user_id FROM pelaksanaan_magang WHERE id = $1`, pelaksanaanID,
                ).Scan(&pesertaID); scanErr == nil {
                        h.notifSvc.KirimKeUser(c.Request.Context(), pesertaID, models.RolePeserta,
                                "Perpanjangan Magang Ditolak",
                                "Pengajuan perpanjangan magang Anda ditolak oleh HRD",
                                "perpanjangan_ditolak", &id)
                }
        }()

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Pengajuan perpanjangan ditolak"})
}
