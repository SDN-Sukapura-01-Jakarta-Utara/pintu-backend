-- Migration: create_contact_table
-- Created: 2026-02-12 07:12:08

BEGIN;

CREATE TABLE contacts (
    id BIGSERIAL PRIMARY KEY,
    alamat TEXT NOT NULL,
    telepon VARCHAR(20) NOT NULL,
    email VARCHAR(100) NOT NULL,
    jam_buka JSONB DEFAULT '[]'::jsonb,
    gmaps TEXT,
    website VARCHAR(255),
    youtube VARCHAR(255),
    instagram VARCHAR(255),
    tiktok VARCHAR(255),
    facebook VARCHAR(255),
    twitter VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

COMMIT;
