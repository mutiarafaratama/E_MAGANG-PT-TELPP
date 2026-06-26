package service

import (
        "context"
        "errors"
        "log"
        "time"

        "github.com/google/uuid"
        "github.com/jackc/pgx/v5"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
        ws "github.com/telpp/emagang/internal/websocket"
)

type ChatService struct {
        repo     *repository.ChatRepository
        notifSvc *NotifikasiService
        userRepo *repository.UserRepository
        hub      *ws.Hub
}

func NewChatService(repo *repository.ChatRepository, notifSvc *NotifikasiService, userRepo *repository.UserRepository, hub *ws.Hub) *ChatService {
        return &ChatService{repo: repo, notifSvc: notifSvc, userRepo: userRepo, hub: hub}
}

func (s *ChatService) BuatTiket(ctx context.Context, userID uuid.UUID, req models.ChatBuatTiketRequest) (*models.ChatTiketDetail, error) {
        existing, err := s.repo.FindAktifByUserID(ctx, userID)
        if err != nil && !errors.Is(err, pgx.ErrNoRows) {
                return nil, err
        }
        if existing != nil {
                return nil, errors.New("kamu masih memiliki tiket aktif (" + existing.NomorTiket + "). Selesaikan atau tunggu tiket tersebut ditutup terlebih dahulu")
        }

        kategori := req.Kategori
        if kategori == "" {
                kategori = string(models.KategoriUmum)
        }

        tiket := &models.ChatTiket{
                UserID:     userID,
                NomorTiket: s.repo.GenerateNomorTiket(ctx),
                Subjek:     &req.Subjek,
                Kategori:   models.KategoriTiket(kategori),
        }

        if err := s.repo.CreateTiket(ctx, tiket); err != nil {
                return nil, err
        }

        // Set timer hangus dari konfigurasi
        hangusJam, _ := s.repo.GetChatTimerConfig(ctx)
        s.repo.SetExpiresAt(ctx, tiket.ID, hangusJam)

        pesan := &models.ChatPesan{
                TiketID:  tiket.ID,
                SenderID: userID,
                Pesan:    req.Pesan,
        }
        s.repo.AddPesan(ctx, pesan)

        // Notif ke semua HRD + Admin
        hrdList, _ := s.userRepo.FindHRDList(ctx)
        for _, h := range hrdList {
                s.notifSvc.KirimKeUser(ctx, h.ID, h.Role,
                        "Tiket Chat Baru",
                        "Tiket "+tiket.NomorTiket+": "+req.Subjek,
                        string(models.NotifChat), &tiket.ID)
        }

        // Push badge_update ke HRD/Admin yang online
        chatBadge := models.BadgeWsPayload{ChatMenunggu: s.repo.CountMenunggu(ctx)}
        s.hub.SendToRoles(
                []models.UserRole{models.RoleHRD, models.RoleAdmin},
                "badge_update",
                chatBadge,
        )

        return s.repo.FindTiketByID(ctx, tiket.ID)
}

func (s *ChatService) AmbilTiket(ctx context.Context, tiketID, hrdID uuid.UUID) (*models.ChatTiketDetail, error) {
        tiket, err := s.repo.FindTiketByID(ctx, tiketID)
        if err != nil {
                return nil, errors.New("tiket tidak ditemukan")
        }
        if tiket.Status == models.TiketSelesai || tiket.Status == models.TiketHangus {
                return nil, errors.New("tiket sudah ditutup")
        }
        if err := s.repo.AssignTiket(ctx, tiketID, hrdID); err != nil {
                return nil, err
        }

        hrd, _ := s.userRepo.FindByID(ctx, hrdID)
        nama := "HRD"
        if hrd != nil {
                nama = hrd.NamaLengkap
        }
        s.notifSvc.KirimKeUser(ctx, tiket.UserID, models.RolePeserta,
                "Tiket Diproses",
                "Tiket "+tiket.NomorTiket+" sedang ditangani oleh "+nama,
                string(models.NotifChat), &tiketID)

        chatBadge := models.BadgeWsPayload{ChatMenunggu: s.repo.CountMenunggu(ctx)}
        s.hub.SendToRoles([]models.UserRole{models.RoleHRD, models.RoleAdmin}, "badge_update", chatBadge)

        hrdIDStr := hrdID.String()
        s.hub.SendToUser(tiket.UserID, "tiket_update", models.TiketUpdateWsPayload{
                TiketID:      tiketID.String(),
                Status:       string(models.TiketDiproses),
                AssignedTo:   &hrdIDStr,
                AssignedNama: &nama,
        })
        s.hub.SendToRoles([]models.UserRole{models.RoleHRD, models.RoleAdmin}, "tiket_update", models.TiketUpdateWsPayload{
                TiketID:      tiketID.String(),
                Status:       string(models.TiketDiproses),
                AssignedTo:   &hrdIDStr,
                AssignedNama: &nama,
        })

        return s.repo.FindTiketByID(ctx, tiketID)
}

