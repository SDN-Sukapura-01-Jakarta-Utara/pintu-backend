-- Migration: create_kutipan_kepsek_table
-- Created: 2026-02-11 10:30:03

BEGIN;

CREATE TABLE kutipan_kepsek (
    id SERIAL PRIMARY KEY,
    nama_kepsek VARCHAR(255) NOT NULL,
    foto_kepsek VARCHAR(512) NOT NULL,
    kutipan_kepsek TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER
);

COMMIT;