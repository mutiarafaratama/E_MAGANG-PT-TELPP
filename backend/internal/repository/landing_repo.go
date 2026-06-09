package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/telpp/emagang/internal/models"
)

type LandingRepository struct {
	db *pgxpool.Pool
}

func NewLandingRepository(db *pgxpool.Pool) *LandingRepository {
	return &LandingRepository{db: db}
}

func (r *LandingRepository) GetAllContent(ctx context.Context) ([]models.LandingContent, error) {
	rows, err := r.db.Query(ctx,
		"SELECT id, kunci, nilai, tipe, updated_by, updated_at FROM landing_content ORDER BY kunci")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.LandingContent
	for rows.Next() {
		var c models.LandingContent
		rows.Scan(&c.ID, &c.Kunci, &c.Nilai, &c.Tipe, &c.UpdatedBy, &c.UpdatedAt)
		list = append(list, c)
	}
	return list, nil
}

func (r *LandingRepository) GetContentByKey(ctx context.Context, kunci string) (*models.LandingContent, error) {
	c := &models.LandingContent{}
	err := r.db.QueryRow(ctx,
		"SELECT id, kunci, nilai, tipe, updated_by, updated_at FROM landing_content WHERE kunci=$1", kunci).
		Scan(&c.ID, &c.Kunci, &c.Nilai, &c.Tipe, &c.UpdatedBy, &c.UpdatedAt)
	return c, err
}

func (r *LandingRepository) UpsertContent(ctx context.Context, kunci, nilai, tipe string, updatedBy uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO landing_content (kunci, nilai, tipe, updated_by, updated_at)
		VALUES ($1,$2,$3,$4,NOW())
		ON CONFLICT (kunci) DO UPDATE SET nilai=$2, tipe=$3, updated_by=$4, updated_at=NOW()`,
		kunci, nilai, tipe, updatedBy)
	return err
}

func (r *LandingRepository) GetAllFAQ(ctx context.Context, activeOnly bool) ([]models.FAQ, error) {
	where := ""
	if activeOnly {
		where = "WHERE is_active=true"
	}
	rows, err := r.db.Query(ctx,
		"SELECT id, pertanyaan, jawaban, urutan, is_active, created_at FROM faq "+where+" ORDER BY urutan ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.FAQ
	for rows.Next() {
		var f models.FAQ
		rows.Scan(&f.ID, &f.Pertanyaan, &f.Jawaban, &f.Urutan, &f.IsActive, &f.CreatedAt)
		list = append(list, f)
	}
	return list, nil
}

func (r *LandingRepository) UpsertFAQ(ctx context.Context, f *models.FAQ) error {
	if f.ID == uuid.Nil {
		return r.db.QueryRow(ctx,
			"INSERT INTO faq (pertanyaan, jawaban, urutan) VALUES ($1,$2,$3) RETURNING id, created_at",
			f.Pertanyaan, f.Jawaban, f.Urutan).Scan(&f.ID, &f.CreatedAt)
	}
	_, err := r.db.Exec(ctx,
		"UPDATE faq SET pertanyaan=$1, jawaban=$2, urutan=$3, is_active=$4 WHERE id=$5",
		f.Pertanyaan, f.Jawaban, f.Urutan, f.IsActive, f.ID)
	return err
}

func (r *LandingRepository) DeleteFAQ(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM faq WHERE id=$1", id)
	return err
}

func (r *LandingRepository) GetAllPeriode(ctx context.Context) ([]models.PeriodeMagang, error) {
	rows, err := r.db.Query(ctx,
		"SELECT id, nama, tanggal_buka, tanggal_tutup, kuota, is_active, created_at FROM periode_magang ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.PeriodeMagang
	for rows.Next() {
		var p models.PeriodeMagang
		rows.Scan(&p.ID, &p.Nama, &p.TanggalBuka, &p.TanggalTutup, &p.Kuota, &p.IsActive, &p.CreatedAt)
		list = append(list, p)
	}
	return list, nil
}

func (r *LandingRepository) GetActivePeriode(ctx context.Context) (*models.PeriodeMagang, error) {
	p := &models.PeriodeMagang{}
	err := r.db.QueryRow(ctx,
		"SELECT id, nama, tanggal_buka, tanggal_tutup, kuota, is_active, created_at FROM periode_magang WHERE is_active=true LIMIT 1").
		Scan(&p.ID, &p.Nama, &p.TanggalBuka, &p.TanggalTutup, &p.Kuota, &p.IsActive, &p.CreatedAt)
	return p, err
}

func (r *LandingRepository) CreatePeriode(ctx context.Context, p *models.PeriodeMagang) error {
	return r.db.QueryRow(ctx,
		"INSERT INTO periode_magang (nama, tanggal_buka, tanggal_tutup, kuota) VALUES ($1,$2,$3,$4) RETURNING id, created_at",
		p.Nama, p.TanggalBuka, p.TanggalTutup, p.Kuota).Scan(&p.ID, &p.CreatedAt)
}

func (r *LandingRepository) SetPeriodeActive(ctx context.Context, id uuid.UUID, active bool) error {
	if active {
		// Nonaktifkan semua dulu
		r.db.Exec(ctx, "UPDATE periode_magang SET is_active=false")
	}
	_, err := r.db.Exec(ctx, "UPDATE periode_magang SET is_active=$1 WHERE id=$2", active, id)
	return err
}
