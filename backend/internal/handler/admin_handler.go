package handler

import (
        "context"
        "fmt"
        "net/http"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/telpp/emagang/internal/config"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
        "github.com/telpp/emagang/internal/service"
        appws "github.com/telpp/emagang/internal/websocket"
        "golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
        userRepo    *repository.UserRepository
        dokumenRepo *repository.DokumenRepository
        emailSvc    *service.EmailService
}

func NewAdminHandler(userRepo *repository.UserRepository, dokumenRepo *repository.DokumenRepository, emailSvc *service.EmailService) *AdminHandler {
        return &AdminHandler{userRepo: userRepo, dokumenRepo: dokumenRepo, emailSvc: emailSvc}
}

// GET /api/admin/stats — statistik dashboard
func (h *AdminHandler) GetStats(c *gin.Context) {
        stats, err := h.userRepo.GetDashboardStats(c.Request.Context())
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: stats})
}

// GET /api/admin/users — daftar semua user
func (h *AdminHandler) GetUsers(c *gin.Context) {
        role := c.Query("role")
        page := queryInt(c, "page", 1)
        limit := queryInt(c, "limit", 20)

        users, total, err := h.userRepo.FindAll(c.Request.Context(), role, page, limit)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.PaginatedResponse{
                Data: users, Total: total, Page: page, Limit: limit,
        })
}

// POST /api/admin/users — buat akun HRD baru
func (h *AdminHandler) CreateUser(c *gin.Context) {
        var req struct {
                NamaLengkap string          `json:"nama_lengkap" binding:"required"`
                Email       string          `json:"email" binding:"required,email"`
                Password    string          `json:"password" binding:"required,min=8"`
                Role        models.UserRole `json:"role" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        if req.Role != models.RoleHRD && req.Role != models.RoleAdmin {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_role", Message: "Role hanya boleh hrd atau admin"})
                return
        }

        hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
        user := &models.User{
                NamaLengkap:  req.NamaLengkap,
                Email:        req.Email,
                PasswordHash: string(hash),
                Role:         req.Role,
        }

        if err := h.userRepo.Create(c.Request.Context(), user); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "create_failed", Message: "Email sudah terdaftar"})
                return
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: "Akun berhasil dibuat",
                Data: models.UserPublic{
                        ID:          user.ID,
                        NamaLengkap: user.NamaLengkap,
                        Email:       user.Email,
                        Role:        user.Role,
                        IsActive:    user.IsActive,
                        CreatedAt:   user.CreatedAt,
                },
        })
}

// PATCH /api/admin/users/:id — update nama dan role user
func (h *AdminHandler) UpdateUser(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        var req struct {
                NamaLengkap string          `json:"nama_lengkap" binding:"required"`
                Role        models.UserRole `json:"role" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        if req.Role != models.RoleHRD && req.Role != models.RoleAdmin {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_role", Message: "Role hanya boleh hrd atau admin"})
                return
        }

        if err := h.userRepo.UpdateUser(c.Request.Context(), id, req.NamaLengkap, req.Role); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "User berhasil diperbarui"})
}

