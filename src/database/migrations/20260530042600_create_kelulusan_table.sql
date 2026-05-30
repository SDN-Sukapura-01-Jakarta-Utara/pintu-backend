-- Migration: create_kelulusan_table
-- Created: 2026-05-30 04:26:00

BEGIN;

CREATE TABLE kelulusan (
    id SERIAL PRIMARY KEY,
    nomor_peserta VARCHAR(50) NOT NULL UNIQUE,
    nisn VARCHAR(25) NOT NULL,
    nama VARCHAR(255) NOT NULL,
    tanggal_lahir DATE NOT NULL,
    nilai JSONB NOT NULL DEFAULT '{}'::jsonb,
    lulus BOOLEAN NOT NULL DEFAULT FALSE,
    skl VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_kelulusan_nisn ON kelulusan(nisn);
CREATE INDEX idx_kelulusan_tanggal_lahir ON kelulusan(tanggal_lahir);

COMMIT;
