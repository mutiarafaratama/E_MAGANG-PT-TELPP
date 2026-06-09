package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/telpp/emagang/internal/models"
)

type DokumenRepository struct {
	db *pgxpool.Pool
}

func NewDokumenRepository(db *pgxpool.Pool) *DokumenRepository {
	return &DokumenRepository{db: db}
}

func (r *DokumenRepository) Save(ctx context.Context, d *models.Dokumen) error {
	query := `
		INSERT INTO dokumen (pengajuan_id, user_id, jenis, nama_file, path_file, ukuran_bytes, mime_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, uploaded_at`
	return r.db.QueryRow(ctx, query,
		d.PengajuanID, d.UserID, d.Jenis, d.NamaFile, d.PathFile, d.UkuranBytes, d.MimeType,
	).Scan(&d.ID, &d.UploadedAt)
}

func (r *DokumenRepository) FindByPengajuanID(ctx context.Context, pengajuanID uuid.UUID) ([]models.Dokumen, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, pengajuan_id, user_id, jenis, nama_file, path_file, ukuran_bytes, mime_type, uploaded_at
		 FROM dokumen WHERE pengajuan_id = $1 ORDER BY uploaded_at ASC`, pengajuanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Dokumen
	for rows.Next() {
		var d models.Dokumen
		if err := rows.Scan(&d.ID, &d.PengajuanID, &d.UserID, &d.Jenis, &d.NamaFile, &d.PathFile, &d.UkuranBytes, &d.MimeType, &d.UploadedAt); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *DokumenRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Dokumen, error) {
	d := &models.Dokumen{}
	err := r.db.QueryRow(ctx,
		`SELECT id, pengajuan_id, user_id, jenis, nama_file, path_file, ukuran_bytes, mime_type, uploaded_at
		 FROM dokumen WHERE id = $1`, id).
		Scan(&d.ID, &d.PengajuanID, &d.UserID, &d.Jenis, &d.NamaFile, &d.PathFile, &d.UkuranBytes, &d.MimeType, &d.UploadedAt)
	return d, err
}

func (r *DokumenRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM dokumen WHERE id = $1", id)
	return err
}
