# Alur Kerja e-Magang TELPP — Semua Role

Sistem magang PT TanjungEnim Lestari Pulp and Paper.  
3 role: **Peserta** · **HRD** · **Admin**

---

## Gambaran Besar Alur

```
LANDING PAGE (publik)
      │
      ▼
[Peserta Daftar Akun]
      │
      ▼
[Isi Form Pengajuan] ──► NOTIF ke HRD ──► [HRD Review]
                                                │
                              ┌─────────────────┼─────────────────┐
                              ▼                 ▼                 ▼
                          DITOLAK          DITERIMA         BUTUH REVISI
                              │                 │
                              ▼                 ▼
                       [Peserta Revisi]   [Pelaksanaan Magang]
                                                │
                                         [Absensi Harian]
                                                │
                                         [Upload Laporan]
                                                │
                                         [HRD Beri Nilai]
                                                │
                                         [Generate Sertifikat PDF]
                                                │
                                         [Peserta Download]
```

---

## 1. PESERTA MAGANG

### 1.1 Pendaftaran & Pengajuan

| Langkah | Aksi | Endpoint | Status DB |
|---------|------|----------|-----------|
| 1 | Buka landing page, cek syarat & FAQ | `GET /api/landing/content` | — |
| 2 | Daftar akun (email + password) | `POST /api/auth/register` | `users` role=peserta |
| 3 | Login | `POST /api/auth/login` | `refresh_tokens` dibuat |
| 4 | Isi form Step 1 (data diri) | `POST /api/pengajuan` (body.step1) | — |
| 5 | Isi form Step 2 (data akademik) | `POST /api/pengajuan` (body.step2) | `pengajuan_magang` status=diajukan |
| 6 | Upload dokumen wajib | `POST /api/dokumen/upload` | `dokumen` |
| 7 | Pantau status | `GET /api/pengajuan/saya` | — |

**Dokumen wajib yang harus diupload:**
- Proposal magang (PDF)
- KTP
- KTM (khusus mahasiswa)
- Pasfoto 3×4
- Kartu BPJS/KIS

**Notifikasi yang diterima peserta:**
- ✅ Pengajuan berhasil dikirim
- 📋 Status diperbarui oleh HRD
- 🎉 Pengajuan diterima → otomatis masuk ke pelaksanaan magang

---

### 1.2 Pelaksanaan Magang

| Langkah | Aksi | Endpoint | Keterangan |
|---------|------|----------|------------|
| 1 | Cek jadwal magang | `GET /api/pelaksanaan/saya` | Divisi, pembimbing, tanggal |
| 2 | Isi absensi harian | `POST /api/absensi` | Setiap hari kerja, 1x/hari |
| 3 | Lihat rekap absensi | `GET /api/absensi` | Bisa filter per bulan |
| 4 | Upload laporan akhir | `POST /api/dokumen/upload` (jenis=laporan_magang) | PDF max 100MB |
| 5 | Buka chat helpdesk | `POST /api/chat/tiket` | Jika ada pertanyaan ke HRD |
| 6 | Cek nilai final | `GET /api/pelaksanaan/saya` | Setelah HRD input nilai |
| 7 | Download sertifikat | `GET /api/sertifikat/download/:id` | PDF auto-generate |

**Format Absensi Harian:**
```
Jam masuk  : 08:00
Jam keluar : 17:00
Keterangan : hadir / izin / sakit
Kegiatan   : [deskripsi kegiatan hari ini]
```

**Alur Status Pelaksanaan:**
```
menunggu_mulai → aktif → upload_laporan → penilaian → selesai
```

---

### 1.3 Chat Helpdesk

| Langkah | Aksi | Endpoint |
|---------|------|----------|
| 1 | Cek knowledge base dulu | `GET /api/chat/knowledge` |
| 2 | Buat tiket jika tidak ketemu jawaban | `POST /api/chat/tiket` |
| 3 | Kirim pesan di tiket | `POST /api/chat/tiket/:id/pesan` |
| 4 | Terima balasan via notif real-time | WebSocket `type: "notifikasi"` |
| 5 | Lihat semua tiket saya | `GET /api/chat/tiket/saya` |

