package service

import (
        "bytes"
        "encoding/base64"
        "encoding/json"
        "fmt"
        "net/http"
        "os"
        "path/filepath"

        "github.com/jung-kurt/gofpdf"
)

// DokumenLampiran — dokumen yang akan dilampirkan ke email backup
type DokumenLampiran struct {
        Jenis    string
        NamaFile string
        PathFile string // absolute path di filesystem
}

// KejuruanItem — satu komponen kejuruan untuk PDF penilaian
type KejuruanItem struct {
        Komponen string
        Nilai    float64
}

// PenilaianLampiranData — data penilaian lengkap untuk generate PDF Lembar Penilaian backup
type PenilaianLampiranData struct {
        NamaLengkap   string
        NomorInduk    string
        AsalInstitusi string
        Jurusan       string
        KelasSemester string
        Divisi        string
        Pembimbing    string
        Periode       string
        ManagerNama   string
        Catatan       string
        DinilaiAt     string

        NilaiMotivasi      float64
        NilaiInisiatif     float64
        NilaiDisiplinWaktu float64
        NilaiKerajinan     float64
        NilaiKreativitas   float64
        NilaiTanggungJawab float64
        NilaiKerjasama     float64
        NilaiAdaptasi      float64
        NilaiKehadiran     float64

        NilaiK3Safety    float64
        NilaiK3Metode    float64
        NilaiK3Manajemen float64
        NilaiK3Volume    float64

        NilaiPrsProses float64
        NilaiPrsTeori  float64
        NilaiPrsJudul  float64
        NilaiPrsData   float64

        NilaiAkhir float64
        Kejuruan   []KejuruanItem
}

