package handler

import (
        "bytes"
        "fmt"
        "log"
        "math"
        "net/http"
        "strings"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/jung-kurt/gofpdf"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
)

type AbsensiHandler struct {
        repo            *repository.AbsensiRepository
        pelaksanaanRepo *repository.PelaksanaanRepository
        pengajuanRepo   *repository.PengajuanRepository
        configRepo      *repository.AbsensiConfigRepository
        divisiRepo      *repository.DivisiRepository
}

func NewAbsensiHandler(
        repo *repository.AbsensiRepository,
        pelaksanaanRepo *repository.PelaksanaanRepository,
        pengajuanRepo *repository.PengajuanRepository,
        configRepo *repository.AbsensiConfigRepository,
        divisiRepo *repository.DivisiRepository,
) *AbsensiHandler {
        return &AbsensiHandler{repo: repo, pelaksanaanRepo: pelaksanaanRepo, pengajuanRepo: pengajuanRepo, configRepo: configRepo, divisiRepo: divisiRepo}
}

// haversineM menghitung jarak antara dua koordinat GPS dalam meter (rumus Haversine)
func haversineM(lat1, lng1, lat2, lng2 float64) float64 {
        const R = 6371000.0
        dLat := (lat2 - lat1) * math.Pi / 180
        dLng := (lng2 - lng1) * math.Pi / 180
        a := math.Sin(dLat/2)*math.Sin(dLat/2) +
                math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
                        math.Sin(dLng/2)*math.Sin(dLng/2)
        return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// wibNow returns current time in WIB (Asia/Jakarta, UTC+7)
func wibNow() time.Time {
        wib, _ := time.LoadLocation("Asia/Jakarta")
        return time.Now().In(wib)
}

// parseWindowTime parses "HH:MM" string into time.Time on the same day as base
func parseWindowTime(hhMM string, base time.Time) time.Time {
        t, err := time.ParseInLocation("2006-01-02 15:04", base.Format("2006-01-02")+" "+hhMM, base.Location())
        if err != nil {
                return base
        }
        return t
}

// GET /api/absensi/geo-status — cek apakah divisi peserta aktif & geofence terkonfigurasi
func (h *AbsensiHandler) GeoStatus(c *gin.Context) {
        userID := middleware.GetUserID(c)
        pelaksanaan, err := h.pelaksanaanRepo.FindByUserID(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusOK, gin.H{
                        "can_checkin": false,
                        "code":        "pelaksanaan_not_found",
                        "message":     "Data magang Anda tidak ditemukan. Hubungi HRD.",
                })
                return
        }
        if strings.TrimSpace(pelaksanaan.Divisi) == "" {
                c.JSON(http.StatusOK, gin.H{
                        "can_checkin": false,
                        "code":        "divisi_not_set",
                        "message":     "Divisi Anda belum diatur oleh HRD. Absensi tidak dapat dilakukan.",
                })
                return
        }
        divisiNama := strings.TrimSpace(pelaksanaan.Divisi)
        divisi, lookupErr := h.divisiRepo.FindByNama(c.Request.Context(), divisiNama)
        if lookupErr != nil || divisi == nil {
                c.JSON(http.StatusOK, gin.H{
                        "can_checkin": false,
                        "code":        "divisi_not_found",
                        "message":     fmt.Sprintf("Divisi '%s' tidak terdaftar dalam sistem. Hubungi HRD.", divisiNama),
                })
                return
        }
        if !divisi.IsActive {
                c.JSON(http.StatusOK, gin.H{
                        "can_checkin": false,
                        "code":        "divisi_inactive",
                        "message":     "Divisi tempat magang Anda saat ini tidak aktif. Hubungi HRD.",
                })
                return
        }
        if divisi.GeoLat == nil || divisi.GeoLng == nil || (*divisi.GeoLat == 0 && *divisi.GeoLng == 0) {
                c.JSON(http.StatusOK, gin.H{
                        "can_checkin": false,
                        "code":        "geofence_not_configured",
                        "message":     "Koordinat lokasi absensi untuk divisi Anda belum dikonfigurasi. Hubungi admin/HRD.",
                })
                return
        }
        geoRadius := float64(1500)
        if divisi.GeoRadius != nil && *divisi.GeoRadius > 0 {
                geoRadius = float64(*divisi.GeoRadius)
        }
        c.JSON(http.StatusOK, gin.H{
                "can_checkin": true,
                "code":        "",
                "message":     "",
                "geo_lat":     *divisi.GeoLat,
                "geo_lng":     *divisi.GeoLng,
                "geo_radius":  geoRadius,
        })
}

