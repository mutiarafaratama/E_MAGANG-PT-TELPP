package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/telpp/emagang/internal/models"
)

type AbsensiRepository struct {
	db *pgxpool.Pool
}

func NewAbsensiRepository(db *pgxpool.Pool) *AbsensiRepository {
	return &AbsensiRepository{db: db}
}

func (r *AbsensiRepository) CheckIn(ctx context.Context, a *models.Absensi) error {
	query := `
		INSERT INTO absensi (pelaksanaan_id, tanggal, jam_masuk, keterangan, kegiatan)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, a.PelaksanaanID, a.Tanggal, a.JamMasuk, a.Keterangan, a.Kegiatan).
		Scan(&a.ID, &a.CreatedAt)
}

func (r *AbsensiRepository) CheckOut(ctx context.Context, pelaksanaanID uuid.UUID, tanggal, jamKeluar string) error {
	_, err := r.db.Exec(ctx,
		"UPDATE absensi SET jam_keluar=$1 WHERE pelaksanaan_id=$2 AND tanggal=$3",
		jamKeluar, pelaksanaanID, tanggal)
	return err
}

func (r *AbsensiRepository) FindByPelaksanaanID(ctx context.Context, pelaksanaanID uuid.UUID) ([]models.Absensi, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, pelaksanaan_id, tanggal, jam_masuk, jam_keluar, keterangan, kegiatan,
		 ttd_pembimbing, approved_by, approved_at, catatan, created_at
		 FROM absensi WHERE pelaksanaan_id=$1 ORDER BY tanggal ASC`, pelaksanaanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Absensi
	for rows.Next() {
		var a models.Absensi
		rows.Scan(&a.ID, &a.PelaksanaanID, &a.Tanggal, &a.JamMasuk, &a.JamKeluar, &a.Keterangan, &a.Kegiatan,
			&a.TTDPembimbing, &a.ApprovedBy, &a.ApprovedAt, &a.Catatan, &a.CreatedAt)
		list = append(list, a)
	}
	return list, nil
}

func (r *AbsensiRepository) FindByDate(ctx context.Context, pelaksanaanID uuid.UUID, tanggal string) (*models.Absensi, error) {
	a := &models.Absensi{}
	err := r.db.QueryRow(ctx,
		`SELECT id, pelaksanaan_id, tanggal, jam_masuk, jam_keluar, keterangan, kegiatan,
		 ttd_pembimbing, approved_by, approved_at, catatan, created_at
		 FROM absensi WHERE pelaksanaan_id=$1 AND tanggal=$2`, pelaksanaanID, tanggal).
		Scan(&a.ID, &a.PelaksanaanID, &a.Tanggal, &a.JamMasuk, &a.JamKeluar, &a.Keterangan, &a.Kegiatan,
			&a.TTDPembimbing, &a.ApprovedBy, &a.ApprovedAt, &a.Catatan, &a.CreatedAt)
	return a, err
}

func (r *AbsensiRepository) Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"UPDATE absensi SET ttd_pembimbing=true, approved_by=$1, approved_at=NOW() WHERE id=$2",
		approvedBy, id)
	return err
}

func (r *AbsensiRepository) CountByPelaksanaan(ctx context.Context, pelaksanaanID uuid.UUID) (hadir, izin, sakit, alpha int) {
	r.db.QueryRow(ctx,
		`SELECT
			COUNT(*) FILTER (WHERE keterangan='hadir'),
			COUNT(*) FILTER (WHERE keterangan='izin'),
			COUNT(*) FILTER (WHERE keterangan='sakit'),
			COUNT(*) FILTER (WHERE keterangan='alpha')
		 FROM absensi WHERE pelaksanaan_id=$1`, pelaksanaanID).
		Scan(&hadir, &izin, &sakit, &alpha)
	return
}
