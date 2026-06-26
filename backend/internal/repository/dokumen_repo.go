package repository

import (
        "context"
        "fmt"

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

// GetPathsByPengajuanID — ambil semua path_file untuk keperluan cleanup fisik sebelum hapus
func (r *DokumenRepository) GetPathsByPengajuanID(ctx context.Context, pengajuanID uuid.UUID) ([]string, error) {
        rows, err := r.db.Query(ctx,
                `SELECT path_file FROM dokumen WHERE pengajuan_id = $1`, pengajuanID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var paths []string
        for rows.Next() {
                var p string
                if err := rows.Scan(&p); err != nil {
                        return nil, err
                }
                paths = append(paths, p)
        }
        return paths, nil
}

// FindAll — admin: daftar semua dokumen dengan info pemilik, paginated + filter jenis + search nama
func (r *DokumenRepository) FindAll(ctx context.Context, jenis, search string, page, limit int) ([]models.DokumenWithUser, int, error) {
        if page < 1 {
                page = 1
        }
        if limit < 1 || limit > 100 {
                limit = 20
        }
        offset := (page - 1) * limit

        conditions := []string{}
        args := []any{}
        argIdx := 1

        if jenis != "" {
                conditions = append(conditions, fmt.Sprintf("d.jenis = $%d", argIdx))
                args = append(args, jenis)
                argIdx++
        }
        if search != "" {
                conditions = append(conditions, fmt.Sprintf(
                        "(u.nama_lengkap ILIKE $%d OR u.email ILIKE $%d OR p.nama_lengkap ILIKE $%d OR p.email ILIKE $%d)",
                        argIdx, argIdx, argIdx, argIdx,
                ))
                args = append(args, "%"+search+"%")
                argIdx++
        }

        where := ""
        if len(conditions) > 0 {
                where = "WHERE " + conditions[0]
                for _, c := range conditions[1:] {
                        where += " AND " + c
                }
        }

        fromClause := `FROM dokumen d
                LEFT JOIN users u ON u.id = d.user_id
                LEFT JOIN pengajuan_magang p ON p.id = d.pengajuan_id`

        var total int
        countQuery := "SELECT COUNT(*) " + fromClause + " " + where
        if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
                return nil, 0, err
        }

        args = append(args, limit, offset)
        query := `
                SELECT d.id, d.pengajuan_id, d.user_id, d.jenis, d.nama_file, d.path_file,
                       d.ukuran_bytes, d.mime_type, d.uploaded_at,
                       COALESCE(u.nama_lengkap, p.nama_lengkap) AS nama_pemilik,
                       COALESCE(u.email, p.email)               AS email_pemilik
                ` + fromClause + `
                ` + where + `
                ORDER BY d.uploaded_at DESC
                LIMIT $` + fmt.Sprintf("%d", argIdx) + ` OFFSET $` + fmt.Sprintf("%d", argIdx+1)

        rows, err := r.db.Query(ctx, query, args...)
        if err != nil {
                return nil, 0, err
        }
        defer rows.Close()

        var list []models.DokumenWithUser
        for rows.Next() {
                var dw models.DokumenWithUser
                if err := rows.Scan(
                        &dw.ID, &dw.PengajuanID, &dw.UserID, &dw.Jenis, &dw.NamaFile, &dw.PathFile,
                        &dw.UkuranBytes, &dw.MimeType, &dw.UploadedAt,
                        &dw.NamaPemilik, &dw.EmailPemilik,
                ); err != nil {
                        return nil, 0, err
                }
                list = append(list, dw)
        }
        return list, total, nil
}

// FindSuratBalasanPath — ambil path_file surat_balasan milik pengajuan tertentu
func (r *DokumenRepository) FindSuratBalasanPath(ctx context.Context, pengajuanID uuid.UUID) string {
        var path string
        r.db.QueryRow(ctx,
                `SELECT path_file FROM dokumen WHERE pengajuan_id = $1 AND jenis = 'surat_balasan' ORDER BY uploaded_at DESC LIMIT 1`,
                pengajuanID,
        ).Scan(&path)
        return path
}

// SavePublik — simpan dokumen dari form publik (user_id NULL, belum punya akun)
func (r *DokumenRepository) SavePublik(ctx context.Context, d *models.Dokumen) error {
        query := `
                INSERT INTO dokumen (pengajuan_id, jenis, nama_file, path_file, ukuran_bytes, mime_type)
                VALUES ($1, $2, $3, $4, $5, $6)
                RETURNING id, uploaded_at`
        return r.db.QueryRow(ctx, query,
                d.PengajuanID, d.Jenis, d.NamaFile, d.PathFile, d.UkuranBytes, d.MimeType,
        ).Scan(&d.ID, &d.UploadedAt)
}