// POST /api/absensi/checkin
func (h *AbsensiHandler) CheckIn(c *gin.Context) {
        var req models.AbsensiCheckInRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "Data tidak valid"})
                return
        }
        if req.Latitude == 0 && req.Longitude == 0 {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "location_required", Message: "Data lokasi GPS diperlukan untuk absensi"})
                return
        }

        userID := middleware.GetUserID(c)
        now := wibNow()

        cfg, err := h.configRepo.Get(c.Request.Context())
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal memuat konfigurasi jam absensi"})
                return
        }

        // Validasi window waktu
        buka := parseWindowTime(cfg.JamMasukBuka, now)
        tutup := parseWindowTime(cfg.JamMasukTutup, now)
        if now.Before(buka) || now.After(tutup) {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{
                        Error:   "outside_window",
                        Message: fmt.Sprintf("Absen masuk hanya bisa dilakukan antara %s - %s WIB", cfg.JamMasukBuka, cfg.JamMasukTutup),
                })
                return
        }

        pelaksanaan, err := h.pelaksanaanRepo.FindByUserID(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "not_found", Message: "Data pelaksanaan magang tidak ditemukan"})
                return
        }
        if pelaksanaan.Status != models.StatusAktif {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_status", Message: "Magang belum aktif"})
                return
        }

        // Catatan: cek akurasi GPS sengaja dihapus.
        // Akurasi di dalam gedung sering > 300 m — cek akurasi memblokir pengguna sah.
        // Geofence jarak (haversine vs radius divisi) sudah cukup sebagai filter lokasi.

        // ── GEOFENCE WAJIB ────────────────────────────────────────────────────────
        // Tidak ada jalan pintas: SETIAP absensi harus lolos semua cek ini.
        if strings.TrimSpace(pelaksanaan.Divisi) == "" {
                c.JSON(http.StatusForbidden, models.ErrorResponse{
                        Error:   "divisi_not_set",
                        Message: "Divisi Anda belum diatur oleh HRD. Absensi tidak dapat dilakukan. Hubungi HRD.",
                })
                return
        }
        divisiNama := strings.TrimSpace(pelaksanaan.Divisi)
        divisi, lookupErr := h.divisiRepo.FindByNama(c.Request.Context(), divisiNama)
        if lookupErr != nil || divisi == nil {
                c.JSON(http.StatusForbidden, models.ErrorResponse{
                        Error:   "divisi_not_found",
                        Message: fmt.Sprintf("Divisi '%s' tidak terdaftar dalam sistem. Hubungi HRD.", divisiNama),
                })
                return
        }
        if !divisi.IsActive {
                c.JSON(http.StatusForbidden, models.ErrorResponse{
                        Error:   "divisi_inactive",
                        Message: "Divisi Anda saat ini tidak aktif. Hubungi HRD.",
                })
                return
        }
        if divisi.GeoLat == nil || divisi.GeoLng == nil || (*divisi.GeoLat == 0 && *divisi.GeoLng == 0) {
                c.JSON(http.StatusForbidden, models.ErrorResponse{
                        Error:   "geofence_not_configured",
                        Message: "Koordinat lokasi absensi untuk divisi Anda belum dikonfigurasi. Hubungi admin/HRD.",
                })
                return
        }
        geoRadius := float64(1500)
        if divisi.GeoRadius != nil && *divisi.GeoRadius > 0 {
                geoRadius = float64(*divisi.GeoRadius)
        }
        namaLokasi := divisi.Nama
        if divisi.NamaLokasi != nil && strings.TrimSpace(*divisi.NamaLokasi) != "" {
                namaLokasi = *divisi.NamaLokasi
        }
        jarakM := haversineM(req.Latitude, req.Longitude, *divisi.GeoLat, *divisi.GeoLng)
        log.Printf("[GEOFENCE CHECKIN] user=%s divisi=%q posisi=(%.6f,%.6f) kantor=(%.6f,%.6f) jarak=%.1fm radius=%.0fm akurasi=%.0fm",
                userID, divisiNama, req.Latitude, req.Longitude, *divisi.GeoLat, *divisi.GeoLng, jarakM, geoRadius, req.Accuracy)
        if jarakM > geoRadius {
                c.JSON(http.StatusForbidden, models.ErrorResponse{
                        Error:   "outside_geofence",
                        Message: fmt.Sprintf("Anda berada %.0f m dari lokasi %s (radius %.0f m). Harap absen dari lokasi kantor.", jarakM, namaLokasi, geoRadius),
                })
                return
        }
        // ── END GEOFENCE ──────────────────────────────────────────────────────────

        tanggal := now.Format("2006-01-02")
        jamMasuk := now.Format("15:04")
        keterangan := "hadir"

        a := &models.Absensi{
                PelaksanaanID: pelaksanaan.ID,
                Tanggal:       tanggal,
                JamMasuk:      &jamMasuk,
                Keterangan:    keterangan,
                Latitude:      &req.Latitude,
                Longitude:     &req.Longitude,
        }

        if err := h.repo.CheckIn(c.Request.Context(), a); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "checkin_failed", Message: "Sudah check-in hari ini atau terjadi kesalahan"})
                return
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: "Absen masuk berhasil",
                Data:    a,
        })
}

