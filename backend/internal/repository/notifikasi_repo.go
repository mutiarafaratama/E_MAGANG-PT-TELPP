package repository

import (
        "context"

        "github.com/google/uuid"
        "github.com/jackc/pgx/v5/pgxpool"
        "github.com/telpp/emagang/internal/models"
)

type NotifikasiRepository struct {
        db *pgxpool.Pool
}

func NewNotifikasiRepository(db *pgxpool.Pool) *NotifikasiRepository {
        return &NotifikasiRepository{db: db}
}

func (r *NotifikasiRepository) Create(ctx context.Context, n *models.Notifikasi) error {
        query := `INSERT INTO notifikasi (user_id, judul, pesan, tipe, referensi_id)
                          VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at`
        return r.db.QueryRow(ctx, query, n.UserID, n.Judul, n.Pesan, n.Tipe, n.ReferensiID).
                Scan(&n.ID, &n.CreatedAt)
}

func (r *NotifikasiRepository) FindByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Notifikasi, int, error) {
        offset := (page - 1) * limit
        var total int
        r.db.QueryRow(ctx, "SELECT COUNT(*) FROM notifikasi WHERE user_id=$1", userID).Scan(&total)

        rows, err := r.db.Query(ctx,
                `SELECT id, user_id, judul, pesan, tipe, referensi_id, is_read, created_at
                 FROM notifikasi WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
                userID, limit, offset)
        if err != nil {
                return nil, 0, err
        }
        defer rows.Close()

        var list []models.Notifikasi
        for rows.Next() {
                var n models.Notifikasi
                rows.Scan(&n.ID, &n.UserID, &n.Judul, &n.Pesan, &n.Tipe, &n.ReferensiID, &n.IsRead, &n.CreatedAt)
                list = append(list, n)
        }
        return list, total, nil
}

func (r *NotifikasiRepository) MarkRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
        _, err := r.db.Exec(ctx,
                "UPDATE notifikasi SET is_read=true WHERE id=$1 AND user_id=$2", id, userID)
        return err
}

func (r *NotifikasiRepository) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
        _, err := r.db.Exec(ctx,
                "UPDATE notifikasi SET is_read=true WHERE user_id=$1 AND is_read=false", userID)
        return err
}

func (r *NotifikasiRepository) CountUnread(ctx context.Context, userID uuid.UUID) int {
        var count int
        r.db.QueryRow(ctx, "SELECT COUNT(*) FROM notifikasi WHERE user_id=$1 AND is_read=false", userID).Scan(&count)
        return count
}

// SaveFCMToken — simpan atau update FCM token milik user
func (r *NotifikasiRepository) SaveFCMToken(ctx context.Context, userID uuid.UUID, token string) error {
        _, err := r.db.Exec(ctx,
                `UPDATE users SET fcm_token=$1, fcm_updated_at=NOW() WHERE id=$2`,
                token, userID)
        return err
}

// GetFCMToken — ambil FCM token milik user (kosong jika belum ada)
func (r *NotifikasiRepository) GetFCMToken(ctx context.Context, userID uuid.UUID) string {
        var token string
        r.db.QueryRow(ctx, `SELECT COALESCE(fcm_token,'') FROM users WHERE id=$1`, userID).Scan(&token)
        return token
}
