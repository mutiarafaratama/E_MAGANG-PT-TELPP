# e-Magang TELPP — Backend Go

Backend API untuk sistem e-Magang PT TanjungEnim Lestari Pulp and Paper.

## Stack
- **Language**: Go 1.22
- **Framework**: Gin
- **Database**: PostgreSQL 15+
- **Auth**: JWT (golang-jwt/jwt/v5)
- **WebSocket**: Gorilla WebSocket (realtime notif & chat)
- **PDF**: gofpdf (sertifikat + rekap absensi)

---

## Prasyarat

Sebelum mulai, pastikan sudah terinstall:

| Tool | Version | Cara Install |
|------|---------|--------------|
| Go | 1.22+ | https://go.dev/dl/ |
| PostgreSQL | 15+ | https://www.postgresql.org/download/ |
| Git | any | https://git-scm.com/ |

---

## LANGKAH 1 — Siapkan Database PostgreSQL

Buka terminal dan jalankan perintah berikut:

```bash
# Masuk ke PostgreSQL
psql -U postgres

# Buat database
CREATE DATABASE emagang_telpp;

# Keluar dari psql
\q
```

Sekarang jalankan migration (buat semua tabel):

```bash
# Dari folder backend/
psql -U postgres -d emagang_telpp -f migrations/001_init.sql
```

Kalau berhasil, akan muncul banyak output `CREATE TYPE`, `CREATE TABLE`, `INSERT`, dll.

**Verifikasi tabel sudah terbuat:**
```bash
psql -U postgres -d emagang_telpp -c "\dt"
```

Harus terlihat 13 tabel:
```
 absensi
 chat_knowledge
 chat_pesan
 chat_tiket
 dokumen
 faq
 landing_content
 notifikasi
 pelaksanaan_magang
 pengajuan_magang
 periode_magang
 refresh_tokens
 status_history
 users
```

---

## LANGKAH 2 — Setup Environment Variables

```bash
# Copy file contoh
cp .env.example .env

# Edit .env dengan text editor favoritmu
nano .env   # atau: code .env / notepad .env
```

Isi file `.env`:
```env
DATABASE_URL=postgres://postgres:PASSWORD_KAMU@localhost:5432/emagang_telpp?sslmode=disable
JWT_SECRET=isi-dengan-string-panjang-acak-minimal-32-karakter-contoh-ini123!
JWT_EXPIRY=24h
REFRESH_EXPIRY=168h
PORT=8080
UPLOAD_DIR=./uploads
MAX_UPLOAD_SIZE=104857600
APP_ENV=development
FRONTEND_URL=http://localhost:5173
```

> ⚠️ **Ganti `PASSWORD_KAMU`** dengan password PostgreSQL kamu.
> ⚠️ **Ganti `JWT_SECRET`** dengan string random yang panjang (minimal 32 karakter).

---

## LANGKAH 3 — Install Dependencies Go

```bash
# Dari folder backend/
go mod tidy
```

Proses ini akan download semua package yang dibutuhkan. Tunggu sampai selesai (~1-2 menit tergantung koneksi internet).

---

## LANGKAH 4 — Jalankan Backend

```bash
# Dari folder backend/
go run cmd/server/main.go
```

Kalau berhasil, output akan seperti ini:
```
✅ Konfigurasi berhasil dimuat (env: development)
✅ Koneksi database PostgreSQL berhasil
🚀 e-Magang TELPP API berjalan di http://localhost:8080
📋 Environment: development
```

---

## LANGKAH 5 — Test API

Buka browser atau gunakan curl untuk test:

```bash
# Health check
curl http://localhost:8080/api/health

# Response yang diharapkan:
# {"app":"e-Magang TELPP","status":"ok"}
```

### Test Login Admin:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@telpp.co.id","password":"Admin@123!"}'
```

### Test Login HRD:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"hrd@telpp.co.id","password":"Admin@123!"}'
```

> ⚠️ **Segera ganti password default** setelah pertama kali login via endpoint `POST /api/auth/change-password`

---

## Struktur Folder