---

### 1.4 Notifikasi & Badge

- Badge merah di lonceng = `GET /api/notifikasi/badge` → `total_unread`
- Klik lonceng = `GET /api/notifikasi` (list)
- Klik notif → navigate ke `payload.route` (sudah ada di payload WS)
- Tandai baca = `PATCH /api/notifikasi/:id/read`
- Tandai semua = `PATCH /api/notifikasi/read-all`

**Event WebSocket yang diterima peserta:**
| Event | Kapan |
|-------|-------|
| `connected` | Saat pertama terhubung |
| `notifikasi` | Ada notif baru (apapun tipenya) |
| `badge_update` | Setelah mark-read |

---

## 2. HRD (Staff)

### 2.1 Kelola Pengajuan Masuk

| Langkah | Aksi | Endpoint | Keterangan |
|---------|------|----------|------------|
| 1 | Lihat daftar pengajuan | `GET /api/hrd/pengajuan` | Filter: status, kategori, search |
| 2 | Buka detail + dokumen | `GET /api/hrd/pengajuan/:id` | Lihat semua file upload |
| 3 | Download dokumen | `GET /api/dokumen/:id/download` | Verifikasi keaslian |
| 4 | Update status | `PATCH /api/hrd/pengajuan/:id/status` | Dengan catatan |
| 5 | Peserta diterima → buat record pelaksanaan | `POST /api/hrd/pelaksanaan` | Set divisi, tanggal, pembimbing |

**Status yang bisa di-set HRD:**
```
diajukan
  → menunggu_verifikasi  (HRD mulai review berkas)
    → diproses           (berkas lengkap, lanjut proses)
      → diterima         (OTOMATIS buat record pelaksanaan_magang)
      → ditolak          (sertakan catatan alasan)
```

Setiap perubahan status → tersimpan di `status_history` → notif otomatis ke peserta.

---

### 2.2 Kelola Pelaksanaan Magang

| Langkah | Aksi | Endpoint | Keterangan |
|---------|------|----------|------------|
| 1 | Lihat semua peserta aktif | `GET /api/hrd/pelaksanaan` | Filter status |
| 2 | TTD absensi peserta | `PATCH /api/hrd/absensi/:id/ttd` | Tandatangan digital |
| 3 | Download rekap absensi | `GET /api/hrd/absensi/rekap/:pelaksanaan_id` | PDF otomatis |
| 4 | Update status pelaksanaan | `PATCH /api/hrd/pelaksanaan/:id/status` | |
| 5 | Input nilai final | `PATCH /api/hrd/pelaksanaan/:id/nilai` | Nilai + catatan |
| 6 | Generate sertifikat | `POST /api/hrd/sertifikat/generate/:pelaksanaan_id` | PDF via gofpdf |

**Alur Penilaian:**
```
1. Status pelaksanaan → "penilaian"
2. HRD buka halaman penilaian
3. Input nilai (0–100) + catatan nilai
4. Generate sertifikat → PDF tersimpan di ./uploads/sertifikat/
5. Status → "selesai"
6. Peserta dapat notif + bisa download sertifikat
```

---

### 2.3 Chat Helpdesk (HRD)

| Langkah | Aksi | Endpoint | Keterangan |
|---------|------|----------|------------|
| 1 | Lihat semua tiket masuk | `GET /api/hrd/chat/tiket` | Filter: menunggu/diproses/selesai |
| 2 | Ambil tiket (assign ke diri sendiri) | `PATCH /api/hrd/chat/tiket/:id/status` | status=diproses, assigned_to=self |
| 3 | Buka & balas tiket | `GET/POST /api/hrd/chat/tiket/:id/pesan` | Balas real-time |
| 4 | Tutup tiket | `PATCH /api/hrd/chat/tiket/:id/status` | status=selesai |
| 5 | Kelola knowledge base | `POST/PUT /api/hrd/chat/knowledge` | FAQ self-service untuk peserta |

