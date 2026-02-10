-- Migration: create_ekstrakurikulers_table
-- Created: 2026-02-10 06:59:23
-- Description: Create ekstrakurikulers table with jsonb for multiple kelas

BEGIN;

CREATE TABLE ekstrakurikulers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    kelas_ids JSONB NOT NULL DEFAULT '[]',
    kategori VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_ekstrakurikulers_name ON ekstrakurikulers(name);
CREATE INDEX idx_ekstrakurikulers_kategori ON ekstrakurikulers(kategori);
CREATE INDEX idx_ekstrakurikulers_kelas_ids ON ekstrakurikulers USING GIN(kelas_ids);

COMMIT;
