-- Migration: create_konfigurasi_mutasi_siswa_table
-- Created: 2026-07-08 19:22:32

BEGIN;

CREATE TABLE IF NOT EXISTS konfigurasi_mutasi_siswa (
    id SERIAL PRIMARY KEY,
    tanggal_buka_pendaftaran DATE NOT NULL,
    tanggal_tutup_pendaftaran DATE NOT NULL,
    nama_kepala_sekolah VARCHAR(255) NOT NULL,
    nip_kepala_sekolah VARCHAR(50) NOT NULL,
    nama_ketua_panitia VARCHAR(255) NOT NULL,
    nip_ketua_panitia VARCHAR(50) NOT NULL,
    template_sptjm VARCHAR(255),
    grup_wa TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

COMMIT;
