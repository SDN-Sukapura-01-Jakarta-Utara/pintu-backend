-- Migration: create_layanan_spmb_table
-- Created: 2026-06-14 19:52:49

BEGIN;

CREATE TABLE IF NOT EXISTS layanan_spmb (
    id SERIAL PRIMARY KEY,
    nama_orang_tua VARCHAR(255) NOT NULL,
    nomor_telepon VARCHAR(20) NOT NULL,
    alamat TEXT NOT NULL,
    nama_lengkap_murid VARCHAR(255) NOT NULL,
    keperluan TEXT NOT NULL,
    tanggal_laporan TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create index for faster queries
CREATE INDEX idx_layanan_spmb_status ON layanan_spmb(status);
CREATE INDEX idx_layanan_spmb_tanggal_laporan ON layanan_spmb(tanggal_laporan);

COMMIT;
