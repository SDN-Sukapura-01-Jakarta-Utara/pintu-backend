-- Migration: create_peserta_didik_ekstrakurikuler_table
-- Created: 2026-08-02 00:00:01
-- Description: Create junction table for student extracurricular registration

BEGIN;

CREATE TABLE peserta_didik_ekstrakurikuler (
    id SERIAL PRIMARY KEY,
    peserta_didik_rombel_id INTEGER NOT NULL,
    ekstrakurikuler_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    
    FOREIGN KEY (peserta_didik_rombel_id) REFERENCES peserta_didik_rombel(id) ON DELETE CASCADE,
    FOREIGN KEY (ekstrakurikuler_id) REFERENCES ekstrakurikuler(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_pd_ekskul_rombel ON peserta_didik_ekstrakurikuler(peserta_didik_rombel_id);
CREATE INDEX idx_pd_ekskul_ekskul ON peserta_didik_ekstrakurikuler(ekstrakurikuler_id);

CREATE UNIQUE INDEX idx_unique_peserta_didik_ekskul_active 
ON peserta_didik_ekstrakurikuler(peserta_didik_rombel_id, ekstrakurikuler_id) 
WHERE deleted_at IS NULL;

COMMIT;
