-- ============================================================
-- Migration: Tabel penilaian & penilaian_kejuruan
-- Jalankan di database emagang_telpp
-- ============================================================

CREATE TABLE IF NOT EXISTS penilaian (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pelaksanaan_id   UUID NOT NULL UNIQUE REFERENCES pelaksanaan_magang(id) ON DELETE CASCADE,
    nilai_disiplin   NUMERIC(5,2),
    nilai_sikap      NUMERIC(5,2),
    nilai_kerjasama  NUMERIC(5,2),
    nilai_inisiatif  NUMERIC(5,2),
    nilai_laporan    NUMERIC(5,2),
    nilai_presentasi NUMERIC(5,2),
    catatan          TEXT,
    nilai_akhir      NUMERIC(5,2),
    dinilai_oleh     UUID REFERENCES users(id),
    dinilai_at       TIMESTAMPTZ,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS penilaian_kejuruan (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    penilaian_id  UUID NOT NULL REFERENCES penilaian(id) ON DELETE CASCADE,
    komponen      VARCHAR(255) NOT NULL,
    nilai         NUMERIC(5,2) NOT NULL,
    urutan        INTEGER DEFAULT 0
);

-- Index untuk performa query
CREATE INDEX IF NOT EXISTS idx_penilaian_pelaksanaan ON penilaian(pelaksanaan_id);
CREATE INDEX IF NOT EXISTS idx_penilaian_kejuruan_penilaian ON penilaian_kejuruan(penilaian_id);
