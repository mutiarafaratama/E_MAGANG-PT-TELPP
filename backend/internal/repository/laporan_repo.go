package repository

import (
        "context"
        "fmt"
        "io"
        "mime/multipart"
        "os"
        "path/filepath"
        "time"

        "github.com/google/uuid"
        "github.com/jackc/pgx/v5"
        "github.com/jackc/pgx/v5/pgxpool"
        "github.com/telpp/emagang/internal/config"
        "github.com/telpp/emagang/internal/models"
)

type LaporanRepository struct {
        db *pgxpool.Pool
}

func NewLaporanRepository(db *pgxpool.Pool) *LaporanRepository {
        return &LaporanRepository{db: db}
}

// saveFile — simpan multipart file ke disk, kembalikan path absolut dan nama asli
func saveFile(fh *multipart.FileHeader) (absPath string, namaFile string, err error) {
        src, err := fh.Open()
        if err != nil {
                return "", "", fmt.Errorf("buka file: %w", err)
        }
        defer src.Close()

        dir := filepath.Join(config.App.UploadDir, "laporan")
        if err := os.MkdirAll(dir, 0755); err != nil {
                return "", "", fmt.Errorf("buat direktori: %w", err)
        }

        ext := filepath.Ext(fh.Filename)
        fname := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
        dst := filepath.Join(dir, fname)

        out, err := os.Create(dst)
        if err != nil {
                return "", "", fmt.Errorf("buat file: %w", err)
        }
        defer out.Close()

        if _, err := io.Copy(out, src); err != nil {
                return "", "", fmt.Errorf("tulis file: %w", err)
        }

        return dst, fh.Filename, nil
}