```
backend/
├── cmd/
│   └── server/
│       └── main.go              ← Entry point, jalankan ini
├── internal/
│   ├── config/
│   │   └── config.go            ← Load .env, konfigurasi app
│   ├── database/
│   │   └── postgres.go          ← Koneksi PostgreSQL
│   ├── middleware/
│   │   ├── auth.go              ← JWT generate/validate, role guard
│   │   └── cors.go              ← CORS untuk Vue frontend
│   ├── models/
│   │   └── models.go            ← Semua struct (User, Pengajuan, dll)
│   ├── repository/              ← Query langsung ke database
│   │   ├── user_repo.go
│   │   ├── pengajuan_repo.go
│   │   ├── dokumen_repo.go
│   │   ├── pelaksanaan_repo.go
│   │   ├── absensi_repo.go
│   │   ├── chat_repo.go
│   │   ├── notifikasi_repo.go
│   │   └── landing_repo.go
│   ├── service/                 ← Business logic
│   │   ├── auth_service.go      ← Register, login, validasi password
│   │   ├── pengajuan_service.go ← Submit pengajuan, update status
│   │   ├── chat_service.go      ← Tiket helpdesk, kirim pesan
│   │   ├── notifikasi_service.go← In-app notif + WebSocket push
│   │   └── sertifikat_service.go← Generate PDF sertifikat otomatis
│   ├── handler/                 ← HTTP request handlers (controller)
│   │   ├── auth_handler.go
│   │   ├── pengajuan_handler.go
│   │   ├── dokumen_handler.go
│   │   ├── pelaksanaan_handler.go
│   │   ├── absensi_handler.go
│   │   ├── chat_handler.go
│   │   ├── notifikasi_handler.go
│   │   ├── landing_handler.go
│   │   └── admin_handler.go
│   ├── websocket/
│   │   └── hub.go               ← WebSocket hub (realtime)
│   └── router/
│       └── router.go            ← Semua route API terdaftar di sini
├── migrations/
│   └── 001_init.sql             ← DDL database + seed data
├── uploads/                     ← File upload tersimpan di sini (gitignored)
├── .env.example                 ← Template env, copy ke .env
├── .gitignore
├── go.mod
└── README.md
```

---

## Daftar Endpoint API Lengkap

### AUTH
| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| POST | `/api/auth/register` | Publik | Daftar akun peserta |
| POST | `/api/auth/login` | Publik | Login semua role |
| GET | `/api/auth/me` | Login | Profil user saat ini |
| POST | `/api/auth/change-password` | Login | Ganti password |

### LANDING PAGE
| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| GET | `/api/landing/content` | Publik | Konten landing page |
| GET | `/api/landing/faq` | Publik | Daftar FAQ aktif |
| GET | `/api/landing/periode` | Publik | Periode magang aktif |
| PUT | `/api/landing/content` | HRD/Admin | Edit konten |
| POST | `/api/landing/faq` | HRD/Admin | Tambah/edit FAQ |
| DELETE | `/api/landing/faq/:id` | HRD/Admin | Hapus FAQ |

### PENGAJUAN MAGANG
| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| POST | `/api/pengajuan` | Peserta | Submit pengajuan (body: {step1, step2}) |
| GET | `/api/pengajuan/saya` | Peserta | Lihat pengajuan sendiri |
| GET | `/api/pengajuan` | HRD/Admin | Semua pengajuan (filter: status, search) |
| GET | `/api/pengajuan/:id` | HRD/Admin | Detail pengajuan |
| PATCH | `/api/pengajuan/:id/status` | HRD/Admin | Update status pengajuan |
| GET | `/api/pengajuan/:id/history` | Login | Riwayat perubahan status |

### DOKUMEN
| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| POST | `/api/dokumen/upload` | Login | Upload file (multipart/form-data) |
| GET | `/api/dokumen/pengajuan/:id` | HRD/Admin | Dokumen per pengajuan |
| GET | `/api/dokumen/:id/download` | Login | Download file |

### PELAKSANAAN MAGANG
| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| POST | `/api/pelaksanaan/pengajuan/:id` | HRD/Admin | Set jadwal magang |
| GET | `/api/pelaksanaan/saya` | Peserta | Info pelaksanaan sendiri |
| GET | `/api/pelaksanaan` | HRD/Admin | Semua pelaksanaan |
| PATCH | `/api/pelaksanaan/:id/nilai` | HRD/Admin | Input nilai akhir |
| POST | `/api/pelaksanaan/:id/sertifikat` | HRD/Admin | Generate sertifikat PDF |
| GET | `/api/pelaksanaan/:id/sertifikat/download` | Peserta | Download sertifikat |

### ABSENSI HARIAN
| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| POST | `/api/absensi/checkin` | Peserta | Check-in + isi kegiatan |
| PATCH | `/api/absensi/checkout` | Peserta | Check-out |
| GET | `/api/absensi/saya` | Peserta | Rekap absensi sendiri |
| GET | `/api/absensi/saya/pdf` | Peserta | Download rekap PDF |
| GET | `/api/absensi/pelaksanaan/:id` | HRD/Admin | Monitor absensi peserta |
| PATCH | `/api/absensi/:id/approve` | HRD/Admin | TTD digital pembimbing |

