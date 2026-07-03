package service

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/telpp/emagang/internal/models"
	"github.com/telpp/emagang/internal/repository"
)

type LaporanService struct {
	repo     *repository.LaporanRepository
	pelRepo  *repository.PelaksanaanRepository
	userRepo *repository.UserRepository
	notifSvc *NotifikasiService
}

func NewLaporanService(
	repo *repository.LaporanRepository,
	pelRepo *repository.PelaksanaanRepository,
	userRepo *repository.UserRepository,
	notifSvc *NotifikasiService,
) *LaporanService {
	return &LaporanService{repo: repo, pelRepo: pelRepo, userRepo: userRepo, notifSvc: notifSvc}
}

// Upload — peserta upload laporan, lalu kirim notif push ke semua HRD
func (s *LaporanService) Upload(ctx context.Context, pelaksanaanID uuid.UUID, fh *multipart.FileHeader) (*models.LaporanMagang, error) {
	laporan, err := s.repo.Upload(ctx, pelaksanaanID, fh)
	if err != nil {
		return nil, err
	}

	// Kirim push notif ke semua HRD
	laporanID := laporan.ID
	hrdList, _ := s.userRepo.FindHRDList(ctx)
	for _, h := range hrdList {
		s.notifSvc.KirimKeUser(ctx, h.ID, h.Role,
			"📄 Laporan Magang Baru",
			"Peserta telah mengunggah laporan magang. Silakan periksa dan berikan review.",
			string(models.NotifLaporan),
			&laporanID,
		)
	}

	return laporan, nil
}

// Review — HRD acc atau minta revisi, kirim notif push ke peserta
func (s *LaporanService) Review(ctx context.Context, id uuid.UUID, status models.StatusLaporan, catatan string, reviewerID uuid.UUID) error {
	// Ambil laporan untuk tahu pelaksanaan_id
	laporan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Review(ctx, id, status, catatan, reviewerID); err != nil {
		return err
	}

	// Ambil pelaksanaan untuk tahu user_id peserta
	pel, err := s.pelRepo.FindByID(ctx, laporan.PelaksanaanID)
	if err != nil {
		return nil
	}

	laporanID := laporan.ID

	switch status {
	case models.StatusLaporanRevisi:
		s.notifSvc.KirimKeUser(ctx, pel.UserID, models.RolePeserta,
			"⚠️ Laporan Perlu Direvisi",
			"HRD meminta revisi pada laporan magang Anda. Segera periksa catatan dan upload ulang laporan.",
			string(models.NotifLaporan),
			&laporanID,
		)
	case models.StatusLaporanDisetujui:
		s.notifSvc.KirimKeUser(ctx, pel.UserID, models.RolePeserta,
			"✅ Laporan Magang Disetujui",
			"Selamat! Laporan magang Anda telah disetujui oleh HRD.",
			string(models.NotifLaporan),
			&laporanID,
		)
	}

	return nil
}

// FindByID — proxy ke repo
func (s *LaporanService) FindByID(ctx context.Context, id uuid.UUID) (*models.LaporanMagang, error) {
	return s.repo.FindByID(ctx, id)
}

// UpdateStatusPelaksanaan — set status pelaksanaan setelah laporan disetujui
func (s *LaporanService) UpdateStatusPelaksanaan(ctx context.Context, pelaksanaanID uuid.UUID, status models.StatusPelaksanaan) error {
	return s.pelRepo.UpdateStatus(ctx, pelaksanaanID, status)
}
