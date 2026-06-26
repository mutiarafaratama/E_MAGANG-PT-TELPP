-- ============================================================
-- VIEW LENGKAP: Menggabungkan SEMUA tabel database
-- Database: emagang_telpp
-- ============================================================

CREATE OR REPLACE VIEW public.newview AS

SELECT 
    -- USERS
    u.id AS user_id,
    u.nama_lengkap,
    u.email,
    u.role,
    u.is_active,
    u.password_changed,
    u.created_at AS user_created_at,
    u.updated_at AS user_updated_at,
    
    -- REFRESH TOKENS (terakhir)
    rt.id AS refresh_token_id,
    rt.expires_at AS token_expires_at,
    rt.created_at AS token_created_at,
    
    -- PENGAJUAN MAGANG
    pm.id AS pengajuan_id,
    pm.tempat_lahir,
    pm.tanggal_lahir,
    pm.jenis_kelamin,
    pm.alamat,
    pm.no_hp,
    pm.kategori_magang,
    pm.nomor_induk,
    pm.asal_institusi,
    pm.jurusan,
    pm.kelas_semester,
    pm.status AS status_pengajuan,
    pm.catatan_hrd,
    pm.verified_by,
    pm.verified_at,
    pm.akun_terkirim_at,
    pm.created_at AS pengajuan_created_at,
    pm.updated_at AS pengajuan_updated_at,
    
    -- STATUS HISTORY (terakhir)
    sh.id AS status_history_id,
    sh.status_lama,
    sh.status_baru,
    sh.changed_by,
    sh.catatan AS status_catatan,
    sh.created_at AS status_changed_at,
    
    -- PERIODE MAGANG
    per.id AS periode_id,
    per.nama AS periode_nama,
    per.tanggal_buka,
    per.tanggal_tutup,
    per.kuota,
    per.is_active AS periode_active,
    
    -- PELAKSANAAN MAGANG
    pl.id AS pelaksanaan_id,
    pl.tanggal_mulai,
    pl.tanggal_selesai,
    pl.divisi,
    pl.pembimbing_id,
    pl.pembimbing_nama,
    pl.status AS status_pelaksanaan,
    pl.nilai,
    pl.catatan_nilai,
    pl.dinilai_oleh,
    pl.dinilai_at,
    pl.sertifikat_generated,
    pl.sertifikat_path,
    pl.sertifikat_generated_at,
    pl.sudah_diperpanjang,
    pl.nama_peserta,
    pl.created_at AS pelaksanaan_created_at,
    pl.updated_at AS pelaksanaan_updated_at,
    
    -- DIVISI
    d.id AS divisi_id,
    d.nama AS divisi_nama,
    d.is_active AS divisi_active,
    d.urutan AS divisi_urutan,
    d.geo_lat,
    d.geo_lng,
    d.geo_radius,
    d.nama_lokasi,
    
    -- ABSENSI CONFIG
    ac.id AS absensi_config_id,
    ac.jam_masuk_buka,
    ac.jam_masuk_tutup,
    ac.jam_pulang_buka,
    ac.jam_pulang_tutup,
    ac.updated_at AS absensi_config_updated,
    
    -- ABSENSI (summary)
    COUNT(DISTINCT abs.id) AS total_absensi,
    COUNT(DISTINCT CASE WHEN abs.keterangan = 'hadir' THEN abs.id END) AS total_hadir,
    COUNT(DISTINCT CASE WHEN abs.keterangan = 'izin' THEN abs.id END) AS total_izin,
    COUNT(DISTINCT CASE WHEN abs.keterangan = 'sakit' THEN abs.id END) AS total_sakit,
    COUNT(DISTINCT CASE WHEN abs.keterangan = 'alpha' THEN abs.id END) AS total_alpha,
    
    -- IZIN SAKIT REQUEST (terakhir)
    isr.id AS izin_sakit_id,
    isr.tanggal AS izin_tanggal,
    isr.jenis AS izin_jenis,
    isr.alasan AS izin_alasan,
    isr.bukti_path,
    isr.status AS izin_status,
    isr.catatan_hrd AS izin_catatan_hrd,
    isr.diproses_oleh,
    isr.diproses_at,
    
    -- DOKUMEN (summary)
    COUNT(DISTINCT dok.id) AS total_dokumen,
    STRING_AGG(DISTINCT dok.jenis, ', ') AS daftar_dokumen,
    
    -- LAPORAN MAGANG (terakhir)
    lm.id AS laporan_id,
    lm.versi AS laporan_versi,
    lm.nama_file AS laporan_nama,
    lm.path_file AS laporan_path,
    lm.ukuran_bytes,
    lm.mime_type,
    lm.status AS status_laporan,
    lm.catatan_hrd AS laporan_catatan_hrd,
    lm.direview_oleh,
    lm.direview_at,
    lm.diupload_at,
    
    -- PENILAIAN
    p.id AS penilaian_id,
    p.nilai_motivasi,
    p.nilai_inisiatif,
    p.nilai_disiplin_waktu,
    p.nilai_kerajinan,
    p.nilai_kreativitas,
    p.nilai_tanggung_jawab,
    p.nilai_kerjasama,
    p.nilai_adaptasi,
    p.nilai_kehadiran,
    p.nilai_k3_safety,
    p.nilai_k3_metode,
    p.nilai_k3_manajemen,
    p.nilai_k3_volume,
    p.nilai_prs_proses,
    p.nilai_prs_teori,
    p.nilai_prs_judul,
    p.nilai_prs_data,
    p.manager_nama,
    p.catatan AS penilaian_catatan,
    p.nilai_akhir,
    p.dinilai_oleh AS penilaian_dinilai_oleh,
    p.dinilai_at AS penilaian_dinilai_at,
    
    -- PENILAIAN KEJURUAN (summary)
    COUNT(DISTINCT pk.id) AS total_kejuruan,
    STRING_AGG(DISTINCT pk.komponen || ': ' || pk.nilai, '; ') AS daftar_kejuruan,
    
    -- PERPANJANGAN MAGANG (terakhir)
    ppm.id AS perpanjangan_id,
    ppm.durasi_hari,
    ppm.tanggal_selesai_lama,
    ppm.tanggal_selesai_baru,
    ppm.alasan AS perpanjangan_alasan,
    ppm.surat_path,
    ppm.status AS status_perpanjangan,
    ppm.catatan_hrd AS perpanjangan_catatan_hrd,
    ppm.diproses_oleh AS perpanjangan_diproses_oleh,
    ppm.diproses_at AS perpanjangan_diproses_at,
    
    -- CHAT TIKET (summary)
    COUNT(DISTINCT ct.id) AS total_tiket,
    COUNT(DISTINCT CASE WHEN ct.status = 'menunggu' THEN ct.id END) AS tiket_menunggu,
    COUNT(DISTINCT CASE WHEN ct.status = 'diproses' THEN ct.id END) AS tiket_diproses,
    COUNT(DISTINCT CASE WHEN ct.status = 'selesai' THEN ct.id END) AS tiket_selesai,
    
    -- CHAT PESAN (summary)
    COUNT(DISTINCT cp.id) AS total_pesan,
    COUNT(DISTINCT CASE WHEN cp.is_read = false THEN cp.id END) AS pesan_unread,
    
    -- NOTIFIKASI (summary)
    COUNT(DISTINCT n.id) AS total_notifikasi,
    COUNT(DISTINCT CASE WHEN n.is_read = false THEN n.id END) AS notifikasi_unread,
    
    -- LANDING CONTENT
    lc.id AS landing_content_id,
    lc.kunci,
    lc.nilai AS landing_nilai,
    lc.tipe AS landing_tipe,
    
    -- FAQ
    f.id AS faq_id,
    f.pertanyaan,
    f.jawaban,
    f.urutan AS faq_urutan,
    f.is_active AS faq_active,
    
    -- ALUR ITEM
    ai.id AS alur_id,
    ai.judul AS alur_judul,
    ai.paragraf AS alur_paragraf,
    ai.gambar_url,
    ai.urutan AS alur_urutan

