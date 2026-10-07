-- Migration: create_anak_inklusi_rombel_table
-- Created: 2026-09-13 13:13:26
-- Description: Create anak_inklusi_rombel table for tracking inclusive students per academic year and class

BEGIN;

CREATE TABLE anak_inklusi_rombel (
    id SERIAL PRIMARY KEY,
    anak_inklusi_id INTEGER NOT NULL,
    peserta_didik_rombel_id INTEGER NOT NULL,
    guru_pendamping_khusus_id INTEGER,
    guru_kelas_id INTEGER,
    catatan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_anak_inklusi_rombel_anak_inklusi 
        FOREIGN KEY (anak_inklusi_id) 
        REFERENCES anak_inklusi(id) 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_anak_inklusi_rombel_peserta_didik_rombel 
        FOREIGN KEY (peserta_didik_rombel_id) 
        REFERENCES peserta_didik_rombel(id) 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_anak_inklusi_rombel_gpk 
        FOREIGN KEY (guru_pendamping_khusus_id) 
        REFERENCES kepegawaian(id) 
        ON DELETE SET NULL,
        
    CONSTRAINT fk_anak_inklusi_rombel_guru_kelas 
        FOREIGN KEY (guru_kelas_id) 
        REFERENCES kepegawaian(id) 
        ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_anak_inklusi_rombel_anak_inklusi_id ON anak_inklusi_rombel(anak_inklusi_id);
CREATE INDEX idx_anak_inklusi_rombel_peserta_didik_rombel_id ON anak_inklusi_rombel(peserta_didik_rombel_id);
CREATE INDEX idx_anak_inklusi_rombel_gpk_id ON anak_inklusi_rombel(guru_pendamping_khusus_id);
CREATE INDEX idx_anak_inklusi_rombel_guru_kelas_id ON anak_inklusi_rombel(guru_kelas_id);

-- Create unique constraint to prevent duplicate entries per tahun ajaran
CREATE UNIQUE INDEX idx_anak_inklusi_rombel_unique ON anak_inklusi_rombel(anak_inklusi_id, peserta_didik_rombel_id) WHERE deleted_at IS NULL;

COMMIT;
