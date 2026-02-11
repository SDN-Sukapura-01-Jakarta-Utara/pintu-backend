-- Migration: create_activity_galleries_table
-- Created: 2026-02-11 21:02:54

BEGIN;

CREATE TABLE activity_galleries (
    id BIGSERIAL PRIMARY KEY,
    judul VARCHAR(255) NOT NULL,
    tanggal DATE NOT NULL,
    foto JSONB DEFAULT '[]'::jsonb,
    status_publikasi VARCHAR(50) DEFAULT 'draft',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_activity_galleries_judul ON activity_galleries(tanggal);

COMMIT;
