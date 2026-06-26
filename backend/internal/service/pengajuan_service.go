package service

import (
        "context"
        "crypto/rand"
        "errors"
        "fmt"
        "math/big"
        "os"

        "github.com/google/uuid"
        "github.com/jackc/pgx/v5"
        "golang.org/x/crypto/bcrypt"

        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
)

type PengajuanService struct {
        repo        *repository.PengajuanRepository
        notifSvc    *NotifikasiService
        userRepo    *repository.UserRepository
        emailSvc    *EmailService
        dokumenRepo *repository.DokumenRepository
}

func NewPengajuanService(
        repo *repository.PengajuanRepository,
        notifSvc *NotifikasiService,
        userRepo *repository.UserRepository,
        emailSvc *EmailService,
        dokumenRepo *repository.DokumenRepository,
) *PengajuanService {
        return &PengajuanService{repo: repo, notifSvc: notifSvc, userRepo: userRepo, emailSvc: emailSvc, dokumenRepo: dokumenRepo}
}

// generatePassword — "Mg" + 6 huruf acak + "1!" = 10 karakter
// Memenuhi regex: huruf (M,g,alpha), angka (1), special (!)
func generatePassword() string {
        const letters = "ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz"
        pwd := "Mg"
        for i := 0; i < 6; i++ {
                n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
                pwd += string(letters[n.Int64()])
        }
        return pwd + "1!"
}

// Submit — peserta yang sudah login mengajukan magang
func (s *PengajuanService) Submit(ctx context.Context, userID uuid.UUID, req1 models.PengajuanStep1Request, req2 models.PengajuanStep2Request) (*models.PengajuanMagang, error) {
        existing, err := s.repo.FindByUserID(ctx, userID)
        if err == nil && existing != nil {
                activeStatuses := []models.StatusPengajuan{
                        models.StatusDiajukan, models.StatusMenungguVerifikasi, models.StatusDiproses,
                }
                for _, st := range activeStatuses {
                        if existing.Status == st {
                                return nil, errors.New("anda sudah memiliki pengajuan yang sedang diproses")
                        }
                }
        }

        tgl, err := parseDate(req1.TanggalLahir)
        if err != nil {
                return nil, errors.New("format tanggal lahir tidak valid (gunakan YYYY-MM-DD)")
        }

        uid := userID
        p := &models.PengajuanMagang{
                UserID:         &uid,
                NamaLengkap:    req1.NamaLengkap,
                TempatLahir:    req1.TempatLahir,
                TanggalLahir:   tgl,
                JenisKelamin:   req1.JenisKelamin,
                Alamat:         req1.Alamat,
                NoHP:           req1.NoHP,
                Email:          req1.Email,
                KategoriMagang: req2.KategoriMagang,
                NomorInduk:     req2.NomorInduk,
                AsalInstitusi:  req2.AsalInstitusi,
                Jurusan:        req2.Jurusan,
                KelasSemester:  req2.KelasSemester,
        }

        if err := s.repo.Create(ctx, p); err != nil {
                return nil, err
        }

        // Notif ke semua HRD
        hrdList, _ := s.userRepo.FindHRDList(ctx)
        for _, h := range hrdList {
                s.notifSvc.KirimKeUser(ctx, h.ID, h.Role,
                        "Pengajuan Magang Baru",
                        p.NamaLengkap+" ("+string(p.KategoriMagang)+") dari "+p.AsalInstitusi,
                        string(models.NotifPengajuan), &p.ID)
        }

        return p, nil
}

// SubmitPublik — pengajuan dari form publik, tanpa login
func (s *PengajuanService) SubmitPublik(ctx context.Context, req1 models.PengajuanStep1Request, req2 models.PengajuanStep2Request) (*models.PengajuanMagang, error) {
        // Cek duplikat email dengan status aktif
        tgl, err := parseDate(req1.TanggalLahir)
        if err != nil {
                return nil, errors.New("format tanggal lahir tidak valid (gunakan YYYY-MM-DD)")
        }

        p := &models.PengajuanMagang{
                NamaLengkap:    req1.NamaLengkap,
                TempatLahir:    req1.TempatLahir,
                TanggalLahir:   tgl,
                JenisKelamin:   req1.JenisKelamin,
                Alamat:         req1.Alamat,
                NoHP:           req1.NoHP,
                Email:          req1.Email,
                KategoriMagang: req2.KategoriMagang,
                NomorInduk:     req2.NomorInduk,
                AsalInstitusi:  req2.AsalInstitusi,
                Jurusan:        req2.Jurusan,
                KelasSemester:  req2.KelasSemester,
        }

        if err := s.repo.CreatePublik(ctx, p); err != nil {
                return nil, err
        }

        // Notif ke semua HRD tentang pengajuan baru
        hrdList, _ := s.userRepo.FindHRDList(ctx)
        for _, h := range hrdList {
                s.notifSvc.KirimKeUser(ctx, h.ID, h.Role,
                        "Pengajuan Magang Baru (Publik)",
                        p.NamaLengkap+" ("+string(p.KategoriMagang)+") dari "+p.AsalInstitusi+" — via form publik",
                        string(models.NotifPengajuan), &p.ID)
        }

        return p, nil
}