**Badge chat di sidebar HRD:**
- `chat_menunggu` dari `GET /api/notifikasi/badge`
- Update real-time via WebSocket `type: "badge_update"`

---

### 2.4 Event WebSocket yang diterima HRD

| Event | Kapan | Payload `route` |
|-------|-------|-----------------|
| `notifikasi` | Pengajuan baru masuk | `/hrd/pengajuan/:id` |
| `notifikasi` | Pesan chat baru dari peserta | `/hrd/chat/:tiket_id` |
| `badge_update` | Badge perlu diperbarui | `.total_unread`, `.chat_menunggu` |

---

## 3. ADMIN

### 3.1 Kelola Pengguna

| Aksi | Endpoint | Keterangan |
|------|----------|------------|
| Lihat semua user | `GET /api/admin/users` | Filter role, search nama/email |
| Buat akun HRD baru | `POST /api/admin/users` | Role default: hrd |
| Aktifkan/nonaktifkan user | `PATCH /api/admin/users/:id/toggle-active` | Blokir akses tanpa hapus data |

---

### 3.2 Kelola Periode Magang

| Aksi | Endpoint | Keterangan |
|------|----------|------------|
| Lihat semua periode | `GET /api/admin/periode` | |
| Buat periode baru | `POST /api/admin/periode` | Nama, tanggal buka-tutup, kuota |
| Aktifkan periode | `PATCH /api/admin/periode/:id/aktifkan` | Nonaktifkan periode sebelumnya otomatis |

**Aturan:**
- Hanya 1 periode yang bisa `is_active = true` dalam satu waktu
- Landing page otomatis tampilkan periode aktif ke calon peserta

---

### 3.3 Kelola Landing Page (CMS)

| Aksi | Endpoint | Keterangan |
|------|----------|------------|
| Lihat semua konten | `GET /api/landing/content` | Publik, tanpa login |
| Edit konten | `PUT /api/admin/landing/:kunci` | Key-value: hero_title, syarat_umum, dll |
| Kelola FAQ | `POST/PUT/DELETE /api/admin/faq` | FAQ yang tampil di halaman publik |

**Key CMS yang tersedia:**
| Kunci | Deskripsi |
|-------|-----------|
| `hero_title` | Judul utama halaman publik |
| `hero_subtitle` | Subtitle/tagline |
| `syarat_umum` | Syarat-syarat umum pendaftaran |
| `ketentuan_magang` | Aturan selama magang |
| `tempat_kumpul` | Lokasi/alamat perusahaan |
| `kontak_hrd` | Email & nomor telepon HRD |

---

### 3.4 Dashboard & Statistik

| Aksi | Endpoint | Data |
|------|----------|------|
| Statistik pengajuan | `GET /api/admin/statistik` | Total per status, per kategori |
| Peserta aktif | — | Via pelaksanaan_magang |
| Aktivitas sistem | — | Via status_history |

---

## 4. Alur Lengkap End-to-End

```
┌─────────────────────────────────────────────────────────┐
│                    HALAMAN PUBLIK                        │
│  Landing Page → Syarat → FAQ → Periode Aktif            │
└──────────────────────────┬──────────────────────────────┘
                           │
                    [Daftar / Login]
                           │
           ┌───────────────┼───────────────┐
           ▼               ▼               ▼
        PESERTA           HRD           ADMIN
           │               │               │
     ┌─────┴──┐      ┌─────┴──┐      ┌─────┴──┐
     │        │      │        │      │        │
  Ajukan   Chat    Review   Chat  Periode  Landing
  Magang   Help   Berkas   Help   Magang    CMS
     │        │      │        │      │
  Upload   Tiket  Setuju/  Balas  Kelola
  Dokumen  Masuk  Tolak   Tiket   User
     │               │
  Magang           Buat
  Diterima       Pelaksanaan
     │               │
  Mulai            TTD
  Absensi        Absensi
  Harian            │
     │            Input
  Upload          Nilai
  Laporan           │
     │            Generate
  Terima          Sertifikat
  Notif             │
     │            Push Notif
  Download    ◄─── ke Peserta
  Sertifikat
```

