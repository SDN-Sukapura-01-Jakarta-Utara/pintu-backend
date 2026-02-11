-- Migration: create_announcements_table
-- Created: 2026-02-11 20:20:33

BEGIN;

CREATE TABLE announcements (
    id BIGSERIAL PRIMARY KEY,
    judul VARCHAR(255) NOT NULL,
    tanggal DATE NOT NULL,
    deskripsi TEXT,
    gambar VARCHAR(255),
    files JSONB DEFAULT '[]'::jsonb,
    penulis VARCHAR(255) NOT NULL,
    status_publikasi VARCHAR(50) DEFAULT 'draft',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_announcements_judul ON announcements(judul);

COMMIT;
