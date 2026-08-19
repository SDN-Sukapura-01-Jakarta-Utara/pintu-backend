-- Migration: create_absensi_ekskul_table
-- Created: 2026-08-10 13:09:39
-- Description: Create absensi_ekskul table for ekstrakurikuler attendance detail per student

BEGIN;

CREATE TABLE absensi_ekskul (
    id SERIAL PRIMARY KEY,
    kegiatan_ekskul_id INTEGER NOT NULL,
    peserta_didik_rombel_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL,
    keterangan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    FOREIGN KEY (kegiatan_ekskul_id) REFERENCES kegiatan_ekskul(id) ON DELETE CASCADE,
    FOREIGN KEY (peserta_didik_rombel_id) REFERENCES peserta_didik_rombel(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_absensi_ekskul_kegiatan ON absensi_ekskul(kegiatan_ekskul_id);
CREATE INDEX idx_absensi_ekskul_peserta ON absensi_ekskul(peserta_didik_rombel_id);
CREATE INDEX idx_absensi_ekskul_status ON absensi_ekskul(status);
CREATE INDEX idx_absensi_ekskul_deleted ON absensi_ekskul(deleted_at);

-- Unique constraint: one attendance record per student per activity session
CREATE UNIQUE INDEX idx_unique_absensi_ekskul_active 
ON absensi_ekskul(kegiatan_ekskul_id, peserta_didik_rombel_id) 
WHERE deleted_at IS NULL;

COMMIT;
