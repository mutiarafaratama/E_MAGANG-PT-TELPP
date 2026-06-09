package models

import "github.com/google/uuid"

// NotifTipe mendefinisikan semua tipe notifikasi yang dikenal sistem
// Dipakai backend saat Kirim() dan dipakai frontend untuk routing
type NotifTipe string

const (
	NotifPengajuan   NotifTipe = "pengajuan"
	NotifDokumen     NotifTipe = "dokumen"
	NotifPelaksanaan NotifTipe = "pelaksanaan"
	NotifAbsensi     NotifTipe = "absensi"
	NotifNilai       NotifTipe = "nilai"
	NotifSertifikat  NotifTipe = "sertifikat"
	NotifChat        NotifTipe = "chat"
	NotifSistem      NotifTipe = "sistem"
)

// NotifWsPayload adalah payload lengkap yang dikirim via WebSocket
// Frontend Vue membaca ini untuk:
//   - Tampilkan toast (judul + pesan)
//   - Klik toast → navigate ke `route`
//   - Update badge count
type NotifWsPayload struct {
	ID          uuid.UUID  `json:"id"`
	Judul       string     `json:"judul"`
	Pesan       string     `json:"pesan"`
	Tipe        string     `json:"tipe"`
	ReferensiID *uuid.UUID `json:"referensi_id,omitempty"`
	Route       string     `json:"route"`       // path Vue Router, contoh: /dashboard/chat/tiket/uuid
	BadgeCount  int        `json:"badge_count"` // total unread setelah notif ini
	CreatedAt   string     `json:"created_at"`
}

// BadgeWsPayload dikirim setiap kali badge perlu diupdate
type BadgeWsPayload struct {
	TotalUnread  int `json:"total_unread"`
	ChatMenunggu int `json:"chat_menunggu"`
}

// RouteForNotif menghitung Vue Router path berdasarkan tipe dan referensi ID
// Dipanggil saat push notifikasi agar frontend tidak perlu logika routing sendiri
func RouteForNotif(tipe string, refID *uuid.UUID) string {
	switch NotifTipe(tipe) {
	case NotifPengajuan:
		if refID != nil {
			return "/dashboard/pengajuan/" + refID.String()
		}
		return "/dashboard/pengajuan"

	case NotifDokumen:
		if refID != nil {
			return "/dashboard/pengajuan/" + refID.String()
		}
		return "/dashboard/pengajuan"

	case NotifPelaksanaan:
		return "/dashboard/pelaksanaan"

	case NotifAbsensi:
		return "/dashboard/absensi"

	case NotifNilai:
		return "/dashboard/pelaksanaan"

	case NotifSertifikat:
		return "/dashboard/sertifikat"

	case NotifChat:
		if refID != nil {
			return "/dashboard/chat/" + refID.String()
		}
		return "/dashboard/chat"

	case NotifSistem:
		return "/dashboard"

	default:
		return "/dashboard"
	}
}

// RouteForRole — HRD/Admin punya path berbeda dari Peserta
// Gunakan ini jika route tergantung role
func RouteForRole(role UserRole, tipe string, refID *uuid.UUID) string {
	switch role {
	case RoleAdmin, RoleHRD:
		switch NotifTipe(tipe) {
		case NotifPengajuan:
			if refID != nil {
				return "/hrd/pengajuan/" + refID.String()
			}
			return "/hrd/pengajuan"

		case NotifChat:
			if refID != nil {
				return "/hrd/chat/" + refID.String()
			}
			return "/hrd/chat"

		case NotifAbsensi:
			return "/hrd/absensi"

		default:
			return "/hrd/dashboard"
		}

	default:
		return RouteForNotif(tipe, refID)
	}
}
