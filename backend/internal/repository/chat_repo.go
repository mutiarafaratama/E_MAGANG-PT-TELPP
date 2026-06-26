package repository

import (
        "context"
        "fmt"
        "time"

        "github.com/google/uuid"
        "github.com/jackc/pgx/v5/pgxpool"
        "github.com/telpp/emagang/internal/models"
)

type ChatRepository struct {
        db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *ChatRepository {
        return &ChatRepository{db: db}
}

func (r *ChatRepository) GenerateNomorTiket(ctx context.Context) string {
        var count int
        today := time.Now().Format("20060102")
        r.db.QueryRow(ctx,
                "SELECT COUNT(*) FROM chat_tiket WHERE DATE(created_at) = CURRENT_DATE").Scan(&count)
        return fmt.Sprintf("TKT-%s-%03d", today, count+1)
}

// FindAktifByUserID — cek apakah user masih punya tiket aktif (status != selesai)
func (r *ChatRepository) FindAktifByUserID(ctx context.Context, userID uuid.UUID) (*models.ChatTiket, error) {
        t := &models.ChatTiket{}
        err := r.db.QueryRow(ctx,
                `SELECT id, user_id, nomor_tiket, subjek, status, assigned_to, kategori, created_at, updated_at
                 FROM chat_tiket WHERE user_id=$1 AND status != 'selesai' LIMIT 1`, userID).
                Scan(&t.ID, &t.UserID, &t.NomorTiket, &t.Subjek, &t.Status, &t.AssignedTo, &t.Kategori, &t.CreatedAt, &t.UpdatedAt)
        if err != nil {
                return nil, err
        }
        return t, nil
}

func (r *ChatRepository) CreateTiket(ctx context.Context, t *models.ChatTiket) error {
        query := `INSERT INTO chat_tiket (user_id, nomor_tiket, subjek, kategori) VALUES ($1,$2,$3,$4)
          RETURNING id, status, created_at, updated_at`
        return r.db.QueryRow(ctx, query, t.UserID, t.NomorTiket, t.Subjek, t.Kategori).
                Scan(&t.ID, &t.Status, &t.CreatedAt, &t.UpdatedAt)
}

func (r *ChatRepository) FindTiketByID(ctx context.Context, id uuid.UUID) (*models.ChatTiketDetail, error) {
        t := &models.ChatTiketDetail{}
        err := r.db.QueryRow(ctx, `
SELECT t.id, t.user_id, t.nomor_tiket, t.subjek, t.status, t.assigned_to, t.kategori, t.created_at, t.updated_at,
       u.nama_lengkap, u.email,
       a.nama_lengkap,
       (SELECT COUNT(*) FROM chat_pesan WHERE tiket_id=t.id AND is_read=false AND sender_id!=t.user_id),
       (SELECT pesan FROM chat_pesan WHERE tiket_id=t.id ORDER BY created_at DESC LIMIT 1)
FROM chat_tiket t
JOIN users u ON u.id = t.user_id
LEFT JOIN users a ON a.id = t.assigned_to
WHERE t.id = $1`, id).
                Scan(&t.ID, &t.UserID, &t.NomorTiket, &t.Subjek, &t.Status, &t.AssignedTo, &t.Kategori, &t.CreatedAt, &t.UpdatedAt,
                        &t.UserNama, &t.UserEmail, &t.AssignedNama, &t.UnreadCount, &t.LastMessage)
        return t, err
}

func (r *ChatRepository) FindTiketByUserID(ctx context.Context, userID uuid.UUID) ([]models.ChatTiketDetail, error) {
        rows, err := r.db.Query(ctx, `
SELECT t.id, t.user_id, t.nomor_tiket, t.subjek, t.status, t.assigned_to, t.kategori, t.created_at, t.updated_at,
       u.nama_lengkap, u.email,
       a.nama_lengkap,
       (SELECT COUNT(*) FROM chat_pesan WHERE tiket_id=t.id AND is_read=false AND sender_id!=t.user_id),
       (SELECT pesan FROM chat_pesan WHERE tiket_id=t.id ORDER BY created_at DESC LIMIT 1)
FROM chat_tiket t
JOIN users u ON u.id = t.user_id
LEFT JOIN users a ON a.id = t.assigned_to
WHERE t.user_id = $1 ORDER BY t.updated_at DESC`, userID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()
        return r.scanTiketRows(rows)
}

func (r *ChatRepository) FindAllTiket(ctx context.Context, status, kategori string, page, limit int) ([]models.ChatTiketDetail, int, error) {
        offset := (page - 1) * limit
        where := "WHERE 1=1"
        var args []interface{}
        argIdx := 1
        if status != "" {
                where += fmt.Sprintf(" AND t.status = $%d", argIdx)
                args = append(args, status)
                argIdx++
        }
        if kategori != "" {
                where += fmt.Sprintf(" AND t.kategori = $%d", argIdx)
                args = append(args, kategori)
                argIdx++
        }
        var total int
        r.db.QueryRow(ctx, "SELECT COUNT(*) FROM chat_tiket t "+where, args...).Scan(&total)
        args = append(args, limit, offset)

        rows, err := r.db.Query(ctx, fmt.Sprintf(`
SELECT t.id, t.user_id, t.nomor_tiket, t.subjek, t.status, t.assigned_to, t.kategori, t.created_at, t.updated_at,
       u.nama_lengkap, u.email,
       a.nama_lengkap,
       (SELECT COUNT(*) FROM chat_pesan WHERE tiket_id=t.id AND is_read=false AND sender_id=t.user_id),
       (SELECT pesan FROM chat_pesan WHERE tiket_id=t.id ORDER BY created_at DESC LIMIT 1)
FROM chat_tiket t
JOIN users u ON u.id = t.user_id
LEFT JOIN users a ON a.id = t.assigned_to
%s ORDER BY t.updated_at DESC LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1), args...)
        if err != nil {
                return nil, 0, err
        }
        defer rows.Close()
        list, err := r.scanTiketRows(rows)
        return list, total, err
}

func (r *ChatRepository) scanTiketRows(rows interface {
        Next() bool
        Scan(...interface{}) error
}) ([]models.ChatTiketDetail, error) {
        var list []models.ChatTiketDetail
        for rows.Next() {
                var t models.ChatTiketDetail
                if err := rows.Scan(
                        &t.ID, &t.UserID, &t.NomorTiket, &t.Subjek, &t.Status, &t.AssignedTo, &t.Kategori,
                        &t.CreatedAt, &t.UpdatedAt,
                        &t.UserNama, &t.UserEmail, &t.AssignedNama, &t.UnreadCount, &t.LastMessage,
                ); err != nil {
                        return nil, err
                }
                list = append(list, t)
        }
        return list, nil
}

func (r *ChatRepository) UpdateTiketStatus(ctx context.Context, id uuid.UUID, status models.StatusTiket, assignedTo *uuid.UUID) error {
        _, err := r.db.Exec(ctx,
                "UPDATE chat_tiket SET status=$1, assigned_to=$2, updated_at=NOW() WHERE id=$3",
                status, assignedTo, id)
        return err
}

// AssignTiket — HRD ambil tiket, set assigned_to + status diproses
func (r *ChatRepository) AssignTiket(ctx context.Context, tiketID, hrdID uuid.UUID) error {
        _, err := r.db.Exec(ctx,
                "UPDATE chat_tiket SET assigned_to=$1, status='diproses', updated_at=NOW() WHERE id=$2",
                hrdID, tiketID)
        return err
}

func (r *ChatRepository) AddPesan(ctx context.Context, p *models.ChatPesan) error {
        query := `INSERT INTO chat_pesan (tiket_id, sender_id, pesan) VALUES ($1,$2,$3)
          RETURNING id, is_read, created_at`
        err := r.db.QueryRow(ctx, query, p.TiketID, p.SenderID, p.Pesan).
                Scan(&p.ID, &p.IsRead, &p.CreatedAt)
        if err != nil {
                return err
        }
        r.db.Exec(ctx, "UPDATE chat_tiket SET updated_at=NOW() WHERE id=$1", p.TiketID)
        return nil
}

func (r *ChatRepository) FindPesanByTiketID(ctx context.Context, tiketID uuid.UUID) ([]models.ChatPesanDetail, error) {
        rows, err := r.db.Query(ctx, `
SELECT p.id, p.tiket_id, p.sender_id, p.pesan, p.is_read, p.created_at,
       u.nama_lengkap, u.role
FROM chat_pesan p
JOIN users u ON u.id = p.sender_id
WHERE p.tiket_id = $1 ORDER BY p.created_at ASC`, tiketID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var list []models.ChatPesanDetail
        for rows.Next() {
                var p models.ChatPesanDetail
                rows.Scan(&p.ID, &p.TiketID, &p.SenderID, &p.Pesan, &p.IsRead, &p.CreatedAt, &p.SenderNama, &p.SenderRole)
                list = append(list, p)
        }
        return list, nil
}

func (r *ChatRepository) MarkPesanRead(ctx context.Context, tiketID uuid.UUID, readerID uuid.UUID) error {
        _, err := r.db.Exec(ctx,
                "UPDATE chat_pesan SET is_read=true WHERE tiket_id=$1 AND sender_id != $2 AND is_read=false",
                tiketID, readerID)
        return err
}

func (r *ChatRepository) CountUnreadForUser(ctx context.Context, userID uuid.UUID) int {
        var count int
        r.db.QueryRow(ctx, `
SELECT COUNT(*) FROM chat_pesan p
JOIN chat_tiket t ON t.id = p.tiket_id
WHERE t.user_id=$1 AND p.sender_id!=$1 AND p.is_read=false`, userID).Scan(&count)
        return count
}

func (r *ChatRepository) CountMenunggu(ctx context.Context) int {
        var count int
        r.db.QueryRow(ctx, "SELECT COUNT(*) FROM chat_tiket WHERE status='menunggu'").Scan(&count)
        return count
}

// SetExpiresAt — set waktu hangus saat tiket dibuat
func (r *ChatRepository) SetExpiresAt(ctx context.Context, tiketID uuid.UUID, jamHangus int) error {
        _, err := r.db.Exec(ctx,
                "UPDATE chat_tiket SET expires_at = NOW() + ($1 * INTERVAL '1 hour') WHERE id=$2",
                jamHangus, tiketID)
        return err
}

// FindTiketKedaluwarsa — tiket menunggu yg sudah melewati expires_at
func (r *ChatRepository) FindTiketKedaluwarsa(ctx context.Context) ([]models.ChatTiket, error) {
        rows, err := r.db.Query(ctx,
                `SELECT id, user_id, nomor_tiket, subjek, status, assigned_to, kategori, expires_at, hangus_at, created_at, updated_at
                 FROM chat_tiket
                 WHERE status = 'menunggu' AND expires_at IS NOT NULL AND expires_at < NOW()`)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var list []models.ChatTiket
        for rows.Next() {
                var t models.ChatTiket
                if err := rows.Scan(&t.ID, &t.UserID, &t.NomorTiket, &t.Subjek, &t.Status, &t.AssignedTo, &t.Kategori, &t.ExpiresAt, &t.HangusAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
                        return nil, err
                }
                list = append(list, t)
        }
        return list, nil
}

// FindTiketIdleSelesai — tiket diproses yg sudah idle > N jam (tidak ada pesan baru)
func (r *ChatRepository) FindTiketIdleSelesai(ctx context.Context, idleJam int) ([]models.ChatTiket, error) {
        rows, err := r.db.Query(ctx,
                `SELECT t.id, t.user_id, t.nomor_tiket, t.subjek, t.status, t.assigned_to, t.kategori, t.expires_at, t.hangus_at, t.created_at, t.updated_at
                 FROM chat_tiket t
                 WHERE t.status = 'diproses'
                   AND t.updated_at < NOW() - ($1 * INTERVAL '1 hour')
                   AND NOT EXISTS (
                       SELECT 1 FROM chat_pesan p
                       WHERE p.tiket_id = t.id
                         AND p.created_at > NOW() - ($1 * INTERVAL '1 hour')
                   )`, idleJam)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var list []models.ChatTiket
        for rows.Next() {
                var t models.ChatTiket
                if err := rows.Scan(&t.ID, &t.UserID, &t.NomorTiket, &t.Subjek, &t.Status, &t.AssignedTo, &t.Kategori, &t.ExpiresAt, &t.HangusAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
                        return nil, err
                }
                list = append(list, t)
        }
        return list, nil
}

// MarkHangus — tandai tiket sebagai hangus
func (r *ChatRepository) MarkHangus(ctx context.Context, tiketID uuid.UUID) error {
        _, err := r.db.Exec(ctx,
                "UPDATE chat_tiket SET status='hangus', hangus_at=NOW(), updated_at=NOW() WHERE id=$1",
                tiketID)
        return err
}

// GetChatTimerConfig — ambil konfigurasi timer dari absensi_config
func (r *ChatRepository) GetChatTimerConfig(ctx context.Context) (hangusJam, idleJam int) {
        hangusJam, idleJam = 24, 48
        r.db.QueryRow(ctx,
                `SELECT chat_timer_hangus_jam, chat_timer_idle_jam FROM absensi_config LIMIT 1`).
                Scan(&hangusJam, &idleJam)
        return
}