// KirimAkun — HRD buat akun peserta dan kirim kredensial via email
func (s *PengajuanService) KirimAkun(ctx context.Context, pengajuanID uuid.UUID) error {
        p, err := s.repo.FindByID(ctx, pengajuanID)
        if err != nil {
                return errors.New("pengajuan tidak ditemukan")
        }
        if p.AkunTerkirimAt != nil {
                return errors.New("akun sudah pernah dikirim untuk pengajuan ini")
        }

        // Ambil path surat_balasan jika sudah diupload HRD
        suratBalasanPath := s.dokumenRepo.FindSuratBalasanPath(ctx, pengajuanID)

        // Cek apakah email sudah terdaftar
        existingUser, _ := s.userRepo.FindByEmail(ctx, p.Email)
        if existingUser != nil {
                // Tautkan saja ke akun yang sudah ada, lalu kirim email tanpa password baru
                if err := s.repo.SetAkunTerkirim(ctx, pengajuanID, existingUser.ID); err != nil {
                        return err
                }
                catatanHRD := ""
                if p.CatatanHRD != nil {
                        catatanHRD = *p.CatatanHRD
                }
                if err := s.emailSvc.KirimKredensial(p.Email, p.NamaLengkap, "", catatanHRD, suratBalasanPath); err != nil {
                        fmt.Printf("[WARN] Gagal kirim email ke %s: %v\n", p.Email, err)
                }
                return nil
        }

        // Buat user baru
        password := generatePassword()
        hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
        if err != nil {
                return fmt.Errorf("gagal hash password: %w", err)
        }

        newUser := &models.User{
                NamaLengkap:  p.NamaLengkap,
                Email:        p.Email,
                PasswordHash: string(hash),
                Role:         models.RolePeserta,
                IsActive:     true,
        }
        if err := s.userRepo.Create(ctx, newUser); err != nil {
                return fmt.Errorf("gagal membuat user: %w", err)
        }

        // Tautkan pengajuan ke user baru
        if err := s.repo.SetAkunTerkirim(ctx, pengajuanID, newUser.ID); err != nil {
                return fmt.Errorf("gagal update pengajuan: %w", err)
        }

        catatanHRD := ""
        if p.CatatanHRD != nil {
                catatanHRD = *p.CatatanHRD
        }
        // Kirim email dengan password + surat balasan (non-fatal)
        if err := s.emailSvc.KirimKredensial(p.Email, p.NamaLengkap, password, catatanHRD, suratBalasanPath); err != nil {
                fmt.Printf("[WARN] Gagal kirim email ke %s: %v\n", p.Email, err)
        }

        return nil
}

func (s *PengajuanService) GetMyPengajuan(ctx context.Context, userID uuid.UUID) (*models.PengajuanMagang, error) {
        p, err := s.repo.FindByUserID(ctx, userID)
        if err != nil {
                if errors.Is(err, pgx.ErrNoRows) {
                        return nil, nil
                }
                return nil, err
        }
        return p, nil
}

func (s *PengajuanService) GetAll(ctx context.Context, status, search string, page, limit int) ([]models.PengajuanMagang, int, error) {
        return s.repo.FindAll(ctx, status, search, page, limit)
}

func (s *PengajuanService) GetDetail(ctx context.Context, id uuid.UUID) (*models.PengajuanMagang, error) {
        return s.repo.FindByID(ctx, id)
}