---

## 5. Peta Notifikasi Otomatis

| Event | Pengirim | Penerima | Tipe | Route |
|-------|----------|----------|------|-------|
| Pengajuan baru submit | Sistem | Semua HRD | `pengajuan` | `/hrd/pengajuan/:id` |
| Status pengajuan berubah | HRD | Peserta | `pengajuan` | `/dashboard/pengajuan/:id` |
| Pengajuan diterima | HRD | Peserta | `pelaksanaan` | `/dashboard/pelaksanaan` |
| Tiket chat baru | Sistem | Semua HRD | `chat` | `/hrd/chat/:tiket_id` |
| Pesan baru dari peserta | Sistem | HRD assigned | `chat` | `/hrd/chat/:tiket_id` |
| Balasan dari HRD | Sistem | Peserta | `chat` | `/dashboard/chat/:tiket_id` |
| Sertifikat siap | Sistem | Peserta | `sertifikat` | `/dashboard/sertifikat` |

---

## 6. Akun Default (Seed)

| Role | Email | Password | Catatan |
|------|-------|----------|---------|
| Admin | admin@telpp.co.id | `Admin@123!` | **Ganti setelah login pertama** |
| HRD | hrd@telpp.co.id | `Admin@123!` | **Ganti setelah login pertama** |

---

## 7. Endpoint API Ringkasan

### Auth (Publik)
```
POST   /api/auth/register
POST   /api/auth/login
GET    /api/auth/me              [🔒 login]
POST   /api/auth/change-password [🔒 login]
```

### Landing (Publik)
```
GET    /api/landing/content
GET    /api/landing/faq
GET    /api/landing/periode
```

### Peserta [🔒 role=peserta]
```
POST   /api/pengajuan
GET    /api/pengajuan/saya
POST   /api/dokumen/upload
GET    /api/dokumen/:id/download
GET    /api/pelaksanaan/saya
POST   /api/absensi
GET    /api/absensi
GET    /api/sertifikat/download/:id
GET    /api/chat/tiket/saya
POST   /api/chat/tiket
POST   /api/chat/tiket/:id/pesan
GET    /api/chat/tiket/:id/pesan
GET    /api/chat/knowledge
```

### HRD [🔒 role=hrd]
```
GET    /api/hrd/pengajuan
GET    /api/hrd/pengajuan/:id
PATCH  /api/hrd/pengajuan/:id/status
GET    /api/hrd/pelaksanaan
POST   /api/hrd/pelaksanaan
PATCH  /api/hrd/pelaksanaan/:id/status
PATCH  /api/hrd/pelaksanaan/:id/nilai
GET    /api/hrd/absensi/:pelaksanaan_id
PATCH  /api/hrd/absensi/:id/ttd
GET    /api/hrd/absensi/rekap/:pelaksanaan_id
POST   /api/hrd/sertifikat/generate/:id
GET    /api/hrd/chat/tiket
GET    /api/hrd/chat/tiket/:id/pesan
POST   /api/hrd/chat/tiket/:id/pesan
PATCH  /api/hrd/chat/tiket/:id/status
POST   /api/hrd/chat/knowledge
PUT    /api/hrd/chat/knowledge/:id
```

### Admin [🔒 role=admin]
```
GET    /api/admin/users
POST   /api/admin/users
PATCH  /api/admin/users/:id/toggle-active
GET    /api/admin/periode
POST   /api/admin/periode
PATCH  /api/admin/periode/:id/aktifkan
GET    /api/admin/statistik
PUT    /api/admin/landing/:kunci
POST   /api/admin/faq
PUT    /api/admin/faq/:id
DELETE /api/admin/faq/:id
```

### Notifikasi [🔒 semua role]
```
GET    /api/notifikasi
GET    /api/notifikasi/badge
PATCH  /api/notifikasi/:id/read
PATCH  /api/notifikasi/read-all
GET    /api/ws?token=<jwt>       ← WebSocket realtime
```
