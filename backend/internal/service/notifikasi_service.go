package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/telpp/emagang/internal/models"
	"github.com/telpp/emagang/internal/repository"
	ws "github.com/telpp/emagang/internal/websocket"
)

type NotifikasiService struct {
	repo *repository.NotifikasiRepository
	hub  *ws.Hub
}

func NewNotifikasiService(repo *repository.NotifikasiRepository, hub *ws.Hub) *NotifikasiService {
	return &NotifikasiService{repo: repo, hub: hub}
}

// Kirim simpan notifikasi ke DB dan push ke WebSocket jika user sedang online.
// Jika offline → notifikasi tetap tersimpan di DB, badge muncul saat user buka web.
func (s *NotifikasiService) Kirim(ctx context.Context, userID uuid.UUID, judul, pesan, tipe string, refID *uuid.UUID) {
	n := &models.Notifikasi{
		UserID:      userID,
		Judul:       judul,
		Pesan:       pesan,
		ReferensiID: refID,
	}
	if tipe != "" {
		n.Tipe = &tipe
	}

	// Simpan ke DB — ini yang muncul sebagai badge saat user kembali buka web
	if err := s.repo.Create(ctx, n); err != nil {
		return
	}

	// Hitung badge terbaru setelah notif ini
	totalUnread := s.repo.CountUnread(ctx, userID)

	// Bangun payload lengkap dengan route navigasi
	route := models.RouteForNotif(tipe, refID)
	payload := models.NotifWsPayload{
		ID:          n.ID,
		Judul:       judul,
		Pesan:       pesan,
		Tipe:        tipe,
		ReferensiID: refID,
		Route:       route,
		BadgeCount:  totalUnread,
		CreatedAt:   n.CreatedAt.Format(time.RFC3339),
	}

	// Push via WebSocket — jika user offline ini no-op, badge sudah tersimpan di DB
	s.hub.SendToUser(userID, "notifikasi", payload)
}

// KirimKeUser alias dengan role-aware routing
func (s *NotifikasiService) KirimKeUser(ctx context.Context, userID uuid.UUID, role models.UserRole, judul, pesan, tipe string, refID *uuid.UUID) {
	n := &models.Notifikasi{
		UserID:      userID,
		Judul:       judul,
		Pesan:       pesan,
		ReferensiID: refID,
	}
	if tipe != "" {
		n.Tipe = &tipe
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return
	}

	totalUnread := s.repo.CountUnread(ctx, userID)

	// Gunakan route yang disesuaikan dengan role
	route := models.RouteForRole(role, tipe, refID)
	payload := models.NotifWsPayload{
		ID:          n.ID,
		Judul:       judul,
		Pesan:       pesan,
		Tipe:        tipe,
		ReferensiID: refID,
		Route:       route,
		BadgeCount:  totalUnread,
		CreatedAt:   n.CreatedAt.Format(time.RFC3339),
	}

	s.hub.SendToUser(userID, "notifikasi", payload)
}

// KirimKeRole kirim notif ke beberapa user + broadcast badge update ke role tersebut
func (s *NotifikasiService) KirimKeRole(ctx context.Context, role models.UserRole, judul, pesan, tipe string, refID *uuid.UUID, userIDs []uuid.UUID) {
	for _, uid := range userIDs {
		s.KirimKeUser(ctx, uid, role, judul, pesan, tipe, refID)
	}

	// Broadcast badge_update ke semua client role tersebut yang sedang online
	// Ini memperbarui badge chat_menunggu di sidebar HRD/Admin
	s.hub.SendToRole(role, "badge_update", map[string]interface{}{
		"type": tipe,
	})
}

func (s *NotifikasiService) GetAll(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Notifikasi, int, error) {
	return s.repo.FindByUserID(ctx, userID, page, limit)
}

func (s *NotifikasiService) MarkRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	err := s.repo.MarkRead(ctx, id, userID)
	if err != nil {
		return err
	}
	// Push badge update ke user setelah read
	s.pushBadgeUpdate(ctx, userID, nil)
	return nil
}

func (s *NotifikasiService) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	err := s.repo.MarkAllRead(ctx, userID)
	if err != nil {
		return err
	}
	s.pushBadgeUpdate(ctx, userID, nil)
	return nil
}

// pushBadgeUpdate kirim badge count terbaru via WebSocket ke user
func (s *NotifikasiService) pushBadgeUpdate(ctx context.Context, userID uuid.UUID, chatRepo *repository.ChatRepository) {
	total := s.repo.CountUnread(ctx, userID)
	chatMenunggu := 0
	if chatRepo != nil {
		chatMenunggu = chatRepo.CountMenunggu(ctx)
	}
	s.hub.SendToUser(userID, "badge_update", models.BadgeWsPayload{
		TotalUnread:  total,
		ChatMenunggu: chatMenunggu,
	})
}

// GetBadgeCount dipakai saat user pertama kali buka web (HTTP polling fallback)
func (s *NotifikasiService) GetBadgeCount(ctx context.Context, userID uuid.UUID, chatRepo *repository.ChatRepository) models.NotifBadgeCount {
	return models.NotifBadgeCount{
		TotalUnread:  s.repo.CountUnread(ctx, userID),
		ChatMenunggu: chatRepo.CountMenunggu(ctx),
	}
}
