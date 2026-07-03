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
        NotifLaporan     NotifTipe = "laporan"
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
        Route       string     `json:"route"`       // path Vue Router, contoh: /dashboard/chat
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
// PENTING: routes harus match persis dengan path di artifacts/frontend/src/router/index.js
func RouteForNotif(tipe string, refID *uuid.UUID) string {
        switch NotifTipe(tipe) {
        case NotifPengajuan, NotifDokumen:
                // Peserta: halaman status pengajuan ada di /dashboard
                return "/dashboard"

        case NotifPelaksanaan:
                return "/dashboard"

        case NotifAbsensi:
                return "/dashboard/absensi"

        case NotifNilai:
                return "/dashboard/nilai"

        case NotifSertifikat:
                return "/dashboard/sertifikat"

        case NotifChat:
                // /dashboard/chat — tidak ada /:id di router peserta
                return "/dashboard/chat"

        case NotifLaporan:
                return "/dashboard/laporan"

        case NotifSistem:
                return "/dashboard"

        default:
                return "/dashboard"
        }
}

// RouteForRole — HRD/Admin punya path berbeda dari Peserta
// Semua routes harus match persis dengan Vue Router di frontend
func RouteForRole(role UserRole, tipe string, refID *uuid.UUID) string {
        switch role {
        case RoleHRD:
                switch NotifTipe(tipe) {
                case NotifPengajuan, NotifDokumen:
                        return "/staff/verifikasi"
                case NotifChat:
                        return "/staff/chat"
                case NotifAbsensi:
                        return "/staff/absen"
                case NotifPelaksanaan:
                        return "/staff/berlangsung"
                case NotifNilai:
                        return "/staff/penilaian"
                case NotifSertifikat:
                        return "/staff/sertifikat"
                case NotifLaporan:
                        return "/staff/laporan"
                default:
                        return "/staff"
                }

        case RoleAdmin:
                // Admin punya halaman operasional yang sama dengan HRD (/staff/*),
                // tapi dashboard utama admin ada di /admin bukan /staff
                switch NotifTipe(tipe) {
                case NotifPengajuan, NotifDokumen:
                        return "/staff/verifikasi"
                case NotifChat:
                        return "/staff/chat"
                case NotifAbsensi:
                        return "/staff/absen"
                case NotifPelaksanaan:
                        return "/staff/berlangsung"
                case NotifNilai:
                        return "/staff/penilaian"
                case NotifSertifikat:
                        return "/staff/sertifikat"
                case NotifLaporan:
                        return "/staff/laporan"
                default:
                        return "/admin"
                }

        default:
                return RouteForNotif(tipe, refID)
        }
}
