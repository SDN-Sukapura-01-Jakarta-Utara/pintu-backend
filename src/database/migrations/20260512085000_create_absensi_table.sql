-- Migration: create_absensi_table
-- Created: 2026-05-12 08:50:00
-- Description: Create absensi table for student attendance tracking from barcode scanner machine

BEGIN;

CREATE TABLE absensi (
    id SERIAL PRIMARY KEY,
    peserta_didik_id INTEGER NOT NULL REFERENCES peserta_didik(id) ON DELETE CASCADE,
    tanggal DATE NOT NULL,
    jam_datang TIME,
    jam_pulang TIME,
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_absensi_harian UNIQUE(peserta_didik_id, tanggal)
);

CREATE INDEX idx_absensi_peserta_didik_id ON absensi(peserta_didik_id);
CREATE INDEX idx_absensi_tanggal ON absensi(tanggal);
CREATE INDEX idx_absensi_peserta_tanggal ON absensi(peserta_didik_id, tanggal);
CREATE INDEX idx_absensi_status ON absensi(status);

COMMIT;
