-- Migration: create_peserta_didik_rombel_table
-- Created: 2026-06-21 01:00:00
-- Description: Create peserta_didik_rombel table for mapping students to rombel per academic year (many-to-many relationship)

BEGIN;

CREATE TABLE peserta_didik_rombel (
    id SERIAL PRIMARY KEY,
    peserta_didik_id INTEGER NOT NULL,
    rombel_id INTEGER NOT NULL,
    tahun_pelajaran_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_peserta_didik_rombel_peserta FOREIGN KEY (peserta_didik_id) REFERENCES peserta_didik(id) ON DELETE CASCADE,
    CONSTRAINT fk_peserta_didik_rombel_rombel FOREIGN KEY (rombel_id) REFERENCES rombel(id) ON DELETE CASCADE,
    CONSTRAINT fk_peserta_didik_rombel_tahun FOREIGN KEY (tahun_pelajaran_id) REFERENCES tahun_pelajaran(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_peserta_didik_rombel_peserta_id ON peserta_didik_rombel(peserta_didik_id);
CREATE INDEX idx_peserta_didik_rombel_rombel_id ON peserta_didik_rombel(rombel_id);
CREATE INDEX idx_peserta_didik_rombel_tahun_id ON peserta_didik_rombel(tahun_pelajaran_id);
CREATE INDEX idx_peserta_didik_rombel_status ON peserta_didik_rombel(status);

COMMIT;
