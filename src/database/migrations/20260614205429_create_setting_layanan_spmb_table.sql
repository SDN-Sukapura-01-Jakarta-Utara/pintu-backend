-- Migration: create_setting_layanan_spmb_table
-- Created: 2026-06-14 20:54:29

BEGIN;

CREATE TABLE IF NOT EXISTS setting_layanan_spmb (
    id SERIAL PRIMARY KEY,
    nama_kepala_sekolah VARCHAR(255),
    nip_kepala_sekolah VARCHAR(50),
    nama_ketua_panitia VARCHAR(255),
    nip_ketua_panitia VARCHAR(50),
    grup_wa TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

COMMIT;
