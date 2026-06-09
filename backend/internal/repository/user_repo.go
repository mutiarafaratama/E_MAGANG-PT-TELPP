package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/telpp/emagang/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (id, nama_lengkap, email, password_hash, role)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, u.NamaLengkap, u.Email, u.PasswordHash, u.Role).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, nama_lengkap, email, password_hash, role, is_active, created_at, updated_at
			  FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.NamaLengkap, &u.Email, &u.PasswordHash,
		&u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, nama_lengkap, email, password_hash, role, is_active, created_at, updated_at
			  FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.NamaLengkap, &u.Email, &u.PasswordHash,
		&u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindAll(ctx context.Context, role string, page, limit int) ([]models.UserPublic, int, error) {
	offset := (page - 1) * limit
	var args []interface{}
	where := "WHERE 1=1"
	argIdx := 1

	if role != "" {
		where += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, role)
		argIdx++
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM users " + where
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	query := fmt.Sprintf(`SELECT id, nama_lengkap, email, role, is_active, created_at
		FROM users %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.UserPublic
	for rows.Next() {
		var u models.UserPublic
		if err := rows.Scan(&u.ID, &u.NamaLengkap, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *UserRepository) UpdateActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	_, err := r.db.Exec(ctx,
		"UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2",
		isActive, id)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error {
	_, err := r.db.Exec(ctx,
		"UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2",
		hash, id)
	return err
}

func (r *UserRepository) FindHRDList(ctx context.Context) ([]models.UserPublic, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, nama_lengkap, email, role, is_active, created_at
		 FROM users WHERE role IN ('hrd', 'admin') AND is_active = true ORDER BY nama_lengkap`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserPublic
	for rows.Next() {
		var u models.UserPublic
		if err := rows.Scan(&u.ID, &u.NamaLengkap, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		 VALUES ($1, $2, NOW() + INTERVAL '7 days')`,
		userID, tokenHash)
	return err
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx,
		"DELETE FROM refresh_tokens WHERE token_hash = $1", tokenHash)
	return err
}

func (r *UserRepository) GetDashboardStats(ctx context.Context) (*models.DashboardAdminStats, error) {
	stats := &models.DashboardAdminStats{}
	err := r.db.QueryRow(ctx, `
		SELECT
			(SELECT COUNT(*) FROM users WHERE role = 'peserta' AND is_active = true),
			(SELECT COUNT(*) FROM pengajuan_magang),
			(SELECT COUNT(*) FROM pengajuan_magang WHERE status IN ('diajukan','menunggu_verifikasi','diproses')),
			(SELECT COUNT(*) FROM pengajuan_magang WHERE status = 'diterima'),
			(SELECT COUNT(*) FROM pengajuan_magang WHERE status = 'ditolak'),
			(SELECT COUNT(*) FROM pelaksanaan_magang WHERE status = 'aktif'),
			(SELECT COUNT(*) FROM pelaksanaan_magang WHERE status = 'selesai'),
			(SELECT COUNT(*) FROM chat_tiket WHERE status = 'menunggu')
	`).Scan(
		&stats.TotalPeserta,
		&stats.TotalPengajuan,
		&stats.PengajuanMenunggu,
		&stats.PengajuanDiterima,
		&stats.PengajuanDitolak,
		&stats.MagangAktif,
		&stats.MagangSelesai,
		&stats.TiketChatMenunggu,
	)
	return stats, err
}
