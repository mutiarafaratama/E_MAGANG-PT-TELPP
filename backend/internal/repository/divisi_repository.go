package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/telpp/emagang/internal/models"
)

type DivisiRepository struct {
	db *pgxpool.Pool
}

func NewDivisiRepository(db *pgxpool.Pool) *DivisiRepository {
	return &DivisiRepository{db: db}
}

const divisiColumns = `id, nama, is_active, urutan, geo_lat, geo_lng, geo_radius, nama_lokasi, created_at`

func scanDivisi(row interface {
	Scan(...interface{}) error
}, d *models.Divisi) error {
	return row.Scan(&d.ID, &d.Nama, &d.IsActive, &d.Urutan,
		&d.GeoLat, &d.GeoLng, &d.GeoRadius, &d.NamaLokasi, &d.CreatedAt)
}

// FindAll — semua divisi, aktif saja jika onlyActive=true
func (r *DivisiRepository) FindAll(ctx context.Context, onlyActive bool) ([]models.Divisi, error) {
	query := `SELECT ` + divisiColumns + ` FROM divisi`
	if onlyActive {
		query += ` WHERE is_active = TRUE`
	}
	query += ` ORDER BY urutan ASC, nama ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Divisi
	for rows.Next() {
		var d models.Divisi
		if err := scanDivisi(rows, &d); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	if list == nil {
		list = []models.Divisi{}
	}
	return list, nil
}

// FindByNama — ambil satu divisi berdasarkan nama (untuk lookup geofence saat absensi)
func (r *DivisiRepository) FindByNama(ctx context.Context, nama string) (*models.Divisi, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+divisiColumns+` FROM divisi WHERE nama = $1 LIMIT 1`, nama)
	var d models.Divisi
	if err := scanDivisi(row, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// Create — buat divisi baru
func (r *DivisiRepository) Create(ctx context.Context, nama string, urutan int, geoLat, geoLng *float64, geoRadius *int, namaLokasi *string) (*models.Divisi, error) {
	d := &models.Divisi{
		ID:         uuid.New(),
		Nama:       nama,
		IsActive:   true,
		Urutan:     urutan,
		GeoLat:     geoLat,
		GeoLng:     geoLng,
		GeoRadius:  geoRadius,
		NamaLokasi: namaLokasi,
		CreatedAt:  time.Now(),
	}
	_, err := r.db.Exec(ctx,
		`INSERT INTO divisi (id, nama, is_active, urutan, geo_lat, geo_lng, geo_radius, nama_lokasi, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		d.ID, d.Nama, d.IsActive, d.Urutan, d.GeoLat, d.GeoLng, d.GeoRadius, d.NamaLokasi, d.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Update — ubah nama, urutan, dan geofencing
func (r *DivisiRepository) Update(ctx context.Context, id uuid.UUID, nama string, urutan int, geoLat, geoLng *float64, geoRadius *int, namaLokasi *string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE divisi SET nama=$1, urutan=$2, geo_lat=$3, geo_lng=$4, geo_radius=$5, nama_lokasi=$6 WHERE id=$7`,
		nama, urutan, geoLat, geoLng, geoRadius, namaLokasi, id,
	)
	return err
}

// ToggleActive — aktifkan/nonaktifkan
func (r *DivisiRepository) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	_, err := r.db.Exec(ctx,
		`UPDATE divisi SET is_active=$1 WHERE id=$2`,
		isActive, id,
	)
	return err
}

// Delete — hapus permanen
func (r *DivisiRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM divisi WHERE id=$1`, id)
	return err
}
