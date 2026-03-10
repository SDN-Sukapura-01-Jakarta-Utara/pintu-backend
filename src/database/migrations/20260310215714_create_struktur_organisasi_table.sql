-- Migration: create_struktur_organisasi_table
-- Created: 2026-03-10 21:57:14
-- Description: Create struktur_organisasi table for organizational structure management

BEGIN;

CREATE TABLE struktur_organisasi (
    id SERIAL PRIMARY KEY,
    pegawai_id INTEGER,
    nama_non_pegawai VARCHAR(255),
    jabatan_non_pegawai VARCHAR(255),
    urutan INTEGER NOT NULL,
    relasi VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_struktur_organisasi_pegawai FOREIGN KEY (pegawai_id) REFERENCES kepegawaian(id) ON DELETE SET NULL
);

COMMIT;
