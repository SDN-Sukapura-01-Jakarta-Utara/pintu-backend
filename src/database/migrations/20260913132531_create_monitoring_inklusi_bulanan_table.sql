-- Migration: create_monitoring_inklusi_bulanan_table
-- Created: 2026-09-13 13:25:31
-- Description: Create monitoring_inklusi_bulanan table for monthly progress monitoring of inclusive students

BEGIN;

CREATE TABLE monitoring_inklusi_bulanan (
    id SERIAL PRIMARY KEY,
    ppi_inklusi_id INTEGER NOT NULL,
    bulan SMALLINT NOT NULL,
    tahun SMALLINT NOT NULL,
    deskripsi_perkembangan TEXT NOT NULL,
    kendala_ditemui TEXT,
    tindak_lanjut TEXT,
    file_pendukung JSONB DEFAULT '[]'::jsonb,
    guru_pengisi_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_monitoring_inklusi_bulanan_ppi 
        FOREIGN KEY (ppi_inklusi_id) 
        REFERENCES ppi_inklusi(id) 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_monitoring_inklusi_bulanan_guru_pengisi 
        FOREIGN KEY (guru_pengisi_id) 
        REFERENCES kepegawaian(id) 
        ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_monitoring_inklusi_bulanan_ppi_id ON monitoring_inklusi_bulanan(ppi_inklusi_id);
CREATE INDEX idx_monitoring_inklusi_bulanan_bulan_tahun ON monitoring_inklusi_bulanan(bulan, tahun);
CREATE INDEX idx_monitoring_inklusi_bulanan_guru_pengisi_id ON monitoring_inklusi_bulanan(guru_pengisi_id);

COMMIT;
