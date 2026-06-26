package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/telpp/emagang/internal/models"
)

type AbsensiConfigRepository struct {
	db *pgxpool.Pool
}

func NewAbsensiConfigRepository(db *pgxpool.Pool) *AbsensiConfigRepository {
	return &AbsensiConfigRepository{db: db}
}

func (r *AbsensiConfigRepository) Get(ctx context.Context) (*models.AbsensiConfig, error) {
	cfg := &models.AbsensiConfig{}
	err := r.db.QueryRow(ctx,
		`SELECT id,
		        to_char(jam_masuk_buka,  'HH24:MI'),
		        to_char(jam_masuk_tutup, 'HH24:MI'),
		        to_char(jam_pulang_buka,  'HH24:MI'),
		        to_char(jam_pulang_tutup, 'HH24:MI'),
		        updated_at, updated_by
		 FROM absensi_config ORDER BY id LIMIT 1`).
		Scan(&cfg.ID, &cfg.JamMasukBuka, &cfg.JamMasukTutup,
			&cfg.JamPulangBuka, &cfg.JamPulangTutup,
			&cfg.UpdatedAt, &cfg.UpdatedBy)
	return cfg, err
}

func (r *AbsensiConfigRepository) Update(ctx context.Context,
	jamMasukBuka, jamMasukTutup, jamPulangBuka, jamPulangTutup string,
	updatedBy uuid.UUID,
) error {
	_, err := r.db.Exec(ctx,
		`UPDATE absensi_config
		 SET jam_masuk_buka=$1, jam_masuk_tutup=$2,
		     jam_pulang_buka=$3, jam_pulang_tutup=$4,
		     updated_at=NOW(), updated_by=$5
		 WHERE id=1`,
		jamMasukBuka, jamMasukTutup, jamPulangBuka, jamPulangTutup, updatedBy)
	return err
}
