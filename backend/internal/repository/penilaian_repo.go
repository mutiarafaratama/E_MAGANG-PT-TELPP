package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/telpp/emagang/internal/models"
)

type PenilaianRepository struct {
	db *pgxpool.Pool
}

func NewPenilaianRepository(db *pgxpool.Pool) *PenilaianRepository {
	return &PenilaianRepository{db: db}
}

// GetByPelaksanaan mengambil penilaian beserta komponen kejuruan-nya.
func (r *PenilaianRepository) GetByPelaksanaan(ctx context.Context, pelaksanaanID uuid.UUID) (*models.Penilaian, error) {
	p := &models.Penilaian{}
	err := r.db.QueryRow(ctx, `
		SELECT id, pelaksanaan_id,
		       nilai_motivasi, nilai_inisiatif, nilai_disiplin_waktu,
		       nilai_kerajinan, nilai_kreativitas, nilai_tanggung_jawab,
		       nilai_kerjasama, nilai_adaptasi, nilai_kehadiran,
		       nilai_k3_safety, nilai_k3_metode, nilai_k3_manajemen, nilai_k3_volume,
		       nilai_prs_proses, nilai_prs_teori, nilai_prs_judul, nilai_prs_data,
		       manager_nama, catatan, nilai_akhir,
		       dinilai_oleh, dinilai_at, created_at, updated_at
		FROM penilaian WHERE pelaksanaan_id = $1`, pelaksanaanID).
		Scan(&p.ID, &p.PelaksanaanID,
			&p.NilaiMotivasi, &p.NilaiInisiatif, &p.NilaiDisiplinWaktu,
			&p.NilaiKerajinan, &p.NilaiKreativitas, &p.NilaiTanggungJawab,
			&p.NilaiKerjasama, &p.NilaiAdaptasi, &p.NilaiKehadiran,
			&p.NilaiK3Safety, &p.NilaiK3Metode, &p.NilaiK3Manajemen, &p.NilaiK3Volume,
			&p.NilaiPrsProses, &p.NilaiPrsTeori, &p.NilaiPrsJudul, &p.NilaiPrsData,
			&p.ManagerNama, &p.Catatan, &p.NilaiAkhir,
			&p.DinilaiOleh, &p.DinilaiAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, penilaian_id, komponen, nilai, urutan
		FROM penilaian_kejuruan WHERE penilaian_id = $1 ORDER BY urutan ASC`, p.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var k models.PenilaianKejuruan
		if err := rows.Scan(&k.ID, &k.PenilaianID, &k.Komponen, &k.Nilai, &k.Urutan); err != nil {
			continue
		}
		p.Kejuruan = append(p.Kejuruan, k)
	}
	return p, nil
}

// Upsert menyimpan atau memperbarui penilaian beserta kejuruan-nya (dalam transaksi).
func (r *PenilaianRepository) Upsert(ctx context.Context, req models.PenilaianUpsertRequest, dinilaiOleh uuid.UUID) (*models.Penilaian, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var penilaianID uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO penilaian (
		  pelaksanaan_id,
		  nilai_motivasi, nilai_inisiatif, nilai_disiplin_waktu,
		  nilai_kerajinan, nilai_kreativitas, nilai_tanggung_jawab,
		  nilai_kerjasama, nilai_adaptasi, nilai_kehadiran,
		  nilai_k3_safety, nilai_k3_metode, nilai_k3_manajemen, nilai_k3_volume,
		  nilai_prs_proses, nilai_prs_teori, nilai_prs_judul, nilai_prs_data,
		  manager_nama, catatan, nilai_akhir,
		  dinilai_oleh, dinilai_at, updated_at
		) VALUES (
		  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,NOW(),NOW()
		)
		ON CONFLICT (pelaksanaan_id) DO UPDATE SET
		  nilai_motivasi        = EXCLUDED.nilai_motivasi,
		  nilai_inisiatif       = EXCLUDED.nilai_inisiatif,
		  nilai_disiplin_waktu  = EXCLUDED.nilai_disiplin_waktu,
		  nilai_kerajinan       = EXCLUDED.nilai_kerajinan,
		  nilai_kreativitas     = EXCLUDED.nilai_kreativitas,
		  nilai_tanggung_jawab  = EXCLUDED.nilai_tanggung_jawab,
		  nilai_kerjasama       = EXCLUDED.nilai_kerjasama,
		  nilai_adaptasi        = EXCLUDED.nilai_adaptasi,
		  nilai_kehadiran       = EXCLUDED.nilai_kehadiran,
		  nilai_k3_safety       = EXCLUDED.nilai_k3_safety,
		  nilai_k3_metode       = EXCLUDED.nilai_k3_metode,
		  nilai_k3_manajemen    = EXCLUDED.nilai_k3_manajemen,
		  nilai_k3_volume       = EXCLUDED.nilai_k3_volume,
		  nilai_prs_proses      = EXCLUDED.nilai_prs_proses,
		  nilai_prs_teori       = EXCLUDED.nilai_prs_teori,
		  nilai_prs_judul       = EXCLUDED.nilai_prs_judul,
		  nilai_prs_data        = EXCLUDED.nilai_prs_data,
		  manager_nama          = EXCLUDED.manager_nama,
		  catatan               = EXCLUDED.catatan,
		  nilai_akhir           = EXCLUDED.nilai_akhir,
		  dinilai_oleh          = EXCLUDED.dinilai_oleh,
		  dinilai_at            = NOW(),
		  updated_at            = NOW()
		RETURNING id`,
		req.PelaksanaanID,
		nilFloat(req.NilaiMotivasi), nilFloat(req.NilaiInisiatif), nilFloat(req.NilaiDisiplinWaktu),
		nilFloat(req.NilaiKerajinan), nilFloat(req.NilaiKreativitas), nilFloat(req.NilaiTanggungJawab),
		nilFloat(req.NilaiKerjasama), nilFloat(req.NilaiAdaptasi), nilFloat(req.NilaiKehadiran),
		nilFloat(req.NilaiK3Safety), nilFloat(req.NilaiK3Metode), nilFloat(req.NilaiK3Manajemen), nilFloat(req.NilaiK3Volume),
		nilFloat(req.NilaiPrsProses), nilFloat(req.NilaiPrsTeori), nilFloat(req.NilaiPrsJudul), nilFloat(req.NilaiPrsData),
		nilString(req.ManagerNama), nilString(req.Catatan), req.NilaiAkhir, dinilaiOleh,
	).Scan(&penilaianID)
	if err != nil {
		return nil, err
	}

	// Hapus kejuruan lama, insert baru
	if _, err = tx.Exec(ctx, `DELETE FROM penilaian_kejuruan WHERE penilaian_id = $1`, penilaianID); err != nil {
		return nil, err
	}
	for i, k := range req.Kejuruan {
		if k.Komponen == "" {
			continue
		}
		if _, err = tx.Exec(ctx,
			`INSERT INTO penilaian_kejuruan (penilaian_id, komponen, nilai, urutan) VALUES ($1,$2,$3,$4)`,
			penilaianID, k.Komponen, k.Nilai, i); err != nil {
			return nil, err
		}
	}

	// Sync nilai_akhir ke pelaksanaan_magang
	if _, err = tx.Exec(ctx,
		`UPDATE pelaksanaan_magang SET nilai=$1, catatan_nilai=$2, dinilai_oleh=$3, dinilai_at=NOW(), updated_at=NOW() WHERE id=$4`,
		req.NilaiAkhir, nilString(req.Catatan), dinilaiOleh, req.PelaksanaanID); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetByPelaksanaan(ctx, req.PelaksanaanID)
}

