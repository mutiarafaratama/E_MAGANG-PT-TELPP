-- Migration untuk menambahkan kolom FCM token ke tabel users
-- Jalankan query ini di database PostgreSQL

-- Cek apakah kolom fcm_token sudah ada, jika belum tambahkan
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name='users' 
        AND column_name='fcm_token'
    ) THEN
        ALTER TABLE users ADD COLUMN fcm_token TEXT;
        RAISE NOTICE 'Kolom fcm_token berhasil ditambahkan';
    ELSE
        RAISE NOTICE 'Kolom fcm_token sudah ada';
    END IF;
END $$;

-- Cek apakah kolom fcm_updated_at sudah ada, jika belum tambahkan
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name='users' 
        AND column_name='fcm_updated_at'
    ) THEN
        ALTER TABLE users ADD COLUMN fcm_updated_at TIMESTAMP;
        RAISE NOTICE 'Kolom fcm_updated_at berhasil ditambahkan';
    ELSE
        RAISE NOTICE 'Kolom fcm_updated_at sudah ada';
    END IF;
END $$;

-- Buat index untuk fcm_token untuk optimasi query
CREATE INDEX IF NOT EXISTS idx_users_fcm_token ON users(fcm_token) WHERE fcm_token IS NOT NULL;
