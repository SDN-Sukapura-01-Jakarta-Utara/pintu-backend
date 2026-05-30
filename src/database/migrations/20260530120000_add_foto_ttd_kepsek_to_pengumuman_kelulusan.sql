-- Migration: add_foto_ttd_kepsek_to_pengumuman_kelulusan
-- Created: 2026-05-30 12:00:00

BEGIN;

ALTER TABLE pengumuman_kelulusan
ADD COLUMN IF NOT EXISTS foto_kepsek VARCHAR(500),
ADD COLUMN IF NOT EXISTS ttd_kepsek VARCHAR(500),
ADD COLUMN IF NOT EXISTS nama_kepsek VARCHAR(255);

COMMIT;
