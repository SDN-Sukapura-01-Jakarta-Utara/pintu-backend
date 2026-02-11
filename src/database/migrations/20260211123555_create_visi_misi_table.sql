-- Migration: create_visi_misi_table
-- Created: 2026-02-11 12:35:55

BEGIN;

CREATE TABLE visi_misi (
    id SERIAL PRIMARY KEY,
    visi TEXT,
    misi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER
);

COMMIT;