-- Migration: create_bidang_studis_table
-- Created: 2026-02-09 20:25:38
-- Description: Create bidang_studis table for field of study management

BEGIN;

CREATE TABLE bidang_studis (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_bidang_studis_name ON bidang_studis(name);
CREATE INDEX idx_bidang_studis_status ON bidang_studis(status);

COMMIT;
