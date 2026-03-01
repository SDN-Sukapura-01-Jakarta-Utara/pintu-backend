-- Migration: create_rombel_table
-- Created: 2026-02-09 21:21:39
-- Description: Create rombel table for study group management with foreign key to kelas

BEGIN;

CREATE TABLE rombel (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    kelas_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_rombel_kelas FOREIGN KEY (kelas_id) REFERENCES kelas(id) ON DELETE RESTRICT
);

-- Create indexes for better query performance
CREATE INDEX idx_rombel_name ON rombel(name);
CREATE INDEX idx_rombel_kelas_id ON rombel(kelas_id);

COMMIT;