func (s *ChatService) KirimPesan(ctx context.Context, tiketID uuid.UUID, senderID uuid.UUID, senderRole models.UserRole, req models.ChatKirimPesanRequest) (*models.ChatPesanDetail, error) {
        tiket, err := s.repo.FindTiketByID(ctx, tiketID)
        if err != nil {
                return nil, errors.New("tiket tidak ditemukan")
        }
        if tiket.Status == models.TiketSelesai || tiket.Status == models.TiketHangus {
                return nil, errors.New("tiket sudah ditutup")
        }
        if (senderRole == models.RoleHRD || senderRole == models.RoleAdmin) && tiket.Status == models.TiketMenunggu {
                return nil, errors.New("ambil tiket terlebih dahulu sebelum membalas")
        }
        if senderRole == models.RolePeserta && tiket.UserID != senderID {
                return nil, errors.New("akses ditolak")
        }

        senderUser, _ := s.userRepo.FindByID(ctx, senderID)
        senderNama := "Pengguna"
        if senderUser != nil {
                senderNama = senderUser.NamaLengkap
        }

        pesan := &models.ChatPesan{
                TiketID:  tiketID,
                SenderID: senderID,
                Pesan:    req.Pesan,
        }
        if err := s.repo.AddPesan(ctx, pesan); err != nil {
                return nil, err
        }

        wsPayload := models.ChatMessageWsPayload{
                TiketID:    tiketID.String(),
                SenderID:   senderID.String(),
                SenderNama: senderNama,
                SenderRole: string(senderRole),
                Pesan:      req.Pesan,
                CreatedAt:  pesan.CreatedAt.Format("2006-01-02T15:04:05Z"),
        }

        if senderRole == models.RolePeserta {
                if tiket.AssignedTo != nil {
                        s.hub.SendToUser(*tiket.AssignedTo, "chat_message", wsPayload)
                } else {
                        s.hub.SendToRoles([]models.UserRole{models.RoleHRD, models.RoleAdmin}, "chat_message", wsPayload)
                }
                if tiket.AssignedTo != nil {
                        s.notifSvc.KirimKeUser(ctx, *tiket.AssignedTo, models.RoleHRD,
                                "Pesan Baru dari Peserta",
                                "["+tiket.NomorTiket+"] "+req.Pesan,
                                string(models.NotifChat), &tiketID)
                } else {
                        hrdList, _ := s.userRepo.FindHRDList(ctx)
                        for _, h := range hrdList {
                                s.notifSvc.KirimKeUser(ctx, h.ID, h.Role,
                                        "Pesan Baru dari Peserta",
                                        "["+tiket.NomorTiket+"] "+req.Pesan,
                                        string(models.NotifChat), &tiketID)
                        }
                }
        } else {
                s.hub.SendToUser(tiket.UserID, "chat_message", wsPayload)
                s.notifSvc.KirimKeUser(ctx, tiket.UserID, models.RolePeserta,
                        "Balasan Chat dari HRD",
                        "["+tiket.NomorTiket+"] "+req.Pesan,
                        string(models.NotifChat), &tiketID)
        }

        pesanList, err := s.repo.FindPesanByTiketID(ctx, tiketID)
        if err != nil || len(pesanList) == 0 {
                return nil, err
        }
        last := pesanList[len(pesanList)-1]
        return &last, nil
}

func (s *ChatService) GetTiketSaya(ctx context.Context, userID uuid.UUID) ([]models.ChatTiketDetail, error) {
        return s.repo.FindTiketByUserID(ctx, userID)
}

func (s *ChatService) GetAllTiket(ctx context.Context, status, kategori string, page, limit int) ([]models.ChatTiketDetail, int, error) {
        return s.repo.FindAllTiket(ctx, status, kategori, page, limit)
}

func (s *ChatService) GetPesan(ctx context.Context, tiketID uuid.UUID, readerID uuid.UUID, readerRole models.UserRole) ([]models.ChatPesanDetail, error) {
        tiket, err := s.repo.FindTiketByID(ctx, tiketID)
        if err != nil {
                return nil, errors.New("tiket tidak ditemukan")
        }
        if readerRole == models.RolePeserta && tiket.UserID != readerID {
                return nil, errors.New("akses ditolak")
        }
        s.repo.MarkPesanRead(ctx, tiketID, readerID)
        return s.repo.FindPesanByTiketID(ctx, tiketID)
}

