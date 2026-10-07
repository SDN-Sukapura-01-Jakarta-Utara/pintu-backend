-- Migration: create_anak_inklusi_table
-- Created: 2026-09-13 13:13:14
-- Description: Create anak_inklusi table for inclusive education student data management

BEGIN;

CREATE TABLE anak_inklusi (
    id SERIAL PRIMARY KEY,
    peserta_didik_id INTEGER NOT NULL,
    jenis_hambatan VARCHAR(100) NOT NULL,
    tanggal_identifikasi DATE,
    tanggal_diagnosa DATE,
    tanggal_kadaluarsa_surat DATE,
    file_surat_dokter VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'identified',
    catatan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_anak_inklusi_peserta_didik 
        FOREIGN KEY (peserta_didik_id) 
        REFERENCES peserta_didik(id) 
        ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_anak_inklusi_peserta_didik_id ON anak_inklusi(peserta_didik_id);
CREATE INDEX idx_anak_inklusi_status ON anak_inklusi(status);
CREATE INDEX idx_anak_inklusi_jenis_hambatan ON anak_inklusi(jenis_hambatan);

-- Create unique constraint to prevent duplicate entries
CREATE UNIQUE INDEX idx_anak_inklusi_unique_peserta_didik ON anak_inklusi(peserta_didik_id) WHERE deleted_at IS NULL;

COMMIT;
