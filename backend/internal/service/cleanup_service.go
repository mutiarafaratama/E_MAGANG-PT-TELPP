package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/telpp/emagang/internal/config"
	"github.com/telpp/emagang/internal/repository"
	appws "github.com/telpp/emagang/internal/websocket"
)

// CleanupService — hapus otomatis akun peserta yang sudah selesai magang > 30 hari
type CleanupService struct {
	userRepo *repository.UserRepository
	emailSvc *EmailService
}

func NewCleanupService(userRepo *repository.UserRepository, emailSvc *EmailService) *CleanupService {
	return &CleanupService{userRepo: userRepo, emailSvc: emailSvc}
}

// Start — jalankan cleanup di background goroutine.
// Pertama kali berjalan setelah delay 5 menit (beri waktu server fully up),
// kemudian setiap 24 jam.
func (s *CleanupService) Start(ctx context.Context) {
	go func() {
		log.Println("[cleanup] layanan hapus otomatis aktif (jalan pertama dalam 5 menit)")
		select {
		case <-time.After(5 * time.Minute):
		case <-ctx.Done():
			return
		}

		s.JalankanHapusOtomatis(ctx)

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.JalankanHapusOtomatis(ctx)
			case <-ctx.Done():
				log.Println("[cleanup] layanan dihentikan")
				return
			}
		}
	}()
}

// JalankanHapusOtomatis — cari + proses semua kandidat hapus otomatis
func (s *CleanupService) JalankanHapusOtomatis(ctx context.Context) {
	log.Println("[cleanup] menjalankan pengecekan hapus otomatis...")

	kandidat, err := s.userRepo.GetKandidatHapusOtomatis(ctx)
	if err != nil {
		log.Printf("[cleanup] gagal ambil kandidat: %v", err)
		return
	}
	if len(kandidat) == 0 {
		log.Println("[cleanup] tidak ada akun yang perlu dihapus otomatis")
		return
	}

	log.Printf("[cleanup] ditemukan %d akun untuk dihapus otomatis", len(kandidat))
	berhasil := 0
	for _, id := range kandidat {
		if s.hapusSatuPeserta(ctx, id) {
			berhasil++
		}
		// Jeda kecil antar proses agar tidak membanjiri email server
		time.Sleep(3 * time.Second)
	}
	log.Printf("[cleanup] selesai: %d/%d akun berhasil dihapus", berhasil, len(kandidat))
}