// generateLembarPenilaianPDF — buat PDF "Lembar Penilaian Magang" menggunakan gofpdf
func generateLembarPenilaianPDF(p *PenilaianLampiranData) ([]byte, error) {
        pdf := gofpdf.New("P", "mm", "A4", "")
        pdf.SetMargins(15, 15, 15)
        pdf.AddPage()

        grade := func(n float64) string {
                switch {
                case n >= 85:
                        return "A — Sangat Baik"
                case n >= 70:
                        return "B — Baik"
                case n >= 55:
                        return "C — Cukup"
                default:
                        return "D — Kurang"
                }
        }
        fmtN := func(n float64) string {
                if n == 0 {
                        return "-"
                }
                return fmt.Sprintf("%.1f", n)
        }

        // ── Header ──
        pdf.SetFillColor(0, 100, 0)
        pdf.Rect(15, 15, 180, 10, "F")
        pdf.SetTextColor(255, 255, 255)
        pdf.SetFont("Helvetica", "B", 11)
        pdf.SetXY(15, 16)
        pdf.CellFormat(180, 8, "PT TANJUNGENIM LESTARI PULP AND PAPER", "", 0, "C", false, 0, "")

        pdf.SetTextColor(0, 0, 0)
        pdf.SetFont("Helvetica", "B", 13)
        pdf.SetXY(15, 28)
        pdf.CellFormat(180, 8, "DAFTAR PENILAIAN PESERTA MAGANG", "", 1, "C", false, 0, "")
        pdf.SetDrawColor(0, 100, 0)
        pdf.SetLineWidth(0.5)
        pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
        pdf.Ln(3)

        // ── Info Peserta ──
        infoRows := [][2]string{
                {"Nama Mahasiswa", p.NamaLengkap},
                {"NIM / NRP", p.NomorInduk},
                {"Program Studi", p.Jurusan},
                {"Kelas / Semester", p.KelasSemester},
                {"Asal Perguruan Tinggi", p.AsalInstitusi},
        }
        infoRight := [][2]string{
                {"Unit Kerja / Divisi", p.Divisi},
                {"Nama Pembimbing", p.Pembimbing},
                {"Periode Magang", p.Periode},
                {"Manager Dept.", func() string { if p.ManagerNama != "" { return p.ManagerNama }; return "—" }()},
        }
        pdf.SetFont("Helvetica", "", 9)
        yStart := pdf.GetY()
        for i, r := range infoRows {
                yRow := yStart + float64(i)*6
                pdf.SetXY(15, yRow)
                pdf.SetFont("Helvetica", "B", 9)
                pdf.CellFormat(42, 5.5, r[0], "", 0, "L", false, 0, "")
                pdf.SetFont("Helvetica", "", 9)
                pdf.CellFormat(4, 5.5, ":", "", 0, "C", false, 0, "")
                pdf.CellFormat(44, 5.5, r[1], "", 0, "L", false, 0, "")
        }
        for i, r := range infoRight {
                yRow := yStart + float64(i)*6
                pdf.SetXY(107, yRow)
                pdf.SetFont("Helvetica", "B", 9)
                pdf.CellFormat(42, 5.5, r[0], "", 0, "L", false, 0, "")
                pdf.SetFont("Helvetica", "", 9)
                pdf.CellFormat(4, 5.5, ":", "", 0, "C", false, 0, "")
                pdf.CellFormat(42, 5.5, r[1], "", 0, "L", false, 0, "")
        }
        pdf.SetY(yStart + float64(len(infoRows))*6 + 4)
        pdf.SetDrawColor(200, 200, 200)
        pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
        pdf.Ln(3)

        // ── Fungsi cetak tabel nilai ──
        tableHeader := func(title string) {
                pdf.SetFillColor(34, 85, 34)
                pdf.SetTextColor(255, 255, 255)
                pdf.SetFont("Helvetica", "B", 9)
                pdf.CellFormat(10, 6, "No", "1", 0, "C", true, 0, "")
                pdf.CellFormat(100, 6, title, "1", 0, "L", true, 0, "")
                pdf.CellFormat(30, 6, "Nilai (Angka)", "1", 0, "C", true, 0, "")
                pdf.CellFormat(40, 6, "Keterangan", "1", 1, "C", true, 0, "")
                pdf.SetTextColor(0, 0, 0)
        }
        tableRow := func(no int, label string, nilai float64) {
                bg := no%2 == 0
                if bg {
                        pdf.SetFillColor(245, 250, 245)
                } else {
                        pdf.SetFillColor(255, 255, 255)
                }
                pdf.SetFont("Helvetica", "", 9)
                pdf.CellFormat(10, 6, fmt.Sprintf("%d", no), "1", 0, "C", bg, 0, "")
                pdf.CellFormat(100, 6, label, "1", 0, "L", bg, 0, "")
                pdf.SetFont("Helvetica", "B", 9)
                pdf.CellFormat(30, 6, fmtN(nilai), "1", 0, "C", bg, 0, "")
                pdf.SetFont("Helvetica", "", 9)
                pdf.CellFormat(40, 6, grade(nilai), "1", 1, "L", bg, 0, "")
        }

        // ── I. Kepribadian / Etos Kerja ──
        pdf.SetFillColor(220, 240, 220)
        pdf.SetTextColor(0, 60, 0)
        pdf.SetFont("Helvetica", "B", 9)
        pdf.CellFormat(180, 6, "  I.  KEPRIBADIAN / ETOS KERJA", "1", 1, "L", true, 0, "")
        pdf.SetTextColor(0, 0, 0)
        tableHeader("Unsur yang Dinilai")
        kepribadian := [][2]interface{}{
                {"Motivasi", p.NilaiMotivasi},
                {"Inisiatif", p.NilaiInisiatif},
                {"Disiplin Waktu", p.NilaiDisiplinWaktu},
                {"Kerajinan", p.NilaiKerajinan},
                {"Kreativitas", p.NilaiKreativitas},
                {"Tanggung Jawab", p.NilaiTanggungJawab},
                {"Kerjasama", p.NilaiKerjasama},
                {"Adaptasi dengan Lingkungan Kerja", p.NilaiAdaptasi},
                {"Kehadiran", p.NilaiKehadiran},
        }
        for i, row := range kepribadian {
                tableRow(i+1, row[0].(string), row[1].(float64))
        }
        pdf.Ln(2)

        // ── II. Kejuruan ──
        if len(p.Kejuruan) > 0 {
                pdf.SetFillColor(220, 240, 220)
                pdf.SetTextColor(0, 60, 0)
                pdf.SetFont("Helvetica", "B", 9)
                pdf.CellFormat(180, 6, "  II.  KEAHLIAN KEJURUAN", "1", 1, "L", true, 0, "")
                pdf.SetTextColor(0, 0, 0)
                tableHeader("Unsur yang Dinilai")
                for i, k := range p.Kejuruan {
                        tableRow(i+1, k.Komponen, k.Nilai)
                }
                pdf.Ln(2)
        }

        // ── III. K3 ──
        pdf.SetFillColor(220, 240, 220)
        pdf.SetTextColor(0, 60, 0)
        pdf.SetFont("Helvetica", "B", 9)
        pdf.CellFormat(180, 6, "  III.  KESELAMATAN & KESEHATAN KERJA (K3)", "1", 1, "L", true, 0, "")
        pdf.SetTextColor(0, 0, 0)
        tableHeader("Unsur yang Dinilai")
        k3Rows := [][2]interface{}{
                {"Safety (Keselamatan Kerja)", p.NilaiK3Safety},
                {"Metode Kerja", p.NilaiK3Metode},
                {"Manajemen K3", p.NilaiK3Manajemen},
                {"Volume Kerja", p.NilaiK3Volume},
        }
        for i, row := range k3Rows {
                tableRow(i+1, row[0].(string), row[1].(float64))
        }
        pdf.Ln(2)

        // ── IV. Presentasi ──
        pdf.SetFillColor(220, 240, 220)
        pdf.SetTextColor(0, 60, 0)
        pdf.SetFont("Helvetica", "B", 9)
        pdf.CellFormat(180, 6, "  IV.  PRESENTASI / LAPORAN AKHIR", "1", 1, "L", true, 0, "")
        pdf.SetTextColor(0, 0, 0)
        tableHeader("Unsur yang Dinilai")
        prsRows := [][2]interface{}{
                {"Proses Presentasi", p.NilaiPrsProses},
                {"Penguasaan Teori", p.NilaiPrsTeori},
                {"Judul / Tema", p.NilaiPrsJudul},
                {"Data & Analisis", p.NilaiPrsData},
        }
        for i, row := range prsRows {
                tableRow(i+1, row[0].(string), row[1].(float64))
        }
        pdf.Ln(4)

        // ── Nilai Akhir ──
        pdf.SetFillColor(0, 100, 0)
        pdf.SetTextColor(255, 255, 255)
        pdf.SetFont("Helvetica", "B", 11)
        pdf.CellFormat(130, 9, "NILAI AKHIR", "1", 0, "L", true, 0, "")
        pdf.SetFont("Helvetica", "B", 13)
        pdf.CellFormat(50, 9, fmt.Sprintf("%.1f  —  %s", p.NilaiAkhir, grade(p.NilaiAkhir)), "1", 1, "C", true, 0, "")
        pdf.SetTextColor(0, 0, 0)
        pdf.Ln(3)

        // ── Catatan & Tanggal ──
        if p.Catatan != "" {
                pdf.SetFont("Helvetica", "I", 9)
                pdf.SetTextColor(80, 80, 80)
                pdf.MultiCell(180, 5, "Catatan: "+p.Catatan, "", "L", false)
                pdf.Ln(2)
        }
        if p.DinilaiAt != "" {
                pdf.SetTextColor(120, 120, 120)
                pdf.SetFont("Helvetica", "", 8)
                pdf.CellFormat(180, 5, "Dinilai pada: "+p.DinilaiAt, "", 1, "R", false, 0, "")
        }

        var buf bytes.Buffer
        if err := pdf.Output(&buf); err != nil {
                return nil, fmt.Errorf("gagal generate PDF penilaian: %w", err)
        }
        return buf.Bytes(), nil
}

