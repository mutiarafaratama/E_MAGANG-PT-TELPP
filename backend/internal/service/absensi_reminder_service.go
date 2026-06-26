package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/telpp/emagang/internal/models"
	"github.com/telpp/emagang/internal/repository"
)

// AbsensiReminderService — kirim push notif reminder absen masuk & pulang ke peserta aktif
type AbsensiReminderService struct {
	absensiRepo *repository.AbsensiRepository
	configRepo  *repository.AbsensiConfigRepository
	notifSvc    *NotifikasiService
}

func NewAbsensiReminderService(
	absensiRepo *repository.AbsensiRepository,
	configRepo *repository.AbsensiConfigRepository,
	notifSvc *NotifikasiService,
) *AbsensiReminderService {
	return &AbsensiReminderService{
		absensiRepo: absensiRepo,
		configRepo:  configRepo,
		notifSvc:    notifSvc,
	}
}

// Start — jalankan scheduler di background goroutine, cek setiap menit
func (s *AbsensiReminderService) Start(ctx context.Context) {
	go func() {
		log.Println("[absensi-reminder] layanan reminder absensi aktif (cek tiap 1 menit)")
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				s.cek(ctx, t)
			case <-ctx.Done():
				log.Println("[absensi-reminder] layanan dihentikan")
				return
			}
		}
	}()
}

// kurangiMenit — mengurangi jam string "HH:MM" sejumlah menit, return "HH:MM"
func kurangiMenit(jamStr string, menit int) string {
	t, err := time.Parse("15:04", jamStr)
	if err != nil {
		return jamStr
	}
	return t.Add(-time.Duration(menit) * time.Minute).Format("15:04")
}

// cek — evaluasi apakah jam sekarang cocok dengan salah satu trigger reminder
func (s *AbsensiReminderService) cek(ctx context.Context, now time.Time) {
	wib, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback ke UTC+7 manual jika zona tidak tersedia
		wib = time.FixedZone("WIB", 7*60*60)
	}
	nowWIB := now.In(wib)

	// Skip Sabtu & Minggu — tidak ada magang di hari libur
	hari := nowWIB.Weekday()
	if hari == time.Saturday || hari == time.Sunday {
		return
	}

	jamSekarang := nowWIB.Format("15:04")
	tanggal := nowWIB.Format("2006-01-02")

	cfg, err := s.configRepo.Get(ctx)
	if err != nil {
		log.Printf("[absensi-reminder] gagal ambil config: %v", err)
		return
	}

	// 4 waktu trigger (dari AbsensiConfig):
	// 1. JamMasukBuka       → reminder masuk (window baru buka)
	// 2. JamMasukTutup-10m  → peringatan terakhir masuk
	// 3. JamPulangBuka      → reminder pulang (window baru buka)
	// 4. JamPulangTutup-10m → peringatan terakhir pulang

	peringatanMasuk  := kurangiMenit(cfg.JamMasukTutup, 10)
	peringatanPulang := kurangiMenit(cfg.JamPulangTutup, 10)

	switch jamSekarang {

	case cfg.JamMasukBuka:
		// Window absen masuk baru dibuka → reminder ke semua peserta aktif
		peserta, err := s.absensiRepo.GetPesertaBelumAbsenMasuk(ctx, tanggal)
		if err != nil {
			log.Printf("[absensi-reminder] gagal query belum masuk: %v", err)
			return
		}
		for _, p := range peserta {
			s.notifSvc.KirimKeUser(ctx, p.UserID, models.RolePeserta,
				"⏰ Waktunya Absen Masuk!",
				fmt.Sprintf("Window absen masuk sudah dibuka sampai jam %s. Jangan lupa absen sekarang.", cfg.JamMasukTutup),
				string(models.NotifAbsensi), nil,
			)
		}
		if len(peserta) > 0 {
			log.Printf("[absensi-reminder] reminder masuk → %d peserta", len(peserta))
		}

	case peringatanMasuk:
		// 10 menit sebelum absen masuk tutup → peringatan terakhir
		peserta, err := s.absensiRepo.GetPesertaBelumAbsenMasuk(ctx, tanggal)
		if err != nil {
			log.Printf("[absensi-reminder] gagal query belum masuk: %v", err)
			return
		}
		for _, p := range peserta {
			s.notifSvc.KirimKeUser(ctx, p.UserID, models.RolePeserta,
				"🚨 Segera Absen Masuk!",
				fmt.Sprintf("Absen masuk tutup dalam 10 menit lagi (jam %s). Segera absen sekarang!", cfg.JamMasukTutup),
				string(models.NotifAbsensi), nil,
			)
		}
		if len(peserta) > 0 {
			log.Printf("[absensi-reminder] peringatan terakhir masuk → %d peserta", len(peserta))
		}

	case cfg.JamPulangBuka:
		// Window absen pulang baru dibuka → reminder ke peserta yang sudah masuk
		peserta, err := s.absensiRepo.GetPesertaBelumAbsenKeluar(ctx, tanggal)
		if err != nil {
			log.Printf("[absensi-reminder] gagal query belum keluar: %v", err)
			return
		}
		for _, p := range peserta {
			s.notifSvc.KirimKeUser(ctx, p.UserID, models.RolePeserta,
				"⏰ Waktunya Absen Pulang!",
				fmt.Sprintf("Window absen pulang sudah dibuka sampai jam %s. Jangan lupa absen pulang.", cfg.JamPulangTutup),
				string(models.NotifAbsensi), nil,
			)
		}
		if len(peserta) > 0 {
			log.Printf("[absensi-reminder] reminder pulang → %d peserta", len(peserta))
		}

	case peringatanPulang:
		// 10 menit sebelum absen pulang tutup → peringatan terakhir
		peserta, err := s.absensiRepo.GetPesertaBelumAbsenKeluar(ctx, tanggal)
		if err != nil {
			log.Printf("[absensi-reminder] gagal query belum keluar: %v", err)
			return
		}
		for _, p := range peserta {
			s.notifSvc.KirimKeUser(ctx, p.UserID, models.RolePeserta,
				"🚨 Segera Absen Pulang!",
				fmt.Sprintf("Absen pulang tutup dalam 10 menit lagi (jam %s). Segera absen sekarang!", cfg.JamPulangTutup),
				string(models.NotifAbsensi), nil,
			)
		}
		if len(peserta) > 0 {
			log.Printf("[absensi-reminder] peringatan terakhir pulang → %d peserta", len(peserta))
		}
	}
}
