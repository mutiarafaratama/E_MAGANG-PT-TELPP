package repository

import (
        "context"
        "fmt"

        "github.com/google/uuid"
        "github.com/jackc/pgx/v5/pgxpool"
        "github.com/telpp/emagang/internal/models"
)

type PerpanjanganRepository struct {
        db *pgxpool.Pool
}

func NewPerpanjanganRepository(db *pgxpool.Pool) *PerpanjanganRepository {
        return &PerpanjanganRepository{db: db}
}

// cekBisaPerpanjang — validasi: pelaksanaan ada dan belum pernah diperpanjang
func (r *PerpanjanganRepository) cekBisaPerpanjang(ctx context.Context, pelaksanaanID uuid.UUID) error {
        var sudah bool
        err := r.db.QueryRow(ctx,
                `SELECT sudah_diperpanjang FROM pelaksanaan_magang WHERE id = $1`,
                pelaksanaanID,
        ).Scan(&sudah)
        if err != nil {
                return fmt.Errorf("pelaksanaan tidak ditemukan")
        }
        if sudah {
                return fmt.Errorf("sudah_diperpanjang")
        }
        return nil
}

// scanPerp — scan satu baris PerpanjanganMagang
func scanPerp(row interface {
        Scan(...any) error
}, p *models.PerpanjanganMagang) error {
        return row.Scan(
                &p.ID, &p.PelaksanaanID, &p.DiajukanOleh, &p.RolePengaju, &p.DurasiHari,
                &p.TanggalSelesaiLama, &p.TanggalSelesaiBaru, &p.Alasan, &p.SuratPath,
                &p.Status, &p.CatatanHRD, &p.DiprosesOleh, &p.DiprosesAt, &p.CreatedAt, &p.UpdatedAt,
        )
}

const perpColumns = `id, pelaksanaan_id, diajukan_oleh, role_pengaju, durasi_hari,
       tanggal_selesai_lama, tanggal_selesai_baru, alasan, surat_path,
       status, catatan_hrd, diproses_oleh, diproses_at, created_at, updated_at`

// CreatePeserta — peserta ajukan perpanjangan (status='menunggu', perlu approval HRD)
func (r *PerpanjanganRepository) CreatePeserta(ctx context.Context, pelaksanaanID, userID uuid.UUID, durasiHari int, alasan string) (*models.PerpanjanganMagang, error) {
        if err := r.cekBisaPerpanjang(ctx, pelaksanaanID); err != nil {
                return nil, err
        }

        var p models.PerpanjanganMagang
        err := scanPerp(r.db.QueryRow(ctx, `
                INSERT INTO perpanjangan_magang
                        (pelaksanaan_id, diajukan_oleh, role_pengaju, durasi_hari,
                         tanggal_selesai_lama, tanggal_selesai_baru, alasan, status)
                SELECT
                        $1, $2, 'peserta', $3,
                        pm.tanggal_selesai,
                        pm.tanggal_selesai + make_interval(days => $3),
                        $4, 'menunggu'
                FROM pelaksanaan_magang pm WHERE pm.id = $1
                RETURNING `+perpColumns,
                pelaksanaanID, userID, durasiHari, alasan,
        ), &p)
        return &p, err
}

// CreateHRDDirect — HRD perpanjang langsung: insert record + update pelaksanaan dalam 1 tx
func (r *PerpanjanganRepository) CreateHRDDirect(ctx context.Context, pelaksanaanID, hrdID uuid.UUID, durasiHari int, alasan string, suratPath *string) (*models.PerpanjanganMagang, error) {
        if err := r.cekBisaPerpanjang(ctx, pelaksanaanID); err != nil {
                return nil, err
        }

        tx, err := r.db.Begin(ctx)
        if err != nil {
                return nil, err
        }
        defer tx.Rollback(ctx)

        var p models.PerpanjanganMagang
        err = scanPerp(tx.QueryRow(ctx, `
                INSERT INTO perpanjangan_magang
                        (pelaksanaan_id, diajukan_oleh, role_pengaju, durasi_hari,
                         tanggal_selesai_lama, tanggal_selesai_baru, alasan, surat_path,
                         status, diproses_oleh, diproses_at)
                SELECT
                        $1, $2, 'hrd', $3,
                        pm.tanggal_selesai,
                        pm.tanggal_selesai + make_interval(days => $3),
                        $4, $5,
                        'disetujui', $2, NOW()
                FROM pelaksanaan_magang pm WHERE pm.id = $1
                RETURNING `+perpColumns,
                pelaksanaanID, hrdID, durasiHari, alasan, suratPath,
        ), &p)
        if err != nil {
                return nil, err
        }

        _, err = tx.Exec(ctx, `
                UPDATE pelaksanaan_magang
                SET tanggal_selesai = $1, sudah_diperpanjang = TRUE, updated_at = NOW()
                WHERE id = $2`,
                p.TanggalSelesaiBaru, pelaksanaanID,
        )
        if err != nil {
                return nil, err
        }

        if err := tx.Commit(ctx); err != nil {
                return nil, err
        }
        return &p, nil
}

