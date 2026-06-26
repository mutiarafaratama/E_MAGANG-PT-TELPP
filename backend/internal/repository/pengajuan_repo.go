package repository

import (
        "context"
        "fmt"

        "github.com/google/uuid"
        "github.com/jackc/pgx/v5/pgxpool"
        "github.com/telpp/emagang/internal/models"
)

type PengajuanRepository struct {
        db *pgxpool.Pool
}

func NewPengajuanRepository(db *pgxpool.Pool) *PengajuanRepository {
        return &PengajuanRepository{db: db}
}

// Create — pengajuan dengan user_id (peserta sudah login)
func (r *PengajuanRepository) Create(ctx context.Context, p *models.PengajuanMagang) error {
        query := `
                INSERT INTO pengajuan_magang (
                        user_id, nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
                        alamat, no_hp, email, kategori_magang, nomor_induk,
                        asal_institusi, jurusan, kelas_semester
                ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
                RETURNING id, status, created_at, updated_at`
        return r.db.QueryRow(ctx, query,
                p.UserID, p.NamaLengkap, p.TempatLahir, p.TanggalLahir, p.JenisKelamin,
                p.Alamat, p.NoHP, p.Email, p.KategoriMagang, p.NomorInduk,
                p.AsalInstitusi, p.Jurusan, p.KelasSemester,
        ).Scan(&p.ID, &p.Status, &p.CreatedAt, &p.UpdatedAt)
}

// CreatePublik — pengajuan publik tanpa user_id
func (r *PengajuanRepository) CreatePublik(ctx context.Context, p *models.PengajuanMagang) error {
        query := `
                INSERT INTO pengajuan_magang (
                        nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
                        alamat, no_hp, email, kategori_magang, nomor_induk,
                        asal_institusi, jurusan, kelas_semester
                ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
                RETURNING id, status, created_at, updated_at`
        return r.db.QueryRow(ctx, query,
                p.NamaLengkap, p.TempatLahir, p.TanggalLahir, p.JenisKelamin,
                p.Alamat, p.NoHP, p.Email, p.KategoriMagang, p.NomorInduk,
                p.AsalInstitusi, p.Jurusan, p.KelasSemester,
        ).Scan(&p.ID, &p.Status, &p.CreatedAt, &p.UpdatedAt)
}

// SetAkunTerkirim — set user_id dan akun_terkirim_at setelah HRD kirim akun
func (r *PengajuanRepository) SetAkunTerkirim(ctx context.Context, pengajuanID, userID uuid.UUID) error {
        _, err := r.db.Exec(ctx,
                `UPDATE pengajuan_magang SET user_id = $1, akun_terkirim_at = NOW(), updated_at = NOW() WHERE id = $2`,
                userID, pengajuanID)
        return err
}

func (r *PengajuanRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PengajuanMagang, error) {
        p := &models.PengajuanMagang{}
        query := `SELECT id,
                COALESCE(user_id, '00000000-0000-0000-0000-000000000000'::uuid),
                nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
                alamat, no_hp, email, kategori_magang, nomor_induk, asal_institusi, jurusan,
                kelas_semester, status, catatan_hrd, verified_by, verified_at, akun_terkirim_at, created_at, updated_at
                FROM pengajuan_magang WHERE id = $1`
        err := r.db.QueryRow(ctx, query, id).Scan(
                &p.ID, &p.UserID, &p.NamaLengkap, &p.TempatLahir, &p.TanggalLahir, &p.JenisKelamin,
                &p.Alamat, &p.NoHP, &p.Email, &p.KategoriMagang, &p.NomorInduk, &p.AsalInstitusi, &p.Jurusan,
                &p.KelasSemester, &p.Status, &p.CatatanHRD, &p.VerifiedBy, &p.VerifiedAt, &p.AkunTerkirimAt,
                &p.CreatedAt, &p.UpdatedAt,
        )
        return p, err
}

func (r *PengajuanRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.PengajuanMagang, error) {
        p := &models.PengajuanMagang{}
        query := `SELECT id,
                COALESCE(user_id, '00000000-0000-0000-0000-000000000000'::uuid),
                nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
                alamat, no_hp, email, kategori_magang, nomor_induk, asal_institusi, jurusan,
                kelas_semester, status, catatan_hrd, verified_by, verified_at, akun_terkirim_at, created_at, updated_at
                FROM pengajuan_magang WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`
        err := r.db.QueryRow(ctx, query, userID).Scan(
                &p.ID, &p.UserID, &p.NamaLengkap, &p.TempatLahir, &p.TanggalLahir, &p.JenisKelamin,
                &p.Alamat, &p.NoHP, &p.Email, &p.KategoriMagang, &p.NomorInduk, &p.AsalInstitusi, &p.Jurusan,
                &p.KelasSemester, &p.Status, &p.CatatanHRD, &p.VerifiedBy, &p.VerifiedAt, &p.AkunTerkirimAt,
                &p.CreatedAt, &p.UpdatedAt,
        )
        return p, err
}