// Upload — simpan file ke disk + buat record laporan baru (versi auto-increment)
func (r *LaporanRepository) Upload(ctx context.Context, pelaksanaanID uuid.UUID, fh *multipart.FileHeader) (*models.LaporanMagang, error) {
        dst, namaFile, err := saveFile(fh)
        if err != nil {
                return nil, err
        }

        var lastVersi int
        r.db.QueryRow(ctx,
                "SELECT COALESCE(MAX(versi),0) FROM laporan_magang WHERE pelaksanaan_id=$1",
                pelaksanaanID).Scan(&lastVersi)

        sz := fh.Size
        mime := fh.Header.Get("Content-Type")

        laporan := &models.LaporanMagang{
                ID:            uuid.New(),
                PelaksanaanID: pelaksanaanID,
                Versi:         lastVersi + 1,
                NamaFile:      namaFile,
                PathFile:      dst,
                UkuranBytes:   &sz,
                Status:        models.StatusLaporanMenunggu,
                DiuploadAt:    time.Now(),
        }
        if mime != "" {
                laporan.MimeType = &mime
        }

        _, err = r.db.Exec(ctx, `
                INSERT INTO laporan_magang
                        (id, pelaksanaan_id, versi, nama_file, path_file, ukuran_bytes, mime_type, status, diupload_at)
                VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
                laporan.ID, laporan.PelaksanaanID, laporan.Versi, laporan.NamaFile,
                laporan.PathFile, laporan.UkuranBytes, laporan.MimeType, laporan.Status, laporan.DiuploadAt,
        )
        if err != nil {
                return nil, err
        }
        return laporan, nil
}

// FindByPelaksanaan — semua laporan untuk satu pelaksanaan, terbaru duluan
func (r *LaporanRepository) FindByPelaksanaan(ctx context.Context, pelaksanaanID uuid.UUID) ([]*models.LaporanMagang, error) {
        rows, err := r.db.Query(ctx, `
                SELECT id, pelaksanaan_id, versi, nama_file, path_file, ukuran_bytes, mime_type,
                       status, catatan_hrd, direview_oleh, direview_at, diupload_at
                FROM laporan_magang
                WHERE pelaksanaan_id = $1
                ORDER BY versi DESC`, pelaksanaanID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var list []*models.LaporanMagang
        for rows.Next() {
                l := &models.LaporanMagang{}
                if err := rows.Scan(
                        &l.ID, &l.PelaksanaanID, &l.Versi, &l.NamaFile, &l.PathFile,
                        &l.UkuranBytes, &l.MimeType, &l.Status, &l.CatatanHRD,
                        &l.DireviewOleh, &l.DireviewAt, &l.DiuploadAt,
                ); err != nil {
                        return nil, err
                }
                list = append(list, l)
        }
        return list, nil
}

// FindByID — ambil satu laporan berdasarkan ID
func (r *LaporanRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.LaporanMagang, error) {
        l := &models.LaporanMagang{}
        err := r.db.QueryRow(ctx, `
                SELECT id, pelaksanaan_id, versi, nama_file, path_file, ukuran_bytes, mime_type,
                       status, catatan_hrd, direview_oleh, direview_at, diupload_at
                FROM laporan_magang WHERE id = $1`, id).
                Scan(&l.ID, &l.PelaksanaanID, &l.Versi, &l.NamaFile, &l.PathFile,
                        &l.UkuranBytes, &l.MimeType, &l.Status, &l.CatatanHRD,
                        &l.DireviewOleh, &l.DireviewAt, &l.DiuploadAt)
        if err != nil {
                return nil, err
        }
        return l, nil
}

// FindAll — semua laporan dengan filter status opsional, untuk halaman HRD
func (r *LaporanRepository) FindAll(ctx context.Context, status string, page, limit int) ([]*models.LaporanMagangDetail, int, error) {
        offset := (page - 1) * limit
        var total int
        if status == "" {
                r.db.QueryRow(ctx, "SELECT COUNT(*) FROM laporan_magang").Scan(&total)
        } else {
                r.db.QueryRow(ctx, "SELECT COUNT(*) FROM laporan_magang WHERE status=$1", status).Scan(&total)
        }

        var rows pgx.Rows
        var err error
        base := `
                SELECT lm.id, lm.pelaksanaan_id, lm.versi, lm.nama_file, lm.path_file,
                       lm.ukuran_bytes, lm.mime_type, lm.status, lm.catatan_hrd,
                       lm.direview_oleh, lm.direview_at, lm.diupload_at,
                       pg.nama_lengkap, pg.asal_institusi, pm.divisi
                FROM laporan_magang lm
                JOIN pelaksanaan_magang pm ON pm.id = lm.pelaksanaan_id
                JOIN pengajuan_magang pg   ON pg.id = pm.pengajuan_id`

        if status == "" {
                rows, err = r.db.Query(ctx, base+` ORDER BY lm.diupload_at DESC LIMIT $1 OFFSET $2`, limit, offset)
        } else {
                rows, err = r.db.Query(ctx, base+` WHERE lm.status = $1 ORDER BY lm.diupload_at DESC LIMIT $2 OFFSET $3`, status, limit, offset)
        }
        if err != nil {
                return nil, 0, err
        }
        defer rows.Close()

        var list []*models.LaporanMagangDetail
        for rows.Next() {
                d := &models.LaporanMagangDetail{}
                if err := rows.Scan(
                        &d.ID, &d.PelaksanaanID, &d.Versi, &d.NamaFile, &d.PathFile,
                        &d.UkuranBytes, &d.MimeType, &d.Status, &d.CatatanHRD,
                        &d.DireviewOleh, &d.DireviewAt, &d.DiuploadAt,
                        &d.NamaPeserta, &d.AsalInstitusi, &d.Divisi,
                ); err != nil {
                        return nil, 0, err
                }
                list = append(list, d)
        }
        return list, total, nil
}

// FindAllMenunggu — semua laporan menunggu review, untuk halaman HRD
func (r *LaporanRepository) FindAllMenunggu(ctx context.Context, page, limit int) ([]*models.LaporanMagangDetail, int, error) {
        offset := (page - 1) * limit
        var total int
        r.db.QueryRow(ctx,
                "SELECT COUNT(*) FROM laporan_magang WHERE status='menunggu_review'").Scan(&total)

        rows, err := r.db.Query(ctx, `
                SELECT lm.id, lm.pelaksanaan_id, lm.versi, lm.nama_file, lm.path_file,
                       lm.ukuran_bytes, lm.mime_type, lm.status, lm.catatan_hrd,
                       lm.direview_oleh, lm.direview_at, lm.diupload_at,
                       pg.nama_lengkap, pg.asal_institusi, pm.divisi
                FROM laporan_magang lm
                JOIN pelaksanaan_magang pm ON pm.id = lm.pelaksanaan_id
                JOIN pengajuan_magang pg   ON pg.id = pm.pengajuan_id
                WHERE lm.status = 'menunggu_review'
                ORDER BY lm.diupload_at ASC
                LIMIT $1 OFFSET $2`, limit, offset)
        if err != nil {
                return nil, 0, err
        }
        defer rows.Close()

        var list []*models.LaporanMagangDetail
        for rows.Next() {
                d := &models.LaporanMagangDetail{}
                if err := rows.Scan(
                        &d.ID, &d.PelaksanaanID, &d.Versi, &d.NamaFile, &d.PathFile,
                        &d.UkuranBytes, &d.MimeType, &d.Status, &d.CatatanHRD,
                        &d.DireviewOleh, &d.DireviewAt, &d.DiuploadAt,
                        &d.NamaPeserta, &d.AsalInstitusi, &d.Divisi,
                ); err != nil {
                        return nil, 0, err
                }
                list = append(list, d)
        }
        return list, total, nil
}

// Review — HRD acc atau minta revisi laporan
func (r *LaporanRepository) Review(ctx context.Context, id uuid.UUID, status models.StatusLaporan, catatan string, reviewerID uuid.UUID) error {
        now := time.Now()
        var cat *string
        if catatan != "" {
                cat = &catatan
        }
        _, err := r.db.Exec(ctx, `
                UPDATE laporan_magang
                SET status=$1, catatan_hrd=$2, direview_oleh=$3, direview_at=$4
                WHERE id=$5`,
                status, cat, reviewerID, now, id)
        return err
}