type EmailService struct{}

func NewEmailService() *EmailService {
        return &EmailService{}
}

type resendAttachment struct {
        Filename string `json:"filename"`
        Content  string `json:"content"`
}

type resendEmailRequest struct {
        From        string             `json:"from"`
        To          []string           `json:"to"`
        Subject     string             `json:"subject"`
        Html        string             `json:"html"`
        Attachments []resendAttachment `json:"attachments,omitempty"`
}

func (s *EmailService) kirimViaResend(to, subject, html string, attachments []resendAttachment) error {
        apiKey := os.Getenv("RESEND_API_KEY")
        if apiKey == "" {
                return fmt.Errorf("RESEND_API_KEY tidak dikonfigurasi")
        }

        from := os.Getenv("RESEND_FROM_EMAIL")
        if from == "" {
                from = "e-Magang TELPP <onboarding@resend.dev>"
        }

        payload := resendEmailRequest{
                From:        from,
                To:          []string{to},
                Subject:     subject,
                Html:        html,
                Attachments: attachments,
        }

        body, _ := json.Marshal(payload)
        req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(body))
        if err != nil {
                return fmt.Errorf("gagal membuat request: %w", err)
        }
        req.Header.Set("Authorization", "Bearer "+apiKey)
        req.Header.Set("Content-Type", "application/json")

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
                return fmt.Errorf("gagal mengirim email: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode >= 300 {
                return fmt.Errorf("Resend API error: status %d", resp.StatusCode)
        }
        return nil
}

func frontendURL() string {
        u := os.Getenv("FRONTEND_URL")
        if u == "" {
                u = "http://localhost:5000"
        }
        return u
}

func buatLampiranSuratBalasan(pathFile string) []resendAttachment {
        if pathFile == "" {
                return nil
        }
        data, err := os.ReadFile(pathFile)
        if err != nil {
                return nil
        }
        return []resendAttachment{
                {
                        Filename: "Surat_Balasan_Magang_TELPP.pdf",
                        Content:  base64.StdEncoding.EncodeToString(data),
                },
        }
}

// buatLampiranSuratUpload — buat lampiran dari relative path di upload dir
func buatLampiranSuratUpload(uploadDir, relPath string) []resendAttachment {
        if relPath == "" {
                return nil
        }
        fullPath := uploadDir + "/" + relPath
        data, err := os.ReadFile(fullPath)
        if err != nil {
                return nil
        }
        ext := ".pdf"
        if len(relPath) > 4 {
                e := relPath[len(relPath)-4:]
                if e == ".jpg" || e == "jpeg" || e == ".png" {
                        ext = e
                }
        }
        return []resendAttachment{
                {
                        Filename: "Surat_Perpanjangan_Magang_TELPP" + ext,
                        Content:  base64.StdEncoding.EncodeToString(data),
                },
        }
}

