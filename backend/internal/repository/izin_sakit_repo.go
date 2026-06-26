package repository

import (
        "context"

        "github.com/google/uuid"
        "github.com/jackc/pgx/v5/pgxpool"
        "github.com/telpp/emagang/internal/models"
)

type IzinSakitRepository struct {
        db *pgxpool.Pool
}

func NewIzinSakitRepository(db *pgxpool.Pool) *IzinSakitRepository {
        return &IzinSakitRepository{db: db}
}

func (r *IzinSakitRepository) Create(ctx context.Context, req *models.IzinSakitRequest) error {
        query := `INSERT INTO izin_sakit_request (pelaksanaan_id, user_id, tanggal, jenis, alasan, bukti_path)
                  VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at`
        return r.db.QueryRow(ctx, query,
                req.PelaksanaanID, req.UserID, req.Tanggal, req.Jenis, req.Alasan, req.BuktiPath,
        ).Scan(&req.ID, &req.CreatedAt)
}

func (r *IzinSakitRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.IzinSakitRequest, error) {
        query := `SELECT id, pelaksanaan_id, user_id, tanggal, jenis, alasan, bukti_path, status,
                  catatan_hrd, diproses_oleh, diproses_at, created_at
                  FROM izin_sakit_request WHERE user_id=$1 ORDER BY tanggal DESC`
        rows, err := r.db.Query(ctx, query, userID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()
        list := []models.IzinSakitRequest{}
        for rows.Next() {
                var req models.IzinSakitRequest
                if err := rows.Scan(
                        &req.ID, &req.PelaksanaanID, &req.UserID, &req.Tanggal, &req.Jenis, &req.Alasan,
                        &req.BuktiPath, &req.Status, &req.CatatanHRD, &req.DiprosesOleh, &req.DiprosesAt, &req.CreatedAt,
                ); err == nil {
                        list = append(list, req)
                }
        }
        return list, nil
}

func (r *IzinSakitRepository) FindAll(ctx context.Context) ([]models.IzinSakitRequestDetail, error) {
        query := `SELECT r.id, r.pelaksanaan_id, r.user_id, r.tanggal, r.jenis, r.alasan, r.bukti_path,
                  r.status, r.catatan_hrd, r.diproses_oleh, r.diproses_at, r.created_at,
                  p.nama_lengkap, pl.divisi
                  FROM izin_sakit_request r
                  JOIN pelaksanaan_magang pl ON pl.id = r.pelaksanaan_id
                  JOIN pengajuan_magang p ON p.id = pl.pengajuan_id
                  ORDER BY r.created_at DESC`
        rows, err := r.db.Query(ctx, query)
        if err != nil {
                return nil, err
        }
        defer rows.Close()
        list := []models.IzinSakitRequestDetail{}
        for rows.Next() {
                var req models.IzinSakitRequestDetail
                if err := rows.Scan(
                        &req.ID, &req.PelaksanaanID, &req.UserID, &req.Tanggal, &req.Jenis, &req.Alasan,
                        &req.BuktiPath, &req.Status, &req.CatatanHRD, &req.DiprosesOleh, &req.DiprosesAt, &req.CreatedAt,
                        &req.NamaPeserta, &req.Divisi,
                ); err == nil {
                        list = append(list, req)
                }
        }
        return list, nil
}

// ApproveAndInsertAbsensi menyetujui request dan sekaligus menyisipkan record absensi.
func (r *IzinSakitRepository) ApproveAndInsertAbsensi(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) error {
        var pelaksanaanID uuid.UUID
        var tanggal, jenis string

        err := r.db.QueryRow(ctx, `
                UPDATE izin_sakit_request
                SET status='disetujui', diproses_oleh=$1, diproses_at=NOW()
                WHERE id=$2 AND status='pending'
                RETURNING pelaksanaan_id, tanggal::text, jenis`,
                approvedBy, id,
        ).Scan(&pelaksanaanID, &tanggal, &jenis)
        if err != nil {
                return err
        }

        _, err = r.db.Exec(ctx, `
                INSERT INTO absensi (pelaksanaan_id, tanggal, keterangan, ttd_pembimbing, approved_by, approved_at)
                VALUES ($1, $2, $3, true, $4, NOW())
                ON CONFLICT (pelaksanaan_id, tanggal) DO UPDATE
                  SET keterangan = EXCLUDED.keterangan,
                      ttd_pembimbing = true,
                      approved_by = EXCLUDED.approved_by,
                      approved_at = EXCLUDED.approved_at`,
                pelaksanaanID, tanggal, jenis, approvedBy,
        )
        return err
}

func (r *IzinSakitRepository) Tolak(ctx context.Context, id uuid.UUID, rejectedBy uuid.UUID, catatan string) error {
        _, err := r.db.Exec(ctx, `
                UPDATE izin_sakit_request
                SET status='ditolak', diproses_oleh=$1, diproses_at=NOW(), catatan_hrd=$2
                WHERE id=$3 AND status='pending'`,
                rejectedBy, catatan, id,
        )
        return err
}
