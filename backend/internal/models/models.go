package models

import (
        "time"

        "github.com/google/uuid"
)

// ============================================================
// ENUMS (sebagai konstanta string)
// ============================================================

type UserRole string

const (
        RoleAdmin   UserRole = "admin"
        RoleHRD     UserRole = "hrd"
        RolePeserta UserRole = "peserta"
)

type KategoriMagang string

const (
        KategoriSMK        KategoriMagang = "smk"
        KategoriD3S1S2     KategoriMagang = "d3_s1_s2"
        KategoriPenelitian KategoriMagang = "penelitian"
        KategoriD3         KategoriMagang = "d3"
        KategoriS1         KategoriMagang = "s1"
        KategoriS2         KategoriMagang = "s2"
        KategoriLainnya    KategoriMagang = "lainnya"
)

type JenisKelamin string

const (
        JenisLakiLaki  JenisKelamin = "laki_laki"
        JenisPerempuan JenisKelamin = "perempuan"
)

type StatusPengajuan string

const (
        StatusDiajukan           StatusPengajuan = "diajukan"
        StatusMenungguVerifikasi StatusPengajuan = "menunggu_verifikasi"
        StatusDiproses           StatusPengajuan = "diproses"
        StatusDiterima           StatusPengajuan = "diterima"
        StatusDitolak            StatusPengajuan = "ditolak"
)

type JenisDokumen string

const (
        DokumenProposal       JenisDokumen = "proposal_magang"
        DokumenSuratPengantar JenisDokumen = "surat_pengantar"
        DokumenKTP            JenisDokumen = "ktp"
        DokumenKTM            JenisDokumen = "ktm"
        DokumenPasfoto        JenisDokumen = "pasfoto"
        DokumenBPJS           JenisDokumen = "bpjs_kis"
        DokumenSuratBalas     JenisDokumen = "surat_balasan"
        DokumenLaporan        JenisDokumen = "laporan_magang"
        DokumenSertifikat     JenisDokumen = "sertifikat"
)

type StatusPelaksanaan string

const (
        StatusMenungguMulai StatusPelaksanaan = "menunggu_mulai"
        StatusAktif         StatusPelaksanaan = "aktif"
        StatusUploadLaporan StatusPelaksanaan = "upload_laporan"
        StatusPenilaian     StatusPelaksanaan = "penilaian"
        StatusSelesai       StatusPelaksanaan = "selesai"
)

type StatusTiket string

const (
        TiketMenunggu StatusTiket = "menunggu"
        TiketDiproses StatusTiket = "diproses"
        TiketSelesai  StatusTiket = "selesai"
        TiketHangus   StatusTiket = "hangus"
)

type KategoriTiket string

const (
        KategoriUmum       KategoriTiket = "umum"
        KategoriAbsensi    KategoriTiket = "absensi"
        KategoriDokumen    KategoriTiket = "dokumen"
        KategoriSertifikat KategoriTiket = "sertifikat"
)

// ============================================================
// USER
// ============================================================

