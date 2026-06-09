package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/telpp/emagang/internal/models"
	"github.com/telpp/emagang/internal/repository"
)

type PengajuanService struct {
	repo      *repository.PengajuanRepository
	notifSvc  *NotifikasiService
	userRepo  *repository.UserRepository
}

func NewPengajuanService(repo *repository.PengajuanRepository, notifSvc *NotifikasiService, userRepo *repository.UserRepository) *PengajuanService {
	return &PengajuanService{repo: repo, notifSvc: notifSvc, userRepo: userRepo}
}

func (s *PengajuanService) Submit(ctx context.Context, userID uuid.UUID, req1 models.PengajuanStep1Request, req2 models.PengajuanStep2Request) (*models.PengajuanMagang, error) {
	// Cek apakah sudah pernah mengajukan dengan status aktif
	existing, err := s.repo.FindByUserID(ctx, userID)
	if err == nil && existing != nil {
		activeStatuses := []models.StatusPengajuan{
			models.StatusDiajukan, models.StatusMenungguVerifikasi, models.StatusDiproses,
		}
		for _, st := range activeStatuses {
			if existing.Status == st {
				return nil, errors.New("anda sudah memiliki pengajuan yang sedang diproses")
			}
		}
	}

	tgl, err := parseDate(req1.TanggalLahir)
	if err != nil {
		return nil, errors.New("format tanggal lahir tidak valid (gunakan YYYY-MM-DD)")
	}

	p := &models.PengajuanMagang{
		UserID:         userID,
		NamaLengkap:    req1.NamaLengkap,
		TempatLahir:    req1.TempatLahir,
		TanggalLahir:   tgl,
		JenisKelamin:   req1.JenisKelamin,
		Alamat:         req1.Alamat,
		NoHP:           req1.NoHP,
		Email:          req1.Email,
		KategoriMagang: req2.KategoriMagang,
		NomorInduk:     req2.NomorInduk,
		AsalInstitusi:  req2.AsalInstitusi,
		Jurusan:        req2.Jurusan,
		KelasSemester:  req2.KelasSemester,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	// Notif ke semua HRD — route langsung ke detail pengajuan
	hrdList, _ := s.userRepo.FindHRDList(ctx)
	for _, h := range hrdList {
		s.notifSvc.KirimKeUser(ctx, h.ID, h.Role,
			"Pengajuan Magang Baru",
			p.NamaLengkap+" ("+string(p.KategoriMagang)+") dari "+p.AsalInstitusi,
			string(models.NotifPengajuan), &p.ID)
	}

	return p, nil
}

func (s *PengajuanService) GetMyPengajuan(ctx context.Context, userID uuid.UUID) (*models.PengajuanMagang, error) {
	p, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (s *PengajuanService) GetAll(ctx context.Context, status, search string, page, limit int) ([]models.PengajuanMagang, int, error) {
	return s.repo.FindAll(ctx, status, search, page, limit)
}

func (s *PengajuanService) GetDetail(ctx context.Context, id uuid.UUID) (*models.PengajuanMagang, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *PengajuanService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.StatusPengajuan, catatan string, changedBy uuid.UUID) error {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("pengajuan tidak ditemukan")
	}

	var catatanPtr *string
	if catatan != "" {
		catatanPtr = &catatan
	}

	if err := s.repo.UpdateStatus(ctx, id, status, catatanPtr, changedBy); err != nil {
		return err
	}

	// Notif ke peserta — route berdasarkan status
	type notifData struct{ judul, pesan, tipe string }
	notifMap := map[models.StatusPengajuan]notifData{
		models.StatusMenungguVerifikasi: {
			"Status Pengajuan Diperbarui",
			"Pengajuan Anda sedang menunggu verifikasi berkas oleh HRD",
			string(models.NotifPengajuan),
		},
		models.StatusDiproses: {
			"Berkas Sedang Diproses",
			"Berkas pengajuan Anda sedang diproses oleh tim HRD",
			string(models.NotifPengajuan),
		},
		models.StatusDiterima: {
			"Pengajuan Magang Diterima! 🎉",
			"Selamat! Permohonan magang Anda di PT TanjungEnim Lestari telah diterima",
			string(models.NotifPelaksanaan),
		},
		models.StatusDitolak: {
			"Pengajuan Tidak Dapat Diterima",
			"Permohonan magang Anda belum dapat kami terima saat ini. Silakan periksa catatan HRD",
			string(models.NotifPengajuan),
		},
	}

	if nd, ok := notifMap[status]; ok {
		s.notifSvc.KirimKeUser(ctx, p.UserID, models.RolePeserta, nd.judul, nd.pesan, nd.tipe, &id)
	}

	return nil
}

func (s *PengajuanService) GetStatusHistory(ctx context.Context, pengajuanID uuid.UUID) ([]models.StatusHistory, error) {
	return s.repo.GetStatusHistory(ctx, pengajuanID)
}
