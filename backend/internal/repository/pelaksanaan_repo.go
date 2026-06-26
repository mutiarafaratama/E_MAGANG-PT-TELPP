package repository

import (
        "context"
        "fmt"

        "github.com/google/uuid"
        "github.com/jackc/pgx/v5/pgxpool"
        "github.com/telpp/emagang/internal/models"
)

type PelaksanaanRepository struct {
        db *pgxpool.Pool
}

func NewPelaksanaanRepository(db *pgxpool.Pool) *PelaksanaanRepository {
        return &PelaksanaanRepository{db: db}
}

func (r *PelaksanaanRepository) Create(ctx context.Context, p *models.PelaksanaanMagang) error {
        query := `
                INSERT INTO pelaksanaan_magang (pengajuan_id, user_id, periode_id, tanggal_mulai, tanggal_selesai, divisi, pembimbing_id, pembimbing)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                RETURNING id, status, created_at, updated_at`
        return r.db.QueryRow(ctx, query,
                p.PengajuanID, p.UserID, p.PeriodeID, p.TanggalMulai, p.TanggalSelesai, p.Divisi, p.PembimbingID, p.PembimbingNama,
        ).Scan(&p.ID, &p.Status, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PelaksanaanRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PelaksanaanMagang, error) {
        p := &models.PelaksanaanMagang{}
        err := r.db.QueryRow(ctx, `
                SELECT id, pengajuan_id, user_id, periode_id, tanggal_mulai, tanggal_selesai,
                divisi, pembimbing_id, status, nilai, catatan_nilai, dinilai_oleh, dinilai_at,
                sertifikat_generated, sertifikat_path, sertifikat_generated_at, created_at, updated_at
                FROM pelaksanaan_magang WHERE id = $1`, id).
                Scan(&p.ID, &p.PengajuanID, &p.UserID, &p.PeriodeID, &p.TanggalMulai, &p.TanggalSelesai,
                        &p.Divisi, &p.PembimbingID, &p.Status, &p.Nilai, &p.CatatanNilai, &p.DinilaiOleh, &p.DinilaiAt,
                        &p.SertifikatGenerated, &p.SertifikatPath, &p.SertifikatGeneratedAt, &p.CreatedAt, &p.UpdatedAt)
        return p, err
}

func (r *PelaksanaanRepository) FindByPengajuanID(ctx context.Context, pengajuanID uuid.UUID) (*models.PelaksanaanMagang, error) {
        p := &models.PelaksanaanMagang{}
        err := r.db.QueryRow(ctx, `
                SELECT id, pengajuan_id, user_id, periode_id, tanggal_mulai, tanggal_selesai,
                divisi, pembimbing_id, status, nilai, catatan_nilai, dinilai_oleh, dinilai_at,
                sertifikat_generated, sertifikat_path, sertifikat_generated_at, created_at, updated_at
                FROM pelaksanaan_magang WHERE pengajuan_id = $1`, pengajuanID).
                Scan(&p.ID, &p.PengajuanID, &p.UserID, &p.PeriodeID, &p.TanggalMulai, &p.TanggalSelesai,
                        &p.Divisi, &p.PembimbingID, &p.Status, &p.Nilai, &p.CatatanNilai, &p.DinilaiOleh, &p.DinilaiAt,
                        &p.SertifikatGenerated, &p.SertifikatPath, &p.SertifikatGeneratedAt, &p.CreatedAt, &p.UpdatedAt)
        return p, err
}

func (r *PelaksanaanRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.PelaksanaanMagang, error) {
        p := &models.PelaksanaanMagang{}
        err := r.db.QueryRow(ctx, `
                SELECT pm.id, pm.pengajuan_id, pm.user_id, pm.periode_id, pm.tanggal_mulai, pm.tanggal_selesai,
                pm.divisi, pm.pembimbing_id, pm.status, pm.nilai, pm.catatan_nilai, pm.dinilai_oleh, pm.dinilai_at,
                pm.sertifikat_generated, pm.sertifikat_path, pm.sertifikat_generated_at, pm.created_at, pm.updated_at,
                COALESCE(pm.pembimbing, u.nama_lengkap) AS pembimbing_nama
                FROM pelaksanaan_magang pm
                LEFT JOIN users u ON u.id = pm.pembimbing_id
                WHERE pm.user_id = $1 ORDER BY pm.created_at DESC LIMIT 1`, userID).
                Scan(&p.ID, &p.PengajuanID, &p.UserID, &p.PeriodeID, &p.TanggalMulai, &p.TanggalSelesai,
                        &p.Divisi, &p.PembimbingID, &p.Status, &p.Nilai, &p.CatatanNilai, &p.DinilaiOleh, &p.DinilaiAt,
                        &p.SertifikatGenerated, &p.SertifikatPath, &p.SertifikatGeneratedAt, &p.CreatedAt, &p.UpdatedAt,
                        &p.PembimbingNama)
        return p, err
}

func (r *PelaksanaanRepository) FindAll(ctx context.Context, status string, page, limit int) ([]models.PelaksanaanMagang, int, error) {
        where := "WHERE 1=1"
        var args []interface{}
        argIdx := 1
        if status != "" {
                where += fmt.Sprintf(" AND pm.status = $%d", argIdx)
                args = append(args, status)
                argIdx++
        }
        offset := (page - 1) * limit
        var total int
        r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pelaksanaan_magang pm "+where, args...).Scan(&total)

        args = append(args, limit, offset)
        query := fmt.Sprintf(`
                SELECT pm.id, pm.pengajuan_id, pm.user_id, pm.periode_id, pm.tanggal_mulai, pm.tanggal_selesai,
                pm.divisi, pm.pembimbing_id, pm.status, pm.nilai, pm.catatan_nilai, pm.dinilai_oleh, pm.dinilai_at,
                pm.sertifikat_generated, pm.sertifikat_path, pm.sertifikat_generated_at, pm.created_at, pm.updated_at,
                COALESCE(pm.pembimbing, u.nama_lengkap) AS pembimbing_nama,
                COALESCE(pm.sudah_diperpanjang, false) AS sudah_diperpanjang,
                pj.nama_lengkap AS nama_peserta
                FROM pelaksanaan_magang pm
                LEFT JOIN users u ON u.id = pm.pembimbing_id
                LEFT JOIN pengajuan_magang pj ON pj.id = pm.pengajuan_id
                %s ORDER BY pm.created_at DESC LIMIT $%d OFFSET $%d`,
                where, argIdx, argIdx+1)

        rows, err := r.db.Query(ctx, query, args...)
        if err != nil {
                return nil, 0, err
        }
        defer rows.Close()
        var list []models.PelaksanaanMagang
        for rows.Next() {
                var p models.PelaksanaanMagang
                rows.Scan(&p.ID, &p.PengajuanID, &p.UserID, &p.PeriodeID, &p.TanggalMulai, &p.TanggalSelesai,
                        &p.Divisi, &p.PembimbingID, &p.Status, &p.Nilai, &p.CatatanNilai, &p.DinilaiOleh, &p.DinilaiAt,
                        &p.SertifikatGenerated, &p.SertifikatPath, &p.SertifikatGeneratedAt, &p.CreatedAt, &p.UpdatedAt,
                        &p.PembimbingNama, &p.SudahDiperpanjang, &p.NamaPeserta)
                list = append(list, p)
        }
        return list, total, nil
}

// AutoUpdateStatuses mengubah status pelaksanaan secara otomatis berdasarkan tanggal.
//   menunggu_mulai → aktif          : saat tanggal_mulai sudah tiba
//   aktif           → upload_laporan : saat tanggal_selesai sudah lewat
func (r *PelaksanaanRepository) AutoUpdateStatuses(ctx context.Context) {
        r.db.Exec(ctx,
                `UPDATE pelaksanaan_magang
                 SET status = 'aktif', updated_at = NOW()
                 WHERE status = 'menunggu_mulai'
                   AND tanggal_mulai <= CURRENT_DATE`)

        r.db.Exec(ctx,
                `UPDATE pelaksanaan_magang
                 SET status = 'upload_laporan', updated_at = NOW()
                 WHERE status = 'aktif'
                   AND tanggal_selesai < CURRENT_DATE`)
}

func (r *PelaksanaanRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.StatusPelaksanaan) error {
        _, err := r.db.Exec(ctx,
                "UPDATE pelaksanaan_magang SET status=$1, updated_at=NOW() WHERE id=$2", status, id)
        return err
}


func (r *PelaksanaanRepository) SetNilai(ctx context.Context, id uuid.UUID, nilai float64, catatan string, dinilaiOleh uuid.UUID) error {
        _, err := r.db.Exec(ctx,
                `UPDATE pelaksanaan_magang SET nilai=$1, catatan_nilai=$2, dinilai_oleh=$3,
                 dinilai_at=NOW(), status='penilaian', updated_at=NOW() WHERE id=$4`,
                nilai, catatan, dinilaiOleh, id)
        return err
}

func (r *PelaksanaanRepository) UpdatePembimbing(ctx context.Context, id uuid.UUID, pembimbing string) error {
        var val *string
        if pembimbing != "" {
                val = &pembimbing
        }
        _, err := r.db.Exec(ctx,
                `UPDATE pelaksanaan_magang SET pembimbing=$1, pembimbing_id=NULL, updated_at=NOW() WHERE id=$2`,
                val, id)
        return err
}

func (r *PelaksanaanRepository) SetSertifikat(ctx context.Context, id uuid.UUID, path string) error {
        _, err := r.db.Exec(ctx,
                `UPDATE pelaksanaan_magang SET sertifikat_generated=true, sertifikat_path=$1,
                 sertifikat_generated_at=NOW(), status='selesai', updated_at=NOW() WHERE id=$2`,
                path, id)
        return err
}