// PATCH /api/absensi/checkout
func (h *AbsensiHandler) CheckOut(c *gin.Context) {
        var req models.AbsensiCheckOutRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        // GPS wajib — sama seperti check-in
        if req.Latitude == 0 && req.Longitude == 0 {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "location_required", Message: "Data lokasi GPS diperlukan untuk absen pulang"})
                return
        }

        userID := middleware.GetUserID(c)
        now := wibNow()

        cfg, err := h.configRepo.Get(c.Request.Context())
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal memuat konfigurasi jam absensi"})
                return
        }

        buka := parseWindowTime(cfg.JamPulangBuka, now)
        tutup := parseWindowTime(cfg.JamPulangTutup, now)
        if now.Before(buka) || now.After(tutup) {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{
                        Error:   "outside_window",
                        Message: fmt.Sprintf("Absen pulang hanya bisa dilakukan antara %s - %s WIB", cfg.JamPulangBuka, cfg.JamPulangTutup),
                })
                return
        }

        pelaksanaan, err := h.pelaksanaanRepo.FindByUserID(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "not_found", Message: "Data pelaksanaan tidak ditemukan"})
                return
        }

        // Catatan: cek akurasi GPS sengaja dihapus — sama alasannya dengan check-in.
        // Geofence jarak yang menentukan, bukan nilai akurasi GPS.

        // ── GEOFENCE WAJIB (CHECKOUT) ─────────────────────────────────────────────
        if strings.TrimSpace(pelaksanaan.Divisi) == "" {
                c.JSON(http.StatusForbidden, models.ErrorResponse{
                        Error:   "divisi_not_set",
                        Message: "Divisi Anda belum diatur oleh HRD. Absensi tidak dapat dilakukan. Hubungi HRD.",
                })
                return
        }
        divisiNamaOut := strings.TrimSpace(pelaksanaan.Divisi)
        divisiOut, lookupErrOut := h.divisiRepo.FindByNama(c.Request.Context(), divisiNamaOut)
        if lookupErrOut != nil || divisiOut == nil {
                c.JSON(http.StatusForbidden, models.ErrorResponse{
                        Error:   "divisi_not_found",
                        Message: fmt.Sprintf("Divisi '%s' tidak terdaftar dalam sistem. Hubungi HRD.", divisiNamaOut),
                })
                return
        }
        if !divisiOut.IsActive {
                c.JSON(http.StatusForbidden, models.ErrorResponse{
                        Error:   "divisi_inactive",
                        Message: "Divisi Anda saat ini tidak aktif. Hubungi HRD.",
                })
                return
        }
        if divisiOut.GeoLat == nil || divisiOut.GeoLng == nil || (*divisiOut.GeoLat == 0 && *divisiOut.GeoLng == 0) {
                c.JSON(http.StatusForbidden, models.ErrorResponse{
                        Error:   "geofence_not_configured",
                        Message: "Koordinat lokasi absensi untuk divisi Anda belum dikonfigurasi. Hubungi admin/HRD.",
                })
                return
        }
        geoRadiusOut := float64(1500)
        if divisiOut.GeoRadius != nil && *divisiOut.GeoRadius > 0 {
                geoRadiusOut = float64(*divisiOut.GeoRadius)
        }
        namaLokasiOut := divisiOut.Nama
        if divisiOut.NamaLokasi != nil && strings.TrimSpace(*divisiOut.NamaLokasi) != "" {
                namaLokasiOut = *divisiOut.NamaLokasi
        }
        jarakMOut := haversineM(req.Latitude, req.Longitude, *divisiOut.GeoLat, *divisiOut.GeoLng)
        log.Printf("[GEOFENCE CHECKOUT] user=%s divisi=%q posisi=(%.6f,%.6f) kantor=(%.6f,%.6f) jarak=%.1fm radius=%.0fm akurasi=%.0fm",
                userID, divisiNamaOut, req.Latitude, req.Longitude, *divisiOut.GeoLat, *divisiOut.GeoLng, jarakMOut, geoRadiusOut, req.Accuracy)
        if jarakMOut > geoRadiusOut {
                c.JSON(http.StatusForbidden, models.ErrorResponse{
                        Error:   "outside_geofence",
                        Message: fmt.Sprintf("Anda berada %.0f m dari lokasi %s (radius %.0f m). Harap absen dari lokasi kantor.", jarakMOut, namaLokasiOut, geoRadiusOut),
                })
                return
        }
        // ── END GEOFENCE ──────────────────────────────────────────────────────────
        latPtr := &req.Latitude
        lngPtr := &req.Longitude

        today := now.Format("2006-01-02")
        jamKeluar := now.Format("15:04")

        if err := h.repo.CheckOut(c.Request.Context(), pelaksanaan.ID, today, jamKeluar, req.Kegiatan, latPtr, lngPtr); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "checkout_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Absen pulang berhasil"})
}

