--
-- PostgreSQL database dump
--

\restrict HphZfkXZymre9dXP9L5Dag0RDKOwqhBcTatrJRqKNb0D8bMTib31VuZaOYShTyO

-- Dumped from database version 18.4
-- Dumped by pg_dump version 18.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: jenis_dokumen; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.jenis_dokumen AS ENUM (
    'proposal_magang',
    'ktp',
    'ktm',
    'pasfoto',
    'bpjs_kis',
    'surat_balasan',
    'laporan_magang',
    'sertifikat',
    'surat_pengantar'
);


--
-- Name: jenis_kelamin; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.jenis_kelamin AS ENUM (
    'laki_laki',
    'perempuan'
);


--
-- Name: kategori_magang; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.kategori_magang AS ENUM (
    'smk',
    'd3_s1_s2',
    'penelitian',
    'd3',
    's1',
    's2',
    'lainnya'
);


--
-- Name: kategori_tiket; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.kategori_tiket AS ENUM (
    'umum',
    'absensi',
    'dokumen',
    'sertifikat'
);


--
-- Name: status_pelaksanaan; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.status_pelaksanaan AS ENUM (
    'menunggu_mulai',
    'aktif',
    'upload_laporan',
    'penilaian',
    'selesai'
);


--
-- Name: status_pengajuan; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.status_pengajuan AS ENUM (
    'diajukan',
    'menunggu_verifikasi',
    'diproses',
    'diterima',
    'ditolak'
);


--
-- Name: status_tiket; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.status_tiket AS ENUM (
    'menunggu',
    'diproses',
    'selesai',
    'hangus'
);


--
-- Name: user_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.user_role AS ENUM (
    'admin',
    'hrd',
    'peserta'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: absensi; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.absensi (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pelaksanaan_id uuid NOT NULL,
    tanggal date NOT NULL,
    jam_masuk time without time zone,
    jam_keluar time without time zone,
    keterangan character varying(20) DEFAULT 'hadir'::character varying,
    kegiatan text,
    ttd_pembimbing boolean DEFAULT false,
    approved_by uuid,
    approved_at timestamp with time zone,
    catatan text,
    created_at timestamp with time zone DEFAULT now(),
    latitude numeric(10,8),
    longitude numeric(11,8),
    lat_keluar numeric(10,7),
    lng_keluar numeric(10,7),
    is_manual boolean DEFAULT false NOT NULL,
    diinput_oleh uuid,
    catatan_manual text
);


--
-- Name: absensi_config; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.absensi_config (
    id integer NOT NULL,
    jam_masuk_buka time without time zone DEFAULT '07:30:00'::time without time zone NOT NULL,
    jam_masuk_tutup time without time zone DEFAULT '08:00:00'::time without time zone NOT NULL,
    jam_pulang_buka time without time zone DEFAULT '15:00:00'::time without time zone NOT NULL,
    jam_pulang_tutup time without time zone DEFAULT '16:00:00'::time without time zone NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by uuid,
    chat_timer_hangus_jam integer DEFAULT 24 NOT NULL,
    chat_timer_idle_jam integer DEFAULT 48 NOT NULL
);


--
-- Name: absensi_config_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.absensi_config_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: absensi_config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.absensi_config_id_seq OWNED BY public.absensi_config.id;


--
-- Name: alur_item; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.alur_item (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    judul character varying(255) NOT NULL,
    paragraf text NOT NULL,
    gambar_url character varying(512),
    urutan integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--
-- Name: TABLE alur_item; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.alur_item IS 'Menyimpan item alur pendaftaran magang untuk landing page';


--
-- Name: chat_knowledge; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_knowledge (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pertanyaan text NOT NULL,
    kata_kunci text[] DEFAULT '{}'::text[] NOT NULL,
    jawaban text NOT NULL,
    kategori character varying(50) DEFAULT 'umum'::character varying NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: chat_pesan; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_pesan (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    tiket_id uuid NOT NULL,
    sender_id uuid NOT NULL,
    pesan text NOT NULL,
    is_read boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: chat_tiket; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_tiket (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    nomor_tiket character varying(20) NOT NULL,
    subjek character varying(255),
    status public.status_tiket DEFAULT 'menunggu'::public.status_tiket,
    assigned_to uuid,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    kategori public.kategori_tiket DEFAULT 'umum'::public.kategori_tiket NOT NULL,
    expires_at timestamp with time zone,
    hangus_at timestamp with time zone
);


--
-- Name: divisi; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.divisi (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    nama character varying(255) NOT NULL,
    is_active boolean DEFAULT true,
    urutan integer DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    geo_lat numeric(10,8),
    geo_lng numeric(11,8),
    geo_radius integer,
    nama_lokasi character varying(255)
);


--
-- Name: dokumen; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.dokumen (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pengajuan_id uuid,
    user_id uuid,
    jenis public.jenis_dokumen NOT NULL,
    nama_file character varying(255) NOT NULL,
    path_file character varying(500) NOT NULL,
    ukuran_bytes bigint,
    mime_type character varying(100),
    uploaded_at timestamp with time zone DEFAULT now()
);


--
-- Name: faq; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.faq (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pertanyaan text NOT NULL,
    jawaban text NOT NULL,
    urutan integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: izin_sakit_request; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.izin_sakit_request (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pelaksanaan_id uuid NOT NULL,
    user_id uuid NOT NULL,
    tanggal date NOT NULL,
    jenis character varying(10) NOT NULL,
    alasan text NOT NULL,
    bukti_path character varying(500),
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    catatan_hrd text,
    diproses_oleh uuid,
    diproses_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT izin_sakit_request_jenis_check CHECK (((jenis)::text = ANY ((ARRAY['izin'::character varying, 'sakit'::character varying])::text[]))),
    CONSTRAINT izin_sakit_request_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'disetujui'::character varying, 'ditolak'::character varying])::text[])))
);