// ListPerluDinilai mengembalikan peserta yang sudah di tahap penilaian atau selesai.
func (r *PenilaianRepository) ListPerluDinilai(ctx context.Context) ([]models.PenilaianListItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT pm.id AS pelaksanaan_id, pm.status,
		       u.nama_lengkap, u.email,
		       COALESCE(pm.divisi,'') AS divisi,
		       COALESCE(pm.pembimbing, u2.nama_lengkap, '') AS pembimbing,
		       pm.tanggal_mulai, pm.tanggal_selesai,
		       p.id AS penilaian_id, p.nilai_akhir, p.dinilai_at
		FROM pelaksanaan_magang pm
		JOIN users u ON u.id = pm.user_id
		LEFT JOIN users u2 ON u2.id = pm.pembimbing_id
		LEFT JOIN penilaian p ON p.pelaksanaan_id = pm.id
		WHERE pm.status IN ('penilaian','selesai')
		ORDER BY pm.updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.PenilaianListItem
	for rows.Next() {
		var item models.PenilaianListItem
		if err := rows.Scan(
			&item.PelaksanaanID, &item.Status,
			&item.NamaLengkap, &item.Email,
			&item.Divisi, &item.Pembimbing,
			&item.TanggalMulai, &item.TanggalSelesai,
			&item.PenilaianID, &item.NilaiAkhir, &item.DinilaiAt,
		); err != nil {
			continue
		}
		list = append(list, item)
	}
	return list, nil
}

// helpers
func nilFloat(v float64) *float64 {
	if v == 0 {
		return nil
	}
	return &v
}
func nilString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
