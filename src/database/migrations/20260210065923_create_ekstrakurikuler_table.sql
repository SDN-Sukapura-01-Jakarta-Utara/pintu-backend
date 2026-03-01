-- Migration: create_ekstrakurikuler_table
-- Created: 2026-02-10 06:59:23
-- Description: Create ekstrakurikuler table with jsonb for multiple kelas

BEGIN;

CREATE TABLE ekstrakurikuler (
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
CREATE INDEX idx_ekstrakurikuler_name ON ekstrakurikuler(name);
CREATE INDEX idx_ekstrakurikuler_kategori ON ekstrakurikuler(kategori);
CREATE INDEX idx_ekstrakurikuler_kelas_ids ON ekstrakurikuler USING GIN(kelas_ids);

COMMIT;