// hapusSatuPeserta — kirim backup email → WS force_logout → hapus data
// Return true jika berhasil hapus, false jika ada error kritis
func (s *CleanupService) hapusSatuPeserta(ctx context.Context, id uuid.UUID) bool {
	info, err := s.userRepo.GetPesertaHapusInfo(ctx, id)
	if err != nil {
		log.Printf("[cleanup] gagal ambil info peserta %s: %v", id, err)
		return false
	}

	// Ambil dokumen PDF (laporan/sertifikat upload manual)
	rawDocs, _ := s.userRepo.GetDokumenPeserta(ctx, id)
	var dokumenLampiran []DokumenLampiran
	for _, d := range rawDocs {
		dokumenLampiran = append(dokumenLampiran, DokumenLampiran{
			Jenis:    d.Jenis,
			NamaFile: d.NamaFile,
			PathFile: d.PathFile,
		})
	}

	// Ambil data penilaian untuk generate PDF Lembar Penilaian
	rawPenilaian, _ := s.userRepo.GetPenilaianPeserta(ctx, id)
	var penilaianLampiran *PenilaianLampiranData
	if rawPenilaian != nil {
		kejuruan := make([]KejuruanItem, 0, len(rawPenilaian.Kejuruan))
		for _, k := range rawPenilaian.Kejuruan {
			kejuruan = append(kejuruan, KejuruanItem{
				Komponen: k.Komponen,
				Nilai:    k.Nilai,
			})
		}
		penilaianLampiran = &PenilaianLampiranData{
			NamaLengkap:        rawPenilaian.NamaLengkap,
			NomorInduk:         rawPenilaian.NomorInduk,
			AsalInstitusi:      rawPenilaian.AsalInstitusi,
			Jurusan:            rawPenilaian.Jurusan,
			KelasSemester:      rawPenilaian.KelasSemester,
			Divisi:             rawPenilaian.Divisi,
			Pembimbing:         rawPenilaian.Pembimbing,
			Periode:            rawPenilaian.Periode,
			ManagerNama:        rawPenilaian.ManagerNama,
			Catatan:            rawPenilaian.Catatan,
			DinilaiAt:          rawPenilaian.DinilaiAt,
			NilaiMotivasi:      rawPenilaian.NilaiMotivasi,
			NilaiInisiatif:     rawPenilaian.NilaiInisiatif,
			NilaiDisiplinWaktu: rawPenilaian.NilaiDisiplinWaktu,
			NilaiKerajinan:     rawPenilaian.NilaiKerajinan,
			NilaiKreativitas:   rawPenilaian.NilaiKreativitas,
			NilaiTanggungJawab: rawPenilaian.NilaiTanggungJawab,
			NilaiKerjasama:     rawPenilaian.NilaiKerjasama,
			NilaiAdaptasi:      rawPenilaian.NilaiAdaptasi,
			NilaiKehadiran:     rawPenilaian.NilaiKehadiran,
			NilaiK3Safety:      rawPenilaian.NilaiK3Safety,
			NilaiK3Metode:      rawPenilaian.NilaiK3Metode,
			NilaiK3Manajemen:   rawPenilaian.NilaiK3Manajemen,
			NilaiK3Volume:      rawPenilaian.NilaiK3Volume,
			NilaiPrsProses:     rawPenilaian.NilaiPrsProses,
			NilaiPrsTeori:      rawPenilaian.NilaiPrsTeori,
			NilaiPrsJudul:      rawPenilaian.NilaiPrsJudul,
			NilaiPrsData:       rawPenilaian.NilaiPrsData,
			NilaiAkhir:         rawPenilaian.NilaiAkhir,
			Kejuruan:           kejuruan,
		}
	}

	// Map field nullable
	divisiStr := ""
	if info.Divisi != nil {
		divisiStr = *info.Divisi
	}
	pembimbingStr := ""
	if info.PembimbingNama != nil {
		pembimbingStr = *info.PembimbingNama
	}
	periode := ""
	if info.TanggalMulai != nil && info.TanggalSelesai != nil {
		periode = *info.TanggalMulai + " \u2013 " + *info.TanggalSelesai
	}
	nilaiStr := ""
	if info.Nilai != nil {
		nilaiStr = fmt.Sprintf("%.1f", *info.Nilai)
	}
	sertRelPath := ""
	if info.SertifikatPath != nil {
		sertRelPath = *info.SertifikatPath
	}

	// Kirim email backup (sync — pastikan email terkirim sebelum hapus)
	if s.emailSvc != nil {
		errEmail := s.emailSvc.KirimBackupHapusAkun(
			info.Email, info.NamaLengkap,
			divisiStr, pembimbingStr, periode, nilaiStr,
			config.App.UploadDir, sertRelPath,
			dokumenLampiran, penilaianLampiran,
		)
		if errEmail != nil {
			log.Printf("[cleanup] gagal kirim backup email ke %s: %v — tetap lanjut hapus", info.Email, errEmail)
		} else {
			log.Printf("[cleanup] backup email terkirim ke %s (%s)", info.Email, info.NamaLengkap)
		}
	}

	// Force logout jika peserta sedang online
	appws.GlobalHub.SendToUser(id, "force_logout", nil)

	// Hapus semua data peserta
	if err := s.userRepo.HapusPesertaData(ctx, id); err != nil {
		log.Printf("[cleanup] GAGAL hapus peserta %s (%s): %v", info.NamaLengkap, info.Email, err)
		return false
	}

	log.Printf("[cleanup] ✅ akun %s (%s) berhasil dihapus otomatis setelah >30 hari selesai magang", info.NamaLengkap, info.Email)
	return true
}