// KirimPerpanjangan — email notifikasi perpanjangan magang (HRD langsung atau approve pengajuan peserta)
// suratRelPath adalah relative path di upload dir (misal "surat_perpanjangan/uuid.pdf"), boleh kosong
func (s *EmailService) KirimPerpanjangan(toEmail, namaLengkap string, durasiHari int, tanggalSelesaiBaru, alasan, uploadDir, suratRelPath string) error {
        var suratNote string
        if suratRelPath != "" {
                suratNote = `
  <tr>
    <td style="padding:16px 40px 0;">
      <p style="margin:0;font-size:13px;color:#6b7280;line-height:1.6;">
        Surat perpanjangan resmi dari PT TELPP terlampir pada email ini.
      </p>
    </td>
  </tr>`
        }

        var alasanSection string
        if alasan != "" {
                alasanSection = fmt.Sprintf(`
  <!-- CATATAN -->
  <tr>
    <td style="padding:20px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0"
             style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:8px;">
        <tr>
          <td style="padding:16px 20px;">
            <p style="margin:0 0 6px;font-size:11px;font-weight:700;color:#166534;text-transform:uppercase;letter-spacing:0.1em;">Keterangan</p>
            <p style="margin:0;font-size:13px;color:#374151;line-height:1.7;">%s</p>
          </td>
        </tr>
      </table>
    </td>
  </tr>`, alasan)
        }

        html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Perpanjangan Magang PT TELPP</title>
</head>
<body style="margin:0;padding:0;background:#f4f5f7;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f5f7;padding:40px 16px;">
<tr><td align="center">
<table width="520" cellpadding="0" cellspacing="0"
       style="max-width:520px;width:100%%;background:#ffffff;border-radius:12px;
              box-shadow:0 2px 12px rgba(0,0,0,.07);overflow:hidden;">

  <!-- TOP ACCENT BAR -->
  <tr><td style="background:#166534;height:4px;font-size:0;">&nbsp;</td></tr>

  <!-- HEADER -->
  <tr>
    <td style="padding:32px 40px 0;text-align:center;">
      <p style="margin:0;font-size:11px;font-weight:700;color:#166534;text-transform:uppercase;letter-spacing:0.12em;">
        PT TanjungEnim Lestari Pulp and Paper
      </p>
      <h1 style="margin:8px 0 0;font-size:22px;font-weight:700;color:#111827;letter-spacing:-0.01em;">
        Perpanjangan Magang Disetujui
      </h1>
    </td>
  </tr>

  <!-- DIVIDER -->
  <tr>
    <td style="padding:24px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0">
        <tr><td style="border-top:1px solid #f3f4f6;font-size:0;">&nbsp;</td></tr>
      </table>
    </td>
  </tr>

  <!-- BODY -->
  <tr>
    <td style="padding:0 40px;">
      <p style="margin:0 0 14px;font-size:14px;color:#374151;line-height:1.8;">
        Yth. <strong>%s</strong>,
      </p>
      <p style="margin:0;font-size:14px;color:#374151;line-height:1.8;">
        Kami informasikan bahwa HRD <strong>PT TanjungEnim Lestari Pulp and Paper</strong>
        telah menyetujui perpanjangan masa magang Anda.
      </p>
    </td>
  </tr>

  <!-- INFO BOX -->
  <tr>
    <td style="padding:20px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0"
             style="background:#f9fafb;border:1px solid #e5e7eb;border-radius:8px;">
        <tr>
          <td style="padding:20px 24px;">
            <p style="margin:0 0 14px;font-size:11px;font-weight:700;color:#9ca3af;
                       text-transform:uppercase;letter-spacing:0.1em;">Rincian Perpanjangan</p>
            <table width="100%%" cellpadding="0" cellspacing="0">
              <tr>
                <td style="padding-bottom:10px;border-bottom:1px solid #e5e7eb;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Durasi Perpanjangan</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%d hari</p>
                </td>
              </tr>
              <tr>
                <td style="padding-top:10px;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Tanggal Selesai Baru</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#166534;">%s</p>
                </td>
              </tr>
            </table>
          </td>
        </tr>
      </table>
    </td>
  </tr>

  %s
  %s

  <!-- FOOTER -->
  <tr>
    <td style="padding:32px 40px 28px;text-align:center;">
      <p style="margin:0;font-size:11px;color:#9ca3af;line-height:1.7;">
        PT TanjungEnim Lestari Pulp and Paper &bull; Muara Enim, Sumatera Selatan<br>
        Email otomatis &mdash; mohon tidak membalas.
      </p>
    </td>
  </tr>

</table>
</td></tr>
</table>
</body>
</html>`, namaLengkap, durasiHari, tanggalSelesaiBaru, alasanSection, suratNote)

        lampiran := buatLampiranSuratUpload(uploadDir, suratRelPath)
        return s.kirimViaResend(toEmail, "Perpanjangan Magang Anda Telah Disetujui — PT TELPP", html, lampiran)
}

// KirimKredensial — email pemberitahuan diterima beserta kredensial login
func (s *EmailService) KirimKredensial(toEmail, namaLengkap, password, catatan, suratBalasanPath string) error {
        loginURL := frontendURL() + "/login"

        var credentialRows string
        if password != "" {
                // Akun baru — tampilkan email + password sementara
                credentialRows = fmt.Sprintf(`
              <tr>
                <td style="padding-bottom:12px;border-bottom:1px solid #e5e7eb;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Email</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
                </td>
              </tr>
              <tr>
                <td style="padding-top:12px;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Password Sementara</p>
                  <p style="margin:0;font-size:24px;font-weight:800;color:#166534;letter-spacing:0.14em;font-family:'Courier New',monospace;">%s</p>
                </td>
              </tr>`, toEmail, password)
        } else {
                // Akun lama — hanya tampilkan email, password tidak berubah
                credentialRows = fmt.Sprintf(`
              <tr>
                <td>
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Email</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
                  <p style="margin:8px 0 0;font-size:13px;color:#6b7280;">Gunakan password akun Anda yang sudah ada.</p>
                </td>
              </tr>`, toEmail)
        }

        var warningSection string
        if password != "" {
                warningSection = `
  <!-- WARNING -->
  <tr>
    <td style="padding:20px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0"
             style="background:#fffbeb;border:1px solid #fde68a;border-radius:6px;">
        <tr>
          <td style="padding:12px 16px;font-size:13px;color:#92400e;line-height:1.6;">
            <strong>Segera ubah password ini</strong> setelah pertama kali masuk melalui menu <em>Pengaturan Akun</em>.
          </td>
        </tr>
      </table>
    </td>
  </tr>`
        }

        var catatanSection string
        if catatan != "" {
                catatanSection = fmt.Sprintf(`
  <!-- CATATAN HRD -->
  <tr>
    <td style="padding:20px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0"
             style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:8px;">
        <tr>
          <td style="padding:16px 20px;">
            <p style="margin:0 0 8px;font-size:11px;font-weight:700;color:#15803d;text-transform:uppercase;letter-spacing:0.1em;">Catatan dari HRD</p>
            <p style="margin:0;font-size:13px;color:#374151;line-height:1.7;">%s</p>
          </td>
        </tr>
      </table>
    </td>
  </tr>`, catatan)
        }

        var suratNote string
        if suratBalasanPath != "" {
                suratNote = `
  <tr>
    <td style="padding:16px 40px 0;">
      <p style="margin:0;font-size:13px;color:#6b7280;line-height:1.6;">
        Surat balasan resmi magang terlampir pada email ini.
      </p>
    </td>
  </tr>`
        }

        html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Akun e-Magang PT TELPP</title>
</head>
<body style="margin:0;padding:0;background:#f4f5f7;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f5f7;padding:40px 16px;">
<tr><td align="center">
<table width="520" cellpadding="0" cellspacing="0"
       style="max-width:520px;width:100%%;background:#ffffff;border-radius:12px;
              box-shadow:0 2px 12px rgba(0,0,0,.07);overflow:hidden;">

  <!-- TOP ACCENT BAR -->
  <tr><td style="background:#166534;height:4px;font-size:0;">&nbsp;</td></tr>

  <!-- HEADER -->
  <tr>
    <td style="padding:32px 40px 0;text-align:center;">
      <p style="margin:0;font-size:11px;font-weight:700;color:#166534;text-transform:uppercase;letter-spacing:0.12em;">
        PT TanjungEnim Lestari Pulp and Paper
      </p>
      <h1 style="margin:8px 0 0;font-size:22px;font-weight:700;color:#111827;letter-spacing:-0.01em;">
        Akun Anda Sudah Siap
      </h1>
      <p style="margin:10px 0 0;font-size:14px;color:#6b7280;line-height:1.6;">
        Halo <strong style="color:#111827;">%s</strong>, akun e-Magang Anda telah dibuat oleh Tim HRD.
      </p>
    </td>
  </tr>

  <!-- DIVIDER -->
  <tr>
    <td style="padding:24px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0">
        <tr><td style="border-top:1px solid #f3f4f6;font-size:0;">&nbsp;</td></tr>
      </table>
    </td>
  </tr>

  <!-- CREDENTIAL BOX -->
  <tr>
    <td style="padding:0 40px;">
      <table width="100%%" cellpadding="0" cellspacing="0"
             style="background:#f9fafb;border:1px solid #e5e7eb;border-radius:8px;">
        <tr>
          <td style="padding:20px 24px;">
            <p style="margin:0 0 14px;font-size:11px;font-weight:700;color:#9ca3af;
                       text-transform:uppercase;letter-spacing:0.1em;">Detail Login</p>
            <table width="100%%" cellpadding="0" cellspacing="0">
              %s
            </table>
          </td>
        </tr>
      </table>
    </td>
  </tr>

  <!-- CTA -->
  <tr>
    <td style="padding:24px 40px 0;text-align:center;">
      <a href="%s"
         style="display:inline-block;background:#166534;color:#ffffff;text-decoration:none;
                border-radius:8px;padding:13px 36px;font-size:14px;font-weight:600;
                letter-spacing:0.02em;">
        Masuk ke e-Magang &rarr;
      </a>
    </td>
  </tr>

  %s
  %s
  %s

  <!-- FOOTER -->
  <tr>
    <td style="padding:28px 40px;text-align:center;">
      <p style="margin:0;font-size:11px;color:#9ca3af;line-height:1.7;">
        PT TanjungEnim Lestari Pulp and Paper &bull; Muara Enim, Sumatera Selatan<br>
        Email otomatis &mdash; mohon tidak membalas.
      </p>
    </td>
  </tr>

</table>
</td></tr>
</table>
</body>
</html>`, namaLengkap, credentialRows, loginURL, warningSection, catatanSection, suratNote)

        return s.kirimViaResend(toEmail, "Akun e-Magang PT TELPP Anda Sudah Siap", html, buatLampiranSuratBalasan(suratBalasanPath))
}

