-- Migration: create_pertanyaan_table
-- Created: 2026-04-28 08:18:24

BEGIN;

CREATE TABLE IF NOT EXISTS pertanyaan (
    id SERIAL PRIMARY KEY,
    id_tiket VARCHAR(50) UNIQUE NOT NULL,
    tanggal_pengajuan TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    nama VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telepon VARCHAR(20),
    kategori VARCHAR(100) NOT NULL,
    prioritas VARCHAR(50) DEFAULT 'Sedang',
    judul VARCHAR(255) NOT NULL,
    deskripsi TEXT NOT NULL,
    file_pertanyaan JSONB DEFAULT '[]'::jsonb,
    judul_jawaban VARCHAR(255),
    deskripsi_jawaban TEXT,
    file_jawaban JSONB DEFAULT '[]'::jsonb,
    tanggal_proses TIMESTAMP,
    email_terkirim BOOLEAN DEFAULT FALSE,
    tanggal_selesai TIMESTAMP,
    status VARCHAR(50) DEFAULT 'pending',
    replied_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted_by_id INTEGER REFERENCES users(id)
);

-- Create index for faster queries
CREATE INDEX idx_pertanyaan_id_tiket ON pertanyaan(id_tiket);
CREATE INDEX idx_pertanyaan_status ON pertanyaan(status);

COMMIT;
