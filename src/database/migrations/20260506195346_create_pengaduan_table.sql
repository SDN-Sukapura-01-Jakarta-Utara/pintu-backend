-- Migration: create_pengaduan_table
-- Created: 2026-05-06 19:53:46

BEGIN;

CREATE TABLE IF NOT EXISTS pengaduan (
    id SERIAL PRIMARY KEY,
    id_tiket VARCHAR(50) UNIQUE NOT NULL,
    tanggal_pengajuan TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    tipe_pelapor VARCHAR(50) DEFAULT 'anonim',
    nama VARCHAR(255),
    email VARCHAR(255),
    telepon VARCHAR(20),
    kategori VARCHAR(100) NOT NULL,
    prioritas VARCHAR(50) DEFAULT 'Sedang',
    judul VARCHAR(255) NOT NULL,
    deskripsi TEXT NOT NULL,
    file_pengaduan JSONB DEFAULT '[]'::jsonb,
    judul_jawaban VARCHAR(255),
    deskripsi_jawaban TEXT,
    file_jawaban JSONB DEFAULT '[]'::jsonb,
    tanggal_proses TIMESTAMP,
    email_terkirim BOOLEAN DEFAULT FALSE,
    tindak_lanjut TEXT,
    file_tindak_lanjut JSONB DEFAULT '[]'::jsonb,
    tanggal_selesai TIMESTAMP,
    status VARCHAR(50) DEFAULT 'pending',
    replied_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted_by_id INTEGER REFERENCES users(id)
);

-- Create index for faster queries
CREATE INDEX idx_pengaduan_id_tiket ON pengaduan(id_tiket);
CREATE INDEX idx_pengaduan_status ON pengaduan(status);
CREATE INDEX idx_pengaduan_tipe_pelapor ON pengaduan(tipe_pelapor);

COMMIT;
