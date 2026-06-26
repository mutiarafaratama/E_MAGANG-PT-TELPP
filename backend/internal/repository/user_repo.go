package repository

import (
        "context"
        "fmt"
        "time"

        "github.com/google/uuid"
        "github.com/jackc/pgx/v5/pgxpool"
        "github.com/telpp/emagang/internal/models"
)

type UserRepository struct {
        db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
        return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
        query := `
                INSERT INTO users (id, nama_lengkap, email, password_hash, role)
                VALUES (gen_random_uuid(), $1, $2, $3, $4)
                RETURNING id, created_at, updated_at`
        return r.db.QueryRow(ctx, query, u.NamaLengkap, u.Email, u.PasswordHash, u.Role).
                Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
        u := &models.User{}
        query := `SELECT id, nama_lengkap, email, password_hash, role, is_active, password_changed, created_at, updated_at
                          FROM users WHERE email = $1`
        err := r.db.QueryRow(ctx, query, email).Scan(
                &u.ID, &u.NamaLengkap, &u.Email, &u.PasswordHash,
                &u.Role, &u.IsActive, &u.PasswordChanged, &u.CreatedAt, &u.UpdatedAt,
        )
        if err != nil {
                return nil, err
        }
        return u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
        u := &models.User{}
        query := `SELECT id, nama_lengkap, email, password_hash, role, is_active, password_changed, created_at, updated_at
                          FROM users WHERE id = $1`
        err := r.db.QueryRow(ctx, query, id).Scan(
                &u.ID, &u.NamaLengkap, &u.Email, &u.PasswordHash,
                &u.Role, &u.IsActive, &u.PasswordChanged, &u.CreatedAt, &u.UpdatedAt,
        )
        if err != nil {
                return nil, err
        }
        return u, nil
}

func (r *UserRepository) FindAll(ctx context.Context, role string, page, limit int) ([]models.UserPublic, int, error) {
        offset := (page - 1) * limit
        var args []interface{}
        var where string
        argIdx := 1

        if role != "" {
                where = fmt.Sprintf("WHERE role = $%d", argIdx)
                args = append(args, role)
                argIdx++
        }

        var total int
        countQuery := "SELECT COUNT(*) FROM users " + where
        if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
                return nil, 0, err
        }

        args = append(args, limit, offset)
        query := fmt.Sprintf(`
                SELECT u.id, u.nama_lengkap, u.email, u.role, u.is_active, u.password_changed, u.created_at,
                       CASE WHEN u.role = 'peserta' THEN (
                           SELECT COALESCE(pm.status::text,
                               COALESCE((SELECT pj.status::text FROM pengajuan_magang pj WHERE pj.user_id = u.id ORDER BY pj.created_at DESC LIMIT 1),
                               'belum_daftar'))
                           FROM pelaksanaan_magang pm WHERE pm.user_id = u.id ORDER BY pm.created_at DESC LIMIT 1
                       ) ELSE NULL END AS status_magang
                FROM users u %s ORDER BY u.created_at DESC LIMIT $%d OFFSET $%d`,
                where, argIdx, argIdx+1)

        rows, err := r.db.Query(ctx, query, args...)
        if err != nil {
                return nil, 0, err
        }
        defer rows.Close()

        var users []models.UserPublic
        for rows.Next() {
                var u models.UserPublic
                if err := rows.Scan(&u.ID, &u.NamaLengkap, &u.Email, &u.Role, &u.IsActive, &u.PasswordChanged, &u.CreatedAt, &u.StatusMagang); err != nil {
                        return nil, 0, err
                }
                users = append(users, u)
        }
        return users, total, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, namaLengkap string, role models.UserRole) error {
        _, err := r.db.Exec(ctx,
                "UPDATE users SET nama_lengkap = $1, role = $2, updated_at = NOW() WHERE id = $3",
                namaLengkap, role, id)
        return err
}