FROM users u
-- LEFT JOIN dengan refresh_tokens
LEFT JOIN refresh_tokens rt ON u.id = rt.user_id
-- LEFT JOIN dengan pengajuan_magang
LEFT JOIN pengajuan_magang pm ON u.id = pm.user_id
-- LEFT JOIN dengan status_history (ambil yang terakhir per pengajuan)
LEFT JOIN LATERAL (
    SELECT DISTINCT ON (pengajuan_id) *
    FROM status_history
    WHERE pengajuan_id = pm.id
    ORDER BY created_at DESC
) sh ON pm.id = sh.pengajuan_id
-- LEFT JOIN dengan pelaksanaan_magang
LEFT JOIN pelaksanaan_magang pl ON pm.id = pl.pengajuan_id
-- LEFT JOIN dengan periode_magang
LEFT JOIN periode_magang per ON pl.periode_id = per.id
-- LEFT JOIN dengan divisi
LEFT JOIN divisi d ON pl.divisi = d.nama
-- LEFT JOIN dengan absensi_config (ambil yang terakhir)
LEFT JOIN LATERAL (
    SELECT DISTINCT ON (id) *
    FROM absensi_config
    ORDER BY updated_at DESC
    LIMIT 1
) ac ON true
-- LEFT JOIN dengan absensi
LEFT JOIN absensi abs ON pl.id = abs.pelaksanaan_id
-- LEFT JOIN dengan izin_sakit_requests (ambil yang terakhir per pelaksanaan)
LEFT JOIN LATERAL (
    SELECT DISTINCT ON (pelaksanaan_id) *
    FROM izin_sakit_requests
    WHERE pelaksanaan_id = pl.id
    ORDER BY created_at DESC
) isr ON pl.id = isr.pelaksanaan_id
-- LEFT JOIN dengan dokumen
LEFT JOIN dokumen dok ON pm.id = dok.pengajuan_id OR u.id = dok.user_id
-- LEFT JOIN dengan laporan_magang (ambil yang terakhir per pelaksanaan)
LEFT JOIN LATERAL (
    SELECT DISTINCT ON (pelaksanaan_id) *
    FROM laporan_magang
    WHERE pelaksanaan_id = pl.id
    ORDER BY versi DESC
) lm ON pl.id = lm.pelaksanaan_id
-- LEFT JOIN dengan penilaian
LEFT JOIN penilaian p ON pl.id = p.pelaksanaan_id
-- LEFT JOIN dengan penilaian_kejuruan
LEFT JOIN penilaian_kejuruan pk ON p.id = pk.penilaian_id
-- LEFT JOIN dengan perpanjangan_magang (ambil yang terakhir per pelaksanaan)
LEFT JOIN LATERAL (
    SELECT DISTINCT ON (pelaksanaan_id) *
    FROM perpanjangan_magang
    WHERE pelaksanaan_id = pl.id
    ORDER BY created_at DESC
) ppm ON pl.id = ppm.pelaksanaan_id
-- LEFT JOIN dengan chat_tiket
LEFT JOIN chat_tiket ct ON u.id = ct.user_id
-- LEFT JOIN dengan chat_pesan
LEFT JOIN chat_pesan cp ON ct.id = cp.tiket_id
-- LEFT JOIN dengan notifikasi
LEFT JOIN notifikasi n ON u.id = n.user_id
-- LEFT JOIN dengan landing_content
LEFT JOIN landing_content lc ON true
-- LEFT JOIN dengan faq
LEFT JOIN faq f ON true
-- LEFT JOIN dengan alur_item
LEFT JOIN alur_item ai ON true

GROUP BY
    u.id, rt.id, pm.id, sh.id, per.id, pl.id, d.id, ac.id, isr.id, lm.id, 
    p.id, ppm.id, lc.id, f.id, ai.id;