// GET /api/absensi/saya — rekap absensi peserta
func (h *AbsensiHandler) GetMy(c *gin.Context) {
        userID := middleware.GetUserID(c)
        pelaksanaan, err := h.pelaksanaanRepo.FindByUserID(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusOK, models.SuccessResponse{Data: map[string]interface{}{"list": []interface{}{}, "rekap": map[string]int{}}})
                return
        }

        list, err := h.repo.FindByPelaksanaanID(c.Request.Context(), pelaksanaan.ID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        hadir, izin, sakit, alpha := h.repo.CountByPelaksanaan(c.Request.Context(), pelaksanaan.ID)

        c.JSON(http.StatusOK, models.SuccessResponse{
                Data: map[string]interface{}{
                        "list":  list,
                        "rekap": map[string]int{"hadir": hadir, "izin": izin, "sakit": sakit, "alpha": alpha},
                },
        })
}

// GET /api/absensi/pelaksanaan/:id — HRD lihat absensi peserta tertentu
func (h *AbsensiHandler) GetByPelaksanaan(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        list, err := h.repo.FindByPelaksanaanID(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        hadir, izin, sakit, alpha := h.repo.CountByPelaksanaan(c.Request.Context(), id)
        c.JSON(http.StatusOK, models.SuccessResponse{
                Data: map[string]interface{}{
                        "list":  list,
                        "rekap": map[string]int{"hadir": hadir, "izin": izin, "sakit": sakit, "alpha": alpha},
                },
        })
}

// GET /api/absensi/rekap — HRD melihat ringkasan absensi semua peserta
func (h *AbsensiHandler) GetRekapHRD(c *gin.Context) {
        list, err := h.repo.GetRekapAll(c.Request.Context())
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// GET /api/absensi/pelaksanaan/:id/pdf — HRD download PDF rekap absensi peserta tertentu
func (h *AbsensiHandler) GetByPelaksanaanPDF(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        pelaksanaan, err := h.pelaksanaanRepo.FindByID(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: "Data pelaksanaan tidak ditemukan"})
                return
        }

        pengajuan, err := h.pengajuanRepo.FindByID(c.Request.Context(), pelaksanaan.PengajuanID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        list, _ := h.repo.FindByPelaksanaanID(c.Request.Context(), id)
        hadir, izin, sakit, alpha := h.repo.CountByPelaksanaan(c.Request.Context(), id)

        pdf := generateAbsensiPDF(pengajuan, pelaksanaan, list, hadir, izin, sakit, alpha)

        var buf bytes.Buffer
        pdf.Output(&buf)

        c.Header("Content-Type", "application/pdf")
        c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="Absensi_%s.pdf"`, pengajuan.NamaLengkap))
        c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}

// PATCH /api/absensi/:id/approve — HRD setujui absensi
func (h *AbsensiHandler) Approve(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        approvedBy := middleware.GetUserID(c)
        if err := h.repo.Approve(c.Request.Context(), id, approvedBy); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Absensi berhasil disetujui"})
}

// GET /api/absensi/saya/pdf — download rekap PDF absensi
func (h *AbsensiHandler) DownloadPDF(c *gin.Context) {
        userID := middleware.GetUserID(c)
        pelaksanaan, err := h.pelaksanaanRepo.FindByUserID(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "not_found", Message: "Data pelaksanaan tidak ditemukan"})
                return
        }

        pengajuan, err := h.pengajuanRepo.FindByID(c.Request.Context(), pelaksanaan.PengajuanID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        list, _ := h.repo.FindByPelaksanaanID(c.Request.Context(), pelaksanaan.ID)
        hadir, izin, sakit, alpha := h.repo.CountByPelaksanaan(c.Request.Context(), pelaksanaan.ID)

        pdf := generateAbsensiPDF(pengajuan, pelaksanaan, list, hadir, izin, sakit, alpha)

        var buf bytes.Buffer
        pdf.Output(&buf)

        c.Header("Content-Type", "application/pdf")
        c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="Absensi_%s.pdf"`, pengajuan.NamaLengkap))
        c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}

