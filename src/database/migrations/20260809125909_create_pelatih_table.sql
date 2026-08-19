-- Migration: create_pelatih_table
-- Created: 2026-08-12 12:59:09
-- Description: Create pelatih table for ekstrakurikuler coach/instructor data management

BEGIN;

CREATE TABLE pelatih (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    username VARCHAR(100) UNIQUE,
    password VARCHAR(255),
    telepon VARCHAR(20),
    alamat TEXT,
    foto_profil VARCHAR(512),
    keahlian TEXT,
    sertifikat JSON,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    
    FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_pelatih_nama ON pelatih(nama);
CREATE INDEX idx_pelatih_username ON pelatih(username);
CREATE INDEX idx_pelatih_status ON pelatih(status);
CREATE INDEX idx_pelatih_deleted ON pelatih(deleted_at);

COMMIT;
