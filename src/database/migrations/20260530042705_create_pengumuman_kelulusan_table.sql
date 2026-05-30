-- Migration: create_pengumuman_kelulusan_table
-- Created: 2026-05-30 04:27:05

BEGIN;

CREATE TABLE pengumuman_kelulusan (
    id SERIAL PRIMARY KEY,
    sambutan_kelulusan TEXT NOT NULL,
    tanggal_pengumuman_nilai TIMESTAMP NOT NULL,
    tanggal_pengumuman_kelulusan TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

COMMIT;
