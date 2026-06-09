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
		INSERT INTO pelaksanaan_magang (pengajuan_id, user_id, periode_id, tanggal_mulai, tanggal_selesai, divisi, pembimbing_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, status, created_at, updated_at`
	return r.db.QueryRow(ctx, query,
		p.PengajuanID, p.UserID, p.PeriodeID, p.TanggalMulai, p.TanggalSelesai, p.Divisi, p.PembimbingID,
	).Scan(&p.ID, &p.Status, &p.CreatedAt, &p.UpdatedAt)
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
		SELECT id, pengajuan_id, user_id, periode_id, tanggal_mulai, tanggal_selesai,
		divisi, pembimbing_id, status, nilai, catatan_nilai, dinilai_oleh, dinilai_at,
		sertifikat_generated, sertifikat_path, sertifikat_generated_at, created_at, updated_at
		FROM pelaksanaan_magang WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`, userID).
		Scan(&p.ID, &p.PengajuanID, &p.UserID, &p.PeriodeID, &p.TanggalMulai, &p.TanggalSelesai,
			&p.Divisi, &p.PembimbingID, &p.Status, &p.Nilai, &p.CatatanNilai, &p.DinilaiOleh, &p.DinilaiAt,
			&p.SertifikatGenerated, &p.SertifikatPath, &p.SertifikatGeneratedAt, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PelaksanaanRepository) FindAll(ctx context.Context, status string, page, limit int) ([]models.PelaksanaanMagang, int, error) {
	where := "WHERE 1=1"
	var args []interface{}
	argIdx := 1
	if status != "" {
		where += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	offset := (page - 1) * limit
	var total int
	r.db.QueryRow(ctx, "SELECT COUNT(*) FROM pelaksanaan_magang "+where, args...).Scan(&total)

	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT id, pengajuan_id, user_id, periode_id, tanggal_mulai, tanggal_selesai,
		divisi, pembimbing_id, status, nilai, catatan_nilai, dinilai_oleh, dinilai_at,
		sertifikat_generated, sertifikat_path, sertifikat_generated_at, created_at, updated_at
		FROM pelaksanaan_magang %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
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
			&p.SertifikatGenerated, &p.SertifikatPath, &p.SertifikatGeneratedAt, &p.CreatedAt, &p.UpdatedAt)
		list = append(list, p)
	}
	return list, total, nil
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

func (r *PelaksanaanRepository) SetSertifikat(ctx context.Context, id uuid.UUID, path string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE pelaksanaan_magang SET sertifikat_generated=true, sertifikat_path=$1,
		 sertifikat_generated_at=NOW(), status='selesai', updated_at=NOW() WHERE id=$2`,
		path, id)
	return err
}