### CHAT / HELPDESK
| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| POST | `/api/chat/tiket` | Peserta | Buat tiket baru |
| GET | `/api/chat/tiket/saya` | Peserta | Tiket milik sendiri |
| GET | `/api/chat/tiket` | HRD/Admin | Semua tiket (filter: status) |
| GET | `/api/chat/tiket/:id/pesan` | Login | Ambil pesan dalam tiket |
| POST | `/api/chat/tiket/:id/pesan` | Login | Kirim pesan |
| PATCH | `/api/chat/tiket/:id/status` | HRD/Admin | Update status tiket |
| GET | `/api/chat/knowledge` | Login | Quick answers |
| POST | `/api/chat/knowledge` | HRD/Admin | Tambah/edit knowledge base |

### NOTIFIKASI
| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| GET | `/api/notifikasi` | Login | Daftar notifikasi |
| PATCH | `/api/notifikasi/:id/read` | Login | Tandai satu sudah dibaca |
| PATCH | `/api/notifikasi/read-all` | Login | Tandai semua sudah dibaca |
| GET | `/api/notifikasi/badge` | Login | Jumlah unread (untuk badge sidebar) |

### WEBSOCKET
| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| GET | `/api/ws` | Login | Koneksi WebSocket realtime |

**Event WebSocket yang dikirim server:**
```
notifikasi     → ada notifikasi baru masuk
chat_badge     → badge jumlah chat update
chat_new_message → ada pesan baru di tiket
badge_update   → update badge notif
```

### ADMIN
| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| GET | `/api/admin/stats` | Admin | Statistik dashboard |
| GET | `/api/admin/users` | Admin | Semua user |
| POST | `/api/admin/users` | Admin | Buat akun HRD/Admin |
| PATCH | `/api/admin/users/:id/toggle` | Admin | Aktif/nonaktif user |
| GET | `/api/admin/hrd-list` | Admin | Daftar HRD untuk assign |
| GET | `/api/admin/periode` | Admin | Semua periode magang |
| POST | `/api/admin/periode` | Admin | Buat periode baru |
| PATCH | `/api/admin/periode/:id/aktif` | Admin | Aktifkan periode |

---

## Cara Pakai di Frontend Vue

Semua request protected wajib pakai header:
```
Authorization: Bearer <access_token>
```

Contoh dengan axios di Vue:
```javascript
// Setelah login, simpan token:
localStorage.setItem('token', response.data.data.access_token)

// Axios interceptor:
axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
```

---

## WebSocket dari Vue

```javascript
// Koneksi WebSocket
const token = localStorage.getItem('token')
const ws = new WebSocket(`ws://localhost:8080/api/ws`)

// Kirim token setelah connect (lewat URL query atau custom header)
// Catatan: gunakan ws://localhost:8080/api/ws?token=${token}

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data)
  
  if (msg.type === 'notifikasi') {
    // Tampilkan toast notifikasi
    showToast(msg.payload.judul)
    // Update badge count
    badgeCount.value++
  }
  
  if (msg.type === 'chat_badge') {
    // Update badge jumlah chat menunggu di sidebar
    chatBadge.value = msg.payload.count
  }
  
  if (msg.type === 'chat_new_message') {
    // Refresh daftar pesan jika tiket sedang dibuka
    if (currentTiketId === msg.payload.tiket_id) {
      loadPesan()
    }
  }
}
```

---

## Troubleshooting

### ❌ `dial tcp: connect: connection refused`
Database PostgreSQL tidak berjalan. Jalankan: `sudo service postgresql start` (Linux) atau buka pgAdmin (Windows).

### ❌ `pq: role "postgres" does not exist`
Ganti `postgres` di DATABASE_URL dengan username PostgreSQL kamu.

### ❌ `JWT_SECRET wajib diisi`
Pastikan file `.env` sudah dibuat dari `.env.example` dan JWT_SECRET sudah diisi.

### ❌ `go: module not found`
Pastikan kamu menjalankan `go mod tidy` sebelum `go run`.

### ❌ Port 8080 sudah dipakai
Ganti `PORT=8080` di `.env` ke port lain, misalnya `PORT=9090`.

---

## Default Accounts (Seed)

| Role | Email | Password |
|------|-------|----------|
| Admin | admin@telpp.co.id | Admin@123! |
| HRD | hrd@telpp.co.id | Admin@123! |

> ⚠️ **WAJIB ganti password ini** setelah pertama kali login!

---

## Build untuk Production

```bash
# Build binary
go build -o emagang-api cmd/server/main.go

# Jalankan binary
APP_ENV=production ./emagang-api
```