type User struct {
        ID              uuid.UUID `json:"id" db:"id"`
        NamaLengkap     string    `json:"nama_lengkap" db:"nama_lengkap"`
        Email           string    `json:"email" db:"email"`
        PasswordHash    string    `json:"-" db:"password_hash"`
        Role            UserRole  `json:"role" db:"role"`
        IsActive        bool      `json:"is_active" db:"is_active"`
        PasswordChanged bool      `json:"password_changed" db:"password_changed"`
        CreatedAt       time.Time `json:"created_at" db:"created_at"`
        UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type UserPublic struct {
        ID              uuid.UUID `json:"id"`
        NamaLengkap     string    `json:"nama_lengkap"`
        Email           string    `json:"email"`
        Role            UserRole  `json:"role"`
        IsActive        bool      `json:"is_active"`
        PasswordChanged bool      `json:"password_changed"`
        CreatedAt       time.Time `json:"created_at"`
        StatusMagang    *string   `json:"status_magang,omitempty"`
}

type RefreshToken struct {
        ID        uuid.UUID `json:"id" db:"id"`
        UserID    uuid.UUID `json:"user_id" db:"user_id"`
        TokenHash string    `json:"-" db:"token_hash"`
        ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

type RegisterRequest struct {
        NamaLengkap string `json:"nama_lengkap" binding:"required"`
        Email       string `json:"email" binding:"required,email"`
        Password    string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
        AccessToken  string     `json:"access_token"`
        RefreshToken string     `json:"refresh_token"`
        User         UserPublic `json:"user"`
}

type NotifBadgeCount struct {
        TotalUnread  int `json:"total_unread"`
        ChatMenunggu int `json:"chat_menunggu"`
}

// ============================================================
// PENGAJUAN MAGANG
// ============================================================

type PengajuanStep1Request struct {
        NamaLengkap  string `json:"nama_lengkap" binding:"required"`
        TempatLahir  string `json:"tempat_lahir" binding:"required"`
        TanggalLahir string `json:"tanggal_lahir" binding:"required"`
        JenisKelamin string `json:"jenis_kelamin" binding:"required"`
        Alamat       string `json:"alamat" binding:"required"`
        NoHP         string `json:"no_hp" binding:"required"`
        Email        string `json:"email" binding:"required,email"`
}

type PengajuanStep2Request struct {
        KategoriMagang string `json:"kategori_magang" binding:"required"`
        NomorInduk     string `json:"nomor_induk" binding:"required"`
        AsalInstitusi  string `json:"asal_institusi" binding:"required"`
        Jurusan        string `json:"jurusan" binding:"required"`
        KelasSemester  string `json:"kelas_semester" binding:"required"`
}

type PengajuanMagang struct {
        ID             uuid.UUID       `json:"id" db:"id"`
        UserID         *uuid.UUID      `json:"user_id" db:"user_id"`
        NamaLengkap    string          `json:"nama_lengkap" db:"nama_lengkap"`
        TempatLahir    string          `json:"tempat_lahir" db:"tempat_lahir"`
        TanggalLahir   time.Time       `json:"tanggal_lahir" db:"tanggal_lahir"`
        JenisKelamin   string          `json:"jenis_kelamin" db:"jenis_kelamin"`
        Alamat         string          `json:"alamat" db:"alamat"`
        NoHP           string          `json:"no_hp" db:"no_hp"`
        Email          string          `json:"email" db:"email"`
        KategoriMagang string          `json:"kategori_magang" db:"kategori_magang"`
        NomorInduk     string          `json:"nomor_induk" db:"nomor_induk"`
        AsalInstitusi  string          `json:"asal_institusi" db:"asal_institusi"`
        Jurusan        string          `json:"jurusan" db:"jurusan"`
        KelasSemester  string          `json:"kelas_semester" db:"kelas_semester"`
        Status         StatusPengajuan `json:"status" db:"status"`
        CatatanHRD     *string         `json:"catatan_hrd" db:"catatan_hrd"`
        VerifiedBy     *uuid.UUID      `json:"verified_by" db:"verified_by"`
        VerifiedAt     *time.Time      `json:"verified_at" db:"verified_at"`
        AkunTerkirimAt *time.Time      `json:"akun_terkirim_at" db:"akun_terkirim_at"`
        CreatedAt      time.Time       `json:"created_at" db:"created_at"`
        UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

type PengajuanDetail struct {
        PengajuanMagang
        NamaPeriode *string `json:"nama_periode"`
        NamaDivisi  *string `json:"nama_divisi"`
}

type PengajuanPublikRequest struct {
        NamaLengkap        string `json:"nama_lengkap" binding:"required"`
        Email              string `json:"email" binding:"required,email"`
        NoTelp             string `json:"no_telp" binding:"required"`
        Institusi          string `json:"institusi" binding:"required"`
        Jurusan            string `json:"jurusan" binding:"required"`
        Kategori           string `json:"kategori" binding:"required"`
        JenisKelamin       string `json:"jenis_kelamin" binding:"required"`
        TanggalLahir       string `json:"tanggal_lahir" binding:"required"`
        Alamat             string `json:"alamat" binding:"required"`
        TglMulaiDiinginkan string `json:"tgl_mulai_diinginkan" binding:"required"`
        DurasiMinggu       int    `json:"durasi_minggu" binding:"required,min=4"`
        BidangMinat        string `json:"bidang_minat" binding:"required"`
        PeriodeID          string `json:"periode_id"`
}

type UpdateStatusRequest struct {
        Status   StatusPengajuan `json:"status" binding:"required"`
        Catatan  string          `json:"catatan"`
        CatatanHRD *string       `json:"catatan_hrd"`
        DivisiID   *string       `json:"divisi_id"`
}

type NilaiRequest struct {
        Nilai        float64 `json:"nilai" binding:"required"`
        CatatanNilai *string `json:"catatan_nilai"`
}

type SertifikatListItem struct {
        PelaksanaanID         string     `json:"pelaksanaan_id"`
        Status                string     `json:"status"`
        TanggalMulai          time.Time  `json:"tanggal_mulai"`
        TanggalSelesai        time.Time  `json:"tanggal_selesai"`
        Divisi                string     `json:"divisi"`
        Pembimbing            string     `json:"pembimbing"`
        SertifikatGenerated   bool       `json:"sertifikat_generated"`
        SertifikatGeneratedAt *time.Time `json:"sertifikat_generated_at"`
        Nilai                 *float64   `json:"nilai"`
        NamaLengkap           string     `json:"nama_lengkap"`
        Email                 string     `json:"email"`
}

type AjukanPerpanjanganRequest struct {
        DurasiHari int    `json:"durasi_hari" binding:"required"`
        Alasan     string `json:"alasan" binding:"required"`
}

type ProsesPerpanjanganRequest struct {
        CatatanHRD string `json:"catatan_hrd"`
}

// ============================================================
// DOKUMEN
// ============================================================

type Dokumen struct {
        ID          uuid.UUID  `json:"id" db:"id"`
        PengajuanID *uuid.UUID `json:"pengajuan_id" db:"pengajuan_id"`
        UserID      *uuid.UUID `json:"user_id" db:"user_id"`
        Jenis       string     `json:"jenis" db:"jenis"`
        NamaFile    string     `json:"nama_file" db:"nama_file"`
        PathFile    string     `json:"path_file" db:"path_file"`
        UkuranBytes int64      `json:"ukuran_bytes" db:"ukuran_bytes"`
        MimeType    string     `json:"mime_type" db:"mime_type"`
        UploadedAt  time.Time  `json:"uploaded_at" db:"uploaded_at"`
}

// ============================================================
// PELAKSANAAN MAGANG
// ============================================================

type PelaksanaanMagang struct {
        ID                    uuid.UUID         `json:"id" db:"id"`
        PengajuanID           uuid.UUID         `json:"pengajuan_id" db:"pengajuan_id"`
        UserID                uuid.UUID         `json:"user_id" db:"user_id"`
        PeriodeID             *uuid.UUID        `json:"periode_id" db:"periode_id"`
        TanggalMulai          time.Time         `json:"tanggal_mulai" db:"tanggal_mulai"`
        TanggalSelesai        time.Time         `json:"tanggal_selesai" db:"tanggal_selesai"`
        Divisi                string            `json:"divisi" db:"divisi"`
        PembimbingID          *uuid.UUID        `json:"pembimbing_id" db:"pembimbing_id"`
        Status                StatusPelaksanaan `json:"status" db:"status"`
        Nilai                 *float64          `json:"nilai" db:"nilai"`
        CatatanNilai          *string           `json:"catatan_nilai" db:"catatan_nilai"`
        DinilaiOleh           *uuid.UUID        `json:"dinilai_oleh" db:"dinilai_oleh"`
        DinilaiAt             *time.Time        `json:"dinilai_at" db:"dinilai_at"`
        SertifikatGenerated   bool              `json:"sertifikat_generated" db:"sertifikat_generated"`
        SertifikatPath        *string           `json:"sertifikat_path" db:"sertifikat_path"`
        SertifikatGeneratedAt *time.Time        `json:"sertifikat_generated_at" db:"sertifikat_generated_at"`
        CreatedAt             time.Time         `json:"created_at" db:"created_at"`
        UpdatedAt             time.Time         `json:"updated_at" db:"updated_at"`
        // join fields (optional)
        PembimbingNama    *string `json:"pembimbing_nama,omitempty" db:"pembimbing_nama"`
        SudahDiperpanjang bool    `json:"sudah_diperpanjang,omitempty" db:"sudah_diperpanjang"`
        NamaPeserta       string  `json:"nama_peserta,omitempty" db:"nama_peserta"`
}

type PelaksanaanDetail struct {
        PelaksanaanMagang
        NamaPeserta    string   `json:"nama_peserta"`
        EmailPeserta   string   `json:"email_peserta"`
        Institusi      string   `json:"institusi"`
        Jurusan        string   `json:"jurusan"`
        NamaPembimbing *string  `json:"nama_pembimbing"`
        NamaDivisi     *string  `json:"nama_divisi"`
        DivisiID       *uuid.UUID `json:"divisi_id"`
}

type SetJadwalRequest struct {
        TglMulai      string `json:"tgl_mulai"`
        TglSelesai    string `json:"tgl_selesai"`
        TanggalMulai  string `json:"tanggal_mulai" binding:"required"`
        TanggalSelesai string `json:"tanggal_selesai" binding:"required"`
        Divisi        string `json:"divisi"`
        Pembimbing    string `json:"pembimbing"`
}

type TolakIzinSakitRequest struct {
        CatatanHRD string `json:"catatan_hrd"`
}

type ReviewLaporanRequest struct {
        Status     StatusLaporan `json:"status" binding:"required"`
        CatatanHRD string        `json:"catatan_hrd"`
}

// ============================================================
// ABSENSI
// ============================================================

type Absensi struct {
        ID             uuid.UUID  `json:"id" db:"id"`
        PelaksanaanID  uuid.UUID  `json:"pelaksanaan_id" db:"pelaksanaan_id"`
        Tanggal        string     `json:"tanggal" db:"tanggal"`
        JamMasuk       *string    `json:"jam_masuk" db:"jam_masuk"`
        JamKeluar      *string    `json:"jam_keluar" db:"jam_keluar"`
        Keterangan     string     `json:"keterangan" db:"keterangan"`
        Kegiatan       string     `json:"kegiatan" db:"kegiatan"`
        Latitude       *float64   `json:"latitude" db:"latitude"`
        Longitude      *float64   `json:"longitude" db:"longitude"`
        TTDPembimbing  bool       `json:"ttd_pembimbing" db:"ttd_pembimbing"`
        ApprovedBy     *uuid.UUID `json:"approved_by" db:"approved_by"`
        ApprovedAt     *time.Time `json:"approved_at" db:"approved_at"`
        Catatan        *string    `json:"catatan" db:"catatan"`
        CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

type AbsensiConfig struct {
        ID              int       `json:"id" db:"id"`
        JamMasukBuka    string    `json:"jam_masuk_buka" db:"jam_masuk_buka"`
        JamMasukTutup   string    `json:"jam_masuk_tutup" db:"jam_masuk_tutup"`
        JamPulangBuka   string    `json:"jam_pulang_buka" db:"jam_pulang_buka"`
        JamPulangTutup  string    `json:"jam_pulang_tutup" db:"jam_pulang_tutup"`
        UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
        UpdatedBy       *uuid.UUID `json:"updated_by" db:"updated_by"`
}

type AbsensiConfigUpdateRequest struct {
        JamMasukBuka   string `json:"jam_masuk_buka" binding:"required"`
        JamMasukTutup  string `json:"jam_masuk_tutup" binding:"required"`
        JamPulangBuka  string `json:"jam_pulang_buka" binding:"required"`
        JamPulangTutup string `json:"jam_pulang_tutup" binding:"required"`
}

type CheckInRequest struct {
        Lat *float64 `json:"lat"`
        Lng *float64 `json:"lng"`
}

type CheckOutRequest struct {
        Lat *float64 `json:"lat"`
        Lng *float64 `json:"lng"`
}

// ============================================================
// IZIN SAKIT
// ============================================================

type IzinSakitRequest struct {
        ID            uuid.UUID  `json:"id" db:"id"`
        PelaksanaanID uuid.UUID  `json:"pelaksanaan_id" db:"pelaksanaan_id"`
        UserID        uuid.UUID  `json:"user_id" db:"user_id"`
        Tanggal       string     `json:"tanggal" db:"tanggal"`
        Jenis         string     `json:"jenis" db:"jenis"`
        Alasan        string     `json:"alasan" db:"alasan"`
        BuktiPath     *string    `json:"bukti_path" db:"bukti_path"`
        Status        string     `json:"status" db:"status"`
        CatatanHRD    *string    `json:"catatan_hrd" db:"catatan_hrd"`
        DiprosesOleh  *uuid.UUID `json:"diproses_oleh" db:"diproses_oleh"`
        DiprosesAt    *time.Time `json:"diproses_at" db:"diproses_at"`
        CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

type IzinSakitDetail struct {
        IzinSakitRequest
        NamaPeserta string `json:"nama_peserta"`
        Divisi      string `json:"divisi"`
}

type AjukanIzinRequest struct {
        Tanggal string  `json:"tanggal" binding:"required"`
        Jenis   string  `json:"jenis" binding:"required"`
        Alasan  string  `json:"alasan" binding:"required"`
        BuktiPath *string `json:"bukti_path"`
}

// ============================================================
// CHAT & TIKET
// ============================================================

type ChatTiket struct {
        ID         uuid.UUID   `json:"id" db:"id"`
        UserID     uuid.UUID   `json:"user_id" db:"user_id"`
        NomorTiket string      `json:"nomor_tiket" db:"nomor_tiket"`
        Subjek     *string     `json:"subjek" db:"subjek"`
        Status     StatusTiket `json:"status" db:"status"`
        AssignedTo *uuid.UUID  `json:"assigned_to" db:"assigned_to"`
        Kategori   KategoriTiket `json:"kategori" db:"kategori"`
        ExpiresAt  *time.Time  `json:"expires_at" db:"expires_at"`
        HangusAt   *time.Time  `json:"hangus_at" db:"hangus_at"`
        CreatedAt  time.Time   `json:"created_at" db:"created_at"`
        UpdatedAt  time.Time   `json:"updated_at" db:"updated_at"`
}

type ChatTiketDetail struct {
        ID            uuid.UUID     `json:"id"`
        UserID        uuid.UUID     `json:"user_id"`
        NomorTiket    string        `json:"nomor_tiket"`
        Subjek        *string       `json:"subjek"`
        Status        StatusTiket   `json:"status"`
        AssignedTo    *uuid.UUID    `json:"assigned_to"`
        Kategori      KategoriTiket `json:"kategori"`
        UserNama      string        `json:"user_nama"`
        UserEmail     string        `json:"user_email"`
        AssignedNama  *string       `json:"assigned_nama"`
        UnreadCount   int           `json:"unread_count"`
        LastMessage   *string       `json:"last_message"`
        ExpiresAt     *time.Time    `json:"expires_at"`
        HangusAt      *time.Time    `json:"hangus_at"`
        CreatedAt     time.Time     `json:"created_at"`
        UpdatedAt     time.Time     `json:"updated_at"`
}

type ChatPesan struct {
        ID        uuid.UUID `json:"id" db:"id"`
        TiketID   uuid.UUID `json:"tiket_id" db:"tiket_id"`
        SenderID  uuid.UUID `json:"sender_id" db:"sender_id"`
        Pesan     string    `json:"pesan" db:"pesan"`
        IsRead    bool      `json:"is_read" db:"is_read"`
        CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ChatPesanDetail struct {
        ChatPesan
        SenderNama string   `json:"sender_nama"`
        SenderRole UserRole `json:"sender_role"`
}

type ChatBuatTiketRequest struct {
        Subjek   string `json:"subjek" binding:"required"`
        Pesan    string `json:"pesan" binding:"required"`
        Kategori string `json:"kategori"`
}

type ChatKirimPesanRequest struct {
        Pesan string `json:"pesan" binding:"required"`
}

// ============================================================
// CHAT KNOWLEDGE BASE
// ============================================================

type ChatKnowledge struct {
        ID         uuid.UUID `json:"id" db:"id"`
        Pertanyaan string    `json:"pertanyaan" db:"pertanyaan"`
        KataKunci  []string  `json:"kata_kunci" db:"kata_kunci"`
        Jawaban    string    `json:"jawaban" db:"jawaban"`
        Kategori   string    `json:"kategori" db:"kategori"`
        IsActive   bool      `json:"is_active" db:"is_active"`
        CreatedAt  time.Time `json:"created_at" db:"created_at"`
        UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ChatKnowledgeRequest struct {
        Pertanyaan string   `json:"pertanyaan" binding:"required"`
        KataKunci  []string `json:"kata_kunci" binding:"required"`
        Jawaban    string   `json:"jawaban" binding:"required"`
        Kategori   string   `json:"kategori"`
        IsActive   *bool    `json:"is_active"`
}

type NLPQueryRequest struct {
        Pesan string `json:"pesan" binding:"required"`
}

type NLPQueryResponse struct {
        Terjawab bool    `json:"terjawab"`
        Jawaban  string  `json:"jawaban"`
        Skor     float64 `json:"skor"`
        Sumber   string  `json:"sumber"` // "knowledge" atau "faq"
}

// ============================================================
// NOTIFIKASI
// ============================================================

type Notifikasi struct {
        ID          uuid.UUID  `json:"id" db:"id"`
        UserID      uuid.UUID  `json:"user_id" db:"user_id"`
        Judul       string     `json:"judul" db:"judul"`
        Pesan       string     `json:"pesan" db:"pesan"`
        Tipe        string     `json:"tipe" db:"tipe"`
        ReferensiID *uuid.UUID `json:"referensi_id" db:"referensi_id"`
        IsRead      bool       `json:"is_read" db:"is_read"`
        Route       string     `json:"route" db:"route"`
        CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// ============================================================
// LANDING PAGE CMS
// ============================================================

type LandingContent struct {
        ID        uuid.UUID  `json:"id" db:"id"`
        Kunci     string     `json:"kunci" db:"kunci"`
        Nilai     string     `json:"nilai" db:"nilai"`
        Tipe      string     `json:"tipe" db:"tipe"`
        UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
        UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type FAQ struct {
        ID         uuid.UUID `json:"id" db:"id"`
        Pertanyaan string    `json:"pertanyaan" db:"pertanyaan"`
        Jawaban    string    `json:"jawaban" db:"jawaban"`
        Urutan     int       `json:"urutan" db:"urutan"`
        IsActive   bool      `json:"is_active" db:"is_active"`
        CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type PeriodeMagang struct {
        ID          uuid.UUID `json:"id" db:"id"`
        Nama        string    `json:"nama" db:"nama"`
        TanggalBuka time.Time `json:"tanggal_buka" db:"tanggal_buka"`
        TanggalTutup time.Time `json:"tanggal_tutup" db:"tanggal_tutup"`
        Kuota       int       `json:"kuota" db:"kuota"`
        IsActive    bool      `json:"is_active" db:"is_active"`
        CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type LandingAlur struct {
        ID        uuid.UUID `json:"id" db:"id"`
        Judul     string    `json:"judul" db:"judul"`
        Paragraf  string    `json:"paragraf" db:"paragraf"`
        GambarURL string    `json:"gambar_url" db:"gambar_url"`
        Urutan    int       `json:"urutan" db:"urutan"`
        CreatedAt time.Time `json:"created_at" db:"created_at"`
        UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================
// DIVISI
// ============================================================

type Divisi struct {
        ID        uuid.UUID  `json:"id" db:"id"`
        Nama      string     `json:"nama" db:"nama"`
        IsActive  bool       `json:"is_active" db:"is_active"`
        Urutan    int        `json:"urutan" db:"urutan"`
        GeoLat    *float64   `json:"geo_lat" db:"geo_lat"`
        GeoLng    *float64   `json:"geo_lng" db:"geo_lng"`
        GeoRadius *int       `json:"geo_radius" db:"geo_radius"`
        NamaLokasi *string   `json:"nama_lokasi" db:"nama_lokasi"`
        CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// ============================================================
// LAPORAN MAGANG
// ============================================================

type LaporanMagang struct {
        ID            uuid.UUID      `json:"id" db:"id"`
        PelaksanaanID uuid.UUID      `json:"pelaksanaan_id" db:"pelaksanaan_id"`
        Versi         int            `json:"versi" db:"versi"`
        NamaFile      string         `json:"nama_file" db:"nama_file"`
        PathFile      string         `json:"path_file" db:"path_file"`
        UkuranBytes   *int64         `json:"ukuran_bytes" db:"ukuran_bytes"`
        MimeType      *string        `json:"mime_type" db:"mime_type"`
        Status        StatusLaporan  `json:"status" db:"status"`
        CatatanHRD    *string        `json:"catatan_hrd" db:"catatan_hrd"`
        DireviewOleh  *uuid.UUID     `json:"direview_oleh" db:"direview_oleh"`
        DireviewAt    *time.Time     `json:"direview_at" db:"direview_at"`
        DiuploadAt    time.Time      `json:"diupload_at" db:"diupload_at"`
}

type LaporanDetail struct {
        LaporanMagang
        NamaPeserta string `json:"nama_peserta"`
        Email       string `json:"email"`
}

// ============================================================
// PENILAIAN
// ============================================================

type Penilaian struct {
        ID                  uuid.UUID  `json:"id" db:"id"`
        PelaksanaanID       uuid.UUID  `json:"pelaksanaan_id" db:"pelaksanaan_id"`
        NilaiMotivasi       *float64   `json:"nilai_motivasi" db:"nilai_motivasi"`
        NilaiInisiatif      *float64   `json:"nilai_inisiatif" db:"nilai_inisiatif"`
        NilaiDisiplinWaktu  *float64   `json:"nilai_disiplin_waktu" db:"nilai_disiplin_waktu"`
        NilaiKerajinan      *float64   `json:"nilai_kerajinan" db:"nilai_kerajinan"`
        NilaiKreativitas    *float64   `json:"nilai_kreativitas" db:"nilai_kreativitas"`
        NilaiTanggungJawab  *float64   `json:"nilai_tanggung_jawab" db:"nilai_tanggung_jawab"`
        NilaiKerjasama      *float64   `json:"nilai_kerjasama" db:"nilai_kerjasama"`
        NilaiAdaptasi       *float64   `json:"nilai_adaptasi" db:"nilai_adaptasi"`
        NilaiKehadiran      *float64   `json:"nilai_kehadiran" db:"nilai_kehadiran"`
        NilaiK3Safety       *float64   `json:"nilai_k3_safety" db:"nilai_k3_safety"`
        NilaiK3Metode       *float64   `json:"nilai_k3_metode" db:"nilai_k3_metode"`
        NilaiK3Manajemen    *float64   `json:"nilai_k3_manajemen" db:"nilai_k3_manajemen"`
        NilaiK3Volume       *float64   `json:"nilai_k3_volume" db:"nilai_k3_volume"`
        NilaiPrsProses      *float64   `json:"nilai_prs_proses" db:"nilai_prs_proses"`
        NilaiPrsTeori       *float64   `json:"nilai_prs_teori" db:"nilai_prs_teori"`
        NilaiPrsJudul       *float64   `json:"nilai_prs_judul" db:"nilai_prs_judul"`
        NilaiPrsData        *float64   `json:"nilai_prs_data" db:"nilai_prs_data"`
        ManagerNama         *string    `json:"manager_nama" db:"manager_nama"`
        Catatan             *string    `json:"catatan" db:"catatan"`
        NilaiAkhir          *float64   `json:"nilai_akhir" db:"nilai_akhir"`
        DinilaiOleh         *uuid.UUID `json:"dinilai_oleh" db:"dinilai_oleh"`
        DinilaiAt           *time.Time `json:"dinilai_at" db:"dinilai_at"`
        CreatedAt           time.Time  `json:"created_at" db:"created_at"`
        UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
        // join field
        Kejuruan []PenilaianKejuruan `json:"kejuruan,omitempty"`
}

type PenilaianKejuruan struct {
        ID          uuid.UUID `json:"id" db:"id"`
        PenilaianID uuid.UUID `json:"penilaian_id" db:"penilaian_id"`
        Komponen    string    `json:"komponen" db:"komponen"`
        Nilai       float64   `json:"nilai" db:"nilai"`
        Urutan      int       `json:"urutan" db:"urutan"`
}

type PenilaianUpsertRequest struct {
        PelaksanaanID      uuid.UUID `json:"pelaksanaan_id"`
        NilaiMotivasi      float64   `json:"nilai_motivasi"`
        NilaiInisiatif     float64   `json:"nilai_inisiatif"`
        NilaiDisiplinWaktu float64   `json:"nilai_disiplin_waktu"`
        NilaiKerajinan     float64   `json:"nilai_kerajinan"`
        NilaiKreativitas   float64   `json:"nilai_kreativitas"`
        NilaiTanggungJawab float64   `json:"nilai_tanggung_jawab"`
        NilaiKerjasama     float64   `json:"nilai_kerjasama"`
        NilaiAdaptasi      float64   `json:"nilai_adaptasi"`
        NilaiKehadiran     float64   `json:"nilai_kehadiran"`
        NilaiK3Safety      float64   `json:"nilai_k3_safety"`
        NilaiK3Metode      float64   `json:"nilai_k3_metode"`
        NilaiK3Manajemen   float64   `json:"nilai_k3_manajemen"`
        NilaiK3Volume      float64   `json:"nilai_k3_volume"`
        NilaiPrsProses     float64   `json:"nilai_prs_proses"`
        NilaiPrsTeori      float64   `json:"nilai_prs_teori"`
        NilaiPrsJudul      float64   `json:"nilai_prs_judul"`
        NilaiPrsData       float64   `json:"nilai_prs_data"`
        ManagerNama        string    `json:"manager_nama"`
        Catatan            string    `json:"catatan"`
        NilaiAkhir         float64   `json:"nilai_akhir"`
        Kejuruan           []struct {
                Komponen string  `json:"komponen"`
                Nilai    float64 `json:"nilai"`
                Urutan   int     `json:"urutan"`
        } `json:"kejuruan"`
}

type PenilaianResponse struct {
        Penilaian   *Penilaian          `json:"penilaian"`
        Kejuruan    []PenilaianKejuruan `json:"kejuruan"`
        NamaHRD     string              `json:"nama_hrd"`
        NamaPeserta string              `json:"nama_peserta"`
}

// ============================================================
// PERPANJANGAN MAGANG
// ============================================================

type PerpanjanganMagang struct {
        ID                 uuid.UUID  `json:"id" db:"id"`
        PelaksanaanID      uuid.UUID  `json:"pelaksanaan_id" db:"pelaksanaan_id"`
        DiajukanOleh       uuid.UUID  `json:"diajukan_oleh" db:"diajukan_oleh"`
        RolePengaju        string     `json:"role_pengaju" db:"role_pengaju"`
        DurasiHari         int        `json:"durasi_hari" db:"durasi_hari"`
        TanggalSelesaiLama time.Time  `json:"tanggal_selesai_lama" db:"tanggal_selesai_lama"`
        TanggalSelesaiBaru time.Time  `json:"tanggal_selesai_baru" db:"tanggal_selesai_baru"`
        Alasan             string     `json:"alasan" db:"alasan"`
        SuratPath          *string    `json:"surat_path" db:"surat_path"`
        Status             string     `json:"status" db:"status"`
        CatatanHRD         *string    `json:"catatan_hrd" db:"catatan_hrd"`
        DiprosesOleh       *uuid.UUID `json:"diproses_oleh" db:"diproses_oleh"`
        DiprosesAt         *time.Time `json:"diproses_at" db:"diproses_at"`
        CreatedAt          time.Time  `json:"created_at" db:"created_at"`
        UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

type PerpanjanganDetail struct {
        PerpanjanganMagang
        NamaPeserta    string    `json:"nama_peserta"`
        Email          string    `json:"email"`
        TglSelesaiLama time.Time `json:"tgl_selesai_lama"`
}

type PerpanjanganAjukanRequest struct {
        TglSelesaiBaru string `json:"tgl_selesai_baru" binding:"required"`
        Alasan         string `json:"alasan" binding:"required"`
}

type PerpanjanganHRDRequest struct {
        PelaksanaanID  string `json:"pelaksanaan_id" binding:"required"`
        TglSelesaiBaru string `json:"tgl_selesai_baru" binding:"required"`
        Alasan         string `json:"alasan" binding:"required"`
}

// ============================================================
// WEBSOCKET PAYLOADS
// ============================================================

type TiketUpdateWsPayload struct {
        TiketID      string  `json:"tiket_id"`
        Status       string  `json:"status"`
        AssignedTo   *string `json:"assigned_to"`
        AssignedNama *string `json:"assigned_nama"`
}

type ChatMessageWsPayload struct {
        TiketID    string `json:"tiket_id"`
        SenderID   string `json:"sender_id"`
        SenderNama string `json:"sender_nama"`
        SenderRole string `json:"sender_role"`
        Pesan      string `json:"pesan"`
        CreatedAt  string `json:"created_at"`
}

// ============================================================
// ERROR & SUCCESS RESPONSE
// ============================================================

type ErrorResponse struct {
        Error   string `json:"error"`
        Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
        Message string      `json:"message,omitempty"`
        Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
        Data       interface{} `json:"data"`
        Total      int         `json:"total"`
        Page       int         `json:"page"`
        Limit      int         `json:"limit"`
        TotalPages int         `json:"total_pages,omitempty"`
}

// ============================================================
// ABSENSI REQUESTS
// ============================================================

type AbsensiCheckInRequest struct {
        Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
        Accuracy  float64 `json:"accuracy"`
}

type AbsensiCheckOutRequest struct {
        Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
        Accuracy  float64 `json:"accuracy"`
        Kegiatan  string  `json:"kegiatan"`
}

// ============================================================
// REKAP ABSENSI
// ============================================================

type RekapAbsensiRow struct {
        PelaksanaanID  uuid.UUID `json:"pelaksanaan_id" db:"pelaksanaan_id"`
        NamaLengkap    string    `json:"nama_lengkap" db:"nama_lengkap"`
        AsalInstitusi  string    `json:"asal_institusi" db:"asal_institusi"`
        KategoriMagang string    `json:"kategori_magang" db:"kategori_magang"`
        Divisi         string    `json:"divisi" db:"divisi"`
        TanggalMulai   time.Time `json:"tanggal_mulai" db:"tanggal_mulai"`
        TanggalSelesai time.Time `json:"tanggal_selesai" db:"tanggal_selesai"`
        Status         string    `json:"status" db:"status"`
        Hadir          int       `json:"hadir" db:"hadir"`
        Izin           int       `json:"izin" db:"izin"`
        Sakit          int       `json:"sakit" db:"sakit"`
        Alpha          int       `json:"alpha" db:"alpha"`
        PendingApproval int      `json:"pending_approval" db:"pending_approval"`
}

// ============================================================
// DOKUMEN WITH USER
// ============================================================

type DokumenWithUser struct {
        ID           uuid.UUID  `json:"id" db:"id"`
        PengajuanID  uuid.UUID  `json:"pengajuan_id" db:"pengajuan_id"`
        UserID       *uuid.UUID `json:"user_id" db:"user_id"`
        Jenis        string     `json:"jenis" db:"jenis"`
        NamaFile     string     `json:"nama_file" db:"nama_file"`
        PathFile     string     `json:"path_file" db:"path_file"`
        UkuranBytes  int64      `json:"ukuran_bytes" db:"ukuran_bytes"`
        MimeType     string     `json:"mime_type" db:"mime_type"`
        UploadedAt   time.Time  `json:"uploaded_at" db:"uploaded_at"`
        NamaPemilik  string     `json:"nama_pemilik" db:"nama_pemilik"`
        EmailPemilik string     `json:"email_pemilik" db:"email_pemilik"`
}

// ============================================================
// IZIN SAKIT REQUEST DETAIL
// ============================================================

type IzinSakitRequestDetail struct {
        ID             uuid.UUID  `json:"id" db:"id"`
        PelaksanaanID  uuid.UUID  `json:"pelaksanaan_id" db:"pelaksanaan_id"`
        UserID         uuid.UUID  `json:"user_id" db:"user_id"`
        Tanggal        time.Time  `json:"tanggal" db:"tanggal"`
        Jenis          string     `json:"jenis" db:"jenis"`
        Alasan         string     `json:"alasan" db:"alasan"`
        BuktiPath      *string    `json:"bukti_path" db:"bukti_path"`
        Status         string     `json:"status" db:"status"`
        CatatanHRD     *string    `json:"catatan_hrd" db:"catatan_hrd"`
        DiprosesOleh   *uuid.UUID `json:"diproses_oleh" db:"diproses_oleh"`
        DiprosesAt     *time.Time `json:"diproses_at" db:"diproses_at"`
        CreatedAt      time.Time  `json:"created_at" db:"created_at"`
        NamaPeserta    string     `json:"nama_peserta" db:"nama_peserta"`
        Divisi         string     `json:"divisi" db:"divisi"`
}

// ============================================================
// ALUR ITEM (Landing Page)
// ============================================================

type AlurItem struct {
        ID        uuid.UUID `json:"id" db:"id"`
        Judul     string    `json:"judul" db:"judul"`
        Paragraf  string    `json:"paragraf" db:"paragraf"`
        GambarURL string    `json:"gambar_url" db:"gambar_url"`
        Urutan    int       `json:"urutan" db:"urutan"`
        CreatedAt time.Time `json:"created_at" db:"created_at"`
        UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================
// LAPORAN MAGANG
// ============================================================

type StatusLaporan string

const (
        StatusLaporanMenunggu  StatusLaporan = "menunggu_review"
        StatusLaporanDisetujui StatusLaporan = "disetujui"
        StatusLaporanDitolak   StatusLaporan = "ditolak"
        StatusLaporanRevisi    StatusLaporan = "revisi"
)

type LaporanMagangDetail struct {
        ID             uuid.UUID      `json:"id" db:"id"`
        PelaksanaanID  uuid.UUID      `json:"pelaksanaan_id" db:"pelaksanaan_id"`
        Versi          int            `json:"versi" db:"versi"`
        NamaFile       string         `json:"nama_file" db:"nama_file"`
        PathFile       string         `json:"path_file" db:"path_file"`
        UkuranBytes    int64          `json:"ukuran_bytes" db:"ukuran_bytes"`
        MimeType       string         `json:"mime_type" db:"mime_type"`
        Status         StatusLaporan  `json:"status" db:"status"`
        CatatanHRD     *string        `json:"catatan_hrd" db:"catatan_hrd"`
        DireviewOleh   *uuid.UUID     `json:"direview_oleh" db:"direview_oleh"`
        DireviewAt     *time.Time     `json:"direview_at" db:"direview_at"`
        DiuploadAt     time.Time      `json:"diupload_at" db:"diupload_at"`
        NamaPeserta    string         `json:"nama_peserta" db:"nama_peserta"`
        AsalInstitusi  string         `json:"asal_institusi" db:"asal_institusi"`
        Divisi         string         `json:"divisi" db:"divisi"`
}

// ============================================================
// STATUS HISTORY PENGAJUAN
// ============================================================

type StatusHistory struct {
        ID          uuid.UUID  `json:"id" db:"id"`
        PengajuanID uuid.UUID  `json:"pengajuan_id" db:"pengajuan_id"`
        StatusLama  string     `json:"status_lama" db:"status_lama"`
        StatusBaru  string     `json:"status_baru" db:"status_baru"`
        ChangedBy   *uuid.UUID `json:"changed_by" db:"changed_by"`
        Catatan     *string    `json:"catatan" db:"catatan"`
        CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// ============================================================
// PENILAIAN LIST ITEM
// ============================================================

type PenilaianListItem struct {
        PelaksanaanID uuid.UUID  `json:"pelaksanaan_id"`
        Status        string     `json:"status"`
        NamaLengkap   string     `json:"nama_lengkap"`
        Email         string     `json:"email"`
        Divisi        string     `json:"divisi"`
        Pembimbing    string     `json:"pembimbing"`
        TanggalMulai  time.Time  `json:"tanggal_mulai"`
        TanggalSelesai time.Time `json:"tanggal_selesai"`
        PenilaianID   *uuid.UUID `json:"penilaian_id"`
        NilaiAkhir    *float64   `json:"nilai_akhir"`
        DinilaiAt     *time.Time `json:"dinilai_at"`
}

// ============================================================
// PERPANJANGAN MAGANG DETAIL
// ============================================================

type PerpanjanganMagangDetail struct {
        ID                 uuid.UUID  `json:"id"`
        PelaksanaanID      uuid.UUID  `json:"pelaksanaan_id"`
        DiajukanOleh       uuid.UUID  `json:"diajukan_oleh"`
        RolePengaju        string     `json:"role_pengaju"`
        DurasiHari         int        `json:"durasi_hari"`
        TanggalSelesaiLama time.Time  `json:"tanggal_selesai_lama"`
        TanggalSelesaiBaru time.Time  `json:"tanggal_selesai_baru"`
        Alasan             string     `json:"alasan"`
        SuratPath          *string    `json:"surat_path"`
        Status             string     `json:"status"`
        CatatanHRD         *string    `json:"catatan_hrd"`
        DiprosesOleh       *uuid.UUID `json:"diproses_oleh"`
        DiprosesAt         *time.Time `json:"diproses_at"`
        CreatedAt          time.Time  `json:"created_at"`
        UpdatedAt          time.Time  `json:"updated_at"`
        NamaPeserta        string     `json:"nama_peserta"`
        AsalInstitusi      string     `json:"asal_institusi"`
        Divisi             string     `json:"divisi"`
        NamaPengaju        string     `json:"nama_pengaju"`
}

// ============================================================
// DASHBOARD ADMIN STATS
// ============================================================

type DashboardAdminStats struct {
        TotalPeserta      int `json:"total_peserta"`
        TotalPengajuan    int `json:"total_pengajuan"`
        PengajuanMenunggu int `json:"pengajuan_menunggu"`
        PengajuanDiterima int `json:"pengajuan_diterima"`
        PengajuanDitolak  int `json:"pengajuan_ditolak"`
        MagangAktif       int `json:"magang_aktif"`
        MagangSelesai     int `json:"magang_selesai"`
        TiketChatMenunggu int `json:"tiket_chat_menunggu"`
}