// formatJamShort memotong format PostgreSQL "HH:MM:SS.ffffff" menjadi "HH:MM"
func formatJamShort(s string) string {
        if len(s) >= 5 {
                return s[:5]
        }
        return s
}

// splitKegiatan memecah teks kegiatan per baris (filter baris kosong)
func splitKegiatan(s string) []string {
        if s == "" {
                return nil
        }
        var result []string
        for _, line := range strings.Split(s, "\n") {
                line = strings.TrimSpace(line)
                if line != "" {
                        result = append(result, line)
                }
        }
        return result
}

// absensiRowPDF mewakili satu baris pada tabel absensi di PDF
type absensiRowPDF struct {
        No        int
        Tanggal   string
        Hari      string
        JamMasuk  string
        JamKeluar string
        Status    string
        Kegiatan  string
        TTD       bool
}

// buildAbsensiRows membuat semua baris dari rentang tanggal mulai–selesai (hanya hari kerja)
func buildAbsensiRows(pel *models.PelaksanaanMagang, list []models.Absensi) []absensiRowPDF {
        absensiMap := make(map[string]models.Absensi)
        for _, a := range list {
                absensiMap[a.Tanggal] = a
        }

        today := wibNow().Format("2006-01-02")
        dayNames := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}

        start := time.Date(pel.TanggalMulai.Year(), pel.TanggalMulai.Month(), pel.TanggalMulai.Day(), 0, 0, 0, 0, time.UTC)
        end := time.Date(pel.TanggalSelesai.Year(), pel.TanggalSelesai.Month(), pel.TanggalSelesai.Day(), 0, 0, 0, 0, time.UTC)

        var rows []absensiRowPDF
        no := 1
        for cur := start; !cur.After(end); cur = cur.AddDate(0, 0, 1) {
                dow := int(cur.Weekday())
                if dow == 0 || dow == 6 {
                        continue
                }
                dateStr := cur.Format("2006-01-02")
                row := absensiRowPDF{
                        No:      no,
                        Tanggal: cur.Format("02/01/2006"),
                        Hari:    dayNames[dow],
                }
                if a, ok := absensiMap[dateStr]; ok {
                        if a.JamMasuk != nil {
                                row.JamMasuk = formatJamShort(*a.JamMasuk)
                        }
                        if a.JamKeluar != nil {
                                row.JamKeluar = formatJamShort(*a.JamKeluar)
                        }
                        switch a.Keterangan {
                        case "hadir":
                                row.Status = "Hadir"
                        case "izin":
                                row.Status = "Izin"
                        case "sakit":
                                row.Status = "Sakit"
                        case "alpha":
                                row.Status = "Alpha"
                        default:
                                row.Status = a.Keterangan
                        }
                        if a.Kegiatan != "" {
                                row.Kegiatan = a.Kegiatan
                        }
                        row.TTD = a.TTDPembimbing
                } else if dateStr < today {
                        row.Status = "Alpha"
                }
                rows = append(rows, row)
                no++
        }
        return rows
}