// PATCH /api/admin/users/:id/toggle — aktifkan/nonaktifkan user
func (h *AdminHandler) ToggleUser(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        var req struct {
                IsActive bool `json:"is_active"`
        }
        c.ShouldBindJSON(&req)

        if err := h.userRepo.UpdateActive(c.Request.Context(), id, req.IsActive); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        status := "dinonaktifkan"
        if req.IsActive {
                status = "diaktifkan"
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Akun berhasil " + status})
}

// GET /api/admin/hrd-list — daftar HRD untuk assign
func (h *AdminHandler) GetHRDList(c *gin.Context) {
        list, err := h.userRepo.FindHRDList(c.Request.Context())
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// DELETE /api/admin/users/:id — hapus akun peserta (hanya status selesai)
// Alur: kirim backup email (sertifikat) → hapus semua data peserta dalam transaksi
func (h *AdminHandler) HapusPeserta(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        // Ambil info peserta + validasi
        info, err := h.userRepo.GetPesertaHapusInfo(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Peserta tidak ditemukan"})
                return
        }

        // Validasi status harus selesai
        if info.StatusMagang == nil || *info.StatusMagang != "selesai" {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{
                        Error:   "status_invalid",
                        Message: "Akun hanya bisa dihapus jika magang sudah berstatus Selesai",
                })
                return
        }

        // Ambil dokumen PDF dari tabel dokumen (laporan/sertifikat upload manual)
        rawDocs, _ := h.userRepo.GetDokumenPeserta(c.Request.Context(), id)
        var dokumenLampiran []service.DokumenLampiran
        for _, d := range rawDocs {
                dokumenLampiran = append(dokumenLampiran, service.DokumenLampiran{
                        Jenis:    d.Jenis,
                        NamaFile: d.NamaFile,
                        PathFile: d.PathFile,
                })
        }

        // Ambil data penilaian untuk generate PDF Lembar Penilaian
        rawPenilaian, _ := h.userRepo.GetPenilaianPeserta(c.Request.Context(), id)
        var penilaianLampiran *service.PenilaianLampiranData
        if rawPenilaian != nil {
                kejuruan := make([]service.KejuruanItem, 0, len(rawPenilaian.Kejuruan))
                for _, k := range rawPenilaian.Kejuruan {
                        kejuruan = append(kejuruan, service.KejuruanItem{
                                Komponen: k.Komponen,
                                Nilai:    k.Nilai,
                        })
                }
                penilaianLampiran = &service.PenilaianLampiranData{
                        NamaLengkap:        rawPenilaian.NamaLengkap,
                        NomorInduk:         rawPenilaian.NomorInduk,
                        AsalInstitusi:      rawPenilaian.AsalInstitusi,
                        Jurusan:            rawPenilaian.Jurusan,
                        KelasSemester:      rawPenilaian.KelasSemester,
                        Divisi:             rawPenilaian.Divisi,
                        Pembimbing:         rawPenilaian.Pembimbing,
                        Periode:            rawPenilaian.Periode,
                        ManagerNama:        rawPenilaian.ManagerNama,
                        Catatan:            rawPenilaian.Catatan,
                        DinilaiAt:          rawPenilaian.DinilaiAt,
                        NilaiMotivasi:      rawPenilaian.NilaiMotivasi,
                        NilaiInisiatif:     rawPenilaian.NilaiInisiatif,
                        NilaiDisiplinWaktu: rawPenilaian.NilaiDisiplinWaktu,
                        NilaiKerajinan:     rawPenilaian.NilaiKerajinan,
                        NilaiKreativitas:   rawPenilaian.NilaiKreativitas,
                        NilaiTanggungJawab: rawPenilaian.NilaiTanggungJawab,
                        NilaiKerjasama:     rawPenilaian.NilaiKerjasama,
                        NilaiAdaptasi:      rawPenilaian.NilaiAdaptasi,
                        NilaiKehadiran:     rawPenilaian.NilaiKehadiran,
                        NilaiK3Safety:      rawPenilaian.NilaiK3Safety,
                        NilaiK3Metode:      rawPenilaian.NilaiK3Metode,
                        NilaiK3Manajemen:   rawPenilaian.NilaiK3Manajemen,
                        NilaiK3Volume:      rawPenilaian.NilaiK3Volume,
                        NilaiPrsProses:     rawPenilaian.NilaiPrsProses,
                        NilaiPrsTeori:      rawPenilaian.NilaiPrsTeori,
                        NilaiPrsJudul:      rawPenilaian.NilaiPrsJudul,
                        NilaiPrsData:       rawPenilaian.NilaiPrsData,
                        NilaiAkhir:         rawPenilaian.NilaiAkhir,
                        Kejuruan:           kejuruan,
                }
        }

        // Siapkan data email backup
        emailSvc := h.emailSvc
        emailCopy := info.Email
        namaCopy := info.NamaLengkap
        divisiStr := ""
        if info.Divisi != nil {
                divisiStr = *info.Divisi
        }
        pembimbingStr := ""
        if info.PembimbingNama != nil {
                pembimbingStr = *info.PembimbingNama
        }
        periode := ""
        if info.TanggalMulai != nil && info.TanggalSelesai != nil {
                periode = *info.TanggalMulai + " \u2013 " + *info.TanggalSelesai
        }
        nilaiStr := ""
        if info.Nilai != nil {
                nilaiStr = fmt.Sprintf("%.1f", *info.Nilai)
        }
        sertRelPath := ""
        if info.SertifikatPath != nil {
                sertRelPath = *info.SertifikatPath
        }
        uploadDir := config.App.UploadDir

        // Kirim email backup (sync, sebelum hapus — agar email pasti terkirim)
        if emailSvc != nil {
                if errEmail := emailSvc.KirimBackupHapusAkun(emailCopy, namaCopy, divisiStr, pembimbingStr, periode, nilaiStr, uploadDir, sertRelPath, dokumenLampiran, penilaianLampiran); errEmail != nil {
                        fmt.Printf("[hapus-peserta] gagal kirim backup email ke %s: %v\n", emailCopy, errEmail)
                        // Lanjutkan hapus meski email gagal (tidak block operasi)
                } else {
                        fmt.Printf("[hapus-peserta] backup email terkirim ke %s\n", emailCopy)
                }
        }

        // Kirim force_logout via WebSocket jika peserta sedang online (sebelum delete)
        appws.GlobalHub.SendToUser(id, "force_logout", nil)

        // Hapus semua data peserta dalam satu transaksi
        bgCtx := context.Background()
        if err := h.userRepo.HapusPesertaData(bgCtx, id); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "delete_failed", Message: "Gagal menghapus data peserta: " + err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Akun peserta dan seluruh datanya berhasil dihapus"})
}

// GET /api/admin/cleanup/preview — daftar akun peserta yang akan dihapus otomatis dalam N hari ke depan
// Query param: days=7 (default 7)
func (h *AdminHandler) GetAkunJatuhTempo(c *gin.Context) {
        days := 7
        if d := c.Query("days"); d != "" {
                if n, err := time.ParseDuration(d + "h"); err == nil {
                        days = int(n.Hours() / 24)
                } else {
                        // coba parse int langsung
                        var n int
                        if _, err2 := fmt.Sscanf(d, "%d", &n); err2 == nil && n > 0 {
                                days = n
                        }
                }
        }

        list, err := h.userRepo.GetAkunJatuhTempo(c.Request.Context(), days)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

func parseTimeStr(s string) (time.Time, error) {
        return time.Parse("2006-01-02", s)
}
