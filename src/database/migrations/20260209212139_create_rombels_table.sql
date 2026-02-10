-- Migration: create_rombels_table
-- Created: 2026-02-09 21:21:39
-- Description: Create rombels table for study group management with foreign key to kelas

BEGIN;

CREATE TABLE rombels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    kelas_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_rombels_kelas FOREIGN KEY (kelas_id) REFERENCES kelas(id) ON DELETE RESTRICT
);

-- Create indexes for better query performance
CREATE INDEX idx_rombels_name ON rombels(name);
CREATE INDEX idx_rombels_kelas_id ON rombels(kelas_id);

COMMIT;
