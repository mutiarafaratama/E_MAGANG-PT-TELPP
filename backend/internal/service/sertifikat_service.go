package service

import (
        "context"
        "errors"
        "fmt"
        "os"
        "path/filepath"
        "time"

        "github.com/google/uuid"
        "github.com/jung-kurt/gofpdf"
        "github.com/telpp/emagang/internal/config"
        "github.com/telpp/emagang/internal/database"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
)

type SertifikatService struct {
        pelaksanaanRepo *repository.PelaksanaanRepository
        pengajuanRepo   *repository.PengajuanRepository
        notifSvc        *NotifikasiService
}

func NewSertifikatService(
        pelaksanaanRepo *repository.PelaksanaanRepository,
        pengajuanRepo *repository.PengajuanRepository,
        notifSvc *NotifikasiService,
) *SertifikatService {
        return &SertifikatService{
                pelaksanaanRepo: pelaksanaanRepo,
                pengajuanRepo:   pengajuanRepo,
                notifSvc:        notifSvc,
        }
}

func (s *SertifikatService) Generate(ctx context.Context, pelaksanaanID uuid.UUID) (string, error) {
        // Ambil data pelaksanaan langsung dari DB
        var p models.PelaksanaanMagang
        err := database.DB.QueryRow(ctx, `
                SELECT id, pengajuan_id, user_id, periode_id, tanggal_mulai, tanggal_selesai,
                divisi, pembimbing_id, status, nilai, catatan_nilai, dinilai_oleh, dinilai_at,
                sertifikat_generated, sertifikat_path, sertifikat_generated_at, created_at, updated_at
                FROM pelaksanaan_magang WHERE id = $1`, pelaksanaanID).
                Scan(&p.ID, &p.PengajuanID, &p.UserID, &p.PeriodeID, &p.TanggalMulai, &p.TanggalSelesai,
                        &p.Divisi, &p.PembimbingID, &p.Status, &p.Nilai, &p.CatatanNilai, &p.DinilaiOleh, &p.DinilaiAt,
                        &p.SertifikatGenerated, &p.SertifikatPath, &p.SertifikatGeneratedAt, &p.CreatedAt, &p.UpdatedAt)
        if err != nil {
                return "", errors.New("data pelaksanaan tidak ditemukan")
        }

        if p.Status != models.StatusPenilaian && p.Status != models.StatusSelesai {
                return "", errors.New("peserta belum mendapatkan penilaian")
        }

        if p.Nilai == nil {
                return "", errors.New("nilai belum diisi oleh pembimbing")
        }

        pengajuan, err := s.pengajuanRepo.FindByID(ctx, p.PengajuanID)
        if err != nil {
                return "", err
        }

        // Generate PDF
        dirPath := filepath.Join(config.App.UploadDir, "sertifikat")
        if err := os.MkdirAll(dirPath, 0755); err != nil {
                return "", fmt.Errorf("gagal membuat folder sertifikat: %w", err)
        }

        filename := fmt.Sprintf("sertifikat_%s_%s.pdf",
                sanitizeFilename(pengajuan.NamaLengkap),
                time.Now().Format("20060102150405"))
        savePath := filepath.Join(dirPath, filename)

        if err := generateSertifikatPDF(pengajuan, &p, savePath); err != nil {
                return "", fmt.Errorf("gagal generate PDF: %w", err)
        }

        // Update path di DB
        if err := s.pelaksanaanRepo.SetSertifikat(ctx, pelaksanaanID, savePath); err != nil {
                return "", err
        }

        // Notif ke peserta — penerima sertifikat selalu peserta
        s.notifSvc.KirimKeUser(ctx, pengajuan.UserID, models.RolePeserta,
                "Sertifikat Magang Siap",
                "Sertifikat magang Anda telah tersedia dan siap untuk didownload",
                string(models.NotifSertifikat), &pelaksanaanID)

        return savePath, nil
}