// KirimSertifikat — email pengiriman sertifikat magang ke peserta (formal, PDF terlampir)
func (s *EmailService) KirimSertifikat(toEmail, namaLengkap, divisi, pembimbing, periode, nilai, sertifikatPath string) error {
        sertifikatURL := frontendURL() + "/dashboard?tab=sertifikat"

        html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Sertifikat Magang PT TELPP</title>
</head>
<body style="margin:0;padding:0;background:#f4f5f7;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f5f7;padding:40px 16px;">
<tr><td align="center">
<table width="520" cellpadding="0" cellspacing="0"
       style="max-width:520px;width:100%%;background:#ffffff;border-radius:12px;
              box-shadow:0 2px 12px rgba(0,0,0,.07);overflow:hidden;">

  <!-- TOP ACCENT BAR -->
  <tr><td style="background:#166534;height:4px;font-size:0;">&nbsp;</td></tr>

  <!-- HEADER -->
  <tr>
    <td style="padding:32px 40px 0;text-align:center;">
      <p style="margin:0;font-size:11px;font-weight:700;color:#166534;text-transform:uppercase;letter-spacing:0.12em;">
        PT TanjungEnim Lestari Pulp and Paper
      </p>
      <h1 style="margin:8px 0 0;font-size:22px;font-weight:700;color:#111827;letter-spacing:-0.01em;">
        Sertifikat Magang Telah Diterbitkan
      </h1>
      <p style="margin:10px 0 0;font-size:14px;color:#6b7280;line-height:1.6;">
        Kepada Yth. <strong style="color:#111827;">%s</strong>,<br>
        dengan bangga kami sampaikan bahwa sertifikat magang Anda<br>
        telah resmi diterbitkan oleh PT TanjungEnim Lestari Pulp and Paper.
      </p>
    </td>
  </tr>

  <!-- DIVIDER -->
  <tr>
    <td style="padding:24px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0">
        <tr><td style="border-top:1px solid #f3f4f6;font-size:0;">&nbsp;</td></tr>
      </table>
    </td>
  </tr>

  <!-- INFO BOX -->
  <tr>
    <td style="padding:0 40px;">
      <table width="100%%" cellpadding="0" cellspacing="0"
             style="background:#f9fafb;border:1px solid #e5e7eb;border-radius:8px;">
        <tr>
          <td style="padding:20px 24px;">
            <p style="margin:0 0 14px;font-size:11px;font-weight:700;color:#9ca3af;
                       text-transform:uppercase;letter-spacing:0.1em;">Ringkasan Magang</p>
            <table width="100%%" cellpadding="0" cellspacing="0">
              <tr>
                <td style="padding-bottom:10px;border-bottom:1px solid #e5e7eb;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Nama Peserta</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
                </td>
              </tr>
              <tr>
                <td style="padding:10px 0;border-bottom:1px solid #e5e7eb;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Unit Kerja / Divisi</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
                </td>
              </tr>
              <tr>
                <td style="padding:10px 0;border-bottom:1px solid #e5e7eb;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Pembimbing Lapangan</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
                </td>
              </tr>
              <tr>
                <td style="padding:10px 0;border-bottom:1px solid #e5e7eb;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Periode Magang</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
                </td>
              </tr>
              <tr>
                <td style="padding-top:10px;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Nilai Akhir</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
                </td>
              </tr>
            </table>
          </td>
        </tr>
      </table>
    </td>
  </tr>

  <!-- NOTE LAMPIRAN -->
  <tr>
    <td style="padding:16px 40px 0;">
      <p style="margin:0;font-size:13px;color:#6b7280;line-height:1.6;">
        Sertifikat resmi magang terlampir pada email ini dalam format PDF.
        Anda juga dapat mengunduhnya kapan saja melalui portal e-Magang.
      </p>
    </td>
  </tr>

  <!-- CTA -->
  <tr>
    <td style="padding:24px 40px 0;text-align:center;">
      <a href="%s"
         style="display:inline-block;background:#166534;color:#ffffff;text-decoration:none;
                border-radius:8px;padding:13px 36px;font-size:14px;font-weight:600;
                letter-spacing:0.02em;">
        Unduh Sertifikat &rarr;
      </a>
    </td>
  </tr>

  <!-- FOOTER -->
  <tr>
    <td style="padding:32px 40px 28px;text-align:center;">
      <p style="margin:0;font-size:11px;color:#9ca3af;line-height:1.7;">
        PT TanjungEnim Lestari Pulp and Paper &bull; Muara Enim, Sumatera Selatan<br>
        Email otomatis &mdash; mohon tidak membalas.
      </p>
    </td>
  </tr>

</table>
</td></tr>
</table>
</body>
</html>`, namaLengkap, namaLengkap, divisi, pembimbing, periode, nilai, sertifikatURL)

        attachments := buatLampiranSertifikat(sertifikatPath)
        return s.kirimViaResend(toEmail, "Sertifikat Magang PT TELPP — "+namaLengkap, html, attachments)
}

func buatLampiranSertifikat(pathFile string) []resendAttachment {
        if pathFile == "" {
                return nil
        }
        data, err := os.ReadFile(pathFile)
        if err != nil {
                return nil
        }
        return []resendAttachment{
                {
                        Filename: "Sertifikat_Magang_PT_TELPP.pdf",
                        Content:  base64.StdEncoding.EncodeToString(data),
                },
        }
}

// KirimDitolak — email pemberitahuan hasil seleksi tidak diterima
func (s *EmailService) KirimDitolak(toEmail, namaLengkap, catatan, suratBalasanPath string) error {
        var catatanSection string
        if catatan != "" {
                catatanSection = fmt.Sprintf(`
  <!-- CATATAN HRD -->
  <tr>
    <td style="padding:20px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0"
             style="background:#f9fafb;border:1px solid #e5e7eb;border-radius:8px;">
        <tr>
          <td style="padding:16px 20px;">
            <p style="margin:0 0 8px;font-size:11px;font-weight:700;color:#9ca3af;text-transform:uppercase;letter-spacing:0.1em;">Keterangan</p>
            <p style="margin:0;font-size:13px;color:#374151;line-height:1.7;">%s</p>
          </td>
        </tr>
      </table>
    </td>
  </tr>`, catatan)
        }

        var suratNote string
        if suratBalasanPath != "" {
                suratNote = `
  <tr>
    <td style="padding:16px 40px 0;">
      <p style="margin:0;font-size:13px;color:#6b7280;line-height:1.6;">
        Surat balasan resmi terlampir pada email ini.
      </p>
    </td>
  </tr>`
        }

        html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Informasi Hasil Seleksi Magang PT TELPP</title>
</head>
<body style="margin:0;padding:0;background:#f4f5f7;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f5f7;padding:40px 16px;">
<tr><td align="center">
<table width="520" cellpadding="0" cellspacing="0"
       style="max-width:520px;width:100%%;background:#ffffff;border-radius:12px;
              box-shadow:0 2px 12px rgba(0,0,0,.07);overflow:hidden;">

  <!-- TOP ACCENT BAR -->
  <tr><td style="background:#166534;height:4px;font-size:0;">&nbsp;</td></tr>

  <!-- HEADER -->
  <tr>
    <td style="padding:32px 40px 0;text-align:center;">
      <p style="margin:0;font-size:11px;font-weight:700;color:#166534;text-transform:uppercase;letter-spacing:0.12em;">
        PT TanjungEnim Lestari Pulp and Paper
      </p>
      <h1 style="margin:8px 0 0;font-size:22px;font-weight:700;color:#111827;letter-spacing:-0.01em;">
        Informasi Hasil Seleksi Magang
      </h1>
    </td>
  </tr>

  <!-- DIVIDER -->
  <tr>
    <td style="padding:24px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0">
        <tr><td style="border-top:1px solid #f3f4f6;font-size:0;">&nbsp;</td></tr>
      </table>
    </td>
  </tr>

  <!-- BODY -->
  <tr>
    <td style="padding:0 40px;">
      <p style="margin:0 0 14px;font-size:14px;color:#374151;line-height:1.8;">
        Yth. <strong>%s</strong>,
      </p>
      <p style="margin:0 0 14px;font-size:14px;color:#374151;line-height:1.8;">
        Terima kasih atas minat Anda untuk bergabung dalam program magang di
        <strong>PT TanjungEnim Lestari Pulp and Paper</strong>.
      </p>
      <p style="margin:0;font-size:14px;color:#374151;line-height:1.8;">
        Setelah melalui proses seleksi berkas, dengan hormat kami sampaikan bahwa pada periode ini
        pengajuan magang Anda belum dapat kami terima. Kami mengundang Anda untuk kembali mendaftar
        pada periode pendaftaran berikutnya.
      </p>
    </td>
  </tr>

  %s
  %s

  <!-- FOOTER -->
  <tr>
    <td style="padding:32px 40px 28px;text-align:center;">
      <p style="margin:0;font-size:11px;color:#9ca3af;line-height:1.7;">
        PT TanjungEnim Lestari Pulp and Paper &bull; Muara Enim, Sumatera Selatan<br>
        Email otomatis &mdash; mohon tidak membalas.
      </p>
    </td>
  </tr>

</table>
</td></tr>
</table>
</body>
</html>`, namaLengkap, catatanSection, suratNote)

        return s.kirimViaResend(toEmail, "Informasi Hasil Seleksi Magang PT TELPP", html, buatLampiranSuratBalasan(suratBalasanPath))
}