func (r *UserRepository) UpdateActive(ctx context.Context, id uuid.UUID, isActive bool) error {
        _, err := r.db.Exec(ctx,
                "UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2",
                isActive, id)
        return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error {
        _, err := r.db.Exec(ctx,
                "UPDATE users SET password_hash = $1, password_changed = true, updated_at = NOW() WHERE id = $2",
                hash, id)
        return err
}

func (r *UserRepository) FindHRDList(ctx context.Context) ([]models.UserPublic, error) {
        rows, err := r.db.Query(ctx,
                `SELECT id, nama_lengkap, email, role, is_active, created_at
                 FROM users WHERE role IN ('hrd', 'admin') AND is_active = true ORDER BY nama_lengkap`)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var users []models.UserPublic
        for rows.Next() {
                var u models.UserPublic
                if err := rows.Scan(&u.ID, &u.NamaLengkap, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt); err != nil {
                        return nil, err
                }
                users = append(users, u)
        }
        return users, nil
}

// DeleteByID — hapus user permanen (cascade ke refresh_tokens & notifikasi)
func (r *UserRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
        _, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
        return err
}

// PesertaHapusInfo — data yang diperlukan untuk kirim backup email sebelum hapus
type PesertaHapusInfo struct {
        UserID          uuid.UUID
        Email           string
        NamaLengkap     string
        PelaksanaanID   *uuid.UUID
        StatusMagang    *string
        SertifikatPath  *string
        Nilai           *float64
        CatatanNilai    *string
        Divisi          *string
        PembimbingNama  *string
        TanggalMulai    *string
        TanggalSelesai  *string
}

// DokumenBackup — info dokumen untuk lampiran email backup
type DokumenBackup struct {
        Jenis    string
        NamaFile string
        PathFile string // absolute path di filesystem
}

// KejuruanBackup — satu komponen kejuruan untuk backup penilaian
type KejuruanBackup struct {
        Komponen string
        Nilai    float64
        Urutan   int
}

// PenilaianBackup — seluruh data penilaian peserta untuk generate PDF backup
type PenilaianBackup struct {
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
        DinilaiAt     string // formatted "02 Jan 2006"

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
        Kejuruan   []KejuruanBackup
}

// GetPenilaianPeserta — ambil data penilaian lengkap untuk generate PDF backup sebelum hapus akun
// Mengembalikan nil jika peserta belum memiliki penilaian
func (r *UserRepository) GetPenilaianPeserta(ctx context.Context, userID uuid.UUID) (*PenilaianBackup, error) {
        deref := func(f *float64) float64 {
                if f != nil {
                        return *f
                }
                return 0
        }

        p := &PenilaianBackup{}
        var penilaianID uuid.UUID
        var managerNama, catatan *string
        var dinilaiAt *time.Time
        var nm, ni, ndw, nkr, nkre, ntj, nkerj, napt, nhd *float64
        var nk3s, nk3m, nk3mn, nk3v *float64
        var nprs1, nprs2, nprs3, nprs4 *float64
        var nilaiAkhir *float64

        err := r.db.QueryRow(ctx, `
                SELECT
                    u.nama_lengkap, pj.nomor_induk, pj.asal_institusi, pj.jurusan, pj.kelas_semester,
                    COALESCE(pm.divisi, '') AS divisi,
                    COALESCE(pm.pembimbing, hu.nama_lengkap, '') AS pembimbing,
                    to_char(pm.tanggal_mulai,'DD Mon YYYY')||' s/d '||to_char(pm.tanggal_selesai,'DD Mon YYYY') AS periode,
                    p.nilai_motivasi, p.nilai_inisiatif, p.nilai_disiplin_waktu,
                    p.nilai_kerajinan, p.nilai_kreativitas, p.nilai_tanggung_jawab,
                    p.nilai_kerjasama, p.nilai_adaptasi, p.nilai_kehadiran,
                    p.nilai_k3_safety, p.nilai_k3_metode, p.nilai_k3_manajemen, p.nilai_k3_volume,
                    p.nilai_prs_proses, p.nilai_prs_teori, p.nilai_prs_judul, p.nilai_prs_data,
                    p.nilai_akhir, p.manager_nama, p.catatan, p.dinilai_at, p.id
                FROM users u
                JOIN pelaksanaan_magang pm ON pm.user_id = u.id
                JOIN pengajuan_magang pj ON pj.id = pm.pengajuan_id
                LEFT JOIN users hu ON hu.id = pm.pembimbing_id
                JOIN penilaian p ON p.pelaksanaan_id = pm.id
                WHERE u.id = $1
                ORDER BY pm.created_at DESC LIMIT 1`, userID).
                Scan(
                        &p.NamaLengkap, &p.NomorInduk, &p.AsalInstitusi, &p.Jurusan, &p.KelasSemester,
                        &p.Divisi, &p.Pembimbing, &p.Periode,
                        &nm, &ni, &ndw, &nkr, &nkre, &ntj, &nkerj, &napt, &nhd,
                        &nk3s, &nk3m, &nk3mn, &nk3v,
                        &nprs1, &nprs2, &nprs3, &nprs4,
                        &nilaiAkhir, &managerNama, &catatan, &dinilaiAt, &penilaianID)
        if err != nil {
                return nil, err
        }

        p.NilaiMotivasi = deref(nm)
        p.NilaiInisiatif = deref(ni)
        p.NilaiDisiplinWaktu = deref(ndw)
        p.NilaiKerajinan = deref(nkr)
        p.NilaiKreativitas = deref(nkre)
        p.NilaiTanggungJawab = deref(ntj)
        p.NilaiKerjasama = deref(nkerj)
        p.NilaiAdaptasi = deref(napt)
        p.NilaiKehadiran = deref(nhd)
        p.NilaiK3Safety = deref(nk3s)
        p.NilaiK3Metode = deref(nk3m)
        p.NilaiK3Manajemen = deref(nk3mn)
        p.NilaiK3Volume = deref(nk3v)
        p.NilaiPrsProses = deref(nprs1)
        p.NilaiPrsTeori = deref(nprs2)
        p.NilaiPrsJudul = deref(nprs3)
        p.NilaiPrsData = deref(nprs4)
        p.NilaiAkhir = deref(nilaiAkhir)

        if managerNama != nil {
                p.ManagerNama = *managerNama
        }
        if catatan != nil {
                p.Catatan = *catatan
        }
        if dinilaiAt != nil {
                p.DinilaiAt = dinilaiAt.Format("02 Jan 2006")
        }

        rows, err := r.db.Query(ctx,
                `SELECT komponen, nilai, urutan FROM penilaian_kejuruan WHERE penilaian_id = $1 ORDER BY urutan ASC`,
                penilaianID)
        if err == nil {
                defer rows.Close()
                for rows.Next() {
                        var k KejuruanBackup
                        if scanErr := rows.Scan(&k.Komponen, &k.Nilai, &k.Urutan); scanErr == nil {
                                p.Kejuruan = append(p.Kejuruan, k)
                        }
                }
        }

        return p, nil
}

// GetPesertaHapusInfo — ambil semua info peserta + pelaksanaan untuk keperluan backup sebelum hapus
func (r *UserRepository) GetPesertaHapusInfo(ctx context.Context, userID uuid.UUID) (*PesertaHapusInfo, error) {
        info := &PesertaHapusInfo{UserID: userID}
        err := r.db.QueryRow(ctx, `
                SELECT u.email, u.nama_lengkap,
                       pm.id, pm.status::text, pm.sertifikat_path, pm.nilai, pm.catatan_nilai,
                       pm.divisi,
                       COALESCE(pm.pembimbing, hu.nama_lengkap),
                       to_char(pm.tanggal_mulai, 'DD Mon YYYY'),
                       to_char(pm.tanggal_selesai, 'DD Mon YYYY')
                FROM users u
                LEFT JOIN pelaksanaan_magang pm ON pm.user_id = u.id
                LEFT JOIN users hu ON hu.id = pm.pembimbing_id
                WHERE u.id = $1 AND u.role = 'peserta'
                ORDER BY pm.created_at DESC LIMIT 1`, userID,
        ).Scan(&info.Email, &info.NamaLengkap,
                &info.PelaksanaanID, &info.StatusMagang, &info.SertifikatPath, &info.Nilai, &info.CatatanNilai,
                &info.Divisi, &info.PembimbingNama, &info.TanggalMulai, &info.TanggalSelesai)
        if err != nil {
                return nil, err
        }
        return info, nil
}

// HapusPesertaData — hapus semua data peserta (cascade melalui user_id/FK)
// Urutan: absensi → laporan → perpanjangan → pelaksanaan → dokumen → izin_sakit → pengajuan → notifikasi → refresh_tokens → user
func (r *UserRepository) HapusPesertaData(ctx context.Context, userID uuid.UUID) error {
        tx, err := r.db.Begin(ctx)
        if err != nil {
                return err
        }
        defer tx.Rollback(ctx)

        steps := []string{
                `DELETE FROM absensi WHERE pelaksanaan_id IN (SELECT id FROM pelaksanaan_magang WHERE user_id = $1)`,
                `DELETE FROM laporan_magang WHERE pelaksanaan_id IN (SELECT id FROM pelaksanaan_magang WHERE user_id = $1)`,
                `DELETE FROM perpanjangan_magang WHERE pelaksanaan_id IN (SELECT id FROM pelaksanaan_magang WHERE user_id = $1)`,
                `DELETE FROM izin_sakit_request WHERE pelaksanaan_id IN (SELECT id FROM pelaksanaan_magang WHERE user_id = $1)`,
                `DELETE FROM pelaksanaan_magang WHERE user_id = $1`,
                `DELETE FROM dokumen WHERE user_id = $1 OR pengajuan_id IN (SELECT id FROM pengajuan_magang WHERE user_id = $1)`,
                `DELETE FROM pengajuan_magang WHERE user_id = $1`,
                `DELETE FROM notifikasi WHERE user_id = $1`,
                `DELETE FROM chat_pesan WHERE sender_id = $1`,
                `DELETE FROM chat_tiket WHERE user_id = $1`,
                `DELETE FROM refresh_tokens WHERE user_id = $1`,
                `DELETE FROM users WHERE id = $1`,
        }
        for _, q := range steps {
                if _, err := tx.Exec(ctx, q, userID); err != nil {
                        return err
                }
        }
        return tx.Commit(ctx)
}

// GetDokumenPeserta — ambil dokumen PDF penting milik peserta untuk backup email sebelum hapus
// Hanya ambil jenis laporan_magang dan sertifikat; dokumen KTP/KTM/foto tidak perlu dikirim
func (r *UserRepository) GetDokumenPeserta(ctx context.Context, userID uuid.UUID) ([]DokumenBackup, error) {
        rows, err := r.db.Query(ctx, `
                SELECT jenis::text, nama_file, path_file
                FROM dokumen
                WHERE (user_id = $1
                   OR pengajuan_id IN (SELECT id FROM pengajuan_magang WHERE user_id = $1))
                  AND jenis IN ('laporan_magang', 'sertifikat')
                ORDER BY uploaded_at ASC`, userID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var docs []DokumenBackup
        for rows.Next() {
                var d DokumenBackup
                if err := rows.Scan(&d.Jenis, &d.NamaFile, &d.PathFile); err != nil {
                        return nil, err
                }
                docs = append(docs, d)
        }
        return docs, nil
}

func (r *UserRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string) error {
        _, err := r.db.Exec(ctx,
                `INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
                 VALUES ($1, $2, NOW() + INTERVAL '7 days')`,
                userID, tokenHash)
        return err
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
        _, err := r.db.Exec(ctx,
                "DELETE FROM refresh_tokens WHERE token_hash = $1", tokenHash)
        return err
}

func (r *UserRepository) GetDashboardStats(ctx context.Context) (*models.DashboardAdminStats, error) {
        stats := &models.DashboardAdminStats{}
        err := r.db.QueryRow(ctx, `
                SELECT
                        (SELECT COUNT(*) FROM users WHERE role = 'peserta' AND is_active = true),
                        (SELECT COUNT(*) FROM pengajuan_magang),
                        (SELECT COUNT(*) FROM pengajuan_magang WHERE status IN ('diajukan','menunggu_verifikasi','diproses')),
                        (SELECT COUNT(*) FROM pengajuan_magang WHERE status = 'diterima'),
                        (SELECT COUNT(*) FROM pengajuan_magang WHERE status = 'ditolak'),
                        (SELECT COUNT(*) FROM pelaksanaan_magang WHERE status = 'aktif'),
                        (SELECT COUNT(*) FROM pelaksanaan_magang WHERE status = 'selesai'),
                        (SELECT COUNT(*) FROM chat_tiket WHERE status = 'menunggu')
        `).Scan(
                &stats.TotalPeserta,
                &stats.TotalPengajuan,
                &stats.PengajuanMenunggu,
                &stats.PengajuanDiterima,
                &stats.PengajuanDitolak,
                &stats.MagangAktif,
                &stats.MagangSelesai,
                &stats.TiketChatMenunggu,
        )
        return stats, err
}

// AkunJatuhTempoItem — info singkat peserta yang akan segera dihapus otomatis
type AkunJatuhTempoItem struct {
        UserID        string  `json:"user_id"`
        NamaLengkap   string  `json:"nama_lengkap"`
        Email         string  `json:"email"`
        AsalInstitusi string  `json:"asal_institusi"`
        Divisi        *string `json:"divisi"`
        TanggalSelesai string `json:"tanggal_selesai"`
        SisaHari      int     `json:"sisa_hari"`
        HapusAt       string  `json:"hapus_at"`
}

// GetAkunJatuhTempo — ambil peserta yang akan dihapus otomatis dalam N hari ke depan
// "jatuh tempo" = tanggal_selesai + 30 hari; default N=7
func (r *UserRepository) GetAkunJatuhTempo(ctx context.Context, dalahNHari int) ([]AkunJatuhTempoItem, error) {
        if dalahNHari <= 0 {
                dalahNHari = 7
        }
        rows, err := r.db.Query(ctx, `
                SELECT
                        u.id::text,
                        u.nama_lengkap,
                        u.email,
                        COALESCE(pj.asal_institusi, '') AS asal_institusi,
                        pm.divisi,
                        to_char(pm.tanggal_selesai, 'DD Mon YYYY') AS tanggal_selesai,
                        GREATEST(0, EXTRACT(DAY FROM (pm.tanggal_selesai + INTERVAL '30 days') - NOW())::int) AS sisa_hari,
                        to_char(pm.tanggal_selesai + INTERVAL '30 days', 'DD Mon YYYY') AS hapus_at
                FROM users u
                JOIN pelaksanaan_magang pm ON pm.user_id = u.id
                LEFT JOIN pengajuan_magang pj ON pj.id = pm.pengajuan_id
                WHERE u.role = 'peserta'
                  AND pm.status = 'selesai'
                  AND pm.tanggal_selesai > NOW() - INTERVAL '30 days'
                  AND pm.tanggal_selesai <= NOW() - INTERVAL '30 days' + ($1 * INTERVAL '1 day')
                  AND pm.created_at = (
                      SELECT MAX(pm2.created_at) FROM pelaksanaan_magang pm2 WHERE pm2.user_id = u.id
                  )
                ORDER BY pm.tanggal_selesai ASC
        `, dalahNHari)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var result []AkunJatuhTempoItem
        for rows.Next() {
                var item AkunJatuhTempoItem
                if err := rows.Scan(
                        &item.UserID, &item.NamaLengkap, &item.Email,
                        &item.AsalInstitusi, &item.Divisi,
                        &item.TanggalSelesai, &item.SisaHari, &item.HapusAt,
                ); err != nil {
                        continue
                }
                result = append(result, item)
        }
        if result == nil {
                result = []AkunJatuhTempoItem{}
        }
        return result, rows.Err()
}

// GetKandidatHapusOtomatis — ambil ID peserta yang magang sudah selesai > 30 hari
// Hanya mempertimbangkan pelaksanaan terbaru per user
func (r *UserRepository) GetKandidatHapusOtomatis(ctx context.Context) ([]uuid.UUID, error) {
        rows, err := r.db.Query(ctx, `
                SELECT u.id
                FROM users u
                JOIN pelaksanaan_magang pm ON pm.user_id = u.id
                WHERE u.role = 'peserta'
                  AND pm.status = 'selesai'
                  AND pm.tanggal_selesai <= NOW() - INTERVAL '30 days'
                  AND pm.created_at = (
                      SELECT MAX(pm2.created_at)
                      FROM pelaksanaan_magang pm2
                      WHERE pm2.user_id = u.id
                  )
                ORDER BY pm.tanggal_selesai ASC
        `)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var ids []uuid.UUID
        for rows.Next() {
                var id uuid.UUID
                if err := rows.Scan(&id); err != nil {
                        continue
                }
                ids = append(ids, id)
        }
        return ids, rows.Err()
}
