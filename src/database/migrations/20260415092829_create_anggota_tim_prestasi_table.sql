-- Migration: create_anggota_tim_prestasi_table
-- Created: 2026-04-15 09:28:29
-- Description: Create anggota_tim_prestasi table for team member management

BEGIN;

CREATE TABLE anggota_tim_prestasi (
    id SERIAL PRIMARY KEY,
    prestasi_id INTEGER NOT NULL,
    peserta_didik_id INTEGER NOT NULL,
    tahun_pelajaran_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    CONSTRAINT fk_anggota_tim_prestasi_prestasi FOREIGN KEY (prestasi_id) REFERENCES prestasi(id) ON DELETE SET NULL,
    CONSTRAINT fk_anggota_tim_prestasi_peserta_didik FOREIGN KEY (peserta_didik_id) REFERENCES peserta_didik(id) ON DELETE SET NULL,
    CONSTRAINT fk_anggota_tim_prestasi_tahun_pelajaran FOREIGN KEY (tahun_pelajaran_id) REFERENCES tahun_pelajaran(id) ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_anggota_tim_prestasi_prestasi_id ON anggota_tim_prestasi(prestasi_id);
CREATE INDEX idx_anggota_tim_prestasi_peserta_didik_id ON anggota_tim_prestasi(peserta_didik_id);
CREATE INDEX idx_anggota_tim_prestasi_tahun_pelajaran_id ON anggota_tim_prestasi(tahun_pelajaran_id);

COMMIT;