--
-- Name: landing_alur; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.landing_alur (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    judul text DEFAULT ''::text NOT NULL,
    paragraf text DEFAULT ''::text NOT NULL,
    gambar_url text DEFAULT ''::text NOT NULL,
    urutan integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: landing_content; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.landing_content (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    kunci character varying(100) NOT NULL,
    nilai text,
    tipe character varying(20) DEFAULT 'text'::character varying,
    updated_by uuid,
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: laporan_magang; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.laporan_magang (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pelaksanaan_id uuid NOT NULL,
    versi integer DEFAULT 1 NOT NULL,
    nama_file text NOT NULL,
    path_file text NOT NULL,
    ukuran_bytes bigint,
    mime_type text,
    status character varying(30) DEFAULT 'menunggu_review'::character varying NOT NULL,
    catatan_hrd text,
    direview_oleh uuid,
    direview_at timestamp with time zone,
    diupload_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: pelaksanaan_magang; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pelaksanaan_magang (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pengajuan_id uuid NOT NULL,
    user_id uuid NOT NULL,
    periode_id uuid,
    tanggal_mulai date NOT NULL,
    tanggal_selesai date NOT NULL,
    divisi character varying(255),
    pembimbing_id uuid,
    status public.status_pelaksanaan DEFAULT 'menunggu_mulai'::public.status_pelaksanaan,
    nilai numeric(5,2),
    catatan_nilai text,
    dinilai_oleh uuid,
    dinilai_at timestamp with time zone,
    sertifikat_generated boolean DEFAULT false,
    sertifikat_path character varying(500),
    sertifikat_generated_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    pembimbing text,
    sudah_diperpanjang boolean DEFAULT false NOT NULL,
    wa_pembimbing text
);


--
-- Name: newview; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.newview AS
 SELECT a.id,
    a.pelaksanaan_id,
    a.tanggal,
    a.jam_masuk,
    a.jam_keluar,
    a.keterangan,
    a.kegiatan,
    a.ttd_pembimbing,
    a.approved_by,
    a.approved_at,
    a.catatan,
    a.created_at,
    a.latitude,
    a.longitude,
    a.lat_keluar,
    a.lng_keluar,
    pm.pembimbing
   FROM (public.absensi a
     LEFT JOIN public.pelaksanaan_magang pm ON ((pm.id = a.pelaksanaan_id)));


--
-- Name: notifikasi; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notifikasi (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    judul character varying(255) NOT NULL,
    pesan text NOT NULL,
    tipe character varying(50),
    referensi_id uuid,
    is_read boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: password_reset_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.password_reset_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token_hash character varying(64) NOT NULL,
    expired_at timestamp with time zone NOT NULL,
    used boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: pengajuan_magang; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pengajuan_magang (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    nama_lengkap character varying(255) NOT NULL,
    tempat_lahir character varying(100) NOT NULL,
    tanggal_lahir date NOT NULL,
    jenis_kelamin public.jenis_kelamin NOT NULL,
    alamat text NOT NULL,
    no_hp character varying(20) NOT NULL,
    email character varying(255) NOT NULL,
    kategori_magang public.kategori_magang NOT NULL,
    nomor_induk character varying(50) NOT NULL,
    asal_institusi character varying(255) NOT NULL,
    jurusan character varying(255) NOT NULL,
    kelas_semester character varying(50) NOT NULL,
    status public.status_pengajuan DEFAULT 'diajukan'::public.status_pengajuan NOT NULL,
    catatan_hrd text,
    verified_by uuid,
    verified_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    akun_terkirim_at timestamp with time zone
);


--
-- Name: penilaian; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.penilaian (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pelaksanaan_id uuid NOT NULL,
    catatan text,
    nilai_akhir numeric(5,2),
    dinilai_oleh uuid,
    dinilai_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    nilai_motivasi numeric(5,2),
    nilai_inisiatif numeric(5,2),
    nilai_disiplin_waktu numeric(5,2),
    nilai_kerajinan numeric(5,2),
    nilai_kreativitas numeric(5,2),
    nilai_tanggung_jawab numeric(5,2),
    nilai_kerjasama numeric(5,2),
    nilai_adaptasi numeric(5,2),
    nilai_kehadiran numeric(5,2),
    nilai_k3_safety numeric(5,2),
    nilai_k3_metode numeric(5,2),
    nilai_k3_manajemen numeric(5,2),
    nilai_k3_volume numeric(5,2),
    nilai_prs_proses numeric(5,2),
    nilai_prs_teori numeric(5,2),
    nilai_prs_judul numeric(5,2),
    nilai_prs_data numeric(5,2),
    manager_nama text
);


--
-- Name: penilaian_kejuruan; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.penilaian_kejuruan (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    penilaian_id uuid NOT NULL,
    komponen text NOT NULL,
    nilai numeric(5,2) NOT NULL,
    urutan integer DEFAULT 0
);


--
-- Name: periode_magang; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.periode_magang (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    nama character varying(255) NOT NULL,
    tanggal_buka date NOT NULL,
    tanggal_tutup date NOT NULL,
    kuota integer,
    is_active boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: perpanjangan_magang; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.perpanjangan_magang (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pelaksanaan_id uuid NOT NULL,
    diajukan_oleh uuid NOT NULL,
    role_pengaju character varying(10) NOT NULL,
    durasi_hari integer NOT NULL,
    tanggal_selesai_lama date NOT NULL,
    tanggal_selesai_baru date NOT NULL,
    alasan text NOT NULL,
    surat_path character varying(500),
    status character varying(20) DEFAULT 'menunggu'::character varying NOT NULL,
    catatan_hrd text,
    diproses_oleh uuid,
    diproses_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT perpanjangan_magang_durasi_hari_check CHECK ((durasi_hari = ANY (ARRAY[7, 14, 30, 60, 90]))),
    CONSTRAINT perpanjangan_magang_role_pengaju_check CHECK (((role_pengaju)::text = ANY ((ARRAY['peserta'::character varying, 'hrd'::character varying])::text[]))),
    CONSTRAINT perpanjangan_magang_status_check CHECK (((status)::text = ANY ((ARRAY['menunggu'::character varying, 'disetujui'::character varying, 'ditolak'::character varying])::text[])))
);


--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.refresh_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token_hash character varying(255) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: status_history; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.status_history (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pengajuan_id uuid NOT NULL,
    status_lama public.status_pengajuan,
    status_baru public.status_pengajuan NOT NULL,
    changed_by uuid,
    catatan text,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    nama_lengkap character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    role public.user_role DEFAULT 'peserta'::public.user_role NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    password_changed boolean DEFAULT false NOT NULL,
    fcm_token text,
    fcm_updated_at timestamp with time zone
);


--
-- Name: absensi_config id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.absensi_config ALTER COLUMN id SET DEFAULT nextval('public.absensi_config_id_seq'::regclass);


--
-- Name: absensi_config absensi_config_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.absensi_config
    ADD CONSTRAINT absensi_config_pkey PRIMARY KEY (id);


--
-- Name: absensi absensi_pelaksanaan_id_tanggal_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.absensi
    ADD CONSTRAINT absensi_pelaksanaan_id_tanggal_key UNIQUE (pelaksanaan_id, tanggal);


--
-- Name: absensi absensi_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.absensi
    ADD CONSTRAINT absensi_pkey PRIMARY KEY (id);


--
-- Name: alur_item alur_item_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.alur_item
    ADD CONSTRAINT alur_item_pkey PRIMARY KEY (id);


--
-- Name: chat_knowledge chat_knowledge_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_knowledge
    ADD CONSTRAINT chat_knowledge_pkey PRIMARY KEY (id);


--
-- Name: chat_pesan chat_pesan_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_pesan
    ADD CONSTRAINT chat_pesan_pkey PRIMARY KEY (id);


--
-- Name: chat_tiket chat_tiket_nomor_tiket_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_tiket
    ADD CONSTRAINT chat_tiket_nomor_tiket_key UNIQUE (nomor_tiket);


--
-- Name: chat_tiket chat_tiket_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_tiket
    ADD CONSTRAINT chat_tiket_pkey PRIMARY KEY (id);


--
-- Name: divisi divisi_nama_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.divisi
    ADD CONSTRAINT divisi_nama_key UNIQUE (nama);


--
-- Name: divisi divisi_nama_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.divisi
    ADD CONSTRAINT divisi_nama_unique UNIQUE (nama);


--
-- Name: divisi divisi_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.divisi
    ADD CONSTRAINT divisi_pkey PRIMARY KEY (id);


--
-- Name: dokumen dokumen_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.dokumen
    ADD CONSTRAINT dokumen_pkey PRIMARY KEY (id);


--
-- Name: faq faq_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.faq
    ADD CONSTRAINT faq_pkey PRIMARY KEY (id);


--
-- Name: izin_sakit_request izin_sakit_request_pelaksanaan_id_tanggal_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.izin_sakit_request
    ADD CONSTRAINT izin_sakit_request_pelaksanaan_id_tanggal_key UNIQUE (pelaksanaan_id, tanggal);


--
-- Name: izin_sakit_request izin_sakit_request_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.izin_sakit_request
    ADD CONSTRAINT izin_sakit_request_pkey PRIMARY KEY (id);


--
-- Name: landing_alur landing_alur_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.landing_alur
    ADD CONSTRAINT landing_alur_pkey PRIMARY KEY (id);


--
-- Name: landing_content landing_content_kunci_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.landing_content
    ADD CONSTRAINT landing_content_kunci_key UNIQUE (kunci);


--
-- Name: landing_content landing_content_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.landing_content
    ADD CONSTRAINT landing_content_pkey PRIMARY KEY (id);


--
-- Name: laporan_magang laporan_magang_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.laporan_magang
    ADD CONSTRAINT laporan_magang_pkey PRIMARY KEY (id);


--
-- Name: notifikasi notifikasi_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifikasi
    ADD CONSTRAINT notifikasi_pkey PRIMARY KEY (id);


--
-- Name: password_reset_tokens password_reset_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_pkey PRIMARY KEY (id);


--
-- Name: password_reset_tokens password_reset_tokens_token_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_token_hash_key UNIQUE (token_hash);


--
-- Name: pelaksanaan_magang pelaksanaan_magang_pengajuan_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pelaksanaan_magang
    ADD CONSTRAINT pelaksanaan_magang_pengajuan_id_key UNIQUE (pengajuan_id);


--
-- Name: pelaksanaan_magang pelaksanaan_magang_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pelaksanaan_magang
    ADD CONSTRAINT pelaksanaan_magang_pkey PRIMARY KEY (id);


--
-- Name: pengajuan_magang pengajuan_magang_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengajuan_magang
    ADD CONSTRAINT pengajuan_magang_pkey PRIMARY KEY (id);


--
-- Name: penilaian_kejuruan penilaian_kejuruan_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.penilaian_kejuruan
    ADD CONSTRAINT penilaian_kejuruan_pkey PRIMARY KEY (id);


--
-- Name: penilaian penilaian_pelaksanaan_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.penilaian
    ADD CONSTRAINT penilaian_pelaksanaan_id_key UNIQUE (pelaksanaan_id);


--
-- Name: penilaian penilaian_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.penilaian
    ADD CONSTRAINT penilaian_pkey PRIMARY KEY (id);


--
-- Name: periode_magang periode_magang_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.periode_magang
    ADD CONSTRAINT periode_magang_pkey PRIMARY KEY (id);


--
-- Name: perpanjangan_magang perpanjangan_magang_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.perpanjangan_magang
    ADD CONSTRAINT perpanjangan_magang_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: status_history status_history_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.status_history
    ADD CONSTRAINT status_history_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_alur_item_urutan; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_alur_item_urutan ON public.alur_item USING btree (urutan);


--
-- Name: idx_chat_knowledge_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_knowledge_active ON public.chat_knowledge USING btree (is_active);


--
-- Name: idx_chat_knowledge_kategori; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_knowledge_kategori ON public.chat_knowledge USING btree (kategori);


--
-- Name: idx_laporan_pelaksanaan; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_laporan_pelaksanaan ON public.laporan_magang USING btree (pelaksanaan_id);


--
-- Name: idx_laporan_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_laporan_status ON public.laporan_magang USING btree (status);


--
-- Name: idx_penilaian_kejuruan_penilaian; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_penilaian_kejuruan_penilaian ON public.penilaian_kejuruan USING btree (penilaian_id);


--
-- Name: idx_penilaian_pelaksanaan; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_penilaian_pelaksanaan ON public.penilaian USING btree (pelaksanaan_id);


--
-- Name: idx_perpanjangan_pelaksanaan; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_perpanjangan_pelaksanaan ON public.perpanjangan_magang USING btree (pelaksanaan_id);


--
-- Name: idx_perpanjangan_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_perpanjangan_status ON public.perpanjangan_magang USING btree (status);


--
-- Name: idx_prt_token_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_prt_token_hash ON public.password_reset_tokens USING btree (token_hash);


--
-- Name: idx_prt_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_prt_user_id ON public.password_reset_tokens USING btree (user_id);


--
-- Name: idx_users_fcm_token; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_fcm_token ON public.users USING btree (fcm_token) WHERE (fcm_token IS NOT NULL);


--
-- Name: uidx_chat_tiket_user_aktif; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uidx_chat_tiket_user_aktif ON public.chat_tiket USING btree (user_id) WHERE (status <> 'selesai'::public.status_tiket);


--
-- Name: absensi absensi_approved_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.absensi
    ADD CONSTRAINT absensi_approved_by_fkey FOREIGN KEY (approved_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: absensi_config absensi_config_updated_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.absensi_config
    ADD CONSTRAINT absensi_config_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: absensi absensi_diinput_oleh_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.absensi
    ADD CONSTRAINT absensi_diinput_oleh_fkey FOREIGN KEY (diinput_oleh) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: absensi absensi_pelaksanaan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.absensi
    ADD CONSTRAINT absensi_pelaksanaan_id_fkey FOREIGN KEY (pelaksanaan_id) REFERENCES public.pelaksanaan_magang(id) ON DELETE CASCADE;


--
-- Name: chat_pesan chat_pesan_sender_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_pesan
    ADD CONSTRAINT chat_pesan_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: chat_pesan chat_pesan_tiket_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_pesan
    ADD CONSTRAINT chat_pesan_tiket_id_fkey FOREIGN KEY (tiket_id) REFERENCES public.chat_tiket(id) ON DELETE CASCADE;


--
-- Name: chat_tiket chat_tiket_assigned_to_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_tiket
    ADD CONSTRAINT chat_tiket_assigned_to_fkey FOREIGN KEY (assigned_to) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: chat_tiket chat_tiket_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_tiket
    ADD CONSTRAINT chat_tiket_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: dokumen dokumen_pengajuan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.dokumen
    ADD CONSTRAINT dokumen_pengajuan_id_fkey FOREIGN KEY (pengajuan_id) REFERENCES public.pengajuan_magang(id) ON DELETE CASCADE;


--
-- Name: dokumen dokumen_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.dokumen
    ADD CONSTRAINT dokumen_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: izin_sakit_request izin_sakit_request_diproses_oleh_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.izin_sakit_request
    ADD CONSTRAINT izin_sakit_request_diproses_oleh_fkey FOREIGN KEY (diproses_oleh) REFERENCES public.users(id);


--
-- Name: izin_sakit_request izin_sakit_request_pelaksanaan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.izin_sakit_request
    ADD CONSTRAINT izin_sakit_request_pelaksanaan_id_fkey FOREIGN KEY (pelaksanaan_id) REFERENCES public.pelaksanaan_magang(id) ON DELETE CASCADE;


--
-- Name: izin_sakit_request izin_sakit_request_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.izin_sakit_request
    ADD CONSTRAINT izin_sakit_request_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: landing_content landing_content_updated_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.landing_content
    ADD CONSTRAINT landing_content_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.users(id);


--
-- Name: laporan_magang laporan_magang_direview_oleh_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.laporan_magang
    ADD CONSTRAINT laporan_magang_direview_oleh_fkey FOREIGN KEY (direview_oleh) REFERENCES public.users(id);


--
-- Name: laporan_magang laporan_magang_pelaksanaan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.laporan_magang
    ADD CONSTRAINT laporan_magang_pelaksanaan_id_fkey FOREIGN KEY (pelaksanaan_id) REFERENCES public.pelaksanaan_magang(id) ON DELETE CASCADE;


--
-- Name: notifikasi notifikasi_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifikasi
    ADD CONSTRAINT notifikasi_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: password_reset_tokens password_reset_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: pelaksanaan_magang pelaksanaan_magang_dinilai_oleh_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pelaksanaan_magang
    ADD CONSTRAINT pelaksanaan_magang_dinilai_oleh_fkey FOREIGN KEY (dinilai_oleh) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: pelaksanaan_magang pelaksanaan_magang_pembimbing_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pelaksanaan_magang
    ADD CONSTRAINT pelaksanaan_magang_pembimbing_id_fkey FOREIGN KEY (pembimbing_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: pelaksanaan_magang pelaksanaan_magang_pengajuan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pelaksanaan_magang
    ADD CONSTRAINT pelaksanaan_magang_pengajuan_id_fkey FOREIGN KEY (pengajuan_id) REFERENCES public.pengajuan_magang(id);


--
-- Name: pelaksanaan_magang pelaksanaan_magang_periode_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pelaksanaan_magang
    ADD CONSTRAINT pelaksanaan_magang_periode_id_fkey FOREIGN KEY (periode_id) REFERENCES public.periode_magang(id);


--
-- Name: pelaksanaan_magang pelaksanaan_magang_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pelaksanaan_magang
    ADD CONSTRAINT pelaksanaan_magang_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: pengajuan_magang pengajuan_magang_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengajuan_magang
    ADD CONSTRAINT pengajuan_magang_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: pengajuan_magang pengajuan_magang_verified_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengajuan_magang
    ADD CONSTRAINT pengajuan_magang_verified_by_fkey FOREIGN KEY (verified_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: penilaian penilaian_dinilai_oleh_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.penilaian
    ADD CONSTRAINT penilaian_dinilai_oleh_fkey FOREIGN KEY (dinilai_oleh) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: penilaian_kejuruan penilaian_kejuruan_penilaian_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.penilaian_kejuruan
    ADD CONSTRAINT penilaian_kejuruan_penilaian_id_fkey FOREIGN KEY (penilaian_id) REFERENCES public.penilaian(id) ON DELETE CASCADE;


--
-- Name: penilaian penilaian_pelaksanaan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.penilaian
    ADD CONSTRAINT penilaian_pelaksanaan_id_fkey FOREIGN KEY (pelaksanaan_id) REFERENCES public.pelaksanaan_magang(id) ON DELETE CASCADE;


--
-- Name: perpanjangan_magang perpanjangan_magang_diajukan_oleh_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.perpanjangan_magang
    ADD CONSTRAINT perpanjangan_magang_diajukan_oleh_fkey FOREIGN KEY (diajukan_oleh) REFERENCES public.users(id);


--
-- Name: perpanjangan_magang perpanjangan_magang_diproses_oleh_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.perpanjangan_magang
    ADD CONSTRAINT perpanjangan_magang_diproses_oleh_fkey FOREIGN KEY (diproses_oleh) REFERENCES public.users(id);


--
-- Name: perpanjangan_magang perpanjangan_magang_pelaksanaan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.perpanjangan_magang
    ADD CONSTRAINT perpanjangan_magang_pelaksanaan_id_fkey FOREIGN KEY (pelaksanaan_id) REFERENCES public.pelaksanaan_magang(id) ON DELETE CASCADE;


--
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: status_history status_history_changed_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.status_history
    ADD CONSTRAINT status_history_changed_by_fkey FOREIGN KEY (changed_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: status_history status_history_pengajuan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.status_history
    ADD CONSTRAINT status_history_pengajuan_id_fkey FOREIGN KEY (pengajuan_id) REFERENCES public.pengajuan_magang(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict HphZfkXZymre9dXP9L5Dag0RDKOwqhBcTatrJRqKNb0D8bMTib31VuZaOYShTyO

