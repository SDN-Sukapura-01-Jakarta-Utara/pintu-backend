-- Migration: create_tahun_pelajaran_table
-- Created: 2026-02-08 22:02:39
-- Description: Create tahun_pelajaran table for academic year management

BEGIN;

CREATE TABLE tahun_pelajaran (
    id SERIAL PRIMARY KEY,
    tahun_pelajaran VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_tahun_pelajaran_tahun_pelajaran ON tahun_pelajaran(tahun_pelajaran);
CREATE INDEX idx_tahun_pelajaran_status ON tahun_pelajaran(status);

COMMIT;