func (r *PengajuanRepository) FindAll(ctx context.Context, status, search string, page, limit int) ([]models.PengajuanMagang, int, error) {
        offset := (page - 1) * limit
        var args []interface{}
        where := "WHERE 1=1"
        argIdx := 1

        if status != "" {
                where += fmt.Sprintf(" AND status = $%d", argIdx)
                args = append(args, status)
                argIdx++
        }
        if search != "" {
                where += fmt.Sprintf(" AND (nama_lengkap ILIKE $%d OR email ILIKE $%d OR asal_institusi ILIKE $%d)", argIdx, argIdx, argIdx)
                args = append(args, "%"+search+"%")
                argIdx++
        }

        var total int
        if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pengajuan_magang "+where, args...).Scan(&total); err != nil {
                return nil, 0, err
        }

        args = append(args, limit, offset)
        query := fmt.Sprintf(`
                SELECT id,
                COALESCE(user_id, '00000000-0000-0000-0000-000000000000'::uuid),
                nama_lengkap, tempat_lahir, tanggal_lahir, jenis_kelamin,
                alamat, no_hp, email, kategori_magang, nomor_induk, asal_institusi, jurusan,
                kelas_semester, status, catatan_hrd, verified_by, verified_at, akun_terkirim_at, created_at, updated_at
                FROM pengajuan_magang %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
                where, argIdx, argIdx+1)

        rows, err := r.db.Query(ctx, query, args...)
        if err != nil {
                return nil, 0, err
        }
        defer rows.Close()

        list := make([]models.PengajuanMagang, 0)
        for rows.Next() {
                var p models.PengajuanMagang
                if err := rows.Scan(
                        &p.ID, &p.UserID, &p.NamaLengkap, &p.TempatLahir, &p.TanggalLahir, &p.JenisKelamin,
                        &p.Alamat, &p.NoHP, &p.Email, &p.KategoriMagang, &p.NomorInduk, &p.AsalInstitusi, &p.Jurusan,
                        &p.KelasSemester, &p.Status, &p.CatatanHRD, &p.VerifiedBy, &p.VerifiedAt, &p.AkunTerkirimAt,
                        &p.CreatedAt, &p.UpdatedAt,
                ); err != nil {
                        return nil, 0, err
                }
                list = append(list, p)
        }
        return list, total, nil
}

func (r *PengajuanRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.StatusPengajuan, catatan *string, changedBy uuid.UUID) error {
        tx, err := r.db.Begin(ctx)
        if err != nil {
                return err
        }
        defer tx.Rollback(ctx)

        var oldStatus models.StatusPengajuan
        if err := tx.QueryRow(ctx, "SELECT status FROM pengajuan_magang WHERE id = $1", id).Scan(&oldStatus); err != nil {
                return err
        }

        _, err = tx.Exec(ctx,
                "UPDATE pengajuan_magang SET status=$1, catatan_hrd=$2, verified_by=$3, verified_at=NOW(), updated_at=NOW() WHERE id=$4",
                status, catatan, changedBy, id)
        if err != nil {
                return err
        }

        _, err = tx.Exec(ctx,
                "INSERT INTO status_history (pengajuan_id, status_lama, status_baru, changed_by, catatan) VALUES ($1,$2,$3,$4,$5)",
                id, oldStatus, status, changedBy, catatan)
        if err != nil {
                return err
        }

        return tx.Commit(ctx)
}

func (r *PengajuanRepository) GetStatusHistory(ctx context.Context, pengajuanID uuid.UUID) ([]models.StatusHistory, error) {
        rows, err := r.db.Query(ctx,
                `SELECT id, pengajuan_id, status_lama, status_baru, changed_by, catatan, created_at
                 FROM status_history WHERE pengajuan_id = $1 ORDER BY created_at ASC`, pengajuanID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var list []models.StatusHistory
        for rows.Next() {
                var h models.StatusHistory
                if err := rows.Scan(&h.ID, &h.PengajuanID, &h.StatusLama, &h.StatusBaru, &h.ChangedBy, &h.Catatan, &h.CreatedAt); err != nil {
                        return nil, err
                }
                list = append(list, h)
        }
        return list, nil
}

// HapusLengkap — hapus pengajuan beserta seluruh data terkait dalam satu transaksi.
//
// Urutan hapus (menghindari FK violation):
//  1. chat_tiket milik peserta (cascade → chat_pesan) — jika hapusUser = true
//  2. pelaksanaan_magang (cascade → absensi_harian)
//  3. pengajuan_magang (cascade → dokumen, status_history)
//  4. users (cascade → refresh_tokens, notifikasi) — jika hapusUser = true
//
// Penghapusan file fisik dilakukan oleh caller setelah fungsi ini sukses.
func (r *PengajuanRepository) HapusLengkap(ctx context.Context, id, userID uuid.UUID, hapusUser bool) error {
        tx, err := r.db.Begin(ctx)
        if err != nil {
                return fmt.Errorf("gagal mulai transaksi: %w", err)
        }
        defer tx.Rollback(ctx)

        // 1. Hapus chat_tiket milik peserta (chat_pesan ter-cascade)
        if hapusUser && userID != uuid.Nil {
                if _, err := tx.Exec(ctx, `DELETE FROM chat_tiket WHERE user_id = $1`, userID); err != nil {
                        return fmt.Errorf("gagal hapus chat tiket: %w", err)
                }
        }

        // 2. Hapus pelaksanaan_magang (absensi_harian ter-cascade)
        if _, err := tx.Exec(ctx, `DELETE FROM pelaksanaan_magang WHERE pengajuan_id = $1`, id); err != nil {
                return fmt.Errorf("gagal hapus pelaksanaan: %w", err)
        }

        // 3. Hapus pengajuan_magang (dokumen & status_history ter-cascade)
        tag, err := tx.Exec(ctx, `DELETE FROM pengajuan_magang WHERE id = $1`, id)
        if err != nil {
                return fmt.Errorf("gagal hapus pengajuan: %w", err)
        }
        if tag.RowsAffected() == 0 {
                return fmt.Errorf("pengajuan tidak ditemukan")
        }

        // 4. Hapus user peserta (refresh_tokens & notifikasi ter-cascade)
        if hapusUser && userID != uuid.Nil {
                if _, err := tx.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID); err != nil {
                        return fmt.Errorf("gagal hapus akun peserta: %w", err)
                }
        }

        return tx.Commit(ctx)
}