func generateSertifikatPDF(p *models.PengajuanMagang, pel *models.PelaksanaanMagang, savePath string) error {
        pdf := gofpdf.New("L", "mm", "A4", "")
        pdf.AddPage()

        // Background
        pdf.SetFillColor(255, 255, 255)
        pdf.Rect(0, 0, 297, 210, "F")

        // Border hijau
        pdf.SetDrawColor(0, 128, 0)
        pdf.SetLineWidth(3)
        pdf.Rect(8, 8, 281, 194, "D")
        pdf.SetLineWidth(1)
        pdf.Rect(11, 11, 275, 188, "D")

        // Header bar hijau
        pdf.SetFillColor(0, 100, 0)
        pdf.Rect(8, 8, 281, 32, "F")

        pdf.SetTextColor(255, 255, 255)
        pdf.SetFont("Helvetica", "B", 18)
        pdf.SetXY(8, 14)
        pdf.CellFormat(281, 9, "PT TANJUNGENIM LESTARI PULP AND PAPER", "", 0, "C", false, 0, "")

        pdf.SetFont("Helvetica", "", 11)
        pdf.SetXY(8, 25)
        pdf.CellFormat(281, 7, "Jl. Raya Muara Enim, Sumatera Selatan  |  Telp: (0734) 123-456  |  www.telpp.co.id", "", 0, "C", false, 0, "")

        // Judul sertifikat
        pdf.SetTextColor(0, 80, 0)
        pdf.SetFont("Helvetica", "B", 30)
        pdf.SetXY(8, 48)
        pdf.CellFormat(281, 14, "SERTIFIKAT MAGANG", "", 0, "C", false, 0, "")

        // Nomor
        pdf.SetFont("Helvetica", "", 10)
        pdf.SetTextColor(120, 120, 120)
        nomor := fmt.Sprintf("No. SRTF/TELPP/%s/%04d", time.Now().Format("2006"), time.Now().UnixNano()%9999+1)
        pdf.SetXY(8, 63)
        pdf.CellFormat(281, 7, nomor, "", 0, "C", false, 0, "")

        // Teks pengantar
        pdf.SetTextColor(60, 60, 60)
        pdf.SetFont("Helvetica", "", 12)
        pdf.SetXY(30, 74)
        pdf.CellFormat(237, 7, "Dengan ini menerangkan bahwa:", "", 0, "C", false, 0, "")

        // Nama peserta besar
        pdf.SetFont("Helvetica", "B", 24)
        pdf.SetTextColor(0, 100, 0)
        pdf.SetXY(8, 84)
        pdf.CellFormat(281, 12, p.NamaLengkap, "", 0, "C", false, 0, "")

        // Garis bawah nama
        pdf.SetDrawColor(0, 128, 0)
        pdf.SetLineWidth(0.8)
        pdf.Line(70, 98, 227, 98)

        // Detail dalam 2 kolom
        pdf.SetFont("Helvetica", "", 11)
        pdf.SetTextColor(50, 50, 50)

        divisi := "Umum"
        if pel.Divisi != nil {
                divisi = *pel.Divisi
        }

        leftDetails := [][]string{
                {"Asal Institusi", p.AsalInstitusi},
                {"Jurusan / Prodi", p.Jurusan},
                {"Nomor Induk", p.NomorInduk},
        }
        rightDetails := [][]string{
                {"Divisi Penempatan", divisi},
                {"Periode Magang", fmt.Sprintf("%s s/d %s",
                        pel.TanggalMulai.Format("02 Jan 2006"),
                        pel.TanggalSelesai.Format("02 Jan 2006"))},
                {"Nilai Akhir", fmt.Sprintf("%.1f — Sangat Baik", *pel.Nilai)},
        }

        yPos := 104.0
        for i := 0; i < 3; i++ {
                pdf.SetXY(20, yPos)
                pdf.SetFont("Helvetica", "B", 10)
                pdf.CellFormat(50, 7, leftDetails[i][0], "", 0, "L", false, 0, "")
                pdf.SetFont("Helvetica", "", 10)
                pdf.CellFormat(4, 7, ":", "", 0, "C", false, 0, "")
                pdf.CellFormat(70, 7, leftDetails[i][1], "", 0, "L", false, 0, "")

                pdf.SetXY(158, yPos)
                pdf.SetFont("Helvetica", "B", 10)
                pdf.CellFormat(50, 7, rightDetails[i][0], "", 0, "L", false, 0, "")
                pdf.SetFont("Helvetica", "", 10)
                pdf.CellFormat(4, 7, ":", "", 0, "C", false, 0, "")
                pdf.CellFormat(75, 7, rightDetails[i][1], "", 0, "L", false, 0, "")

                yPos += 8
        }

        // Kalimat penutup
        pdf.SetFont("Helvetica", "I", 11)
        pdf.SetTextColor(60, 60, 60)
        pdf.SetXY(20, yPos+4)
        pdf.CellFormat(257, 7, "Telah menyelesaikan Program Magang dengan baik dan memuaskan.", "", 0, "C", false, 0, "")

        // Tanda tangan
        tanggal := time.Now().Format("02 Januari 2006")
        pdf.SetFont("Helvetica", "", 10)
        pdf.SetTextColor(50, 50, 50)
        pdf.SetXY(190, 168)
        pdf.CellFormat(80, 6, "Muara Enim, "+tanggal, "", 0, "C", false, 0, "")
        pdf.SetXY(190, 175)
        pdf.CellFormat(80, 6, "Manager HRD", "", 0, "C", false, 0, "")

        pdf.SetDrawColor(0, 100, 0)
        pdf.SetLineWidth(0.5)
        pdf.Line(195, 193, 265, 193)
        pdf.SetFont("Helvetica", "B", 10)
        pdf.SetXY(190, 195)
        pdf.CellFormat(80, 6, "PT TanjungEnim Lestari", "", 0, "C", false, 0, "")

        // QR Code area (placeholder)
        pdf.SetDrawColor(200, 200, 200)
        pdf.SetLineWidth(0.3)
        pdf.Rect(20, 168, 22, 22, "D")
        pdf.SetFont("Helvetica", "", 7)
        pdf.SetTextColor(150, 150, 150)
        pdf.SetXY(20, 191)
        pdf.CellFormat(22, 4, "Scan verify", "", 0, "C", false, 0, "")

        return pdf.OutputFileAndClose(savePath)
}

func sanitizeFilename(name string) string {
        result := ""
        for _, r := range name {
                if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
                        result += string(r)
                } else {
                        result += "_"
                }
        }
        return result
}
