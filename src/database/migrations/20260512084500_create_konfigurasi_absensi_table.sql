-- Migration: create_konfigurasi_absensi_table
-- Created: 2026-07-06 21:35:44
-- Description: Create konfigurasi_absensi table for attendance time configuration and validation

BEGIN;

CREATE TABLE konfigurasi_absensi (
    id SERIAL PRIMARY KEY,

    -- Rentang waktu absen datang
    jam_datang_mulai TIME NOT NULL,      -- Contoh: 05:30 (mulai bisa scan)
    jam_max_datang TIME NOT NULL,        -- Contoh: 06:30 (batas tepat waktu, lewat dari ini = terlambat)
    jam_datang_selesai TIME NOT NULL,    -- Contoh: 08:00 (batas akhir scan datang)
    
    -- Rentang waktu absen pulang
    jam_pulang_mulai TIME NOT NULL,      -- Contoh: 14:00
    jam_pulang_selesai TIME NOT NULL,    -- Contoh: 16:00

    -- Data Kepala Sekolah untuk PDF export
    nama_kepsek VARCHAR(200),
    nip_kepsek VARCHAR(50),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