func (s *ChatService) UpdateStatusTiket(ctx context.Context, tiketID uuid.UUID, status models.StatusTiket, assignedTo *uuid.UUID) error {
        tiket, err := s.repo.FindTiketByID(ctx, tiketID)
        if err != nil {
                return s.repo.UpdateTiketStatus(ctx, tiketID, status, assignedTo)
        }
        if err := s.repo.UpdateTiketStatus(ctx, tiketID, status, assignedTo); err != nil {
                return err
        }

        if status == models.TiketSelesai {
                s.notifSvc.KirimKeUser(ctx, tiket.UserID, models.RolePeserta,
                        "Tiket Chat Selesai",
                        "Tiket "+tiket.NomorTiket+" telah diselesaikan oleh HRD",
                        string(models.NotifChat), &tiketID)
        }

        var assignedToStr *string
        if assignedTo != nil {
                str := assignedTo.String()
                assignedToStr = &str
        }
        payload := models.TiketUpdateWsPayload{
                TiketID:    tiketID.String(),
                Status:     string(status),
                AssignedTo: assignedToStr,
        }
        s.hub.SendToUser(tiket.UserID, "tiket_update", payload)
        s.hub.SendToRoles([]models.UserRole{models.RoleHRD, models.RoleAdmin}, "tiket_update", payload)

        return nil
}

// StartTimerGoroutine — jalankan background timer pengecekan tiket kedaluwarsa
// Dipanggil sekali saat server start dari router.go
func (s *ChatService) StartTimerGoroutine(ctx context.Context) {
        go func() {
                ticker := time.NewTicker(15 * time.Minute)
                defer ticker.Stop()
                log.Println("[chat-timer] layanan timer tiket aktif (cek tiap 15 menit)")

                // Jalankan sekali langsung saat start
                s.prosesTimerHangus(ctx)

                for range ticker.C {
                        s.prosesTimerHangus(ctx)
                }
        }()
}

func (s *ChatService) prosesTimerHangus(ctx context.Context) {
        hangusJam, idleJam := s.repo.GetChatTimerConfig(ctx)

        // 1. Tiket menunggu yang melewati batas waktu → hangus
        kedaluwarsa, err := s.repo.FindTiketKedaluwarsa(ctx)
        if err == nil {
                for _, t := range kedaluwarsa {
                        if err := s.repo.MarkHangus(ctx, t.ID); err != nil {
                                continue
                        }
                        log.Printf("[chat-timer] tiket %s hangus (HRD tidak reply dalam %d jam)\n", t.NomorTiket, hangusJam)

                        // Notif ke peserta
                        s.notifSvc.KirimKeUser(ctx, t.UserID, models.RolePeserta,
                                "Tiket Chat Kedaluwarsa",
                                "Tiket "+t.NomorTiket+" hangus karena HRD belum merespon dalam "+
                                        itoa(hangusJam)+" jam. Silakan buat tiket baru.",
                                string(models.NotifChat), &t.ID)

                        // Push WS ke peserta
                        s.hub.SendToUser(t.UserID, "tiket_update", models.TiketUpdateWsPayload{
                                TiketID: t.ID.String(),
                                Status:  string(models.TiketHangus),
                        })
                }
        }

        // 2. Tiket diproses yang idle > idleJam → selesai otomatis
        idle, err := s.repo.FindTiketIdleSelesai(ctx, idleJam)
        if err == nil {
                for _, t := range idle {
                        if err := s.repo.UpdateTiketStatus(ctx, t.ID, models.TiketSelesai, t.AssignedTo); err != nil {
                                continue
                        }
                        log.Printf("[chat-timer] tiket %s ditutup otomatis (idle %d jam)\n", t.NomorTiket, idleJam)

                        s.notifSvc.KirimKeUser(ctx, t.UserID, models.RolePeserta,
                                "Tiket Ditutup Otomatis",
                                "Tiket "+t.NomorTiket+" ditutup otomatis karena tidak ada aktivitas selama "+
                                        itoa(idleJam)+" jam.",
                                string(models.NotifChat), &t.ID)

                        s.hub.SendToUser(t.UserID, "tiket_update", models.TiketUpdateWsPayload{
                                TiketID: t.ID.String(),
                                Status:  string(models.TiketSelesai),
                        })
                }
        }
}

func itoa(n int) string {
        if n == 0 {
                return "0"
        }
        result := ""
        for n > 0 {
                result = string(rune('0'+n%10)) + result
                n /= 10
        }
        return result
}