// labelDokumen — konversi jenis enum ke nama file yang human-readable dengan ekstensi asli
func labelDokumen(jenis, namaFile string) string {
        labels := map[string]string{
                "proposal_magang": "Proposal_Magang",
                "surat_pengantar": "Surat_Pengantar",
                "ktp":             "KTP",
                "ktm":             "KTM",
                "pasfoto":         "Pasfoto",
                "bpjs_kis":        "BPJS_KIS",
                "surat_balasan":   "Surat_Balasan",
                "laporan_magang":  "Laporan_Magang",
                "sertifikat":      "Sertifikat_Magang",
        }
        label := jenis
        if l, ok := labels[jenis]; ok {
                label = l
        }
        ext := filepath.Ext(namaFile)
        if ext == "" {
                ext = ".pdf"
        }
        return label + ext
}

// KirimBackupHapusAkun — kirim email backup dokumen sebelum akun peserta dihapus
// sertifikatRelPath boleh "" jika belum ada; penilaian boleh nil jika belum ada penilaian
func (s *EmailService) KirimBackupHapusAkun(toEmail, namaLengkap, divisi, pembimbing, periode, nilai, uploadDir, sertifikatRelPath string, dokumen []DokumenLampiran, penilaian *PenilaianLampiranData) error {
        var infoNilai string
        if nilai != "" {
                infoNilai = fmt.Sprintf(`
  <tr>
    <td style="padding:10px 0;border-bottom:1px solid #e5e7eb;">
      <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Nilai Akhir</p>
      <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
    </td>
  </tr>`, nilai)
        }

        // Bangun semua lampiran: PDF penilaian (generate langsung) + sertifikat PDF
        var lampiran []resendAttachment

        // 1. Generate PDF Lembar Penilaian jika data penilaian tersedia
        if penilaian != nil {
                pdfBytes, pdfErr := generateLembarPenilaianPDF(penilaian)
                if pdfErr == nil {
                        lampiran = append(lampiran, resendAttachment{
                                Filename: fmt.Sprintf("Lembar_Penilaian_%s.pdf", penilaian.NamaLengkap),
                                Content:  base64.StdEncoding.EncodeToString(pdfBytes),
                        })
                }
        }

        // 2. Lampirkan sertifikat PDF dari disk
        // sertifikat_path di DB sudah merupakan path final yang disimpan oleh sertifikat_service
        // Pakai langsung seperti yang dilakukan DownloadSertifikat handler
        if sertifikatRelPath != "" {
                data, err := os.ReadFile(sertifikatRelPath)
                if err != nil && uploadDir != "" {
                        // Fallback: coba join dengan uploadDir jika path tidak ditemukan langsung
                        data, err = os.ReadFile(filepath.Join(uploadDir, sertifikatRelPath))
                }
                if err == nil {
                        lampiran = append(lampiran, resendAttachment{
                                Filename: "Sertifikat_Magang_TELPP.pdf",
                                Content:  base64.StdEncoding.EncodeToString(data),
                        })
                } else {
                        fmt.Printf("[backup-email] sertifikat tidak dapat dibaca: %s — %v\n", sertifikatRelPath, err)
                }
        }

        // 3. Dokumen tambahan dari tabel dokumen (laporan, sertifikat upload manual)
        for _, d := range dokumen {
                data, err := os.ReadFile(d.PathFile)
                if err != nil {
                        continue
                }
                lampiran = append(lampiran, resendAttachment{
                        Filename: labelDokumen(d.Jenis, d.NamaFile),
                        Content:  base64.StdEncoding.EncodeToString(data),
                })
        }

        var lampiranNote string
        if len(lampiran) > 0 {
                lampiranNote = fmt.Sprintf(`
  <tr>
    <td style="padding:16px 40px 0;">
      <p style="margin:0;font-size:13px;color:#6b7280;line-height:1.6;">
        <strong>%d file terlampir</strong> pada email ini (lembar penilaian &amp; sertifikat magang).
        Harap simpan file-file ini dengan baik karena akun Anda telah dihapus dari sistem.
      </p>
    </td>
  </tr>`, len(lampiran))
        }

        html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Backup Dokumen Magang PT TELPP</title>
</head>
<body style="margin:0;padding:0;background:#f4f5f7;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f5f7;padding:40px 16px;">
<tr><td align="center">
<table width="520" cellpadding="0" cellspacing="0"
       style="max-width:520px;width:100%%;background:#ffffff;border-radius:12px;
              box-shadow:0 2px 12px rgba(0,0,0,.07);overflow:hidden;">

  <!-- TOP ACCENT BAR -->
  <tr><td style="background:#166534;height:4px;font-size:0;">&nbsp;</td></tr>

  <!-- HEADER -->
  <tr>
    <td style="padding:32px 40px 0;text-align:center;">
      <p style="margin:0;font-size:11px;font-weight:700;color:#166534;text-transform:uppercase;letter-spacing:0.12em;">
        PT TanjungEnim Lestari Pulp and Paper
      </p>
      <h1 style="margin:8px 0 0;font-size:22px;font-weight:700;color:#111827;letter-spacing:-0.01em;">
        Backup Dokumen Magang
      </h1>
      <p style="margin:10px 0 0;font-size:13px;color:#6b7280;line-height:1.6;">
        Akun Anda di sistem e-Magang TELPP telah dihapus.<br>
        Email ini berisi salinan dokumen magang Anda sebagai arsip pribadi.
      </p>
    </td>
  </tr>

  <!-- DIVIDER -->
  <tr>
    <td style="padding:24px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0">
        <tr><td style="border-top:1px solid #f3f4f6;font-size:0;">&nbsp;</td></tr>
      </table>
    </td>
  </tr>

  <!-- BODY -->
  <tr>
    <td style="padding:0 40px;">
      <p style="margin:0;font-size:14px;color:#374151;line-height:1.8;">
        Yth. <strong>%s</strong>,<br>
        berikut adalah ringkasan data magang Anda di PT TanjungEnim Lestari Pulp and Paper.
      </p>
    </td>
  </tr>

  <!-- INFO BOX -->
  <tr>
    <td style="padding:20px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0"
             style="background:#f9fafb;border:1px solid #e5e7eb;border-radius:8px;">
        <tr>
          <td style="padding:20px 24px;">
            <p style="margin:0 0 14px;font-size:11px;font-weight:700;color:#9ca3af;
                       text-transform:uppercase;letter-spacing:0.1em;">Ringkasan Magang</p>
            <table width="100%%" cellpadding="0" cellspacing="0">
              <tr>
                <td style="padding-bottom:10px;border-bottom:1px solid #e5e7eb;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Nama Peserta</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
                </td>
              </tr>
              <tr>
                <td style="padding:10px 0;border-bottom:1px solid #e5e7eb;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Unit Kerja / Divisi</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
                </td>
              </tr>
              <tr>
                <td style="padding:10px 0;border-bottom:1px solid #e5e7eb;">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Pembimbing Lapangan</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
                </td>
              </tr>
              <tr>
                <td style="padding:10px 0%s">
                  <p style="margin:0 0 2px;font-size:11px;color:#9ca3af;">Periode Magang</p>
                  <p style="margin:0;font-size:14px;font-weight:600;color:#111827;">%s</p>
                </td>
              </tr>
              %s
            </table>
          </td>
        </tr>
      </table>
    </td>
  </tr>

  %s

  <!-- WARNING NOTE -->
  <tr>
    <td style="padding:20px 40px 0;">
      <table width="100%%" cellpadding="0" cellspacing="0"
             style="background:#fffbeb;border:1px solid #fde68a;border-radius:8px;">
        <tr>
          <td style="padding:14px 18px;">
            <p style="margin:0;font-size:12px;color:#92400e;line-height:1.7;">
              <strong>Perhatian:</strong> Akun dan seluruh data Anda telah dihapus dari sistem e-Magang TELPP.
              Simpan email ini sebagai arsip rekam jejak magang Anda.
            </p>
          </td>
        </tr>
      </table>
    </td>
  </tr>

  <!-- FOOTER -->
  <tr>
    <td style="padding:32px 40px 28px;text-align:center;">
      <p style="margin:0;font-size:11px;color:#9ca3af;line-height:1.7;">
        PT TanjungEnim Lestari Pulp and Paper &bull; Muara Enim, Sumatera Selatan<br>
        Email otomatis &mdash; mohon tidak membalas.
      </p>
    </td>
  </tr>

</table>
</td></tr>
</table>
</body>
</html>`, namaLengkap, namaLengkap, divisi, pembimbing,
                func() string {
                        if nilai != "" {
                                return ";border-bottom:1px solid #e5e7eb;"
                        }
                        return ";"
                }(),
                periode, infoNilai, lampiranNote)

        return s.kirimViaResend(toEmail, "Backup Dokumen Magang — PT TanjungEnim Lestari Pulp and Paper", html, lampiran)
}
