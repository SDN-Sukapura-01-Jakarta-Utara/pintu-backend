-- Migration: create_absensi_pelatih_ekskul_table
-- Created: 2026-08-12 12:54:55
-- Description: Create absensi_pelatih_ekskul table for coach attendance tracking per activity session

BEGIN;

CREATE TABLE absensi_pelatih_ekskul (
    id SERIAL PRIMARY KEY,
    kegiatan_ekskul_id INTEGER NOT NULL,
    pelatih_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    FOREIGN KEY (kegiatan_ekskul_id) REFERENCES kegiatan_ekskul(id) ON DELETE CASCADE,
    FOREIGN KEY (pelatih_id) REFERENCES pelatih(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_absensi_pelatih_kegiatan ON absensi_pelatih_ekskul(kegiatan_ekskul_id);
CREATE INDEX idx_absensi_pelatih_pelatih ON absensi_pelatih_ekskul(pelatih_id);
CREATE INDEX idx_absensi_pelatih_deleted ON absensi_pelatih_ekskul(deleted_at);

-- Unique constraint: one attendance record per coach per activity session
CREATE UNIQUE INDEX idx_unique_absensi_pelatih_active 
ON absensi_pelatih_ekskul(kegiatan_ekskul_id, pelatih_id) 
WHERE deleted_at IS NULL;

COMMIT;