func (s *PengajuanService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.StatusPengajuan, catatan string, changedBy uuid.UUID) error {
        p, err := s.repo.FindByID(ctx, id)
        if err != nil {
                return errors.New("pengajuan tidak ditemukan")
        }

        var catatanPtr *string
        if catatan != "" {
                catatanPtr = &catatan
        }

        if err := s.repo.UpdateStatus(ctx, id, status, catatanPtr, changedBy); err != nil {
                return err
        }

        // ── Otomatisasi saat DITERIMA ─────────────────────────────────────────
        if status == models.StatusDiterima {
                password := ""

                if p.AkunTerkirimAt == nil {
                        // Cek apakah email sudah terdaftar sebagai user
                        existingUser, _ := s.userRepo.FindByEmail(ctx, p.Email)
                        if existingUser != nil {
                                // Tautkan ke akun yang sudah ada (tanpa kirim password)
                                _ = s.repo.SetAkunTerkirim(ctx, id, existingUser.ID)
                        } else {
                                // Buat akun baru
                                password = generatePassword()
                                hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
                                if err == nil {
                                        newUser := &models.User{
                                                NamaLengkap:  p.NamaLengkap,
                                                Email:        p.Email,
                                                PasswordHash: string(hash),
                                                Role:         models.RolePeserta,
                                                IsActive:     true,
                                        }
                                        if err := s.userRepo.Create(ctx, newUser); err == nil {
                                                _ = s.repo.SetAkunTerkirim(ctx, id, newUser.ID)
                                                // Update p.UserID agar notif terkirim
                                                uid := newUser.ID
                                                p.UserID = &uid
                                        }
                                }
                        }
                }

                // Kirim email diterima + kredensial + catatan + surat balasan (non-fatal)
                suratBalasanPath := s.dokumenRepo.FindSuratBalasanPath(ctx, id)
                go func() {
                        if err := s.emailSvc.KirimKredensial(p.Email, p.NamaLengkap, password, catatan, suratBalasanPath); err != nil {
                                fmt.Printf("[WARN] Gagal kirim email diterima ke %s: %v\n", p.Email, err)
                        }
                }()
        }

        // ── Otomatisasi saat DITOLAK ─────────────────────────────────────────
        if status == models.StatusDitolak {
                suratBalasanPath := s.dokumenRepo.FindSuratBalasanPath(ctx, id)
                go func() {
                        if err := s.emailSvc.KirimDitolak(p.Email, p.NamaLengkap, catatan, suratBalasanPath); err != nil {
                                fmt.Printf("[WARN] Gagal kirim email ditolak ke %s: %v\n", p.Email, err)
                        }
                }()
        }

        // ── Notifikasi in-app ke peserta (jika sudah punya akun) ─────────────
        if p.UserID != nil && *p.UserID != uuid.Nil {
                type notifData struct{ judul, pesan, tipe string }
                notifMap := map[models.StatusPengajuan]notifData{
                        models.StatusMenungguVerifikasi: {
                                "Status Pengajuan Diperbarui",
                                "Pengajuan Anda sedang menunggu verifikasi berkas oleh HRD",
                                string(models.NotifPengajuan),
                        },
                        models.StatusDiproses: {
                                "Berkas Sedang Diproses",
                                "Berkas pengajuan Anda sedang diproses oleh tim HRD",
                                string(models.NotifPengajuan),
                        },
                        models.StatusDiterima: {
                                "Pengajuan Magang Diterima! 🎉",
                                "Selamat! Permohonan magang Anda di PT TanjungEnim Lestari telah diterima",
                                string(models.NotifPelaksanaan),
                        },
                        models.StatusDitolak: {
                                "Pengajuan Tidak Dapat Diterima",
                                "Permohonan magang Anda belum dapat kami terima saat ini. Silakan periksa catatan HRD",
                                string(models.NotifPengajuan),
                        },
                }

                if nd, ok := notifMap[status]; ok {
                        s.notifSvc.KirimKeUser(ctx, *p.UserID, models.RolePeserta, nd.judul, nd.pesan, nd.tipe, &id)
                }
        }

        return nil
}

func (s *PengajuanService) GetStatusHistory(ctx context.Context, pengajuanID uuid.UUID) ([]models.StatusHistory, error) {
        return s.repo.GetStatusHistory(ctx, pengajuanID)
}

// HapusPengajuan — hapus pengajuan + akun peserta (jika sudah dibuat) + file fisik.
//
// Alur:
//  1. Ambil data pengajuan (pastikan ada)
//  2. Ambil semua path file dokumen sebelum hapus DB
//  3. Hapus seluruh data DB dalam satu transaksi (HapusLengkap)
//  4. Hapus file fisik — non-fatal: kegagalan hanya dicatat di log
func (s *PengajuanService) HapusPengajuan(ctx context.Context, id uuid.UUID) error {
        // 1. Ambil pengajuan
        p, err := s.repo.FindByID(ctx, id)
        if err != nil {
                return errors.New("pengajuan tidak ditemukan")
        }

        // 2. Ambil semua path file dokumen sebelum data terhapus
        filePaths, err := s.dokumenRepo.GetPathsByPengajuanID(ctx, id)
        if err != nil {
                // Non-fatal: lanjutkan meskipun gagal query path
                fmt.Printf("[WARN] Gagal ambil path dokumen untuk pengajuan %s: %v\n", id, err)
                filePaths = nil
        }

        // Tentukan apakah akun peserta perlu dihapus:
        // user_id ada (bukan zero UUID) dan akun_terkirim_at tidak nil
        hapusUser := p.UserID != nil && *p.UserID != uuid.Nil && p.AkunTerkirimAt != nil

        // 3. Hapus seluruh data DB dalam satu transaksi
        var userIDVal uuid.UUID
        if p.UserID != nil {
                userIDVal = *p.UserID
        }
        if err := s.repo.HapusLengkap(ctx, id, userIDVal, hapusUser); err != nil {
                return err
        }

        // 4. Hapus file fisik — kegagalan non-fatal (data DB sudah bersih)
        for _, path := range filePaths {
                if path == "" {
                        continue
                }
                if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
                        fmt.Printf("[WARN] Gagal hapus file fisik %q: %v\n", path, err)
                }
        }

        return nil
}