func generateAbsensiPDF(pengajuan *models.PengajuanMagang, pel *models.PelaksanaanMagang, list []models.Absensi, hadir, izin, sakit, alpha int) *gofpdf.Fpdf {
        pdf := gofpdf.New("P", "mm", "A4", "")
        pdf.SetMargins(18, 14, 18)
        pdf.SetAutoPageBreak(true, 18)
        pdf.AddPage()

        const (
                leftX = 18.0
                pageW = 174.0 // 210 - 18 - 18
        )

        // ── LETTERHEAD ────────────────────────────────────────────
        // Kotak hijau tua di kiri sebagai aksen
        pdf.SetFillColor(22, 101, 52)
        pdf.Rect(leftX, 14, 6, 24, "F")

        // Nama perusahaan
        pdf.SetTextColor(22, 101, 52)
        pdf.SetFont("Helvetica", "B", 13)
        pdf.SetXY(27, 15)
        pdf.CellFormat(165, 7, "PT TANJUNGENIM LESTARI PULP AND PAPER", "", 2, "L", false, 0, "")

        // Sub-judul
        pdf.SetTextColor(75, 85, 99)
        pdf.SetFont("Helvetica", "", 9)
        pdf.SetX(27)
        pdf.CellFormat(165, 5, "Sistem Manajemen Magang  —  e-Magang PT TELPP", "", 2, "L", false, 0, "")
        pdf.SetX(27)
        pdf.CellFormat(165, 5, "Muara Enim, Sumatera Selatan", "", 0, "L", false, 0, "")

        // Garis pembatas hijau
        pdf.SetDrawColor(22, 101, 52)
        pdf.SetLineWidth(0.6)
        pdf.Line(leftX, 40, leftX+pageW, 40)

        // ── JUDUL DOKUMEN ─────────────────────────────────────────
        pdf.SetTextColor(17, 24, 39)
        pdf.SetFont("Helvetica", "B", 13)
        pdf.SetXY(leftX, 43)
        pdf.CellFormat(pageW, 8, "REKAP ABSENSI MAGANG", "", 1, "C", false, 0, "")

        pdf.SetDrawColor(209, 213, 219)
        pdf.SetLineWidth(0.2)
        pdf.Line(leftX, 52, leftX+pageW, 52)

        // ── DATA PESERTA ─────────────────────────────────────────
        yPos := 56.0

        pdf.SetTextColor(107, 114, 128)
        pdf.SetFont("Helvetica", "B", 8)
        pdf.SetXY(leftX, yPos)
        pdf.CellFormat(pageW, 5, "DATA PESERTA", "", 1, "L", false, 0, "")
        yPos += 6

        divisiVal := "-"
        if pel.Divisi != "" {
                divisiVal = pel.Divisi
        }

        infoRows := [][]string{
                {"Nama Peserta", pengajuan.NamaLengkap},
                {"Asal Institusi", pengajuan.AsalInstitusi},
                {"Jurusan / Program Studi", pengajuan.Jurusan},
                {"Divisi / Bagian", divisiVal},
                {"Periode Magang", fmt.Sprintf("%s  s/d  %s",
                        pel.TanggalMulai.Format("02 Januari 2006"),
                        pel.TanggalSelesai.Format("02 Januari 2006"))},
                {"Tanggal Dicetak", wibNow().Format("02 Januari 2006")},
        }

        for _, row := range infoRows {
                pdf.SetFont("Helvetica", "", 9)
                pdf.SetTextColor(75, 85, 99)
                pdf.SetXY(leftX, yPos)
                pdf.CellFormat(54, 5.5, row[0], "", 0, "L", false, 0, "")
                pdf.SetTextColor(55, 65, 81)
                pdf.CellFormat(5, 5.5, ":", "", 0, "C", false, 0, "")
                pdf.SetFont("Helvetica", "B", 9)
                pdf.SetTextColor(17, 24, 39)
                pdf.CellFormat(115, 5.5, row[1], "", 1, "L", false, 0, "")
                yPos += 5.5
        }

        yPos += 5

        // ── RINGKASAN KEHADIRAN ───────────────────────────────────
        pdf.SetTextColor(107, 114, 128)
        pdf.SetFont("Helvetica", "B", 8)
        pdf.SetXY(leftX, yPos)
        pdf.CellFormat(pageW, 5, "RINGKASAN KEHADIRAN", "", 1, "L", false, 0, "")
        yPos += 6

        type statBox struct {
                label     string
                val       int
                fr, fg, fb int // fill color
                tr, tg, tb int // text color
        }
        boxes := []statBox{
                {"Hadir", hadir, 240, 253, 244, 22, 101, 52},
                {"Izin", izin, 254, 252, 232, 161, 98, 7},
                {"Sakit", sakit, 239, 246, 255, 29, 78, 216},
                {"Alpha", alpha, 255, 241, 242, 190, 18, 60},
        }
        boxW := pageW / 4
        for i, b := range boxes {
                x := leftX + float64(i)*boxW
                pdf.SetFillColor(b.fr, b.fg, b.fb)
                pdf.SetDrawColor(209, 213, 219)
                pdf.SetLineWidth(0.2)
                pdf.Rect(x, yPos, boxW-0.5, 14, "FD")
                pdf.SetTextColor(b.tr, b.tg, b.tb)
                pdf.SetFont("Helvetica", "B", 15)
                pdf.SetXY(x, yPos+1)
                pdf.CellFormat(boxW-0.5, 7, fmt.Sprintf("%d", b.val), "", 1, "C", false, 0, "")
                pdf.SetFont("Helvetica", "", 8)
                pdf.SetTextColor(107, 114, 128)
                pdf.SetXY(x, yPos+8)
                pdf.CellFormat(boxW-0.5, 5, b.label, "", 0, "C", false, 0, "")
        }
        yPos += 20

        // ── TABEL DETAIL ABSENSI ──────────────────────────────────
        pdf.SetTextColor(107, 114, 128)
        pdf.SetFont("Helvetica", "B", 8)
        pdf.SetXY(leftX, yPos)
        pdf.CellFormat(pageW, 5, "DETAIL ABSENSI", "", 1, "L", false, 0, "")
        yPos += 6

        colWidths := []float64{8, 22, 16, 15, 15, 17, 71, 10}
        tblHeaders := []string{"No", "Tanggal", "Hari", "Masuk", "Pulang", "Status", "Kegiatan", "TTD"}

        printTableHeader := func(y float64) {
                pdf.SetFillColor(22, 101, 52)
                pdf.SetTextColor(255, 255, 255)
                pdf.SetFont("Helvetica", "B", 8)
                pdf.SetXY(leftX, y)
                for i, h := range tblHeaders {
                        align := "C"
                        if i == 6 {
                                align = "L"
                        }
                        pdf.CellFormat(colWidths[i], 6, h, "1", 0, align, true, 0, "")
                }
        }
        printTableHeader(yPos)
        yPos += 6

        rows := buildAbsensiRows(pel, list)
        pdf.SetFont("Helvetica", "", 8)

        const (
                lineH   = 4.5  // tinggi satu baris kegiatan
                minRowH = 5.5  // tinggi minimum baris (1 baris kosong)
        )
        // kegiatanX = leftX + lebar kolom 0-5
        kegiatanX := leftX + colWidths[0] + colWidths[1] + colWidths[2] +
                colWidths[3] + colWidths[4] + colWidths[5]
        ttdX := kegiatanX + colWidths[6]

        pdf.SetDrawColor(209, 213, 219)
        pdf.SetLineWidth(0.2)

        for rowIdx, row := range rows {
                poinList := splitKegiatan(row.Kegiatan)
                numLines := len(poinList)
                if numLines == 0 {
                        numLines = 1
                }
                rowH := math.Max(minRowH, float64(numLines)*lineH)

                if yPos+rowH > 262 {
                        pdf.AddPage()
                        yPos = 18
                        printTableHeader(yPos)
                        yPos += 6
                        pdf.SetFont("Helvetica", "", 8)
                        pdf.SetDrawColor(209, 213, 219)
                        pdf.SetLineWidth(0.2)
                }

                // Latar belakang baris
                if rowIdx%2 == 0 {
                        pdf.SetFillColor(249, 250, 251)
                } else {
                        pdf.SetFillColor(255, 255, 255)
                }
                pdf.Rect(leftX, yPos, pageW, rowH, "F")

                // Warna teks berdasarkan status
                switch row.Status {
                case "Alpha":
                        pdf.SetTextColor(190, 18, 60)
                case "Izin":
                        pdf.SetTextColor(161, 98, 7)
                case "Sakit":
                        pdf.SetTextColor(29, 78, 216)
                default:
                        pdf.SetTextColor(55, 65, 81)
                }

                ttdStr := ""
                if row.TTD {
                        ttdStr = "V"
                }

                // Kolom 0–5 (No, Tanggal, Hari, Masuk, Pulang, Status) — tinggi penuh
                nonKeg := []string{
                        fmt.Sprintf("%d", row.No),
                        row.Tanggal, row.Hari,
                        row.JamMasuk, row.JamKeluar,
                        row.Status,
                }
                pdf.SetXY(leftX, yPos)
                for i, cell := range nonKeg {
                        pdf.CellFormat(colWidths[i], rowH, cell, "1", 0, "C", false, 0, "")
                }

                // Kolom Kegiatan — border kotak, tiap poin per baris
                pdf.Rect(kegiatanX, yPos, colWidths[6], rowH, "D")
                pdf.SetTextColor(55, 65, 81)
                maxKegW := colWidths[6] - 3.0
                for j, poin := range poinList {
                        lineText := fmt.Sprintf("%d. %s", j+1, poin)
                        // Potong teks jika terlalu panjang
                        for pdf.GetStringWidth(lineText) > maxKegW && len(lineText) > 4 {
                                lineText = lineText[:len(lineText)-1]
                        }
                        pdf.SetXY(kegiatanX+1.5, yPos+float64(j)*lineH+0.5)
                        pdf.CellFormat(maxKegW, lineH, lineText, "", 0, "L", false, 0, "")
                }

                // Kolom TTD — kembalikan warna status
                switch row.Status {
                case "Alpha":
                        pdf.SetTextColor(190, 18, 60)
                case "Izin":
                        pdf.SetTextColor(161, 98, 7)
                case "Sakit":
                        pdf.SetTextColor(29, 78, 216)
                default:
                        pdf.SetTextColor(55, 65, 81)
                }
                pdf.SetXY(ttdX, yPos)
                pdf.CellFormat(colWidths[7], rowH, ttdStr, "1", 0, "C", false, 0, "")

                yPos += rowH
        }

        // ── TANDA TANGAN ─────────────────────────────────────────
        yPos += 14
        if yPos > 256 {
                pdf.AddPage()
                yPos = 30
        }

        pdf.SetTextColor(55, 65, 81)
        pdf.SetFont("Helvetica", "", 10)
        pdf.SetXY(128, yPos)
        pdf.CellFormat(64, 5, fmt.Sprintf("Muara Enim, %s", wibNow().Format("02 Januari 2006")), "", 1, "C", false, 0, "")
        pdf.SetXY(128, yPos+6)
        pdf.CellFormat(64, 5, "Mengetahui,", "", 1, "C", false, 0, "")
        pdf.SetXY(128, yPos+11)
        pdf.SetFont("Helvetica", "I", 9)
        pdf.SetTextColor(107, 114, 128)
        pdf.CellFormat(64, 5, "Pembimbing Magang", "", 0, "C", false, 0, "")

        pdf.SetDrawColor(107, 114, 128)
        pdf.SetLineWidth(0.3)
        pdf.Line(133, yPos+34, 190, yPos+34)
        pdf.SetXY(128, yPos+36)
        pdf.SetFont("Helvetica", "", 9)
        pdf.SetTextColor(55, 65, 81)
        pdf.CellFormat(64, 5, "PT TanjungEnim Lestari Pulp and Paper", "", 0, "C", false, 0, "")

        return pdf
}

func parseTimeOnly(s string) time.Time {
        t, _ := time.Parse("2006-01-02", s)
        return t
}