// FindByPelaksanaan — peserta lihat status perpanjangan miliknya (1 record terakhir)
func (r *PerpanjanganRepository) FindByPelaksanaan(ctx context.Context, pelaksanaanID uuid.UUID) (*models.PerpanjanganMagang, error) {
        var p models.PerpanjanganMagang
        err := scanPerp(r.db.QueryRow(ctx, `
                SELECT `+perpColumns+`
                FROM perpanjangan_magang
                WHERE pelaksanaan_id = $1
                ORDER BY created_at DESC LIMIT 1`,
                pelaksanaanID,
        ), &p)
        if err != nil {
                return nil, err
        }
        return &p, nil
}

// FindAll — HRD lihat semua pengajuan perpanjangan
func (r *PerpanjanganRepository) FindAll(ctx context.Context) ([]models.PerpanjanganMagangDetail, error) {
        rows, err := r.db.Query(ctx, `
                SELECT
                        pp.id, pp.pelaksanaan_id, pp.diajukan_oleh, pp.role_pengaju, pp.durasi_hari,
                        pp.tanggal_selesai_lama, pp.tanggal_selesai_baru, pp.alasan, pp.surat_path,
                        pp.status, pp.catatan_hrd, pp.diproses_oleh, pp.diproses_at, pp.created_at, pp.updated_at,
                        pj.nama_lengkap AS nama_peserta,
                        pj.asal_institusi,
                        pm.divisi,
                        u.nama_lengkap  AS nama_pengaju
                FROM perpanjangan_magang pp
                JOIN pelaksanaan_magang pm ON pm.id  = pp.pelaksanaan_id
                JOIN pengajuan_magang   pj ON pj.id  = pm.pengajuan_id
                JOIN users               u  ON u.id   = pp.diajukan_oleh
                ORDER BY pp.created_at DESC`)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var list []models.PerpanjanganMagangDetail
        for rows.Next() {
                var d models.PerpanjanganMagangDetail
                if err := rows.Scan(
                        &d.ID, &d.PelaksanaanID, &d.DiajukanOleh, &d.RolePengaju, &d.DurasiHari,
                        &d.TanggalSelesaiLama, &d.TanggalSelesaiBaru, &d.Alasan, &d.SuratPath,
                        &d.Status, &d.CatatanHRD, &d.DiprosesOleh, &d.DiprosesAt, &d.CreatedAt, &d.UpdatedAt,
                        &d.NamaPeserta, &d.AsalInstitusi, &d.Divisi, &d.NamaPengaju,
                ); err == nil {
                        list = append(list, d)
                }
        }
        return list, nil
}

// CountPending — untuk badge HRD
func (r *PerpanjanganRepository) CountPending(ctx context.Context) int {
        var count int
        r.db.QueryRow(ctx, `SELECT COUNT(*) FROM perpanjangan_magang WHERE status = 'menunggu'`).Scan(&count)
        return count
}

// Approve — HRD setujui pengajuan peserta (update record + update pelaksanaan dalam 1 tx)
func (r *PerpanjanganRepository) Approve(ctx context.Context, id, approvedBy uuid.UUID, suratPath *string, catatan string) (uuid.UUID, error) {
        tx, err := r.db.Begin(ctx)
        if err != nil {
                return uuid.Nil, err
        }
        defer tx.Rollback(ctx)

        var pelaksanaanID uuid.UUID
        var tanggalBaru interface{}
        var durasiHari int
        err = tx.QueryRow(ctx, `
                UPDATE perpanjangan_magang
                SET status       = 'disetujui',
                    diproses_oleh = $1,
                    diproses_at   = NOW(),
                    catatan_hrd   = NULLIF($2, ''),
                    surat_path    = COALESCE($3, surat_path),
                    updated_at    = NOW()
                WHERE id = $4 AND status = 'menunggu'
                RETURNING pelaksanaan_id, tanggal_selesai_baru, durasi_hari`,
                approvedBy, catatan, suratPath, id,
        ).Scan(&pelaksanaanID, &tanggalBaru, &durasiHari)
        if err != nil {
                return uuid.Nil, fmt.Errorf("pengajuan tidak ditemukan atau sudah diproses")
        }

        _, err = tx.Exec(ctx, `
                UPDATE pelaksanaan_magang
                SET tanggal_selesai      = $1,
                    sudah_diperpanjang   = TRUE,
                    updated_at           = NOW()
                WHERE id = $2`,
                tanggalBaru, pelaksanaanID,
        )
        if err != nil {
                return uuid.Nil, err
        }

        if err := tx.Commit(ctx); err != nil {
                return uuid.Nil, err
        }
        return pelaksanaanID, nil
}

// Tolak — HRD tolak pengajuan peserta
func (r *PerpanjanganRepository) Tolak(ctx context.Context, id, rejectedBy uuid.UUID, catatan string) (uuid.UUID, error) {
        var pelaksanaanID uuid.UUID
        err := r.db.QueryRow(ctx, `
                UPDATE perpanjangan_magang
                SET status        = 'ditolak',
                    diproses_oleh  = $1,
                    diproses_at    = NOW(),
                    catatan_hrd    = NULLIF($2, ''),
                    updated_at     = NOW()
                WHERE id = $3 AND status = 'menunggu'
                RETURNING pelaksanaan_id`,
                rejectedBy, catatan, id,
        ).Scan(&pelaksanaanID)
        return pelaksanaanID, err
}
