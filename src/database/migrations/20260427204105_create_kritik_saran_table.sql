-- Migration: create_kritik_saran_table
-- Created: 2026-04-27 20:41:05

BEGIN;

CREATE TABLE IF NOT EXISTS kritik_saran (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    kritik_saran TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

COMMIT;
